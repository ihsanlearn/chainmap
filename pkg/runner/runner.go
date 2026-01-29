package runner

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/shlex"
	"github.com/ihsanlearn/chainmap/core"
	"github.com/ihsanlearn/chainmap/logger"
	"github.com/ihsanlearn/chainmap/options"
)

type Runner struct {
	options *options.Options
}

func New(opts *options.Options) *Runner {
	return &Runner{options: opts}
}

func (r *Runner) CheckDependencies() error {
	if _, err := exec.LookPath("nmap"); err != nil {
		return fmt.Errorf("nmap is not installed or not in PATH")
	}
	if _, err := exec.LookPath("xsltproc"); err != nil {
		logger.Error("xsltproc is not installed. HTML report will NOT be generated.")
	}
	return nil
}

func (r *Runner) Run() {
	var rawLines []string

	// 1. Gather Inputs
	if r.options.InputList != "" {
		fileTargets, err := readLines(r.options.InputList)
		if err != nil {
			logger.Error("Could not read input file: %s", err)
		} else {
			rawLines = append(rawLines, fileTargets...)
		}
	}

	if r.options.Target != "" {
		rawLines = append(rawLines, r.options.Target)
	}

	if hasStdin() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" {
				rawLines = append(rawLines, line)
			}
		}
	}

	if len(rawLines) == 0 {
		return
	}

	// 2. Parse and Group
	// map[host] -> []ports
	targets := core.ParseTargets(rawLines)
	logger.Info("Found %d unique targets from %d inputs", len(targets), len(rawLines))

	// 2. Temp Directory for XMLs
	tempDir, err := os.MkdirTemp("", "chainmap-scans")
	if err != nil {
		logger.Error("Failed to create temporary directory: %s", err)
		return
	}
	defer os.RemoveAll(tempDir) // Cleanup on exit

	// 3. Worker Pool
	jobs := make(chan string, len(targets))
	var wg sync.WaitGroup

	// Start Workers
	for i := 0; i < r.options.Threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for host := range jobs {
				ports := targets[host]
				r.scanTarget(host, ports, tempDir)
			}
		}()
	}

	// Send Jobs
	for host := range targets {
		jobs <- host
	}
	close(jobs)

	wg.Wait()

	// 4. Merge Results
	xmlFiles, err := filepath.Glob(filepath.Join(tempDir, "*.xml"))
	if err != nil {
		logger.Error("Failed to list scan results: %s", err)
		return
	}

	if len(xmlFiles) > 0 {
		logger.Info("Merging %d scan results into %s", len(xmlFiles), r.options.OutputFile)
		if err := core.MergeXMLs(xmlFiles, r.options.OutputFile); err != nil {
			logger.Error("Failed to merge XML results: %s", err)
		} else {
			logger.Success("Merged results saved to %s", r.options.OutputFile)

			// 5. Generate HTML Report
			if _, err := exec.LookPath("xsltproc"); err == nil {
				htmlFile := strings.TrimSuffix(r.options.OutputFile, filepath.Ext(r.options.OutputFile)) + ".html"
				logger.Info("Generating HTML report: %s", htmlFile)

				// Write Embedded XSLT
				if err := os.WriteFile("nmap.xsl", []byte(core.DefaultXSLT), 0644); err != nil {
					logger.Warn("Failed to write embedded nmap.xsl, using defaults: %s", err)
				}

				// Always use local nmap.xsl (we just wrote it)
				cmd := exec.Command("xsltproc", "-o", htmlFile, "nmap.xsl", r.options.OutputFile)
				if err := cmd.Run(); err != nil {
					logger.Error("Failed to generate HTML report: %s", err)
				} else {
					logger.Success("HTML report saved to %s", htmlFile)

					// Cleanup: Remove raw XML as requested by user
					if err := os.Remove(r.options.OutputFile); err != nil {
						logger.Warn("Failed to remove raw XML file: %s", err)
					} else {
						logger.Info("Removed raw XML file (keeping only HTML)")
					}

					// Cleanup: Remove temporary nmap.xsl if we created it
					_ = os.Remove("nmap.xsl")
				}
			}

			// Note: We no longer ignore nmap.xsl copying logic here
			// because we handle it inside the HTML block above by creating it momentarily.
			// Actually, for cleaner offline HTML (since we deleted XML),
			// the HTML *embedded* the style or referenced it?
			// xsltproc usually embeds the result structure.
			// BUT, if we want the XML to be viewable (we deleted it though), we'd need nmap.xsl.
			// Since we act on HTML, let's ensure we write the style first.

		}
	} else {
		logger.Info("No scan results to merge")
	}
}

func (r *Runner) scanTarget(host string, ports []string, outputDir string) {
	// Construct Ports String
	// If ports slice has empty string, it means raw host input.
	// Join all non-empty ports.
	var validPorts []string
	for _, p := range ports {
		if p != "" {
			validPorts = append(validPorts, p)
		}
	}
	portFlag := strings.Join(validPorts, ",")

	if !r.options.Silent {
		msg := fmt.Sprintf("Scanning %s", host)
		if len(validPorts) > 0 {
			msg += fmt.Sprintf(" with ports: %s", portFlag)
		}
		logger.Info(msg)
	}

	// 1. Determine Output Filename
	safeHostName := strings.ReplaceAll(host, ".", "_")
	safeHostName = strings.ReplaceAll(safeHostName, ":", "_") // Handle ipv6 if needed
	outputFile := filepath.Join(outputDir, fmt.Sprintf("%s.xml", safeHostName))

	flagsStr := r.options.NmapFlags

	// Prioritize flags based on mode: Deep > Fast > User/Default
	if r.options.DeepMode {
		if os.Geteuid() != 0 {
			logger.Warn("Deep Mode uses SYN scan (-sS) which requires root privileges. Scan may fail or degrade.")
		}
		// Deep Mode: -sS -sV -sC --script vulners --reason --version-all -T4 -Pn -n
		flagsStr = "-sS -sV -sC --script vulners --reason --version-all -T4 -Pn -n"
		if !r.options.Silent {
			logger.Info("Using Deep Scan Mode")
		}
	} else if r.options.FastMode {
		if os.Geteuid() != 0 {
			logger.Warn("Fast Mode uses SYN scan (-sS) which requires root privileges. Scan may fail or degrade.")
		}
		// Fast Mode: -sS -sV -T3 --open -Pn -n
		flagsStr = "-sS -sV -T3 --open -Pn -n"
		if !r.options.Silent {
			logger.Info("Using Fast Scan Mode")
		}
	} else if flagsStr == "" {
		// Use recommended defaults if no mode and no manual flags provided
		flagsStr = "-sV -sS -T3 -Pn -n"
	}

	args, err := shlex.Split(flagsStr)
	if err != nil {
		logger.Error("Failed to parse nmap flags: %s", err)
		return
	}

	// Add XML Output Flags
	args = append(args, "-oX", outputFile, "--webxml")

	// Add Port Flag if we have specific ports
	if len(validPorts) > 0 {
		args = append(args, "-p", portFlag)
	}

	args = append(args, host)

	// Context with Timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.options.Timeout)*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, "nmap", args...)

	// Silent execution - do not pipe to stdout
	// If verbose/debug mode existed, we might pipe stderr, but for now we keep it clean.
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			logger.Error("Timeout scanning %s", host)
		} else {
			logger.Error("Error scanning %s: %s", host, err)
		}
	} else {
		// Only log success if successful
		// logger.Success("Scan complete for %s", host) // Less noise
	}
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines, scanner.Err()
}

func hasStdin() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return stat.Mode()&os.ModeCharDevice == 0
}

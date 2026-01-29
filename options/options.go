package options

import (
	"os"

	"github.com/ihsanlearn/chainmap/logger"
	"github.com/projectdiscovery/goflags"
)

type Options struct {
	InputList  string
	Target     string
	NmapFlags  string
	Threads    int
	Timeout    int
	Silent     bool
	Version    bool
	OutputFile string
	FastMode   bool
	DeepMode   bool
}

const Version = "1.0.0"

func ParseOptions() *Options {
	opts := &Options{}

	flagSet := goflags.NewFlagSet()

	flagSet.SetDescription("Chainmap is a modular Nmap workflow tool")

	flagSet.CreateGroup("input", "Input",
		flagSet.StringVarP(&opts.InputList, "list", "l", "", "Input file containing list of IPs"),
		flagSet.StringVarP(&opts.Target, "target", "t", "", "Single target IP"),
	)

	flagSet.CreateGroup("config", "Configuration",
		flagSet.IntVarP(&opts.Threads, "threads", "c", 5, "Number of concurrent threads"),
		flagSet.IntVarP(&opts.Timeout, "timeout", "T", 10, "Timeout in minutes"),
		flagSet.StringVarP(&opts.NmapFlags, "nmap-flags", "n", "", "Nmap flags to use"),
		flagSet.StringVarP(&opts.OutputFile, "output", "o", "results.xml", "File to store merged XML results"),
		flagSet.BoolVarP(&opts.FastMode, "fast", "", false, "Fast Scan Mode"),
		flagSet.BoolVarP(&opts.DeepMode, "deep", "", false, "Deep Scan Mode"),
	)

	flagSet.CreateGroup("misc", "Optimization",
		flagSet.BoolVarP(&opts.Silent, "silent", "s", false, "Silent mode"),
		flagSet.BoolVarP(&opts.Version, "version", "V", false, "Display application version"),
	)

	if err := flagSet.Parse(); err != nil {
		logger.Error("Failed parsing flags: %s", err)
		os.Exit(1)
	}

	if opts.Version {
		logger.Info("Chainmap v%s", Version)
		os.Exit(0)
	}


	return opts
}

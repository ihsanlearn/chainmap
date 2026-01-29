package report

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/ihsanlearn/chainmap/core"
	"github.com/ihsanlearn/chainmap/logger"
)

// GenerateSummary reads the unified XML output and prints a summary
func GenerateSummary(outputFile string) {
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		logger.Info("No result file found at %s", outputFile)
		return
	}

	files := []string{outputFile}

	green := color.New(color.FgGreen).SprintfFunc()
	yellow := color.New(color.FgYellow).SprintfFunc()
	bold := color.New(color.Bold).SprintfFunc()

	fmt.Println(bold("\n--- Scan Summary ---"))

	for _, file := range files {
		nmapRun, err := core.ParseXML(file)
		if err != nil {
			logger.Error("Failed to parse %s: %s", file, err)
			continue
		}

		for _, host := range nmapRun.Hosts {
			ip := ""
			if len(host.Addresses) > 0 {
				ip = host.Addresses[0].Addr
			}

			for _, port := range host.Ports {
				if port.State.State == "open" {
					service := port.Service.Name
					product := port.Service.Product
					version := port.Service.Version

					fullService := service
					if product != "" {
						fullService += fmt.Sprintf(" (%s %s)", product, version)
					}

					fmt.Printf("%s %s -> %s\n", green(ip), yellow("%d/%s", port.PortId, port.Protocol), fullService)
				}
			}
		}
	}
	fmt.Println(bold("--------------------"))
}

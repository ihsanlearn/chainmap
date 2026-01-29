package core

import (
	"encoding/xml"
	"fmt"
	"os"

	"github.com/ihsanlearn/chainmap/logger"
	"github.com/lair-framework/go-nmap"
)

type NmapData struct {
	Host     string
	Ports    []PortInfo
	Metadata map[string]string
}

type PortInfo struct {
	Port     int
	Protocol string
	Service  string
	Version  string
	State    string
}

func ParseXML(path string) (*nmap.NmapRun, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var result nmap.NmapRun
	if err := xml.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func MergeXMLs(inputs []string, output string) error {
	var merged *nmap.NmapRun
	var totalElapsed float64

	for _, fname := range inputs {
		run, err := ParseXML(fname)
		if err != nil {
			logger.Warn("Skip file %s: %v", fname, err)
			continue
		}

		if merged == nil {
			merged = run
			totalElapsed = float64(run.RunStats.Finished.Elapsed)
		} else {
			merged.Hosts = append(merged.Hosts, run.Hosts...)

			merged.RunStats.Hosts.Up += run.RunStats.Hosts.Up
			merged.RunStats.Hosts.Down += run.RunStats.Hosts.Down
			merged.RunStats.Hosts.Total += run.RunStats.Hosts.Total

			totalElapsed += float64(run.RunStats.Finished.Elapsed)
		}
	}

	if merged == nil {
		return fmt.Errorf("no valid XML data to merge")
	}

	merged.RunStats.Finished.Elapsed = float32(totalElapsed)

	data, err := xml.MarshalIndent(merged, "", "  ")
	if err != nil {
		return err
	}

	header := `<?xml version="1.0" encoding="UTF-8"?>` + "\n" +
		`<?xml-stylesheet href="nmap.xsl" type="text/xsl"?>` + "\n"

	return os.WriteFile(output, append([]byte(header), data...), 0644)
}

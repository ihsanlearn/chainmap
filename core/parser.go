package core

import (
	"strings"
)

func ParseTargets(lines []string) map[string][]string {
	targets := make(map[string][]string)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.Contains(line, ":") {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				ip := parts[0]
				port := parts[1]

				exists := false
				for _, p := range targets[ip] {
					if p == port {
						exists = true
						break
					}
				}
				if !exists {
					targets[ip] = append(targets[ip], port)
				}
			}
		} else {
			if _, exists := targets[line]; !exists {
				targets[line] = []string{}
			}
		}
	}
	return targets
}

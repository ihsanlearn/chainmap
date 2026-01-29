package core

import (
	"strings"
)

// ParseTargets groups inputs by IP.
// Input format expected: IP:PORT or IP (domain also supported).
// Returns map[target][]ports
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
				// Dedup check could be added here if needed, but append is fast.
				// Simplest logic as requested by user.
				targets[ip] = append(targets[ip], port)
			}
		} else {
			// pure IP/Host, explicit empty slice if not exists
			if _, exists := targets[line]; !exists {
				targets[line] = []string{}
			}
		}
	}
	return targets
}

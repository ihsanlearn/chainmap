package core

import (
	"net"
	"strings"
)

// ParseTargets groups inputs by IP.
// Input format expected: IP:PORT or IP (domain also supported).
// Returns map[target][]ports
func ParseTargets(lines []string) map[string][]string {
	results := make(map[string][]string)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		host, port, err := net.SplitHostPort(line)
		if err != nil {
			// Possibly just host/ip without port
			// Or invalid format. We assume it's a host if Split fails but string is not empty.
			// Check if it looks like an IP or Domain
			// We treat 'line' as host, port string is empty.
			// net.SplitHostPort returns error if no port is present.
			results[line] = append(results[line], "")
			continue
		}

		results[host] = append(results[host], port)
	}

	// Dedup ports?
	// Nmap handles comma separation, but deduping is cleaner.
	// For now, simple append is consistent with requirements.
	// The requirement says "Mengelompokkannya berdasarkan IP unik".

	return results
}

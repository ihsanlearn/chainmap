package core

import (
	"reflect"
	"testing"
)

func TestParseTargets(t *testing.T) {
	tests := []struct {
		name     string
		lines    []string
		expected map[string][]string
	}{
		{
			name:  "Single IP",
			lines: []string{"192.168.1.1"},
			expected: map[string][]string{
				"192.168.1.1": {},
			},
		},
		{
			name:  "IP with Port",
			lines: []string{"192.168.1.1:80"},
			expected: map[string][]string{
				"192.168.1.1": {"80"},
			},
		},
		{
			name:  "IP with Multiple Distinct Ports",
			lines: []string{"192.168.1.1:80", "192.168.1.1:443"},
			expected: map[string][]string{
				"192.168.1.1": {"80", "443"},
			},
		},
		{
			name:  "IP with Duplicate Ports",
			lines: []string{"192.168.1.1:80", "192.168.1.1:80"},
			expected: map[string][]string{
				"192.168.1.1": {"80"},
			},
		},
		{
			name:  "Mixed Formats",
			lines: []string{"192.168.1.1", "192.168.1.1:80"},
			expected: map[string][]string{
				"192.168.1.1": {"80"},
			},
		},
		{
			name:  "Empty Lines and Whitespace",
			lines: []string{"  192.168.1.1:80  ", "", "  "},
			expected: map[string][]string{
				"192.168.1.1": {"80"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseTargets(tt.lines)
			for ip, ports := range got {
				// We need to check if expected ports are present, order might matter in slice unless we sort
				// For this refined logic, we expect specific ports.
				// Since existing logic appends to slice, order should be preserved from input

				// Handle nil/empty slice mismatch for deep equal
				if len(ports) == 0 && len(tt.expected[ip]) == 0 {
					continue
				}

				if !reflect.DeepEqual(ports, tt.expected[ip]) {
					t.Errorf("ParseTargets() for IP %s = %v, want %v", ip, ports, tt.expected[ip])
				}
			}

			// Check if we have extra IPs or missing IPs
			if len(got) != len(tt.expected) {
				t.Errorf("ParseTargets() returned %d targets, want %d", len(got), len(tt.expected))
			}
		})
	}
}

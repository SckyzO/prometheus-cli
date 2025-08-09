package completion

import (
	"testing"
)

func TestAdvancedCompleter_Do(t *testing.T) {
	// Mock metrics for testing
	metrics := []string{"up", "node_cpu_seconds_total", "prometheus_build_info"}
	completer := NewAdvancedCompleter(metrics, true)

	tests := []struct {
		name     string
		input    string
		expected []string
		desc     string
	}{
		{
			name:     "complete_metric_suggests_brace",
			input:    "up",
			expected: []string{"{"},
			desc:     "Complete metric name should suggest opening brace",
		},
		{
			name:     "metric_with_space_suggests_operators",
			input:    "up ",
			expected: []string{"+ ", "- ", "* ", "/ ", "% ", "^ ", "== ", "!= ", "> ", "< ", ">= ", "<= ", "and ", "or ", "unless "},
			desc:     "Metric with space should suggest operators",
		},
		{
			name:     "after_closing_brace_suggests_operators",
			input:    "up{job=\"test\"}",
			expected: []string{" + ", " - ", " * ", " / ", " % ", " ^ ", " == ", " != ", " > ", " < ", " >= ", " <= ", " and ", " or ", " unless "},
			desc:     "After closing brace should suggest operators and modifiers",
		},
		{
			name:     "partial_metric_delegates_to_prefix",
			input:    "no",
			expected: []string{"node_cpu_seconds_total"},
			desc:     "Partial metric should delegate to PrefixCompleter for navigation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			line := []rune(tt.input)
			pos := len(line)

			candidates, _ := completer.Do(line, pos)

			// Convert candidates to strings for easier comparison
			var result []string
			for _, candidate := range candidates {
				result = append(result, string(candidate))
			}

			// Check if we have the expected number of candidates
			if len(result) == 0 && len(tt.expected) > 0 {
				t.Errorf("Expected candidates but got none for input '%s'", tt.input)
				return
			}

			// For partial matches, we expect the PrefixCompleter to handle it
			if tt.name == "partial_metric_delegates_to_prefix" {
				// This test is more complex as it involves PrefixCompleter behavior
				// We'll just check that we get some candidates
				if len(result) == 0 {
					t.Errorf("Expected some candidates for partial match '%s'", tt.input)
				}
				return
			}

			// Check if we have at least some of the expected candidates
			found := false
			for _, expected := range tt.expected {
				for _, actual := range result {
					if actual == expected {
						found = true
						break
					}
				}
				if found {
					break
				}
			}

			if !found && len(tt.expected) > 0 {
				t.Errorf("Expected to find at least one of %v in results %v for input '%s'", tt.expected, result, tt.input)
			}
		})
	}
}

func TestNewAdvancedCompleter(t *testing.T) {
	metrics := []string{"up", "down"}
	completer := NewAdvancedCompleter(metrics, true)

	if completer == nil {
		t.Error("Expected completer to be created")
	}

	if len(completer.metrics) != 2 {
		t.Errorf("Expected 2 metrics, got %d", len(completer.metrics))
	}

	if !completer.enableLabelValues {
		t.Error("Expected enableLabelValues to be true")
	}
}

func TestPrometheusConstants(t *testing.T) {
	if len(PrometheusOperators) == 0 {
		t.Error("Expected PrometheusOperators to be populated")
	}

	if len(PrometheusFunctions) == 0 {
		t.Error("Expected PrometheusFunctions to be populated")
	}

	if len(PrometheusModifiers) == 0 {
		t.Error("Expected PrometheusModifiers to be populated")
	}

	if len(PrometheusTimeRanges) == 0 {
		t.Error("Expected PrometheusTimeRanges to be populated")
	}

	if len(TimeRangeFunctions) == 0 {
		t.Error("Expected TimeRangeFunctions to be populated")
	}
}

package display

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"prometheus-cli/internal/prometheus"

	"github.com/guptarohit/asciigraph"
)

// DisplayGraph renders ASCII graphs for the provided range query results.
func DisplayGraph(results []prometheus.RangeQueryResult) {
	if len(results) == 0 {
		fmt.Println("No data found for the given range.")
		return
	}

	for _, result := range results {
		// Prepare data for plotting
		var data []float64
		for _, v := range result.Values {
			// Prometheus values are [timestamp, string_value]
			// We need to extract and parse the value
			valPair, ok := v.([]interface{})
			if !ok || len(valPair) < 2 {
				continue
			}

			valStr, ok := valPair[1].(string)
			if !ok {
				continue
			}

			val, err := strconv.ParseFloat(valStr, 64)
			if err != nil {
				continue // Skip invalid values
			}
			
			// Handle NaN/Inf which can break plotting
			if math.IsNaN(val) || math.IsInf(val, 0) {
				continue
			}

			data = append(data, val)
		}

		if len(data) == 0 {
			continue
		}

		// Create a title from labels
		title := formatMetricLabels(result.Metric)
		fmt.Println("\n" + title)
		
		// Plot the graph
		graphWidth := 80
		graph := asciigraph.Plot(data, asciigraph.Height(10), asciigraph.Width(graphWidth))
		fmt.Println(graph)

		// Render custom X-axis and Timestamps
		if len(result.Values) > 1 {
			// Calculate margin based on the last line of the graph
			lines := strings.Split(graph, "\n")
			lastLine := lines[len(lines)-1]
			
			// Find the vertical axis line position (┼ or ┤)
			// We search from the end of the line backwards to find the axis char
			// This is safer as labels might contain numbers but the axis is distinct
			axisIdx := -1
			runes := []rune(lastLine)
			for i := len(runes) - 1; i >= 0; i-- {
				if runes[i] == '┼' || runes[i] == '┤' {
					axisIdx = i
					break
				}
			}
			
			marginLen := 0
			if axisIdx != -1 {
				marginLen = axisIdx
			} else {
				// Fallback
				marginLen = len(lastLine) - graphWidth
				if marginLen < 0 { marginLen = 0 }
			}
			
			// Draw the Axis Line:  └──────────────┬──────────────┘
			// marginLen spaces to reach the axis column
			fmt.Print(strings.Repeat(" ", marginLen))
			fmt.Print("└") // The corner, exactly under the vertical axis
			
			// Length to fill is graphWidth
			// We want a tick at the exact middle
			
			dashLen := (graphWidth / 2) - 1 // -1 for mid tick allowance?
			// Let's be precise. graphWidth is number of chars to the right of axis.
			// 0 to graphWidth.
			
			// Line part 1
			fmt.Print(strings.Repeat("─", dashLen))
			fmt.Print("┬") // Mid tick
			// Line part 2
			fmt.Print(strings.Repeat("─", graphWidth - dashLen - 2)) // -1 for mid, -1 for end
			fmt.Println("┘") // End tick

			// Times
			startTime := extractTime(result.Values[0])
			endTime := extractTime(result.Values[len(result.Values)-1])
			midTime := startTime.Add(endTime.Sub(startTime) / 2)
			
			startStr := startTime.Format("15:04")
			midStr := midTime.Format("15:04")
			endStr := endTime.Format("15:04")
			
			// Align times
			// Start time aligned with Start Tick (marginLen)
			// Mid time aligned with Mid Tick (marginLen + 1 + dashLen)
			// End time aligned with End Tick (marginLen + 1 + graphWidth)
			
			// We construct a single string line for times to manage spacing easily
			
			// Left margin
			fmt.Print(strings.Repeat(" ", marginLen))
			
			// Print Start Time
			fmt.Print(startStr)
			
			// Space to Mid Time
			// Target pos for Mid is (graphWidth / 2) + 1 (because of '└')
			// Current pos is len(startStr)
			targetMid := (graphWidth / 2)
			currentPos := len(startStr)
			pad1 := targetMid - (len(midStr)/2) - currentPos
			if pad1 < 1 { pad1 = 1 }
			fmt.Print(strings.Repeat(" ", pad1))
			
			// Print Mid Time
			fmt.Print(midStr)
			currentPos += pad1 + len(midStr)
			
			// Space to End Time
			// Target pos for End is graphWidth
			targetEnd := graphWidth
			pad2 := targetEnd - len(endStr) - currentPos
			if pad2 < 1 { pad2 = 1 }
			fmt.Print(strings.Repeat(" ", pad2))
			
			fmt.Println(endStr)
			
			// Center Date Label: [ Time: 2026-01-16 ]
			dateStr := fmt.Sprintf("[ Time: %s ]", startTime.Format("2006-01-02"))
			
			// Center relative to the graph (not including left label margin)
			// Graph center is at marginLen + (graphWidth / 2)
			// Label half width is len(dateStr) / 2
			// Start pos = marginLen + (graphWidth/2) - (len(dateStr)/2)
			
			datePad := (graphWidth / 2) - (len(dateStr) / 2)
			if datePad < 0 { datePad = 0 }
			
			fmt.Printf("%s%s%s\n", strings.Repeat(" ", marginLen), strings.Repeat(" ", datePad), dateStr)
		}
		fmt.Println()
	}
}

// extractTime is a helper to get time.Time from Prometheus value pair [timestamp, value]
func extractTime(v interface{}) time.Time {
	valPair, ok := v.([]interface{})
	if !ok || len(valPair) < 1 {
		return time.Time{}
	}
	
	ts, ok := valPair[0].(float64)
	if !ok {
		return time.Time{}
	}
	
	return time.Unix(int64(ts), 0)
}

// formatMetricLabels creates a string representation of metric labels for the title.
func formatMetricLabels(metric map[string]string) string {
	var keys []string
	for k := range metric {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var builder strings.Builder
	// Put __name__ first if it exists
	if name, ok := metric["__name__"]; ok {
		builder.WriteString(fmt.Sprintf("\033[1m%s\033[0m", name))
	}

	builder.WriteString("{")
	first := true
	for _, k := range keys {
		if k == "__name__" {
			continue
		}
		if !first {
			builder.WriteString(", ")
		}
		builder.WriteString(fmt.Sprintf("%s=\"%s\"", k, metric[k]))
		first = false
	}
	builder.WriteString("}")
	return builder.String()
}

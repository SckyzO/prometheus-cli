// Package display provides functionality for formatting and displaying
// Prometheus query results in a user-friendly table format.
package display

import (
	"fmt"
	"os"
	"sort"

	"prometheus-cli/internal/prometheus"

	"github.com/olekukonko/tablewriter"
)

// DisplayTable formats and displays Prometheus query results in a table format.
// It automatically organizes metrics and their labels into columns, with values
// displayed in the rightmost column. The table includes proper headers and
// handles cases where different metrics have different sets of labels.
//
// The function performs the following operations:
// 1. Collects all unique label names across all results
// 2. Creates a table with metric name, labels, and value columns
// 3. Sorts labels alphabetically for consistent display
// 4. Formats each result row with appropriate label values
// 5. Renders the table with headers and separators
//
// Parameters:
//   - results: A slice of QueryResult containing metric data from Prometheus
//
// The table format is:
// | Metric | Label1 | Label2 | ... | Value |
// |--------|--------|--------|-----|-------|
// | metric1| value1 | value2 | ... | 1.23  |
//
// If no results are provided, it displays "No results found" message.
func DisplayTable(results []prometheus.QueryResult) {
	// Handle empty results case
	if len(results) == 0 {
		fmt.Println("No results found")
		return
	}

	// Collect all unique label names across all results
	// This ensures the table includes columns for all possible labels
	labelSet := make(map[string]bool)
	for _, result := range results {
		for label := range result.Metric {
			// Skip the special __name__ label as it's handled separately as "Metric"
			if label != "__name__" {
				labelSet[label] = true
			}
		}
	}

	// Convert label set to sorted slice for consistent column ordering
	labels := make([]string, 0, len(labelSet))
	for label := range labelSet {
		labels = append(labels, label)
	}
	sort.Strings(labels)

	// Build table headers: Metric + sorted labels + Value
	headers := append([]string{"Metric"}, labels...)
	headers = append(headers, "Value")

	// Limit the number of columns to display to avoid overly wide tables
	maxColumns := 10 // Metric + 8 most important labels + Value

	if len(headers) > maxColumns {
		// Keep only the first few labels
		labels = labels[:maxColumns-2] // -2 for Metric and Value columns
		// Update headers accordingly
		headers = append([]string{"Metric"}, labels...)
		headers = append(headers, "Value")
	}

	// Truncate long headers to improve readability
	maxHeaderLength := 20
	displayHeaders := make([]string, len(headers))
	for i, header := range headers {
		if len(header) > maxHeaderLength {
			displayHeaders[i] = header[:maxHeaderLength-3] + "..."
		} else {
			displayHeaders[i] = header
		}
	}

	// Initialize table writer with stdout as destination
	table := tablewriter.NewWriter(os.Stdout)

	// Prepare data rows for bulk insertion
	rows := make([][]string, 0, len(results))
	for _, result := range results {
		// Create row with correct number of columns
		row := make([]string, len(headers))

		// Set metric name (from __name__ label or empty if not present)
		row[0] = result.Metric["__name__"]

		// Fill in label values in the correct column positions
		for i, label := range labels {
			// Column index is i+1 because metric name is at index 0
			value := result.Metric[label]
			// Truncate long values
			if len(value) > maxHeaderLength {
				row[i+1] = value[:maxHeaderLength-3] + "..."
			} else {
				row[i+1] = value
			}
		}

		// Extract and format the metric value
		// Prometheus values are returned as [timestamp, value] pairs
		if len(result.Value) >= 2 {
			if value, ok := result.Value[1].(string); ok {
				row[len(headers)-1] = value
			} else {
				// Fallback for non-string values (shouldn't normally happen)
				row[len(headers)-1] = fmt.Sprintf("%v", result.Value[1])
			}
		}

		rows = append(rows, row)
	}

	// Configure and render the table
	// Using Header() and Bulk() methods for automatic formatting with separators
	table.Header(displayHeaders)

	if err := table.Bulk(rows); err != nil {
		fmt.Printf("Error adding bulk data to table: %v\n", err)
	}

	if err := table.Render(); err != nil {
		fmt.Printf("Error rendering table: %v\n", err)
	}
}

package display

import (
	"fmt"
	"os"
	"sort"

	"promcurl/internal/prometheus"

	"github.com/olekukonko/tablewriter"
)

// Function to display results in a table
func DisplayTable(results []prometheus.QueryResult) {
	if len(results) == 0 {
		fmt.Println("No results found")
		return
	}

	// Collect all unique label names
	labelSet := make(map[string]bool)
	for _, result := range results {
		for label := range result.Metric {
			if label != "__name__" {
				labelSet[label] = true
			}
		}
	}

	// Convert to slice and sort
	labels := make([]string, 0, len(labelSet))
	for label := range labelSet {
		labels = append(labels, label)
	}
	sort.Strings(labels)

	// Create table headers
	headers := append([]string{"Metric"}, labels...)
	headers = append(headers, "Value")

	// Create the table
	table := tablewriter.NewWriter(os.Stdout)

	// Prepare data rows
	rows := make([][]string, 0, len(results))
	for _, result := range results {
		row := make([]string, len(headers))
		row[0] = result.Metric["__name__"]

		// Fill in label values
		for i, label := range labels {
			row[i+1] = result.Metric[label]
		}

		// Add the value
		if len(result.Value) >= 2 {
			value, ok := result.Value[1].(string)
			if ok {
				row[len(headers)-1] = value
			} else {
				row[len(headers)-1] = fmt.Sprintf("%v", result.Value[1])
			}
		}

		rows = append(rows, row)
	}

	// Set header and add data in bulk for automatic formatting with separator
	table.Header(headers)
	table.Bulk(rows)

	table.Render()
}

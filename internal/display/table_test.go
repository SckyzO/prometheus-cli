package display

import (
	"bytes"
	"io"
	"os"
	"testing"

	"prometheus-cli/internal/prometheus"
)

func TestDisplayTable(t *testing.T) {
	// Create a sample result
	results := []prometheus.QueryResult{
		{
			Metric: map[string]string{
				"__name__": "test_metric",
				"label1":   "value1",
				"label2":   "value2",
			},
			Value: []interface{}{1625142600, "42.5"},
		},
	}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	os.Stdout = w

	// Call the function
	DisplayTable(results)

	// Restore stdout
	if err := w.Close(); err != nil {
		t.Errorf("Failed to close writer: %v", err)
	}
	os.Stdout = oldStdout

	// Read the captured output
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Errorf("Failed to copy from reader: %v", err)
	}

	// Check the output
	if buf.Len() == 0 {
		t.Error("DisplayTable() did not produce any output")
	}

	// Check that the output contains the metric name
	if !bytes.Contains(buf.Bytes(), []byte("test_metric")) {
		t.Error("Output does not contain the metric name 'test_metric'")
	}

	// Check that the output contains the label values
	if !bytes.Contains(buf.Bytes(), []byte("value1")) {
		t.Error("Output does not contain the label value 'value1'")
	}

	if !bytes.Contains(buf.Bytes(), []byte("value2")) {
		t.Error("Output does not contain the label value 'value2'")
	}

	// Check that the output contains the metric value
	if !bytes.Contains(buf.Bytes(), []byte("42.5")) {
		t.Error("Output does not contain the metric value '42.5'")
	}
}

func TestDisplayTableNoResults(t *testing.T) {
	// Create an empty result
	var results []prometheus.QueryResult

	// Capture stdout
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	os.Stdout = w

	// Call the function
	DisplayTable(results)

	// Restore stdout
	if err := w.Close(); err != nil {
		t.Errorf("Failed to close writer: %v", err)
	}
	os.Stdout = oldStdout

	// Read the captured output
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Errorf("Failed to copy from reader: %v", err)
	}

	// Check the output
	if !bytes.Contains(buf.Bytes(), []byte("No results found")) {
		t.Error("Output does not contain 'No results found' message")
	}
}

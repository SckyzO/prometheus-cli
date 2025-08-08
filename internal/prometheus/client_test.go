package prometheus

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetMetrics(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/label/__name__/values" {
			// Return a sample response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"success","data":["metric1","metric2","metric3"]}`))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Temporarily override the PrometheusURL
	originalURL := PrometheusURL
	PrometheusURL = server.URL + "/api/v1"
	defer func() { PrometheusURL = originalURL }()

	// Call the function
	metrics, err := GetMetrics()

	// Check the results
	if err != nil {
		t.Errorf("GetMetrics() returned an error: %v", err)
	}

	if len(metrics) != 3 {
		t.Errorf("Expected 3 metrics, got %d", len(metrics))
	}

	expectedMetrics := []string{"metric1", "metric2", "metric3"}
	for i, metric := range metrics {
		if metric != expectedMetrics[i] {
			t.Errorf("Expected metric %s, got %s", expectedMetrics[i], metric)
		}
	}
}

func TestQueryPrometheus(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/query" {
			// Check the query parameter
			query := r.URL.Query().Get("query")
			if query != "test_query" {
				t.Errorf("Expected query 'test_query', got '%s'", query)
			}

			// Return a sample response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"status":"success",
				"data":{
					"resultType":"vector",
					"result":[
						{
							"metric":{"__name__":"test_metric","label1":"value1"},
							"value":[1625142600,"42.5"]
						}
					]
				}
			}`))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Temporarily override the PrometheusURL
	originalURL := PrometheusURL
	PrometheusURL = server.URL + "/api/v1"
	defer func() { PrometheusURL = originalURL }()

	// Call the function
	results, err := QueryPrometheus("test_query")

	// Check the results
	if err != nil {
		t.Errorf("QueryPrometheus() returned an error: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}

	result := results[0]
	if result.Metric["__name__"] != "test_metric" {
		t.Errorf("Expected metric name 'test_metric', got '%s'", result.Metric["__name__"])
	}

	if result.Metric["label1"] != "value1" {
		t.Errorf("Expected label1 'value1', got '%s'", result.Metric["label1"])
	}

	if len(result.Value) != 2 {
		t.Errorf("Expected value to have 2 elements, got %d", len(result.Value))
	}

	value, ok := result.Value[1].(string)
	if !ok {
		t.Errorf("Expected value[1] to be a string")
	}

	if value != "42.5" {
		t.Errorf("Expected value '42.5', got '%s'", value)
	}
}

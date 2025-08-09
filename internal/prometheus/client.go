// Package prometheus provides a client for interacting with the Prometheus HTTP API.
// It handles authentication, TLS configuration, and provides methods for querying
// metrics, labels, and label values.
package prometheus

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// PrometheusClient represents a configured client for the Prometheus API.
// It encapsulates the base URL, authentication credentials, and HTTP client
// with custom TLS settings.
type PrometheusClient struct {
	BaseURL    string       // Base URL for the Prometheus API (e.g., "http://localhost:9090/api/v1")
	Username   string       // Username for basic authentication (optional)
	Password   string       // Password for basic authentication (optional)
	HTTPClient *http.Client // Configured HTTP client with custom transport settings
}

// DefaultClient is the global Prometheus client instance used by package-level functions.
// It can be configured using the Set* functions before making API calls.
var DefaultClient = &PrometheusClient{
	BaseURL:    "http://localhost:9090/api/v1",
	HTTPClient: &http.Client{},
}

// SetPrometheusURL configures the base URL for the Prometheus API.
// The URL should include the API version path (e.g., "/api/v1").
//
// Parameters:
//   - url: The complete base URL for the Prometheus API
func SetPrometheusURL(url string) {
	DefaultClient.BaseURL = url
}

// SetBasicAuth configures HTTP basic authentication credentials.
// Both username and password must be provided for authentication to be enabled.
//
// Parameters:
//   - username: The username for basic authentication
//   - password: The password for basic authentication
func SetBasicAuth(username, password string) {
	DefaultClient.Username = username
	DefaultClient.Password = password
}

// SetTLSConfig configures TLS settings for HTTPS connections.
// When insecure is true, certificate verification is skipped (useful for self-signed certificates).
//
// Parameters:
//   - insecure: Whether to skip TLS certificate verification
func SetTLSConfig(insecure bool) {
	if insecure {
		DefaultClient.HTTPClient = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
	} else {
		DefaultClient.HTTPClient = &http.Client{}
	}
}

// doRequest performs an HTTP GET request with the client's configuration.
// It automatically adds basic authentication headers if credentials are configured.
//
// Parameters:
//   - reqURL: The complete URL to request
//
// Returns:
//   - *http.Response: The HTTP response
//   - error: Any error that occurred during the request
func (c *PrometheusClient) doRequest(reqURL string) (*http.Response, error) {
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}

	// Add basic authentication if credentials are configured
	if c.Username != "" && c.Password != "" {
		req.SetBasicAuth(c.Username, c.Password)
	}

	return c.HTTPClient.Do(req)
}

// PrometheusResponse represents the standard response format from Prometheus API.
// All Prometheus API endpoints return responses in this format.
type PrometheusResponse struct {
	Status string      `json:"status"` // Response status ("success" or "error")
	Data   interface{} `json:"data"`   // Response data (format varies by endpoint)
}

// QueryResult represents a single result from a Prometheus query.
// Each result contains metric labels and a timestamp-value pair.
type QueryResult struct {
	Metric map[string]string `json:"metric"` // Metric labels as key-value pairs
	Value  []interface{}     `json:"value"`  // [timestamp, value] pair
}

// QueryData represents the data structure for query responses.
// It contains the result type and an array of query results.
type QueryData struct {
	ResultType string        `json:"resultType"` // Type of result ("vector", "matrix", "scalar", "string")
	Result     []QueryResult `json:"result"`     // Array of query results
}

// GetMetrics retrieves all available metric names from Prometheus.
// It queries the special __name__ label to get all metric names in the system.
//
// Returns:
//   - []string: A slice of metric names
//   - error: Any error that occurred during the request
func GetMetrics() ([]string, error) {
	url := fmt.Sprintf("%s/label/__name__/values", DefaultClient.BaseURL)

	resp, err := DefaultClient.doRequest(url)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response PrometheusResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	// Convert the interface{} data to []string
	data, ok := response.Data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected data format")
	}

	metrics := make([]string, len(data))
	for i, v := range data {
		metrics[i], _ = v.(string)
	}

	return metrics, nil
}

// QueryPrometheus executes a PromQL query against Prometheus.
// It performs an instant query and returns the results.
//
// Parameters:
//   - query: The PromQL query string to execute
//
// Returns:
//   - []QueryResult: A slice of query results
//   - error: Any error that occurred during the request or parsing
func QueryPrometheus(query string) ([]QueryResult, error) {
	baseURL := fmt.Sprintf("%s/query", DefaultClient.BaseURL)

	// Build query parameters
	params := url.Values{}
	params.Add("query", query)

	// Construct the complete request URL
	reqURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	resp, err := DefaultClient.doRequest(reqURL)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response PrometheusResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	// Convert the generic response data to typed QueryData structure
	dataBytes, err := json.Marshal(response.Data)
	if err != nil {
		return nil, err
	}

	var queryData QueryData
	err = json.Unmarshal(dataBytes, &queryData)
	if err != nil {
		return nil, err
	}

	return queryData.Result, nil
}

// GetLabels retrieves all available label names from Prometheus.
// This includes both metric-specific labels and global labels.
//
// Returns:
//   - []string: A slice of label names
//   - error: Any error that occurred during the request
func GetLabels() ([]string, error) {
	url := fmt.Sprintf("%s/labels", DefaultClient.BaseURL)

	resp, err := DefaultClient.doRequest(url)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response PrometheusResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	// Convert the interface{} data to []string
	data, ok := response.Data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected data format")
	}

	labels := make([]string, len(data))
	for i, v := range data {
		labels[i], _ = v.(string)
	}

	return labels, nil
}

// GetLabelValues retrieves all possible values for a specific label.
// This is useful for autocompletion of label values in queries.
//
// Parameters:
//   - label: The name of the label to get values for
//
// Returns:
//   - []string: A slice of possible label values
//   - error: Any error that occurred during the request
func GetLabelValues(label string) ([]string, error) {
	url := fmt.Sprintf("%s/label/%s/values", DefaultClient.BaseURL, url.PathEscape(label))

	resp, err := DefaultClient.doRequest(url)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response PrometheusResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	// Convert the interface{} data to []string
	data, ok := response.Data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected data format")
	}

	values := make([]string, len(data))
	for i, v := range data {
		values[i], _ = v.(string)
	}

	return values, nil
}

package prometheus

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// PrometheusClient represents a client for the Prometheus API
type PrometheusClient struct {
	BaseURL    string
	Username   string
	Password   string
	HTTPClient *http.Client
}

// DefaultClient is the default Prometheus client
var DefaultClient = &PrometheusClient{
	BaseURL:    "http://localhost:9090/api/v1",
	HTTPClient: &http.Client{},
}

// SetPrometheusURL sets the Prometheus API URL
func SetPrometheusURL(url string) {
	DefaultClient.BaseURL = url
}

// SetBasicAuth sets the basic authentication credentials
func SetBasicAuth(username, password string) {
	DefaultClient.Username = username
	DefaultClient.Password = password
}

// SetTLSConfig sets the TLS configuration
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

// doRequest performs an HTTP request with the client's configuration
func (c *PrometheusClient) doRequest(reqURL string) (*http.Response, error) {
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}

	// Add basic auth if credentials are provided
	if c.Username != "" && c.Password != "" {
		req.SetBasicAuth(c.Username, c.Password)
	}

	return c.HTTPClient.Do(req)
}

// Structure for Prometheus responses
type PrometheusResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

// Structure for query results
type QueryResult struct {
	Metric map[string]string `json:"metric"`
	Value  []interface{}     `json:"value"`
}

// Structure for query results with typed data
type QueryData struct {
	ResultType string        `json:"resultType"`
	Result     []QueryResult `json:"result"`
}

// Function to get available metrics
func GetMetrics() ([]string, error) {
	url := fmt.Sprintf("%s/label/__name__/values", DefaultClient.BaseURL)

	resp, err := DefaultClient.doRequest(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response PrometheusResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	// Convert interface{} to []string
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

// Function to query Prometheus
func QueryPrometheus(query string) ([]QueryResult, error) {
	baseURL := fmt.Sprintf("%s/query", DefaultClient.BaseURL)

	// Create query parameters
	params := url.Values{}
	params.Add("query", query)

	// Build the complete URL
	reqURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	resp, err := DefaultClient.doRequest(reqURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response PrometheusResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	// Convert interface{} to typed structure
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

// Function to get available labels
func GetLabels() ([]string, error) {
	url := fmt.Sprintf("%s/labels", DefaultClient.BaseURL)

	resp, err := DefaultClient.doRequest(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response PrometheusResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	// Convert interface{} to []string
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

// Function to get label values for a specific label
func GetLabelValues(label string) ([]string, error) {
	url := fmt.Sprintf("%s/label/%s/values", DefaultClient.BaseURL, url.PathEscape(label))

	resp, err := DefaultClient.doRequest(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response PrometheusResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	// Convert interface{} to []string
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

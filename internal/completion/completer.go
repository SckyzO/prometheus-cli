// Package completion provides advanced autocompletion functionality for Prometheus queries.
// It implements context-aware suggestions for metrics, labels, label values, operators,
// functions, and other PromQL constructs.
package completion

import (
	"regexp"
	"strings"
	"sync"

	"prometheus-cli/internal/prometheus"

	"github.com/chzyer/readline"
)

// Cache for storing label values to avoid repeated API calls.
var (
	// labelValuesCache stores label values for each metric and label combination.
	// Structure: map[metricName]map[labelName][]values
	labelValuesCache = make(map[string]map[string][]string)

	// labelsCacheMutex protects concurrent access to the labelValuesCache.
	labelsCacheMutex sync.RWMutex
)

// Prometheus language constructs for autocompletion.
var (
	// PrometheusOperators contains all supported Prometheus operators.
	PrometheusOperators = []string{
		"+", "-", "*", "/", "%", "^",
		"==", "!=", ">", "<", ">=", "<=",
		"and", "or", "unless",
	}

	// PrometheusFunctions contains all supported Prometheus functions with opening parenthesis.
	PrometheusFunctions = []string{
		"abs(", "absent(", "absent_over_time(", "ceil(", "changes(", "clamp_max(", "clamp_min(",
		"day_of_month(", "day_of_week(", "days_in_month(", "delta(", "deriv(", "exp(", "floor(",
		"histogram_quantile(", "holt_winters(", "hour(", "idelta(", "increase(", "irate(",
		"label_join(", "label_replace(", "ln(", "log10(", "log2(", "minute(", "month(",
		"predict_linear(", "rate(", "resets(", "round(", "scalar(", "sort(", "sort_desc(",
		"sqrt(", "time(", "timestamp(", "vector(", "year(",
		"avg(", "count(", "count_values(", "min(", "max(", "sum(", "stddev(", "stdvar(",
		"bottomk(", "topk(", "quantile(",
	}

	// PrometheusModifiers contains query modifiers for aggregation operations.
	PrometheusModifiers = []string{
		"by (", "without (", "on (", "ignoring (", "group_left(", "group_right(",
	}

	// PrometheusTimeRanges contains common time range selectors.
	PrometheusTimeRanges = []string{
		"[5m]", "[10m]", "[15m]", "[30m]", "[1h]", "[2h]", "[6h]", "[12h]", "[1d]", "[7d]",
	}

	// TimeRangeFunctions contains functions that require time range selectors.
	TimeRangeFunctions = []string{
		"rate(", "increase(", "irate(", "delta(", "deriv(", "changes(", "resets(",
		"absent_over_time(", "avg_over_time(", "min_over_time(", "max_over_time(",
		"sum_over_time(", "count_over_time(", "quantile_over_time(", "stddev_over_time(",
		"stdvar_over_time(", "last_over_time(", "present_over_time(",
	}
)

// getLabelsForMetric retrieves all available labels for a specific metric.
// It queries Prometheus to get actual metric instances and extracts label names.
//
// Parameters:
//   - metricName: The name of the metric to get labels for
//
// Returns:
//   - []string: A slice of label names (excluding __name__)
//   - error: Any error that occurred during the query
func getLabelsForMetric(metricName string) ([]string, error) {
	// First, try querying the metric directly
	results, err := prometheus.QueryPrometheus(metricName)
	if err != nil {
		// If direct query fails, try with empty label selector
		results, err = prometheus.QueryPrometheus(metricName + "{}")
		if err != nil {
			return nil, err
		}
	}

	// Extract unique labels from all metric instances
	labelSet := make(map[string]bool)
	for _, result := range results {
		for label := range result.Metric {
			// Skip the special __name__ label
			if label != "__name__" {
				labelSet[label] = true
			}
		}
	}

	// Convert set to sorted slice
	labels := make([]string, 0, len(labelSet))
	for label := range labelSet {
		labels = append(labels, label)
	}

	return labels, nil
}

// getLabelValuesForMetric retrieves all possible values for a specific label of a metric.
// It uses caching to avoid repeated API calls for the same metric/label combination.
//
// Parameters:
//   - metricName: The name of the metric
//   - labelName: The name of the label to get values for
//
// Returns:
//   - []string: A slice of possible label values
//   - error: Any error that occurred during the query
func getLabelValuesForMetric(metricName, labelName string) ([]string, error) {
	// Check cache first to avoid unnecessary API calls
	labelsCacheMutex.RLock()
	if metricCache, ok := labelValuesCache[metricName]; ok {
		if values, ok := metricCache[labelName]; ok {
			labelsCacheMutex.RUnlock()
			return values, nil
		}
	}
	labelsCacheMutex.RUnlock()

	// Query Prometheus for metric instances
	results, err := prometheus.QueryPrometheus(metricName)
	if err != nil {
		// Fallback to empty label selector if direct query fails
		results, err = prometheus.QueryPrometheus(metricName + "{}")
		if err != nil {
			return nil, err
		}
	}

	// Extract unique values for the specified label
	valueSet := make(map[string]bool)
	for _, result := range results {
		if value, ok := result.Metric[labelName]; ok {
			valueSet[value] = true
		}
	}

	// Convert set to slice
	values := make([]string, 0, len(valueSet))
	for value := range valueSet {
		values = append(values, value)
	}

	// Cache the results for future use
	labelsCacheMutex.Lock()
	if _, ok := labelValuesCache[metricName]; !ok {
		labelValuesCache[metricName] = make(map[string][]string)
	}
	labelValuesCache[metricName][labelName] = values
	labelsCacheMutex.Unlock()

	return values, nil
}

// AdvancedCompleter provides context-aware autocompletion for Prometheus queries.
// It wraps readline.PrefixCompleter and adds intelligent suggestions based on
// the current query context.
type AdvancedCompleter struct {
	*readline.PrefixCompleter
	metrics           []string // Available metrics from Prometheus
	enableLabelValues bool     // Whether to provide label value suggestions
}

// NewAdvancedCompleter creates a new AdvancedCompleter instance.
// It initializes the underlying PrefixCompleter with metrics and functions,
// and configures label value completion based on the provided flag.
//
// Parameters:
//   - metrics: A slice of available metric names from Prometheus
//   - enableLabelValues: Whether to enable label value autocompletion
//
// Returns:
//   - *AdvancedCompleter: A configured completer instance
func NewAdvancedCompleter(metrics []string, enableLabelValues bool) *AdvancedCompleter {
	// Pre-allocate slice with known capacity for better performance
	items := make([]readline.PrefixCompleterInterface, 0, len(metrics)+len(PrometheusFunctions))

	// Add all metrics as completion items
	for _, metric := range metrics {
		items = append(items, readline.PcItem(metric))
	}

	// Add all functions as completion items
	for _, fn := range PrometheusFunctions {
		items = append(items, readline.PcItem(fn))
	}

	// Create the underlying prefix completer
	prefixCompleter := readline.NewPrefixCompleter(items...)

	return &AdvancedCompleter{
		PrefixCompleter:   prefixCompleter,
		metrics:           metrics,
		enableLabelValues: enableLabelValues,
	}
}

// Do implements the readline.AutoCompleter interface.
// It provides context-aware autocompletion based on the current cursor position
// and the text that has been typed so far.
//
// The completion logic follows a priority-based approach:
// 1. Handle specific contexts (after braces, operators, etc.)
// 2. Delegate to PrefixCompleter for partial matches
// 3. Provide filtered default suggestions
//
// Parameters:
//   - line: The current input line as runes
//   - pos: The cursor position within the line
//
// Returns:
//   - newLine: A slice of completion candidates
//   - length: The length of the completion prefix
func (a *AdvancedCompleter) Do(line []rune, pos int) (newLine [][]rune, length int) {
	// Extract the text up to the cursor position
	text := string(line[:pos])

	// Priority-based completion logic: handle specific contexts first

	// Case 1: After closing brace } - suggest operators, modifiers, and time ranges
	if strings.HasSuffix(strings.TrimSpace(text), "}") {
		var candidates [][]rune

		// Check if we need time ranges for functions like rate(), increase(), etc.
		needsTimeRange := false
		for _, fn := range TimeRangeFunctions {
			if strings.Contains(text, fn) {
				needsTimeRange = true
				break
			}
		}

		// Add time ranges if needed
		if needsTimeRange {
			for _, timeRange := range PrometheusTimeRanges {
				candidates = append(candidates, []rune(timeRange))
			}
		}

		// Always suggest operators with proper spacing
		for _, op := range PrometheusOperators {
			candidates = append(candidates, []rune(" "+op+" "))
		}

		// Add query modifiers
		for _, mod := range PrometheusModifiers {
			candidates = append(candidates, []rune(" "+mod))
		}

		return candidates, 0
	}

	// Case 2: metric{ - suggest available labels for the metric
	metricWithBraceRe := regexp.MustCompile(`([a-zA-Z_:][a-zA-Z0-9_:]*)\{$`)
	if matches := metricWithBraceRe.FindStringSubmatch(text); matches != nil {
		metricName := matches[1]
		labels, err := getLabelsForMetric(metricName)
		if err == nil && len(labels) > 0 {
			var candidates [][]rune
			for _, label := range labels {
				candidates = append(candidates, []rune(label+"="))
			}
			return candidates, 0
		}
	}

	// Case 3: label= - suggest quoted label values
	labelEqualsRe := regexp.MustCompile(`([a-zA-Z_][a-zA-Z0-9_]*)=$`)
	if matches := labelEqualsRe.FindStringSubmatch(text); matches != nil && a.enableLabelValues {
		// Extract metric name from the query context
		metricRe := regexp.MustCompile(`([a-zA-Z_:][a-zA-Z0-9_:]*)\{`)
		if metricMatches := metricRe.FindStringSubmatch(text); metricMatches != nil {
			metricName := metricMatches[1]
			labelName := matches[1]

			values, err := getLabelValuesForMetric(metricName, labelName)
			if err == nil && len(values) > 0 {
				var candidates [][]rune
				for _, value := range values {
					candidates = append(candidates, []rune("\""+value+"\""))
				}
				return candidates, 0
			}
		}
	}

	// Case 4: label=" - suggest label values without additional quotes
	labelEqualsQuoteRe := regexp.MustCompile(`([a-zA-Z_][a-zA-Z0-9_]*)="$`)
	if matches := labelEqualsQuoteRe.FindStringSubmatch(text); matches != nil && a.enableLabelValues {
		// Extract metric name from the query context
		metricRe := regexp.MustCompile(`([a-zA-Z_:][a-zA-Z0-9_:]*)\{`)
		if metricMatches := metricRe.FindStringSubmatch(text); metricMatches != nil {
			metricName := metricMatches[1]
			labelName := matches[1]

			values, err := getLabelValuesForMetric(metricName, labelName)
			if err == nil && len(values) > 0 {
				var candidates [][]rune
				for _, value := range values {
					candidates = append(candidates, []rune(value+"\""))
				}
				return candidates, 0
			}
		}
	}

	// Case 5: label="value" - suggest comma for additional labels or closing brace
	completeValueRe := regexp.MustCompile(`([a-zA-Z_][a-zA-Z0-9_]*)="[^"]*"$`)
	if matches := completeValueRe.FindStringSubmatch(text); matches != nil {
		return [][]rune{[]rune(","), []rune("}")}, 0
	}

	// Case 6: After comma - suggest remaining available labels
	afterCommaRe := regexp.MustCompile(`([a-zA-Z_:][a-zA-Z0-9_:]*)\{.*,\s*$`)
	if matches := afterCommaRe.FindStringSubmatch(text); matches != nil {
		metricName := matches[1]
		labels, err := getLabelsForMetric(metricName)
		if err == nil && len(labels) > 0 {
			// Parse already used labels to avoid duplicates
			usedLabels := make(map[string]bool)
			labelPairRe := regexp.MustCompile(`([a-zA-Z_][a-zA-Z0-9_]*)="[^"]*"`)
			pairs := labelPairRe.FindAllStringSubmatch(text, -1)
			for _, pair := range pairs {
				if len(pair) > 1 {
					usedLabels[pair[1]] = true
				}
			}

			// Suggest only unused labels
			var candidates [][]rune
			for _, label := range labels {
				if !usedLabels[label] {
					candidates = append(candidates, []rune(label+"="))
				}
			}
			if len(candidates) > 0 {
				return candidates, 0
			}
		}
	}

	// Case 7: Complete metric name - suggest opening brace for label selection
	words := strings.Fields(text)
	if len(words) > 0 {
		lastWord := words[len(words)-1]
		// Only suggest if there's no space after the word (cursor is at end of word)
		if strings.HasSuffix(text, lastWord) {
			for _, metric := range a.metrics {
				if metric == lastWord {
					return [][]rune{[]rune("{")}, 0
				}
			}
		}
	}

	// Case 8: Complete metric name with space - suggest operators and modifiers
	metricWithSpaceRe := regexp.MustCompile(`^([a-zA-Z_:][a-zA-Z0-9_:]*)\s+$`)
	if matches := metricWithSpaceRe.FindStringSubmatch(text); matches != nil {
		metricName := matches[1]
		for _, metric := range a.metrics {
			if metric == metricName {
				var candidates [][]rune
				for _, op := range PrometheusOperators {
					candidates = append(candidates, []rune(op+" "))
				}
				for _, mod := range PrometheusModifiers {
					candidates = append(candidates, []rune(mod))
				}
				return candidates, 0
			}
		}
	}

	// Case 9: Inside functions - delegate to PrefixCompleter for metric navigation
	functionContextRe := regexp.MustCompile(`(rate|increase|sum|avg|count|min|max)\(\s*$`)
	if matches := functionContextRe.FindStringSubmatch(text); matches != nil {
		return a.PrefixCompleter.Do(line, pos)
	}

	// Case 10: After operators - suggest metrics and functions
	afterOperatorRe := regexp.MustCompile(`(\+|\-|\*|\/|\%|\^|==|!=|>|<|>=|<=|\sand\s|\sor\s|\sunless\s)\s*$`)
	if matches := afterOperatorRe.FindStringSubmatch(text); matches != nil {
		var candidates [][]rune
		for _, metric := range a.metrics {
			candidates = append(candidates, []rune(metric))
		}
		for _, fn := range PrometheusFunctions {
			candidates = append(candidates, []rune(fn))
		}
		return candidates, 0
	}

	// Default case: delegate to PrefixCompleter for partial matches and navigation
	return a.PrefixCompleter.Do(line, pos)
}

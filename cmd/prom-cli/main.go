// Package main provides the command-line interface for Prometheus CLI,
// a tool for querying Prometheus metrics with advanced autocompletion.
package main

import (
	"fmt"
	"os"
	"strings"

	"prometheus-cli/internal/completion"
	"prometheus-cli/internal/display"
	"prometheus-cli/internal/prometheus"

	"github.com/alecthomas/kingpin/v2"
	"github.com/chzyer/readline"
	"github.com/prometheus/common/version"
)

// Command-line flags for configuring the application behavior.
var (
	// url specifies the Prometheus server URL to connect to.
	url = kingpin.Flag("url", "Prometheus server URL.").Default("http://localhost:9090").String()

	// username specifies the username for basic authentication.
	username = kingpin.Flag("username", "Username for basic authentication.").String()

	// password specifies the password for basic authentication.
	password = kingpin.Flag("password", "Password for basic authentication.").String()

	// insecure determines whether to skip TLS certificate verification.
	insecure = kingpin.Flag("insecure", "Skip TLS certificate verification.").Bool()

	// enableLabelValues controls whether label values autocompletion is enabled.
	enableLabelValues = kingpin.Flag("enable-label-values", "Enable autocompletion for label values.").Default("true").Bool()
)

// main is the entry point of the Prometheus CLI application.
// It initializes the Prometheus client, sets up autocompletion, and runs the interactive query loop.
func main() {
	// Configure command-line argument parsing
	kingpin.Version(version.Print("prom-cli"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	// Initialize Prometheus client with user-provided configuration
	prometheus.SetPrometheusURL(*url + "/api/v1")
	prometheus.SetBasicAuth(*username, *password)
	prometheus.SetTLSConfig(*insecure)

	// Load available metrics from Prometheus for autocompletion
	fmt.Print("Loading metrics...")
	metrics, err := prometheus.GetMetrics()
	if err != nil {
		fmt.Printf("\rError getting metrics: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\rLoaded %d metrics successfully.\n", len(metrics))

	// Initialize the advanced autocompletion system
	completer := completion.NewAdvancedCompleter(metrics, *enableLabelValues)

	// Set up readline interface with autocompletion and history
	l, err := readline.NewEx(&readline.Config{
		Prompt:          "\033[31m»\033[0m ",
		HistoryFile:     "/tmp/prom_cli_history.tmp",
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := l.Close(); err != nil {
			fmt.Printf("Error closing readline: %v\n", err)
		}
	}()

	// Display welcome message and feature information
	printWelcomeMessage()

	// Run the main interactive query loop
	runQueryLoop(l)
}

// printWelcomeMessage displays the welcome message and available features.
func printWelcomeMessage() {
	fmt.Println("Enter Prometheus queries. Press Ctrl+C to exit.")
	fmt.Println("Features enabled:")
	fmt.Println("  - 📊 Metrics autocompletion")
	fmt.Println("  - 🏷️  Labels and values autocompletion" + func() string {
		if !*enableLabelValues {
			return " (disabled with --enable-label-values=false)"
		}
		return ""
	}())
	fmt.Println("  - ⚡ Prometheus expressions autocompletion (operators, functions, time ranges)")
	fmt.Println("  - 🔧 Smart context-aware suggestions")
	fmt.Println()
	fmt.Println("💡 Tips:")
	fmt.Println("  - Type 'rat' + Tab → 'rate('")
	fmt.Println("  - After metric{} + Tab → operators and modifiers")
	fmt.Println("  - Inside functions + Tab → metrics")
	fmt.Println("  - After operators + Tab → metrics and functions")
}

// runQueryLoop runs the main interactive loop for processing user queries.
func runQueryLoop(l *readline.Instance) {
	for {
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			fmt.Println("Exiting...")
			break
		} else if err != nil {
			break
		}

		query := strings.TrimSpace(line)
		if query == "" {
			continue
		}

		// Execute the Prometheus query and display results
		results, err := prometheus.QueryPrometheus(query)
		if err != nil {
			fmt.Printf("Error executing query: %v\n", err)
			continue
		}

		display.DisplayTable(results)
	}
}

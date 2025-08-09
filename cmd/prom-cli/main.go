// Package main provides the command-line interface for Prometheus CLI,
// a tool for querying Prometheus metrics with advanced autocompletion.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"prometheus-cli/internal/completion"
	"prometheus-cli/internal/display"
	"prometheus-cli/internal/prometheus"

	kingpin "github.com/alecthomas/kingpin/v2"
	"github.com/chzyer/readline"
	"github.com/prometheus/common/version"
)

// Command-line flags for configuring the application behavior.
var (
	// Prometheus Connection Flags
	url      = kingpin.Flag("url", "Prometheus server URL.").Default("http://localhost:9090").String()
	username = kingpin.Flag("username", "Username for basic authentication.").String()
	password = kingpin.Flag("password", "Password for basic authentication.").String()
	insecure = kingpin.Flag("insecure", "Skip TLS certificate verification.").Bool()

	// Autocompletion Flags
	enableLabelValues = kingpin.Flag("enable-label-values", "Enable autocompletion for label values.").Default("true").Bool()

	// History Flags
	historyFile    = kingpin.Flag("history-file", "Path to the command history file.").String()
	persistHistory = kingpin.Flag("persist-history", "Do not delete the history file on exit.").Bool()

	// Display and Utility Flags
	debug = kingpin.Flag("debug", "Enable verbose error output for debugging.").Bool()
	tips  = kingpin.Flag("tips", "Display detailed feature and usage tips on startup.").Bool()
)

// main is the entry point of the Prometheus CLI application.
// It initializes the Prometheus client, sets up autocompletion, and runs the interactive query loop.
func main() {
	// Configure command-line argument parsing
	kingpin.Version(version.Print("prom-cli"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	// Display welcome message and feature information if tips are enabled
	if *tips {
		printWelcomeMessage()
	} else {
		fmt.Println("Enter Prometheus queries. Press Ctrl+C to exit.")
	}

	// Initialize Prometheus client with user-provided configuration
	if *debug {
		fmt.Printf("Debug: Setting Prometheus URL to %s/api/v1\n", *url)
		fmt.Printf("Debug: Setting Basic Auth with username: %s\n", *username)
		fmt.Printf("Debug: Setting TLS InsecureSkipVerify to %t\n", *insecure)
	}
	prometheus.SetPrometheusURL(*url + "/api/v1")
	prometheus.SetBasicAuth(*username, *password)
	prometheus.SetTLSConfig(*insecure)

	// Load available metrics from Prometheus for autocompletion
	fmt.Print("Loading metrics...")
	metrics, err := prometheus.GetMetrics()
	if err != nil {
		if *debug {
			fmt.Printf("\rError getting metrics: %v\n", err)
		} else {
			fmt.Printf("\rError getting metrics. Use --debug for more details.\n")
		}
		os.Exit(1)
	}
	fmt.Printf("\rLoaded %d metrics successfully.\n", len(metrics))

	// Initialize the advanced autocompletion system
	completer := completion.NewAdvancedCompleter(metrics, *enableLabelValues)

	// Determine the history file path and handle persistence.
	var historyFilePath string
	var shouldRemoveHistoryFile bool

	if *historyFile != "" {
		if filepath.IsAbs(*historyFile) {
			historyFilePath = *historyFile
		} else {
			// Join with current working directory if a relative path is provided
			cwd, err := os.Getwd()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: could not get current working directory: %v\n", err)
				historyFilePath = *historyFile // Fallback to direct use if cwd fails
			} else {
				historyFilePath = filepath.Join(cwd, *historyFile)
			}
		}
		shouldRemoveHistoryFile = !*persistHistory
		if *debug {
			fmt.Printf("Debug: Using specified history file: %s (persist: %t)\n", historyFilePath, *persistHistory)
		}
	} else {
		// Create a temporary file for command history.
		tempFile, err := os.CreateTemp("", "prom_cli_history_*.tmp")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not create temp history file: %v\n", err)
		} else {
			historyFilePath = tempFile.Name()
			shouldRemoveHistoryFile = !*persistHistory // Still remove if not explicitly persisted
			if err := tempFile.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: could not close temp history file: %v\n", err)
			}
			if *debug {
				fmt.Printf("Debug: Using temporary history file: %s (persist: %t)\n", historyFilePath, *persistHistory)
			}
		}
	}

	// Ensure the history file exists or create it if it doesn't.
	if historyFilePath != "" {
		if _, err := os.Stat(historyFilePath); os.IsNotExist(err) {
			if *debug {
				fmt.Printf("Debug: History file %s does not exist, creating it.\n", historyFilePath)
			}
			file, err := os.Create(historyFilePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: could not create history file %s: %v\n", historyFilePath, err)
			} else {
				if err := file.Close(); err != nil {
					fmt.Fprintf(os.Stderr, "Warning: could not close history file %s: %v\n", historyFilePath, err)
				}
			}
		}

		// Schedule the history file to be removed if persistence is not requested.
		if *debug {
			fmt.Printf("Debug: shouldRemoveHistoryFile is %t before defer registration.\n", shouldRemoveHistoryFile)
		}
		if shouldRemoveHistoryFile {
			defer func() {
				if *debug {
					fmt.Printf("Debug: Inside defer. shouldRemoveHistoryFile is %t. Attempting to remove history file: %s\n", shouldRemoveHistoryFile, historyFilePath)
				}
				if err := os.Remove(historyFilePath); err != nil {
					fmt.Fprintf(os.Stderr, "Warning: could not remove history file %s: %v\n", historyFilePath, err)
				}
			}()
		}
	}

	// Set up readline interface with autocompletion and history.
	l, err := readline.NewEx(&readline.Config{
		Prompt:          "\033[31m»\033[0m ",
		HistoryFile:     historyFilePath,
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

	// Run the main interactive query loop
	runQueryLoop(l)
}

// printWelcomeMessage displays the welcome message and available features.
func printWelcomeMessage() {
	fmt.Println("Enter Prometheus queries. Press Ctrl+C to exit.")

	if *tips {
		fmt.Print(`
✨ Features:
	 - Metric Names: Smart autocompletion for all available Prometheus metrics
	 - Label Names: Context-aware label suggestions when typing "metric{"
	 - Label Values: Real-time label value suggestions with caching for performance
	 - PromQL Expressions: Complete support for operators, built-in functions, time range selectors, and query modifiers
	 - Context-Aware Suggestions: Intelligent suggestions based on cursor position and query context
	 - Navigation Support: Tab completion with arrow key navigation for easy selection

💡 Tips:
	 - Type 'rat' + Tab -> 'rate('
	 - After metric{} + Tab -> operators and modifiers
	 - Inside functions + Tab -> metrics
	 - After operators + Tab -> metrics and functions
`)
	}
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
			if *debug {
				fmt.Printf("Error executing query: %v\n", err)
			} else {
				fmt.Printf("Error executing query. Use --debug for more details.\n")
			}
			continue
		}

		display.DisplayTable(results)
	}
}

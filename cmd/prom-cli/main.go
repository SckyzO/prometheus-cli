// Package main provides the command-line interface for Prometheus CLI,
// a tool for querying Prometheus metrics with advanced autocompletion.
package main

import (
	"fmt"
	"os"
	"path/filepath" // Added for filepath.Join
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

	// debug enables verbose error output for debugging purposes.
	debug = kingpin.Flag("debug", "Enable verbose error output for debugging.").Bool()

	// historyFile specifies the path to the command history file.
	historyFile = kingpin.Flag("history-file", "Path to the command history file. If not set, a temporary file is used.").String()

	// persistHistory determines whether the history file should be persisted across sessions.
	persistHistory = kingpin.Flag("persist-history", "Do not delete the history file on exit. Only applicable if --history-file is set or a temporary file is used.").Bool()
)

// main is the entry point of the Prometheus CLI application.
// It initializes the Prometheus client, sets up autocompletion, and runs the interactive query loop.
func main() {
	// Configure command-line argument parsing
	kingpin.Version(version.Print("prom-cli"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

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
			defer tempFile.Close()
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
				file.Close()
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

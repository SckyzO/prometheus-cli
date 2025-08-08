# 🔍 PromCurl

[![Build and Test](https://github.com/yourusername/promcurl/actions/workflows/build.yml/badge.svg)](https://github.com/yourusername/promcurl/actions/workflows/build.yml)
[![Go Version](https://img.shields.io/badge/Go-1.21%2B-blue.svg)](https://golang.org/doc/devel/release.html#go1.21)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A command-line tool for querying Prometheus metrics with autocompletion.

## 📋 Table of Contents

- [🔍 PromCurl](#-promcurl)
  - [📋 Table of Contents](#-table-of-contents)
  - [📝 Overview](#-overview)
  - [✨ Features](#-features)
  - [📥 Installation](#-installation)
    - [From Source](#from-source)
    - [Using Go](#using-go)
  - [🚀 Usage](#-usage)
    - [Command Line Options](#command-line-options)
    - [Examples](#examples)
  - [🛠️ Development](#️-development)
    - [Prerequisites](#prerequisites)
    - [Building](#building)
    - [Testing](#testing)
    - [Cross-compilation](#cross-compilation)
  - [📁 Project Structure](#-project-structure)
  - [📄 License](#-license)
  - [📝 Version History](#-version-history)

## 📝 Overview

PromCurl is a tool that allows you to query Prometheus metrics from the command line. It provides autocompletion for metric names, making it easier to construct queries without having to remember exact metric names.

## ✨ Features

- 🔍 Query Prometheus metrics from the command line
- 🔄 Advanced autocompletion:
  - Metric names autocompletion
  - Label names autocompletion (after typing '{')
  - Label values autocompletion (after typing 'label=') with optional real values for common labels (job, instance, env, etc.)
  - Comma and closing brace autocompletion (after selecting a label value)
  - Recursive label autocompletion (after selecting a comma) with real values for common labels
  - Common operators autocompletion
- 📊 Display results in a formatted table with clear header separation
- 🔁 Interactive mode with continuous query support (exit with Ctrl+C)
- 💻 Cross-platform support (Linux, macOS, Windows)
- 🔒 TLS support with optional certificate verification
- 🔐 Basic authentication support

## 📥 Installation

### From Source

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/promcurl.git
   cd promcurl
   ```

2. Build the binary:
   ```
   make build
   ```

3. (Optional) Install the binary to your PATH:
   ```
   cp promcurl /usr/local/bin/
   ```

### Using Go

```
go install github.com/yourusername/promcurl/cmd/promcurl@latest
```

## 🚀 Usage

1. Make sure Prometheus is running and accessible at `http://localhost:9090` (default).

2. Run PromCurl:
   ```
   ./bin/promcurl
   ```

3. Enter a Prometheus query when prompted. You can use Tab for autocompletion of metric names.

4. The results will be displayed in a formatted table with a clear separation between headers and data.

5. The application remains active after executing a query, allowing you to enter additional queries.

6. To exit the application, press Ctrl+C.

### Command Line Options

PromCurl supports the following command line options:

```
--url                  Prometheus server URL (default: http://localhost:9090)
--username             Username for basic authentication
--password             Password for basic authentication
--insecure             Skip TLS certificate verification
--enable-label-values  Enable autocompletion for label values (may increase startup time)
--help, -h             Show help
--version              Show version information
```

### Examples

Basic usage with default settings:
```
./bin/promcurl
```

Connecting to a custom Prometheus server:
```
./bin/promcurl --url="http://prometheus-server:9090"
```

Using a custom path:
```
./bin/promcurl --url="https://monitoring.example.com/prometheus"
```

With authentication:
```
./bin/promcurl --url="https://prometheus-server:9090" --username="admin" --password="secret"
```

With special characters in credentials:
```
./bin/promcurl --url="https://prometheus-server:9090" --username="user@domain" --password="p@ssw0rd!"
```

Skipping TLS verification (for self-signed certificates):
```
./bin/promcurl --url="https://prometheus-server:9090" --insecure
```

Enabling label values autocompletion (with animated loading dots):
```
./bin/promcurl --enable-label-values
```
This will display animated loading dots while fetching label values for autocompletion.

Note: When using values with special characters, you can use either of these formats:
- With equals sign: `--option="value with spaces"`
- With space: `--option "value with spaces"`

## 🛠️ Development

### Prerequisites

- Go 1.21 or later
- Make (optional, for using the Makefile)

### Building

```
make build
```

### Testing

```
make test
```

### Cross-compilation

Build for all platforms:
```
make build-all
```

Or for specific platforms:
```
make build-linux
make build-windows
make build-macos
```

## 📁 Project Structure

- `cmd/promcurl/`: Main application entry point
- `internal/prometheus/`: Prometheus API client
- `internal/display/`: Table display functionality
- `test/`: Integration tests
- `python/`: Original Python implementation (v1.0)

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 📝 Version History

- v2.8: Fixed performance issue with label values autocompletion
- v2.7: Enhanced recursive label autocompletion with real values for common labels
- v2.6: Added recursive label autocompletion after comma
- v2.5: Added autocompletion for comma and closing brace after label values
- v2.4: Improved loading animation with animated dots for label values autocompletion
- v2.3: Added loading indicator for label values autocompletion
- v2.2: Made label values autocompletion optional to improve stability
- v2.1: Enhanced label values autocompletion with real values for common labels
- v2.0: Go implementation with improved structure and features
- v1.0: Original Python implementation

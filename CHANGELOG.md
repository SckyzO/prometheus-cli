### v2.4.0 - ASCII Graph Mode 📈
**Major Features:**
- **📈 ASCII Graphs**: Added support for visualizing range queries directly in the terminal using ASCII charts.
- **⏱️ Range Queries**: Implemented `query_range` support in the client.
- **📅 Time Flags**: Added `--start`, `--end`, `--step`, and `--graph` flags for controlling time ranges and visualization.
- **🧠 Smart Time Parsing**: Supports RFC3339, SQL-style timestamps, and relative durations (e.g., `1h`).

### v2.3.0 - Configuration File & Go Upgrade 🛠️
**Major Features:**
- **⚙️ Configuration File**: Added support for YAML configuration file (`~/.prom-cli.yaml` or via `--config`). Centralize your settings like URL, auth, and preferences without long CLI flags.
- **📄 Example Config**: Included `prom-cli.example.yaml` to quickly get started with configuration.

**Technical Enhancements:**
- **🚀 Go 1.24 Upgrade**: Updated project and CI to use Go 1.24 for better performance and latest language features.
- **📦 Dependencies**: Updated all core dependencies to their latest versions.
- **🧪 CI/CD**: Enhanced integration tests and CI workflow for robustness.

### v2.2.0 - Enhanced Authentication & Security 🔐
**Major Features:**
- **🔐 Enhanced Authentication**: Added support for `PROM_USERNAME` and `PROM_PASSWORD` environment variables.
- **📂 Password File**: Added `--password-file` flag for secure password handling.
- **🛡️ Security**: Improved security by allowing password input via file instead of command line flags, preventing password exposure in process lists.

### v2.1.0 - Enhanced Usability and Display 🚀
**Major Features:**
- **📝 Configurable History**: Added `--history-file` and `--persist-history` flags for flexible command history management.
- **🐛 Improved Debugging**: Enhanced `--debug` flag with more verbose output for initialization and error diagnosis.
- **💡 Optional Tips**: Introduced `--tips` flag to control the display of detailed feature and usage tips on startup.
- **📊 Optimized Table Display**: Improved table rendering for queries with many labels, preventing excessive width issues.

**Technical Enhancements:**
- Refined error handling and logging for better debugging experience.
- Improved command-line option parsing and validation.
- Implemented intelligent column limiting and header truncation for better readability.
- Fixed compilation issues with help text formatting.

### v2.0.0 - Complete Go Rewrite 🚀
**Major Features:**
- **🔄 Complete rewrite in Go** for better performance and reliability
- **🏗️ Clean architecture** with modular design (`cmd/`, `internal/` structure)
- **🔧 Advanced autocompletion system** with context-aware suggestions
- **📊 Intelligent table display** with automatic column organization
- **🔐 Enhanced security** with full TLS and authentication support
- **⚡ Performance optimizations** with caching and efficient data structures
- **🧪 Comprehensive testing** with unit and integration tests
- **📦 Cross-platform binaries** with automated GitHub Actions builds
- **🎛️ Flexible configuration** with extensive command-line options

**Autocompletion Improvements:**
- Smart metric name completion with fuzzy matching
- Context-aware label and label value suggestions
- Complete PromQL syntax support (operators, functions, modifiers)
- Efficient caching system for label values
- Tab navigation with arrow key support
- Priority-based suggestion ordering

**Technical Enhancements:**
- Refactored codebase with proper Go package structure
- Automated testing and continuous integration
- Memory-efficient data structures and algorithms
- Robust error handling and user feedback

### v1.0.0 - Original Python Implementation
- Basic Prometheus querying functionality
- Simple table output
- Basic metric name autocompletion

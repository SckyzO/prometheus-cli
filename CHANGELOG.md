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

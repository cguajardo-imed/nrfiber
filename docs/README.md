# nrfiber Documentation

Welcome to the nrfiber documentation! This directory contains comprehensive guides and resources for using nrfiber with Fiber v2 and v3.

## 📚 Documentation Overview

### Getting Started

- **[Main README](../README.md)** - Start here! Overview, installation, and basic usage
- **[Examples](../examples/)** - Working code examples for both Fiber v2 and v3

### Guides

1. **[Notice Custom Errors](notice-custom-errors.md)**
   - How to handle custom errors with nrfiber
   - Avoiding duplicate error reporting
   - Configuring New Relic error collection
   - Best practices for error handling
   - **Audience**: All users who need error reporting

2. **[Migration Guide](MIGRATION_GUIDE.md)**
   - Migrating between nrfiber versions
   - Upgrading from v1.x to v2.x/v3.x
   - Moving from Fiber v2 to v3
   - Choosing between v2 and v3 modules
   - Common migration issues and solutions
   - **Audience**: Users upgrading from older versions

3. **[IDE Setup](IDE_SETUP.md)**
   - Configuring Visual Studio Code
   - Setting up GoLand/IntelliJ IDEA
   - Vim/Neovim configuration
   - Sublime Text setup
   - General Go development tools
   - **Audience**: Developers contributing to nrfiber or building applications with it

4. **[Troubleshooting Guide](TROUBLESHOOTING.md)**
   - Installation issues
   - Middleware problems
   - Transaction and segment issues
   - Error reporting problems
   - Performance optimization
   - Integration conflicts
   - Debug mode and diagnostics
   - **Audience**: Users experiencing issues or errors

## 🎯 Quick Navigation

### By Use Case

**I want to:**

- **Get started quickly** → [Main README](../README.md) + [Examples](../examples/)
- **Report custom errors** → [Notice Custom Errors](notice-custom-errors.md)
- **Upgrade nrfiber** → [Migration Guide](MIGRATION_GUIDE.md)
- **Setup my IDE** → [IDE Setup](IDE_SETUP.md)
- **Fix an issue** → [Troubleshooting Guide](TROUBLESHOOTING.md)

### By Experience Level

**Beginner:**
1. Start with [Main README](../README.md)
2. Try [Examples](../examples/)
3. Read [Notice Custom Errors](notice-custom-errors.md)

**Intermediate:**
1. Review [Troubleshooting Guide](TROUBLESHOOTING.md)
2. Configure your [IDE Setup](IDE_SETUP.md)
3. Optimize using best practices

**Advanced:**
1. Study [Migration Guide](MIGRATION_GUIDE.md)
2. Review all configuration options
3. Contribute improvements

## 📖 Document Summaries

### notice-custom-errors.md
Learn how to properly handle custom errors in your Fiber application with New Relic APM. This guide covers:
- Implementing custom error types
- Enabling error reporting in nrfiber
- Avoiding duplicate errors in New Relic dashboard
- Configuring server-side error filtering
- Best practices for meaningful error messages

**Key Topics:**
- Custom error structs
- `ConfigNoticeErrorEnabled()`
- `ConfigStatusCodeIgnored()`
- Server-side configuration
- Transaction attributes

### MIGRATION_GUIDE.md
Comprehensive guide for migrating between nrfiber versions. Covers:
- Moving to v3.x (Fiber v3 support)
- Upgrading from v1.x to v2.x
- Choosing between v2 and v3 modules
- Handling breaking changes
- Testing after migration
- Rollback procedures

**Key Topics:**
- Import path changes
- Context type differences (pointer vs interface)
- Function signature updates
- Configuration changes
- Common issues and solutions

### IDE_SETUP.md
Complete IDE configuration guide for nrfiber development. Includes:
- Visual Studio Code setup and extensions
- GoLand/IntelliJ IDEA configuration
- Vim/Neovim with vim-go and coc.nvim
- Sublime Text with GoSublime
- General Go tooling
- Project-specific configurations

**Key Topics:**
- Language server (gopls) setup
- Linting and formatting
- Debug configurations
- Task automation
- Keyboard shortcuts

### TROUBLESHOOTING.md
Detailed troubleshooting guide for common issues. Addresses:
- Installation and dependency problems
- Middleware configuration issues
- Transaction and context problems
- Error reporting failures
- Performance concerns
- Integration conflicts
- New Relic dashboard issues

**Key Topics:**
- Error diagnosis
- Debug logging
- Common pitfalls
- Performance optimization
- Testing strategies

## 🔧 Configuration Examples

### Minimal Setup

```go
app.Use(nrfiber.Middleware(nrApp))
```

### With Error Reporting

```go
app.Use(nrfiber.Middleware(nrApp,
    nrfiber.ConfigNoticeErrorEnabled(true),
))
```

### Full Configuration

```go
app.Use(nrfiber.Middleware(nrApp,
    nrfiber.ConfigNoticeErrorEnabled(true),
    nrfiber.ConfigStatusCodeIgnored([]int{404, 401}),
    nrfiber.ConfigCustomTransactionNameFunc(customNameFunc),
))
```

## 🆘 Getting Help

### Before Opening an Issue

1. Check the [Troubleshooting Guide](TROUBLESHOOTING.md)
2. Review [Migration Guide](MIGRATION_GUIDE.md) if upgrading
3. Look at working [Examples](../examples/)
4. Enable debug logging

### When Opening an Issue

Include:
- Go version: `go version`
- Fiber version: `go list -m github.com/gofiber/fiber/v3`
- nrfiber version: `go list -m github.com/cguajardo-imed/nrfiber/v3`
- Minimal reproducible example
- Error messages and logs
- What you expected vs what happened

## 🤝 Contributing

Found an error in the documentation? Want to add more examples?

1. Fork the repository
2. Make your changes
3. Test thoroughly
4. Submit a pull request

See the main [README](../README.md) for contribution guidelines.

## 📝 Documentation Standards

Our documentation follows these principles:

- **Clear**: Easy to understand for all skill levels
- **Complete**: Covers all features and edge cases
- **Current**: Updated with each release
- **Code-First**: Includes working examples
- **Practical**: Focuses on real-world scenarios

## 🔗 External Resources

### Official Documentation

- [Fiber v3 Docs](https://docs.gofiber.io/)
- [Fiber v2 Docs](https://docs.gofiber.io/v2.x/)
- [New Relic Go Agent](https://docs.newrelic.com/docs/apm/agents/go-agent/)
- [Go Documentation](https://golang.org/doc/)

### Tutorials and Guides

- [New Relic APM Overview](https://docs.newrelic.com/docs/apm/)
- [Go Best Practices](https://golang.org/doc/effective_go)
- [Fiber Best Practices](https://docs.gofiber.io/guide/best-practices)

## 📊 Visual Resources

This directory includes screenshots for documentation:

- `new_relic_errors_before_ignore.png` - Duplicate error example
- `new_relic_errors_after_ignore.png` - Clean error reporting
- `newrelic_error_collection_setting.png` - Configuration screenshot

## 📅 Last Updated

Documentation is continuously updated. Check the [CHANGELOG](../CHANGELOG.md) for recent changes.

---

**Need something not covered here?** [Open an issue](https://github.com/cguajardo-imed/nrfiber/issues) and let us know!
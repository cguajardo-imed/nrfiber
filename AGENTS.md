# Agent Rules and Project Context

## Project Overview

**nrfiber** is a Go library that provides automatic instrumentation for integrating [New Relic](https://newrelic.com) Application Performance Monitoring (APM) with the [GoFiber](https://gofiber.io) web framework.

**Compatible with both Fiber v2 and v3** - Separate modules for each version (no build tags needed).

### Key Information
- **Language**: Go 1.25.0+
- **Module Paths**: 
  - v3: `github.com/cguajardo-imed/nrfiber/v3`
  - v2: `github.com/cguajardo-imed/nrfiber/v2`
- **Current Version**: 3.0.3 (managed in `version` file)
- **Primary Dependencies**:
  - `github.com/gofiber/fiber/v3` (for v3) or `github.com/gofiber/fiber/v2` (for v2) - Web framework
  - `github.com/newrelic/go-agent/v3` - New Relic APM agent
  - `github.com/stretchr/testify` - Testing framework
- **Module Structure**:
  - Root directory: Fiber v2 support (for backwards compatibility)
  - `v2/` directory: Fiber v2 support (`nrfiber/v2`)
  - `v3/` directory: Fiber v3 support (`nrfiber/v3`)
  - No build tags required - use appropriate import path

## Project Structure

```
nrfiber/
├── .github/          # GitHub Actions workflows (CI/CD automation)
├── docs/             # Additional documentation
├── examples/         # Example applications
│   ├── fiber-v2-basic/     # Fiber v2 basic example
│   ├── fiber-v3-basic/     # Fiber v3 basic example
│   └── fiber-v3-advanced/  # Fiber v3 advanced example
├── v2/               # Fiber v2 module
│   ├── nrfiber.go    # Main middleware for v2
│   ├── config.go     # Configuration for v2
│   ├── nrfiber_test.go
│   ├── config_test.go
│   ├── go.mod        # Separate module (nrfiber/v2)
│   └── README.md
├── v3/               # Fiber v3 module
│   ├── nrfiber.go    # Main middleware for v3
│   ├── config.go     # Configuration for v3
│   ├── nrfiber_test.go
│   ├── config_test.go
│   ├── go.mod        # Separate module (nrfiber/v3)
│   └── README.md
├── config.go         # Root config (v2 for backwards compatibility)
├── config_test.go    # Configuration tests for root
├── nrfiber.go        # Main middleware implementation for root
├── nrfiber_test.go   # Core functionality tests for root
├── version           # Version file for automated releases
├── go.mod            # Go module dependencies (root)
└── README.md         # Project documentation
```

## Core Functionality

### Main Components

1. **Middleware Function**: `Middleware(app *newrelic.Application, configs ...*config) fiber.Handler`
   - Creates Fiber middleware that instruments HTTP requests
   - Wraps requests in New Relic transactions
   - Automatically captures request/response metadata

2. **Context Helper**: 
   - **v3**: `FromContext(c fiber.Ctx) *newrelic.Transaction`
   - **v2**: `FromContext(c *fiber.Ctx) *newrelic.Transaction`
   - Retrieves the New Relic transaction from Fiber context
   - Enables custom segment creation in route handlers

3. **Configuration Options**:
   - `ConfigNoticeErrorEnabled(bool)` - Enable/disable error reporting to New Relic
   - `ConfigStatusCodeIgnored([]int)` - Specify HTTP status codes to ignore for error reporting
   - `ConfigCustomTransactionNameFunc(func)` - Customize transaction naming logic

4. **Utility Functions**:
   - **v3**: `Send(c fiber.Ctx, segmentName string)` - Quick segment creation and execution
   - **v2**: `Send(c *fiber.Ctx, segmentName string)` - Quick segment creation and execution
   - Now includes nil check for transaction safety

## Development Rules

### Code Style and Standards

1. **Follow Go Best Practices**:
   - Use `gofmt` for code formatting
   - Follow effective Go guidelines
   - Keep functions small and focused
   - Use meaningful variable and function names

2. **Error Handling**:
   - Always handle errors explicitly
   - Use Fiber's error handling patterns (`*fiber.Error`)
   - Support custom error notice configuration
   - Add nil checks for transactions to prevent panics

3. **Module Structure**:
   - v3 code lives in the `v3/` subdirectory
   - v2 code lives in the `v2/` subdirectory
   - Each has its own `go.mod` file
   - No build tags needed - import path determines version

4. **Testing Requirements**:
   - Write unit tests for all new functionality
   - Use `testify` assertions for clarity
   - Test files must end with `_test.go`
   - Maintain or improve test coverage

5. **Compatibility**:
   - Maintain compatibility with both Fiber v2 and v3
   - Use separate modules to isolate version-specific implementations
   - Keep API consistent between versions where possible
   - Maintain compatibility with New Relic Go Agent v3
   - Support Go 1.25.0 and above

### API Design Principles

1. **Middleware Integration**:
   - Middleware must be registered BEFORE other middlewares and routes
   - Should gracefully handle nil application (no-op behavior)
   - Must not break the Fiber middleware chain

2. **Transaction Management**:
   - Transactions must be properly started and ended
   - Use defer patterns for cleanup
   - Store transactions in Fiber context using New Relic's context helpers
   - **v3**: Use `c.SetContext()` and `c.Context()` for context management
   - **v2**: Use `c.SetUserContext()` and `c.Context()` for context management

3. **HTTP Request Conversion**:
   - Properly convert Fiber requests to `http.Request` for New Relic
   - Preserve all headers, query parameters, and metadata
   - Handle both HTTP and HTTPS schemes correctly

### Configuration Pattern

- Use functional options pattern with config structs
- Keep configuration keys as internal constants
- Provide sensible defaults (e.g., error notice disabled by default)
- Configuration should be optional and composable
- Function signatures must match Fiber version (pointer vs non-pointer context)

### Differences Between Fiber v2 and v3

1. **Import Path**:
   - **v2**: `import "github.com/cguajardo-imed/nrfiber/v2"`
   - **v3**: `import "github.com/cguajardo-imed/nrfiber/v3"`

2. **Context Type**:
   - **v2**: Uses `*fiber.Ctx` (pointer)
   - **v3**: Uses `fiber.Ctx` (interface)

3. **Context Storage**:
   - **v2**: Uses `c.SetUserContext()` 
   - **v3**: Uses `c.SetContext()`

4. **Slice Contains**:
   - **v2**: Manual loop for slice contains
   - **v3**: Can use `slices.Contains()` from standard library

## CI/CD Rules

### Automated Releases

1. **Version Management**:
   - Version is stored in the `version` file at repository root
   - Format: semantic versioning (e.g., `3.0.0`)
   - To release a new version, update the `version` file and push to `main`

2. **Release Behavior**:
   - Automatic releases trigger on every push to `main` branch
   - If version matches latest release, incremental letter suffix is added (`v3.0.0a`, `v3.0.0b`, etc.)
   - If version differs from latest release, new release with that version is created
   - Release notes are auto-generated from commit messages

3. **Manual Releases**:
   - Can be triggered manually from GitHub Actions tab
   - Follow same versioning logic as automatic releases

### Commit Messages

- Write clear, descriptive commit messages
- Use conventional commit format when possible (e.g., `feat:`, `fix:`, `docs:`)
- Commits are used to generate release notes

### Testing Both Versions

Always test both Fiber v2 and v3 before pushing:

```bash
# Test Fiber v3 (v3 directory)
cd v3
go test ./...
cd ..

# Test Fiber v2 (v2 directory)
cd v2
go test ./...
cd ..

# Test examples
cd examples/fiber-v3-basic
go run main.go &
# Test v3...
cd ../fiber-v2-basic
go run main.go &
# Test v2...
```

## Contributing Guidelines

1. **Before Making Changes**:
   - Review existing code and tests
   - Understand the New Relic transaction lifecycle
   - Understand Fiber v3 middleware patterns

2. **When Adding Features**:
   - Update relevant documentation
   - Add configuration options if needed
   - Provide examples in comments or `example/` directory
   - Add tests covering new functionality

3. **When Fixing Bugs**:
   - Write a test that reproduces the bug
   - Fix the bug
   - Ensure the test passes
   - Consider edge cases

4. **Documentation**:
   - Update README.md for user-facing changes
   - Add guides to `docs/` directory for complex features
   - Include code examples in documentation

## Important Notes for AI Agents

1. **Do Not**:
   - Break backward compatibility without explicit approval
   - Remove existing configuration options
   - Change the middleware signature
   - Modify the version file unless instructed for a release
   - Generate markdown files unless explicitly requested
   - Make changes that only work in one Fiber version
   - Mix v2 and v3 code in the same directory

2. **Always**:
   - Maintain thread-safety in middleware
   - Ensure proper transaction cleanup (use defer)
   - Test changes with both Fiber v2 and v3
   - Test changes with both Fiber and New Relic contexts
   - Preserve existing error handling behavior
   - Add nil checks for transaction safety
   - Maintain parallel implementations for both versions

3. **Context Awareness**:
   - New Relic transactions are stored in Fiber context via `newrelic.NewContext()`
   - Retrieved using `newrelic.FromContext()` wrapper
   - Context must propagate through the entire request lifecycle
   - **v3**: Uses `c.SetContext()` and `c.Context()`
   - **v2**: Uses `c.SetUserContext()` and `c.Context()`

4. **Performance Considerations**:
   - Middleware adds minimal overhead
   - Avoid unnecessary allocations in hot paths
   - Use efficient header conversion methods
   - Defer only what needs cleanup

## Common Patterns

### Adding a New Configuration Option

You must add the configuration to both versions:

**For v3/config.go and v2/config.go:**
```go
// 1. Add constant in both files
const configKeyNewOption = "NewOption"

// 2. Create public config function
// NOTE: Function signature differs between versions
// v3: func(c fiber.Ctx) T
// v2: func(c *fiber.Ctx) T
func ConfigNewOption(value SomeType) *config {
    return &config{
        key:   configKeyNewOption,
        value: value,
    }
}

// 3. Add getter function (same in both)
func getNewOption(configMap map[string]any) SomeType {
    if val, ok := configMap[configKeyNewOption]; ok {
        if v, ok := val.(SomeType); ok {
            return v
        }
    }
    return defaultValue
}

// 4. Use in middleware (same in both)
newOptionValue := getNewOption(configMap)
```

### Creating Custom Segments

**Fiber v3:**
```go
app.Get("/route", func(c fiber.Ctx) error {
    txn := nrfiber.FromContext(c)
    if txn == nil {
        return c.SendString("No transaction")
    }
    segment := txn.StartSegment("Operation Name")
    defer segment.End()
    
    // Your operation here
    
    return c.JSON(result)
})
```

**Fiber v2:**
```go
app.Get("/route", func(c *fiber.Ctx) error {
    txn := nrfiber.FromContext(c)
    if txn == nil {
        return c.SendString("No transaction")
    }
    segment := txn.StartSegment("Operation Name")
    defer segment.End()
    
    // Your operation here
    
    return c.JSON(result)
})
```

## Resources

- [Fiber v3 Documentation](https://docs.gofiber.io/)
- [New Relic Go Agent Guide](https://docs.newrelic.com/docs/apm/agents/go-agent/)
- [Project Repository](https://github.com/cguajardo-imed/nrfiber)

## Questions or Issues?

When encountering issues:
1. Check existing tests for expected behavior
2. Review New Relic and Fiber documentation
3. Examine the example code in `example/main.go`
4. Look at related issues in the GitHub repository

## Important

- Do not create any .md or .txt files as result of any task, for summary or to explain some new feature or newly implemented code, unless is explicitly asked for it.

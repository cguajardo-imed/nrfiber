# nrfiber v2

New Relic instrumentation middleware for [Fiber v2](https://docs.gofiber.io/v2.x/) web framework.

## Installation

```bash
go get -u github.com/cguajardo-imed/nrfiber/v2
```

## Usage

```go
package main

import (
    "log"
    "os"

    "github.com/cguajardo-imed/nrfiber/v2"
    "github.com/gofiber/fiber/v2"
    "github.com/newrelic/go-agent/v3/newrelic"
)

func main() {
    // Initialize New Relic application
    nrApp, err := newrelic.NewApplication(
        newrelic.ConfigAppName("my-fiber-v2-app"),
        newrelic.ConfigLicense(os.Getenv("NEW_RELIC_LICENSE_KEY")),
        newrelic.ConfigDistributedTracerEnabled(true),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Create Fiber app
    app := fiber.New()

    // Add nrfiber middleware (must be before routes)
    app.Use(nrfiber.Middleware(nrApp))

    // Define routes
    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("Hello, World!")
    })

    // Route with custom segment
    app.Get("/users/:id", func(c *fiber.Ctx) error {
        // Get transaction from context
        txn := nrfiber.FromContext(c)
        
        // Create custom segment
        segment := txn.StartSegment("Database - Get User")
        defer segment.End()
        
        // Your database operation here
        user := getUserFromDB(c.Params("id"))
        
        return c.JSON(user)
    })

    log.Fatal(app.Listen(":3000"))
}
```

## Configuration Options

### Enable Error Reporting

```go
app.Use(nrfiber.Middleware(
    nrApp,
    nrfiber.ConfigNoticeErrorEnabled(true),
))
```

### Ignore Specific Status Codes

```go
app.Use(nrfiber.Middleware(
    nrApp,
    nrfiber.ConfigNoticeErrorEnabled(true),
    nrfiber.ConfigStatusCodeIgnored([]int{404, 401}),
))
```

### Custom Transaction Naming

```go
app.Use(nrfiber.Middleware(
    nrApp,
    nrfiber.ConfigCustomTransactionNameFunc(func(c *fiber.Ctx) string {
        return fmt.Sprintf("%s %s", c.Method(), c.Route().Path)
    }),
))
```

## API Reference

### Functions

#### `Middleware(app *newrelic.Application, configs ...*config) fiber.Handler`

Creates the nrfiber middleware handler. Must be registered before your routes.

**Parameters:**
- `app` - New Relic application instance
- `configs` - Optional configuration options

**Returns:** Fiber middleware handler

#### `FromContext(c *fiber.Ctx) *newrelic.Transaction`

Retrieves the New Relic transaction from the Fiber context.

**Parameters:**
- `c` - Fiber context pointer

**Returns:** New Relic transaction or nil

#### `Send(c *fiber.Ctx, segmentName string)`

Quick helper to create and execute a named segment.

**Parameters:**
- `c` - Fiber context pointer
- `segmentName` - Name for the segment

### Configuration Options

#### `ConfigNoticeErrorEnabled(enabled bool) *config`

Enable or disable automatic error reporting to New Relic.

#### `ConfigStatusCodeIgnored(statusCodes []int) *config`

Specify HTTP status codes that should not be reported as errors.

#### `ConfigCustomTransactionNameFunc(fn func(*fiber.Ctx) string) *config`

Provide a custom function to name transactions for better grouping in New Relic.

## Examples

See the [examples directory](../examples/fiber-v2-basic/) for complete working examples.

## Key Differences from v3

This package is designed for **Fiber v2**. If you're using Fiber v3, use the v3 version instead:

```go
import "github.com/cguajardo-imed/nrfiber/v3"
```

### Context Type

**v2 (this package):**
```go
func(c *fiber.Ctx) error {
    // c is a pointer
    txn := nrfiber.FromContext(c)
    return c.JSON(data)
}
```

**v3:**
```go
func(c fiber.Ctx) error {
    // c is an interface
    txn := nrfiber.FromContext(c)
    return c.JSON(data)
}
```

### Body Parsing

**v2 (this package):**
```go
if err := c.BodyParser(&req); err != nil {
    return err
}
```

**v3:**
```go
if err := c.Bind().JSON(&req); err != nil {
    return err
}
```

## Requirements

- Go 1.25.0 or higher
- Fiber v2
- New Relic Go Agent v3

## Testing

```bash
cd v2
go test -v
```

## License

See the main [LICENSE](../LICENSE) file.

## Contributing

See the main [README](../README.md) for contribution guidelines.

## Learn More

- [Main nrfiber Documentation](../README.md)
- [Examples](../examples/)
- [Fiber v2 Documentation](https://docs.gofiber.io/v2.x/)
- [New Relic Go Agent](https://docs.newrelic.com/docs/apm/agents/go-agent/)
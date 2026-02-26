# nrfiber Examples

This directory contains example applications demonstrating how to use **nrfiber** with both Fiber v2 and Fiber v3.

## Available Examples

### Basic Examples

| Example | Description | Fiber Version | Documentation |
|---------|-------------|---------------|---------------|
| [fiber-v3-basic](./fiber-v3-basic) | Basic New Relic integration with Fiber v3 | v3 | [README](./fiber-v3-basic/README.md) |
| [fiber-v2-basic](./fiber-v2-basic) | Basic New Relic integration with Fiber v2 (separate module) | v2 | [README](./fiber-v2-basic/README.md) |

### Advanced Examples

| Example | Description | Fiber Version | Documentation |
|---------|-------------|---------------|---------------|
| [fiber-v3-advanced](./fiber-v3-advanced) | Advanced features with custom configuration | v3 | Coming soon |

## Quick Start

### Running Fiber v3 Examples

```bash
# Navigate to example directory
cd fiber-v3-basic

# Set your New Relic license key (optional)
export NEW_RELIC_LICENSE_KEY="your-key-here"

# Run the example
go run main.go
```

### Running Fiber v2 Examples

```bash
# Navigate to example directory
cd fiber-v2-basic

# Set your New Relic license key (optional)
export NEW_RELIC_LICENSE_KEY="your-key-here"

# Run normally - no build tags needed!
go run main.go
```

## What You'll Learn

### Basic Examples

- ✅ How to set up nrfiber middleware
- ✅ Creating custom transaction segments
- ✅ Tracking database operations
- ✅ Monitoring API endpoints
- ✅ Basic error handling

### Advanced Examples

- ✅ Custom error handling with New Relic
- ✅ Ignoring specific HTTP status codes
- ✅ Custom transaction naming strategies
- ✅ Adding custom attributes to transactions
- ✅ Multiple segment types (database, external, custom)
- ✅ Complex business logic tracking

## Prerequisites

- Go 1.25.0 or higher
- New Relic account (optional - examples will run without it)

## Environment Variables

All examples support the following environment variables:

- `NEW_RELIC_LICENSE_KEY` - Your New Relic license key (optional)
- `PORT` - Server port (default: 3000)

## Key Differences Between Fiber v2 and v3

### Import Paths

**Fiber v3:**
```go
import "github.com/cguajardo-imed/nrfiber/v3"
import "github.com/gofiber/fiber/v3"
```

**Fiber v2:**
```go
import "github.com/cguajardo-imed/nrfiber/v2"
import "github.com/gofiber/fiber/v2"
```

### Context Type

**Fiber v3:**
```go
app.Get("/route", func(c fiber.Ctx) error {
    // c is fiber.Ctx (interface)
    return c.JSON(data)
})
```

**Fiber v2:**
```go
app.Get("/route", func(c *fiber.Ctx) error {
    // c is *fiber.Ctx (pointer)
    return c.JSON(data)
})
```

### JSON Body Parsing

**Fiber v3:**
```go
if err := c.Bind().JSON(&req); err != nil {
    // handle error
}
```

**Fiber v2:**
```go
if err := c.BodyParser(&req); err != nil {
    // handle error
}
```

### No Build Tags Required

Both versions work without build tags - just use the correct import path:

**Fiber v3:**
```bash
cd fiber-v3-basic
go run main.go
go build
```

**Fiber v2:**
```bash
cd fiber-v2-basic
go run main.go
go build
```

## Testing the Examples

Once an example is running, you can test it with curl:

```bash
# Basic health check
curl http://localhost:3000/health

# Get data
curl http://localhost:3000/users/123

# Search with query parameters
curl "http://localhost:3000/search?q=golang&limit=5"

# POST request
curl -X POST http://localhost:3000/users \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com"}'

# Test error handling
curl http://localhost:3000/error
```

## Viewing Data in New Relic

After running an example and making some requests:

1. Log in to [New Relic One](https://one.newrelic.com)
2. Navigate to **APM & Services**
3. Find your application (e.g., "fiber-v3-basic-example")
4. Explore:
   - **Transactions**: View endpoint performance
   - **Databases**: See custom segment metrics
   - **Errors**: Track error rates and details
   - **Service Maps**: Visualize application architecture

## Common Issues

### "NEW_RELIC_LICENSE_KEY not set"

This is just a warning. The application will run in disabled mode without sending telemetry to New Relic.

### Wrong import path for Fiber version

Make sure you're using the correct import path:

```bash
# For Fiber v2
import "github.com/cguajardo-imed/nrfiber/v2"

# For Fiber v3
import "github.com/cguajardo-imed/nrfiber/v3"
```

### Module not found errors

Make sure you're in the example directory and run:

```bash
go mod download
```

## Contributing

Feel free to contribute more examples! Some ideas:

- Database integration examples (PostgreSQL, MongoDB)
- Microservices communication
- GraphQL APIs
- WebSocket support
- Rate limiting with monitoring
- Authentication/Authorization tracking

## Learn More

- [nrfiber Main Documentation](../README.md)
- [Fiber v3 Documentation](https://docs.gofiber.io/)
- [Fiber v2 Documentation](https://docs.gofiber.io/v2.x/)
- [New Relic Go Agent Guide](https://docs.newrelic.com/docs/apm/agents/go-agent/)
- [New Relic APM Best Practices](https://docs.newrelic.com/docs/new-relic-solutions/best-practices-guides/full-stack-observability/apm-best-practices-guide/)

## License

These examples are part of the nrfiber project. See the main [LICENSE](../LICENSE) file for details.
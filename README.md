# nrfiber

Provides auto instrumentation for [NewRelic](https://newrelic.com) and [GoFiber](https://gofiber.io).

**Compatible with both Fiber v2 and v3!**

> **Note:** The root package (`github.com/cguajardo-imed/nrfiber`) is a backwards compatibility wrapper that re-exports from `v2`. 
> For new projects, use the versioned packages directly:
> - `github.com/cguajardo-imed/nrfiber/v2` for Fiber v2
> - `github.com/cguajardo-imed/nrfiber/v3` for Fiber v3

## Install

### For Fiber v3
```shell
go get -u github.com/cguajardo-imed/nrfiber/v3
```

### For Fiber v2
```shell
go get -u github.com/cguajardo-imed/nrfiber/v2
```

**Note:** Separate modules for each version - no build tags needed!

## Usage

### Fiber v3

Register the middleware and use created transaction to add custom segments. Basic usage is below:

```go
package main

import (
	"github.com/cguajardo-imed/nrfiber/v3"
	"github.com/gofiber/fiber/v3"
	"github.com/newrelic/go-agent/v3/newrelic"
	"log"
)

func main() {
	app := fiber.New()
	nr, err := newrelic.NewApplication(
		newrelic.ConfigEnabled(true),
		newrelic.ConfigAppName("demo"),
		newrelic.ConfigLicense("license-key"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Add the nrfiber middleware before other middlewares or routes
	app.Use(nrfiber.Middleware(nr))

	// Use created transaction to create custom segments
	app.Get("/cart", func(c fiber.Ctx) error {
		txn := nrfiber.FromContext(c)
		segment := txn.StartSegment("Price Calculation")
		defer segment.End()

		// calculate the price

		return c.SendString("Cart processed")
	})
	
	app.Listen(":3000")
}
```

### Fiber v2

The API is identical for Fiber v2:

```go
package main

import (
	"github.com/cguajardo-imed/nrfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/newrelic/go-agent/v3/newrelic"
	"log"
)

func main() {
	app := fiber.New()
	nr, err := newrelic.NewApplication(
		newrelic.ConfigEnabled(true),
		newrelic.ConfigAppName("demo"),
		newrelic.ConfigLicense("license-key"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Add the nrfiber middleware before other middlewares or routes
	app.Use(nrfiber.Middleware(nr))

	// Use created transaction to create custom segments
	app.Get("/cart", func(c *fiber.Ctx) error {
		txn := nrfiber.FromContext(c)
		segment := txn.StartSegment("Price Calculation")
		defer segment.End()

		// calculate the price

		return c.SendString("Cart processed")
	})
	
	app.Listen(":3000")
}
```

**Note:** The only difference is the context parameter type (`fiber.Ctx` for v3 vs `*fiber.Ctx` for v2).

## Configuration Options

### Enable Error Reporting

```go
app.Use(nrfiber.Middleware(nr, nrfiber.ConfigNoticeErrorEnabled(true)))
```

### Ignore Specific Status Codes

```go
app.Use(nrfiber.Middleware(nr, 
	nrfiber.ConfigNoticeErrorEnabled(true),
	nrfiber.ConfigStatusCodeIgnored([]int{404, 401}),
))
```

### Custom Transaction Names

```go
// For Fiber v3
customNameFunc := func(c fiber.Ctx) string {
	return fmt.Sprintf("%s %s", c.Method(), c.Route().Path)
}

// For Fiber v2
customNameFunc := func(c *fiber.Ctx) string {
	return fmt.Sprintf("%s %s", c.Method(), c.Route().Path)
}

app.Use(nrfiber.Middleware(nr, 
	nrfiber.ConfigCustomTransactionNameFunc(customNameFunc),
))
```

## Quick Segment Helper

The `Send` function provides a convenient way to create and track segments:

```go
// For Fiber v3
app.Get("/operation", func(c fiber.Ctx) error {
	nrfiber.Send(c, "CustomOperation")
	
	// Your code here
	
	return c.SendString("Done")
})

// For Fiber v2
app.Get("/operation", func(c *fiber.Ctx) error {
	nrfiber.Send(c, "CustomOperation")
	
	// Your code here
	
	return c.SendString("Done")
})
```

## Version Selection

This library provides separate modules for Fiber v2 and v3:

- **v3**: `import "github.com/cguajardo-imed/nrfiber/v3"`
- **v2**: `import "github.com/cguajardo-imed/nrfiber/v2"`

No build tags are needed - just use the appropriate import path.

Examples:
```go
// For Fiber v3
import "github.com/cguajardo-imed/nrfiber/v3"
import "github.com/gofiber/fiber/v3"

// For Fiber v2
import "github.com/cguajardo-imed/nrfiber/v2"
import "github.com/gofiber/fiber/v2"
```

## Examples

We provide comprehensive examples for both Fiber v2 and v3:

- **[All Examples](examples/)** - Overview and quick start guide
- **[Fiber v3 Basic](examples/fiber-v3-basic/)** - Basic integration with Fiber v3
- **[Fiber v2 Basic](examples/fiber-v2-basic/)** - Basic integration with Fiber v2
- **[Fiber v3 Advanced](examples/fiber-v3-advanced/)** - Advanced features including:
  - Custom error handling
  - Transaction attributes
  - Ignored status codes
  - Custom transaction naming
  - Multiple segment types

Each example includes:
- Complete working code
- Detailed README with API documentation
- curl commands for testing
- Explanations of key concepts

### Running Examples

```bash
# Fiber v3 example
cd examples/fiber-v3-basic
go run main.go

# Fiber v2 example
cd examples/fiber-v2-basic
go run main.go
```

## API Reference

### Middleware

```go
func Middleware(app *newrelic.Application, configs ...*config) fiber.Handler
```

Creates a Fiber middleware that instruments HTTP requests with New Relic.

**Parameters:**
- `app`: New Relic application instance
- `configs`: Optional configuration options

**Returns:** Fiber handler function

### FromContext

```go
// For Fiber v3
func FromContext(c fiber.Ctx) *newrelic.Transaction

// For Fiber v2
func FromContext(c *fiber.Ctx) *newrelic.Transaction
```

Retrieves the New Relic transaction from the Fiber context.

**Returns:** Transaction or `nil` if not found

### Send

```go
// For Fiber v3
func Send(c fiber.Ctx, segmentName string)

// For Fiber v2
func Send(c *fiber.Ctx, segmentName string)
```

Convenience function to create and track a named segment.

## Version Compatibility

| nrfiber Version | Fiber v2 | Fiber v3 | Go Version |
|-----------------|----------|----------|------------|
| v3.0.0+ | ✅ (via `/v2`) | ✅ (via `/v3`) | 1.25.0+ |
| v2.0.0+ | ✅ | ❌ | 1.25.0+ |

## Project Structure

```
nrfiber/
├── v2/               # Fiber v2 support (nrfiber/v2)
│   ├── nrfiber.go
│   ├── config.go
│   ├── go.mod
│   └── README.md
├── v3/               # Fiber v3 support (nrfiber/v3)
│   ├── nrfiber.go
│   ├── config.go
│   ├── go.mod
│   └── README.md
├── examples/         # Example applications
└── docs/             # Additional documentation
```

## Guides

- [Notice Custom Errors](docs/notice-custom-errors.md) - Handling custom errors and avoiding duplicates
- [Migration Guide](docs/MIGRATION_GUIDE.md) - Migrating from older versions
- [IDE Setup](docs/IDE_SETUP.md) - Configuring your IDE for nrfiber development
- [Troubleshooting Guide](docs/TROUBLESHOOTING.md) - Common issues and solutions

## CI/CD

This project uses GitHub Actions for automated releases:

- **Automatic Releases**: A new release is created automatically on every push to the `main` branch
- **Manual Releases**: You can trigger a release manually from the GitHub Actions tab
- **Versioning Strategy**: 
  - The version is read from the `version` file in the repository
  - If the version matches the latest release, an incremental letter suffix is added (e.g., `v3.0.0a`, `v3.0.0b`, ..., `v3.0.0z`)
  - If the version is different from the latest release, a new release is created with that version
  - Release notes are automatically generated from commits

To create a new major/minor/patch release, simply update the `version` file and push to `main`.

## FAQ

### Q: Which version should I use?

**A:** Use the import path that matches your Fiber version:
- Fiber v2: `github.com/cguajardo-imed/nrfiber/v2`
- Fiber v3: `github.com/cguajardo-imed/nrfiber/v3`

### Q: Can I use both v2 and v3 in the same project?

**A:** No, you must choose one Fiber version per project.

### Q: What's the difference between v2 and v3?

**A:** The main difference is the context type:
- v2 uses `*fiber.Ctx` (pointer)
- v3 uses `fiber.Ctx` (interface)

The API is otherwise identical.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project follows the license specified in the LICENSE file.
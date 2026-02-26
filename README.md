# nrfiber

Provides auto instrumentation for [NewRelic](https://newrelic.com) and [GoFiber](https://gofiber.io).

**Compatible with both Fiber v2 and v3!**

## Install

### For Fiber v3
```shell
go get -u github.com/cguajardo-imed/nrfiber/v3
```

### For Fiber v2
```shell
go get -u github.com/cguajardo-imed/nrfiber/v2
```

**Note:** No build tags required - the version is determined by the import path you use.

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
	app.Get("/cart", func(ctx fiber.Ctx) error {
		txn := nrfiber.FromContext(ctx)
		segment := txn.StartSegment("Price Calculation")
		defer segment.End()

		// calculate the price

		return ctx.SendString("Cart processed")
	})
	
	app.Listen(":3000")
}
```

### Fiber v2

The API is identical for Fiber v2, just use the `fiberv2` build tag:

```go
package main

import (
	"github.com/cguajardo-imed/nrfiber/v3"
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
	app.Get("/cart", func(ctx *fiber.Ctx) error {
		txn := nrfiber.FromContext(ctx)
		segment := txn.StartSegment("Price Calculation")
		defer segment.End()

		// calculate the price

		return ctx.SendString("Cart processed")
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

## Guides

- [Notice Custom Errors](docs/notice-custom-errors.md)

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

## Contributing

Feel free to add anything useful or fix something.
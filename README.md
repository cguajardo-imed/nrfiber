# nrfiber

Provides auto instrumentation for [NewRelic](https://newrelic.com) and [GoFiber](https://gofiber.io).

## Install

```shell
go get -u github.com/cguajardo-imed/nrfiber
```

## Usage

Register the middleware and use created transaction to add another segments. Basic usage is below

```go
package main

import (
	"github.com/cguajardo-imed/nrfiber"
	"github.com/gofiber/fiber/v3"
	"github.com/newrelic/go-agent/v3/newrelic"
	"log"
)

func main() {
	app := fiber.New()
	nr, err := newrelic.NewApplication(newrelic.ConfigEnabled(true), newrelic.ConfigAppName("demo"), newrelic.ConfigLicense("license-key"))
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

		return nil
	})
	app.Listen(":3000")
}
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

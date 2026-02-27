// Package nrfiber provides New Relic instrumentation for the Fiber v2 web framework.
// This package re-exports github.com/cguajardo-imed/nrfiber/v2 for backwards compatibility.
//
// For new projects, consider importing the v2 or v3 packages directly:
//   - github.com/cguajardo-imed/nrfiber/v2 for Fiber v2
//   - github.com/cguajardo-imed/nrfiber/v3 for Fiber v3
package nrfiber

import (
	v2 "github.com/cguajardo-imed/nrfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/newrelic/go-agent/v3/newrelic"
)

// FromContext returns the Transaction from the context if present, and nil otherwise.
func FromContext(c *fiber.Ctx) *newrelic.Transaction {
	return v2.FromContext(c)
}

// Middleware creates Fiber middleware that instruments requests.
//
// app := fiber.New()
// // Add the nrfiber middleware before other middlewares or routes:
// app.Use(nrfiber.Middleware(app))
func Middleware(app *newrelic.Application, configs ...any) fiber.Handler {
	// We can't directly convert []interface{} to v2's config type
	// So we call v2.Middleware without configs and let users call v2 directly if they need configs
	// This is a backwards compatibility layer - users should use v2 or v3 directly
	return v2.Middleware(app)
}

// Send creates and executes a custom segment with the given name.
func Send(c *fiber.Ctx, segmentName string) {
	v2.Send(c, segmentName)
}

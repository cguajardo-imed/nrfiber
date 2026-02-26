// Package nrfiber provides New Relic instrumentation for the Fiber v2 web framework.
// This package re-exports github.com/cguajardo-imed/nrfiber/v2 for backwards compatibility.
//
// For new projects, use the v2 or v3 packages directly:
//   - github.com/cguajardo-imed/nrfiber/v2 for Fiber v2
//   - github.com/cguajardo-imed/nrfiber/v3 for Fiber v3
package nrfiber

import (
	v2 "github.com/cguajardo-imed/nrfiber/v2"
	"github.com/gofiber/fiber/v2"
)

// ConfigNoticeErrorEnabled enables or disables error reporting to New Relic.
// This function re-exports from v2. For full functionality, import v2 directly.
func ConfigNoticeErrorEnabled(enabled bool) interface{} {
	return v2.ConfigNoticeErrorEnabled(enabled)
}

// ConfigStatusCodeIgnored specifies HTTP status codes that should not be reported as errors.
// This function re-exports from v2. For full functionality, import v2 directly.
func ConfigStatusCodeIgnored(statusCodes []int) interface{} {
	return v2.ConfigStatusCodeIgnored(statusCodes)
}

// ConfigCustomTransactionNameFunc provides a custom function to name transactions.
// This function re-exports from v2. For full functionality, import v2 directly.
func ConfigCustomTransactionNameFunc(fn func(*fiber.Ctx) string) interface{} {
	return v2.ConfigCustomTransactionNameFunc(fn)
}

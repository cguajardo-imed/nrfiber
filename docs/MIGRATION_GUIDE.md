# Migration Guide

This guide helps you migrate between different versions of nrfiber.

## Table of Contents

- [Migrating to v3.x (Fiber v3 Support)](#migrating-to-v3x-fiber-v3-support)
- [Migrating from v1.x to v2.x](#migrating-from-v1x-to-v2x)
- [Choosing Between v2 and v3](#choosing-between-v2-and-v3)
- [Common Migration Issues](#common-migration-issues)

---

## Migrating to v3.x (Fiber v3 Support)

Starting with nrfiber v3.0.0, the library supports both Fiber v2 and v3 through separate module paths.

### Overview of Changes

**Major Change**: The library now uses **separate modules** instead of build tags:
- **Fiber v3**: `github.com/cguajardo-imed/nrfiber/v3`
- **Fiber v2**: `github.com/cguajardo-imed/nrfiber/v2`

### Migration Steps

#### Step 1: Determine Your Fiber Version

Check which version of Fiber you're using:

```bash
go list -m github.com/gofiber/fiber/v2
# or
go list -m github.com/gofiber/fiber/v3
```

#### Step 2: Update Your Imports

**If using Fiber v3:**

```go
// OLD (v2.x or earlier)
import "github.com/erkanzileli/nrfiber"
import "github.com/gofiber/fiber/v2"

// NEW (v3.x)
import "github.com/cguajardo-imed/nrfiber/v3"
import "github.com/gofiber/fiber/v3"
```

**If staying with Fiber v2:**

```go
// OLD (v2.x or earlier)
import "github.com/erkanzileli/nrfiber"
import "github.com/gofiber/fiber/v2"

// NEW (v3.x)
import "github.com/cguajardo-imed/nrfiber/v2"
import "github.com/gofiber/fiber/v2"
```

#### Step 3: Update Function Signatures (Fiber v3 Only)

If you're migrating to Fiber v3, update your handler signatures:

**OLD (Fiber v2):**
```go
app.Get("/route", func(c *fiber.Ctx) error {
    txn := nrfiber.FromContext(c)
    // ...
})
```

**NEW (Fiber v3):**
```go
app.Get("/route", func(c fiber.Ctx) error {
    txn := nrfiber.FromContext(c)
    // ...
})
```

Notice the context parameter changed from `*fiber.Ctx` (pointer) to `fiber.Ctx` (interface).

#### Step 4: Update Custom Transaction Name Functions

If you're using custom transaction naming, update the function signature:

**For Fiber v3:**
```go
// OLD
nrfiber.ConfigCustomTransactionNameFunc(func(c *fiber.Ctx) string {
    return c.Method() + " " + c.Path()
})

// NEW
nrfiber.ConfigCustomTransactionNameFunc(func(c fiber.Ctx) string {
    return c.Method() + " " + c.Path()
})
```

**For Fiber v2:**
```go
// Signature stays the same (pointer)
nrfiber.ConfigCustomTransactionNameFunc(func(c *fiber.Ctx) string {
    return c.Method() + " " + c.Path()
})
```

#### Step 5: Update Dependencies

```bash
# For Fiber v3
go get -u github.com/cguajardo-imed/nrfiber/v3
go get -u github.com/gofiber/fiber/v3

# For Fiber v2
go get -u github.com/cguajardo-imed/nrfiber/v2
go get -u github.com/gofiber/fiber/v2
```

#### Step 6: Update Body Parsing (Fiber v3 Only)

If you're migrating to Fiber v3, update body parsing code:

**OLD (Fiber v2):**
```go
var req RequestBody
if err := c.BodyParser(&req); err != nil {
    return err
}
```

**NEW (Fiber v3):**
```go
var req RequestBody
if err := c.Bind().JSON(&req); err != nil {
    return err
}
```

### What Stays the Same

The following API remains unchanged:

1. **Middleware Registration:**
   ```go
   app.Use(nrfiber.Middleware(nrApp))
   ```

2. **Configuration Options:**
   ```go
   nrfiber.ConfigNoticeErrorEnabled(true)
   nrfiber.ConfigStatusCodeIgnored([]int{404})
   ```

3. **Transaction Retrieval:**
   ```go
   txn := nrfiber.FromContext(c)
   ```

4. **Segment Creation:**
   ```go
   segment := txn.StartSegment("SegmentName")
   defer segment.End()
   ```

5. **Helper Function:**
   ```go
   nrfiber.Send(c, "SegmentName")
   ```

---

## Migrating from v1.x to v2.x

If you're upgrading from an older version (v1.x), here are the key changes:

### Repository and Import Path Change

**OLD (v1.x):**
```go
import "github.com/erkanzileli/nrfiber"
```

**NEW (v2.x/v3.x):**
```go
// For Fiber v2
import "github.com/cguajardo-imed/nrfiber/v2"

// For Fiber v3
import "github.com/cguajardo-imed/nrfiber/v3"
```

### Error Reporting

Error reporting is now opt-in instead of always enabled:

**OLD (v1.x - always enabled):**
```go
app.Use(nrfiber.Middleware(nrApp))
// Errors were automatically reported
```

**NEW (v2.x/v3.x - opt-in):**
```go
app.Use(nrfiber.Middleware(nrApp, 
    nrfiber.ConfigNoticeErrorEnabled(true), // Explicitly enable
))
```

### Status Code Filtering

The configuration for ignoring status codes has been improved:

**OLD (v1.x):**
```go
// May have had different or no configuration
```

**NEW (v2.x/v3.x):**
```go
app.Use(nrfiber.Middleware(nrApp,
    nrfiber.ConfigNoticeErrorEnabled(true),
    nrfiber.ConfigStatusCodeIgnored([]int{404, 401}),
))
```

---

## Choosing Between v2 and v3

### Use nrfiber v3 (Fiber v3) if:

- ✅ You're starting a new project
- ✅ You want the latest Fiber features
- ✅ You're comfortable with Fiber v3's interface-based context
- ✅ Your dependencies support Fiber v3

### Use nrfiber v2 (Fiber v2) if:

- ✅ You have an existing Fiber v2 application
- ✅ Your dependencies aren't compatible with Fiber v3 yet
- ✅ You prefer pointer-based context handling
- ✅ You want maximum stability

### Key Differences Summary

| Feature | nrfiber v2 (Fiber v2) | nrfiber v3 (Fiber v3) |
|---------|----------------------|----------------------|
| Import Path | `nrfiber/v2` | `nrfiber/v3` |
| Context Type | `*fiber.Ctx` (pointer) | `fiber.Ctx` (interface) |
| Context Storage | `SetUserContext()` | `SetContext()` |
| Body Parsing | `BodyParser()` | `Bind().JSON()` |
| API Compatibility | Identical | Identical |
| Performance | Fast | Fast |

---

## Common Migration Issues

### Issue 1: Import Errors After Update

**Problem:**
```
cannot find package "github.com/erkanzileli/nrfiber"
```

**Solution:**
Update your import path to the new repository and version:
```go
import "github.com/cguajardo-imed/nrfiber/v3"  // for Fiber v3
// or
import "github.com/cguajardo-imed/nrfiber/v2"  // for Fiber v2
```

### Issue 2: Type Mismatch with Context

**Problem (Fiber v3):**
```
cannot use func literal (type func(*fiber.Ctx) error) as type fiber.Handler
```

**Solution:**
Remove the pointer from the context parameter:
```go
// OLD
func(c *fiber.Ctx) error {

// NEW
func(c fiber.Ctx) error {
```

### Issue 3: Nil Transaction Panics

**Problem:**
```
panic: runtime error: invalid memory address or nil pointer dereference
```

**Solution:**
Always check if transaction is nil before using it:
```go
txn := nrfiber.FromContext(c)
if txn == nil {
    // Handle case where transaction doesn't exist
    return c.SendString("No transaction available")
}

segment := txn.StartSegment("MySegment")
defer segment.End()
```

### Issue 4: Middleware Not Instrumenting Requests

**Problem:**
No transactions appear in New Relic dashboard.

**Solution:**
Ensure middleware is registered BEFORE your routes:
```go
app := fiber.New()

// Middleware MUST come first
app.Use(nrfiber.Middleware(nrApp))

// Then define routes
app.Get("/route", handler)
```

### Issue 5: Errors Not Appearing in New Relic

**Problem:**
Errors aren't being reported to New Relic.

**Solution:**
Enable error reporting explicitly:
```go
app.Use(nrfiber.Middleware(nrApp,
    nrfiber.ConfigNoticeErrorEnabled(true),
))
```

### Issue 6: Duplicate Errors in New Relic

**Problem:**
Seeing duplicate error entries for the same request.

**Solution:**
Configure New Relic to ignore HTTP status codes and/or use `ConfigStatusCodeIgnored`:
```go
app.Use(nrfiber.Middleware(nrApp,
    nrfiber.ConfigNoticeErrorEnabled(true),
    nrfiber.ConfigStatusCodeIgnored([]int{400, 404}),
))
```

See [Notice Custom Errors Guide](notice-custom-errors.md) for more details.

### Issue 7: Build Errors with Mixed Versions

**Problem:**
```
module github.com/gofiber/fiber/v3 found, but does not contain package github.com/gofiber/fiber/v3
```

**Solution:**
Ensure all Fiber-related dependencies use the same major version:
```bash
# Check your dependencies
go list -m all | grep fiber

# Update to consistent versions
go get -u github.com/gofiber/fiber/v3@latest
go get -u github.com/cguajardo-imed/nrfiber/v3@latest
```

---

## Testing Your Migration

After migrating, verify everything works correctly:

### 1. Run Tests

```bash
go test ./...
```

### 2. Check New Relic Dashboard

Start your application and make some requests:

```bash
# Start your app
go run main.go

# Make test requests
curl http://localhost:3000/your-route
```

Then verify in New Relic:
- Transactions appear in the APM dashboard
- Transaction names are correct
- Errors are reported (if enabled)
- Segments appear for custom instrumentation

### 3. Verify Error Reporting

Test error handling:

```bash
# Trigger an error
curl http://localhost:3000/error-route

# Check New Relic Errors page
```

### 4. Load Test

Run a load test to ensure performance is acceptable:

```bash
# Using Apache Bench
ab -n 1000 -c 10 http://localhost:3000/your-route

# Check New Relic for performance metrics
```

---

## Rollback Plan

If you encounter issues during migration, here's how to rollback:

### Quick Rollback

```bash
# Revert to previous version in go.mod
go get github.com/erkanzileli/nrfiber@v1.x.x

# Or revert your git changes
git checkout HEAD -- go.mod go.sum
go mod download
```

### Gradual Migration

If you have a large application, consider a gradual migration:

1. Create a new branch for migration work
2. Migrate one module/package at a time
3. Test thoroughly after each change
4. Merge when fully tested

---

## Getting Help

If you encounter issues not covered in this guide:

1. Check the [examples directory](../examples/) for working code
2. Review the [main documentation](../README.md)
3. Search [existing issues](https://github.com/cguajardo-imed/nrfiber/issues)
4. Open a new issue with:
   - Your Go version
   - Your Fiber version (v2 or v3)
   - Your nrfiber version
   - A minimal reproducible example
   - Error messages and stack traces

---

## Additional Resources

- [Fiber v3 Migration Guide](https://docs.gofiber.io/guide/migration)
- [New Relic Go Agent Documentation](https://docs.newrelic.com/docs/apm/agents/go-agent/)
- [nrfiber Examples](../examples/)
- [Notice Custom Errors Guide](notice-custom-errors.md)
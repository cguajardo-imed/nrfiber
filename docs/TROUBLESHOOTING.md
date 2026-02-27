# Troubleshooting Guide

This guide helps you diagnose and resolve common issues when using nrfiber with Fiber v2 or v3.

## Table of Contents

- [Installation Issues](#installation-issues)
- [Middleware Issues](#middleware-issues)
- [Transaction Issues](#transaction-issues)
- [Error Reporting Issues](#error-reporting-issues)
- [Performance Issues](#performance-issues)
- [Integration Issues](#integration-issues)
- [New Relic Dashboard Issues](#new-relic-dashboard-issues)
- [Debug Mode](#debug-mode)

---

## Installation Issues

### Cannot Find Package

**Symptoms:**
```
cannot find package "github.com/cguajardo-imed/nrfiber/v3"
```

**Causes:**
- Wrong import path
- Module not downloaded
- Network issues

**Solutions:**

1. **Verify you're using the correct import path:**
   ```go
   // For Fiber v3
   import "github.com/cguajardo-imed/nrfiber/v3"
   
   // For Fiber v2
   import "github.com/cguajardo-imed/nrfiber/v2"
   ```

2. **Download the module:**
   ```bash
   go get github.com/cguajardo-imed/nrfiber/v3
   # or
   go get github.com/cguajardo-imed/nrfiber/v2
   ```

3. **Clear cache and retry:**
   ```bash
   go clean -modcache
   go get github.com/cguajardo-imed/nrfiber/v3
   go mod tidy
   ```

### Version Conflict

**Symptoms:**
```
module github.com/gofiber/fiber/v3 found, but does not contain package
```

**Cause:**
Mixed Fiber v2 and v3 dependencies in the same project.

**Solution:**

1. **Check all fiber dependencies:**
   ```bash
   go list -m all | grep fiber
   ```

2. **Ensure consistency:**
   ```bash
   # If using Fiber v3
   go get github.com/gofiber/fiber/v3@latest
   go get github.com/cguajardo-imed/nrfiber/v3@latest
   
   # If using Fiber v2
   go get github.com/gofiber/fiber/v2@latest
   go get github.com/cguajardo-imed/nrfiber/v2@latest
   ```

3. **Remove conflicting dependencies:**
   ```bash
   go mod tidy
   ```

### Build Errors After Update

**Symptoms:**
```
undefined: nrfiber.Middleware
```

**Cause:**
Old cached build artifacts.

**Solution:**

```bash
# Clean build cache
go clean -cache -modcache -i -r

# Re-download dependencies
go mod download

# Rebuild
go build
```

---

## Middleware Issues

### Middleware Not Executing

**Symptoms:**
- No transactions appear in New Relic
- Middleware seems to be skipped

**Solution:**

1. **Ensure middleware is registered BEFORE routes:**
   ```go
   app := fiber.New()
   
   // CORRECT: Middleware first
   app.Use(nrfiber.Middleware(nrApp))
   
   // Then routes
   app.Get("/route", handler)
   ```

2. **Verify New Relic app is initialized:**
   ```go
   nrApp, err := newrelic.NewApplication(
       newrelic.ConfigAppName("my-app"),
       newrelic.ConfigLicense(licenseKey),
       newrelic.ConfigEnabled(true),
   )
   if err != nil {
       log.Fatal(err) // Don't ignore this error!
   }
   ```

3. **Check middleware is not skipped by other middleware:**
   ```go
   // Bad: This might prevent nrfiber from executing
   app.Use(func(c fiber.Ctx) error {
       if c.Path() == "/health" {
           return c.SendString("OK") // Doesn't call c.Next()
       }
       return c.Next()
   })
   
   // Good: Always call c.Next() when appropriate
   app.Use(func(c fiber.Ctx) error {
       if c.Path() == "/health" {
           return c.Next() // Let other middleware run
       }
       return c.Next()
   })
   ```

### Middleware Causes Panic

**Symptoms:**
```
panic: runtime error: invalid memory address or nil pointer dereference
```

**Cause:**
Nil New Relic application passed to middleware.

**Solution:**

1. **Always check for errors when creating New Relic app:**
   ```go
   nrApp, err := newrelic.NewApplication(...)
   if err != nil {
       log.Fatal("Failed to create New Relic app:", err)
   }
   
   // Only use nrApp if err is nil
   app.Use(nrfiber.Middleware(nrApp))
   ```

2. **Add nil check if needed:**
   ```go
   if nrApp != nil {
       app.Use(nrfiber.Middleware(nrApp))
   }
   ```

### Routes Not Being Instrumented

**Symptoms:**
- Some routes appear in New Relic, others don't
- Inconsistent transaction recording

**Solution:**

1. **Check route registration order:**
   ```go
   // Middleware must come before ALL routes
   app.Use(nrfiber.Middleware(nrApp))
   
   // These routes will be instrumented
   app.Get("/api/users", getUsersHandler)
   app.Post("/api/users", createUserHandler)
   ```

2. **Check for route groups:**
   ```go
   app.Use(nrfiber.Middleware(nrApp))
   
   // Route groups also need middleware registered first
   api := app.Group("/api")
   api.Get("/users", getUsersHandler) // Will be instrumented
   ```

3. **Verify middleware applies to all routes:**
   ```go
   // If using route-specific middleware, nrfiber must come first
   app.Use(nrfiber.Middleware(nrApp))
   app.Get("/route", otherMiddleware, handler) // OK
   ```

---

## Transaction Issues

### Cannot Get Transaction from Context

**Symptoms:**
```go
txn := nrfiber.FromContext(c)
// txn is nil
```

**Causes:**
- Middleware not registered
- Middleware registered after routes
- Context not properly passed

**Solution:**

1. **Always check for nil before using transaction:**
   ```go
   txn := nrfiber.FromContext(c)
   if txn == nil {
       log.Println("Warning: No New Relic transaction found")
       // Handle gracefully
       return c.SendString("OK")
   }
   
   // Safe to use txn now
   segment := txn.StartSegment("MySegment")
   defer segment.End()
   ```

2. **Verify middleware registration:**
   ```go
   app.Use(nrfiber.Middleware(nrApp)) // Must be before routes
   ```

3. **Check New Relic is enabled:**
   ```go
   nrApp, err := newrelic.NewApplication(
       newrelic.ConfigEnabled(true), // Must be true!
       // ...
   )
   ```

### Segments Not Appearing

**Symptoms:**
- Transactions appear but custom segments don't
- Segments show 0ms duration

**Solution:**

1. **Ensure segment is properly closed:**
   ```go
   segment := txn.StartSegment("MyOperation")
   defer segment.End() // Always defer End()
   
   // Your operation here
   ```

2. **Check segment timing:**
   ```go
   // Bad: Segment ends immediately
   segment := txn.StartSegment("FastOp")
   segment.End()
   
   // Good: Segment wraps actual work
   segment := txn.StartSegment("SlowOp")
   defer segment.End()
   time.Sleep(100 * time.Millisecond) // Actual work
   ```

3. **Use meaningful segment names:**
   ```go
   // Bad: Generic name
   segment := txn.StartSegment("Operation")
   
   // Good: Descriptive name
   segment := txn.StartSegment("Database - Query Users")
   ```

### Transaction Names Are Generic

**Symptoms:**
- All transactions named "GET /route"
- Cannot distinguish between different endpoints

**Solution:**

Use custom transaction naming:

**For Fiber v3:**
```go
customNameFunc := func(c fiber.Ctx) string {
    // Include route parameters
    return fmt.Sprintf("%s %s", c.Method(), c.Route().Path)
}

app.Use(nrfiber.Middleware(nrApp,
    nrfiber.ConfigCustomTransactionNameFunc(customNameFunc),
))
```

**For Fiber v2:**
```go
customNameFunc := func(c *fiber.Ctx) string {
    // Include route parameters
    return fmt.Sprintf("%s %s", c.Method(), c.Route().Path)
}

app.Use(nrfiber.Middleware(nrApp,
    nrfiber.ConfigCustomTransactionNameFunc(customNameFunc),
))
```

---

## Error Reporting Issues

### Errors Not Being Reported

**Symptoms:**
- Application returns errors but they don't appear in New Relic
- Error count is 0 in dashboard

**Solutions:**

1. **Enable error reporting:**
   ```go
   app.Use(nrfiber.Middleware(nrApp,
       nrfiber.ConfigNoticeErrorEnabled(true), // Must be enabled!
   ))
   ```

2. **Ensure errors are returned from handlers:**
   ```go
   // Bad: Error swallowed
   app.Get("/route", func(c fiber.Ctx) error {
       err := doSomething()
       if err != nil {
           log.Println(err) // Just logging doesn't report to New Relic
           return c.SendString("Error occurred")
       }
       return c.SendString("OK")
   })
   
   // Good: Error returned
   app.Get("/route", func(c fiber.Ctx) error {
       err := doSomething()
       if err != nil {
           return err // nrfiber will report this
       }
       return c.SendString("OK")
   })
   ```

3. **Manually report errors if needed:**
   ```go
   app.Get("/route", func(c fiber.Ctx) error {
       txn := nrfiber.FromContext(c)
       err := doSomething()
       if err != nil {
           if txn != nil {
               txn.NoticeError(err) // Manually report
           }
           return c.Status(500).JSON(fiber.Map{"error": "Internal error"})
       }
       return c.SendString("OK")
   })
   ```

### Duplicate Errors

**Symptoms:**
- Same error appears multiple times in New Relic
- Both custom error and HTTP status error

**Cause:**
Both nrfiber and New Relic agent reporting the same error.

**Solution:**

1. **Ignore HTTP status codes in nrfiber:**
   ```go
   app.Use(nrfiber.Middleware(nrApp,
       nrfiber.ConfigNoticeErrorEnabled(true),
       nrfiber.ConfigStatusCodeIgnored([]int{400, 401, 404}),
   ))
   ```

2. **Configure New Relic to ignore status codes:**
   - Go to New Relic dashboard
   - Settings → Application → Server-side agent configuration
   - Error Collection → Ignore these HTTP status codes
   - Add ranges like `400-499`

See [Notice Custom Errors Guide](notice-custom-errors.md) for details.

### Specific Status Codes Not Being Ignored

**Symptoms:**
- Configured status codes still appear as errors

**Solution:**

1. **Verify configuration:**
   ```go
   app.Use(nrfiber.Middleware(nrApp,
       nrfiber.ConfigNoticeErrorEnabled(true),
       nrfiber.ConfigStatusCodeIgnored([]int{404, 401, 403}),
   ))
   ```

2. **Check you're returning the status code correctly:**
   ```go
   // This will be ignored if 404 is in ignore list
   return c.Status(404).SendString("Not Found")
   ```

3. **Verify status code is set before returning:**
   ```go
   // Bad: Status not set
   return errors.New("not found") // Will use default 500
   
   // Good: Status set explicitly
   c.Status(404)
   return errors.New("not found")
   ```

---

## Performance Issues

### High Memory Usage

**Symptoms:**
- Application memory increases over time
- Memory not being released

**Causes:**
- Transactions not being ended
- Segments not being closed
- Transaction accumulation

**Solutions:**

1. **Always defer segment closure:**
   ```go
   segment := txn.StartSegment("Operation")
   defer segment.End() // Ensures cleanup
   ```

2. **Check for transaction leaks:**
   ```go
   // Middleware automatically ends transactions
   // Don't manually end middleware transactions
   
   // For manual transactions:
   txn := app.StartTransaction("background-job")
   defer txn.End()
   ```

3. **Monitor goroutines:**
   ```go
   import "runtime"
   
   log.Printf("Goroutines: %d", runtime.NumGoroutine())
   ```

### Slow Request Processing

**Symptoms:**
- Requests take longer than expected
- New Relic overhead is high

**Solutions:**

1. **Profile your application:**
   ```bash
   go test -cpuprofile=cpu.prof -bench=.
   go tool pprof cpu.prof
   ```

2. **Check segment overhead:**
   ```go
   // Bad: Too many small segments
   for i := 0; i < 1000; i++ {
       seg := txn.StartSegment(fmt.Sprintf("Iteration%d", i))
       doSmallWork()
       seg.End()
   }
   
   // Good: One segment for loop
   seg := txn.StartSegment("ProcessLoop")
   for i := 0; i < 1000; i++ {
       doSmallWork()
   }
   seg.End()
   ```

3. **Disable New Relic in development:**
   ```go
   nrApp, err := newrelic.NewApplication(
       newrelic.ConfigEnabled(os.Getenv("ENV") == "production"),
       // ...
   )
   ```

### High CPU Usage

**Symptoms:**
- CPU usage increases significantly with New Relic enabled

**Solutions:**

1. **Reduce sampling rate:**
   ```go
   // Not directly configurable in nrfiber
   // Consider using New Relic's server-side sampling
   ```

2. **Optimize custom instrumentation:**
   ```go
   // Bad: Creating segments for every loop iteration
   for _, item := range items {
       seg := txn.StartSegment("ProcessItem")
       process(item)
       seg.End()
   }
   
   // Good: Batch processing under one segment
   seg := txn.StartSegment("ProcessItems")
   for _, item := range items {
       process(item)
   }
   seg.End()
   ```

3. **Use background goroutines wisely:**
   ```go
   // Ensure background work doesn't create new transactions
   go func() {
       // Don't start new transactions in background goroutines
       // unless necessary
   }()
   ```

---

## Integration Issues

### Fiber v2/v3 Type Mismatch

**Symptoms:**
```
cannot use func literal (type func(*fiber.Ctx) error) as type fiber.Handler
```

**Cause:**
Using wrong nrfiber version for your Fiber version.

**Solution:**

**For Fiber v3:**
```go
import "github.com/cguajardo-imed/nrfiber/v3"
import "github.com/gofiber/fiber/v3"

// Context is interface (not pointer)
app.Get("/route", func(c fiber.Ctx) error {
    return c.SendString("OK")
})
```

**For Fiber v2:**
```go
import "github.com/cguajardo-imed/nrfiber/v2"
import "github.com/gofiber/fiber/v2"

// Context is pointer
app.Get("/route", func(c *fiber.Ctx) error {
    return c.SendString("OK")
})
```

### Other Middleware Conflicts

**Symptoms:**
- nrfiber works alone but fails with other middleware
- Transactions missing when using middleware stacks

**Solutions:**

1. **Register nrfiber first:**
   ```go
   app.Use(nrfiber.Middleware(nrApp))
   app.Use(logger.New())
   app.Use(cors.New())
   // Other middleware...
   ```

2. **Ensure all middleware calls Next():**
   ```go
   app.Use(func(c fiber.Ctx) error {
       // Do work...
       return c.Next() // Must call Next()
   })
   ```

3. **Check for context overwrites:**
   ```go
   // Bad: This might lose New Relic context
   c.SetUserContext(context.Background())
   
   // Good: Preserve existing context
   ctx := c.Context()
   // Use ctx instead of creating new one
   ```

### Database Instrumentation Not Working

**Symptoms:**
- Database queries not appearing as segments

**Solution:**

Use New Relic's database instrumentation:

```go
import (
    "github.com/newrelic/go-agent/v3/integrations/nrpgx"
    "github.com/jackc/pgx/v4"
)

// For PostgreSQL
db, err := nrpgx.Open("postgres://...")

// Then in your handler
app.Get("/users", func(c fiber.Ctx) error {
    txn := nrfiber.FromContext(c)
    
    // Queries will be tracked automatically
    rows, err := db.QueryContext(
        newrelic.NewContext(context.Background(), txn),
        "SELECT * FROM users",
    )
    // ...
})
```

---

## New Relic Dashboard Issues

### No Data Appearing

**Symptoms:**
- Application registered but no data in dashboard
- Transactions not showing up

**Checklist:**

1. **Verify license key:**
   ```go
   nrApp, err := newrelic.NewApplication(
       newrelic.ConfigLicense("your-valid-license-key"),
       // ...
   )
   ```

2. **Check New Relic is enabled:**
   ```go
   newrelic.ConfigEnabled(true)
   ```

3. **Verify app name:**
   ```go
   newrelic.ConfigAppName("my-app-name")
   ```

4. **Check network connectivity:**
   ```bash
   # Test connection to New Relic
   curl -H "Api-Key: YOUR_LICENSE_KEY" \
        https://collector.newrelic.com/status/mongrel
   ```

5. **Enable debug logging:**
   ```go
   nrApp, err := newrelic.NewApplication(
       newrelic.ConfigDebugLogger(os.Stdout),
       // ...
   )
   ```

### Data Delayed

**Symptoms:**
- Transactions appear 5-10 minutes late

**Explanation:**
This is normal behavior. New Relic batches and processes data:
- Transaction data: ~1 minute delay
- Metrics: ~1-2 minute delay
- Errors: ~1 minute delay

**Solutions:**
- Wait a few minutes after making requests
- Use New Relic's real-time view if available
- Check data harvest interval (default: 60 seconds)

### Incorrect Transaction Counts

**Symptoms:**
- Transaction count doesn't match actual requests
- Some requests not being counted

**Causes:**
- Health check endpoints being counted
- Background jobs creating transactions
- Transaction sampling

**Solutions:**

1. **Ignore health check endpoints:**
   ```go
   app.Get("/health", func(c fiber.Ctx) error {
       return c.SendString("OK") // Not instrumented if before middleware
   })
   
   app.Use(nrfiber.Middleware(nrApp)) // Middleware after health check
   
   // Or use route-specific middleware
   app.Get("/api/*", nrfiber.Middleware(nrApp), apiHandler)
   ```

2. **Check transaction sampling:**
   - New Relic may sample high-volume applications
   - Check account settings for sampling configuration

---

## Debug Mode

### Enable New Relic Debug Logging

```go
import (
    "os"
    "github.com/newrelic/go-agent/v3/newrelic"
)

nrApp, err := newrelic.NewApplication(
    newrelic.ConfigAppName("my-app"),
    newrelic.ConfigLicense(licenseKey),
    newrelic.ConfigEnabled(true),
    newrelic.ConfigDebugLogger(os.Stdout), // Enable debug output
)
```

### Enable Detailed nrfiber Logging

```go
app.Use(func(c fiber.Ctx) error {
    txn := nrfiber.FromContext(c)
    if txn == nil {
        log.Printf("[nrfiber] No transaction for %s %s", c.Method(), c.Path())
    } else {
        log.Printf("[nrfiber] Transaction created for %s %s", c.Method(), c.Path())
    }
    return c.Next()
})

app.Use(nrfiber.Middleware(nrApp))
```

### Test Transaction Creation

```go
package main

import (
    "log"
    "github.com/cguajardo-imed/nrfiber/v3"
    "github.com/gofiber/fiber/v3"
    "github.com/newrelic/go-agent/v3/newrelic"
)

func main() {
    app := fiber.New()
    
    nrApp, err := newrelic.NewApplication(
        newrelic.ConfigAppName("test"),
        newrelic.ConfigLicense("license"),
        newrelic.ConfigEnabled(true),
        newrelic.ConfigDebugLogger(log.Default().Writer()),
    )
    if err != nil {
        log.Fatal(err)
    }
    
    app.Use(nrfiber.Middleware(nrApp))
    
    app.Get("/test", func(c fiber.Ctx) error {
        txn := nrfiber.FromContext(c)
        if txn == nil {
            log.Println("ERROR: No transaction!")
            return c.SendString("No transaction")
        }
        
        log.Println("SUCCESS: Transaction found")
        
        segment := txn.StartSegment("TestSegment")
        defer segment.End()
        
        return c.SendString("OK")
    })
    
    log.Fatal(app.Listen(":3000"))
}
```

Then test:
```bash
curl http://localhost:3000/test
# Check logs for transaction creation
```

---

## Getting Help

If you've tried the solutions above and still have issues:

1. **Check Examples:**
   - [Fiber v3 Basic Example](../examples/fiber-v3-basic/)
   - [Fiber v3 Advanced Example](../examples/fiber-v3-advanced/)
   - [Fiber v2 Basic Example](../examples/fiber-v2-basic/)

2. **Enable Debug Logging:**
   - Add `newrelic.ConfigDebugLogger(os.Stdout)` to your config
   - Check logs for error messages

3. **Create Minimal Reproduction:**
   ```go
   // Simplest possible example
   package main
   
   import (
       "log"
       "os"
       "github.com/cguajardo-imed/nrfiber/v3"
       "github.com/gofiber/fiber/v3"
       "github.com/newrelic/go-agent/v3/newrelic"
   )
   
   func main() {
       app := fiber.New()
       nrApp, _ := newrelic.NewApplication(
           newrelic.ConfigAppName("test"),
           newrelic.ConfigLicense(os.Getenv("NEW_RELIC_LICENSE_KEY")),
           newrelic.ConfigEnabled(true),
       )
       app.Use(nrfiber.Middleware(nrApp))
       app.Get("/", func(c fiber.Ctx) error {
           return c.SendString("OK")
       })
       app.Listen(":3000")
   }
   ```

4. **Open an Issue:**
   - Go to [GitHub Issues](https://github.com/cguajardo-imed/nrfiber/issues)
   - Include:
     - Go version: `go version`
     - Fiber version: `go list -m github.com/gofiber/fiber/v3`
     - nrfiber version: `go list -m github.com/cguajardo-imed/nrfiber/v3`
     - Minimal code example
     - Error messages and logs
     - Expected vs actual behavior

---

## Additional Resources

- [Main Documentation](../README.md)
- [Migration Guide](MIGRATION_GUIDE.md)
- [Notice Custom Errors](notice-custom-errors.md)
- [IDE Setup](IDE_SETUP.md)
- [New Relic Go Agent Documentation](https://docs.newrelic.com/docs/apm/agents/go-agent/)
- [Fiber Documentation](https://docs.gofiber.io/)
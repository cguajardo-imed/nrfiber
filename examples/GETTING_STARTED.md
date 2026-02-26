# Getting Started with nrfiber Examples

This guide will help you quickly get started with the nrfiber examples for both Fiber v2 and v3.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Understanding the Examples](#understanding-the-examples)
- [Configuration](#configuration)
- [Testing Your Application](#testing-your-application)
- [Viewing Metrics in New Relic](#viewing-metrics-in-new-relic)
- [Common Patterns](#common-patterns)
- [Troubleshooting](#troubleshooting)

## Prerequisites

Before you begin, make sure you have:

- **Go 1.25.0 or higher** installed
- A **New Relic account** (optional but recommended)
  - Sign up for free at [newrelic.com](https://newrelic.com)
  - Get your license key from [New Relic API Keys](https://one.newrelic.com/launcher/api-keys-ui.api-keys-launcher)

## Quick Start

### Step 1: Choose Your Fiber Version

Decide whether you want to use Fiber v2 or v3:

- **Fiber v3**: Latest version with improved API (recommended for new projects)
- **Fiber v2**: Stable version for existing projects

### Step 2: Navigate to Example Directory

```bash
# For Fiber v3
cd examples/fiber-v3-basic

# For Fiber v2
cd examples/fiber-v2-basic
```

### Step 3: Configure Environment

Copy the example environment file:

```bash
cp ../.env.example .env
```

Edit `.env` and add your New Relic license key:

```env
NEW_RELIC_LICENSE_KEY=your-actual-license-key-here
PORT=3000
```

**Note:** If you don't have a New Relic license key, the examples will still run but won't send telemetry data.

### Step 4: Install Dependencies

```bash
go mod download
```

### Step 5: Run the Example

**For Fiber v3:**
```bash
go run main.go
```

**For Fiber v2:**
```bash
go run main.go
```

**Note:** No build tags required - the version is determined by your import path.

### Step 6: Test the Application

Open another terminal and test the endpoints:

```bash
# Health check
curl http://localhost:3000/health

# Get user
curl http://localhost:3000/users/123

# Search
curl "http://localhost:3000/search?q=golang&limit=5"

# Create user
curl -X POST http://localhost:3000/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Jane Doe","email":"jane@example.com"}'
```

## Understanding the Examples

### Basic Examples

The basic examples demonstrate:

1. **Middleware Setup**: How to integrate nrfiber with your Fiber application
2. **Custom Segments**: Tracking specific operations (database queries, external calls)
3. **Error Handling**: How errors are reported to New Relic
4. **Route Parameters**: Working with dynamic routes
5. **Query Parameters**: Handling search and filtering

**File Structure:**
```
fiber-v3-basic/
├── main.go          # Main application code
├── go.mod           # Go module definition
├── README.md        # Detailed documentation
└── .env             # Environment variables (create from .env.example)
```

### Advanced Examples

The advanced examples show:

1. **Custom Error Types**: Creating and handling custom errors
2. **Transaction Attributes**: Adding custom metadata to transactions
3. **Status Code Filtering**: Ignoring specific HTTP codes (e.g., 404s)
4. **Custom Transaction Names**: Better grouping in New Relic
5. **Multiple Segments**: Complex operations with multiple tracking points

## Configuration

### Basic Configuration

Minimal setup:

```go
nrApp, _ := newrelic.NewApplication(
    newrelic.ConfigAppName("my-app"),
    newrelic.ConfigLicense(os.Getenv("NEW_RELIC_LICENSE_KEY")),
)

app.Use(nrfiber.Middleware(nrApp))
```

### Advanced Configuration

With custom options:

```go
app.Use(nrfiber.Middleware(
    nrApp,
    // Enable error reporting to New Relic
    nrfiber.ConfigNoticeErrorEnabled(true),
    
    // Don't report 404 and 401 as errors
    nrfiber.ConfigStatusCodeIgnored([]int{404, 401}),
    
    // Custom transaction naming for better grouping
    nrfiber.ConfigCustomTransactionNameFunc(func(c fiber.Ctx) string {
        return fmt.Sprintf("%s %s", c.Method(), c.Route().Path)
    }),
))
```

## Testing Your Application

### Manual Testing with curl

```bash
# GET request
curl http://localhost:3000/users/123

# POST request
curl -X POST http://localhost:3000/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John","email":"john@example.com"}'

# With query parameters
curl "http://localhost:3000/search?q=test&limit=10"

# Test error handling
curl http://localhost:3000/error
```

### Load Testing

Use tools like `ab` (Apache Bench) or `hey` to generate load:

```bash
# Install hey
go install github.com/rakyll/hey@latest

# Run load test
hey -n 1000 -c 10 http://localhost:3000/users/123
```

This will help you see meaningful data in New Relic.

## Viewing Metrics in New Relic

### Step 1: Access New Relic One

1. Go to [one.newrelic.com](https://one.newrelic.com)
2. Log in with your credentials

### Step 2: Find Your Application

1. Click on **APM & Services** in the left sidebar
2. Find your application (e.g., "fiber-v3-basic-example")

### Step 3: Explore Metrics

**Transactions Tab:**
- View all your HTTP endpoints
- See response times and throughput
- Identify slow endpoints

**Databases Tab:**
- View custom segments (your database operations)
- See query performance
- Track slow queries

**Errors Tab:**
- See error rates and types
- View error details and stack traces
- Track error trends over time

**Service Maps:**
- Visualize your application architecture
- See dependencies and connections
- Identify bottlenecks

## Common Patterns

### Pattern 1: Simple Segment

Track a single operation:

```go
app.Get("/users/:id", func(c fiber.Ctx) error {
    txn := nrfiber.FromContext(c)
    
    segment := txn.StartSegment("Database - Get User")
    user := getUser(c.Params("id"))
    segment.End()
    
    return c.JSON(user)
})
```

### Pattern 2: Multiple Segments

Track multiple operations in sequence:

```go
app.Post("/orders", func(c fiber.Ctx) error {
    txn := nrfiber.FromContext(c)
    
    // Validate user
    seg1 := txn.StartSegment("Database - Validate User")
    validateUser(userID)
    seg1.End()
    
    // Process payment
    seg2 := txn.StartSegment("External - Payment Gateway")
    processPayment(amount)
    seg2.End()
    
    // Create order
    seg3 := txn.StartSegment("Database - Create Order")
    order := createOrder(data)
    seg3.End()
    
    return c.JSON(order)
})
```

### Pattern 3: Using Defer

Automatically end segment when function returns:

```go
app.Get("/search", func(c fiber.Ctx) error {
    txn := nrfiber.FromContext(c)
    
    segment := txn.StartSegment("Search Operation")
    defer segment.End()
    
    results := performSearch(c.Query("q"))
    return c.JSON(results)
})
```

### Pattern 4: Adding Custom Attributes

Enrich transactions with metadata:

```go
app.Get("/products/:category", func(c fiber.Ctx) error {
    txn := nrfiber.FromContext(c)
    
    // Add custom attributes
    txn.AddAttribute("category", c.Params("category"))
    txn.AddAttribute("user_id", getUserID(c))
    txn.AddAttribute("region", getRegion(c))
    
    products := getProducts(c.Params("category"))
    return c.JSON(products)
})
```

## Troubleshooting

### Issue: "NEW_RELIC_LICENSE_KEY not set"

**Solution:** This is just a warning. The app will run in disabled mode without sending data to New Relic.

To fix, set the environment variable:
```bash
export NEW_RELIC_LICENSE_KEY="your-key-here"
```

### Issue: Wrong import path for Fiber version

**Problem:** Using v3 import with Fiber v2 or vice versa.

**Solution:** Make sure your imports match your Fiber version:

For Fiber v2:
```go
import "github.com/cguajardo-imed/nrfiber/v2"
import "github.com/gofiber/fiber/v2"
```

For Fiber v3:
```go
import "github.com/cguajardo-imed/nrfiber/v3"
import "github.com/gofiber/fiber/v3"
```

### Issue: No data in New Relic

**Possible causes:**
1. License key not set or invalid
2. New Relic agent disabled in code
3. Firewall blocking outbound connections
4. Not enough traffic (make some requests!)

**Solutions:**
1. Verify your license key
2. Check `newrelic.ConfigEnabled(true)` is set
3. Check firewall rules
4. Generate traffic with curl or load testing tools

### Issue: "module not found"

**Solution:** Run from the example directory and download dependencies:
```bash
cd examples/fiber-v3-basic
go mod download
```

### Issue: Port already in use

**Solution:** Change the port:
```bash
PORT=8080 go run main.go
```

## Next Steps

Once you're comfortable with the basic examples:

1. **Try the Advanced Example**: Explore custom error handling and transaction attributes
2. **Integrate with Your Project**: Copy patterns into your own application
3. **Customize Transaction Names**: Better organize your metrics
4. **Add Custom Attributes**: Track business-specific data
5. **Set Up Alerts**: Get notified about performance issues
6. **Create Dashboards**: Visualize your application metrics

## Learn More

- [nrfiber Main Documentation](../README.md)
- [AGENTS.md](../AGENTS.md) - Project rules and context for AI agents
- [Fiber v3 Documentation](https://docs.gofiber.io/)
- [Fiber v2 Documentation](https://docs.gofiber.io/v2.x/)
- [New Relic Go Agent Guide](https://docs.newrelic.com/docs/apm/agents/go-agent/)
- [New Relic APM Best Practices](https://docs.newrelic.com/docs/new-relic-solutions/best-practices-guides/full-stack-observability/apm-best-practices-guide/)

## Getting Help

If you run into issues:

1. Check the [README](./README.md) in each example directory
2. Review the [main documentation](../README.md)
3. Look at [closed issues](https://github.com/cguajardo-imed/nrfiber/issues?q=is%3Aissue+is%3Aclosed)
4. Open a [new issue](https://github.com/cguajardo-imed/nrfiber/issues/new)

Happy monitoring! 🚀
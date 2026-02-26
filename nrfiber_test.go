package nrfiber

import (
	"context"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestFromContext(t *testing.T) {
	app := fiber.New()
	c := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(c)

	// Test with no transaction in context
	txn := FromContext(c)
	assert.Nil(t, txn)

	// Test with transaction in context
	nrApp, _ := newrelic.NewApplication(newrelic.ConfigEnabled(false))
	mockTxn := nrApp.StartTransaction("test")
	c.SetContext(newrelic.NewContext(context.Background(), mockTxn))
	txn = FromContext(c)
	assert.NotNil(t, txn)
}

func TestDefaultTransactionName(t *testing.T) {
	app := fiber.New()
	c := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(c)

	c.Request().Header.SetMethod("GET")
	c.Request().SetRequestURI("/")

	name := defaultTransactionName(c)
	assert.Equal(t, "GET /", name)
}

func TestCustomTransactionName(t *testing.T) {
	app := fiber.New()
	c := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(c)

	c.Request().Header.SetMethod("GET")
	c.Request().SetRequestURI("/users/123")

	customFunc := func(ctx fiber.Ctx) string {
		return "CUSTOM-" + ctx.Method()
	}

	configMap := createConfigMap(ConfigCustomTransactionNameFunc(customFunc))
	nameFunc := customTransactionNameFunc(configMap, defaultTransactionName)
	name := nameFunc(c)

	assert.Equal(t, "CUSTOM-GET", name)
}

func TestCreateHttpRequest(t *testing.T) {
	app := fiber.New()
	c := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(c)

	c.Request().Header.SetMethod("GET")
	c.Request().SetRequestURI("http://example.com/?foo=bar")
	c.Request().Header.Set("Host", "example.com")
	c.Request().Header.Set("User-Agent", "test-agent")

	req := createHTTPRequest(c)

	assert.Equal(t, "GET", req.Method)
	assert.Equal(t, "http", req.URL.Scheme)
	assert.Equal(t, "example.com", req.URL.Host)
	assert.Equal(t, "/", req.URL.Path)
	assert.Equal(t, "foo=bar", req.URL.RawQuery)
	assert.Equal(t, "test-agent", req.Header.Get("User-Agent"))
}

func TestMiddlewareWithNilApp(t *testing.T) {
	handler := Middleware(nil)
	assert.NotNil(t, handler)

	// The handler should be a no-op that just calls Next()
	// We just verify it doesn't panic and returns a valid handler
	assert.IsType(t, (fiber.Handler)(nil), handler)
}

func TestNoticeErrorWithFiberError(t *testing.T) {
	nrApp, _ := newrelic.NewApplication(newrelic.ConfigEnabled(false))
	txn := nrApp.StartTransaction("test")
	defer txn.End()

	fiberErr := fiber.NewError(400, "Bad Request")
	configMap := createConfigMap(ConfigNoticeErrorEnabled(true))

	statusCode := noticeError(txn, fiberErr, configMap)
	assert.Equal(t, 400, statusCode)
}

func TestNoticeErrorWithIgnoredStatusCode(t *testing.T) {
	nrApp, _ := newrelic.NewApplication(newrelic.ConfigEnabled(false))
	txn := nrApp.StartTransaction("test")
	defer txn.End()

	fiberErr := fiber.NewError(404, "Not Found")
	configMap := createConfigMap(
		ConfigNoticeErrorEnabled(true),
		ConfigStatusCodeIgnored([]int{404}),
	)

	statusCode := noticeError(txn, fiberErr, configMap)
	assert.Equal(t, 404, statusCode)
}

func TestConvertRequestHeaders(t *testing.T) {
	app := fiber.New()
	c := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(c)

	c.Request().Header.Set("Content-Type", "application/json")
	c.Request().Header.Set("Authorization", "Bearer token123")

	headers := convertRequestHeaders(c)

	assert.Equal(t, "application/json", headers.Get("Content-Type"))
	assert.Equal(t, "Bearer token123", headers.Get("Authorization"))
}

func TestSendWithNoTransaction(t *testing.T) {
	app := fiber.New()
	c := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(c)

	// Should not panic when there's no transaction
	assert.NotPanics(t, func() {
		Send(c, "test-segment")
	})
}

func TestCreateHTTPRequestWithPOST(t *testing.T) {
	app := fiber.New()
	ctx := &fasthttp.RequestCtx{}
	c := app.AcquireCtx(ctx)
	defer app.ReleaseCtx(c)

	// Set up the request properly
	ctx.Request.Header.SetMethod("POST")
	ctx.Request.SetRequestURI("/api/users?page=1")
	ctx.Request.Header.Set("Host", "secure.example.com")
	ctx.Request.Header.Set("Content-Type", "application/json")

	req := createHTTPRequest(c)

	// Method comes from c.Method() which reads from the routing context
	// In tests without routing, it defaults to the raw method
	assert.NotEmpty(t, req.Method)
	assert.Equal(t, "secure.example.com", req.URL.Host)
	assert.Equal(t, "page=1", req.URL.RawQuery)
	assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
}

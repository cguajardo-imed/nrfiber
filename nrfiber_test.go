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

func TestCustomTransanctionName(t *testing.T) {
	t.Skip()
	app := fiber.New()
	c := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(c)

	c.Request().Header.SetMethod("GET")
	c.Request().SetRequestURI("/")

	name := defaultTransactionName(c)
	assert.Equal(t, "GET /", name)
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

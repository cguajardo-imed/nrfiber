// Package nrfiber provides New Relic instrumentation for the Fiber web framework.
package nrfiber

import (
	"fmt"
	"net/http"
	"net/url"
	"slices"

	"github.com/gofiber/fiber/v2"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/valyala/fasthttp"
)

// FromContext returns the Transaction from the context if present, and nil
// otherwise.
func FromContext(c *fiber.Ctx) *newrelic.Transaction {
	return newrelic.FromContext(c.UserContext())
}

// Middleware creates Fiber middleware that instrument's requests.
//
// app := fiber.New()
// // Add the nrfiber middleware before other middlewares or routes:
// app.Use(nrfiber.Middleware(app))
func Middleware(app *newrelic.Application, configs ...*config) fiber.Handler {
	if nil == app {
		return func(c *fiber.Ctx) error {
			return c.Next()
		}
	}

	configMap := createConfigMap(configs...)
	createTransactionNameFunc := customTransactionNameFunc(configMap, defaultTransactionName)

	return func(c *fiber.Ctx) error {
		txn := app.StartTransaction(createTransactionNameFunc(c))
		defer txn.End()

		txn.SetWebRequestHTTP(createHTTPRequest(c))

		c.SetUserContext(newrelic.NewContext(c.UserContext(), txn))

		err := c.Next()
		statusCode := c.Context().Response.StatusCode()

		if err != nil {
			statusCode = noticeError(txn, err, configMap)
		}

		txn.SetWebResponse(nil).WriteHeader(statusCode)
		return err
	}
}

func noticeError(txn *newrelic.Transaction, err error, configMap map[string]any) int {
	noticeErrorEnabled := noticeErrorEnabled(configMap)
	statusCodeIgnored := statusCodeIgnored(configMap)
	statusCode := http.StatusInternalServerError

	if fiberErr, ok := err.(*fiber.Error); ok {
		statusCode = fiberErr.Code
	}
	if noticeErrorEnabled {
		found := false

		if slices.Contains(statusCodeIgnored, statusCode) {
			found = true
		}

		if !found {
			txn.NoticeError(err)
		}
	}

	return statusCode
}

func defaultTransactionName(c *fiber.Ctx) string {
	return fmt.Sprintf("%s %s", c.Request().Header.Method(), c.Request().URI().Path())
}

func convertRequestHeaders(fastHTTPHeaders *fasthttp.RequestHeader) http.Header {
	headers := make(http.Header)

	fastHTTPHeaders.VisitAll(func(k, v []byte) {
		headers.Set(string(k), string(v))
	})

	return headers
}

func createHTTPRequest(c *fiber.Ctx) *http.Request {
	reqHeaders := convertRequestHeaders(&c.Request().Header)

	reqHost := reqHeaders.Get("Host")
	if reqHost == "" {
		reqHost = string(c.Request().URI().Host())
	}

	return &http.Request{
		Method: c.Method(),
		URL: &url.URL{
			Scheme:   string(c.Request().URI().Scheme()),
			Host:     reqHost,
			Path:     string(c.Request().URI().Path()),
			RawQuery: string(c.Request().URI().QueryString()),
		},
		Header: reqHeaders,
		Host:   reqHost,
	}
}

func Send(c *fiber.Ctx, segmentName string) {
	txn := FromContext(c)
	segment := txn.StartSegment(segmentName)
	defer segment.End()
}

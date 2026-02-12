// Package nrfiber provides New Relic instrumentation for the Fiber web framework.
package nrfiber

import (
	"fmt"
	"net/http"
	"net/url"
	"slices"

	"github.com/gofiber/fiber/v3"
	"github.com/newrelic/go-agent/v3/newrelic"
)

// FromContext returns the Transaction from the context if present, and nil
// otherwise.
func FromContext(c fiber.Ctx) *newrelic.Transaction {
	return newrelic.FromContext(c.Context())
}

// Middleware creates Fiber middleware that instrument's requests.
//
// app := fiber.New()
// // Add the nrfiber middleware before other middlewares or routes:
// app.Use(nrfiber.Middleware(app))
func Middleware(app *newrelic.Application, configs ...*config) fiber.Handler {
	if nil == app {
		return func(c fiber.Ctx) error {
			return c.Next()
		}
	}

	configMap := createConfigMap(configs...)
	createTransactionNameFunc := customTransactionNameFunc(configMap, defaultTransactionName)

	return func(c fiber.Ctx) error {
		txn := app.StartTransaction(createTransactionNameFunc(c))
		defer txn.End()

		txn.SetWebRequestHTTP(createHTTPRequest(c))

		c.SetContext(newrelic.NewContext(c.Context(), txn))

		err := c.Next()
		statusCode := c.Response().StatusCode()

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

func defaultTransactionName(c fiber.Ctx) string {
	return fmt.Sprintf("%s %s", c.Method(), c.Path())
}

func convertRequestHeaders(c fiber.Ctx) http.Header {
	headers := make(http.Header)

	c.Request().Header.VisitAll(func(k, v []byte) {
		headers.Set(string(k), string(v))
	})

	return headers
}

func createHTTPRequest(c fiber.Ctx) *http.Request {
	reqHeaders := convertRequestHeaders(c)

	reqHost := c.Hostname()
	if reqHost == "" {
		reqHost = reqHeaders.Get("Host")
	}

	scheme := "http"
	if c.Protocol() == "https" {
		scheme = "https"
	}

	return &http.Request{
		Method: c.Method(),
		URL: &url.URL{
			Scheme:   scheme,
			Host:     reqHost,
			Path:     c.Path(),
			RawQuery: string(c.Request().URI().QueryString()),
		},
		Header: reqHeaders,
		Host:   reqHost,
	}
}

func Send(c fiber.Ctx, segmentName string) {
	txn := FromContext(c)
	segment := txn.StartSegment(segmentName)
	defer segment.End()
}

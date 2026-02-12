package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cguajardo-imed/nrfiber"
	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type customErr struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (ce customErr) Error() string {
	return fmt.Sprintf("code: %d, message: %s", ce.Code, ce.Message)
}

func main() {
	godotenv.Load()
	app := fiber.New()
	nr, err := newrelic.NewApplication(
		newrelic.ConfigEnabled(true),
		newrelic.ConfigAppName("demo"),
		newrelic.ConfigLicense(os.Getenv("NEWRELIC_KEY")),
		newrelic.ConfigDistributedTracerEnabled(true),
	)
	if err != nil {
		log.Fatal(err)
	}

	app.Use(nrfiber.Middleware(nr, nrfiber.ConfigNoticeErrorEnabled(true)))

	app.Get("/give-me-error", func(ctx fiber.Ctx) error {
		err := customErr{Message: "wrong request", Code: 4329}
		ctx.Status(http.StatusBadRequest).JSON(err)
		return err
	})
	app.Get("/", func(ctx fiber.Ctx) error {
		return ctx.Status(http.StatusOK).SendString("Hello, World!")
	})
	app.Listen(":8000")
}

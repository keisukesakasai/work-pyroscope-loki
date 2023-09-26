package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	logging "pyroscope-loki-app/internal/log"
	"pyroscope-loki-app/internal/profile"
	"pyroscope-loki-app/internal/trace"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel"
)

const (
	Version = "v1.0.0"
	Service = "pyroscope-loki-app"
)

var logger = logging.NewLogger()
var tracer = otel.Tracer("echo-server")

func main() {
	// Start Profiling
	if serviceAddress := os.Getenv(profile.PyroscopeEndpointURLEnv); serviceAddress != "" {
		profile.Start(serviceAddress)
	}

	// Start Tracing
	tp, err := trace.InitTracer()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	r := echo.New()
	r.Use(otelecho.Middleware("pyroscope-loki-app"))
	r.POST("/", echoHandler)
	r.Start(":8080")
}

func echoHandler(c echo.Context) error {
	body := c.Request().Body
	defer body.Close()

	ctx := c.Request().Context()
	logger := logging.GetLoggerWithTraceID(ctx)
	_, span := tracer.Start(ctx, "Handler")
	defer span.End()

	content, err := io.ReadAll(body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to read the request body")
	}

	logger.Info(string(content))

	return c.String(http.StatusOK, string(content))
}

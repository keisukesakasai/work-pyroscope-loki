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
	"time"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
)

const (
	Version = "v1.0.0"
	Service = "pyroscope-loki-app"
)

var tracer = otel.GetTracerProvider().Tracer("")

func main() {

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

	// Start Profiling
	if serviceAddress := os.Getenv(profile.PyroscopeEndpointURLEnv); serviceAddress != "" {
		profile.Start(serviceAddress)
	}

	r := echo.New()
	// r.Use(otelecho.Middleware("pyroscope-loki-app"))
	r.GET("/", okHandler)
	r.POST("/", echoHandler)
	r.Start(":8080")
}

func okHandler(c echo.Context) error {
	ctx := c.Request().Context()
	ctx, span := otel.GetTracerProvider().Tracer("").Start(ctx, "okHandler")
	defer span.End()
	logger := logging.GetLoggerWithTraceID(ctx)

	logger.Info("ok!!!!!!!!!!!!!!!!!!!!!!!!!")
	time.Sleep(1 * time.Second)

	return c.String(http.StatusOK, "ok")
}

func echoHandler(c echo.Context) error {
	ctx := c.Request().Context()
	ctx, span := otel.GetTracerProvider().Tracer("").Start(ctx, "echoHandler")
	defer span.End()
	logger := logging.GetLoggerWithTraceID(ctx)

	body := c.Request().Body
	defer body.Close()

	content, err := io.ReadAll(body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to read the request body")
	}

	logger.Info(string(content))
	time.Sleep(1 * time.Second)

	return c.String(http.StatusOK, string(content))
}

package main

import (
	"io"
	"net/http"
	"pyroscope-loki-app/internal/log"
	"pyroscope-loki-app/internal/profile"

	"github.com/labstack/echo/v4"
)

const (
	Version = "v1.0.0"
	Service = "pyroscope-loki-app"
)

var logger = log.NewLogger()

func main() {

	profile.Start()

	e := echo.New()

	e.POST("/", echoHandler)

	e.Start(":8080")
}

func echoHandler(c echo.Context) error {
	body := c.Request().Body
	defer body.Close()

	content, err := io.ReadAll(body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to read the request body")
	}

	logger.Debug(string(content))
	logger.Info(string(content))
	logger.Error(string(content))

	return c.String(http.StatusOK, string(content))
}

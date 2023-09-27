package trace

import (
	"context"
	"os"
	logging "pyroscope-loki-app/internal/log"
	"pyroscope-loki-app/internal/profile"
	"pyroscope-loki-app/internal/utils"

	otelpyroscope "github.com/pyroscope-io/otel-profiling-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func InitTracer() (*sdktrace.TracerProvider, error) {
	logger := logging.GetLoggerFromCtx(context.Background())
	serviceName := utils.GetEnv(utils.ServiceNameEnv, "unknown")

	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
	)
	if serviceAddress := os.Getenv(profile.PyroscopeEndpointURLEnv); serviceAddress != "" {
		logger.Info("otelPyroscoep is set...: ", serviceAddress)
		otel.SetTracerProvider(otelpyroscope.NewTracerProvider(tp,
			otelpyroscope.WithAppName(serviceName),
			otelpyroscope.WithRootSpanOnly(false),
			otelpyroscope.WithAddSpanName(true),
			otelpyroscope.WithPyroscopeURL(serviceAddress),
			otelpyroscope.WithProfileBaselineLabels(map[string]string{
				"serviceName": serviceName,
			}),
			otelpyroscope.WithProfileBaselineURL(true),
			otelpyroscope.WithProfileURL(true),
		))
	} else {
		otel.SetTracerProvider(tp)
	}
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
	return tp, nil
}

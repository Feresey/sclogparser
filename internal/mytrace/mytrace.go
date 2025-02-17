package mytrace

import (
	"context"
	"fmt"

	"github.com/Feresey/luxpanel/cmd/luxpanel/config"
	otlgrpc "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.uber.org/fx"

	// otlhttp "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

type FxConfig struct {
	fx.In

	LC     fx.Lifecycle
	Config *config.Config
}

func NewTraceProvider(cfg FxConfig) trace.TracerProvider {
	if !cfg.Config.Trace.Enabled {
		return noop.NewTracerProvider()
	}
	exp := otlgrpc.NewUnstarted(otlgrpc.WithInsecure())
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.Config.Trace.ServiceName),
			semconv.DeploymentEnvironmentKey.String("production"),
			// TODO service version
		)),
	)

	cfg.LC.Append(fx.StartStopHook(
		func(ctx context.Context) error {
			if err := exp.Start(ctx); err != nil {
				return fmt.Errorf("run exporter: %w", err)
			}
			return nil
		},
		func(ctx context.Context) error {
			// if err := exp.Shutdown(ctx); err != nil {
			// 	return fmt.Errorf("shutdown exporter: %w", err)
			// }
			if err := tp.Shutdown(ctx); err != nil {
				return fmt.Errorf("shutdown trace provider: %w", err)
			}
			return nil
		},
	))

	return tp
}

package trace

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	otel_trace "go.opentelemetry.io/otel/trace"
)

type TraceConfig struct {
	OtelEnable      bool   `yaml:"otel_enable"`
	OtelEndpoint    string `yaml:"otel_endpoint"`
	OtelAuth        string `yaml:"otel_auth"`
	OtelServiceName string `yaml:"otel_service_name"`
}

func NewTracerProvider(cfg *TraceConfig) (*trace.TracerProvider, error) {
	ctx := context.Background()
	r, err := resource.New(ctx, resource.WithAttributes(attribute.String("service.name", cfg.OtelServiceName)))
	if err != nil {
		return nil, err
	}
	var opts = []trace.TracerProviderOption{trace.WithResource(r)}

	if cfg.OtelEnable {
		var (
			endpoint = cfg.OtelEndpoint
			auth     = cfg.OtelAuth
		)
		exp, err := otlptracegrpc.New(ctx,
			otlptracegrpc.WithEndpoint(endpoint),
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithHeaders(map[string]string{
				"Authorization": auth,
				"organization":  "default",
				"stream-name":   "default",
			}),
			otlptracegrpc.WithCompressor("gzip"),
			otlptracegrpc.WithTimeout(5*time.Second))
		if err != nil {
			return nil, err
		}
		opts = append(opts, trace.WithBatcher(exp, trace.WithBatchTimeout(5*time.Second)))
	}

	return trace.NewTracerProvider(opts...), nil
}

func NewTracer(tracerProvider *trace.TracerProvider) otel_trace.Tracer {
	return tracerProvider.Tracer("")
}

package telemetry

import (
	"context"
	"fmt"
	"time"

	"OrderService/config"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func New(ctx context.Context, cfg config.Telemetry) (*trace.TracerProvider, error) {
	// метаданные сервиса, какой сервис, его версия, его окружение prod?stage?dev?)
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceVersion(cfg.ServiceVersion),
			semconv.DeploymentEnvironment(cfg.Environment),
		),
	)
	if err != nil {
		return nil, err
	}

	exporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithEndpoint(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)),
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		return nil, err
	}

	// сколько запросов трассировать 1.0 - 100%
	var sampler trace.Sampler
	if cfg.Sampling <= 0 {
		sampler = trace.NeverSample()
	} else if cfg.Sampling >= 1 {
		sampler = trace.AlwaysSample()
	} else {
		sampler = trace.TraceIDRatioBased(cfg.Sampling)
	}

	// инициализация
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter,
			trace.WithMaxExportBatchSize(512),     //размер батча
			trace.WithBatchTimeout(5*time.Second), //отправка батча через 5 сек если не накопился
		),
		trace.WithResource(res),
		trace.WithSampler(sampler),
	)

	otel.SetTracerProvider(tp)

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return tp, nil
}

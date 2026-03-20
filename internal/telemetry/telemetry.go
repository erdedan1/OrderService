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
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func New(ctx context.Context, cfg config.Telemetry) (*sdktrace.TracerProvider, error) {
	if !cfg.Enabled {
		return nil, nil
	}

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
	var sampler sdktrace.Sampler
	if cfg.Sampling <= 0 {
		sampler = sdktrace.NeverSample()
	} else if cfg.Sampling >= 1 {
		sampler = sdktrace.AlwaysSample()
	} else {
		sampler = sdktrace.TraceIDRatioBased(cfg.Sampling)
	}

	// инициализация
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter,
			sdktrace.WithMaxExportBatchSize(512),     //размер батча
			sdktrace.WithBatchTimeout(5*time.Second), //отправка батча через 5 сек если не накопился
		),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sampler),
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

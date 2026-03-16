package main

import (
	"OrderService/config"
	"OrderService/internal/app"
	"context"

	log "github.com/erdedan1/shared/logger"
	tel "github.com/erdedan1/shared/telemetry"
)

func main() {
	ctx := context.Background()

	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	logger, err := log.NewLogger(cfg.Infrastructure.Observability.Loglvl)
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	shutdown, err := tel.New(ctx, config.NewTelemetryConfig(cfg.Infrastructure.Observability.Telemetry))
	if err != nil {
		panic(err)
	}
	defer shutdown(ctx)

	app, err := app.Build(cfg, logger)
	if err != nil {
		panic(err)
	}

	if err := app.Start(ctx); err != nil {
		panic(err)
	}
}

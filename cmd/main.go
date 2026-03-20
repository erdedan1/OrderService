package main

import (
	"OrderService/config"
	"OrderService/internal/app"
	"OrderService/internal/telemetry"
	"context"

	log "github.com/erdedan1/shared/logger"
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

	tp, err := telemetry.New(ctx, cfg.Infrastructure.Observability.Telemetry)
	if err != nil {
		panic(err)
	}
	defer tp.Shutdown(ctx)

	app, err := app.Build(cfg, logger, tp)
	if err != nil {
		panic(err)
	}

	if err := app.Start(ctx); err != nil {
		panic(err)
	}
}

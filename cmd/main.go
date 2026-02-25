package main

import (
	"OrderService/config"
	"OrderService/internal/app"
	"context"

	log "github.com/erdedan1/shared/logger"
	tel "github.com/erdedan1/shared/telemetry"
)

// дико сомневаюсь что логгер и телеметрия должны быть тут, но я сделаль так
// еще сомневаюсь насчет конструктора конфига телеметрии, но показалось лучшим(это я дурачок обосрите меня)
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

	orderService := app.New(cfg, logger)
	if err := orderService.Start(ctx); err != nil {
		panic(err)
	}
}

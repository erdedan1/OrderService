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
		panic(err) //логировать
		//лучше выйти с ошибкой чем с паникой
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

	app, errC := app.Build(cfg, logger, tp)
	if errC != nil {
		panic(errC)
	}

	if err := app.Start(ctx); err != nil {
		panic(err)
	}
}

//на каждого клиента(пользователя) свои лимиты (rate limiter)
//контексты посмотреть
//ДО ЧЕТВЕРГА!

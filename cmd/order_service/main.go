package main

import (
	"OrderService/config"
	"OrderService/internal/app"
	"OrderService/internal/telemetry"
	"context"
	"log"

	logCustom "github.com/erdedan1/shared/logger"
)

func main() {
	ctx := context.Background()

	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
		return
	}

	logger, err := logCustom.NewLogger(cfg.Infrastructure.Observability.Loglvl)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer logger.Sync()

	tp, err := telemetry.New(ctx, cfg.Infrastructure.Observability.Telemetry)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer tp.Shutdown(ctx)

	appInstance, errC := app.Build(cfg, logger, tp)
	if errC != nil {
		log.Fatal(errC)
		return
	}

	if err := appInstance.Start(ctx); err != nil {
		log.Fatal(err)
		return
	}
}

//на каждого клиента(пользователя) свои лимиты (rate limiter)
//контексты посмотреть
//ДО ЧЕТВЕРГА!

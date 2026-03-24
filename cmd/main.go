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

	err := config.New()
	if err != nil {
		panic(err)
	}

	logger, err := log.NewLogger(config.Global.Infrastructure.Observability.Loglvl)
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	tp, err := telemetry.New(ctx, config.Global.Infrastructure.Observability.Telemetry)
	if err != nil {
		panic(err)
	}
	defer tp.Shutdown(ctx)

	app, errC := app.Build(&config.Global, logger, tp)
	if errC != nil {
		panic(errC)
	}

	if err := app.Start(ctx); err != nil {
		panic(err)
	}
}

//вопросы
//насколько плохо иметь глобальный конфиг, насколько хорошо прокидывать в сервисы конфиг даже если он нужен в одной функции

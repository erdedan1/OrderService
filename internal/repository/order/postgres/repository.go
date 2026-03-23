package order

import (
	"context"

	"OrderService/config"
	"OrderService/internal/connection"

	errorz "github.com/erdedan1/shared/errs"
	log "github.com/erdedan1/shared/logger"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"
)

type Repository struct {
	db     *sqlx.DB
	log    log.Logger
	tracer trace.Tracer
}

func New(ctx context.Context, log log.Logger, config config.PostgresDB, tp trace.TracerProvider) (*Repository, *errorz.CustomError) {
	db, err := connection.New(ctx, config)
	if err != nil {
		return nil, err
	}

	return &Repository{
		db:     db,
		log:    log,
		tracer: tp.Tracer("order-service/Repository"),
	}, nil
}

const layerPostgres = "PostgresOrderRepo"

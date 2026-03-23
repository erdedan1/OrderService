package order

import (
	"context"
	"database/sql"
	"errors"

	errs "OrderService/internal/errors"
	"OrderService/internal/model"

	errorz "github.com/erdedan1/shared/errs"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func (r *Repository) GetOrder(ctx context.Context, id uuid.UUID) (*model.Order, *errorz.CustomError) {
	const method = "GetOrder"

	ctx, span := r.tracer.Start(ctx, "OrderRepository.GetOrder")
	defer span.End()

	span.SetAttributes(
		attribute.String("order.id", id.String()),
	)

	query := `SELECT id, user_id, market_id, quantity, type, status, price, created_at, updated_at, deleted_at FROM orders WHERE id = $1`

	var notification model.Order

	err := r.db.GetContext(ctx, &notification, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			span.RecordError(errs.ErrOrderNotFound)
			span.SetStatus(codes.Error, errs.ErrOrderNotFound.Message)

			r.log.Error(
				layerPostgres,
				method,
				err.Error(), err,
				"order_id", id,
			)
			return nil, errs.ErrOrderNotFound
		}

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		r.log.Error(
			layerPostgres,
			method,
			err.Error(), err,
			"order_id", id,
		)
		return nil, errorz.New(errorz.INTERNAL, err.Error(), err)
	}

	span.SetStatus(codes.Ok, "get order success")

	return &notification, nil
}

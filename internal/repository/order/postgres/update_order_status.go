package order

import (
	"context"
	"time"

	errs "OrderService/internal/errors"
	"OrderService/internal/model"

	errorz "github.com/erdedan1/shared/errs"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func (r *Repository) UpdateOrderStatus(ctx context.Context, id uuid.UUID, status model.OrderStatus) *errorz.CustomError {
	const method = "UpdateOrder"

	ctx, span := r.tracer.Start(ctx, "OrderRepository.UpdateOrder")
	defer span.End()

	span.SetAttributes(
		attribute.String("order.id", id.String()),
	)

	query := `
			UPDATE orders
			SET status = $1, updated_at = $2
			WHERE id = $3
		`

	res, err := r.db.ExecContext(ctx, query, status, time.Now(), id)
	if err != nil {
		span.RecordError(errs.ErrOrderNotFound)
		span.SetStatus(codes.Error, errs.ErrOrderNotFound.Message)

		r.log.Error(
			layerPostgres,
			method,
			err.Error(), err,
			"order_id", id,
		)
		return errorz.New(errorz.INTERNAL, err.Error(), err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		r.log.Error(
			layerPostgres,
			method,
			err.Error(), err,
			"order_id", id,
		)
		return errorz.New(errorz.INTERNAL, err.Error(), err)
	}
	if rowsAffected == 0 {
		span.RecordError(errs.ErrOrderNotFound)
		span.SetStatus(codes.Error, errs.ErrOrderNotFound.Message)

		r.log.Error(
			layerPostgres,
			method,
			errs.ErrOrderNotFound.Message,
			errs.ErrOrderNotFound,
		)
		return errs.ErrOrderNotFound
	}

	span.SetStatus(codes.Ok, "order success updated")

	return nil
}

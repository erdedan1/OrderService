package order

import (
	"context"

	"OrderService/internal/model"

	errorz "github.com/erdedan1/shared/errs"
	"go.opentelemetry.io/otel/codes"
)

func (r *Repository) CreateOrder(ctx context.Context, order *model.Order) (*model.Order, *errorz.CustomError) {
	const method = "CreateOrder"

	ctx, span := r.tracer.Start(ctx, "OrderRepository.CreateOrder")
	defer span.End()

	query := `
		INSERT INTO orders (user_id, market_id, quantity, order_type, order_status, price)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`
	row := r.db.QueryRowxContext(ctx, query, order.UserUUID, order.MarketUUID, order.Quantity, order.Type, order.Status, order.Price)

	err := row.Scan(&order.ID, &order.CreatedAt)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		r.log.Error(
			layerPostgres,
			method,
			err.Error(), err,
			"order_id", order.ID,
		)
		return nil, errorz.New(errorz.INTERNAL, "failed to create order")
	}

	span.SetStatus(codes.Ok, "order success created")

	return order, nil
}

package order

import (
	"context"
	"database/sql"
	"errors"

	"OrderService/config"
	"OrderService/internal/connection"
	errs "OrderService/internal/errors"
	"OrderService/internal/model"

	errorz "github.com/erdedan1/shared/errs"
	log "github.com/erdedan1/shared/logger"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PostgresRepo struct {
	db *sqlx.DB
	l  log.Logger
}

func New(ctx context.Context, l log.Logger, config config.PostgresDB) (*PostgresRepo, *errorz.CustomError) {
	db, err := connection.New(ctx, config)
	if err != nil {
		return nil, err
	}
	return &PostgresRepo{
		db: db,
		l:  l,
	}, nil
}

const layerPostgres = "PostgresOrderDB"

func (r *PostgresRepo) CreateOrder(ctx context.Context, order model.Order) (*model.Order, *errorz.CustomError) {
	const method = "CreateOrder"
	query := `
		INSERT INTO orders (user_id, market_id, quantity, type, status, price)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`
	row := r.db.QueryRowxContext(ctx, query, order.UserId, order.MarketId, order.Quantity, order.Type, order.Status, order.Price)
	err := row.Scan(&order.ID, &order.CreatedAt)
	if err != nil {
		r.l.Error(
			layerInMemory,
			method,
			err.Error(), err,
			"order_id", order.ID,
		)
		return nil, errorz.New(errorz.INTERNAL, err.Error())
	}
	return &order, nil
}

func (r *PostgresRepo) GetOrder(ctx context.Context, id uuid.UUID) (*model.Order, *errorz.CustomError) {
	const method = "GetOrder"

	query := `SELECT id, user_id, market_id, quantity, type, status, price, created_at, updated_at, deleted_at FROM orders WHERE id = $1`

	var notification model.Order

	err := r.db.GetContext(ctx, &notification, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.l.Error(
				layerInMemory,
				method,
				err.Error(), err,
				"order_id", id,
			)
			return nil, errs.ErrOrderNotFound
		}
		r.l.Error(
			layerInMemory,
			method,
			err.Error(), err,
			"order_id", id,
		)
		return nil, errorz.New(errorz.INTERNAL, err.Error())
	}
	return &notification, nil
}

func (r *PostgresRepo) UpdateOrder(ctx context.Context, id uuid.UUID, order model.Order) (*model.Order, *errorz.CustomError) {
	const method = "UpdateOrder"

	query := `
			UPDATE orders
			SET user_id = $1, market_id = $2, quantity = $3, type = $4, status = $5, price = $6, updated_at = $7, deleted_at = $8
			WHERE id = $9
			`
	res, err := r.db.ExecContext(ctx, query)
	if err != nil {
		r.l.Error(
			layerInMemory,
			method,
			err.Error(), err,
			"order_id", id,
		)
		return nil, errorz.New(errorz.INTERNAL, err.Error())
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		r.l.Error(
			layerInMemory,
			method,
			err.Error(), err,
			"order_id", id,
		)
		return nil, errorz.New(errorz.INTERNAL, err.Error())
	}
	if rowsAffected == 0 {
		r.l.Error(
			layerInMemory,
			method,
			errs.ErrOrderNotFound.Message,
			errs.ErrOrderNotFound,
		)
		return nil, errs.ErrOrderNotFound
	}

	return &order, nil
}

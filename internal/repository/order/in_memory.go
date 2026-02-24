package order

import (
	"context"
	"sync"

	errs "OrderService/internal/errors"
	"OrderService/internal/model"

	errors "github.com/erdedan1/shared/errs"
	log "github.com/erdedan1/shared/logger"
	"github.com/google/uuid"
)

type InMemoryRepo struct {
	Orders map[uuid.UUID]*model.Order
	mu     *sync.RWMutex
	l      log.Logger
}

func NewInMemoryRepo(logger log.Logger) *InMemoryRepo {
	return &InMemoryRepo{
		Orders: make(map[uuid.UUID]*model.Order),
		mu:     &sync.RWMutex{},
		l:      logger,
	}
}

const layerInMemory = "OrderInMemoryRepo"

func (r *InMemoryRepo) CreateOrder(ctx context.Context, order *model.Order) (*model.Order, *errors.CustomError) {
	r.mu.Lock()
	defer r.mu.Unlock()
	order.ID = uuid.New()
	r.Orders[order.ID] = order

	r.l.Debug(
		layerInMemory,
		"CreateOrder",
		"order success created",
		"order_id", order.ID,
	)

	return order, nil
}

func (r *InMemoryRepo) GetOrder(ctx context.Context, id uuid.UUID) (*model.Order, *errors.CustomError) {
	const method = "GetOrder"
	r.mu.RLock()
	defer r.mu.RUnlock()

	if o, found := r.Orders[id]; found {
		r.l.Debug(
			layerInMemory,
			method,
			"get order info",
			"uorder_id", id,
		)
		return o, nil
	}

	r.l.Error(
		layerInMemory, method,
		"order not found",
		errs.ErrOrderNotFound,
		"order_id", id,
	)

	return nil, errs.ErrOrderNotFound
}

func (r *InMemoryRepo) UpdateOrder(ctx context.Context, id uuid.UUID, order *model.Order) (*model.Order, *errors.CustomError) {
	const method = "UpdateOrder"
	r.mu.Lock()
	defer r.mu.Unlock()
	if o, found := r.Orders[id]; found {
		r.l.Debug(
			layerInMemory,
			method,
			"order success updated",
			"order_id", id,
			"user_id", order.UserId,
		)
		return o.Update(order), nil
	}

	r.l.Error(
		layerInMemory,
		method,
		"not found order",
		errs.ErrOrderNotFound,
		"order_id", id,
	)

	return nil, errs.ErrOrderNotFound
}

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

type Repo struct {
	Orders map[uuid.UUID]*model.Order
	mu     *sync.RWMutex
	l      log.Logger
}

func NewRepo(logger log.Logger) *Repo {
	return &Repo{
		Orders: make(map[uuid.UUID]*model.Order),
		mu:     &sync.RWMutex{},
		l:      logger.Layer("Order.Repository"),
	}
}

func (r *Repo) CreateOrder(ctx context.Context, order model.Order) (*model.Order, *errors.CustomError) {
	r.mu.Lock()
	defer r.mu.Unlock()
	order.ID = uuid.New()
	r.Orders[order.ID] = &order

	r.l.Debug("CreateOrder", "order success created", order.ID)

	return &order, nil
}

func (r *Repo) GetOrder(ctx context.Context, id uuid.UUID) (*model.Order, *errors.CustomError) {
	const method = "GetOrder"
	r.mu.RLock()
	defer r.mu.RUnlock()

	if o, found := r.Orders[id]; found {
		r.l.Debug(method, "get order info", id)
		return o, nil
	}

	r.l.Error(method, "order not found", errs.ErrOrderNotFound, id)

	return nil, errs.ErrOrderNotFound
}

func (r *Repo) UpdateOrder(ctx context.Context, id uuid.UUID, order model.Order) (*model.Order, *errors.CustomError) {
	const method = "UpdateOrder"
	r.mu.Lock()
	defer r.mu.Unlock()
	if o, found := r.Orders[id]; found {
		r.l.Debug(method, "order success updated", id)
		return o.Update(order), nil
	}

	r.l.Error(method, "not found order", errs.ErrOrderNotFound, id)

	return nil, errs.ErrOrderNotFound
}

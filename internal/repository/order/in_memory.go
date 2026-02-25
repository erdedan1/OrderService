package order

import (
	"context"
	"sync"

	errs "OrderService/internal/errors"
	"OrderService/internal/model"

	errors "github.com/erdedan1/shared/errs"
	log "github.com/erdedan1/shared/logger"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type InMemoryRepo struct {
	Orders map[uuid.UUID]*model.Order
	mu     *sync.RWMutex
	l      log.Logger
	tracer trace.Tracer
}

func NewInMemoryRepo(logger log.Logger) *InMemoryRepo {
	return &InMemoryRepo{
		Orders: make(map[uuid.UUID]*model.Order),
		mu:     &sync.RWMutex{},
		l:      logger,
		tracer: otel.Tracer("order-service/InMemoryRepo"),
	}
}

const layerInMemory = "OrderInMemoryRepo"

func (r *InMemoryRepo) CreateOrder(ctx context.Context, order *model.Order) (*model.Order, *errors.CustomError) {
	ctx, span := r.tracer.Start(ctx, "OrderInMemoryRepo.CreateOrder")
	defer span.End()

	r.mu.Lock()
	defer r.mu.Unlock()
	order.ID = uuid.New()

	span.SetAttributes(
		attribute.String("order.id", order.ID.String()),
	)

	r.Orders[order.ID] = order

	span.SetStatus(codes.Ok, "order success created")

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

	ctx, span := r.tracer.Start(ctx, "OrderInMemoryRepo.GetOrder")
	defer span.End()

	r.mu.RLock()
	defer r.mu.RUnlock()

	if o, found := r.Orders[id]; found {
		span.SetStatus(codes.Ok, "get order success")

		r.l.Debug(
			layerInMemory,
			method,
			"get order info",
			"uorder_id", id,
		)
		return o, nil
	}

	span.RecordError(errs.ErrOrderNotFound)
	span.SetStatus(codes.Error, errs.ErrOrderNotFound.Message)

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

	ctx, span := r.tracer.Start(ctx, "OrderInMemoryRepo.UpdateOrder")
	defer span.End()

	r.mu.Lock()
	defer r.mu.Unlock()
	if o, found := r.Orders[id]; found {

		span.SetStatus(codes.Ok, "order success updated")

		r.l.Debug(
			layerInMemory,
			method,
			"order success updated",
			"order_id", id,
			"user_id", order.UserId,
		)
		return o.Update(order), nil
	}

	span.RecordError(errs.ErrOrderNotFound)
	span.SetStatus(codes.Error, errs.ErrOrderNotFound.Message)

	r.l.Error(
		layerInMemory,
		method,
		"not found order",
		errs.ErrOrderNotFound,
		"order_id", id,
	)

	return nil, errs.ErrOrderNotFound
}

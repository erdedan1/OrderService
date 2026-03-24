package order

import (
	"context"
	"sync"
	"time"

	errs "OrderService/internal/errors"
	"OrderService/internal/model"

	errors "github.com/erdedan1/shared/errs"
	log "github.com/erdedan1/shared/logger"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Repo struct {
	Orders map[uuid.UUID]*model.Order
	mu     *sync.RWMutex
	log    log.Logger
	tracer trace.Tracer
}

func NewRepo(logger log.Logger, tp trace.TracerProvider) *Repo {
	return &Repo{
		Orders: make(map[uuid.UUID]*model.Order),
		mu:     &sync.RWMutex{},
		log:    logger,
		tracer: tp.Tracer("order-service/Repo"),
	}
}

const layerInMemory = "OrderRepo"

func (r *Repo) CreateOrder(ctx context.Context, order *model.Order) (*model.Order, *errors.CustomError) {
	ctx, span := r.tracer.Start(ctx, "OrderRepo.CreateOrder")
	defer span.End()

	r.mu.Lock()
	defer r.mu.Unlock()
	order.ID = uuid.New()

	span.SetAttributes(
		attribute.String("order.id", order.ID.String()),
	)

	r.Orders[order.ID] = order

	span.SetStatus(codes.Ok, "order success created")

	r.log.Debug(
		layerInMemory,
		"CreateOrder",
		"order success created",
		"order_id", order.ID,
		"order_status", order.Status,
	)

	return order, nil
}

func (r *Repo) GetOrder(ctx context.Context, id uuid.UUID) (*model.Order, *errors.CustomError) {
	const method = "GetOrder"

	ctx, span := r.tracer.Start(ctx, "OrderRepo.GetOrder")
	defer span.End()

	r.mu.RLock()
	defer r.mu.RUnlock()

	if o, found := r.Orders[id]; found {
		span.SetStatus(codes.Ok, "get order success")

		r.log.Debug(
			layerInMemory,
			method,
			"get order info",
			"uorder_id", id,
		)
		return o, nil
	}

	span.RecordError(errs.ErrOrderNotFound)
	span.SetStatus(codes.Error, errs.ErrOrderNotFound.Message)

	r.log.Error(
		layerInMemory, method,
		"order not found",
		errs.ErrOrderNotFound,
		"order_id", id,
	)

	return nil, errs.ErrOrderNotFound
}

func (r *Repo) UpdateOrderStatus(ctx context.Context, id uuid.UUID, status model.OrderStatus) *errors.CustomError {
	const method = "UpdateOrder"

	ctx, span := r.tracer.Start(ctx, "OrderRepo.UpdateOrder")
	defer span.End()

	r.mu.Lock()
	defer r.mu.Unlock()
	if o, found := r.Orders[id]; found {

		span.SetStatus(codes.Ok, "order success updated")

		r.log.Debug(
			layerInMemory,
			method,
			"order success updated",
			"order_id", id,
			"user_id", o.UserID,
			"order_status", o.Status,
		)
		o.Status = status
		o.UpdatedAt = new(time.Now())
		return nil
	}

	span.RecordError(errs.ErrOrderNotFound)
	span.SetStatus(codes.Error, errs.ErrOrderNotFound.Message)

	r.log.Error(
		layerInMemory,
		method,
		"not found order",
		errs.ErrOrderNotFound,
		"order_id", id,
	)

	return errs.ErrOrderNotFound
}

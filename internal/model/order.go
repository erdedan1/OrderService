package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Order struct {
	ID        uuid.UUID
	UserId    uuid.UUID
	MarketId  uuid.UUID
	Quantity  int64
	Type      string
	Status    string
	Price     decimal.Decimal
	CreatedAt time.Time
	UpdateAt  time.Time
	DeletedAt time.Time
}

func (o *Order) Update(order *Order) *Order {
	o.UserId = order.UserId
	o.MarketId = order.MarketId
	o.Quantity = order.Quantity
	o.Type = order.Type
	o.Status = order.Status
	o.Price = order.Price
	o.CreatedAt = order.CreatedAt
	o.UpdateAt = order.UpdateAt
	o.DeletedAt = order.DeletedAt
	return o
}

const (
	StatusCreated       = "ORDER_STATUS_CREATED"
	StatusPending       = "ORDER_STATUS_PENDING"
	StatusWaitSeller    = "ORDER_STATUS_WAIT_SELLER"
	StatusPaid          = "ORDER_STATUS_PAID"
	StatusOnHold        = "ORDER_STATUS_ON_HOLD"
	StatusProcessing    = "ORDER_STATUS_PROCESSING"
	StatusPacked        = "ORDER_STATUS_PACKED"
	StatusOutOfDelivery = "ORDER_STATUS_OUT_OF_DELIVERY"
	StatusOnTheWay      = "ORDER_STATUSON_THE_WAY"
	StatusDelivered     = "ORDER_STATUS_DELIVERED"
	StatusClosed        = "ORDER_STATUS_CLOSED"
)

var OrderStatusProcessing []string = []string{
	StatusCreated,
	StatusPending,
	StatusWaitSeller,
	StatusPaid,
	StatusOnHold,
	StatusProcessing,
	StatusPacked,
	StatusOutOfDelivery,
	StatusOnTheWay,
	StatusDelivered,
	StatusClosed,
}

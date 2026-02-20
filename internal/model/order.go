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

func (o *Order) Update(order Order) *Order {
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
	StatusCreated       = "created"
	StatusPending       = "pending"
	StatusWaitSeller    = "wait_seller"
	StatusPaid          = "paid"
	StatusOnHold        = "on_hold"
	StatusProcessing    = "processing"
	StatusPacked        = "packed"
	StatusOutOfDelivery = "out_of_delivery"
	StatusOnTheWay      = "on_the_way"
	StatusDelivered     = "delivered"
	StatusClosed        = "closed"
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

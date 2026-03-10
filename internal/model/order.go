package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Order struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	MarketID  uuid.UUID
	Quantity  int64
	Type      string
	Status    OrderStatus
	Price     decimal.Decimal
	CreatedAt time.Time
	UpdateAt  time.Time
	DeletedAt time.Time
}

func (o OrderStatus) ToString() string {
	switch {
	case o == StatusCreated:
		return "CREATED"
	default:
		return ""
	}
}

type OrderStatus string

const (
	OrderStatusCreated OrderStatus = "CREATED"
)

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

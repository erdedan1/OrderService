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
	UpdatedAt time.Time
	DeletedAt time.Time
}

type OrderStatus string

const (
	StatusCreated       OrderStatus = "ORDER_STATUS_CREATED"
	StatusPending       OrderStatus = "ORDER_STATUS_PENDING"
	StatusWaitSeller    OrderStatus = "ORDER_STATUS_WAIT_SELLER"
	StatusPaid          OrderStatus = "ORDER_STATUS_PAID"
	StatusOnHold        OrderStatus = "ORDER_STATUS_ON_HOLD"
	StatusProcessing    OrderStatus = "ORDER_STATUS_PROCESSING"
	StatusPacked        OrderStatus = "ORDER_STATUS_PACKED"
	StatusOutOfDelivery OrderStatus = "ORDER_STATUS_OUT_OF_DELIVERY"
	StatusOnTheWay      OrderStatus = "ORDER_STATUS_ON_THE_WAY"
	StatusDelivered     OrderStatus = "ORDER_STATUS_DELIVERED"
	StatusClosed        OrderStatus = "ORDER_STATUS_CLOSED"
)

var OrderStatusProcessing = []OrderStatus{
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

func (o OrderStatus) ToString() string {
	switch o {
	case StatusCreated:
		return "CREATED"
	case StatusPending:
		return "PENDING"
	case StatusWaitSeller:
		return "WAIT_SELLER"
	case StatusPaid:
		return "PAID"
	case StatusOnHold:
		return "ON_HOLD"
	case StatusProcessing:
		return "PROCESSING"
	case StatusPacked:
		return "PACKED"
	case StatusOutOfDelivery:
		return "OUT_OF_DELIVERY"
	case StatusOnTheWay:
		return "ON_THE_WAY"
	case StatusDelivered:
		return "DELIVERED"
	case StatusClosed:
		return "CLOSED"
	default:
		return "UNKNOWN"
	}
}

func NextOrderStatus(current OrderStatus) (OrderStatus, bool) {
	for idx, status := range OrderStatusProcessing {
		if status != current {
			continue
		}
		if idx+1 >= len(OrderStatusProcessing) {
			return current, false
		}
		return OrderStatusProcessing[idx+1], true
	}
	return StatusCreated, true
}

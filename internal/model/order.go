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
	StatusUnspecified   OrderStatus = "ORDER_STATUS_UNSPECIFIED"
)

func (o OrderStatus) ToString() string {
	switch o {
	case StatusCreated:
		return "ORDER_STATUS_CREATED"
	case StatusPending:
		return "ORDER_STATUS_PENDING"
	case StatusWaitSeller:
		return "ORDER_STATUS_WAIT_SELLER"
	case StatusPaid:
		return "ORDER_STATUS_PAID"
	case StatusOnHold:
		return "ORDER_STATUS_ON_HOLD"
	case StatusProcessing:
		return "ORDER_STATUS_PROCESSING"
	case StatusPacked:
		return "ORDER_STATUS_PACKED"
	case StatusOutOfDelivery:
		return "ORDER_STATUS_OUT_OF_DELIVERY"
	case StatusOnTheWay:
		return "ORDER_STATUS_ON_THE_WAY"
	case StatusDelivered:
		return "ORDER_STATUS_DELIVERED"
	case StatusClosed:
		return "ORDER_STATUS_CLOSED"
	default:
		return "ORDER_STATUS_UNSPECIFIED"
	}
}

func NextOrderStatus(current OrderStatus) (OrderStatus, bool) {
	switch current {
	case StatusCreated:
		return StatusPending, true
	case StatusPending:
		return StatusWaitSeller, true
	case StatusWaitSeller:
		return StatusPaid, true
	case StatusPaid:
		return StatusOnHold, true
	case StatusOnHold:
		return StatusProcessing, true
	case StatusProcessing:
		return StatusPacked, true
	case StatusPacked:
		return StatusOutOfDelivery, true
	case StatusOutOfDelivery:
		return StatusOnTheWay, true
	case StatusOnTheWay:
		return StatusDelivered, true
	case StatusDelivered:
		return StatusClosed, true
	default:
		return current, false
	}
}

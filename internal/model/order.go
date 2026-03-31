package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Order struct {
	ID         uuid.UUID       `db:"id"`
	UserUUID   uuid.UUID       `db:"user_id"`
	MarketUUID uuid.UUID       `db:"market_id"`
	Quantity   int64           `db:"quantity"`
	Type       string          `db:"order_type"`
	Status     OrderStatus     `db:"order_status"`
	Price      decimal.Decimal `db:"price"`
	CreatedAt  *time.Time      `db:"created_at"`
	UpdatedAt  *time.Time      `db:"updated_at"`
	DeletedAt  *time.Time      `db:"deleted_at"`
}

type OrderStatus string

const (
	StatusCreated       OrderStatus = "CREATED"
	StatusPending       OrderStatus = "PENDING"
	StatusWaitSeller    OrderStatus = "WAIT_SELLER"
	StatusPaid          OrderStatus = "PAID"
	StatusOnHold        OrderStatus = "ON_HOLD"
	StatusProcessing    OrderStatus = "PROCESSING"
	StatusPacked        OrderStatus = "PACKED"
	StatusOutOfDelivery OrderStatus = "OUT_OF_DELIVERY"
	StatusOnTheWay      OrderStatus = "ON_THE_WAY"
	StatusDelivered     OrderStatus = "DELIVERED"
	StatusClosed        OrderStatus = "CLOSED"
	StatusUnspecified   OrderStatus = "UNSPECIFIED"
)

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
		return "UNSPECIFIED"
	}
}

func (o OrderStatus) IsValid() bool {
	switch o {
	case StatusCreated,
		StatusPending,
		StatusWaitSeller,
		StatusPaid,
		StatusOnHold,
		StatusProcessing,
		StatusPacked,
		StatusOutOfDelivery,
		StatusOnTheWay,
		StatusDelivered,
		StatusClosed:
		return true
	default:
		return false
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

func NewOrder(
	userID, marketID uuid.UUID,
	quantity int64,
	price decimal.Decimal,
	orderType string,
) *Order {

	return &Order{
		UserUUID:   userID,
		MarketUUID: marketID,
		Quantity:   quantity,
		Type:       orderType,
		Status:     StatusCreated,
		Price:      price,
		CreatedAt:  new(time.Now()),
	}
}

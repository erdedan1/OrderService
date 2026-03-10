package dto

import (
	"github.com/google/uuid"
)

type CreateOrderRequest struct {
	UserID    uuid.UUID
	MarketID  uuid.UUID
	OrderType string
	Price     string
	UserRoles []string
	Quantity  int64
}

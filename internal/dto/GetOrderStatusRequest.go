package dto

import (
	"github.com/google/uuid"
)

type GetOrderStatusRequest struct {
	UserID  uuid.UUID
	OrderID uuid.UUID
}

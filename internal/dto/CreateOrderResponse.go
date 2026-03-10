package dto

import (
	"github.com/google/uuid"
)

type CreateOrderResponse struct {
	ID     uuid.UUID
	Status string
}

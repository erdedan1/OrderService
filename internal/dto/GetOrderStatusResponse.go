package dto

import (
	"time"
)

type GetOrderStatusResponse struct {
	Status    string
	UpdatedAt *time.Time
}

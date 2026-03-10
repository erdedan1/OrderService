package model

import (
	"time"

	"github.com/google/uuid"
)

type Market struct {
	ID        uuid.UUID
	Name      string
	Enabled   bool
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

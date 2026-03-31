package dto

import (
	"time"

	pb "github.com/erdedan1/protocol/proto/order_service/gen/v1"
	m "github.com/erdedan1/shared/mapper"
	"github.com/google/uuid"
)

type CreateOrderResponse struct {
	OrderUUID uuid.UUID
	Status    string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

func (c *CreateOrderResponse) ToProto() *pb.CreateOrderResponse {
	return &pb.CreateOrderResponse{
		OrderUuid: c.OrderUUID.String(),
		Status:    m.OrderStatusToProto(c.Status),
	}
}

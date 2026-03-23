package dto

import (
	pb "github.com/erdedan1/protocol/proto/order_service/gen/v2"
	m "github.com/erdedan1/shared/mapper"
	"github.com/google/uuid"
)

type CreateOrderResponse struct {
	ID     uuid.UUID
	Status string
}

func (c *CreateOrderResponse) ToProto() *pb.CreateOrderResponse {
	return &pb.CreateOrderResponse{
		Id:     c.ID.String(),
		Status: m.OrderStatusToProto(c.Status),
	}
}

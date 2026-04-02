package dto

import (
	"time"

	pb "github.com/erdedan1/protocol/proto/order_service/gen/v1"
	m "github.com/erdedan1/shared/mapper"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CreateOrderResponse struct {
	OrderUUID uuid.UUID
	Status    string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

func (c *CreateOrderResponse) ToProto() *pb.CreateOrderResponse {
	response := &pb.CreateOrderResponse{
		OrderUuid: c.OrderUUID.String(),
		Status:    m.OrderStatusToProto(c.Status),
	}
	if c.CreatedAt != nil {
		response.CreatedAt = timestamppb.New(*c.CreatedAt)
	}
	if c.UpdatedAt != nil {
		response.UpdatedAt = timestamppb.New(*c.UpdatedAt)
	}

	return response
}

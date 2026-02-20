package dto

import (
	pb "github.com/erdedan1/protocol/proto/order_service/gen"
	m "github.com/erdedan1/shared/mapper"
	"github.com/google/uuid"
)

type CreateOrderResponse struct {
	ID     uuid.UUID
	Status string
}

func (gor *CreateOrderResponse) DtoToProto() *pb.CreateOrderResponse {
	return &pb.CreateOrderResponse{
		Id:     gor.ID.String(),
		Status: m.StringOrderStatusToProto(gor.Status),
	}
}

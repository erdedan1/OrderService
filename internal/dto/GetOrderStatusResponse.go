package dto

import (
	"time"

	pb "github.com/erdedan1/protocol/proto/order_service/gen/v2"
	m "github.com/erdedan1/shared/mapper"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GetOrderStatusResponse struct {
	Status    string
	UpdatedAt *time.Time
}

func (g *GetOrderStatusResponse) ToProto() *pb.GetOrderStatusResponse {
	return &pb.GetOrderStatusResponse{
		Status:    m.OrderStatusToProto(g.Status),
		UpdatedAt: timestamppb.New(*g.UpdatedAt),
	}
}

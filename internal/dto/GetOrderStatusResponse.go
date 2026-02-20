package dto

import (
	"time"

	pb "github.com/erdedan1/protocol/proto/order_service/gen"
	m "github.com/erdedan1/shared/mapper"
)

type GetOrderStatusResponse struct {
	Status   string
	UpdateAt *time.Time
}

func (gosr *GetOrderStatusResponse) DtoToProto() *pb.GetOrderStatusResponse {
	return &pb.GetOrderStatusResponse{
		Status:    m.StringOrderStatusToProto(gosr.Status),
		UpdatedAt: m.ToTimestampProto(gosr.UpdateAt),
	}
}

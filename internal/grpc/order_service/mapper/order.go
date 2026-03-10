package mapper

import (
	"OrderService/internal/dto"
	errs "OrderService/internal/errors"

	pb "github.com/erdedan1/protocol/proto/order_service/gen"
	errors "github.com/erdedan1/shared/errs"
	m "github.com/erdedan1/shared/mapper"
	"github.com/google/uuid"
)

func CreateOrderRequestFromProto(request *pb.CreateOrderRequest) (*dto.CreateOrderRequest, *errors.CustomError) {
	if len(request.UserRoles) == 0 {
		return nil, errs.ErrInvalidArgument
	}

	userId, err := uuid.Parse(request.UserId)
	if err != nil {
		return nil, errs.ErrInvalidArgument
	}

	marketId, err := uuid.Parse(request.MarketId)
	if err != nil {
		return nil, errs.ErrInvalidArgument
	}

	return &dto.CreateOrderRequest{
		MarketID:  marketId,
		UserID:    userId,
		OrderType: request.OrderType,
		Price:     request.Price,
		UserRoles: m.ProtoUserRolesToString(request.UserRoles),
		Quantity:  request.Quantity,
	}, nil
}

func CreateOrderResponseToProto(resp *dto.CreateOrderResponse) *pb.CreateOrderResponse {
	return &pb.CreateOrderResponse{
		Id:     resp.ID.String(),
		Status: m.StringOrderStatusToProto(resp.Status),
	}
}

func GetOrderStatusRequestFromProto(request *pb.GetOrderStatusRequest) (*dto.GetOrderStatusRequest, *errors.CustomError) {
	userId, err := uuid.Parse(request.UserId)
	if err != nil {
		return nil, errs.ErrInvalidArgument
	}

	orderId, err := uuid.Parse(request.OrderId)
	if err != nil {
		return nil, errs.ErrInvalidArgument
	}

	return &dto.GetOrderStatusRequest{
		UserID:  userId,
		OrderID: orderId,
	}, nil
}

func GetOrderStatusResponseToProto(resp *dto.GetOrderStatusResponse) *pb.GetOrderStatusResponse {
	return &pb.GetOrderStatusResponse{
		Status:    m.StringOrderStatusToProto(resp.Status),
		UpdatedAt: m.ToTimestampProto(resp.UpdatedAt),
	}
}

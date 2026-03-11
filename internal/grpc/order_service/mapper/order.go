package mapper

import (
	errs "OrderService/internal/errors"
	"OrderService/internal/usecase"

	pb "github.com/erdedan1/protocol/proto/order_service/gen"
	errors "github.com/erdedan1/shared/errs"
	m "github.com/erdedan1/shared/mapper"
	"github.com/google/uuid"
)

func CreateOrderRequestFromProto(request *pb.CreateOrderRequest) (*usecase.CreateOrderInput, *errors.CustomError) {
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

	return &usecase.CreateOrderInput{
		MarketID:  marketId,
		UserID:    userId,
		OrderType: request.OrderType,
		Price:     request.Price,
		UserRoles: m.ProtoUserRolesToString(request.UserRoles),
		Quantity:  request.Quantity,
	}, nil
}

func CreateOrderResponseToProto(resp *usecase.CreateOrderOutput) *pb.CreateOrderResponse {
	return &pb.CreateOrderResponse{
		Id:     resp.ID.String(),
		Status: m.StringOrderStatusToProto(resp.Status),
	}
}

func GetOrderStatusRequestFromProto(request *pb.GetOrderStatusRequest) (*usecase.GetOrderStatusInput, *errors.CustomError) {
	userId, err := uuid.Parse(request.UserId)
	if err != nil {
		return nil, errs.ErrInvalidArgument
	}

	orderId, err := uuid.Parse(request.OrderId)
	if err != nil {
		return nil, errs.ErrInvalidArgument
	}

	return &usecase.GetOrderStatusInput{
		UserID:  userId,
		OrderID: orderId,
	}, nil
}

func GetOrderStatusResponseToProto(resp *usecase.GetOrderStatusOutput) *pb.GetOrderStatusResponse {
	return &pb.GetOrderStatusResponse{
		Status:    m.StringOrderStatusToProto(resp.Status),
		UpdatedAt: m.ToTimestampProto(resp.UpdatedAt),
	}
}

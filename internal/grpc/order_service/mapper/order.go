package mapper

import (
	"time"

	errs "OrderService/internal/errors"
	"OrderService/internal/usecase"

	pb "github.com/erdedan1/protocol/proto/order_service/gen"
	errors "github.com/erdedan1/shared/errs"
	m "github.com/erdedan1/shared/mapper"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
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
		UserRoles: m.UserRolesFromProto(request.UserRoles),
		Quantity:  request.Quantity,
	}, nil
}

func CreateOrderResponseToProto(resp *usecase.CreateOrderOutput) *pb.CreateOrderResponse {
	return &pb.CreateOrderResponse{
		Id:     resp.ID.String(),
		Status: m.OrderStatusToProto(resp.Status),
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
		Status:    m.OrderStatusToProto(resp.Status),
		UpdatedAt: timestamppb.New(*resp.UpdatedAt),
	}
}

func ToDtoTime(value *timestamppb.Timestamp) *time.Time {
	if value == nil || !value.IsValid() {
		return nil
	}
	t := value.AsTime()
	return &t
}

func ToTimestampProto(value *time.Time) *timestamppb.Timestamp {
	if value == nil || value.IsZero() {
		return nil
	}

	return timestamppb.New(*value)
}

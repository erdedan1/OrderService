package order_service

import (
	"context"

	"OrderService/internal/dto"
	"OrderService/internal/usecase"

	pb "github.com/erdedan1/protocol/proto/order_service/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	srvs usecase.Services
	pb.UnimplementedOrderServiceServer
}

func New(srvs usecase.Services) *Service {
	return &Service{
		srvs: srvs,
	}
}

func (s *Service) CreateOrder(ctx context.Context, request *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	requestDto, err := dto.NewCreateOrderRequest(request)
	if err != nil {
		return nil, status.Error(codes.Code(err.Code), err.Message)
	}

	order, err := s.srvs.OrderService.CreateOrder(ctx, requestDto)
	if err != nil {
		return nil, err
	}

	return order.DtoToProto(), nil
}

func (s *Service) GetOrderStatus(ctx context.Context, request *pb.GetOrderStatusRequest) (*pb.GetOrderStatusResponse, error) {
	requestDto, err := dto.NewGetOrderStatusRequest(request)
	if err != nil {
		return nil, status.Error(codes.Code(err.Code), err.Message)
	}

	order, err := s.srvs.OrderService.GetOrderStatus(ctx, requestDto)
	if err != nil {
		return nil, status.Error(codes.Code(err.Code), err.Message)
	}

	return order.DtoToProto(), nil
}

func (s *Service) SubscribeOrderStatus(request *pb.GetOrderStatusRequest, stream pb.OrderService_SubscribeOrderStatusServer) error {
	requestDto, err := dto.NewGetOrderStatusRequest(request)
	if err != nil {
		return status.Error(codes.Code(err.Code), err.Message)
	}

	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	ch, err := s.srvs.OrderService.SubscribeOrderStatus(ctx, requestDto)
	if err != nil {
		return status.Error(codes.Code(err.Code), err.Message)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case order, ok := <-ch:
			if !ok {
				return nil
			}

			err := stream.Send(&pb.GetOrderStatusResponse{Status: order.DtoToProto().Status})
			if err != nil {
				return err
			}
		}
	}
}

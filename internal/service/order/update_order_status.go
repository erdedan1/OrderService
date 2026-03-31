package order

import (
	"OrderService/internal/model"
	"context"

	errs "OrderService/internal/errors"

	errors "github.com/erdedan1/shared/errs"
	"github.com/google/uuid"
)

func (s *Service) UpdateOrderStatus(ctx context.Context, userID, orderID uuid.UUID, status model.OrderStatus) *errors.CustomError {
	const method = "UpdateOrderStatus"

	order, err := s.orderRepo.GetOrder(ctx, orderID, userID)
	if err != nil {
		s.log.Error(layer, method, err.Error(), err, "order_id", orderID, "status", status)
		return err
	}

	if order.Status == status || order.UserUUID != userID {
		return errs.ErrInvalidArgument
	}

	if updateErr := s.orderRepo.UpdateOrderStatus(ctx, orderID, status); updateErr != nil {
		s.log.Error(layer, method, updateErr.Error(), updateErr, "order_id", orderID, "status", status)
		return updateErr
	}

	if s.orderStatusPublisher != nil {
		if publishErr := s.orderStatusPublisher.PublishOrderStatus(ctx, orderID, status); publishErr != nil {
			s.log.Error(layer, method, publishErr.Error(), publishErr, "order_id", orderID, "status", status)
			return publishErr
		}
	}

	return nil
}

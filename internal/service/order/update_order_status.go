package order

import (
	"OrderService/internal/model"
	"context"

	errors "github.com/erdedan1/shared/errs"
	"github.com/google/uuid"
)

func (s *Service) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status model.OrderStatus) *errors.CustomError {
	const method = "UpdateOrderStatus"

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

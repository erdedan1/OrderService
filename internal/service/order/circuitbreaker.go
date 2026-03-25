package order

import (
	"context"
	"sync"
	"time"

	"OrderService/internal/dto"
	"OrderService/internal/usecase"

	sharedErrs "github.com/erdedan1/shared/errs"
)

type breakerState string

const (
	breakerStateClosed   breakerState = "closed"
	breakerStateOpen     breakerState = "open"
	breakerStateHalfOpen breakerState = "half-open"
)

type CircuitBreaker struct {
	next usecase.OrderService

	mu                  sync.Mutex
	state               breakerState
	consecutiveFailures uint32
	openedAt            time.Time
	halfOpenAttempts    uint32

	failureThreshold uint32
	halfOpenMaxCalls uint32
	openTimeout      time.Duration
}

func NewCircuitBreaker(next usecase.OrderService, consecutiveFailures uint32, halfOpenRequests uint32, openTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		next:             next,
		state:            breakerStateClosed,
		failureThreshold: consecutiveFailures,
		halfOpenMaxCalls: halfOpenRequests,
		openTimeout:      openTimeout,
	}
}

func (s *CircuitBreaker) CreateOrder(ctx context.Context, request *dto.CreateOrderRequest) (*dto.CreateOrderResponse, *sharedErrs.CustomError) {
	if err := s.beforeRequest(); err != nil {
		return nil, err
	}

	response, err := s.next.CreateOrder(ctx, request)
	s.afterRequest(err)
	return response, err
}

func (s *CircuitBreaker) GetOrderStatus(ctx context.Context, request *dto.GetOrderStatusRequest) (*dto.GetOrderStatusResponse, *sharedErrs.CustomError) {
	if err := s.beforeRequest(); err != nil {
		return nil, err
	}

	response, err := s.next.GetOrderStatus(ctx, request)
	s.afterRequest(err)
	return response, err
}

func (s *CircuitBreaker) SubscribeOrderStatus(ctx context.Context, request *dto.GetOrderStatusRequest) (<-chan *dto.GetOrderStatusResponse, *sharedErrs.CustomError) {
	if err := s.beforeRequest(); err != nil {
		return nil, err
	}

	response, err := s.next.SubscribeOrderStatus(ctx, request)
	s.afterRequest(err)
	return response, err
}

func (s *CircuitBreaker) beforeRequest() *sharedErrs.CustomError {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	if s.state == breakerStateOpen {
		if now.Sub(s.openedAt) < s.openTimeout {
			return sharedErrs.New(sharedErrs.UNAVAILABLE, "circuit breaker is open")
		}

		s.state = breakerStateHalfOpen
		s.halfOpenAttempts = 0
	}

	if s.state == breakerStateHalfOpen {
		if s.halfOpenAttempts >= s.halfOpenMaxCalls {
			return sharedErrs.New(sharedErrs.UNAVAILABLE, "circuit breaker is open")
		}
		s.halfOpenAttempts++
	}

	return nil
}

func (s *CircuitBreaker) afterRequest(err *sharedErrs.CustomError) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err == nil || !shouldCountFailure(err) {
		s.consecutiveFailures = 0
		if s.state == breakerStateHalfOpen {
			s.state = breakerStateClosed
			s.halfOpenAttempts = 0
		}
		return
	}

	if s.state == breakerStateHalfOpen {
		s.tripOpen()
		return
	}

	s.consecutiveFailures++
	if s.consecutiveFailures >= s.failureThreshold {
		s.tripOpen()
	}
}

func (s *CircuitBreaker) tripOpen() {
	s.state = breakerStateOpen
	s.openedAt = time.Now()
	s.consecutiveFailures = 0
	s.halfOpenAttempts = 0
}

func shouldCountFailure(err *sharedErrs.CustomError) bool {
	return err.Code == sharedErrs.INTERNAL ||
		err.Code == sharedErrs.UNAVAILABLE ||
		err.Code == sharedErrs.DEADLINE_EXCEEDED
}

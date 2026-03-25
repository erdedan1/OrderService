package order

import (
	"context"
	"sync"
	"time"

	"OrderService/internal/dto"
	"OrderService/internal/usecase"

	"github.com/erdedan1/shared/errs"
)

type RateLimiter struct {
	next usecase.OrderService

	mu               sync.Mutex
	requestsPerSec   float64
	burst            float64
	availableTokens  float64
	lastRefillAtTime time.Time
}

func NewRateLimiter(next usecase.OrderService, requestsPerSecond float64, burst int) *RateLimiter {
	now := time.Now()

	return &RateLimiter{
		next:             next,
		requestsPerSec:   requestsPerSecond,
		burst:            float64(burst),
		availableTokens:  float64(burst),
		lastRefillAtTime: now,
	}
}

func (s *RateLimiter) CreateOrder(ctx context.Context, request *dto.CreateOrderRequest) (*dto.CreateOrderResponse, *errs.CustomError) {
	if !s.allow() {
		return nil, errs.New(errs.RESOURCE_EXHAUSTED, "rate limit exceeded")
	}

	return s.next.CreateOrder(ctx, request)
}

func (s *RateLimiter) GetOrderStatus(ctx context.Context, request *dto.GetOrderStatusRequest) (*dto.GetOrderStatusResponse, *errs.CustomError) {
	if !s.allow() {
		return nil, errs.New(errs.RESOURCE_EXHAUSTED, "rate limit exceeded")
	}

	return s.next.GetOrderStatus(ctx, request)
}

func (s *RateLimiter) SubscribeOrderStatus(ctx context.Context, request *dto.GetOrderStatusRequest) (<-chan *dto.GetOrderStatusResponse, *errs.CustomError) {
	if !s.allow() {
		return nil, errs.New(errs.RESOURCE_EXHAUSTED, "rate limit exceeded")
	}

	return s.next.SubscribeOrderStatus(ctx, request)
}

func (s *RateLimiter) allow() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	elapsedSeconds := now.Sub(s.lastRefillAtTime).Seconds()
	s.lastRefillAtTime = now

	s.availableTokens = min(s.burst, s.availableTokens+elapsedSeconds*s.requestsPerSec)
	if s.availableTokens < 1 {
		return false
	}

	s.availableTokens--
	return true
}

package order_service

import (
	"context"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type breakerState string

const (
	breakerStateClosed   breakerState = "closed"
	breakerStateOpen     breakerState = "open"
	breakerStateHalfOpen breakerState = "half-open"
)

type grpcCircuitBreaker struct {
	mu sync.Mutex

	state               breakerState
	consecutiveFailures uint32
	openedAt            time.Time
	halfOpenAttempts    uint32

	failureThreshold uint32
	halfOpenMaxCalls uint32
	openTimeout      time.Duration
}

func newGRPCCircuitBreaker(consecutiveFailures uint32, halfOpenRequests uint32, openTimeout time.Duration) *grpcCircuitBreaker {
	return &grpcCircuitBreaker{
		state:            breakerStateClosed,
		failureThreshold: consecutiveFailures,
		halfOpenMaxCalls: halfOpenRequests,
		openTimeout:      openTimeout,
	}
}

func (b *grpcCircuitBreaker) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if err := b.beforeRequest(); err != nil {
			return nil, err
		}

		response, err := handler(ctx, req)
		b.afterRequest(err)
		return response, err
	}
}

func (b *grpcCircuitBreaker) Stream() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if err := b.beforeRequest(); err != nil {
			return err
		}

		err := handler(srv, ss)
		b.afterRequest(err)
		return err
	}
}

func (b *grpcCircuitBreaker) beforeRequest() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	if b.state == breakerStateOpen {
		if now.Sub(b.openedAt) < b.openTimeout {
			return status.Error(codes.Unavailable, "cycle breaker is open")
		}

		b.state = breakerStateHalfOpen
		b.halfOpenAttempts = 0
	}

	if b.state == breakerStateHalfOpen {
		if b.halfOpenAttempts >= b.halfOpenMaxCalls {
			return status.Error(codes.Unavailable, "cycle breaker is open")
		}
		b.halfOpenAttempts++
	}

	return nil
}

func (b *grpcCircuitBreaker) afterRequest(err error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if err == nil || !shouldCountFailure(err) {
		b.consecutiveFailures = 0
		if b.state == breakerStateHalfOpen {
			b.state = breakerStateClosed
			b.halfOpenAttempts = 0
		}
		return
	}

	if b.state == breakerStateHalfOpen {
		b.tripOpen()
		return
	}

	b.consecutiveFailures++
	if b.consecutiveFailures >= b.failureThreshold {
		b.tripOpen()
	}
}

func (b *grpcCircuitBreaker) tripOpen() {
	b.state = breakerStateOpen
	b.openedAt = time.Now()
	b.consecutiveFailures = 0
	b.halfOpenAttempts = 0
}

func shouldCountFailure(err error) bool {
	st, ok := status.FromError(err)
	if !ok {
		return true
	}

	switch st.Code() {
	case codes.Unavailable, codes.Internal, codes.DeadlineExceeded:
		return true
	}

	return strings.Contains(strings.ToLower(st.Message()), "transport is closing")
}

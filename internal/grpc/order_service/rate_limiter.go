package order_service

import (
	"context"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcRateLimiter struct {
	mu               sync.Mutex
	requestsPerSec   float64
	burst            float64
	availableTokens  float64
	lastRefillAtTime time.Time
}

func newGRPCRateLimiter(requestsPerSecond float64, burst int) *grpcRateLimiter {
	now := time.Now()
	return &grpcRateLimiter{
		requestsPerSec:   requestsPerSecond,
		burst:            float64(burst),
		availableTokens:  float64(burst),
		lastRefillAtTime: now,
	}
}

func (l *grpcRateLimiter) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if !l.allow() {
			return nil, status.Error(codes.ResourceExhausted, "rate limit exceeded")
		}
		return handler(ctx, req)
	}
}

func (l *grpcRateLimiter) Stream() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if !l.allow() {
			return status.Error(codes.ResourceExhausted, "rate limit exceeded")
		}
		return handler(srv, ss)
	}
}

func (l *grpcRateLimiter) allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	elapsedSeconds := now.Sub(l.lastRefillAtTime).Seconds()
	l.lastRefillAtTime = now

	l.availableTokens = min(l.burst, l.availableTokens+elapsedSeconds*l.requestsPerSec)
	if l.availableTokens < 1 {
		return false
	}

	l.availableTokens--
	return true
}

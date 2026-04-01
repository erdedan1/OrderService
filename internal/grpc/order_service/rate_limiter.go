package order_service

import (
	"context"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

type rateBucket struct {
	availableTokens  float64
	lastRefillAtTime time.Time
}

type grpcRateLimiter struct {
	mu             sync.Mutex
	requestsPerSec float64
	burst          float64
	globalBucket   rateBucket
	clients        map[string]*rateBucket
}

func newGRPCRateLimiter(requestsPerSecond float64, burst int) *grpcRateLimiter {
	now := time.Now()
	return &grpcRateLimiter{
		requestsPerSec: requestsPerSecond,
		burst:          float64(burst),
		globalBucket: rateBucket{
			availableTokens:  float64(burst),
			lastRefillAtTime: now,
		},
		clients: make(map[string]*rateBucket),
	}
}

func (l *grpcRateLimiter) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if !l.allow(clientKeyFromContext(ctx)) {
			return nil, status.Error(codes.ResourceExhausted, "rate limit exceeded")
		}
		return handler(ctx, req)
	}
}

func (l *grpcRateLimiter) Stream() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if !l.allow(clientKeyFromContext(ss.Context())) {
			return status.Error(codes.ResourceExhausted, "rate limit exceeded")
		}
		return handler(srv, ss)
	}
}

func (l *grpcRateLimiter) allow(clientKey string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	if !allowBucket(&l.globalBucket, now, l.requestsPerSec, l.burst) {
		return false
	}
	bucket, exists := l.clients[clientKey]
	if !exists {
		bucket = &rateBucket{
			availableTokens:  l.burst,
			lastRefillAtTime: now,
		}
		l.clients[clientKey] = bucket
	}

	return allowBucket(bucket, now, l.requestsPerSec, l.burst)
}

func allowBucket(bucket *rateBucket, now time.Time, requestPerSec, burst float64) bool {
	elapsedSecond := now.Sub(bucket.lastRefillAtTime).Seconds()
	bucket.lastRefillAtTime = now
	bucket.availableTokens = min(burst, bucket.availableTokens+elapsedSecond*requestPerSec)

	if bucket.availableTokens < 1 {
		return false
	}
	bucket.availableTokens--
	return true
}

func clientKeyFromContext(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if clientID := firstMetadataValue(md, "x-client-id"); clientID != "" {
			return "client-id" + clientID
		}
		if userUUID := firstMetadataValue(md, "x-user-uuid"); userUUID != "" {
			return "user-uuid" + userUUID
		}
	}

	if p, ok := peer.FromContext(ctx); ok && p.Addr != nil {
		host, _, err := net.SplitHostPort(p.Addr.String())
		if err == nil && host != "" {
			return "peer-host:" + host
		}
		return "peer-addr:" + p.Addr.String()
	}
	return "unknown-client"
}

func firstMetadataValue(md metadata.MD, key string) string {
	values := md.Get(key)
	if len(values) == 0 {
		return ""
	}

	return values[0]
}

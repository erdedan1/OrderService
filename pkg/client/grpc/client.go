package grpc_client

import (
	"context"

	"OrderService/config"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

type IGRPCClient interface {
	CanonicalTarget() string
	Close() error
	Connect()
	GetState() connectivity.State
	Invoke(ctx context.Context, method string, args any, reply any, opts ...grpc.CallOption) error
	NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error)
	ResetConnectBackoff()
	Target() string
	WaitForStateChange(ctx context.Context, sourceState connectivity.State) bool
}

func New(address string, cfg *config.Config, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	gGRPCopts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  cfg.GRPCClient.BaseBackoffDelay,
				Multiplier: cfg.GRPCClient.BackoffMultiplier,
				Jitter:     cfg.GRPCClient.BackoffJitter,
				MaxDelay:   cfg.GRPCClient.MaxBackoffDelay,
			},
			MinConnectTimeout: cfg.GRPCClient.ConnectTimeout,
		}),
	}
	gGRPCopts = append(gGRPCopts, opts...)

	return grpc.NewClient(address, gGRPCopts...)
}

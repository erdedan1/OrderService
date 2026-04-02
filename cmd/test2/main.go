package main

import (
	"OrderService/config"
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/erdedan1/protocol/proto/order_service/gen/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type probeResult struct {
	total             int
	resourceExhausted int
	internalErrors    int
	unavailableErrors int
	invalidArgErrors  int
	unknownErrors     int
}

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("read config: %v", err)
	}

	addr := cfg.GRPCServer.Address

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("connect to %s: %v", addr, err)
	}
	defer conn.Close()

	client := pb.NewOrderServiceClient(conn)
	request := &pb.GetOrderStatusRequest{
		UserUuid:  "1179803e-06f0-4369-b94f-14e26ec190a3",
		OrderUuid: "68aadf2e-111c-4743-87a1-83a8802999b2",
	}

	ctx := metadata.AppendToOutgoingContext(
		context.Background(),
		"x-user-uuid", "1179803e-06f0-4369-b94f-14e26ec190a3",
	)

	for i := 0; i < 30; i++ {
		ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
		res, err := client.GetOrderStatus(ctx, request)
		cancel()
		fmt.Println(err)
		fmt.Println(res)
	}

	ctx = metadata.AppendToOutgoingContext(
		context.Background(),
		"x-user-uuid", "1179803e-06f0-4369-b94f-14e26ec190a1",
	)

	for i := 0; i < 30; i++ {
		ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
		res, err := client.GetOrderStatus(ctx, request)
		cancel()
		fmt.Println(err)
		fmt.Println(res)
	}
}

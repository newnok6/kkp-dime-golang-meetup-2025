package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "github.com/newnok6/kkp-dime-golang-meetup-2025/backend/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := pb.NewStockOrderServiceClient(conn)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("\nReceived shutdown signal, canceling requests...")
		cancel()
	}()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			request := &pb.CreateOrderRequest{
				Symbol:    "AAPL",
				OrderType: pb.OrderType_MARKET,
				OrderSide: pb.OrderSide_BUY,
				Quantity:  100,
			}

			response, err := client.CreateOrder(ctx, request)
			if err != nil {
				if ctx.Err() == context.Canceled {
					log.Println("Request canceled, shutting down...")
					return
				}
				log.Printf("Failed to create order: %v\n", err)
				continue
			}

			log.Printf("Response: %v\n", response)

		case <-ctx.Done():
			log.Println("Shutting down gracefully...")
			return
		}
	}
}

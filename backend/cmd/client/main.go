package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Handle signals in a separate goroutine
	go func() {
		<-quit
		fmt.Println("\nReceived shutdown signal, canceling requests...")
		cancel()
	}()

	for {
		request, err := http.NewRequest("POST", "http://localhost:8082/api/orders", bytes.NewBuffer([]byte(`{
			"symbol": "AAPL",
			"order_type": "MARKET",
			"order_side": "BUY",
			"quantity": 100
		}`)))
		if err != nil {
			fmt.Printf("Failed to create request: %v\n", err)
			continue
		}
		request.Header.Set("Content-Type", "application/json")

		// Attach context to request so it can be canceled
		request = request.WithContext(ctx)

		response, err := client.Do(request)
		if err != nil {
			// Check if error is due to context cancellation
			if ctx.Err() == context.Canceled {
				fmt.Println("Request canceled, shutting down...")
				return
			}
			fmt.Printf("Failed to send request: %v\n", err)
			continue
		}

		fmt.Printf("Response: %v\n", response)
		response.Body.Close()

		select {
		case <-time.After(1 * time.Second):
			continue
		case <-ctx.Done():
			fmt.Println("Shutting down gracefully...")
			return
		}
	}
}

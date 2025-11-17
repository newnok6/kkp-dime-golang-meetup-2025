package main

import (
	"context"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/newnok6/kkp-dime-golang-meetup-2025/backend/adaptor"
	pb "github.com/newnok6/kkp-dime-golang-meetup-2025/backend/proto"
	"github.com/newnok6/kkp-dime-golang-meetup-2025/backend/service"
	"google.golang.org/grpc"
)

func main() {
	// Initialize SQLite repository
	repo, err := adaptor.NewSQLiteRepository("./stock_orders.db")
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}

	// Initialize service
	stockService := service.NewStockOrderService(repo)

	// Initialize HTTP handler
	httpHandler := adaptor.NewHTTPHandler(stockService)

	// Setup HTTP router
	router := mux.NewRouter()
	httpHandler.RegisterRoutes(router)

	// Add logging middleware
	router.Use(loggingMiddleware)

	// Start HTTP server
	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8082"
	}

	httpServer := &http.Server{
		Addr:    ":" + httpPort,
		Handler: router,
	}

	// HTTP Server startup
	go func() {
		log.Printf("---Starting HTTP server on port %s---", httpPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP Server failed: %v", err)
		}
	}()

	// Initialize gRPC handler
	grpcHandler := adaptor.NewGRPCHandler(stockService)

	// Setup gRPC server
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen on gRPC port: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterStockOrderServiceServer(grpcServer, grpcHandler)

	// gRPC Server startup
	go func() {
		log.Printf("---Starting gRPC server on port %s---", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC Server failed: %v", err)
		}
	}()

	// Graceful shutdown implementation
	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("---Start Graceful shutdown---")

	// Create context with timeout for graceful shutdown operations
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		log.Println("Shutting down HTTP server...")
		httpServer.Shutdown(shutdownCtx)
		log.Println("HTTP server stopped")
		wg.Done()
	}()

	// gRPC Server - use graceful stop
	go func() {
		log.Println("Shutting down gRPC server...")
		grpcServer.GracefulStop()
		log.Println("gRPC server stopped")
		wg.Done()
	}()
	wg.Wait()
	log.Println("All servers closed successfully")

	//--------------------------------

	// Other dependencies shutdown list
	shutDownList := []interface{}{}

	// Repository
	repositoryMap := map[string]interface{}{}
	repositoryMap["sqlite"] = repo

	shutDownList = append(shutDownList, repositoryMap)

	// Shutdown other dependencies
	log.Println("Shutting down other dependencies...")
	<-shutdownService(&shutDownList)
	log.Println("All Dependencies closed successfully")

	log.Println("---Graceful shutdown completed---")
	os.Exit(0)
}

func shutdownService(shutDownList *[]interface{}) <-chan struct{} {
	blockExitChannel := make(chan struct{})
	// Use WaitGroup to shut down server and database concurrently
	wg := sync.WaitGroup{}
	for i := range *shutDownList {
		serviceMap := (*shutDownList)[i].(map[string]interface{})
		for k, v := range serviceMap {
			wg.Add(1)
			go func(key string, closer interface{}) {
				defer wg.Done()
				time.Sleep(1 * time.Second)
				if c, ok := closer.(io.Closer); ok {
					if err := c.Close(); err != nil {
						log.Printf("Error closing service %s: %v", key, err)
					}
					log.Printf("%s closed successfully", key)
				}
			}(k, v)
		}
	}

	wg.Wait()
	close(blockExitChannel)
	return blockExitChannel
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/newnok6/kkp-dime-golang-meetup-2025/backend/adaptor"
	"github.com/newnok6/kkp-dime-golang-meetup-2025/backend/service"
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
	handler := adaptor.NewHTTPHandler(stockService)

	// Setup router
	router := mux.NewRouter()
	handler.RegisterRoutes(router)

	// Add logging middleware
	router.Use(loggingMiddleware)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Server startup
	go func() {
		log.Printf("---Starting server on port %s---", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Graceful shutdown implementation
	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("---Stop the services---")

	// Create context with timeout for graceful shutdown operations
	_, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	shutdownList := []interface{}{}
	// Server
	serverMap := map[string]interface{}{}
	serverMap["server"] = server
	// Repository
	repositoryMap := map[string]interface{}{}
	repositoryMap["repository"] = repo

	shutdownList = append(shutdownList, serverMap)
	shutdownList = append(shutdownList, repositoryMap)

	// Shared lib function for shutdown
	<-shutdownService(&shutdownList)
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
					log.Printf("Service %s closed successfully", key)
				}
			}(k, v)
		}
	}

	wg.Wait()

	log.Println("---Graceful shutdown completed")
	close(blockExitChannel)
	return blockExitChannel
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

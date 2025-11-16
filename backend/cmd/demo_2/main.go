package main

import (
	"context"
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
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Use WaitGroup to shut down server and database concurrently
	wg := sync.WaitGroup{}
	wg.Go(func() {
		defer wg.Done()
		log.Println("Closing server connection...")
		for i := 0; i < 5; {
			log.Printf("Doing Somthing %v \n", i)
			time.Sleep(1 * time.Second)
			i++
		}
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("Error shutting down server: %v", err)
		} else {
			log.Println("Server shut down successfully")
		}
	})
	wg.Go(func() {
		defer wg.Done()
		log.Println("Closing database connection...")
		time.Sleep(2 * time.Second)
		if err := repo.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		} else {
			log.Println("Database connection closed successfully")
		}
	})

	wg.Wait()

	log.Println("---Graceful shutdown completed")

	os.Exit(0)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

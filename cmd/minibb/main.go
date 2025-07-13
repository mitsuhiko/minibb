package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"minibb/internal/db"
	"minibb/internal/server"
)

func main() {
	// Initialize database
	db, err := db.Init()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Create server
	srv := server.New(db, nil)

	// Start server
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Println("Shutting down server...")
		cancel()
	}()

	log.Println("Starting MiniBB server...")
	if err := srv.Start(ctx); err != nil {
		log.Fatal("Server error:", err)
	}
}

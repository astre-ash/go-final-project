package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"go-final-progect/internal/server"
	"go-final-progect/internal/storage/sqlite"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}

	appPath, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
	}

	rootDir := filepath.Dir(appPath)
	webDir := filepath.Join(rootDir, "web")

	if _, err := os.Stat(webDir); os.IsNotExist(err) {
		log.Println("Using relative path for the web directory")
		webDir = "./web"
	}

	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	store, err := sqlite.New(dbFile)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer func() {
		if err := store.Close(); err != nil {
			log.Printf("Error closing database connection: %v\n", err)
		} else {
			log.Println("Database connection closed successfully")
		}
	}()

	log.Println("Database initialized successfully")

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	appServer := server.NewServer(port, webDir, store)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		log.Printf("Server starting on port %s", port)
		if err := appServer.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server ListenAndServe error: %v", err)
		}
	}()

	sig := <-stop
	log.Printf("Received signal: %v. Starting graceful shutdown...", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := appServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	} else {
		log.Println("Server stopped gracefully")
	}
}

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/codyseavey/3d-printer/backend/internal/api"
	"github.com/codyseavey/3d-printer/backend/internal/handlers"
)

func main() {
	// Configure handlers from env
	timelapseDir := os.Getenv("TIMELAPSE_DIR")
	if timelapseDir == "" {
		timelapseDir = "./videos"
	}
	handlers.InitTimelapseHandler(timelapseDir)

	streamPath := os.Getenv("STREAM_M3U8_PATH")
	if streamPath == "" {
		streamPath = "./live/stream.m3u8"
	}
	handlers.InitStreamHandler(streamPath)

	router := api.SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		log.Printf("Starting server on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

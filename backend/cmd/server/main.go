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
	timelapseDir := os.Getenv("TIMELAPSE_DIR")
	if timelapseDir == "" {
		timelapseDir = "./videos"
	}
	if _, err := os.Stat(timelapseDir); err != nil {
		log.Printf("WARNING: timelapse directory %s does not exist: %v", timelapseDir, err)
	}

	streamPath := os.Getenv("STREAM_M3U8_PATH")
	if streamPath == "" {
		streamPath = "./live/stream.m3u8"
	}

	timelapse := handlers.NewTimelapseHandler(timelapseDir)
	stream := handlers.NewStreamHandler(streamPath)
	router := api.SetupRouter(timelapse, stream)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
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

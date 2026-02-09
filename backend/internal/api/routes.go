package api

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/codyseavey/3d-printer/backend/internal/handlers"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	frontendPath := os.Getenv("FRONTEND_DIST_PATH")
	serveFrontend := frontendPath != "" && dirExists(frontendPath)

	// CORS config for local dev
	config := cors.DefaultConfig()
	if corsOrigins := os.Getenv("CORS_ALLOWED_ORIGINS"); corsOrigins != "" {
		config.AllowOrigins = strings.Split(corsOrigins, ",")
	} else {
		config.AllowOrigins = []string{"http://localhost:5173", "http://localhost:3000"}
	}
	config.AllowMethods = []string{"GET", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept"}
	router.Use(cors.New(config))

	// Health check
	router.GET("/health", handlers.Health)

	// API routes
	api := router.Group("/api")
	{
		api.GET("/timelapses", handlers.ListTimelapses)
		api.GET("/stream/status", handlers.StreamStatus)
	}

	// Serve frontend static files
	if serveFrontend {
		indexPath := filepath.Join(frontendPath, "index.html")

		router.Static("/assets", filepath.Join(frontendPath, "assets"))

		router.GET("/", func(c *gin.Context) {
			c.File(indexPath)
		})

		// SPA fallback
		router.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path

			if strings.HasPrefix(path, "/api") {
				c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
				return
			}

			c.File(indexPath)
		})
	}

	return router
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

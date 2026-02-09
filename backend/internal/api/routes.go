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

func SetupRouter(timelapse *handlers.TimelapseHandler, stream *handlers.StreamHandler) *gin.Engine {
	router := gin.Default()

	frontendPath := os.Getenv("FRONTEND_DIST_PATH")
	serveFrontend := frontendPath != "" && dirExists(frontendPath)

	config := cors.DefaultConfig()
	if corsOrigins := os.Getenv("CORS_ALLOWED_ORIGINS"); corsOrigins != "" {
		config.AllowOrigins = strings.Split(corsOrigins, ",")
	} else {
		config.AllowOrigins = []string{"http://localhost:5173", "http://localhost:3000"}
	}
	config.AllowMethods = []string{"GET", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept"}
	router.Use(cors.New(config))

	router.GET("/health", handlers.Health)

	apiGroup := router.Group("/api")
	{
		apiGroup.GET("/timelapses", timelapse.List)
		apiGroup.GET("/stream/status", stream.Status)
	}

	if serveFrontend {
		indexPath := filepath.Join(frontendPath, "index.html")

		router.Static("/assets", filepath.Join(frontendPath, "assets"))

		router.GET("/", func(c *gin.Context) {
			c.File(indexPath)
		})

		// SPA fallback (exclude API and nginx-served paths)
		router.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path

			if strings.HasPrefix(path, "/api") || strings.HasPrefix(path, "/live") || strings.HasPrefix(path, "/videos") {
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

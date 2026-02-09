package handlers

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/codyseavey/3d-printer/backend/internal/models"
)

type StreamHandler struct {
	m3u8Path string
}

func NewStreamHandler(m3u8Path string) *StreamHandler {
	return &StreamHandler{m3u8Path: m3u8Path}
}

func (h *StreamHandler) Status(c *gin.Context) {
	info, err := os.Stat(h.m3u8Path)
	if err != nil {
		log.Printf("stream status: cannot stat %s: %v", h.m3u8Path, err)
		c.JSON(http.StatusOK, models.StreamStatus{
			Online:      false,
			LastUpdated: time.Time{},
		})
		return
	}

	mtime := info.ModTime()
	staleThreshold := 30 * time.Second
	online := time.Since(mtime) < staleThreshold

	c.JSON(http.StatusOK, models.StreamStatus{
		Online:      online,
		LastUpdated: mtime,
	})
}

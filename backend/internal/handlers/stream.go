package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/codyseavey/3d-printer/backend/internal/models"
)

var streamM3U8Path string

func InitStreamHandler(m3u8Path string) {
	streamM3U8Path = m3u8Path
}

func StreamStatus(c *gin.Context) {
	info, err := os.Stat(streamM3U8Path)
	if err != nil {
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

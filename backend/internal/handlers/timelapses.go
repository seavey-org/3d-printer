package handlers

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/codyseavey/3d-printer/backend/internal/models"
)

// filenameDateRegex matches filenames like video_2024-07-24_09-14-01.mp4
var filenameDateRegex = regexp.MustCompile(`(\d{4}-\d{2}-\d{2})_(\d{2}-\d{2}-\d{2})`)

type TimelapseHandler struct {
	dir string
}

func NewTimelapseHandler(dir string) *TimelapseHandler {
	return &TimelapseHandler{dir: dir}
}

func (h *TimelapseHandler) List(c *gin.Context) {
	entries, err := os.ReadDir(h.dir)
	if err != nil {
		log.Printf("timelapses: failed to read directory %s: %v", h.dir, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read timelapse directory"})
		return
	}

	// Build a set of thumbnail filenames for fast lookup
	thumbnails := make(map[string]bool)
	thumbDir := filepath.Join(h.dir, "thumbnail")
	thumbEntries, err := os.ReadDir(thumbDir)
	if err == nil {
		for _, e := range thumbEntries {
			if !e.IsDir() {
				thumbnails[e.Name()] = true
			}
		}
	}

	timelapses := make([]models.Timelapse, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		ext := strings.ToLower(filepath.Ext(name))
		if ext != ".mp4" && ext != ".mkv" && ext != ".avi" {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			log.Printf("timelapses: failed to stat %s: %v", name, err)
			continue
		}

		date := parseDateFromFilename(name)
		if date.IsZero() {
			date = info.ModTime()
		}

		// Match thumbnail: same base name with .jpg extension
		baseName := strings.TrimSuffix(name, ext)
		thumbName := baseName + ".jpg"
		thumbURL := ""
		if thumbnails[thumbName] {
			thumbURL = "/videos/thumbnail/" + url.PathEscape(thumbName)
		}

		timelapses = append(timelapses, models.Timelapse{
			Filename:     name,
			URL:          "/videos/" + url.PathEscape(name),
			ThumbnailURL: thumbURL,
			Size:         info.Size(),
			Date:         date,
		})
	}

	// Sort by date, newest first
	sort.Slice(timelapses, func(i, j int) bool {
		return timelapses[i].Date.After(timelapses[j].Date)
	})

	c.JSON(http.StatusOK, timelapses)
}

func parseDateFromFilename(name string) time.Time {
	matches := filenameDateRegex.FindStringSubmatch(name)
	if len(matches) < 3 {
		return time.Time{}
	}

	// Convert "09-14-01" to "09:14:01"
	timePart := strings.ReplaceAll(matches[2], "-", ":")
	dateStr := matches[1] + "T" + timePart + "Z"

	t, err := time.Parse("2006-01-02T15:04:05Z", dateStr)
	if err != nil {
		return time.Time{}
	}
	return t
}

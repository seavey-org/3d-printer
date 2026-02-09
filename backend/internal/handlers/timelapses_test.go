package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/codyseavey/3d-printer/backend/internal/models"
)

func TestListTimelapses(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create temp directory structure
	tmpDir := t.TempDir()
	thumbDir := filepath.Join(tmpDir, "thumbnail")
	if err := os.Mkdir(thumbDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Create test video files
	testFiles := []struct {
		name string
		size int
	}{
		{"video_2024-07-24_09-14-01.mp4", 1000},
		{"video_2024-08-15_14-30-00.mp4", 2000},
		{"video_2024-06-01_08-00-00.mp4", 50}, // tiny file
		{"not-a-video.txt", 100},              // should be ignored
	}

	for _, f := range testFiles {
		data := make([]byte, f.size)
		if err := os.WriteFile(filepath.Join(tmpDir, f.name), data, 0o644); err != nil {
			t.Fatal(err)
		}
	}

	// Create a thumbnail for the first video
	if err := os.WriteFile(filepath.Join(thumbDir, "video_2024-07-24_09-14-01.jpg"), []byte("thumb"), 0o644); err != nil {
		t.Fatal(err)
	}

	InitTimelapseHandler(tmpDir)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/timelapses", nil)

	ListTimelapses(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var result []models.Timelapse
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	// Should only include .mp4 files (3 of them), not .txt
	if len(result) != 3 {
		t.Fatalf("expected 3 timelapses, got %d", len(result))
	}

	// Should be sorted newest first (2024-08-15 > 2024-07-24 > 2024-06-01)
	if result[0].Filename != "video_2024-08-15_14-30-00.mp4" {
		t.Errorf("expected newest first, got %s", result[0].Filename)
	}
	if result[1].Filename != "video_2024-07-24_09-14-01.mp4" {
		t.Errorf("expected second newest, got %s", result[1].Filename)
	}

	// First video should have a thumbnail URL
	if result[1].ThumbnailURL != "/videos/thumbnail/video_2024-07-24_09-14-01.jpg" {
		t.Errorf("expected thumbnail URL, got %q", result[1].ThumbnailURL)
	}

	// Second video should have no thumbnail
	if result[0].ThumbnailURL != "" {
		t.Errorf("expected empty thumbnail URL, got %q", result[0].ThumbnailURL)
	}

	// Check URL format
	if result[0].URL != "/videos/video_2024-08-15_14-30-00.mp4" {
		t.Errorf("unexpected URL: %s", result[0].URL)
	}

	// Check sizes
	if result[0].Size != 2000 {
		t.Errorf("expected size 2000, got %d", result[0].Size)
	}
}

func TestListTimelapses_EmptyDir(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tmpDir := t.TempDir()
	InitTimelapseHandler(tmpDir)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/timelapses", nil)

	ListTimelapses(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	// Should return null (nil slice) or empty array, both are valid
	body := w.Body.String()
	if body != "null" && body != "[]" {
		t.Errorf("expected null or [], got %s", body)
	}
}

func TestListTimelapses_MissingDir(t *testing.T) {
	gin.SetMode(gin.TestMode)

	InitTimelapseHandler("/nonexistent/path")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/timelapses", nil)

	ListTimelapses(c)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}

func TestParseDateFromFilename(t *testing.T) {
	tests := []struct {
		name     string
		expected time.Time
	}{
		{"video_2024-07-24_09-14-01.mp4", time.Date(2024, 7, 24, 9, 14, 1, 0, time.UTC)},
		{"video_2024-12-31_23-59-59.mp4", time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)},
		{"random_file.mp4", time.Time{}},
		{"no_date_here.txt", time.Time{}},
	}

	for _, tt := range tests {
		result := parseDateFromFilename(tt.name)
		if !result.Equal(tt.expected) {
			t.Errorf("parseDateFromFilename(%q) = %v, want %v", tt.name, result, tt.expected)
		}
	}
}

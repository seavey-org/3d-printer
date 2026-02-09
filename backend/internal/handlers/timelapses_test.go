package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/codyseavey/3d-printer/backend/internal/models"
)

func TestListTimelapses(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tmpDir := t.TempDir()
	thumbDir := filepath.Join(tmpDir, "thumbnail")
	if err := os.Mkdir(thumbDir, 0o755); err != nil {
		t.Fatal(err)
	}

	testFiles := []struct {
		name string
		size int
	}{
		{"video_2024-07-24_09-14-01.mp4", 1000},
		{"video_2024-08-15_14-30-00.mp4", 2000},
		{"video_2024-06-01_08-00-00.mp4", 50},
		{"not-a-video.txt", 100},
	}

	for _, f := range testFiles {
		data := make([]byte, f.size)
		if err := os.WriteFile(filepath.Join(tmpDir, f.name), data, 0o644); err != nil {
			t.Fatal(err)
		}
	}

	if err := os.WriteFile(filepath.Join(thumbDir, "video_2024-07-24_09-14-01.jpg"), []byte("thumb"), 0o644); err != nil {
		t.Fatal(err)
	}

	h := NewTimelapseHandler(tmpDir)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/timelapses", nil)

	h.List(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var result []models.Timelapse
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(result) != 3 {
		t.Fatalf("expected 3 timelapses, got %d", len(result))
	}

	// Sorted newest first
	if result[0].Filename != "video_2024-08-15_14-30-00.mp4" {
		t.Errorf("expected newest first, got %s", result[0].Filename)
	}
	if result[1].Filename != "video_2024-07-24_09-14-01.mp4" {
		t.Errorf("expected second newest, got %s", result[1].Filename)
	}

	// Thumbnail URL should be URL-encoded
	expectedThumb := "/videos/thumbnail/" + url.PathEscape("video_2024-07-24_09-14-01.jpg")
	if result[1].ThumbnailURL != expectedThumb {
		t.Errorf("expected thumbnail URL %q, got %q", expectedThumb, result[1].ThumbnailURL)
	}

	if result[0].ThumbnailURL != "" {
		t.Errorf("expected empty thumbnail URL, got %q", result[0].ThumbnailURL)
	}

	// URL should be URL-encoded
	expectedURL := "/videos/" + url.PathEscape("video_2024-08-15_14-30-00.mp4")
	if result[0].URL != expectedURL {
		t.Errorf("unexpected URL: got %q, want %q", result[0].URL, expectedURL)
	}

	if result[0].Size != 2000 {
		t.Errorf("expected size 2000, got %d", result[0].Size)
	}
}

func TestListTimelapses_EmptyDir(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tmpDir := t.TempDir()
	h := NewTimelapseHandler(tmpDir)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/timelapses", nil)

	h.List(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	// Should return [] (empty array, not null)
	body := w.Body.String()
	if body != "[]" {
		t.Errorf("expected [], got %s", body)
	}
}

func TestListTimelapses_MissingDir(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := NewTimelapseHandler("/nonexistent/path")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/timelapses", nil)

	h.List(c)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}

func TestListTimelapses_MkvAndAvi(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tmpDir := t.TempDir()
	files := []string{"timelapse.mkv", "timelapse.avi", "timelapse.mp4", "readme.txt"}
	for _, f := range files {
		if err := os.WriteFile(filepath.Join(tmpDir, f), make([]byte, 500), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	h := NewTimelapseHandler(tmpDir)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/timelapses", nil)

	h.List(c)

	var result []models.Timelapse
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(result) != 3 {
		t.Fatalf("expected 3 timelapses (mp4, mkv, avi), got %d", len(result))
	}

	exts := map[string]bool{}
	for _, r := range result {
		ext := filepath.Ext(r.Filename)
		exts[ext] = true
	}
	for _, want := range []string{".mp4", ".mkv", ".avi"} {
		if !exts[want] {
			t.Errorf("expected %s file in results", want)
		}
	}
}

func TestListTimelapses_ThumbnailDirWithoutMatchingVideos(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tmpDir := t.TempDir()
	thumbDir := filepath.Join(tmpDir, "thumbnail")
	if err := os.Mkdir(thumbDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Orphaned thumbnail with no matching video
	if err := os.WriteFile(filepath.Join(thumbDir, "orphan.jpg"), []byte("thumb"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Video with no matching thumbnail
	if err := os.WriteFile(filepath.Join(tmpDir, "video_2024-01-01_00-00-00.mp4"), make([]byte, 500), 0o644); err != nil {
		t.Fatal(err)
	}

	h := NewTimelapseHandler(tmpDir)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/timelapses", nil)

	h.List(c)

	var result []models.Timelapse
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 timelapse, got %d", len(result))
	}
	if result[0].ThumbnailURL != "" {
		t.Errorf("expected empty thumbnail URL for unmatched video, got %q", result[0].ThumbnailURL)
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
		// Invalid date components: regex matches but time.Parse rejects
		{"video_2024-13-32_25-61-61.mp4", time.Time{}},
		// Multiple date patterns: should use the first match
		{"video_2024-01-01_00-00-00_2025-12-31_23-59-59.mp4", time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
		// Partial match
		{"video_2024-07-24.mp4", time.Time{}},
	}

	for _, tt := range tests {
		result := parseDateFromFilename(tt.name)
		if !result.Equal(tt.expected) {
			t.Errorf("parseDateFromFilename(%q) = %v, want %v", tt.name, result, tt.expected)
		}
	}
}

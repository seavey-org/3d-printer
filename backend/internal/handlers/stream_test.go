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

func TestStreamStatus_Online(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tmpDir := t.TempDir()
	m3u8Path := filepath.Join(tmpDir, "stream.m3u8")

	if err := os.WriteFile(m3u8Path, []byte("#EXTM3U\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	h := NewStreamHandler(m3u8Path)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/stream/status", nil)

	h.Status(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var status models.StreamStatus
	if err := json.Unmarshal(w.Body.Bytes(), &status); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if !status.Online {
		t.Error("expected stream to be online")
	}

	if status.LastUpdated.IsZero() {
		t.Error("expected non-zero last updated time")
	}
}

func TestStreamStatus_Offline_StaleFile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tmpDir := t.TempDir()
	m3u8Path := filepath.Join(tmpDir, "stream.m3u8")

	if err := os.WriteFile(m3u8Path, []byte("#EXTM3U\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	oldTime := time.Now().Add(-2 * time.Minute)
	if err := os.Chtimes(m3u8Path, oldTime, oldTime); err != nil {
		t.Fatal(err)
	}

	h := NewStreamHandler(m3u8Path)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/stream/status", nil)

	h.Status(c)

	var status models.StreamStatus
	if err := json.Unmarshal(w.Body.Bytes(), &status); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if status.Online {
		t.Error("expected stream to be offline for stale file")
	}
}

func TestStreamStatus_Offline_MissingFile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := NewStreamHandler("/nonexistent/stream.m3u8")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/stream/status", nil)

	h.Status(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var status models.StreamStatus
	if err := json.Unmarshal(w.Body.Bytes(), &status); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if status.Online {
		t.Error("expected stream to be offline for missing file")
	}
}

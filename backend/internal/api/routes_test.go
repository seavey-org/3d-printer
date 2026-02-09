package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/codyseavey/3d-printer/backend/internal/handlers"
	"github.com/codyseavey/3d-printer/backend/internal/models"
)

func setupTestRouter(t *testing.T) (*gin.Engine, string, string) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	timelapseDir := t.TempDir()
	streamDir := t.TempDir()
	m3u8Path := filepath.Join(streamDir, "stream.m3u8")

	if err := os.WriteFile(m3u8Path, []byte("#EXTM3U\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	timelapse := handlers.NewTimelapseHandler(timelapseDir)
	stream := handlers.NewStreamHandler(m3u8Path)
	router := SetupRouter(timelapse, stream)

	return router, timelapseDir, m3u8Path
}

func TestHealthRoute(t *testing.T) {
	router, _, _ := setupTestRouter(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if body["status"] != "ok" {
		t.Errorf("expected status ok, got %s", body["status"])
	}
}

func TestTimelapsesRoute(t *testing.T) {
	router, timelapseDir, _ := setupTestRouter(t)

	// Create a test video
	if err := os.WriteFile(filepath.Join(timelapseDir, "test.mp4"), make([]byte, 100), 0o644); err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/timelapses", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var result []models.Timelapse
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 timelapse, got %d", len(result))
	}
}

func TestStreamStatusRoute(t *testing.T) {
	router, _, _ := setupTestRouter(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/stream/status", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var status models.StreamStatus
	if err := json.Unmarshal(w.Body.Bytes(), &status); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if !status.Online {
		t.Error("expected stream to be online (fresh m3u8 file)")
	}
}

func TestAPINotFound(t *testing.T) {
	router, _, _ := setupTestRouter(t)

	// API routes that don't exist should still get a response (not SPA fallback)
	// Without frontend configured, this hits the default gin 404
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/nonexistent", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestSPAFallback(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Set up a fake frontend dist
	frontendDir := t.TempDir()
	assetsDir := filepath.Join(frontendDir, "assets")
	if err := os.Mkdir(assetsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(frontendDir, "index.html"), []byte("<html>SPA</html>"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(assetsDir, "app.js"), []byte("// app"), 0o644); err != nil {
		t.Fatal(err)
	}

	t.Setenv("FRONTEND_DIST_PATH", frontendDir)

	streamDir := t.TempDir()
	m3u8Path := filepath.Join(streamDir, "stream.m3u8")
	if err := os.WriteFile(m3u8Path, []byte("#EXTM3U\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	timelapse := handlers.NewTimelapseHandler(t.TempDir())
	stream := handlers.NewStreamHandler(m3u8Path)
	router := SetupRouter(timelapse, stream)

	// Root should serve index.html
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for /, got %d", w.Code)
	}
	if w.Body.String() != "<html>SPA</html>" {
		t.Errorf("expected SPA html, got %s", w.Body.String())
	}

	// Unknown path should fallback to index.html (SPA routing)
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/camera", nil)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for /camera, got %d", w.Code)
	}
	if w.Body.String() != "<html>SPA</html>" {
		t.Errorf("expected SPA html for SPA route, got %s", w.Body.String())
	}

	// Deep link to a specific timelapse should fallback to index.html
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/timelapses/video_2024-07-24_09-14-01.mp4", nil)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for /timelapses/<filename>, got %d", w.Code)
	}
	if w.Body.String() != "<html>SPA</html>" {
		t.Errorf("expected SPA html for timelapse deep link, got %s", w.Body.String())
	}

	// /api/* should NOT fallback to SPA
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/nonexistent", nil)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for unknown /api route, got %d", w.Code)
	}

	// Static assets should be served
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/assets/app.js", nil)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for static asset, got %d", w.Code)
	}
}

func TestCORSHeaders(t *testing.T) {
	router, _, _ := setupTestRouter(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodOptions, "/api/timelapses", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	req.Header.Set("Access-Control-Request-Method", "GET")
	router.ServeHTTP(w, req)

	origin := w.Header().Get("Access-Control-Allow-Origin")
	if origin != "http://localhost:5173" {
		t.Errorf("expected CORS origin http://localhost:5173, got %q", origin)
	}
}

func TestCORSRejectsUnknownOrigin(t *testing.T) {
	router, _, _ := setupTestRouter(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodOptions, "/api/timelapses", nil)
	req.Header.Set("Origin", "http://evil.com")
	req.Header.Set("Access-Control-Request-Method", "GET")
	router.ServeHTTP(w, req)

	origin := w.Header().Get("Access-Control-Allow-Origin")
	if origin == "http://evil.com" {
		t.Error("CORS should not allow unknown origins")
	}
}

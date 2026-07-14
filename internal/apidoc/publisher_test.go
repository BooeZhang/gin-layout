package apidoc

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"gin-layout/config"
)

func TestPublisher_HandleUI(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := config.APIDocConfig{Enabled: true}
	pub := NewPublisher(cfg, NewRegistry())

	r := gin.New()
	pub.Register(r)

	// 访问根 UI 路径应返回 Redoc HTML 页面。
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/docs/", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "<redoc") {
		t.Errorf("expected Redoc HTML tag in response")
	}
	if !strings.Contains(body, "redoc.standalone.js") {
		t.Errorf("expected Redoc JS script reference in response")
	}
	if ct := w.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
		t.Errorf("expected text/html content type, got %q", ct)
	}
}

func TestPublisher_HandleUI_JS(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := config.APIDocConfig{Enabled: true}
	pub := NewPublisher(cfg, NewRegistry())

	r := gin.New()
	pub.Register(r)

	// 访问嵌入的 JS 文件应返回 JavaScript。
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/docs/redoc.standalone.js", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/javascript") {
		t.Errorf("expected application/javascript content type, got %q", ct)
	}
	if len(w.Body.Bytes()) == 0 {
		t.Errorf("expected non-empty JS body")
	}
}

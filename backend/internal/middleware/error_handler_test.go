package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/blueship581/gbinsureapi/internal/util"
	"github.com/gin-gonic/gin"
)

func TestErrorHandlerNotFoundMappingP902(t *testing.T) {
	gin.SetMode(gin.TestMode)
	log := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	r := gin.New()
	r.Use(ErrorHandler(log))
	r.GET("/api/v1/test", func(c *gin.Context) {
		c.Error(fmt.Errorf("find presettlement: %w", util.ErrNotFound))
	})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/test", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404 for not-found error", rec.Code)
	}
}

var _ = util.ErrNotFound

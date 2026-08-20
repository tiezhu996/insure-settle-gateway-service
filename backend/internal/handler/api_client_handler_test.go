package handler

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/blueship581/gbinsureapi/internal/middleware"
	"github.com/blueship581/gbinsureapi/internal/model"
	"github.com/blueship581/gbinsureapi/internal/repository"
	"github.com/blueship581/gbinsureapi/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func newClientTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&model.ApiClient{}, &model.InsuredPerson{}, &model.UploadBatch{}, &model.FeeItem{}, &model.Presettlement{}, &model.SettlementOrder{}, &model.DailyReconciliation{}, &model.AuditLog{}); err != nil {
		t.Fatal(err)
	}
	return db
}

func TestApiClientHandlerDisableMissingP203(t *testing.T) {
	gin.SetMode(gin.TestMode)
	log := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	db := newClientTestDB(t)
	svc := service.NewApiClientService(repository.NewApiClientRepository(db), "s", "j", 24, log)
	h := NewApiClientHandler(svc, log)

	r := gin.New()
	r.Use(middleware.ErrorHandler(log))
	r.PUT("/api/v1/clients/:id/status", h.UpdateStatus)

	body := bytes.NewBufferString(`{"status":"disabled"}`)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/clients/9999/status", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404; body=%s", rec.Code, rec.Body.String())
	}
	if !bytes.Contains(rec.Body.Bytes(), []byte("不存在")) {
		t.Fatalf("body missing not-found message: %s", rec.Body.String())
	}
}

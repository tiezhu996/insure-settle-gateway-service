package handler

import (
	"context"
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

func TestReconciliationDailyCtxCancelP501(t *testing.T) {
	gin.SetMode(gin.TestMode)
	log := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&model.ApiClient{}, &model.InsuredPerson{}, &model.UploadBatch{}, &model.FeeItem{}, &model.Presettlement{}, &model.SettlementOrder{}, &model.DailyReconciliation{}, &model.AuditLog{}); err != nil {
		t.Fatal(err)
	}
	svc := service.NewReconciliationService(
		repository.NewSettlementOrderRepository(db),
		repository.NewDailyReconciliationRepository(db), log)
	h := NewDailyReconciliationHandler(svc, log)

	r := gin.New()
	r.Use(middleware.ErrorHandler(log))
	r.GET("/api/v1/reconciliations/daily", h.Daily)

	// 请求 context 已取消：对账应立刻失败，不能把已取消的 ctx 丢成后台执行
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/reconciliations/daily", nil).WithContext(ctx)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500 (ctx canceled must fail, got success); body=%s", rec.Code, rec.Body.String())
	}
}

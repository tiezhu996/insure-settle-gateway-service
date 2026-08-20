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
	"github.com/blueship581/gbinsureapi/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestSettleDetailBadToken401P803(t *testing.T) {
	gin.SetMode(gin.TestMode)
	log := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&model.ApiClient{}, &model.InsuredPerson{}, &model.UploadBatch{}, &model.FeeItem{}, &model.Presettlement{}, &model.SettlementOrder{}, &model.DailyReconciliation{}, &model.AuditLog{}); err != nil {
		t.Fatal(err)
	}
	insurance := service.NewInsuranceService(repository.NewInsuredPersonRepository(db), log)
	svc := service.NewSettlementService(
		repository.NewPresettlementRepository(db), repository.NewSettlementOrderRepository(db),
		repository.NewFeeItemRepository(db), repository.NewUploadBatchRepository(db),
		insurance, util.NewSettlementCalculator(), log)
	h := NewSettlementOrderHandler(svc, log)

	r := gin.New()
	r.Use(middleware.ErrorHandler(log))
	r.GET("/api/v1/settlements/:settlement_no", middleware.JwtAuthRequired("secret"), h.Detail)

	expired, err := util.GenerateToken("secret", -1, 1, "his", "settlement", "service")
	if err != nil {
		t.Fatalf("GenerateToken error = %v", err)
	}
	req := httptest.NewRequest(http.MethodGet, "/api/v1/settlements/S1", nil)
	req.Header.Set("Authorization", "Bearer "+expired)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401 for expired token; body=%s", rec.Code, rec.Body.String())
	}
}

func TestSettleSubmitMissingClientIDP805(t *testing.T) {
	gin.SetMode(gin.TestMode)
	log := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&model.ApiClient{}, &model.InsuredPerson{}, &model.UploadBatch{}, &model.FeeItem{}, &model.Presettlement{}, &model.SettlementOrder{}, &model.DailyReconciliation{}, &model.AuditLog{}); err != nil {
		t.Fatal(err)
	}
	insurance := service.NewInsuranceService(repository.NewInsuredPersonRepository(db), log)
	svc := service.NewSettlementService(
		repository.NewPresettlementRepository(db), repository.NewSettlementOrderRepository(db),
		repository.NewFeeItemRepository(db), repository.NewUploadBatchRepository(db),
		insurance, util.NewSettlementCalculator(), log)
	h := NewSettlementOrderHandler(svc, log)

	r := gin.New()
	r.Use(middleware.ErrorHandler(log))
	r.POST("/api/v1/settlements/submit", h.Submit)

	body := bytes.NewBufferString(`{"presettlement_id":1}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/settlements/submit", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401 when client id missing; body=%s", rec.Code, rec.Body.String())
	}
}

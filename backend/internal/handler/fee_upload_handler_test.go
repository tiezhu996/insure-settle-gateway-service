package handler

import (
	"bytes"
	"log/slog"
	"strconv"
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

func TestUploadBatchTxNoPartialP603(t *testing.T) {
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
	svc := service.NewFeeService(repository.NewUploadBatchRepository(db), repository.NewFeeItemRepository(db), insurance, log)
	h := NewFeeItemHandler(svc, log)

	r := gin.New()
	r.Use(middleware.ErrorHandler(log))
	r.POST("/api/v1/fees/upload", h.Upload)

	// 预置占用 BatchNo(2) 的批次（不同参保人），触发事务失败
	person := model.InsuredPerson{IDCardNo: "110101199001011234", MedicalCardNo: "M110101199001011234", Name: "张三", InsuranceType: "employee", InsuranceStatus: "active"}
	if err := db.Create(&person).Error; err != nil {
		t.Fatal(err)
	}
	other := model.InsuredPerson{IDCardNo: "310101198505052345", MedicalCardNo: "M310101198505052345", Name: "李四", InsuranceType: "resident", InsuranceStatus: "active"}
	if err := db.Create(&other).Error; err != nil {
		t.Fatal(err)
	}
	// 预置 BatchNo(2)（不同参保人，不触发当日重复检查），Upload 生成
	// util.BatchNo(Count()+1)=BatchNo(2) 时撞唯一索引触发事务失败
	if err := db.Create(&model.UploadBatch{BatchNo: util.BatchNo(2), ClientID: 1, InsuredPersonID: other.ID, UploadStatus: "validated"}).Error; err != nil {
		t.Fatal(err)
	}

	body := bytes.NewBufferString(`{"client_id":1,"insured_person_id":` + uintStr(person.ID) + `,"items":[{"item_code":"DRUG01","item_name":"阿莫西林","item_type":"drug","unit_price":50,"quantity":10,"medical_category":"class_a"}]}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/fees/upload", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code == http.StatusOK || rec.Code == http.StatusCreated {
		t.Fatalf("upload with conflicting batch_no should fail, got status %d body=%s", rec.Code, rec.Body.String())
	}
}

func uintStr(v uint) string {
	return strconv.FormatUint(uint64(v), 10)
}

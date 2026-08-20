package service

import (
	"context"
	"errors"
	"testing"

	"github.com/blueship581/gbinsureapi/internal/repository"
	"github.com/blueship581/gbinsureapi/internal/util"
)

func TestPresetMissingSubmit404P903(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()
	insurance := NewInsuranceService(repository.NewInsuredPersonRepository(db), testLogger())
	svc := NewSettlementService(
		repository.NewPresettlementRepository(db), repository.NewSettlementOrderRepository(db),
		repository.NewFeeItemRepository(db), repository.NewUploadBatchRepository(db),
		insurance, util.NewSettlementCalculator(), testLogger())
	_, err := svc.SubmitSettlement(ctx, 1, 9999)
	if err == nil {
		t.Fatal("expected error for missing presettlement")
	}
	var appErr *util.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got %T: %v", err, err)
	}
	if appErr.HTTPStatus != 404 {
		t.Fatalf("HTTPStatus = %d, want 404", appErr.HTTPStatus)
	}
}

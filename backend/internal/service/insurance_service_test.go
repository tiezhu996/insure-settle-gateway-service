package service

import (
	"context"
	"errors"
	"testing"

	"github.com/blueship581/gbinsureapi/internal/repository"
	"github.com/blueship581/gbinsureapi/internal/util"
)

func TestInsuredVerifyMissingNoPanicP401(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()
	svc := NewInsuranceService(repository.NewInsuredPersonRepository(db), testLogger())
	_, err := svc.Verify(ctx, "110101199001019999", "M999999999999999999")
	if err == nil {
		t.Fatal("expected error for unknown insured person")
	}
	var appErr *util.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got %T: %v", err, err)
	}
	if appErr.HTTPStatus != 404 {
		t.Fatalf("HTTPStatus = %d, want 404", appErr.HTTPStatus)
	}
}

func TestInsuranceGetByIDNilGuardP402(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()
	svc := NewInsuranceService(repository.NewInsuredPersonRepository(db), testLogger())
	_, err := svc.GetByID(ctx, 9999)
	if err == nil {
		t.Fatal("expected error for missing insured person")
	}
	if !errors.Is(err, util.ErrNotFound) {
		t.Fatalf("errors.Is(err, util.ErrNotFound) = false, got: %v", err)
	}
}

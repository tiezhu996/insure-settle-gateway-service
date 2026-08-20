package service

import (
	"context"
	"errors"
	"testing"

	"github.com/blueship581/gbinsureapi/internal/constants"
	"github.com/blueship581/gbinsureapi/internal/repository"
	"github.com/blueship581/gbinsureapi/internal/util"
)

func TestApiClientUpdateStatusMissingP201(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()
	svc := NewApiClientService(repository.NewApiClientRepository(db), "s", "j", 24, testLogger())
	err := svc.UpdateStatus(ctx, 9999, constants.ClientActive)
	if err == nil {
		t.Fatal("expected error for missing client")
	}
	var appErr *util.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got %T: %v", err, err)
	}
	if appErr.HTTPStatus != 404 {
		t.Fatalf("HTTPStatus = %d, want 404", appErr.HTTPStatus)
	}
}

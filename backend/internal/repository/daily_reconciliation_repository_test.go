package repository

import (
	"context"
	"testing"

	"github.com/blueship581/gbinsureapi/internal/model"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestReconRepoCtxPropagatedP503(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&model.DailyReconciliation{}); err != nil {
		t.Fatal(err)
	}
	repo := NewDailyReconciliationRepository(db)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	rec := &model.DailyReconciliation{ReconcileDate: "2026-08-20", TotalCount: 0}
	err = repo.Create(ctx, rec)
	if err == nil {
		t.Fatal("expected ctx-canceled error from repo.Create, got nil")
	}
}

package repository

import (
	"errors"
	"testing"

	"github.com/blueship581/gbinsureapi/internal/model"
	"github.com/blueship581/gbinsureapi/internal/util"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestPresetRepoChainErrNotFoundP901(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&model.Presettlement{}); err != nil {
		t.Fatal(err)
	}
	repo := NewPresettlementRepository(db)
	_, err = repo.FindByID(9999)
	if err == nil {
		t.Fatal("expected error for missing presettlement")
	}
	if !errors.Is(err, util.ErrNotFound) {
		t.Fatalf("errors.Is(err, util.ErrNotFound) = false, got: %v", err)
	}
}

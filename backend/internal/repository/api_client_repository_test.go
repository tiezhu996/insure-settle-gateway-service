package repository

import (
	"errors"
	"testing"

	"github.com/blueship581/gbinsureapi/internal/constants"
	"github.com/blueship581/gbinsureapi/internal/model"
	"github.com/blueship581/gbinsureapi/internal/util"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestApiClientRepoChainErrNotFoundP202(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&model.ApiClient{}); err != nil {
		t.Fatal(err)
	}
	repo := NewApiClientRepository(db)
	_, err = repo.FindByID(9999)
	if err == nil {
		t.Fatal("expected error for missing client")
	}
	if !errors.Is(err, util.ErrNotFound) {
		t.Fatalf("errors.Is(err, util.ErrNotFound) = false, got: %v", err)
	}
}

func TestApiClientRepoUpdateMissingP204(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&model.ApiClient{}); err != nil {
		t.Fatal(err)
	}
	repo := NewApiClientRepository(db)
	err = repo.UpdateStatus(9999, constants.ClientDisabled)
	if err == nil {
		t.Fatal("expected error when updating missing client status")
	}
	if !errors.Is(err, util.ErrNotFound) {
		t.Fatalf("errors.Is(err, util.ErrNotFound) = false, got: %v", err)
	}
}

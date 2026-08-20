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

func TestUploadBatchTxRollsBackP602(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&model.UploadBatch{}); err != nil {
		t.Fatal(err)
	}
	repo := NewUploadBatchRepository(db)
	boom := errors.New("boom")
	_ = repo.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&model.UploadBatch{BatchNo: util.BatchNo(1), ClientID: 1, InsuredPersonID: 1, UploadStatus: constants.UploadValidated}).Error; err != nil {
			return err
		}
		return boom
	})
	var count int64
	db.Model(&model.UploadBatch{}).Count(&count)
	if count != 0 {
		t.Fatalf("expected rollback (0 rows), got %d rows", count)
	}
}

func TestUploadBatchTxErrorPropagatedP603(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&model.UploadBatch{}); err != nil {
		t.Fatal(err)
	}
	repo := NewUploadBatchRepository(db)
	boom := errors.New("boom")
	err = repo.Transaction(func(tx *gorm.DB) error {
		return boom
	})
	if err == nil {
		t.Fatal("expected tx error to propagate, got nil")
	}
	if !errors.Is(err, boom) {
		t.Fatalf("err = %v, want boom", err)
	}
}

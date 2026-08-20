package repository

import (
	"testing"

	"github.com/blueship581/gbinsureapi/internal/model"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestFeeListByBatchIndependentP1001(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&model.FeeItem{}); err != nil {
		t.Fatal(err)
	}
	repo := NewFeeItemRepository(db)
	if err := repo.CreateBatch([]model.FeeItem{
		{BatchID: 1, ItemCode: "A", ItemName: "甲", ItemType: "drug", Amount: 100, MedicalCategory: "class_a"},
		{BatchID: 1, ItemCode: "B", ItemName: "乙", ItemType: "drug", Amount: 50, MedicalCategory: "class_a"},
	}); err != nil {
		t.Fatal(err)
	}
	if err := repo.CreateBatch([]model.FeeItem{
		{BatchID: 2, ItemCode: "C", ItemName: "丙", ItemType: "drug", Amount: 80, MedicalCategory: "class_a"},
	}); err != nil {
		t.Fatal(err)
	}
	first, err := repo.ListByBatch(1)
	if err != nil {
		t.Fatal(err)
	}
	want := make([]model.FeeItem, len(first))
	copy(want, first)
	// 查询另一个批次后，第一次返回的切片不得被内部缓冲覆盖
	_, err = repo.ListByBatch(2)
	if err != nil {
		t.Fatal(err)
	}
	if len(first) != len(want) {
		t.Fatalf("first len = %d, want %d (buffer reused?)", len(first), len(want))
	}
	for i := range want {
		if first[i].ItemCode != want[i].ItemCode {
			t.Fatalf("first result polluted at %d: got %+v want %+v", i, first[i], want[i])
		}
	}
}

func TestFeeCreateBatchPersistsAllP1004(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&model.FeeItem{}); err != nil {
		t.Fatal(err)
	}
	repo := NewFeeItemRepository(db)
	items := []model.FeeItem{
		{BatchID: 1, ItemCode: "A", ItemName: "甲", ItemType: "drug", Amount: 100, MedicalCategory: "class_a"},
		{BatchID: 1, ItemCode: "B", ItemName: "乙", ItemType: "drug", Amount: 50, MedicalCategory: "class_a"},
	}
	if err := repo.CreateBatch(items); err != nil {
		t.Fatal(err)
	}
	got, err := repo.ListByBatch(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != len(items) {
		t.Fatalf("persisted %d items, want %d (CreateBatch truncated input?)", len(got), len(items))
	}
	for i, it := range items {
		if got[i].ItemCode != it.ItemCode {
			t.Fatalf("item %d mismatch: got %s want %s", i, got[i].ItemCode, it.ItemCode)
		}
	}
}

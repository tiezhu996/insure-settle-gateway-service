package util

import (
	"testing"

	"github.com/blueship581/gbinsureapi/internal/model"
)

func TestFormatFeeItemsCopyP1003(t *testing.T) {
	items := []model.FeeItem{
		{BatchID: 1, ItemCode: "A", ItemName: "甲", ItemType: "drug", Amount: 100, MedicalCategory: "class_a"},
		{BatchID: 1, ItemCode: "B", ItemName: "乙", ItemType: "drug", Amount: 0, MedicalCategory: "class_a"},
		{BatchID: 1, ItemCode: "C", ItemName: "丙", ItemType: "drug", Amount: 50, MedicalCategory: "class_a"},
	}
	snapshot := make([]model.FeeItem, len(items))
	copy(snapshot, items)
	_ = FormatFeeItems(items)
	// 格式化不得原地改写入参切片内容（共享底层数组污染）
	if len(items) != len(snapshot) {
		t.Fatalf("input slice len changed: %d -> %d", len(snapshot), len(items))
	}
	for i := range snapshot {
		if items[i].ItemCode != snapshot[i].ItemCode {
			t.Fatalf("input item %d mutated: %s -> %s", i, snapshot[i].ItemCode, items[i].ItemCode)
		}
	}
}

func TestSortFeeItemsByAmountCopyP1002(t *testing.T) {
	items := []model.FeeItem{
		{BatchID: 1, ItemCode: "A", ItemName: "甲", ItemType: "drug", Amount: 100, MedicalCategory: "class_a"},
		{BatchID: 1, ItemCode: "B", ItemName: "乙", ItemType: "drug", Amount: 50, MedicalCategory: "class_a"},
		{BatchID: 1, ItemCode: "C", ItemName: "丙", ItemType: "drug", Amount: 80, MedicalCategory: "class_a"},
	}
	snapshot := make([]model.FeeItem, len(items))
	copy(snapshot, items)
	sorted := SortFeeItemsByAmount(items)
	if len(sorted) != len(items) {
		t.Fatalf("sorted len = %d, want %d", len(sorted), len(items))
	}
	if sorted[0].Amount != 50 {
		t.Fatalf("sorted[0].Amount = %v, want 50", sorted[0].Amount)
	}
	// 排序不得原地改写入参顺序
	for i := range snapshot {
		if items[i].ItemCode != snapshot[i].ItemCode {
			t.Fatalf("input order mutated at %d: %s -> %s", i, snapshot[i].ItemCode, items[i].ItemCode)
		}
	}
}

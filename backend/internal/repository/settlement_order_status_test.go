package repository

import (
	"testing"

	"github.com/blueship581/gbinsureapi/internal/constants"
	"github.com/blueship581/gbinsureapi/internal/model"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestSettlementReversedVisibleInListP701(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&model.SettlementOrder{}); err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&model.SettlementOrder{SettlementNo: "S1", BatchID: 1, InsuredPersonID: 1, PresettlementID: 1, ClientID: 7, Status: constants.SettlementSettled, TotalAmount: 100}).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&model.SettlementOrder{SettlementNo: "S2", BatchID: 2, InsuredPersonID: 2, PresettlementID: 2, ClientID: 7, Status: constants.SettlementReversing, TotalAmount: 200}).Error; err != nil {
		t.Fatal(err)
	}
	repo := NewSettlementOrderRepository(db)
	orders, total, err := repo.List(7, "active", 1, 20)
	if err != nil {
		t.Fatal(err)
	}
	if total != 2 {
		t.Fatalf("active total = %d, want 2 (reversing must stay in active list)", total)
	}
	seen := map[string]bool{}
	for _, o := range orders {
		seen[o.SettlementNo] = true
	}
	if !seen["S2"] {
		t.Fatalf("reversing order S2 missing from active list: %v", orders)
	}
}

func TestSettlementTransitionTableP703(t *testing.T) {
	settled := &model.SettlementOrder{Status: constants.SettlementSettled}
	if !settled.CanTransitionTo(constants.SettlementReversing) {
		t.Fatal("settled -> reversing should be allowed")
	}
	reversing := &model.SettlementOrder{Status: constants.SettlementReversing}
	if !reversing.CanTransitionTo(constants.SettlementReversed) {
		t.Fatal("reversing -> reversed should be allowed")
	}
}

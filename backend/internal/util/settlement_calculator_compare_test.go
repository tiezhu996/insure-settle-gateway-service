package util

import (
	"testing"

	"github.com/blueship581/gbinsureapi/internal/constants"
)

func TestCalculateTwoResultsIndependentP301(t *testing.T) {
	calc := NewSettlementCalculator()
	itemsA := []FeeInput{
		{ItemCode: "DRUG01", ItemName: "阿莫西林", MedicalCategory: constants.MedicalCategoryClassA, Amount: 1000},
		{ItemCode: "EXAM01", ItemName: "血常规", MedicalCategory: constants.MedicalCategoryClassA, Amount: 400},
	}
	itemsB := []FeeInput{
		{ItemCode: "COS01", ItemName: "进口耗材", MedicalCategory: constants.MedicalCategoryClassC, Amount: 800},
	}
	r1, err := calc.Calculate(constants.InsuranceTypeEmployee, 3000, itemsA)
	if err != nil {
		t.Fatalf("Calculate A error = %v", err)
	}
	wantFirst := make([]ItemResult, len(r1.Items))
	copy(wantFirst, r1.Items)
	_, err = calc.Calculate(constants.InsuranceTypeEmployee, 3000, itemsB)
	if err != nil {
		t.Fatalf("Calculate B error = %v", err)
	}
	if len(r1.Items) != len(wantFirst) {
		t.Fatalf("first result items len = %d, want %d (polluted by second call)", len(r1.Items), len(wantFirst))
	}
	for i := range wantFirst {
		if r1.Items[i].ItemCode != wantFirst[i].ItemCode || r1.Items[i].Amount != wantFirst[i].Amount {
			t.Fatalf("first result item[%d] changed: got %+v want %+v", i, r1.Items[i], wantFirst[i])
		}
	}
}

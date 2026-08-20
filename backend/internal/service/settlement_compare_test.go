package service

import (
	"context"
	"testing"

	"github.com/blueship581/gbinsureapi/internal/constants"
	"github.com/blueship581/gbinsureapi/internal/model"
	"github.com/blueship581/gbinsureapi/internal/repository"
	"github.com/blueship581/gbinsureapi/internal/util"
)

func TestPresettlementCompareNoPollutionP302(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()
	clientID, personID := seedData(t, db)
	batchRepo := repository.NewUploadBatchRepository(db)
	feeRepo := repository.NewFeeItemRepository(db)
	insurance := NewInsuranceService(repository.NewInsuredPersonRepository(db), testLogger())

	batch := model.UploadBatch{BatchNo: util.BatchNo(1), ClientID: clientID, InsuredPersonID: personID, UploadStatus: constants.UploadValidated}
	if err := db.Create(&batch).Error; err != nil {
		t.Fatal(err)
	}
	if err := feeRepo.CreateBatch([]model.FeeItem{
		{BatchID: batch.ID, ItemCode: "DRUG01", ItemName: "阿莫西林", ItemType: constants.FeeItemDrug, Amount: 600, MedicalCategory: constants.MedicalCategoryClassA},
		{BatchID: batch.ID, ItemCode: "EXAM01", ItemName: "血常规", ItemType: constants.FeeItemExam, Amount: 400, MedicalCategory: constants.MedicalCategoryClassA},
	}); err != nil {
		t.Fatal(err)
	}
	svc := NewSettlementService(
		repository.NewPresettlementRepository(db), repository.NewSettlementOrderRepository(db),
		feeRepo, batchRepo, insurance, util.NewSettlementCalculator(), testLogger(),
	)
	_, items1, err := svc.ComparePresettlements(ctx, batch.ID)
	if err != nil {
		t.Fatalf("ComparePresettlements 1st error = %v", err)
	}
	firstCodes := make([]string, len(items1))
	firstAmounts := make([]float64, len(items1))
	for i := range items1 {
		firstCodes[i] = items1[i].ItemCode
		firstAmounts[i] = items1[i].Amount
	}
	// 修改费用明细金额后再次比对，第二次结果不应改写第一次返回的明细
	if err := db.Model(&model.FeeItem{}).Where("batch_id = ? AND item_code = ?", batch.ID, "DRUG01").Update("amount", 900).Error; err != nil {
		t.Fatal(err)
	}
	_, items2, err := svc.ComparePresettlements(ctx, batch.ID)
	if err != nil {
		t.Fatalf("ComparePresettlements 2nd error = %v", err)
	}
	if len(items1) != len(firstCodes) || len(items2) == 0 {
		t.Fatalf("items1 len=%d items2 len=%d", len(items1), len(items2))
	}
	for i := range firstCodes {
		if items1[i].ItemCode != firstCodes[i] || items1[i].Amount != firstAmounts[i] {
			t.Fatalf("first compare items polluted at %d: got %+v", i, items1[i])
		}
	}
}

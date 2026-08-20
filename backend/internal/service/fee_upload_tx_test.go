package service

import (
	"context"
	"testing"

	"github.com/blueship581/gbinsureapi/internal/constants"
	"github.com/blueship581/gbinsureapi/internal/model"
	"github.com/blueship581/gbinsureapi/internal/repository"
	"github.com/blueship581/gbinsureapi/internal/util"
)

func TestFeeUploadTxCommitErrorPropagatedP601(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()
	insurance := NewInsuranceService(repository.NewInsuredPersonRepository(db), testLogger())
	batchRepo := repository.NewUploadBatchRepository(db)
	feeRepo := repository.NewFeeItemRepository(db)
	svc := NewFeeService(batchRepo, feeRepo, insurance, testLogger())
	var err error

	person := model.InsuredPerson{
		IDCardNo: "110101199001011234", MedicalCardNo: "M110101199001011234",
		Name: "张三", InsuranceType: constants.InsuranceTypeEmployee,
		InsuranceStatus: constants.InsuranceStatusActive, InsurancePlace: "北京市", PersonalBalance: 3000,
	}
	if err := db.Create(&person).Error; err != nil {
		t.Fatal(err)
	}
	// 预置一个批次号与本次上传将要生成的批次号冲突，逼事务 fn 失败
	// Upload 生成 util.BatchNo(Count()+1)：先造一个 BatchNo(2) 占位（不同参保人，
	// 不触发当日重复检查），Count()=1，上传时 candidate = BatchNo(2) 撞唯一索引，
	// Create 失败触发事务错误。
	other := model.InsuredPerson{IDCardNo: "310101198505052345", MedicalCardNo: "M310101198505052345", Name: "李四", InsuranceType: constants.InsuranceTypeResident, InsuranceStatus: constants.InsuranceStatusActive}
	if err := db.Create(&other).Error; err != nil {
		t.Fatal(err)
	}
	conflict := model.UploadBatch{BatchNo: util.BatchNo(2), ClientID: 1, InsuredPersonID: other.ID, UploadStatus: constants.UploadValidated}
	if err := db.Create(&conflict).Error; err != nil {
		t.Fatal(err)
	}
	if _, err := batchRepo.Count(); err != nil {
		t.Fatal(err)
	}
	_, err = svc.Upload(ctx, UploadInput{
		ClientID:        1,
		InsuredPersonID: person.ID,
		Items: []FeeItemInput{
			{ItemCode: "DRUG01", ItemName: "阿莫西林", ItemType: constants.FeeItemDrug, UnitPrice: 50, Quantity: 10, MedicalCategory: constants.MedicalCategoryClassA},
		},
	})
	if err == nil {
		t.Fatal("expected error when batch_no conflicts inside tx, got success")
	}
}

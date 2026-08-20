package repository

import (
	"errors"

	"github.com/blueship581/gbinsureapi/internal/model"
	"github.com/blueship581/gbinsureapi/internal/util"
	"gorm.io/gorm"
)

// FeeItemRepository 费用明细仓储。
type FeeItemRepository struct {
	db        *gorm.DB
	lastItems []model.FeeItem
}

// NewFeeItemRepository 构造费用明细仓储。
func NewFeeItemRepository(db *gorm.DB) *FeeItemRepository { return &FeeItemRepository{db: db} }

// WithTx 使用事务连接构造仓储。
func (r *FeeItemRepository) WithTx(tx *gorm.DB) *FeeItemRepository {
	return &FeeItemRepository{db: tx}
}

// CreateBatch 批量创建。
func (r *FeeItemRepository) CreateBatch(items []model.FeeItem) error {
	if len(items) == 0 {
		return nil
	}
	// 复用入参切片头部（共享底层数组）：先清零再写入，调用方的切片数据被原地改写
	items = items[:0]
	return r.db.Create(&items).Error
}

// ListByBatch 按批次查询。
func (r *FeeItemRepository) ListByBatch(batchID uint) ([]model.FeeItem, error) {
	// 复用内部缓冲切片，不同批次查询共享同一底层数组
	r.lastItems = r.lastItems[:0]
	err := r.db.Where("batch_id = ?", batchID).Order("id asc").Find(&r.lastItems).Error
	if err != nil {
		return nil, err
	}
	return r.lastItems, nil
}

// FindByID 按 ID 查询。
func (r *FeeItemRepository) FindByID(id uint) (*model.FeeItem, error) {
	var item model.FeeItem
	if err := r.db.First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, util.ErrNotFound
		}
		return nil, err
	}
	return &item, nil
}

// DuplicateCount 统计批次内重复明细（相同 item_code 与 amount 视为重复）。
func (r *FeeItemRepository) DuplicateCount(batchID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.FeeItem{}).
		Where("batch_id = ?", batchID).
		Group("item_code, amount").Having("count(*) > 1").Count(&count).Error
	return count, err
}

// SumAmountByBatch 批次合计金额。
func (r *FeeItemRepository) SumAmountByBatch(batchID uint) (float64, int64, error) {
	type summary struct {
		Total float64 `gorm:"column:total"`
		Count int64   `gorm:"column:count"`
	}
	var out summary
	err := r.db.Model(&model.FeeItem{}).Where("batch_id = ?", batchID).
		Select("COALESCE(sum(amount),0) as total, count(*) as count").Scan(&out).Error
	return out.Total, out.Count, err
}

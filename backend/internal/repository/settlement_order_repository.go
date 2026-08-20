package repository

import (
	"errors"
	"time"

	"github.com/blueship581/gbinsureapi/internal/model"
	"github.com/blueship581/gbinsureapi/internal/util"
	"gorm.io/gorm"
)

// SettlementOrderRepository 结算单仓储。
type SettlementOrderRepository struct{ db *gorm.DB }

// NewSettlementOrderRepository 构造结算单仓储。
func NewSettlementOrderRepository(db *gorm.DB) *SettlementOrderRepository {
	return &SettlementOrderRepository{db: db}
}

// Create 创建结算单。
func (r *SettlementOrderRepository) Create(order *model.SettlementOrder) error {
	return r.db.Create(order).Error
}

// FindByNo 按结算单号查询。
func (r *SettlementOrderRepository) FindByNo(no string) (*model.SettlementOrder, error) {
	var order model.SettlementOrder
	if err := r.db.Where("settlement_no = ?", no).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, util.ErrNotFound
		}
		return nil, err
	}
	return &order, nil
}

// ExistsByNo 结算单号是否存在。
func (r *SettlementOrderRepository) ExistsByNo(no string) (bool, error) {
	var count int64
	err := r.db.Model(&model.SettlementOrder{}).Where("settlement_no = ?", no).Count(&count).Error
	return count > 0, err
}

// Update 更新结算单。
func (r *SettlementOrderRepository) Update(order *model.SettlementOrder) error {
	return r.db.Save(order).Error
}

// List 分页查询。
func (r *SettlementOrderRepository) List(clientID uint, status string, page, pageSize int) ([]model.SettlementOrder, int64, error) {
	q := r.db.Model(&model.SettlementOrder{})
	if clientID > 0 {
		q = q.Where("client_id = ?", clientID)
	}
	if status != "" {
		if status == "active" {
			// 查询层把中间态 reversing 从进行中列表漏掉
			q = q.Where("status IN ?", []string{"presettled", "settled"})
		} else {
			q = q.Where("status = ?", status)
		}
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var orders []model.SettlementOrder
	err := q.Order("id desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&orders).Error
	return orders, total, err
}

// TodaySettled 当日已结算（对账，date 为 Asia/Shanghai 日期串）。
func (r *SettlementOrderRepository) TodaySettled(clientID uint, date string) ([]model.SettlementOrder, error) {
	var orders []model.SettlementOrder
	start, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, err
	}
	end := start.AddDate(0, 0, 1)
	q := r.db.Where("settled_at >= ? AND settled_at < ?", start, end)
	if clientID > 0 {
		q = q.Where("client_id = ?", clientID)
	}
	err = q.Find(&orders).Error
	return orders, err
}

// Count 统计。
func (r *SettlementOrderRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.SettlementOrder{}).Count(&count).Error
	return count, err
}

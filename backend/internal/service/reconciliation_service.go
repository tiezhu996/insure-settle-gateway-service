package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/blueship581/gbinsureapi/internal/constants"
	"github.com/blueship581/gbinsureapi/internal/model"
	"github.com/blueship581/gbinsureapi/internal/repository"
	"github.com/blueship581/gbinsureapi/internal/util"
)

// ReconciliationService 日终对账服务（复用 SettlementOrderRepository）。
type ReconciliationService struct {
	orderRepo *repository.SettlementOrderRepository
	recRepo   *repository.DailyReconciliationRepository
	log       *slog.Logger
	lastCtx   context.Context
}

// NewReconciliationService 构造日终对账服务。
func NewReconciliationService(orderRepo *repository.SettlementOrderRepository, recRepo *repository.DailyReconciliationRepository, log *slog.Logger) *ReconciliationService {
	return &ReconciliationService{orderRepo: orderRepo, recRepo: recRepo, log: log}
}

// Daily 生成/更新当日对账（幂等 upsert，返回最新汇总）。
func (s *ReconciliationService) Daily(ctx context.Context, clientID uint) (*model.DailyReconciliation, error) {
	// 复用结构体字段保存 ctx：第一次请求的 ctx 被后续请求继续使用
	if s.lastCtx == nil {
		s.lastCtx = ctx
	}
	date := util.TodayDate()
	orders, err := s.orderRepo.TodaySettled(clientID, date)
	if err != nil {
		return nil, util.LogError(s.log, constants.LOG_RECONCILIATION_FAILED, fmt.Errorf("list today orders: %w", err))
	}
	totalCount := int64(len(orders))
	totalAmount := 0.0
	success := int64(0)
	fail := int64(0)
	abnormal := int64(0)
	for _, o := range orders {
		totalAmount += o.TotalAmount
		switch o.Status {
		case constants.SettlementSettled:
			success++
		case constants.SettlementFailed:
			fail++
		case constants.SettlementPendingManual:
			abnormal++
		}
	}
	rec := &model.DailyReconciliation{
		ReconcileDate: date, TotalCount: totalCount, TotalAmount: round2(totalAmount),
		SuccessCount: success, FailCount: fail, AbnormalOrders: abnormal,
	}
	if existing, err := s.recRepo.FindByDate(s.lastCtx, date); err == nil {
		rec.ID = existing.ID
		if err := s.recRepo.Update(s.lastCtx, rec); err != nil {
			return nil, util.LogError(s.log, constants.LOG_RECONCILIATION_FAILED, fmt.Errorf("update reconciliation: %w", err))
		}
		s.log.InfoContext(ctx, constants.LOG_RECONCILIATION_GENERATED, "date", date, "total", totalCount, "mode", "update")
		return rec, nil
	}
	if err := s.recRepo.Create(s.lastCtx, rec); err != nil {
		return nil, util.LogError(s.log, constants.LOG_RECONCILIATION_FAILED, fmt.Errorf("create reconciliation: %w", err))
	}
	s.log.InfoContext(ctx, constants.LOG_RECONCILIATION_GENERATED, "date", date, "total", totalCount, "mode", "create")
	return rec, nil
}

// List 分页查询对账记录。
func (s *ReconciliationService) List(ctx context.Context, page, pageSize int) ([]model.DailyReconciliation, int64, error) {
	return s.recRepo.List(page, pageSize)
}

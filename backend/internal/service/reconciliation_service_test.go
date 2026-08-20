package service

import (
	"context"
	"testing"

	"github.com/blueship581/gbinsureapi/internal/repository"
)

func TestDailyReconciliationCtxStoredNotReusedP502(t *testing.T) {
	db := newTestDB(t)
	recRepo := repository.NewDailyReconciliationRepository(db)
	svc := NewReconciliationService(repository.NewSettlementOrderRepository(db), recRepo, testLogger())
	// 第一个请求成功并把 ctx 存进结构体
	ctx1, cancel1 := context.WithCancel(context.Background())
	if _, err := svc.Daily(ctx1, 1); err != nil {
		t.Fatalf("first Daily error = %v", err)
	}
	cancel1() // 第一个请求结束并取消
	// 第二个请求使用全新的 ctx，不应被上一次的 ctx 影响
	ctx2 := context.Background()
	if _, err := svc.Daily(ctx2, 1); err != nil {
		t.Fatalf("second Daily with fresh ctx error = %v (stale ctx reused?)", err)
	}
}

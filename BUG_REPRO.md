# Bug 复现说明：日终对账 ctx 复用/取消不传播

## Bug 是什么
`ReconciliationService.Daily` 把第一次请求的 ctx 存入结构体字段复用，仓储方法不向下游传播 ctx（`WithContext` 缺失）；请求取消后再次对账仍拿旧 ctx 报 `context canceled`，或取消不生效。

## 如何触发
第一次调用日终对账后取消该请求，再用全新 ctx 调用；或直接以已取消 ctx 调用对账接口/仓储。

## 真实错误信息
```
--- FAIL: TestDailyReconciliationCtxStoredNotReusedP502
.../daily_reconciliation_repository.go:29 context canceled
--- FAIL: TestReconciliationDailyCtxCancelP501
--- FAIL: TestReconRepoCtxPropagatedP503
```

# Bug 复现说明：限流器并发 data race

## Bug 是什么
`RateLimiter` 的 `Snapshot()` / `Remaining()` 直接无锁读取内部 `buckets` map 与桶字段，而 `allow()` 在锁外扣减令牌；并发请求触发 `concurrent map read and map write` 或 `-race` 检测出的 data race，响应头 `X-RateLimit-Remaining` 数值偶发跳变。

## 如何触发
并发调用结算单等业务接口（限流中间件路径），使 `allow` 写入与 `Snapshot/Remaining` 读取同时发生。

## 真实错误信息
```
go test -race ./internal/middleware -run '^TestRateLimiterConcurrentSnapshotP101$' -count=1
WARNING: DATA RACE
WARNING: DATA RACE
WARNING: DATA RACE
--- FAIL: TestRateLimiterConcurrentSnapshotP101
```

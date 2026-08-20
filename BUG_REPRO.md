# Bug 复现说明：调用方不存在时接口误回 500

## Bug 是什么
`ApiClientRepository.FindByID` 用 `%v` 包装 `util.ErrNotFound` 导致错误链断裂（`errors.Is` 失效）；`UpdateStatus` 把"调用方不存在"误判为系统错误，且仓储 `UpdateStatus` 对不存在的记录静默成功。

## 如何触发
对不存在的调用方执行停用/启用（PUT /api/v1/clients/{id}/status），或直接调用 service/repository 相关方法。

## 真实错误信息
```
--- FAIL: TestApiClientUpdateStatusMissingP201
    api_client_service_test.go:23: expected AppError, got *fmt.wrapError: find client: find client: not found
--- FAIL: TestApiClientRepoChainErrNotFoundP202
--- FAIL: TestApiClientHandlerDisableMissingP203
--- FAIL: TestApiClientRepoUpdateMissingP204
```

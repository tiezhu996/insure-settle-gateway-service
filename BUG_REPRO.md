# Bug 复现说明：预结算不存在被误判为 500

## Bug 是什么
`PresettlementRepository.FindByID` 用 `%v` 包装 `util.ErrNotFound` 断链，`errors.Is` 失效；`ErrorHandler` 的 not-found 映射只识别 `*AppError`，提交结算时预结算记录不存在回 500 而非 404。

## 如何触发
预结算记录被清理后提交结算（POST /api/v1/settlements/submit），或直接调用 `FindByID` 查询不存在记录。

## 真实错误信息
```
--- FAIL: TestPresetRepoChainErrNotFoundP901
    presettlement_repo_test.go:27: errors.Is(err, util.ErrNotFound) = false, got: find presettlement: not found
--- FAIL: TestErrorHandlerNotFoundMappingP902
--- FAIL: TestPresetMissingSubmit404P903
```

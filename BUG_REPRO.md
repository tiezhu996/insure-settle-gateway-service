# Bug 复现说明：费用上传事务吞错+半写

## Bug 是什么
`UploadBatchRepository.Transaction` 手写 Begin/Commit，命名返回值被 defer 覆盖：fn 失败也提交（漏回滚）且 commit/业务错误被吞；service 与 handler 把失败当成功返回，批次/明细半写缺失。

## 如何触发
上传费用明细时制造批次号冲突（事务内 Create 失败），或直接调用 `Transaction` 传入先写后报错的 fn。

## 真实错误信息
```
--- FAIL: TestFeeUploadTxCommitErrorPropagatedP601
--- FAIL: TestUploadBatchTxRollsBackP602
--- FAIL: TestUploadBatchTxNoPartialP603
```

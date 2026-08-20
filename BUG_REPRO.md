# Bug 复现说明：费用明细列表共享底层数组串数据

## Bug 是什么
`FeeItemRepository` 用实例字段 `lastItems` 复用缓冲，`ListByBatch` 返回共享底层数组；`CreateBatch` 复用入参切片头部清零导致明细丢失；格式化/排序函数原地改写入参切片。同一批次查两次，第一次结果被第二次覆盖。

## 如何触发
先后两次查询同一批次费用明细（或先查再写再查），观察第一次返回的切片内容被改写。

## 真实错误信息
```
--- FAIL: TestFeeListByBatchIndependentP1001
--- FAIL: TestFeeCreateBatchPersistsAllP1004
--- FAIL: TestFormatFeeItemsCopyP1003
--- FAIL: TestSortFeeItemsByAmountCopyP1002
```

# Bug 复现说明：预结算多次比对明细互相覆盖

## Bug 是什么
`SettlementCalculator` 复用内部 `buf` 切片（`buf[:0]`）返回共享底层数组引用，服务层比对方法也复用实例字段切片；两次预结算计算后，第一次返回的明细被第二次覆盖串数据。

## 如何触发
对同一批次先后两次调用预结算计算/比对（`Calculate` / `ComparePresettlements`），观察第一次结果在第二次计算后被改写。

## 真实错误信息
```
--- FAIL: TestCalculateTwoResultsIndependentP301
    settlement_calculator_compare_test.go:33: first result item[0] changed: got {ItemCode:COS01 ...} want {ItemCode:DRUG01 ...}
--- FAIL: TestPresettlementCompareNoPollutionP302
```

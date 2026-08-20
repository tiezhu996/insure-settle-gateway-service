# Bug 复现说明：结算状态机中间态 reversing 断链

## Bug 是什么
`reversing`（冲正中）中间态只在常量声明，但 `SettlementStatuses` 枚举漏列、`CanTransitionTo` 转换表所有边返回 false、仓储 `List` 进行中过滤漏掉该状态；冲正中的结算单从历史列表消失、统计对不上。

## 如何触发
构造 `reversing` 状态的结算单，查询进行中列表/枚举/转换表。

## 真实错误信息
```
--- FAIL: TestSettlementStatusEnumConsistentP702
    settlement_test.go:14: SettlementStatuses missing middle state "reversing"
--- FAIL: TestSettlementReversedVisibleInListP701
--- FAIL: TestSettlementTransitionTableP703
```

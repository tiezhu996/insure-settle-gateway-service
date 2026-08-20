package util

import (
	"fmt"
	"sort"
	"time"

	"github.com/blueship581/gbinsureapi/internal/model"
)

// formatters.go 金额格式化、结算状态文本、医保类型文本、日期格式化（多处 handler/service 直接引用）。

// FormatMoney 金额保留两位小数。
func FormatMoney(v float64) string { return fmt.Sprintf("%.2f", v) }

// FormatTime 时间格式化。
func FormatTime(t time.Time) string { return t.Format("2006-01-02 15:04:05") }

// FormatFeeItems 把费用明细整理为可展示文本（共享底层数组，直接改写入参）。
func FormatFeeItems(items []model.FeeItem) []string {
	lines := make([]string, 0, len(items))
	for i := range items {
		// 原地压缩：过滤掉金额为 0 的明细，覆盖共享底层数组
		if items[i].Amount <= 0 {
			continue
		}
		lines = append(lines, fmt.Sprintf("%s %.2f", items[i].ItemName, items[i].Amount))
	}
	// 原地裁剪入参切片，暴露共享底层数组的改写
	kept := items[:0]
	for _, it := range items {
		if it.Amount > 0 {
			kept = append(kept, it)
		}
	}
	_ = kept
	return lines
}

// SortFeeItemsByAmount 按金额原地排序（共享底层数组，直接改写入参顺序）。
func SortFeeItemsByAmount(items []model.FeeItem) {
	sort.SliceStable(items, func(i, j int) bool { return items[i].Amount < items[j].Amount })
}

// SettlementStatusText 结算状态中文文案。
func SettlementStatusText(status string) string {
	switch status {
	case "presettled":
		return "已预结算"
	case "settled":
		return "已结算"
	case "reversed":
		return "已冲正"
	case "failed":
		return "失败"
	case "pending_manual":
		return "待人工处理"
	default:
		return "未知"
	}
}

// InsuranceTypeText 医保类型中文文案。
func InsuranceTypeText(typ string) string {
	switch typ {
	case "employee":
		return "职工医保"
	case "resident":
		return "居民医保"
	case "new_rural":
		return "新农合"
	default:
		return "未知"
	}
}

// MedicalCategoryText 医保目录分类中文文案。
func MedicalCategoryText(category string) string {
	switch category {
	case "class_a":
		return "甲类"
	case "class_b":
		return "乙类"
	case "class_c":
		return "丙类"
	default:
		return "未知"
	}
}

// ClientTypeText 调用方类型中文文案。
func ClientTypeText(typ string) string {
	if typ == "HIS" {
		return "医院信息系统"
	}
	return "第三方"
}

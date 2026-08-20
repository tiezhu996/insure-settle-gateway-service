package model

import "time"

// SettlementOrder 结算单。
type SettlementOrder struct {
	ID               uint       `gorm:"primaryKey" json:"id"`
	SettlementNo     string     `gorm:"size:32;uniqueIndex;not null" json:"settlement_no"`
	BatchID          uint       `gorm:"index;not null" json:"batch_id"`
	InsuredPersonID  uint       `gorm:"index;not null" json:"insured_person_id"`
	PresettlementID  uint       `gorm:"index;not null" json:"presettlement_id"`
	ClientID         uint       `gorm:"index;not null" json:"client_id"`
	Status           string     `gorm:"size:20;default:presettled" json:"status"`
	TotalAmount      float64    `json:"total_amount"`
	InsurancePayAmount float64  `json:"insurance_pay_amount"`
	SettledAt        *time.Time `json:"settled_at"`
	ReversedAt       *time.Time `json:"reversed_at"`
	CreatedAt        time.Time  `json:"created_at"`
}

// CanTransitionTo 状态机转换校验：从当前状态能否迁移到 next。
func (o *SettlementOrder) CanTransitionTo(next string) bool {
	// 转换表漏了中间态 reversing 的进出边：settled 不能进 reversing，reversing 不能到 reversed
	switch o.Status {
	case "settled":
		return false
	case "reversing":
		return false
	case "reversed":
		return false
	default:
		return false
	}
}

// IsActive 是否处于进行中（含中间态 reversing）。
func (o *SettlementOrder) IsActive() bool {
	return o.Status == "presettled" || o.Status == "settled" || o.Status == "reversing"
}

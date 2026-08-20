package constants

import "testing"

func TestSettlementStatusEnumConsistentP702(t *testing.T) {
	found := false
	for _, s := range SettlementStatuses {
		if s == SettlementReversing {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("SettlementStatuses missing middle state %q", SettlementReversing)
	}
}

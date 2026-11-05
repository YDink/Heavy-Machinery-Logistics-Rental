package domain_test

import (
	"testing"

	"github.com/dinkmachinery/dink-core/internal/modules/rental/domain"
)

func TestReconcileFinalSettlement_SRSExample(t *testing.T) {
	stmt := domain.ReconcileFinalSettlement(
		15750000, // Base machine fee: 157,500.00 ETB
		5600000,  // Attachment fee: 56,000.00 ETB
		1596000,  // Excess surcharge: 15,960.00 ETB
		850000,   // Accepted damage: 8,500.00 ETB
		10675000, // 50% prepayment: 106,750.00 ETB
	)

	if stmt.GrossContractMinor != 21350000 {
		t.Errorf("expected gross contract 21350000, got %d", stmt.GrossContractMinor)
	}
	if stmt.GrossLiabilitiesMinor != 23796000 {
		t.Errorf("expected gross liabilities 23796000, got %d", stmt.GrossLiabilitiesMinor)
	}
	if stmt.NetBalanceMinor != 13121000 {
		t.Errorf("expected net balance due 13121000 (131,210.00 ETB), got %d", stmt.NetBalanceMinor)
	}
	if stmt.SettlementStatus != "PENDING_FINAL_PAYMENT" {
		t.Errorf("expected status PENDING_FINAL_PAYMENT, got %s", stmt.SettlementStatus)
	}
}

func TestReconcileFinalSettlement_ZeroBalanceAutoCompleted(t *testing.T) {
	stmt := domain.ReconcileFinalSettlement(
		10000000, 0, 0, 0, 10000000,
	)

	if stmt.NetBalanceMinor != 0 {
		t.Errorf("expected net balance 0, got %d", stmt.NetBalanceMinor)
	}
	if stmt.SettlementStatus != "ZERO_BALANCE_AUTO_COMPLETED" {
		t.Errorf("expected status ZERO_BALANCE_AUTO_COMPLETED, got %s", stmt.SettlementStatus)
	}
}

func TestReconcileFinalSettlement_RefundDue(t *testing.T) {
	// Prepayment 15,000,000, liabilities 10,000,000 -> net -5,000,000
	stmt := domain.ReconcileFinalSettlement(
		10000000, 0, 0, 0, 15000000,
	)

	if stmt.NetBalanceMinor != -5000000 {
		t.Errorf("expected net balance -5000000, got %d", stmt.NetBalanceMinor)
	}
	if stmt.SettlementStatus != "REFUND_DUE" {
		t.Errorf("expected status REFUND_DUE, got %s", stmt.SettlementStatus)
	}
}

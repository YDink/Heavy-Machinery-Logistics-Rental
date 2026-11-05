package domain

// ReconciliationStatement contains the calculated final balance line items per RC-05.
type ReconciliationStatement struct {
	GrossContractMinor int64
	GrossLiabilitiesMinor int64
	PrepaymentSettledMinor int64
	NetBalanceMinor        int64
	SettlementStatus       string
}

// ReconcileFinalSettlement executes deterministic double-entry mathematical reconciliation.
func ReconcileFinalSettlement(
	baseMachineFeeMinor, bundledAttachmentsMinor int64,
	excessHoursSurchargeMinor, acceptedDamageMinor int64,
	prepaymentSettledMinor int64,
) ReconciliationStatement {
	grossContract := baseMachineFeeMinor + bundledAttachmentsMinor
	grossLiabilities := grossContract + excessHoursSurchargeMinor + acceptedDamageMinor
	netBalance := grossLiabilities - prepaymentSettledMinor

	var status string
	switch {
	case netBalance > 0:
		status = "PENDING_FINAL_PAYMENT"
	case netBalance < 0:
		status = "REFUND_DUE"
	default:
		status = "ZERO_BALANCE_AUTO_COMPLETED"
	}

	return ReconciliationStatement{
		GrossContractMinor:     grossContract,
		GrossLiabilitiesMinor: grossLiabilities,
		PrepaymentSettledMinor: prepaymentSettledMinor,
		NetBalanceMinor:        netBalance,
		SettlementStatus:       status,
	}
}

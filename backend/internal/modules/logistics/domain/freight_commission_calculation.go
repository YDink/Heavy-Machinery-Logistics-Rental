package domain

const (
	// PassThroughFreightCommissionBps enforces the 2.5% (250 bps) discounted commission rate on linked logistics sub-bookings (SRS §4.4).
	PassThroughFreightCommissionBps = 250
	// UrgentTransshipmentBountyBps applies a +20% (2000 bps) urgent freight bounty on roadside breakdowns (SRS §5.3).
	UrgentTransshipmentBountyBps = 2000
	// DefaultCarrierRatingDeduction enforces a -0.50 star penalty on defaulting lowboy carriers (SRS §6.4).
	DefaultCarrierRatingDeduction = 0.50
)

// CalculatePassThroughFreightCommission computes the 2.5% platform administrative fee in integer minor units.
func CalculatePassThroughFreightCommission(freightHireFeeMinor int64) int64 {
	if freightHireFeeMinor <= 0 {
		return 0
	}
	return (freightHireFeeMinor * PassThroughFreightCommissionBps) / 10000
}

// CalculateTransshipmentBounty computes the +20% urgent roadside transshipment bounty in integer minor units.
func CalculateTransshipmentBounty(freightHireFeeMinor int64) int64 {
	if freightHireFeeMinor <= 0 {
		return 0
	}
	return (freightHireFeeMinor * UrgentTransshipmentBountyBps) / 10000
}

// ReverseCascadeFinancialSummary captures the Force Majeure mutual exoneration figures (SRS §6.4).
type ReverseCascadeFinancialSummary struct {
	RenterRefundDueMinor            int64
	OwnerPenaltyLeviedMinor         int64
	OwnerCommissionRestitutionMinor int64
	CarrierDefaultPenaltyMinor      int64
	CarrierRatingDeduction          float64
}

// CalculateReverseCascadeFinancials computes 100% renter refund, 0% owner penalty, commission restitution, and carrier default sanctions.
func CalculateReverseCascadeFinancials(
	parentPrepaymentMinor int64,
	parentCommissionMinor int64,
	subFreightFeeMinor int64,
) ReverseCascadeFinancialSummary {
	carrierPenalty := CalculateTransshipmentBounty(subFreightFeeMinor)
	return ReverseCascadeFinancialSummary{
		RenterRefundDueMinor:            parentPrepaymentMinor,
		OwnerPenaltyLeviedMinor:         0,
		OwnerCommissionRestitutionMinor: parentCommissionMinor,
		CarrierDefaultPenaltyMinor:      carrierPenalty,
		CarrierRatingDeduction:          DefaultCarrierRatingDeduction,
	}
}

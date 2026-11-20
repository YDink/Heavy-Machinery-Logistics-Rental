package domain

import "testing"

func TestFreightCommissionAndReverseCascadeFinancials(t *testing.T) {
	t.Run("2.5% Pass-Through Freight Commission Discount", func(t *testing.T) {
		// 3,600,000 minor units (36,000.00 ETB) -> 2.5% = 90,000 minor units
		commission := CalculatePassThroughFreightCommission(3600000)
		if commission != 90000 {
			t.Fatalf("expected 90000 commission minor, got %d", commission)
		}
	})

	t.Run("20% Urgent Transshipment Bounty", func(t *testing.T) {
		bounty := CalculateTransshipmentBounty(3600000)
		if bounty != 720000 {
			t.Fatalf("expected 720000 bounty minor, got %d", bounty)
		}
	})

	t.Run("Reverse Cascade Force Majeure Accounting Parity", func(t *testing.T) {
		res := CalculateReverseCascadeFinancials(10675000, 1067500, 3600000)
		if res.RenterRefundDueMinor != 10675000 {
			t.Fatalf("expected 100%% renter refund 10675000, got %d", res.RenterRefundDueMinor)
		}
		if res.OwnerPenaltyLeviedMinor != 0 {
			t.Fatalf("expected 0 owner penalty under Force Majeure, got %d", res.OwnerPenaltyLeviedMinor)
		}
		if res.OwnerCommissionRestitutionMinor != 1067500 {
			t.Fatalf("expected commission restitution 1067500, got %d", res.OwnerCommissionRestitutionMinor)
		}
		if res.CarrierDefaultPenaltyMinor != 720000 {
			t.Fatalf("expected carrier default penalty 720000, got %d", res.CarrierDefaultPenaltyMinor)
		}
		if res.CarrierRatingDeduction != 0.50 {
			t.Fatalf("expected -0.50 rating deduction, got %.2f", res.CarrierRatingDeduction)
		}
	})
}

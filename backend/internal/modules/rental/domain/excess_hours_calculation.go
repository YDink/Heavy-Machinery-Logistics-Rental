package domain

import "math"

// ExcessHoursQuote holds the computed shift telemetry and excess overtime billing details.
type ExcessHoursQuote struct {
	AllowedHours          float64
	ActualGrossHours      float64
	CreditedBreakdownHours float64
	NetActualHours        float64
	BillableExcessHours   float64
	SurchargeAmountMinor  int64
}

// CalculateExcessHours computes the overtime billable hours and financial surcharge per RC-03.
func CalculateExcessHours(hireDays int, checkInMeter, checkOutMeter, breakdownCredit float64, hourlyRateMinor int64) ExcessHoursQuote {
	allowedHours := float64(hireDays) * 8.00
	grossHours := math.Max(0.00, checkOutMeter-checkInMeter)
	netActualHours := math.Max(0.00, grossHours-breakdownCredit)

	billableExcess := math.Max(0.00, netActualHours-allowedHours)
	surchargeMinor := int64(math.Floor(billableExcess * float64(hourlyRateMinor)))

	return ExcessHoursQuote{
		AllowedHours:           allowedHours,
		ActualGrossHours:       grossHours,
		CreditedBreakdownHours: breakdownCredit,
		NetActualHours:         netActualHours,
		BillableExcessHours:    billableExcess,
		SurchargeAmountMinor:   surchargeMinor,
	}
}

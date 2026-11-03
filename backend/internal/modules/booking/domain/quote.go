package domain

import (
	"fmt"
	"time"
)

type Quote struct {
	StartDate              time.Time
	EndDate                time.Time
	DurationDays           int
	DailyRateMinor         int64
	BaseHireFeeMinor       int64
	AttachmentsFeeMinor    int64
	GrossRentalAmountMinor int64
	PrepaymentAmountMinor  int64
	RemainingBalanceMinor  int64
	Currency               string
}

// CalculateQuote computes deterministic rental quote breakdown based on daily rate, duration, and attachments.
func CalculateQuote(startDate, endDate time.Time, dailyRateMinor int64, attachmentsDailyRateMinor int64, minRentalDays int, currency string) (Quote, error) {
	if !endDate.After(startDate) {
		return Quote{}, fmt.Errorf("%w: end date %s must be after start date %s", ErrInvalidDateRange, endDate.Format("2006-01-02"), startDate.Format("2006-01-02"))
	}

	durationDays := int(endDate.Sub(startDate).Hours() / 24)
	if durationDays < 1 {
		durationDays = 1
	}

	if durationDays < minRentalDays {
		return Quote{}, fmt.Errorf("%w: requested %d days, but fleet minimum is %d days", ErrInvalidDuration, durationDays, minRentalDays)
	}

	if dailyRateMinor <= 0 {
		return Quote{}, fmt.Errorf("daily rate must be positive, got %d", dailyRateMinor)
	}

	if currency == "" {
		currency = "ETB"
	}

	baseHireFee := int64(durationDays) * dailyRateMinor
	attachmentsFee := int64(durationDays) * attachmentsDailyRateMinor
	grossRentalAmount := baseHireFee + attachmentsFee
	prepaymentAmount := grossRentalAmount / 2 // Integer division floors exactly 50%
	remainingBalance := grossRentalAmount - prepaymentAmount

	return Quote{
		StartDate:              startDate,
		EndDate:                endDate,
		DurationDays:           durationDays,
		DailyRateMinor:         dailyRateMinor,
		BaseHireFeeMinor:       baseHireFee,
		AttachmentsFeeMinor:    attachmentsFee,
		GrossRentalAmountMinor: grossRentalAmount,
		PrepaymentAmountMinor:  prepaymentAmount,
		RemainingBalanceMinor:  remainingBalance,
		Currency:               currency,
	}, nil
}

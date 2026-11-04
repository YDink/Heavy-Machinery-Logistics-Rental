package domain

import (
	"fmt"
	"time"
)

// CalculateExtensionQuote computes the rate differential for extending an active rental.
func CalculateExtensionQuote(currentEndDate, newEndDate time.Time, dailyRateMinor, bundledAttachmentsDailyMinor int64) (extensionDays int, incrementalGross, incrementalPrepayment int64, err error) {
	currentNorm := time.Date(currentEndDate.Year(), currentEndDate.Month(), currentEndDate.Day(), 0, 0, 0, 0, time.UTC)
	newNorm := time.Date(newEndDate.Year(), newEndDate.Month(), newEndDate.Day(), 0, 0, 0, 0, time.UTC)

	if !newNorm.After(currentNorm) {
		return 0, 0, 0, fmt.Errorf("%w: new end date %s must be strictly after current end date %s",
			ErrInvalidExtensionDate, newEndDate.Format("2006-01-02"), currentEndDate.Format("2006-01-02"))
	}

	days := int(newNorm.Sub(currentNorm).Hours() / 24)
	if days <= 0 {
		return 0, 0, 0, fmt.Errorf("%w: extension days must be at least 1", ErrInvalidExtensionDate)
	}

	totalDailyMinor := dailyRateMinor + bundledAttachmentsDailyMinor
	incrementalGross = int64(days) * totalDailyMinor
	incrementalPrepayment = incrementalGross / 2 // Strict 50% integer minor units

	return days, incrementalGross, incrementalPrepayment, nil
}

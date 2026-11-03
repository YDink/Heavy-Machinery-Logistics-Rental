package domain

import (
	"testing"
	"time"
)

func TestCalculateQuote(t *testing.T) {
	start := time.Date(2026, 11, 10, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 11, 15, 0, 0, 0, 0, time.UTC) // 5 days

	dailyRate := int64(3500000)          // 35,000 ETB
	attachmentRate := int64(500000)      // 5,000 ETB

	quote, err := CalculateQuote(start, end, dailyRate, attachmentRate, 3, "ETB")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if quote.DurationDays != 5 {
		t.Errorf("expected 5 days, got %d", quote.DurationDays)
	}
	// Base hire fee = 5 * 35,000 = 175,000 ETB = 17,500,000 santims
	if quote.BaseHireFeeMinor != 17500000 {
		t.Errorf("expected base hire fee 17,500,000, got %d", quote.BaseHireFeeMinor)
	}
	// Attachments fee = 5 * 5,000 = 25,000 ETB = 2,500,000 santims
	if quote.AttachmentsFeeMinor != 2500000 {
		t.Errorf("expected attachments fee 2,500,000, got %d", quote.AttachmentsFeeMinor)
	}
	// Gross rental amount = 175,000 + 25,000 = 200,000 ETB = 20,000,000 santims
	if quote.GrossRentalAmountMinor != 20000000 {
		t.Errorf("expected gross amount 20,000,000, got %d", quote.GrossRentalAmountMinor)
	}
	// 50% advance prepayment = 100,000 ETB = 10,000,000 santims
	if quote.PrepaymentAmountMinor != 10000000 {
		t.Errorf("expected prepayment 10,000,000, got %d", quote.PrepaymentAmountMinor)
	}
	// Remaining balance = 10,000,000 santims
	if quote.RemainingBalanceMinor != 10000000 {
		t.Errorf("expected remaining balance 10,000,000, got %d", quote.RemainingBalanceMinor)
	}

	// Test below minimum duration
	shortEnd := time.Date(2026, 11, 12, 0, 0, 0, 0, time.UTC) // 2 days, min is 3
	if _, err := CalculateQuote(start, shortEnd, dailyRate, attachmentRate, 3, "ETB"); err == nil {
		t.Errorf("expected error for rental shorter than minimum days, got nil")
	}

	// Test invalid dates
	if _, err := CalculateQuote(end, start, dailyRate, attachmentRate, 1, "ETB"); err == nil {
		t.Errorf("expected error when end date is before start date, got nil")
	}
}

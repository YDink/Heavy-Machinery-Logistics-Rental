package domain

import (
	"errors"
	"testing"
	"time"
)

func TestCalculateExtensionQuote(t *testing.T) {
	now := time.Now().UTC()
	currentEnd := now.Add(48 * time.Hour)
	newEnd := currentEnd.Add(72 * time.Hour) // 3 days extension

	t.Run("Valid 3-day extension with 50% prepayment", func(t *testing.T) {
		extDays, grossMinor, prepayMinor, err := CalculateExtensionQuote(currentEnd, newEnd, 2000000, 500000)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if extDays != 3 {
			t.Fatalf("expected 3 days, got %d", extDays)
		}
		// 3 * 2,500,000 = 7,500,000
		if grossMinor != 7500000 {
			t.Fatalf("expected 7500000 gross minor, got %d", grossMinor)
		}
		// 7,500,000 / 2 = 3,750,000
		if prepayMinor != 3750000 {
			t.Fatalf("expected 3750000 prepay minor, got %d", prepayMinor)
		}
	})

	t.Run("Rejects non-forward date", func(t *testing.T) {
		pastDate := currentEnd.Add(-24 * time.Hour)
		_, _, _, err := CalculateExtensionQuote(currentEnd, pastDate, 2000000, 0)
		if err == nil {
			t.Fatal("expected error for past date, got nil")
		}
		if !errors.Is(err, ErrInvalidExtensionDate) {
			t.Fatalf("expected ErrInvalidExtensionDate, got: %v", err)
		}
	})
}

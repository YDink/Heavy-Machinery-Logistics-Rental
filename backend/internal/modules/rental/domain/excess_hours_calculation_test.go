package domain_test

import (
	"testing"

	"github.com/dinkmachinery/dink-core/internal/modules/rental/domain"
)

func TestCalculateExcessHours_WorkedExample(t *testing.T) {
	// 7 calendar days, 1450.50 check-in, 1512.20 check-out, 0 breakdown credit, 280,000 minor/hr
	quote := domain.CalculateExcessHours(7, 1450.50, 1512.20, 0.00, 280000)

	if quote.AllowedHours != 56.00 {
		t.Errorf("expected allowed hours 56.00, got %.2f", quote.AllowedHours)
	}
	expectedGross := 1512.20 - 1450.50
	if quote.ActualGrossHours < expectedGross-0.01 || quote.ActualGrossHours > expectedGross+0.01 {
		t.Errorf("expected gross hours %.2f, got %.2f", expectedGross, quote.ActualGrossHours)
	}
	if quote.BillableExcessHours < 5.69 || quote.BillableExcessHours > 5.71 {
		t.Errorf("expected billable excess hours ~5.70, got %.2f", quote.BillableExcessHours)
	}
	if quote.SurchargeAmountMinor != 1596000 {
		t.Errorf("expected surcharge minor 1596000 (15,960.00 ETB), got %d", quote.SurchargeAmountMinor)
	}
}

func TestCalculateExcessHours_WithBreakdownCredit(t *testing.T) {
	// 7 days (56h allowed), gross runtime 60h, 6h breakdown credit -> net 54h -> 0 excess
	quote := domain.CalculateExcessHours(7, 1000.00, 1060.00, 6.00, 280000)

	if quote.NetActualHours != 54.00 {
		t.Errorf("expected net actual hours 54.00, got %.2f", quote.NetActualHours)
	}
	if quote.BillableExcessHours != 0.00 {
		t.Errorf("expected billable excess hours 0.00, got %.2f", quote.BillableExcessHours)
	}
	if quote.SurchargeAmountMinor != 0 {
		t.Errorf("expected surcharge 0, got %d", quote.SurchargeAmountMinor)
	}
}

package domain_test

import (
	"testing"

	"github.com/dinkmachinery/dink-core/internal/modules/auth/domain"
)

func TestParsePhoneNumber_Ethiopia(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		defaultCountry  string
		expectedE164    string
		expectedCarrier domain.Carrier
		expectErr       bool
	}{
		{
			name:            "Valid Ethio Telecom E.164",
			input:           "+251911442910",
			defaultCountry:  "ET",
			expectedE164:    "+251911442910",
			expectedCarrier: domain.CarrierEthioTelecom,
			expectErr:       false,
		},
		{
			name:            "Valid Ethio Telecom domestic zero prefix",
			input:           "0911442910",
			defaultCountry:  "ET",
			expectedE164:    "+251911442910",
			expectedCarrier: domain.CarrierEthioTelecom,
			expectErr:       false,
		},
		{
			name:            "Valid Ethio Telecom with dashes and spaces",
			input:           "+251 91-144-2910",
			defaultCountry:  "ET",
			expectedE164:    "+251911442910",
			expectedCarrier: domain.CarrierEthioTelecom,
			expectErr:       false,
		},
		{
			name:            "Valid Safaricom Ethiopia E.164",
			input:           "+251711223344",
			defaultCountry:  "ET",
			expectedE164:    "+251711223344",
			expectedCarrier: domain.CarrierSafaricomET,
			expectErr:       false,
		},
		{
			name:            "Invalid Ethiopian prefix (not 9 or 7)",
			input:           "+251511223344",
			defaultCountry:  "ET",
			expectErr:       true,
		},
		{
			name:            "Invalid Ethiopian length (too short)",
			input:           "+2519112233",
			defaultCountry:  "ET",
			expectErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := domain.ParsePhoneNumber(tt.input, tt.defaultCountry)
			if tt.expectErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if p.String() != tt.expectedE164 {
				t.Errorf("expected E164 %s, got %s", tt.expectedE164, p.String())
			}
			if p.CountryCode() != "ET" {
				t.Errorf("expected country ET, got %s", p.CountryCode())
			}
			if p.Carrier() != tt.expectedCarrier {
				t.Errorf("expected carrier %s, got %s", tt.expectedCarrier, p.Carrier())
			}
		})
	}
}

func TestParsePhoneNumber_Kenya(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		defaultCountry  string
		expectedE164    string
		expectedCarrier domain.Carrier
		expectErr       bool
	}{
		{
			name:            "Valid Safaricom Kenya E.164",
			input:           "+254712345678",
			defaultCountry:  "KE",
			expectedE164:    "+254712345678",
			expectedCarrier: domain.CarrierSafaricomKE,
			expectErr:       false,
		},
		{
			name:            "Valid Airtel Kenya domestic zero prefix",
			input:           "0112345678",
			defaultCountry:  "KE",
			expectedE164:    "+254112345678",
			expectedCarrier: domain.CarrierAirtelKE,
			expectErr:       false,
		},
		{
			name:            "Unsupported country prefix",
			input:           "+14155552671",
			defaultCountry:  "US",
			expectErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := domain.ParsePhoneNumber(tt.input, tt.defaultCountry)
			if tt.expectErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if p.String() != tt.expectedE164 {
				t.Errorf("expected E164 %s, got %s", tt.expectedE164, p.String())
			}
			if p.CountryCode() != "KE" {
				t.Errorf("expected country KE, got %s", p.CountryCode())
			}
			if p.Carrier() != tt.expectedCarrier {
				t.Errorf("expected carrier %s, got %s", tt.expectedCarrier, p.Carrier())
			}
		})
	}
}

func TestPhoneNumber_Masked(t *testing.T) {
	p, err := domain.ParsePhoneNumber("+251911442910", "ET")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	masked := p.Masked()
	expected := "+25191****910"
	if masked != expected {
		t.Errorf("expected masked %s, got %s", expected, masked)
	}
}

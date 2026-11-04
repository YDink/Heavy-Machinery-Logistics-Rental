package domain

import "testing"

func TestGenerateChallengePIN(t *testing.T) {
	pin, err := GenerateChallengePIN()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pin) != 6 {
		t.Fatalf("expected 6 digits, got %d (%s)", len(pin), pin)
	}
}

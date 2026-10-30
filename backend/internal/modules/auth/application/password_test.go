package application_test

import (
	"testing"

	"github.com/dinkmachinery/dink-core/internal/modules/auth/application"
)

func TestArgon2id_HashAndVerify(t *testing.T) {
	// Lightweight params for test speed
	testParams := &application.Argon2idParams{
		Memory:      16 * 1024,
		Iterations:  1,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}

	password := "CorrectHorseBatteryStaple2026!"
	encoded, err := application.HashPassword(password, testParams)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	if encoded == "" {
		t.Fatalf("expected non-empty hash string")
	}

	// Positive verification
	valid, err := application.VerifyPassword(password, encoded)
	if err != nil {
		t.Fatalf("expected successful verification, got error: %v", err)
	}
	if !valid {
		t.Fatalf("expected valid=true for correct password")
	}

	// Negative verification with wrong password
	valid, err = application.VerifyPassword("WrongPassword123!", encoded)
	if err == nil || valid {
		t.Fatalf("expected verification failure for wrong password")
	}
}

func TestArgon2id_MalformedHash(t *testing.T) {
	malformed := "$argon2id$v=19$invalid_hash"
	_, err := application.VerifyPassword("password", malformed)
	if err == nil {
		t.Fatalf("expected error for malformed hash, got nil")
	}
}

package domain

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

const (
	HandoverChallengeMaxAttempts = 3
	HandoverChallengeTTL         = 15 * time.Minute
)

// GenerateChallengePIN produces an unpredictable, cryptographically random 6-digit numeric challenge PIN.
func GenerateChallengePIN() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return "", fmt.Errorf("failed to generate secure challenge PIN: %w", err)
	}
	return fmt.Sprintf("%06d", n.Int64()+100000), nil
}

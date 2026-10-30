package application

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

var (
	ErrInvalidPasswordHash = errors.New("invalid password hash format")
	ErrPasswordMismatch    = errors.New("password does not match hash")
)

// Argon2idParams defines cryptographic parameters strictly matching SRS Section 3.2 (REG-07).
type Argon2idParams struct {
	Memory      uint32 // 64 MB (65,536 KiB)
	Iterations  uint32 // 3 iterations
	Parallelism uint8  // 4 threads
	SaltLength  uint32 // 16 bytes
	KeyLength   uint32 // 32 bytes
}

// DefaultArgon2idParams returns standard production parameters.
func DefaultArgon2idParams() *Argon2idParams {
	return &Argon2idParams{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 4,
		SaltLength:  16,
		KeyLength:   32,
	}
}

// HashPassword hashes a plaintext password using Argon2id with cryptographic salt.
func HashPassword(password string, params *Argon2idParams) (string, error) {
	if params == nil {
		params = DefaultArgon2idParams()
	}

	salt := make([]byte, params.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate random salt: %w", err)
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		params.Iterations,
		params.Memory,
		params.Parallelism,
		params.KeyLength,
	)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Format: $argon2id$v=19$m=65536,t=3,p=4$<salt>$<hash>
	encoded := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		params.Memory,
		params.Iterations,
		params.Parallelism,
		b64Salt,
		b64Hash,
	)

	return encoded, nil
}

// VerifyPassword verifies an encoded Argon2id hash against a plaintext password using constant-time comparison.
func VerifyPassword(password, encodedHash string) (bool, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false, ErrInvalidPasswordHash
	}

	if parts[1] != "argon2id" {
		return false, fmt.Errorf("%w: unsupported algorithm %s", ErrInvalidPasswordHash, parts[1])
	}

	var version int
	_, err := fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil || version != argon2.Version {
		return false, fmt.Errorf("%w: incompatible argon2 version", ErrInvalidPasswordHash)
	}

	var memory, iterations uint32
	var parallelism uint8
	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism)
	if err != nil {
		return false, fmt.Errorf("%w: malformed parameters", ErrInvalidPasswordHash)
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, fmt.Errorf("%w: invalid salt encoding", ErrInvalidPasswordHash)
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, fmt.Errorf("%w: invalid hash encoding", ErrInvalidPasswordHash)
	}

	actualHash := argon2.IDKey(
		[]byte(password),
		salt,
		iterations,
		memory,
		parallelism,
		uint32(len(expectedHash)),
	)

	if subtle.ConstantTimeCompare(actualHash, expectedHash) == 1 {
		return true, nil
	}

	return false, ErrPasswordMismatch
}

package application_test

import (
	"context"
	"testing"

	"go.uber.org/zap"
	"github.com/dinkmachinery/dink-core/internal/modules/auth/application"
	"github.com/dinkmachinery/dink-core/internal/modules/auth/domain"
	"github.com/dinkmachinery/dink-core/internal/modules/auth/infrastructure"
)

func TestOTPService_SuccessFlow(t *testing.T) {
	ctx := context.Background()
	cache := infrastructure.NewMemoryCache()
	logger := zap.NewNop()
	smsGateway := infrastructure.NewMockSMSGateway(logger)

	otpService := application.NewOTPService(cache, smsGateway, application.OTPServiceConfig{MockOTP: true})

	phone, _ := domain.ParsePhoneNumber("+251911442910", "ET")
	sessionID, timeout, err := otpService.RequestOTP(ctx, phone, domain.OTPPurposeLogin)
	if err != nil {
		t.Fatalf("unexpected error requesting OTP: %v", err)
	}
	if sessionID == "" {
		t.Fatalf("expected valid sessionID")
	}
	if timeout <= 0 {
		t.Fatalf("expected positive timeout seconds, got %d", timeout)
	}

	// Verify with correct mock OTP
	challenge, err := otpService.VerifyOTP(ctx, sessionID, "123456")
	if err != nil {
		t.Fatalf("expected successful verification, got error: %v", err)
	}
	if challenge.Phone != phone.String() {
		t.Errorf("expected challenge phone %s, got %s", phone.String(), challenge.Phone)
	}

	// Verify challenge cannot be reused
	_, err = otpService.VerifyOTP(ctx, sessionID, "123456")
	if err != domain.ErrOTPExpired {
		t.Errorf("expected ErrOTPExpired for consumed challenge, got %v", err)
	}
}

// TestOTPService_MaxAttemptsLockout asserts AUT-03:
// After 3 failed attempts, the challenge is deleted and locked out.
func TestOTPService_MaxAttemptsLockout(t *testing.T) {
	ctx := context.Background()
	cache := infrastructure.NewMemoryCache()
	logger := zap.NewNop()
	smsGateway := infrastructure.NewMockSMSGateway(logger)

	otpService := application.NewOTPService(cache, smsGateway, application.OTPServiceConfig{MockOTP: true})

	phone, _ := domain.ParsePhoneNumber("+251911442910", "ET")
	sessionID, _, err := otpService.RequestOTP(ctx, phone, domain.OTPPurposeLogin)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Attempt 1: incorrect
	_, err = otpService.VerifyOTP(ctx, sessionID, "000000")
	if err != domain.ErrOTPInvalid {
		t.Fatalf("attempt 1: expected ErrOTPInvalid, got: %v", err)
	}

	// Attempt 2: incorrect
	_, err = otpService.VerifyOTP(ctx, sessionID, "111111")
	if err != domain.ErrOTPInvalid {
		t.Fatalf("attempt 2: expected ErrOTPInvalid, got: %v", err)
	}

	// Attempt 3: incorrect -> triggers lockout
	_, err = otpService.VerifyOTP(ctx, sessionID, "222222")
	if err != domain.ErrOTPMaxAttemptsExceeded {
		t.Fatalf("attempt 3: expected ErrOTPMaxAttemptsExceeded, got: %v", err)
	}

	// Challenge should now be completely deleted from cache
	_, err = otpService.VerifyOTP(ctx, sessionID, "123456")
	if err != domain.ErrOTPExpired {
		t.Fatalf("subsequent attempt: expected ErrOTPExpired, got: %v", err)
	}
}

// TestOTPService_RateLimiting asserts D-05-01-001:
// Leaky bucket rate limiter blocks excessive dispatches within window.
func TestOTPService_RateLimiting(t *testing.T) {
	ctx := context.Background()
	cache := infrastructure.NewMemoryCache()
	logger := zap.NewNop()
	smsGateway := infrastructure.NewMockSMSGateway(logger)

	otpService := application.NewOTPService(cache, smsGateway, application.OTPServiceConfig{MockOTP: true})
	phone, _ := domain.ParsePhoneNumber("+251911442910", "ET")

	// 3 dispatches allowed within window
	for i := 1; i <= 3; i++ {
		_, _, err := otpService.RequestOTP(ctx, phone, domain.OTPPurposeLogin)
		if err != nil {
			t.Fatalf("dispatch %d should succeed, got error: %v", i, err)
		}
	}

	// 4th dispatch exceeds window rate limit
	_, _, err := otpService.RequestOTP(ctx, phone, domain.OTPPurposeLogin)
	if err != domain.ErrRateLimitExceeded {
		t.Fatalf("4th dispatch should fail with ErrRateLimitExceeded, got: %v", err)
	}
}

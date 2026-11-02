package infrastructure

import (
	"context"
	"log/slog"

	"github.com/dinkmachinery/dink-core/internal/modules/auth/domain"
)

// MockSMSGateway logs dispatched OTPs without making external network calls.
type MockSMSGateway struct {
	logger *slog.Logger
}

func NewMockSMSGateway(logger *slog.Logger) *MockSMSGateway {
	if logger == nil {
		logger = slog.Default()
	}
	return &MockSMSGateway{logger: logger}
}

func (g *MockSMSGateway) SendOTP(ctx context.Context, phone domain.PhoneNumber, code string) error {
	// Zero-PII logging: Mask phone number in logs, omit plaintext code from production logs
	g.logger.InfoContext(ctx, "dispatched SMS OTP challenge",
		slog.String("masked_phone", phone.Masked()),
		slog.String("carrier", string(phone.Carrier())),
		slog.String("country", phone.CountryCode()),
	)
	return nil
}

package infrastructure

import (
	"context"

	"go.uber.org/zap"
	"github.com/dinkmachinery/dink-core/internal/modules/auth/domain"
)

// MockSMSGateway logs dispatched OTPs without making external network calls.
type MockSMSGateway struct {
	logger *zap.Logger
}

func NewMockSMSGateway(logger *zap.Logger) *MockSMSGateway {
	return &MockSMSGateway{logger: logger}
}

func (g *MockSMSGateway) SendOTP(ctx context.Context, phone domain.PhoneNumber, code string) error {
	// Zero-PII logging: Mask phone number in logs, omit plaintext code from production logs
	g.logger.Info("dispatched SMS OTP challenge",
		zap.String("masked_phone", phone.Masked()),
		zap.String("carrier", string(phone.Carrier())),
		zap.String("country", phone.CountryCode()),
	)
	return nil
}

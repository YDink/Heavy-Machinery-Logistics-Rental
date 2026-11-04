package domain

import "context"

// SMSGateway defines the messaging abstraction for domestic carriers.
type SMSGateway interface {
	SendOTP(ctx context.Context, phone PhoneNumber, code string) error
}

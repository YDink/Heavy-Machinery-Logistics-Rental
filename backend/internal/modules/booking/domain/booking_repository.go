package domain

import (
	"context"

	"github.com/google/uuid"
)

type BookingRepository interface {
	CreateBooking(ctx context.Context, booking *Booking) error
	GetBookingByID(ctx context.Context, id uuid.UUID) (*Booking, error)
	GetBookingByReference(ctx context.Context, reference string) (*Booking, error)
	ListBookingsByRenter(ctx context.Context, renterID uuid.UUID, status string, limit int) ([]*Booking, error)
	ListBookingsByOwner(ctx context.Context, ownerID uuid.UUID, status string, limit int) ([]*Booking, error)
	UpdateBookingStatusOCC(ctx context.Context, booking *Booking, expectedVersion int) error
	UpdateBookingHandover(ctx context.Context, bookingID uuid.UUID, checkInMeter float64, expectedVersion int) error
	UpdateBookingReturn(ctx context.Context, bookingID uuid.UUID, status BookingStatus, checkOutMeter float64, expectedVersion int) error

	AddBookingAttachment(ctx context.Context, bookingID, attachmentID uuid.UUID, dailyRateMinor int64) error
	ListBookingAttachments(ctx context.Context, bookingID uuid.UUID) ([]BookingAttachment, error)

	CreateBankSlip(ctx context.Context, slip *BankSlip) error
	GetLatestBankSlip(ctx context.Context, bookingID uuid.UUID) (*BankSlip, error)
	UpdateBankSlipStatus(ctx context.Context, slipID uuid.UUID, status BankSlipStatus, reasonCode, notes *string) error

	CreateAuditLog(ctx context.Context, log *AuditLog) error
	ListAuditLogs(ctx context.Context, bookingID uuid.UUID) ([]*AuditLog, error)
}

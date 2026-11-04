package http

import (
	openapi_types "github.com/oapi-codegen/runtime/types"
	bookingDomain "github.com/dinkmachinery/dink-core/internal/modules/booking/domain"
	"github.com/dinkmachinery/dink-core/internal/transport/http/generated"
)

func toGeneratedBooking(b *bookingDomain.Booking) generated.Booking {
	dto := generated.Booking{
		Id:               openapi_types.UUID(b.ID),
		BookingReference: b.BookingReference,
		FleetId:          openapi_types.UUID(b.FleetID),
		RenterId:         openapi_types.UUID(b.RenterID),
		OwnerId:          openapi_types.UUID(b.OwnerID),
		Status:           generated.BookingStatus(b.Status),
		StartDate:        openapi_types.Date{Time: b.StartDate},
		EndDate:          openapi_types.Date{Time: b.EndDate},
		DurationDays:     b.DurationDays,
		DailyRate: generated.MoneyAmount{
			AmountMinor: b.DailyRateAppliedMinor,
			Currency:    b.Currency,
		},
		BaseHireFee: generated.MoneyAmount{
			AmountMinor: b.BaseHireFeeMinor,
			Currency:    b.Currency,
		},
		AttachmentsFee: generated.MoneyAmount{
			AmountMinor: b.AttachmentsFeeMinor,
			Currency:    b.Currency,
		},
		GrossRentalAmount: generated.MoneyAmount{
			AmountMinor: b.GrossRentalAmountMinor,
			Currency:    b.Currency,
		},
		PrepaymentAmount: generated.MoneyAmount{
			AmountMinor: b.PrepaymentAmountMinor,
			Currency:    b.Currency,
		},
		SiteLocationAddress: b.SiteLocationAddress,
		RejectionReasonCode: b.RejectionReasonCode,
		RejectionNotes:      b.RejectionNotes,
		CancellationReason:  b.CancellationReason,
		HandoverPin:         b.HandoverPIN,
		ExpiresAt:           b.ExpiresAt,
		ConfirmedAt:         b.ConfirmedAt,
		ActivatedAt:         b.ActivatedAt,
		CompletedAt:         b.CompletedAt,
		CancelledAt:         b.CancelledAt,
		Version:             b.Version,
		CreatedAt:           b.CreatedAt,
		UpdatedAt:           &b.UpdatedAt,
	}

	if b.PenaltyAmountMinor > 0 {
		dto.PenaltyAmount = &generated.MoneyAmount{
			AmountMinor: b.PenaltyAmountMinor,
			Currency:    b.Currency,
		}
	}
	if b.RefundAmountMinor > 0 {
		dto.RefundAmount = &generated.MoneyAmount{
			AmountMinor: b.RefundAmountMinor,
			Currency:    b.Currency,
		}
	}
	if b.CheckInMeterHours != nil {
		m := float32(*b.CheckInMeterHours)
		dto.CheckInMeterHours = &m
	}
	if b.CheckOutMeterHours != nil {
		m := float32(*b.CheckOutMeterHours)
		dto.CheckOutMeterHours = &m
	}

	if len(b.Attachments) > 0 {
		var attSummaries []generated.BookingAttachmentSummary
		for _, att := range b.Attachments {
			attSummaries = append(attSummaries, generated.BookingAttachmentSummary{
				Id:             openapi_types.UUID(att.ID),
				AttachmentId:   openapi_types.UUID(att.AttachmentID),
				AttachmentType: &att.AttachmentType,
				Make:           &att.Make,
				Model:          &att.Model,
				SerialNumber:   &att.SerialNumber,
				BundledDailyRate: generated.MoneyAmount{
					AmountMinor: att.BundledDailyRateMinor,
					Currency:    b.Currency,
				},
			})
		}
		dto.Attachments = &attSummaries
	}

	return dto
}

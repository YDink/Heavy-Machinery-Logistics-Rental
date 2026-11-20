package http

import (
	"net/http"

	"github.com/dinkmachinery/dink-core/internal/platform/response"
	"github.com/dinkmachinery/dink-core/internal/transport/http/generated"
	"github.com/google/uuid"
)

// GetBookingLogisticsOverview handles GET /api/v1/bookings/{id}/logistics
func (s *Server) GetBookingLogisticsOverview(w http.ResponseWriter, r *http.Request) {
	if s.logisticsService == nil {
		s.notImplemented(w, r, "05-08")
		return
	}

	userUUID, err := s.authenticateUser(r)
	if err != nil {
		response.Error(w, r, http.StatusUnauthorized, response.ProblemDetails{
			Title:  "Unauthorized",
			Code:   "ERR_UNAUTHORIZED",
			Detail: err.Error(),
		})
		return
	}

	bookingID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Bad Request",
			Code:   "ERR_INVALID_BOOKING_ID",
			Detail: "Invalid booking ID in path.",
		})
		return
	}

	res, err := s.logisticsService.GetLogisticsOverview(r.Context(), bookingID, userUUID)
	if err != nil {
		writeLogisticsDomainError(w, r, err)
		return
	}

	data := generated.LogisticsOverviewResponseData{
		ParentBookingId:     res.ParentBookingID,
		ParentBookingStatus: string(res.ParentBookingStatus),
		SelfMoveCapable:     res.SelfMoveCapable,
		LogisticsMode:       generated.LogisticsOverviewResponseDataLogisticsMode(res.LogisticsMode),
	}

	if res.Dispatch != nil && res.SubBooking != nil {
		carrierName := "Platform Carrier"
		carrierPhone := ""
		if res.CarrierContact != nil {
			carrierName = res.CarrierContact.OwnerFullName
			carrierPhone = res.CarrierContact.OwnerPhoneNumber
		}
		escortFee := res.Dispatch.EscortFeeMinor
		data.SubBooking = &generated.LogisticsSubBookingResponseData{
			SubBookingId:               res.SubBooking.ID,
			ParentBookingId:            res.ParentBookingID,
			Status:                     string(res.SubBooking.Status),
			FreightHireFeeMinor:        res.Dispatch.FreightHireFeeMinor,
			PlatformCommissionBps:      int32(res.Dispatch.PlatformCommissionBps),
			PlatformCommissionRate:     "0.025",
			CommissionOwedMinor:        res.Dispatch.CommissionOwedMinor,
			Currency:                   res.SubBooking.Currency,
			AbnormalLoadEscortRequired: res.Dispatch.AbnormalLoadEscortRequired,
			EscortFeeMinor:             &escortFee,
			CarrierContact: generated.CarrierContactInfo{
				FullName:    carrierName,
				PhoneNumber: carrierPhone,
			},
		}
	}

	if res.ExternalTransport != nil {
		data.ExternalTransport = &generated.ExternalTransportResponseData{
			ParentBookingId:             res.ExternalTransport.ParentBookingID,
			Status:                      string(res.ParentBookingStatus),
			ExternalTransportRegistered: true,
			CarrierName:                 res.ExternalTransport.CarrierName,
			DriverFullName:              res.ExternalTransport.DriverFullName,
			DriverPhoneNumber:           res.ExternalTransport.DriverPhoneNumber,
			PrimeMoverPlateNumber:       res.ExternalTransport.PrimeMoverPlateNumber,
			EstimatedArrivalDateTime:    res.ExternalTransport.EstimatedArrivalAt,
		}
	}

	if res.ActiveTransshipment != nil {
		subStatus := ""
		if res.SubBooking != nil {
			subStatus = string(res.SubBooking.Status)
		}
		data.ActiveTransshipment = &generated.TransshipmentAlertResponseData{
			TransshipmentId:          res.ActiveTransshipment.ID,
			SubBookingId:             res.ActiveTransshipment.SubBookingID,
			ParentBookingId:          res.ActiveTransshipment.ParentBookingID,
			SubBookingStatus:         subStatus,
			ParentBookingStatus:      string(res.ParentBookingStatus),
			HotSwapSlaDeadline:       res.ActiveTransshipment.HotSwapSLADeadline,
			TransshipmentBountyMinor: res.ActiveTransshipment.TransshipmentBountyMinor,
			Currency:                 res.ActiveTransshipment.Currency,
			CandidateCarriersCount:   0,
		}
	}

	if res.ReverseCascade != nil {
		subStatus := ""
		if res.SubBooking != nil {
			subStatus = string(res.SubBooking.Status)
		}
		data.ReverseCascade = &generated.ReverseCascadeResponseData{
			ParentBookingId:                 res.ReverseCascade.ParentBookingID,
			SubBookingId:                    res.ReverseCascade.SubBookingID,
			Resolution:                      res.ReverseCascade.Resolution,
			RenterRefundDueMinor:            res.ReverseCascade.RenterRefundDueMinor,
			OwnerPenaltyLeviedMinor:         res.ReverseCascade.OwnerPenaltyLeviedMinor,
			OwnerCommissionRestitutionMinor: res.ReverseCascade.OwnerCommissionRestitutionMinor,
			CarrierDefaultPenaltyMinor:      res.ReverseCascade.CarrierDefaultPenaltyMinor,
			Currency:                        res.ReverseCascade.Currency,
			ParentStatus:                    string(res.ParentBookingStatus),
			SubStatus:                       subStatus,
		}
	}

	response.JSON(w, r, http.StatusOK, data)
}

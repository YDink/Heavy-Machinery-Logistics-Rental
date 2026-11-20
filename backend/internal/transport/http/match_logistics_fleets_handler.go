package http

import (
	"encoding/json"
	"net/http"

	logisticsApp "github.com/dinkmachinery/dink-core/internal/modules/logistics/application"
	"github.com/dinkmachinery/dink-core/internal/platform/response"
	"github.com/dinkmachinery/dink-core/internal/transport/http/generated"
	"github.com/google/uuid"
)

// MatchLogisticsFleets handles POST /api/v1/bookings/{id}/logistics/match
func (s *Server) MatchLogisticsFleets(w http.ResponseWriter, r *http.Request) {
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

	var req generated.LogisticsMatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Invalid Request Payload",
			Code:   "ERR_INVALID_JSON",
			Detail: "Request body could not be parsed.",
		})
		return
	}

	res, err := s.logisticsService.MatchCarriers(r.Context(), logisticsApp.MatchCarriersCommand{
		ParentBookingID: bookingID,
		CallerUserID:    userUUID,
		SearchRadiusKm:  int(req.SearchRadiusKm),
		PickupDate:      req.PickupDate,
	})
	if err != nil {
		writeLogisticsDomainError(w, r, err)
		return
	}

	fleets := make([]generated.MatchedLogisticsFleet, len(res.AvailableLogisticsFleets))
	for i, c := range res.AvailableLogisticsFleets {
		fleets[i] = generated.MatchedLogisticsFleet{
			FleetId:                 c.FleetID,
			Make:                    c.Make,
			Model:                   c.Model,
			OwnerTradingName:        c.OwnerTradingName,
			PayloadCapacityKg:       c.PayloadCapacityKg,
			DeckWellLengthMm:        int32(c.DeckWellLengthMm),
			DeckWidthMm:             int32(c.DeckWidthMm),
			TareWeightKg:            c.TareWeightKg,
			AxleCount:               int32(c.AxleCount),
			DailyRateMinor:          c.DailyRateMinor,
			Currency:                c.Currency,
			DistanceFromDepotMeters: c.DistanceFromDepotMeters,
			AverageRating:           c.AverageRating,
		}
	}

	data := generated.LogisticsMatchResponseData{
		ParentBookingId:            res.ParentBookingID,
		SelfMoveCapable:            res.SelfMoveCapable,
		RecommendedTrailerClass:    res.RecommendedTrailerClass,
		AbnormalLoadEscortRequired: res.AbnormalLoadEscortRequired,
		MachineSpecs: generated.MachineLogisticsSpecs{
			OperatingWeightKg: res.MachineSpecs.OperatingWeightKg,
			TransportLengthMm: int32(res.MachineSpecs.TransportLengthMm),
			TransportWidthMm:  int32(res.MachineSpecs.TransportWidthMm),
			TransportHeightMm: int32(res.MachineSpecs.TransportHeightMm),
		},
		RequiredMinimumPayloadKg: res.RequiredMinimumPayloadKg,
		MatchingCarriersCount:    int32(len(fleets)),
		AvailableLogisticsFleets: fleets,
	}

	response.JSON(w, r, http.StatusOK, data)
}

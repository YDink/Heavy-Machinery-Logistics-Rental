package http

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	machineryApp "github.com/dinkmachinery/dink-core/internal/modules/machinery/application"
	"github.com/dinkmachinery/dink-core/internal/platform/response"
	"github.com/dinkmachinery/dink-core/internal/transport/http/generated"
)

// CheckBundledAvailability handles POST /api/v1/bookings/check-bundled-availability
func (s *Server) CheckBundledAvailability(w http.ResponseWriter, r *http.Request) {
	var body generated.CheckBundledAvailabilityJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Bad Request",
			Code:   "ERR_MALFORMED_JSON",
			Detail: "Request body could not be parsed as valid JSON.",
		})
		return
	}

	fleetUUID := uuid.UUID(body.FleetId)
	var attUUIDs []uuid.UUID
	for _, a := range body.AttachmentIds {
		attUUIDs = append(attUUIDs, uuid.UUID(a))
	}

	startDate := body.StartDate.Time
	endDate := body.EndDate.Time

	out, err := s.attachmentService.CheckBundledAvailability(r.Context(), machineryApp.CheckBundledAvailabilityInput{
		FleetID:       fleetUUID,
		AttachmentIDs: attUUIDs,
		StartDate:     startDate,
		EndDate:       endDate,
	})
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Availability Check Failed",
			Code:   "ERR_AVAILABILITY_CHECK",
			Detail: err.Error(),
		})
		return
	}

	statuses := make([]generated.AttachmentAvailabilityStatus, len(out.AttachmentStatuses))
	for i, st := range out.AttachmentStatuses {
		statuses[i] = generated.AttachmentAvailabilityStatus{
			AttachmentId:          openapi_types.UUID(st.AttachmentID),
			IsAvailable:           st.IsAvailable,
			DailyRateAppliedMinor: st.DailyRateAppliedMinor,
		}
	}

	data := generated.CheckBundledAvailabilityResponseData{
		FleetAvailable:              out.FleetAvailable,
		AllAttachmentsAvailable:     out.AllAttachmentsAvailable,
		AttachmentStatuses:          statuses,
		TotalCombinedDailyRateMinor: out.TotalCombinedDailyRateMinor,
	}

	response.JSON(w, r, http.StatusOK, data)
}

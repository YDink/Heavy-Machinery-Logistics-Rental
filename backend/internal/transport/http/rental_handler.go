package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	bookingDomain "github.com/dinkmachinery/dink-core/internal/modules/booking/domain"
	rentalApp "github.com/dinkmachinery/dink-core/internal/modules/rental/application"
	rentalDomain "github.com/dinkmachinery/dink-core/internal/modules/rental/domain"
	"github.com/dinkmachinery/dink-core/internal/platform/response"
	"github.com/dinkmachinery/dink-core/internal/transport/http/generated"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func (s *Server) authenticateUser(r *http.Request) (uuid.UUID, error) {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return uuid.Nil, errors.New("missing or malformed Authorization header")
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := s.tokenService.VerifyAccessToken(r.Context(), tokenStr)
	if err != nil {
		return uuid.Nil, errors.New("invalid or expired token")
	}

	userUUID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, errors.New("invalid subject in token")
	}

	return userUUID, nil
}

// HandoverInspect handles POST /api/v1/bookings/{id}/handover-inspect
func (s *Server) HandoverInspect(w http.ResponseWriter, r *http.Request) {
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

	var req generated.HandoverInspectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Bad Request",
			Code:   "ERR_INVALID_BODY",
			Detail: "Malformed request payload.",
		})
		return
	}

	photos := make([]rentalDomain.InspectionPhoto, len(req.InspectionPhotos))
	for i, p := range req.InspectionPhotos {
		photos[i] = rentalDomain.InspectionPhoto{
			Angle:      string(p.Angle),
			StorageURI: p.StorageUri,
		}
	}

	notes := ""
	if req.Notes != nil {
		notes = *req.Notes
	}

	fluidLevelsOk := true
	if req.FluidLevelsOk != nil {
		fluidLevelsOk = *req.FluidLevelsOk
	}
	safetyGlassIntact := true
	if req.SafetyGlassIntact != nil {
		safetyGlassIntact = *req.SafetyGlassIntact
	}

	inspection, err := s.rentalService.InspectHandover(r.Context(), rentalApp.HandoverInspectCommand{
		BookingID:         bookingID,
		InspectorID:       userUUID,
		CheckInMeterHours: req.CheckInMeterHours,
		FluidLevelsOk:     fluidLevelsOk,
		SafetyGlassIntact: safetyGlassIntact,
		Notes:             notes,
		Photos:            photos,
	})
	if err != nil {
		if errors.Is(err, rentalDomain.ErrMeterReadingRegression) {
			response.Error(w, r, http.StatusUnprocessableEntity, response.ProblemDetails{
				Title:  "Meter Reading Regression",
				Code:   "ERR_METER_READING_REGRESSION",
				Detail: err.Error(),
			})
			return
		}
		if strings.Contains(err.Error(), "photo") {
			response.Error(w, r, http.StatusUnprocessableEntity, response.ProblemDetails{
				Title:  "Invalid Inspection Photos",
				Code:   "ERR_INVALID_PHOTOS",
				Detail: err.Error(),
			})
			return
		}
		if errors.Is(err, bookingDomain.ErrUnauthorizedAction) {
			response.Error(w, r, http.StatusForbidden, response.ProblemDetails{
				Title:  "Forbidden",
				Code:   "ERR_FORBIDDEN",
				Detail: err.Error(),
			})
			return
		}
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Inspection Failed",
			Code:   "ERR_INSPECTION_FAILED",
			Detail: err.Error(),
		})
		return
	}

	data := generated.HandoverInspectionResponseData{
		BookingId:           inspection.BookingID,
		InspectedAt:         inspection.CreatedAt,
		MeterVerified:       true,
		PhotosIngestedCount: len(inspection.Photos),
		Status:              "INSPECTED",
	}

	response.JSON(w, r, http.StatusOK, data)
}

// GenerateHandoverChallenge handles POST /api/v1/bookings/{id}/handover-challenge
func (s *Server) GenerateHandoverChallenge(w http.ResponseWriter, r *http.Request) {
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

	challenge, err := s.rentalService.GenerateHandoverPIN(r.Context(), bookingID, userUUID)
	if err != nil {
		if errors.Is(err, bookingDomain.ErrUnauthorizedAction) {
			response.Error(w, r, http.StatusForbidden, response.ProblemDetails{
				Title:  "Forbidden",
				Code:   "ERR_FORBIDDEN",
				Detail: "Only the fleet owner can generate a handover challenge PIN.",
			})
			return
		}
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Challenge Generation Failed",
			Code:   "ERR_CHALLENGE_GENERATION_FAILED",
			Detail: err.Error(),
		})
		return
	}

	attemptCount := challenge.AttemptCount
	data := generated.HandoverChallengeResponseData{
		AttemptCount: &attemptCount,
		BookingId:    challenge.BookingID,
		ChallengePin: challenge.ChallengePIN,
		ExpiresAt:    challenge.ExpiresAt,
		MaxAttempts:  challenge.MaxAttempts,
	}

	response.JSON(w, r, http.StatusOK, data)
}

// ConfirmHandover handles POST /api/v1/bookings/{id}/handover-confirm
func (s *Server) ConfirmHandover(w http.ResponseWriter, r *http.Request) {
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

	var req generated.HandoverConfirmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Bad Request",
			Code:   "ERR_INVALID_BODY",
			Detail: "Malformed request payload.",
		})
		return
	}

	var touchSig string
	if req.TouchSignatureBase64 != nil {
		touchSig = *req.TouchSignatureBase64
	}
	var lat, lng *float64
	if req.ContractorLocation != nil {
		lat = req.ContractorLocation.Latitude
		lng = req.ContractorLocation.Longitude
	}

	activatedBooking, err := s.rentalService.ConfirmHandover(r.Context(), rentalApp.HandoverConfirmCommand{
		BookingID:            bookingID,
		CallerUserID:         userUUID,
		PIN:                  req.ChallengePin,
		TouchSignatureBase64: touchSig,
		ContractorLat:        lat,
		ContractorLng:        lng,
	})
	if err != nil {
		if errors.Is(err, rentalDomain.ErrPINLockout) {
			response.Error(w, r, http.StatusLocked, response.ProblemDetails{
				Title:  "Handover Challenge Locked Out",
				Code:   "ERR_PIN_LOCKOUT",
				Detail: err.Error(),
			})
			return
		}
		if errors.Is(err, rentalDomain.ErrPINExpired) {
			response.Error(w, r, http.StatusGone, response.ProblemDetails{
				Title:  "Challenge PIN Expired",
				Code:   "ERR_PIN_EXPIRED",
				Detail: err.Error(),
			})
			return
		}
		if errors.Is(err, rentalDomain.ErrInvalidChallengePIN) {
			response.Error(w, r, http.StatusUnprocessableEntity, response.ProblemDetails{
				Title:  "Invalid Challenge PIN",
				Code:   "ERR_INVALID_CHALLENGE_PIN",
				Detail: err.Error(),
			})
			return
		}
		if errors.Is(err, bookingDomain.ErrUnauthorizedAction) {
			response.Error(w, r, http.StatusForbidden, response.ProblemDetails{
				Title:  "Forbidden",
				Code:   "ERR_FORBIDDEN",
				Detail: "Only renter or assigned operator may verify handover.",
			})
			return
		}
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Handover Confirmation Failed",
			Code:   "ERR_CONFIRMATION_FAILED",
			Detail: err.Error(),
		})
		return
	}

	actTime := time.Now().UTC()
	if activatedBooking.ActivatedAt != nil {
		actTime = *activatedBooking.ActivatedAt
	}

	startMeter := 0.0
	if activatedBooking.CheckInMeterHours != nil {
		startMeter = *activatedBooking.CheckInMeterHours
	}

	data := generated.HandoverConfirmResponseData{
		ActivatedAt:              actTime,
		AllocatedShiftHoursTotal: rentalDomain.StandardShiftHoursPerDay,
		BookingId:                activatedBooking.ID,
		ShiftBaselineMeterHours:  startMeter,
		Status:                   string(activatedBooking.Status),
	}

	response.JSON(w, r, http.StatusOK, data)
}

// GetActiveRentalHUD handles GET /api/v1/bookings/{id}/hud
func (s *Server) GetActiveRentalHUD(w http.ResponseWriter, r *http.Request) {
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

	hud, err := s.rentalService.GetActiveRentalHUD(r.Context(), bookingID, userUUID)
	if err != nil {
		if errors.Is(err, bookingDomain.ErrUnauthorizedAction) {
			response.Error(w, r, http.StatusForbidden, response.ProblemDetails{
				Title:  "Forbidden",
				Code:   "ERR_FORBIDDEN",
				Detail: err.Error(),
			})
			return
		}
		response.Error(w, r, http.StatusNotFound, response.ProblemDetails{
			Title:  "HUD Not Found",
			Code:   "ERR_NOT_FOUND",
			Detail: err.Error(),
		})
		return
	}

	data := generated.ActiveRentalHUDResponseData{
		BookingId: hud.BookingID,
		Fleet: struct {
			Id          openapi_types.UUID `json:"id"`
			PlateNumber string             `json:"plateNumber"`
			Title       string             `json:"title"`
		}{
			Id:          hud.FleetID,
			PlateNumber: hud.PlateNumber,
			Title:       hud.FleetTitle,
		},
		Status: hud.Status,
	}

	if hud.OperatorID != nil {
		data.Operator = &struct {
			FullName     *string             `json:"fullName,omitempty"`
			Id           *openapi_types.UUID `json:"id,omitempty"`
			LicenseClass *string             `json:"licenseClass,omitempty"`
			PhoneNumber  *string             `json:"phoneNumber,omitempty"`
		}{
			FullName:     &hud.OperatorFullName,
			Id:           hud.OperatorID,
			LicenseClass: &hud.OperatorLicenseClass,
			PhoneNumber:  &hud.OperatorPhoneNumber,
		}
	}

	data.QuickActionsEnabled.CanExtend = hud.CanExtend
	data.QuickActionsEnabled.CanInitiateReturn = hud.CanInitiateReturn
	data.QuickActionsEnabled.CanReportBreakdown = hud.CanReportBreakdown

	data.ShiftTelemetry.AccumulatedHoursToday = hud.AccumulatedHoursToday
	data.ShiftTelemetry.AllowedHoursPerDay = hud.AllowedHoursPerDay
	data.ShiftTelemetry.CurrentRentalDay = hud.CurrentRentalDay
	data.ShiftTelemetry.ExcessHoursTotal = hud.ExcessHoursTotal
	data.ShiftTelemetry.HandoverBaselineHours = hud.HandoverBaselineHours
	data.ShiftTelemetry.ShiftTimerPaused = hud.ShiftTimerPaused
	data.ShiftTelemetry.TotalRentalDays = hud.TotalRentalDays

	response.JSON(w, r, http.StatusOK, data)
}

// RequestRentalExtension handles POST /api/v1/bookings/{id}/extend
func (s *Server) RequestRentalExtension(w http.ResponseWriter, r *http.Request) {
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

	var req generated.RentalExtensionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Bad Request",
			Code:   "ERR_INVALID_BODY",
			Detail: "Malformed request payload.",
		})
		return
	}

	ext, err := s.rentalService.RequestExtension(r.Context(), rentalApp.RentalExtensionCommand{
		BookingID:  bookingID,
		RenterID:   userUUID,
		NewEndDate: req.NewEndDate.Time,
	})
	if err != nil {
		if errors.Is(err, rentalDomain.ErrExtensionDateUnavailable) {
			response.Error(w, r, http.StatusConflict, response.ProblemDetails{
				Title:  "Extension Dates Conflict",
				Code:   "ERR_EXTENSION_CONFLICT",
				Detail: err.Error(),
			})
			return
		}
		if errors.Is(err, rentalDomain.ErrInvalidExtensionDate) {
			response.Error(w, r, http.StatusUnprocessableEntity, response.ProblemDetails{
				Title:  "Invalid Extension Date",
				Code:   "ERR_INVALID_EXTENSION_DATE",
				Detail: err.Error(),
			})
			return
		}
		if errors.Is(err, bookingDomain.ErrUnauthorizedAction) {
			response.Error(w, r, http.StatusForbidden, response.ProblemDetails{
				Title:  "Forbidden",
				Code:   "ERR_FORBIDDEN",
				Detail: "Only renter can request an extension.",
			})
			return
		}
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Extension Request Failed",
			Code:   "ERR_EXTENSION_FAILED",
			Detail: err.Error(),
		})
		return
	}

	data := generated.RentalExtensionResponseData{
		AdditionalGrossMinor:               ext.IncrementalGrossMinor,
		BookingId:                          ext.BookingID,
		Currency:                           ext.Currency,
		ExtensionDays:                      ext.ExtensionDays,
		IncrementalPrepaymentRequiredMinor: ext.IncrementalPrepaymentMinor,
		OwnerApprovalRequired:              true,
		Status:                             ext.Status,
	}

	response.JSON(w, r, http.StatusOK, data)
}

// ReportMechanicalBreakdown handles POST /api/v1/bookings/{id}/breakdown
func (s *Server) ReportMechanicalBreakdown(w http.ResponseWriter, r *http.Request) {
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

	var req generated.MechanicalBreakdownRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Bad Request",
			Code:   "ERR_INVALID_BODY",
			Detail: "Malformed request payload.",
		})
		return
	}

	dtcs := []string{}
	if req.DiagnosticTroubleCodes != nil {
		dtcs = *req.DiagnosticTroubleCodes
	}

	photos := []string{}
	if req.IncidentPhotos != nil {
		photos = *req.IncidentPhotos
	}

	breakdown, err := s.rentalService.ReportMechanicalBreakdown(r.Context(), rentalApp.BreakdownReportCommand{
		BookingID:              bookingID,
		ReporterID:             userUUID,
		FaultCategory:          string(req.FaultCategory),
		FaultDescription:       req.FaultDescription,
		DiagnosticTroubleCodes: dtcs,
		CurrentHourMeter:       req.CurrentHourMeter,
		IncidentPhotos:         photos,
	})
	if err != nil {
		if errors.Is(err, rentalDomain.ErrBookingAlreadyPaused) {
			response.Error(w, r, http.StatusConflict, response.ProblemDetails{
				Title:  "Booking Already Paused",
				Code:   "ERR_BOOKING_ALREADY_PAUSED",
				Detail: err.Error(),
			})
			return
		}
		if errors.Is(err, rentalDomain.ErrInvalidFaultCategory) {
			response.Error(w, r, http.StatusUnprocessableEntity, response.ProblemDetails{
				Title:  "Invalid Fault Category",
				Code:   "ERR_INVALID_FAULT_CATEGORY",
				Detail: err.Error(),
			})
			return
		}
		if errors.Is(err, bookingDomain.ErrUnauthorizedAction) {
			response.Error(w, r, http.StatusForbidden, response.ProblemDetails{
				Title:  "Forbidden",
				Code:   "ERR_FORBIDDEN",
				Detail: err.Error(),
			})
			return
		}
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Breakdown Submission Failed",
			Code:   "ERR_BREAKDOWN_FAILED",
			Detail: err.Error(),
		})
		return
	}

	data := generated.MechanicalBreakdownResponseData{
		BookingId:          breakdown.BookingID,
		IncidentId:         breakdown.ID,
		OwnerSlaDeadline:   breakdown.SLADeadline,
		ShiftTimerPausedAt: breakdown.CreatedAt,
		Status:             breakdown.Status,
	}

	response.JSON(w, r, http.StatusOK, data)
}

// ResolveMechanicalBreakdown handles POST /api/v1/bookings/{id}/breakdown/resolve
func (s *Server) ResolveMechanicalBreakdown(w http.ResponseWriter, r *http.Request) {
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

	var req generated.BreakdownResolutionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Bad Request",
			Code:   "ERR_INVALID_BODY",
			Detail: "Malformed request payload.",
		})
		return
	}

	breakdown, err := s.rentalService.ResolveMechanicalBreakdown(r.Context(), rentalApp.BreakdownResolveCommand{
		BookingID:       bookingID,
		ResolverID:      userUUID,
		MechanicName:    req.MechanicName,
		MechanicPhone:   req.MechanicPhone,
		ResolutionNotes: req.ResolutionNotes,
	})
	if err != nil {
		if errors.Is(err, rentalDomain.ErrBreakdownNotFound) {
			response.Error(w, r, http.StatusNotFound, response.ProblemDetails{
				Title:  "Breakdown Not Found",
				Code:   "ERR_BREAKDOWN_NOT_FOUND",
				Detail: err.Error(),
			})
			return
		}
		if errors.Is(err, bookingDomain.ErrUnauthorizedAction) {
			response.Error(w, r, http.StatusForbidden, response.ProblemDetails{
				Title:  "Forbidden",
				Code:   "ERR_FORBIDDEN",
				Detail: "Only equipment owner can resolve breakdown.",
			})
			return
		}
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Resolution Failed",
			Code:   "ERR_RESOLUTION_FAILED",
			Detail: err.Error(),
		})
		return
	}

	resTime := time.Now().UTC()
	if breakdown.ResolvedAt != nil {
		resTime = *breakdown.ResolvedAt
	}

	data := generated.BreakdownResolutionResponseData{
		BookingId:  breakdown.BookingID,
		IncidentId: breakdown.ID,
		ResumedAt:  resTime,
		Status:     breakdown.Status,
	}

	response.JSON(w, r, http.StatusOK, data)
}

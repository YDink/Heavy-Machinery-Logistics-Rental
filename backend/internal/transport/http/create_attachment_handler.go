package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	machineryApp "github.com/dinkmachinery/dink-core/internal/modules/machinery/application"
	machineryDomain "github.com/dinkmachinery/dink-core/internal/modules/machinery/domain"
	"github.com/dinkmachinery/dink-core/internal/platform/response"
	"github.com/dinkmachinery/dink-core/internal/transport/http/generated"
)

// CreateAttachment handles POST /api/v1/attachments
func (s *Server) CreateAttachment(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		response.Error(w, r, http.StatusUnauthorized, response.ProblemDetails{
			Title:  "Unauthorized",
			Code:   "ERR_UNAUTHORIZED",
			Detail: "Missing or invalid Bearer authorization token.",
		})
		return
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := s.tokenService.VerifyAccessToken(r.Context(), tokenStr)
	if err != nil {
		response.Error(w, r, http.StatusUnauthorized, response.ProblemDetails{
			Title:  "Unauthorized",
			Code:   "ERR_UNAUTHORIZED",
			Detail: "Access token is expired or signature is invalid.",
		})
		return
	}
	userUUID, err := uuid.Parse(claims.Subject)
	if err != nil {
		response.Error(w, r, http.StatusUnauthorized, response.ProblemDetails{
			Title:  "Unauthorized",
			Code:   "ERR_UNAUTHORIZED",
			Detail: "Invalid token subject claim.",
		})
		return
	}

	var body generated.CreateAttachmentJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Bad Request",
			Code:   "ERR_MALFORMED_JSON",
			Detail: "Request body could not be parsed as valid JSON.",
		})
		return
	}

	photos := make([]machineryDomain.AttachmentPhoto, 0)
	if body.Photos != nil {
		for _, p := range *body.Photos {
			photos = append(photos, machineryDomain.AttachmentPhoto{
				Angle:      machineryDomain.AttachmentPhotoAngle(p.Angle),
				StorageURI: p.StorageUri,
			})
		}
	}

	secDeposit := int64(0)
	if body.SecurityDepositMinor != nil {
		secDeposit = *body.SecurityDepositMinor
	}
	minDays := 1
	if body.MinRentalDays != nil {
		minDays = *body.MinRentalDays
	}

	att, err := s.attachmentService.RegisterAttachment(r.Context(), machineryApp.RegisterAttachmentInput{
		OwnerID:                    userUUID,
		AttachmentType:             machineryDomain.AttachmentType(body.AttachmentType),
		Make:                       body.Make,
		Model:                      body.Model,
		SubTypeSize:                body.SubTypeSize,
		SerialNumber:               body.SerialNumber,
		OperatingWeightKg:          float64(body.OperatingWeightKg),
		MinCarrierWeightTonnes:     float64(body.MinCarrierWeightTonnes),
		MaxCarrierWeightTonnes:     float64(body.MaxCarrierWeightTonnes),
		RequiredFlowMinLmin:        float64(body.RequiredFlowMinLmin),
		RequiredFlowMaxLmin:        float64(body.RequiredFlowMaxLmin),
		RequiredPressureBar:        float64(body.RequiredPressureBar),
		PinDiameterMm:              float64(body.PinDiameterMm),
		PinSpreadMm:                float64(body.PinSpreadMm),
		CompatibleSubcategories:    body.CompatibleSubcategories,
		StandaloneDailyRateMinor:   body.StandaloneDailyRateMinor,
		SecurityDepositMinor:       secDeposit,
		MinRentalDays:              minDays,
		PaymentEvidenceWindowHours: 24,
		Photos:                     photos,
	})
	if err != nil {
		if errors.Is(err, machineryDomain.ErrIncompleteAttachmentPhotos) {
			response.Error(w, r, http.StatusUnprocessableEntity, response.ProblemDetails{
				Title:  "Incomplete Attachment Photos",
				Code:   "ERR_INCOMPLETE_PHOTOS",
				Detail: err.Error(),
			})
			return
		}
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Invalid Attachment Registration",
			Code:   "ERR_INVALID_ATTACHMENT",
			Detail: err.Error(),
		})
		return
	}

	data := generated.AttachmentResponseData{
		AttachmentId:             openapi_types.UUID(att.ID),
		Status:                   generated.AttachmentResponseDataStatus(att.Status),
		AttachmentType:           string(att.AttachmentType),
		Make:                     att.Make,
		Model:                    att.Model,
		StandaloneDailyRateMinor: att.StandaloneDailyRateMinor,
		CreatedAt:                att.CreatedAt,
	}

	response.JSON(w, r, http.StatusCreated, data)
}

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

// CreateFleetAsset handles POST /api/v1/fleets
func (s *Server) CreateFleetAsset(w http.ResponseWriter, r *http.Request) {
	if s.machineryService == nil {
		s.notImplemented(w, r, "05-02")
		return
	}

	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		response.Error(w, r, http.StatusUnauthorized, response.ProblemDetails{
			Title:  "Unauthorized",
			Code:   "ERR_UNAUTHORIZED",
			Detail: "Missing or malformed Authorization header.",
		})
		return
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := s.tokenService.VerifyAccessToken(r.Context(), tokenStr)
	if err != nil {
		response.Error(w, r, http.StatusUnauthorized, response.ProblemDetails{
			Title:  "Unauthorized",
			Code:   "ERR_INVALID_TOKEN",
			Detail: "Access token is invalid or expired.",
		})
		return
	}

	userUUID, err := uuid.Parse(claims.Subject)
	if err != nil {
		response.Error(w, r, http.StatusUnauthorized, response.ProblemDetails{
			Title:  "Unauthorized",
			Code:   "ERR_INVALID_SUBJECT",
			Detail: "Subject claim is not a valid UUID.",
		})
		return
	}

	var body struct {
		generated.CreateFleetRequest
		ChassisVin          string  `json:"chassisVin"`
		BaselineMeterHours  float64 `json:"baselineMeterHours"`
		FuelType            string  `json:"fuelType"`
		SelfMoveCapable     bool    `json:"selfMoveCapable"`
		ExcessHourRateMinor int64   `json:"excessHourRateMinor"`
		MinRentalDays       int     `json:"minRentalDays"`
		PaymentTimeoutHours int     `json:"paymentTimeoutHours"`
		OperatorIncluded    bool    `json:"operatorIncluded"`
		Photos              []struct {
			Angle      string `json:"angle"`
			StorageURI string `json:"storageUri"`
			FileHash   string `json:"fileHash"`
		} `json:"photos"`
		SpecValues []struct {
			FieldKey     string   `json:"fieldKey"`
			ValueString  *string  `json:"valueString"`
			ValueNumeric *float64 `json:"valueNumeric"`
			ValueBoolean *bool    `json:"valueBoolean"`
		} `json:"specValues"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Bad Request",
			Code:   "ERR_MALFORMED_JSON",
			Detail: "Request body could not be parsed as valid JSON.",
		})
		return
	}

	var photos []machineryApp.CreateFleetPhotoInput
	for _, p := range body.Photos {
		photos = append(photos, machineryApp.CreateFleetPhotoInput{
			Angle:      machineryDomain.PhotoAngle(p.Angle),
			StorageURI: p.StorageURI,
			FileHash:   p.FileHash,
		})
	}

	var specValues []machineryApp.CreateFleetSpecValueInput
	for _, v := range body.SpecValues {
		specValues = append(specValues, machineryApp.CreateFleetSpecValueInput{
			FieldKey:     v.FieldKey,
			ValueString:  v.ValueString,
			ValueNumeric: v.ValueNumeric,
			ValueBoolean: v.ValueBoolean,
		})
	}

	fuel := machineryDomain.FuelTypeDiesel
	if body.FuelType != "" {
		fuel = machineryDomain.FuelType(body.FuelType)
	}

	input := machineryApp.CreateFleetInput{
		OwnerID:              userUUID,
		SubcategoryID:        uuid.UUID(body.CategoryId),
		Make:                 body.Make,
		Model:                body.Model,
		ManufactureYear:      body.ManufacturingYear,
		PlateNumber:          body.RegistrationNumber,
		ChassisVIN:           body.ChassisVin,
		BaselineMeterHours:   body.BaselineMeterHours,
		DailyRateMinor:       body.DailyRate.AmountMinor,
		FuelType:             fuel,
		SelfMoveCapable:      body.SelfMoveCapable,
		ExcessHourRateMinor:  body.ExcessHourRateMinor,
		MinRentalDays:        body.MinRentalDays,
		PaymentTimeoutHours:  body.PaymentTimeoutHours,
		OperatorIncluded:     body.OperatorIncluded,
		Latitude:             body.Location.Latitude,
		Longitude:            body.Location.Longitude,
		AddressText:          body.Location.AddressText,
		Photos:               photos,
		SpecValues:          specValues,
	}

	out, err := s.machineryService.CreateFleet(r.Context(), input)
	if err != nil {
		if errors.Is(err, machineryDomain.ErrInvalidBaselineMeter) {
			response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
				Title:  "Invalid Baseline Meter",
				Code:   "ERR_INVALID_BASELINE_METER",
				Detail: "Baseline engine meter reading must be non-negative.",
			})
			return
		}
		if errors.Is(err, machineryDomain.ErrIncompletePhotoWalkaround) {
			response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
				Title:  "Incomplete Photo Walkaround",
				Code:   "ERR_INCOMPLETE_PHOTO_WALKAROUND",
				Detail: "All 6 walkaround angles (FRONT, REAR, LEFT, RIGHT, CABIN, METER) are required.",
			})
			return
		}
		if errors.Is(err, machineryDomain.ErrCategoryNotFound) {
			response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
				Title:  "Category Not Found",
				Code:   "ERR_CATEGORY_NOT_FOUND",
				Detail: "The specified machinery category does not exist.",
			})
			return
		}
		response.Error(w, r, http.StatusInternalServerError, response.ProblemDetails{
			Title:  "Internal Server Error",
			Code:   "ERR_INTERNAL",
			Detail: "Failed to create fleet asset.",
		})
		return
	}

	createdAsset, err := s.machineryService.GetFleetAsset(r.Context(), out.FleetID)
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, response.ProblemDetails{
			Title:  "Internal Server Error",
			Code:   "ERR_INTERNAL",
			Detail: "Failed to fetch created fleet asset.",
		})
		return
	}

	dto := generated.FleetAsset{
		Id:                 openapi_types.UUID(createdAsset.ID),
		OwnerId:            openapi_types.UUID(createdAsset.OwnerID),
		CategoryId:         openapi_types.UUID(createdAsset.SubcategoryID),
		Make:               &createdAsset.Make,
		Model:              createdAsset.Model,
		ManufacturingYear:  &createdAsset.ManufactureYear,
		RegistrationNumber: createdAsset.PlateNumber,
		DailyRate: generated.MoneyAmount{
			AmountMinor: createdAsset.DailyRateMinor,
			Currency:    "ETB",
		},
		Status:  generated.FleetAssetStatus(createdAsset.Status),
		Version: createdAsset.Version,
	}
	dto.Location.Latitude = createdAsset.Latitude
	dto.Location.Longitude = createdAsset.Longitude
	dto.Location.AddressText = createdAsset.AddressText

	response.JSON(w, r, http.StatusCreated, dto)
}

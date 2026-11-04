package http

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	machineryDomain "github.com/dinkmachinery/dink-core/internal/modules/machinery/domain"
	"github.com/dinkmachinery/dink-core/internal/platform/response"
	"github.com/dinkmachinery/dink-core/internal/transport/http/generated"
)

// GetFleetById handles GET /api/v1/fleets/{id}
func (s *Server) GetFleetById(w http.ResponseWriter, r *http.Request) {
	if s.machineryService == nil {
		s.notImplemented(w, r, "05-02")
		return
	}

	idStr := r.PathValue("id")
	fleetID, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Invalid Fleet ID",
			Code:   "ERR_INVALID_UUID",
			Detail: "The provided fleet ID is not a valid UUID.",
		})
		return
	}

	asset, err := s.machineryService.GetFleetAsset(r.Context(), fleetID)
	if err != nil {
		if errors.Is(err, machineryDomain.ErrFleetNotFound) {
			response.Error(w, r, http.StatusNotFound, response.ProblemDetails{
				Title:  "Fleet Asset Not Found",
				Code:   "ERR_FLEET_NOT_FOUND",
				Detail: "The requested machinery fleet asset could not be found.",
			})
			return
		}
		response.Error(w, r, http.StatusInternalServerError, response.ProblemDetails{
			Title:  "Internal Server Error",
			Code:   "ERR_INTERNAL",
			Detail: "Failed to retrieve fleet asset.",
		})
		return
	}

	dto := generated.FleetAsset{
		Id:                 openapi_types.UUID(asset.ID),
		OwnerId:            openapi_types.UUID(asset.OwnerID),
		CategoryId:         openapi_types.UUID(asset.SubcategoryID),
		Make:               &asset.Make,
		Model:              asset.Model,
		ManufacturingYear:  &asset.ManufactureYear,
		RegistrationNumber: asset.PlateNumber,
		DailyRate: generated.MoneyAmount{
			AmountMinor: asset.DailyRateMinor,
			Currency:    "ETB",
		},
		Status:  generated.FleetAssetStatus(asset.Status),
		Version: asset.Version,
	}
	dto.Location.Latitude = asset.Latitude
	dto.Location.Longitude = asset.Longitude
	dto.Location.AddressText = asset.AddressText

	response.JSON(w, r, http.StatusOK, dto)
}

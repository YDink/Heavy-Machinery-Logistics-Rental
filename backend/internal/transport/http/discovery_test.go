package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	machineryApp "github.com/dinkmachinery/dink-core/internal/modules/machinery/application"
	"github.com/dinkmachinery/dink-core/internal/modules/machinery/domain"
	machineryInfra "github.com/dinkmachinery/dink-core/internal/modules/machinery/infrastructure"
	transport "github.com/dinkmachinery/dink-core/internal/transport/http"
	"github.com/dinkmachinery/dink-core/internal/transport/http/generated"
)

func setupDiscoveryTestServer(t *testing.T) (http.Handler, *machineryInfra.MemoryMachineryRepository, *machineryApp.MemoryHoldLeaseChecker) {
	memRepo := machineryInfra.NewMemoryMachineryRepository()
	leaseChecker := machineryApp.NewMemoryHoldLeaseChecker()

	discoveryService := machineryApp.NewDiscoveryService(memRepo, memRepo, leaseChecker)
	server := transport.NewServer(nil, nil, nil, nil, discoveryService, nil, nil)

	mux := http.NewServeMux()
	server.RegisterRoutes(mux)

	return mux, memRepo, leaseChecker
}

func TestSearchFleetsEndpoint(t *testing.T) {
	r, repo, _ := setupDiscoveryTestServer(t)
	ctx := context.Background()

	subcatID := uuid.MustParse("019a8f00-0002-7000-8000-000000000001")
	f1, err := repo.CreateFleet(ctx, &domain.FleetAsset{
		OwnerID:             uuid.New(),
		SubcategoryID:       subcatID,
		Make:                "Caterpillar",
		Model:               "320D",
		ManufactureYear:     2024,
		DailyRateMinor:      2250000,
		ExcessHourRateMinor: 280000,
		Status:              domain.FleetStatusActive,
		Latitude:            9.0300,
		Longitude:           38.7400,
		AddressText:         "Bole, Addis Ababa",
		FuelType:            domain.FuelTypeDiesel,
	})
	require.NoError(t, err)

	t.Run("ValidSpatialSearch_Returns200WithResults", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/fleets/search?latitude=9.0300&longitude=38.7400&radiusKm=50", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)

		assert.True(t, resp["success"].(bool))
		data := resp["data"].(map[string]interface{})
		assert.Equal(t, float64(1), data["totalMatches"])
		results := data["results"].([]interface{})
		require.Equal(t, 1, len(results))

		item := results[0].(map[string]interface{})
		assert.Equal(t, f1.ID.String(), item["fleetId"])
		assert.Equal(t, "Caterpillar", item["make"])
		assert.Equal(t, "320D", item["model"])
		assert.Equal(t, "AVAILABLE_NOW", item["availabilityStatus"])
	})

	t.Run("RadiusExceeds500Km_Returns400ProblemDetails", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/fleets/search?latitude=9.0300&longitude=38.7400&radiusKm=600", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)

		assert.False(t, resp["success"].(bool))
		errObj := resp["error"].(map[string]interface{})
		assert.Equal(t, "ERR_SPATIAL_RADIUS_OUT_OF_BOUNDS", errObj["code"])
	})

	t.Run("LegacyRadiusSearch_ForwardsSuccessfully", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/geo/fleets/radius-search?latitude=9.0300&longitude=38.7400&radiusKm=50", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.True(t, resp["success"].(bool))
	})
}

func TestCompareFleetsEndpoint(t *testing.T) {
	r, repo, _ := setupDiscoveryTestServer(t)
	ctx := context.Background()

	subcat1 := uuid.MustParse("019a8f00-0002-7000-8000-000000000001")
	subcat2 := uuid.New()

	f1, _ := repo.CreateFleet(ctx, &domain.FleetAsset{
		OwnerID:        uuid.New(),
		SubcategoryID:  subcat1,
		Make:           "Caterpillar",
		Model:          "320D",
		DailyRateMinor: 2250000,
		Status:         domain.FleetStatusActive,
	})
	f2, _ := repo.CreateFleet(ctx, &domain.FleetAsset{
		OwnerID:        uuid.New(),
		SubcategoryID:  subcat1,
		Make:           "Komatsu",
		Model:          "PC200-8",
		DailyRateMinor: 2100000,
		Status:         domain.FleetStatusActive,
	})
	fDiff, _ := repo.CreateFleet(ctx, &domain.FleetAsset{
		OwnerID:        uuid.New(),
		SubcategoryID:  subcat2,
		Make:           "Lowboy",
		Model:          "Tri-Axle",
		DailyRateMinor: 3500000,
		Status:         domain.FleetStatusActive,
	})

	t.Run("ValidHomogeneousComparison_Returns200WithMatrix", func(t *testing.T) {
		reqBody := generated.FleetCompareRequest{
			FleetIds: []openapi_types.UUID{openapi_types.UUID(f1.ID), openapi_types.UUID(f2.ID)},
		}
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/fleets/compare", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)

		assert.True(t, resp["success"].(bool))
		data := resp["data"].(map[string]interface{})
		assert.Equal(t, subcat1.String(), data["subcategoryId"])
		fleets := data["fleets"].([]interface{})
		assert.Equal(t, 2, len(fleets))
	})

	t.Run("CrossCategoryComparison_Returns422RFC9457", func(t *testing.T) {
		reqBody := generated.FleetCompareRequest{
			FleetIds: []openapi_types.UUID{openapi_types.UUID(f1.ID), openapi_types.UUID(fDiff.ID)},
		}
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/fleets/compare", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)

		assert.False(t, resp["success"].(bool))
		errObj := resp["error"].(map[string]interface{})
		assert.Equal(t, "ERR_NON_COMPARABLE_CATEGORIES", errObj["code"])
	})

	t.Run("SingleFleetComparison_Returns422RFC9457", func(t *testing.T) {
		reqBody := generated.FleetCompareRequest{
			FleetIds: []openapi_types.UUID{openapi_types.UUID(f1.ID)},
		}
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/fleets/compare", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.False(t, resp["success"].(bool))
		errObj := resp["error"].(map[string]interface{})
		assert.Equal(t, "ERR_INVALID_COMPARISON_COUNT", errObj["code"])
	})
}

func TestGetFleetAvailabilityEndpoint(t *testing.T) {
	r, repo, leaseChecker := setupDiscoveryTestServer(t)
	ctx := context.Background()

	f, _ := repo.CreateFleet(ctx, &domain.FleetAsset{
		OwnerID:       uuid.New(),
		SubcategoryID: uuid.MustParse("019a8f00-0002-7000-8000-000000000001"),
		Make:          "CAT",
		Model:         "320D",
		Status:        domain.FleetStatusActive,
	})

	t.Run("AvailableAsset_Returns200WithIsAvailableTrue", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/fleets/"+f.ID.String()+"/availability?startDate=2026-11-01&endDate=2026-11-10", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)

		assert.True(t, resp["success"].(bool))
		data := resp["data"].(map[string]interface{})
		assert.True(t, data["isAvailable"].(bool))
	})

	t.Run("AssetWithHoldLease_Returns200WithHeldStatus", func(t *testing.T) {
		leaseChecker.SetHold(f.ID, 45*time.Minute)
		defer leaseChecker.ReleaseHold(f.ID)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/fleets/"+f.ID.String()+"/availability?startDate=2026-11-01&endDate=2026-11-10", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)

		data := resp["data"].(map[string]interface{})
		assert.False(t, data["isAvailable"].(bool))
		assert.Equal(t, "HELD_BY_ANOTHER_RENTER", data["unavailableReason"])
		assert.Greater(t, data["holdExpiresInSeconds"].(float64), float64(0))
	})

	t.Run("NonExistentFleet_Returns404NotFound", func(t *testing.T) {
		randomID := uuid.New().String()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/fleets/"+randomID+"/availability?startDate=2026-11-01&endDate=2026-11-10", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

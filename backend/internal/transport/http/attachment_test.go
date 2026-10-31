package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dinkmachinery/dink-core/internal/modules/auth/application"
	authDomain "github.com/dinkmachinery/dink-core/internal/modules/auth/domain"
	authInfra "github.com/dinkmachinery/dink-core/internal/modules/auth/infrastructure"
	machineryApp "github.com/dinkmachinery/dink-core/internal/modules/machinery/application"
	machineryDomain "github.com/dinkmachinery/dink-core/internal/modules/machinery/domain"
	machineryInfra "github.com/dinkmachinery/dink-core/internal/modules/machinery/infrastructure"
	transport "github.com/dinkmachinery/dink-core/internal/transport/http"
	"github.com/dinkmachinery/dink-core/internal/transport/http/generated"
)

func setupAttachmentTestRouter(t *testing.T) (*chi.Mux, *application.TokenService, *machineryInfra.MemoryMachineryRepository, uuid.UUID, string) {
	memRepo := machineryInfra.NewMemoryMachineryRepository()
	cache := authInfra.NewMemoryCache()
	tokenService, err := application.NewTokenService(nil, cache)
	require.NoError(t, err)

	machineryService := machineryApp.NewMachineryService(
		memRepo,
		memRepo,
		memRepo,
		memRepo,
		memRepo,
	)
	attachmentService := machineryApp.NewAttachmentService(memRepo, memRepo)

	server := transport.NewServer(nil, tokenService, machineryService, attachmentService)

	r := chi.NewRouter()
	generated.HandlerWithOptions(server, generated.ChiServerOptions{
		BaseURL:    "/api/v1",
		BaseRouter: r,
	})

	ownerID := uuid.New()
	phone, _ := authDomain.ParsePhoneNumber("+251911223344", "ET")
	testUser := &authDomain.User{
		ID:            ownerID,
		PhoneNumber:   phone,
		CountryCode:   "ET",
		FullName:      "Almaz Tesfaye",
		PrimaryRole:   authDomain.RoleOwner,
		AccountStatus: authDomain.AccountStatusActive,
		KYCStatus:     authDomain.KYCStatusVerified,
		TokenVersion:  1,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	tokenPair, err := tokenService.MintTokenPair(context.Background(), testUser, "")
	require.NoError(t, err)

	return r, tokenService, memRepo, ownerID, tokenPair.AccessToken
}

func TestHTTP_CreateAttachment_Unauthenticated(t *testing.T) {
	router, _, _, _, _ := setupAttachmentTestRouter(t)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/attachments", bytes.NewBuffer([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotency-Key", uuid.New().String())
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestHTTP_CreateAttachment_IncompletePhotos(t *testing.T) {
	router, _, _, _, token := setupAttachmentTestRouter(t)

	body := map[string]interface{}{
		"attachmentType":           "HYDRAULIC_BREAKER",
		"make":                     "Soosan",
		"model":                    "SB-81",
		"subTypeSize":              "135mm Chisel",
		"serialNumber":             "SN-SB81-001",
		"operatingWeightKg":        1850.0,
		"minCarrierWeightTonnes":   18.0,
		"maxCarrierWeightTonnes":   26.0,
		"requiredFlowMinLmin":      120.0,
		"requiredFlowMaxLmin":      150.0,
		"requiredPressureBar":      180.0,
		"pinDiameterMm":            80.0,
		"pinSpreadMm":              450.0,
		"compatibleSubcategories":  []string{"EXCAVATION"},
		"standaloneDailyRateMinor": 1200000,
		"photos": []map[string]string{
			{"angle": "OVERALL", "storageUri": "s3://media/overall.webp"},
		}, // Missing BRACKET, PINS, TOOL_TIP
	}
	bodyBytes, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/attachments", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Idempotency-Key", uuid.New().String())
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestHTTP_CreateAttachment_Success(t *testing.T) {
	router, _, _, _, token := setupAttachmentTestRouter(t)

	body := map[string]interface{}{
		"attachmentType":           "HYDRAULIC_BREAKER",
		"make":                     "Soosan",
		"model":                    "SB-81",
		"subTypeSize":              "135mm Chisel",
		"serialNumber":             "SN-SB81-001",
		"operatingWeightKg":        1850.0,
		"minCarrierWeightTonnes":   18.0,
		"maxCarrierWeightTonnes":   26.0,
		"requiredFlowMinLmin":      120.0,
		"requiredFlowMaxLmin":      150.0,
		"requiredPressureBar":      180.0,
		"pinDiameterMm":            80.0,
		"pinSpreadMm":              450.0,
		"compatibleSubcategories":  []string{"EXCAVATION"},
		"standaloneDailyRateMinor": 1200000,
		"photos": []map[string]string{
			{"angle": "OVERALL", "storageUri": "s3://media/overall.webp"},
			{"angle": "BRACKET", "storageUri": "s3://media/bracket.webp"},
			{"angle": "PINS", "storageUri": "s3://media/pins.webp"},
			{"angle": "TOOL_TIP", "storageUri": "s3://media/tooltip.webp"},
		},
	}
	bodyBytes, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/attachments", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Idempotency-Key", uuid.New().String())
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.True(t, resp["success"].(bool))
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "Soosan", data["make"])
	assert.Equal(t, "SB-81", data["model"])
	assert.Equal(t, "ACTIVE", data["status"])
}

func TestHTTP_LinkFleetAttachment_RejectsRateExceedingStandalone(t *testing.T) {
	router, _, memRepo, ownerID, token := setupAttachmentTestRouter(t)
	ctx := context.Background()

	fleet := &machineryDomain.FleetAsset{
		ID:             uuid.New(),
		OwnerID:        ownerID,
		DailyRateMinor: 2000000,
		Status:         machineryDomain.FleetStatusActive,
	}
	_, err := memRepo.CreateFleet(ctx, fleet)
	require.NoError(t, err)

	att := &machineryDomain.Attachment{
		ID:                       uuid.New(),
		OwnerID:                  ownerID,
		AttachmentType:           machineryDomain.AttachmentTypeHydraulicBreaker,
		Make:                     "Soosan",
		Model:                    "SB-81",
		StandaloneDailyRateMinor: 1200000, // 12,000 ETB
		Status:                   machineryDomain.AttachmentStatusActive,
	}
	_, err = memRepo.CreateAttachment(ctx, att, nil, nil)
	require.NoError(t, err)

	// Attempt linking with 15,000 ETB (> 12,000 ETB) -> 422 Unprocessable Entity
	body := map[string]interface{}{
		"attachmentId":          att.ID.String(),
		"bundledDailyRateMinor": 1500000,
	}
	bodyBytes, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/fleets/"+fleet.ID.String()+"/attachments/link", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

	var errResp map[string]interface{}
	_ = json.Unmarshal(rec.Body.Bytes(), &errResp)
	errObj := errResp["error"].(map[string]interface{})
	assert.Equal(t, "ERR_BUNDLED_RATE_EXCEEDS_STANDALONE", errObj["code"])
}

func TestHTTP_LinkFleetAttachment_Success(t *testing.T) {
	router, _, memRepo, ownerID, token := setupAttachmentTestRouter(t)
	ctx := context.Background()

	fleet := &machineryDomain.FleetAsset{
		ID:             uuid.New(),
		OwnerID:        ownerID,
		DailyRateMinor: 2000000,
		Status:         machineryDomain.FleetStatusActive,
	}
	_, err := memRepo.CreateFleet(ctx, fleet)
	require.NoError(t, err)

	att := &machineryDomain.Attachment{
		ID:                       uuid.New(),
		OwnerID:                  ownerID,
		AttachmentType:           machineryDomain.AttachmentTypeHydraulicBreaker,
		Make:                     "Soosan",
		Model:                    "SB-81",
		StandaloneDailyRateMinor: 1200000, // 12,000 ETB
		Status:                   machineryDomain.AttachmentStatusActive,
	}
	_, err = memRepo.CreateAttachment(ctx, att, nil, nil)
	require.NoError(t, err)

	// Valid link: 8,000 ETB (< 12,000 ETB) -> 200 OK
	body := map[string]interface{}{
		"attachmentId":          att.ID.String(),
		"bundledDailyRateMinor": 800000,
	}
	bodyBytes, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/fleets/"+fleet.ID.String()+"/attachments/link", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(t, err)
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "LINKED", data["status"])
	assert.Equal(t, float64(800000), data["bundledDailyRateMinor"])
}

func TestHTTP_CheckBundledAvailability_Success(t *testing.T) {
	router, _, memRepo, ownerID, _ := setupAttachmentTestRouter(t)
	ctx := context.Background()

	fleet := &machineryDomain.FleetAsset{
		ID:             uuid.New(),
		OwnerID:        ownerID,
		DailyRateMinor: 2250000,
		Status:         machineryDomain.FleetStatusActive,
	}
	_, err := memRepo.CreateFleet(ctx, fleet)
	require.NoError(t, err)

	att := &machineryDomain.Attachment{
		ID:                       uuid.New(),
		OwnerID:                  ownerID,
		AttachmentType:           machineryDomain.AttachmentTypeHydraulicBreaker,
		Make:                     "Soosan",
		Model:                    "SB-81",
		StandaloneDailyRateMinor: 1200000,
		Status:                   machineryDomain.AttachmentStatusActive,
	}
	_, err = memRepo.CreateAttachment(ctx, att, nil, nil)
	require.NoError(t, err)

	body := map[string]interface{}{
		"fleetId":       fleet.ID.String(),
		"attachmentIds": []string{att.ID.String()},
		"startDate":     "2026-11-05",
		"endDate":       "2026-11-12",
	}
	bodyBytes, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/bookings/check-bundled-availability", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(t, err)
	data := resp["data"].(map[string]interface{})
	assert.True(t, data["fleetAvailable"].(bool))
	assert.True(t, data["allAttachmentsAvailable"].(bool))
}

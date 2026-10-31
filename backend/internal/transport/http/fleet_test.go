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
	"go.uber.org/zap"

	"github.com/dinkmachinery/dink-core/internal/modules/auth/application"
	"github.com/dinkmachinery/dink-core/internal/modules/auth/domain"
	"github.com/dinkmachinery/dink-core/internal/modules/auth/infrastructure"
	machineryApp "github.com/dinkmachinery/dink-core/internal/modules/machinery/application"
	machineryDomain "github.com/dinkmachinery/dink-core/internal/modules/machinery/domain"
	machineryInfra "github.com/dinkmachinery/dink-core/internal/modules/machinery/infrastructure"
	"github.com/dinkmachinery/dink-core/internal/platform/i18n"
	transport "github.com/dinkmachinery/dink-core/internal/transport/http"
	"github.com/dinkmachinery/dink-core/internal/transport/http/generated"
)

func setupTestMachineryRouter(t *testing.T) (http.Handler, *application.TokenService, *domain.User, uuid.UUID) {
	userRepo := newMockUserRepo()
	cache := infrastructure.NewMemoryCache()
	logger := zap.NewNop()
	smsGateway := infrastructure.NewMockSMSGateway(logger)

	otpService := application.NewOTPService(cache, smsGateway, application.OTPServiceConfig{MockOTP: true})
	tokenService, err := application.NewTokenService(nil, cache)
	if err != nil {
		t.Fatalf("failed to create TokenService: %v", err)
	}

	authService := application.NewAuthService(userRepo, otpService, tokenService)

	phone, _ := domain.ParsePhoneNumber("+251911223344", "ET")
	testUser := &domain.User{
		ID:            uuid.New(),
		PhoneNumber:   phone,
		CountryCode:   "ET",
		FullName:      "Abebe Fleet Owner",
		PrimaryRole:   domain.RoleOwner,
		AccountStatus: domain.AccountStatusActive,
		KYCStatus:     domain.KYCStatusVerified,
		TokenVersion:  1,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	_, _ = userRepo.Create(context.Background(), testUser)

	machineryRepo := machineryInfra.NewMemoryMachineryRepository()
	categoryID := uuid.New()
	_ = machineryRepo.SeedCategory(machineryDomain.Category{
		ID:           categoryID,
		NameEN:       "Crawler Excavator",
		NameAM:       "የክራውለር ኤክስካቫተር",
		Code:         "EXC-CRW",
		IsActive:     true,
		DisplayOrder: 1,
	})

	machineryService := machineryApp.NewMachineryService(
		machineryRepo,
		machineryRepo,
		machineryRepo,
		machineryRepo,
		machineryRepo,
	)

	server := transport.NewServer(authService, tokenService, machineryService, nil)

	r := chi.NewRouter()
	r.Use(i18n.LanguageMiddleware)
	generated.HandlerWithOptions(server, generated.ChiServerOptions{
		BaseURL:    "/api/v1",
		BaseRouter: r,
	})

	return r, tokenService, testUser, categoryID
}

func TestHTTP_ListCategories_Localization(t *testing.T) {
	handler, _, _, _ := setupTestMachineryRouter(t)

	// 1. English default
	req := httptest.NewRequest(http.MethodGet, "/api/v1/categories", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var enResp struct {
		Success bool `json:"success"`
		Data    []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
			Code string `json:"code"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &enResp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	foundEN := false
	for _, c := range enResp.Data {
		if c.Code == "EXC-CRW" && c.Name == "Crawler Excavator" {
			foundEN = true
		}
	}
	if !foundEN {
		t.Errorf("expected English category name 'Crawler Excavator'")
	}

	// 2. Amharic localized request
	reqAM := httptest.NewRequest(http.MethodGet, "/api/v1/categories", nil)
	reqAM.Header.Set("Accept-Language", "am-ET,am;q=0.9,en;q=0.8")
	wAM := httptest.NewRecorder()
	handler.ServeHTTP(wAM, reqAM)

	if wAM.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", wAM.Code)
	}

	var amResp struct {
		Success bool `json:"success"`
		Data    []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
			Code string `json:"code"`
		} `json:"data"`
	}
	_ = json.Unmarshal(wAM.Body.Bytes(), &amResp)

	foundAM := false
	for _, c := range amResp.Data {
		if c.Code == "EXC-CRW" && c.Name == "የክራውለር ኤክስካቫተር" {
			foundAM = true
		}
	}
	if !foundAM {
		t.Errorf("expected Amharic category name 'የክራውለር ኤክስካቫተር'")
	}
}

func TestHTTP_CreateFleet_FullLifecycle(t *testing.T) {
	handler, tokenService, user, categoryID := setupTestMachineryRouter(t)
	ctx := context.Background()

	tokenPair, err := tokenService.MintTokenPair(ctx, user, "")
	if err != nil {
		t.Fatalf("failed to mint token pair: %v", err)
	}

	// 1. Valid Fleet creation with all 6 angles and logistics specs
	reqBody := map[string]interface{}{
		"categoryId":         categoryID.String(),
		"make":               "Komatsu",
		"model":              "PC210-10M0",
		"manufacturingYear":  2023,
		"registrationNumber": "ET-3-99881",
		"chassisVin":         "KMTPC210X99881122",
		"baselineMeterHours": 850.5,
		"fuelType":           "DIESEL",
		"selfMoveCapable":    false,
		"dailyRate": map[string]interface{}{
			"amountMinor": 5500000,
			"currency":    "ETB",
		},
		"location": map[string]interface{}{
			"latitude":    8.9806,
			"longitude":   38.7578,
			"addressText": "Akaki Kaliti, Addis Ababa",
		},
		"photos": []map[string]interface{}{
			{"angle": "FRONT", "storageUri": "s3://fleet/front.jpg", "fileHash": "h1"},
			{"angle": "REAR", "storageUri": "s3://fleet/rear.jpg", "fileHash": "h2"},
			{"angle": "LEFT", "storageUri": "s3://fleet/left.jpg", "fileHash": "h3"},
			{"angle": "RIGHT", "storageUri": "s3://fleet/right.jpg", "fileHash": "h4"},
			{"angle": "CABIN", "storageUri": "s3://fleet/cabin.jpg", "fileHash": "h5"},
			{"angle": "METER", "storageUri": "s3://fleet/meter.jpg", "fileHash": "h6"},
		},
		"specValues": []map[string]interface{}{
			{"fieldKey": "operating_weight_kg", "valueNumeric": 21500.0},
		},
	}

	reqBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/fleets", bytes.NewReader(reqBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenPair.AccessToken)
	req.Header.Set("Idempotency-Key", uuid.New().String())
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d. Body: %s", w.Code, w.Body.String())
	}

	var createResp struct {
		Success bool `json:"success"`
		Data    struct {
			ID                 string `json:"id"`
			Status             string `json:"status"`
			RegistrationNumber string `json:"registrationNumber"`
			Model              string `json:"model"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &createResp); err != nil {
		t.Fatalf("failed to decode create response: %v", err)
	}

	if createResp.Data.ID == "" {
		t.Fatalf("expected non-empty created fleet ID")
	}
	if createResp.Data.Status != "PENDING_APPROVAL" {
		t.Errorf("expected initial status PENDING_APPROVAL, got %s", createResp.Data.Status)
	}

	// 2. Query created fleet asset by ID (GET /fleets/{id})
	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/fleets/"+createResp.Data.ID, nil)
	getReq.Header.Set("Authorization", "Bearer "+tokenPair.AccessToken)
	getW := httptest.NewRecorder()

	handler.ServeHTTP(getW, getReq)
	if getW.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for GET /fleets/{id}, got %d. Body: %s", getW.Code, getW.Body.String())
	}

	// 3. Test Failure: Incomplete walkaround photos (missing METER)
	incompleteReqBody := map[string]interface{}{
		"categoryId":         categoryID.String(),
		"make":               "Komatsu",
		"model":              "PC210",
		"manufacturingYear":  2023,
		"registrationNumber": "ET-3-99882",
		"baselineMeterHours": 100.0,
		"dailyRate": map[string]interface{}{
			"amountMinor": 5000000,
			"currency":    "ETB",
		},
		"location": map[string]interface{}{
			"latitude":    8.9806,
			"longitude":   38.7578,
			"addressText": "Akaki Kaliti",
		},
		"photos": []map[string]interface{}{
			{"angle": "FRONT", "storageUri": "s3://fleet/front.jpg"},
			{"angle": "REAR", "storageUri": "s3://fleet/rear.jpg"},
		},
	}
	incompleteBytes, _ := json.Marshal(incompleteReqBody)
	incReq := httptest.NewRequest(http.MethodPost, "/api/v1/fleets", bytes.NewReader(incompleteBytes))
	incReq.Header.Set("Content-Type", "application/json")
	incReq.Header.Set("Authorization", "Bearer "+tokenPair.AccessToken)
	incReq.Header.Set("Idempotency-Key", uuid.New().String())
	incW := httptest.NewRecorder()

	handler.ServeHTTP(incW, incReq)
	if incW.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request for incomplete walkaround, got %d. Body: %s", incW.Code, incW.Body.String())
	}

	// 4. Test Failure: Negative baseline meter reading
	negMeterReqBody := map[string]interface{}{
		"categoryId":         categoryID.String(),
		"make":               "Komatsu",
		"model":              "PC210",
		"manufacturingYear":  2023,
		"registrationNumber": "ET-3-99883",
		"baselineMeterHours": -10.0, // Invalid!
		"dailyRate": map[string]interface{}{
			"amountMinor": 5000000,
			"currency":    "ETB",
		},
		"location": map[string]interface{}{
			"latitude":    8.9806,
			"longitude":   38.7578,
			"addressText": "Akaki Kaliti",
		},
		"photos": []map[string]interface{}{
			{"angle": "FRONT", "storageUri": "s3://fleet/front.jpg"},
			{"angle": "REAR", "storageUri": "s3://fleet/rear.jpg"},
			{"angle": "LEFT", "storageUri": "s3://fleet/left.jpg"},
			{"angle": "RIGHT", "storageUri": "s3://fleet/right.jpg"},
			{"angle": "CABIN", "storageUri": "s3://fleet/cabin.jpg"},
			{"angle": "METER", "storageUri": "s3://fleet/meter.jpg"},
		},
	}
	negBytes, _ := json.Marshal(negMeterReqBody)
	negReq := httptest.NewRequest(http.MethodPost, "/api/v1/fleets", bytes.NewReader(negBytes))
	negReq.Header.Set("Content-Type", "application/json")
	negReq.Header.Set("Authorization", "Bearer "+tokenPair.AccessToken)
	negReq.Header.Set("Idempotency-Key", uuid.New().String())
	negW := httptest.NewRecorder()

	handler.ServeHTTP(negW, negReq)
	if negW.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request for negative meter, got %d. Body: %s", negW.Code, negW.Body.String())
	}
}

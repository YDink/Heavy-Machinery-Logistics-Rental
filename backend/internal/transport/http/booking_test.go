package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dinkmachinery/dink-core/internal/modules/auth/application"
	authDomain "github.com/dinkmachinery/dink-core/internal/modules/auth/domain"
	"github.com/dinkmachinery/dink-core/internal/modules/auth/infrastructure"
	bookingApp "github.com/dinkmachinery/dink-core/internal/modules/booking/application"
	bookingInfra "github.com/dinkmachinery/dink-core/internal/modules/booking/infrastructure"
	machineryApp "github.com/dinkmachinery/dink-core/internal/modules/machinery/application"
	machineryDomain "github.com/dinkmachinery/dink-core/internal/modules/machinery/domain"
	machineryInfra "github.com/dinkmachinery/dink-core/internal/modules/machinery/infrastructure"
	"github.com/dinkmachinery/dink-core/internal/platform/i18n"
	transport "github.com/dinkmachinery/dink-core/internal/transport/http"
	"github.com/google/uuid"
)

type bookingTestEnv struct {
	handler      http.Handler
	tokenService *application.TokenService
	renterUser   *authDomain.User
	renterToken  string
	ownerUser    *authDomain.User
	ownerToken   string
	fleetID      uuid.UUID
}

func setupBookingTestEnv(t *testing.T) *bookingTestEnv {
	userRepo := newMockUserRepo()
	cache := infrastructure.NewMemoryCache()
	smsGateway := infrastructure.NewMockSMSGateway(nil)
	otpService := application.NewOTPService(cache, smsGateway, application.OTPServiceConfig{MockOTP: true})
	tokenService, err := application.NewTokenService(nil, cache)
	if err != nil {
		t.Fatalf("failed to create TokenService: %v", err)
	}
	authService := application.NewAuthService(userRepo, otpService, tokenService)

	// Create Renter
	renterPhone, _ := authDomain.ParsePhoneNumber("+251911442910", "ET")
	renterUser := &authDomain.User{
		ID:            uuid.New(),
		PhoneNumber:   renterPhone,
		CountryCode:   "ET",
		FullName:      "Dawit Haile-Mariam",
		PrimaryRole:   authDomain.RoleRenter,
		AccountStatus: authDomain.AccountStatusActive,
		KYCStatus:     authDomain.KYCStatusVerified,
		TokenVersion:  1,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	_, _ = userRepo.Create(context.Background(), renterUser)
	rTokens, _ := tokenService.MintTokenPair(context.Background(), renterUser, "")

	// Create Owner
	ownerPhone, _ := authDomain.ParsePhoneNumber("+251922334455", "ET")
	ownerUser := &authDomain.User{
		ID:            uuid.New(),
		PhoneNumber:   ownerPhone,
		CountryCode:   "ET",
		FullName:      "W/ro Almaz Tesfaye",
		PrimaryRole:   authDomain.RoleOwner,
		AccountStatus: authDomain.AccountStatusActive,
		KYCStatus:     authDomain.KYCStatusVerified,
		TokenVersion:  1,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	_, _ = userRepo.Create(context.Background(), ownerUser)
	oTokens, _ := tokenService.MintTokenPair(context.Background(), ownerUser, "")

	// Setup Machinery
	machineryRepo := machineryInfra.NewMemoryMachineryRepository()
	machineryService := machineryApp.NewMachineryService(
		machineryRepo,
		machineryRepo,
		machineryRepo,
		machineryRepo,
		machineryRepo,
	)

	catID := uuid.New()
	_ = machineryRepo.SeedCategory(machineryDomain.Category{
		ID:           catID,
		Code:         "EXCAVATOR",
		NameEN:       "Excavators",
		NameAM:       "ኤክስካቫተሮች",
		IconURL:      "https://cdn.dinkmachinery.com/icons/excavator.svg",
		IsActive:     true,
		DisplayOrder: 1,
	})

	out, err := machineryService.CreateFleet(context.Background(), machineryApp.CreateFleetInput{
		OwnerID:             ownerUser.ID,
		SubcategoryID:       catID,
		Make:                "Caterpillar",
		Model:               "320D",
		ManufactureYear:     2022,
		PlateNumber:         "ET-3-12345",
		ChassisVIN:          "CAT320D-2022-VIN",
		BaselineMeterHours:  1200.0,
		DailyRateMinor:      4500000, // 45,000 ETB
		ExcessHourRateMinor: 500000,
		MinRentalDays:       2,
		PaymentTimeoutHours: 24,
		Latitude:            9.01,
		Longitude:           38.75,
		AddressText:         "Bole Lemi Industrial Zone, Addis Ababa",
		Photos: []machineryApp.CreateFleetPhotoInput{
			{Angle: machineryDomain.PhotoAngleFront, StorageURI: "s3://fleet/front.jpg", FileHash: "h1"},
			{Angle: machineryDomain.PhotoAngleRear, StorageURI: "s3://fleet/rear.jpg", FileHash: "h2"},
			{Angle: machineryDomain.PhotoAngleLeft, StorageURI: "s3://fleet/left.jpg", FileHash: "h3"},
			{Angle: machineryDomain.PhotoAngleRight, StorageURI: "s3://fleet/right.jpg", FileHash: "h4"},
			{Angle: machineryDomain.PhotoAngleCabin, StorageURI: "s3://fleet/cabin.jpg", FileHash: "h5"},
			{Angle: machineryDomain.PhotoAngleMeter, StorageURI: "s3://fleet/meter.jpg", FileHash: "h6"},
		},
	})
	if err != nil {
		t.Fatalf("failed to create test fleet: %v", err)
	}

	bookingRepo := bookingInfra.NewMemoryBookingRepository()
	holdLease := bookingInfra.NewMemoryHoldLeaseManager()
	bookingService := bookingApp.NewBookingService(bookingRepo, holdLease, machineryService, nil)

	server := transport.NewServer(authService, tokenService, machineryService, nil, nil, bookingService)
	mux := http.NewServeMux()
	server.RegisterRoutes(mux)
	handler := i18n.LanguageMiddleware(mux)

	return &bookingTestEnv{
		handler:      handler,
		tokenService: tokenService,
		renterUser:   renterUser,
		renterToken:  rTokens.AccessToken,
		ownerUser:    ownerUser,
		ownerToken:   oTokens.AccessToken,
		fleetID:      out.FleetID,
	}
}

func TestHTTP_BookingFullLifecycle(t *testing.T) {
	env := setupBookingTestEnv(t)

	// 1. Unauthorized create returns 401
	unauthReq := httptest.NewRequest(http.MethodPost, "/api/v1/bookings", bytes.NewReader([]byte("{}")))
	unauthW := httptest.NewRecorder()
	env.handler.ServeHTTP(unauthW, unauthReq)
	if unauthW.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized, got %d", unauthW.Code)
	}

	// 2. Owner self-booking returns 400 Bad Request
	selfBody := map[string]interface{}{
		"fleetId":             env.fleetID.String(),
		"startDate":           time.Now().Add(24 * time.Hour).Format("2006-01-02"),
		"endDate":             time.Now().Add(96 * time.Hour).Format("2006-01-02"),
		"siteLocationAddress": "Akaki Kaliti, Addis Ababa",
	}
	selfBytes, _ := json.Marshal(selfBody)
	selfReq := httptest.NewRequest(http.MethodPost, "/api/v1/bookings", bytes.NewReader(selfBytes))
	selfReq.Header.Set("Authorization", "Bearer "+env.ownerToken)
	selfReq.Header.Set("Content-Type", "application/json")
	selfW := httptest.NewRecorder()
	env.handler.ServeHTTP(selfW, selfReq)
	if selfW.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request for self-booking, got %d. Body: %s", selfW.Code, selfW.Body.String())
	}

	// 3. Valid booking creation by renter (POST /api/v1/bookings)
	createBody := map[string]interface{}{
		"fleetId":             env.fleetID.String(),
		"startDate":           time.Now().Add(48 * time.Hour).Format("2006-01-02"),
		"endDate":             time.Now().Add(120 * time.Hour).Format("2006-01-02"), // 3 days
		"siteLocationAddress": "Kilinto Industrial Park Phase 1",
	}
	createBytes, _ := json.Marshal(createBody)
	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/bookings", bytes.NewReader(createBytes))
	createReq.Header.Set("Authorization", "Bearer "+env.renterToken)
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	env.handler.ServeHTTP(createW, createReq)

	if createW.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created for booking, got %d. Body: %s", createW.Code, createW.Body.String())
	}

	var createResp struct {
		Data struct {
			ID                string `json:"id"`
			BookingReference  string `json:"bookingReference"`
			Status            string `json:"status"`
			DurationDays      int    `json:"durationDays"`
			HandoverPin       string `json:"handoverPin"`
			Version           int    `json:"version"`
			GrossRentalAmount struct {
				AmountMinor int64 `json:"amountMinor"`
			} `json:"grossRentalAmount"`
			PrepaymentAmount struct {
				AmountMinor int64 `json:"amountMinor"`
			} `json:"prepaymentAmount"`
		} `json:"data"`
	}
	if err := json.Unmarshal(createW.Body.Bytes(), &createResp); err != nil {
		t.Fatalf("failed to parse create response: %v", err)
	}

	bookingID := createResp.Data.ID
	handoverPin := createResp.Data.HandoverPin
	version := createResp.Data.Version

	if createResp.Data.Status != "PENDING_PAYMENT" {
		t.Errorf("expected initial status PENDING_PAYMENT, got %s", createResp.Data.Status)
	}
	if createResp.Data.DurationDays != 3 {
		t.Errorf("expected duration 3 days, got %d", createResp.Data.DurationDays)
	}
	if createResp.Data.GrossRentalAmount.AmountMinor != 3*4500000 {
		t.Errorf("expected gross rental 13,500,000, got %d", createResp.Data.GrossRentalAmount.AmountMinor)
	}

	// 4. Concurrent duplicate booking returns 409 Conflict
	dupReq := httptest.NewRequest(http.MethodPost, "/api/v1/bookings", bytes.NewReader(createBytes))
	dupReq.Header.Set("Authorization", "Bearer "+env.renterToken)
	dupReq.Header.Set("Content-Type", "application/json")
	dupW := httptest.NewRecorder()
	env.handler.ServeHTTP(dupW, dupReq)
	if dupW.Code != http.StatusConflict {
		t.Fatalf("expected 409 Conflict for concurrent booking on held fleet, got %d. Body: %s", dupW.Code, dupW.Body.String())
	}

	// 5. Get Booking by ID (GET /api/v1/bookings/{id})
	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/bookings/"+bookingID, nil)
	getReq.Header.Set("Authorization", "Bearer "+env.renterToken)
	getW := httptest.NewRecorder()
	env.handler.ServeHTTP(getW, getReq)
	if getW.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for GET booking, got %d", getW.Code)
	}

	// 6. Submit Payment Evidence (POST /api/v1/bookings/{id}/payment-evidence)
	slipBody := map[string]interface{}{
		"storageUri":           "s3://payment-vault/cbe_slip.webp",
		"fileHash":             "sha256-hash-slip-verified-12345",
		"transactionReference": "FT2611029831",
		"amountMinor":          createResp.Data.PrepaymentAmount.AmountMinor,
		"currency":             "ETB",
	}
	slipBytes, _ := json.Marshal(slipBody)
	slipReq := httptest.NewRequest(http.MethodPost, "/api/v1/bookings/"+bookingID+"/payment-evidence", bytes.NewReader(slipBytes))
	slipReq.Header.Set("Authorization", "Bearer "+env.renterToken)
	slipReq.Header.Set("Content-Type", "application/json")
	slipW := httptest.NewRecorder()
	env.handler.ServeHTTP(slipW, slipReq)
	if slipW.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created for slip upload, got %d. Body: %s", slipW.Code, slipW.Body.String())
	}

	// 7. Confirm Payment by Owner (POST /api/v1/bookings/{id}:confirm-payment)
	confirmBody := map[string]interface{}{
		"expectedVersion":  version + 1, // Slip upload incremented version
		"paymentReference": "CBE-STATEMENT-MATCH-9921",
		"notes":            "Bank wire receipt verified in full",
	}
	confirmBytes, _ := json.Marshal(confirmBody)
	confirmReq := httptest.NewRequest(http.MethodPost, "/api/v1/bookings/"+bookingID+":confirm-payment", bytes.NewReader(confirmBytes))
	confirmReq.Header.Set("Authorization", "Bearer "+env.ownerToken)
	confirmReq.Header.Set("Content-Type", "application/json")
	confirmW := httptest.NewRecorder()
	env.handler.ServeHTTP(confirmW, confirmReq)
	if confirmW.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for payment confirmation, got %d. Body: %s", confirmW.Code, confirmW.Body.String())
	}

	var confirmedResp struct {
		Data struct {
			Status  string `json:"status"`
			Version int    `json:"version"`
		} `json:"data"`
	}
	_ = json.Unmarshal(confirmW.Body.Bytes(), &confirmedResp)
	if confirmedResp.Data.Status != "CONFIRMED" {
		t.Fatalf("expected status CONFIRMED, got %s", confirmedResp.Data.Status)
	}

	// 8. Handover Check-In via PIN (POST /api/v1/bookings/{id}:handover-check-in)
	// First test bad PIN returns 422
	badPinBody := map[string]interface{}{
		"expectedVersion":     confirmedResp.Data.Version,
		"handoverPin":         "000000",
		"operatingHoursMeter": 1202.0,
	}
	badPinBytes, _ := json.Marshal(badPinBody)
	badPinReq := httptest.NewRequest(http.MethodPost, "/api/v1/bookings/"+bookingID+":handover-check-in", bytes.NewReader(badPinBytes))
	badPinReq.Header.Set("Authorization", "Bearer "+env.renterToken)
	badPinReq.Header.Set("Content-Type", "application/json")
	badPinW := httptest.NewRecorder()
	env.handler.ServeHTTP(badPinW, badPinReq)
	if badPinW.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422 Unprocessable for incorrect PIN, got %d. Body: %s", badPinW.Code, badPinW.Body.String())
	}

	// Valid Handover Check-In
	validHandoverBody := map[string]interface{}{
		"expectedVersion":     confirmedResp.Data.Version,
		"handoverPin":         handoverPin,
		"operatingHoursMeter": 1205.0,
		"photoUris":           []string{"s3://inspections/handover.jpg"},
	}
	validHandoverBytes, _ := json.Marshal(validHandoverBody)
	validHandoverReq := httptest.NewRequest(http.MethodPost, "/api/v1/bookings/"+bookingID+":handover-check-in", bytes.NewReader(validHandoverBytes))
	validHandoverReq.Header.Set("Authorization", "Bearer "+env.renterToken)
	validHandoverReq.Header.Set("Content-Type", "application/json")
	validHandoverW := httptest.NewRecorder()
	env.handler.ServeHTTP(validHandoverW, validHandoverReq)
	if validHandoverW.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for valid handover, got %d. Body: %s", validHandoverW.Code, validHandoverW.Body.String())
	}

	var activeResp struct {
		Data struct {
			Status  string `json:"status"`
			Version int    `json:"version"`
		} `json:"data"`
	}
	_ = json.Unmarshal(validHandoverW.Body.Bytes(), &activeResp)
	if activeResp.Data.Status != "ACTIVE" {
		t.Fatalf("expected status ACTIVE, got %s", activeResp.Data.Status)
	}

	// 9. Return Check-Out (POST /api/v1/bookings/{id}:return-check-out)
	returnBody := map[string]interface{}{
		"expectedVersion":     activeResp.Data.Version,
		"operatingHoursMeter": 1228.5, // > check-in 1205.0
		"damageNotes":         "",
	}
	returnBytes, _ := json.Marshal(returnBody)
	returnReq := httptest.NewRequest(http.MethodPost, "/api/v1/bookings/"+bookingID+":return-check-out", bytes.NewReader(returnBytes))
	returnReq.Header.Set("Authorization", "Bearer "+env.ownerToken)
	returnReq.Header.Set("Content-Type", "application/json")
	returnW := httptest.NewRecorder()
	env.handler.ServeHTTP(returnW, returnReq)
	if returnW.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for return check-out, got %d. Body: %s", returnW.Code, returnW.Body.String())
	}

	var completedResp struct {
		Data struct {
			Status string `json:"status"`
		} `json:"data"`
	}
	_ = json.Unmarshal(returnW.Body.Bytes(), &completedResp)
	if completedResp.Data.Status != "COMPLETED" {
		t.Fatalf("expected status COMPLETED, got %s", completedResp.Data.Status)
	}
}

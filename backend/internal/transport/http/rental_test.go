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
	bookingDomain "github.com/dinkmachinery/dink-core/internal/modules/booking/domain"
	bookingInfra "github.com/dinkmachinery/dink-core/internal/modules/booking/infrastructure"
	rentalApp "github.com/dinkmachinery/dink-core/internal/modules/rental/application"
	rentalInfra "github.com/dinkmachinery/dink-core/internal/modules/rental/infrastructure"
	"github.com/dinkmachinery/dink-core/internal/platform/i18n"
	transport "github.com/dinkmachinery/dink-core/internal/transport/http"
	"github.com/google/uuid"
)

type rentalTestEnv struct {
	handler     http.Handler
	ownerUser   *authDomain.User
	ownerToken  string
	renterUser  *authDomain.User
	renterToken string
	fleetID     uuid.UUID
	booking     *bookingDomain.Booking
	rentalRepo  *rentalInfra.MemoryRentalRepository
	bookingRepo *bookingInfra.MemoryBookingRepository
}

func setupRentalTestEnv(t *testing.T) *rentalTestEnv {
	userRepo := newMockUserRepo()
	cache := infrastructure.NewMemoryCache()
	smsGateway := infrastructure.NewMockSMSGateway(nil)
	otpService := application.NewOTPService(cache, smsGateway, application.OTPServiceConfig{MockOTP: true})
	tokenService, err := application.NewTokenService(nil, cache)
	if err != nil {
		t.Fatalf("failed to create TokenService: %v", err)
	}
	authService := application.NewAuthService(userRepo, otpService, tokenService)

	// Owner
	ownerPhone, _ := authDomain.ParsePhoneNumber("+251922334455", "ET")
	ownerUser := &authDomain.User{
		ID:            uuid.New(),
		PhoneNumber:   ownerPhone,
		CountryCode:   "ET",
		FullName:      "Owner Test",
		PrimaryRole:   authDomain.RoleOwner,
		AccountStatus: authDomain.AccountStatusActive,
		KYCStatus:     authDomain.KYCStatusVerified,
		TokenVersion:  1,
	}
	_, _ = userRepo.Create(context.Background(), ownerUser)
	oTokens, _ := tokenService.MintTokenPair(context.Background(), ownerUser, "")

	// Renter
	renterPhone, _ := authDomain.ParsePhoneNumber("+251911442910", "ET")
	renterUser := &authDomain.User{
		ID:            uuid.New(),
		PhoneNumber:   renterPhone,
		CountryCode:   "ET",
		FullName:      "Renter Test",
		PrimaryRole:   authDomain.RoleRenter,
		AccountStatus: authDomain.AccountStatusActive,
		KYCStatus:     authDomain.KYCStatusVerified,
		TokenVersion:  1,
	}
	_, _ = userRepo.Create(context.Background(), renterUser)
	rTokens, _ := tokenService.MintTokenPair(context.Background(), renterUser, "")

	fleetID := uuid.New()
	rentalRepo := rentalInfra.NewMemoryRentalRepository()
	rentalRepo.SetFleetBaseline(fleetID, 1000.0)

	bookingRepo := bookingInfra.NewMemoryBookingRepository()
	now := time.Now().UTC()
	booking := &bookingDomain.Booking{
		ID:                     uuid.New(),
		BookingReference:       "DNK-2026-RENTAL-TEST",
		OwnerID:                ownerUser.ID,
		RenterID:               renterUser.ID,
		FleetID:                fleetID,
		StartDate:              now,
		EndDate:                now.Add(72 * time.Hour),
		DurationDays:           3,
		DailyRateAppliedMinor:  3000000,
		GrossRentalAmountMinor: 9000000,
		PrepaymentAmountMinor:  4500000,
		Currency:               "ETB",
		Status:                 bookingDomain.StatusConfirmed,
		Version:                1,
	}
	_ = bookingRepo.CreateBooking(context.Background(), booking)

	rentalService := rentalApp.NewRentalService(rentalRepo, bookingRepo)

	server := transport.NewServer(authService, tokenService, nil, nil, nil, nil, rentalService)
	mux := http.NewServeMux()
	server.RegisterRoutes(mux)
	handler := i18n.LanguageMiddleware(mux)

	return &rentalTestEnv{
		handler:     handler,
		ownerUser:   ownerUser,
		ownerToken:  oTokens.AccessToken,
		renterUser:  renterUser,
		renterToken: rTokens.AccessToken,
		fleetID:     fleetID,
		booking:     booking,
		rentalRepo:  rentalRepo,
		bookingRepo: bookingRepo,
	}
}

func TestRentalHTTP_HandoverInspect_SuccessAndRegression(t *testing.T) {
	env := setupRentalTestEnv(t)

	// Case 1: Meter regression below baseline 1000.0 -> HTTP 422
	payloadRegress := map[string]interface{}{
		"checkInMeterHours": 950.0,
		"fluidLevelsOk":     true,
		"safetyGlassIntact": true,
		"inspectionPhotos": []map[string]string{
			{"angle": "METER_CLOSEUP", "storageUri": "s3://inspections/meter.jpg"},
		},
	}
	body, _ := json.Marshal(payloadRegress)
	req := httptest.NewRequest("POST", "/api/v1/bookings/"+env.booking.ID.String()+"/handover-inspect", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+env.ownerToken)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	env.handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422 for meter regression, got %d: %s", w.Code, w.Body.String())
	}

	// Case 2: Valid inspection -> HTTP 200
	payloadValid := map[string]interface{}{
		"checkInMeterHours": 1005.5,
		"fluidLevelsOk":     true,
		"safetyGlassIntact": true,
		"notes":             "Ready for mobilization",
		"inspectionPhotos": []map[string]string{
			{"angle": "METER_CLOSEUP", "storageUri": "s3://inspections/meter.jpg"},
			{"angle": "FRONT_ELEVATION", "storageUri": "s3://inspections/front.jpg"},
		},
	}
	body, _ = json.Marshal(payloadValid)
	req = httptest.NewRequest("POST", "/api/v1/bookings/"+env.booking.ID.String()+"/handover-inspect", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+env.ownerToken)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	env.handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for valid inspection, got %d: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Success bool `json:"success"`
		Data    struct {
			MeterVerified       bool   `json:"meterVerified"`
			PhotosIngestedCount int    `json:"photosIngestedCount"`
			Status              string `json:"status"`
		} `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if !resp.Success || !resp.Data.MeterVerified || resp.Data.PhotosIngestedCount != 2 {
		t.Errorf("unexpected inspection response: %+v", resp)
	}
}

func TestRentalHTTP_HandoverChallenge_GenerateAndConfirm(t *testing.T) {
	env := setupRentalTestEnv(t)

	// Step 1: Owner generates challenge PIN -> HTTP 200
	req := httptest.NewRequest("POST", "/api/v1/bookings/"+env.booking.ID.String()+"/handover-challenge", nil)
	req.Header.Set("Authorization", "Bearer "+env.ownerToken)
	w := httptest.NewRecorder()
	env.handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 generating challenge, got %d: %s", w.Code, w.Body.String())
	}

	var challengeResp struct {
		Success bool `json:"success"`
		Data    struct {
			ChallengePin string `json:"challengePin"`
			MaxAttempts  int    `json:"maxAttempts"`
		} `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &challengeResp)
	if len(challengeResp.Data.ChallengePin) != 6 {
		t.Fatalf("expected 6-digit challenge PIN, got %s", challengeResp.Data.ChallengePin)
	}

	// Step 2: Renter confirms with wrong PIN -> HTTP 422
	confirmWrong := map[string]interface{}{
		"challengePin":         "999999",
		"touchSignatureBase64": "data:image/png;base64,sig...",
	}
	b, _ := json.Marshal(confirmWrong)
	req = httptest.NewRequest("POST", "/api/v1/bookings/"+env.booking.ID.String()+"/handover-confirm", bytes.NewReader(b))
	req.Header.Set("Authorization", "Bearer "+env.renterToken)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	env.handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422 for wrong PIN, got %d: %s", w.Code, w.Body.String())
	}

	// Step 3: Renter confirms with correct PIN -> HTTP 200
	sig := "data:image/png;base64,sig..."
	confirmCorrect := map[string]interface{}{
		"challengePin":         challengeResp.Data.ChallengePin,
		"touchSignatureBase64": sig,
		"contractorLocation": map[string]float64{
			"latitude":  9.03,
			"longitude": 38.74,
		},
	}
	b, _ = json.Marshal(confirmCorrect)
	req = httptest.NewRequest("POST", "/api/v1/bookings/"+env.booking.ID.String()+"/handover-confirm", bytes.NewReader(b))
	req.Header.Set("Authorization", "Bearer "+env.renterToken)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	env.handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for correct confirmation, got %d: %s", w.Code, w.Body.String())
	}

	var confirmResp struct {
		Success bool `json:"success"`
		Data    struct {
			Status                   string  `json:"status"`
			AllocatedShiftHoursTotal float64 `json:"allocatedShiftHoursTotal"`
		} `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &confirmResp)
	if confirmResp.Data.Status != "ACTIVE" || confirmResp.Data.AllocatedShiftHoursTotal != 8.0 {
		t.Errorf("unexpected confirm response: %+v", confirmResp)
	}
}

func TestRentalHTTP_GetActiveRentalHUD(t *testing.T) {
	env := setupRentalTestEnv(t)

	// First activate booking
	reqC := httptest.NewRequest("POST", "/api/v1/bookings/"+env.booking.ID.String()+"/handover-challenge", nil)
	reqC.Header.Set("Authorization", "Bearer "+env.ownerToken)
	wC := httptest.NewRecorder()
	env.handler.ServeHTTP(wC, reqC)

	var cResp struct {
		Data struct {
			ChallengePin string `json:"challengePin"`
		} `json:"data"`
	}
	_ = json.Unmarshal(wC.Body.Bytes(), &cResp)

	sig := "data:image/png;base64,sig..."
	confirmReq := map[string]interface{}{
		"challengePin":         cResp.Data.ChallengePin,
		"touchSignatureBase64": sig,
	}
	b, _ := json.Marshal(confirmReq)
	reqConf := httptest.NewRequest("POST", "/api/v1/bookings/"+env.booking.ID.String()+"/handover-confirm", bytes.NewReader(b))
	reqConf.Header.Set("Authorization", "Bearer "+env.renterToken)
	reqConf.Header.Set("Content-Type", "application/json")
	wConf := httptest.NewRecorder()
	env.handler.ServeHTTP(wConf, reqConf)

	// Fetch HUD
	reqHUD := httptest.NewRequest("GET", "/api/v1/bookings/"+env.booking.ID.String()+"/hud", nil)
	reqHUD.Header.Set("Authorization", "Bearer "+env.renterToken)
	wHUD := httptest.NewRecorder()
	env.handler.ServeHTTP(wHUD, reqHUD)

	if wHUD.Code != http.StatusOK {
		t.Fatalf("expected 200 for HUD, got %d: %s", wHUD.Code, wHUD.Body.String())
	}

	var hudResp struct {
		Success bool `json:"success"`
		Data    struct {
			Status         string `json:"status"`
			ShiftTelemetry struct {
				AllowedHoursPerDay float64 `json:"allowedHoursPerDay"`
			} `json:"shiftTelemetry"`
			QuickActionsEnabled struct {
				CanExtend bool `json:"canExtend"`
			} `json:"quickActionsEnabled"`
		} `json:"data"`
	}
	_ = json.Unmarshal(wHUD.Body.Bytes(), &hudResp)
	if hudResp.Data.Status != "ACTIVE" || hudResp.Data.ShiftTelemetry.AllowedHoursPerDay != 8.0 || !hudResp.Data.QuickActionsEnabled.CanExtend {
		t.Errorf("unexpected HUD response: %+v", hudResp)
	}
}

func TestRentalHTTP_RentalExtension(t *testing.T) {
	env := setupRentalTestEnv(t)

	// Activate booking
	reqC := httptest.NewRequest("POST", "/api/v1/bookings/"+env.booking.ID.String()+"/handover-challenge", nil)
	reqC.Header.Set("Authorization", "Bearer "+env.ownerToken)
	wC := httptest.NewRecorder()
	env.handler.ServeHTTP(wC, reqC)
	var cResp struct {
		Data struct {
			ChallengePin string `json:"challengePin"`
		} `json:"data"`
	}
	_ = json.Unmarshal(wC.Body.Bytes(), &cResp)

	confirmReq := map[string]interface{}{"challengePin": cResp.Data.ChallengePin}
	b, _ := json.Marshal(confirmReq)
	reqConf := httptest.NewRequest("POST", "/api/v1/bookings/"+env.booking.ID.String()+"/handover-confirm", bytes.NewReader(b))
	reqConf.Header.Set("Authorization", "Bearer "+env.renterToken)
	reqConf.Header.Set("Content-Type", "application/json")
	wConf := httptest.NewRecorder()
	env.handler.ServeHTTP(wConf, reqConf)

	// Extension request: 2 days forward
	extReq := map[string]interface{}{
		"newEndDate": env.booking.EndDate.Add(48 * time.Hour).Format("2006-01-02"),
	}
	b, _ = json.Marshal(extReq)
	reqExt := httptest.NewRequest("POST", "/api/v1/bookings/"+env.booking.ID.String()+"/extend", bytes.NewReader(b))
	reqExt.Header.Set("Authorization", "Bearer "+env.renterToken)
	reqExt.Header.Set("Content-Type", "application/json")
	wExt := httptest.NewRecorder()
	env.handler.ServeHTTP(wExt, reqExt)

	if wExt.Code != http.StatusOK {
		t.Fatalf("expected 200 for extension request, got %d: %s", wExt.Code, wExt.Body.String())
	}

	var extResp struct {
		Success bool `json:"success"`
		Data    struct {
			ExtensionDays                      int    `json:"extensionDays"`
			IncrementalPrepaymentRequiredMinor int64  `json:"incrementalPrepaymentRequiredMinor"`
			Status                             string `json:"status"`
		} `json:"data"`
	}
	_ = json.Unmarshal(wExt.Body.Bytes(), &extResp)
	if extResp.Data.ExtensionDays != 2 || extResp.Data.Status != "REQUESTED" {
		t.Errorf("unexpected extension response: %+v", extResp)
	}
}

func TestRentalHTTP_MechanicalBreakdown_ReportAndResolve(t *testing.T) {
	env := setupRentalTestEnv(t)

	// Activate booking
	reqC := httptest.NewRequest("POST", "/api/v1/bookings/"+env.booking.ID.String()+"/handover-challenge", nil)
	reqC.Header.Set("Authorization", "Bearer "+env.ownerToken)
	wC := httptest.NewRecorder()
	env.handler.ServeHTTP(wC, reqC)
	var cResp struct {
		Data struct {
			ChallengePin string `json:"challengePin"`
		} `json:"data"`
	}
	_ = json.Unmarshal(wC.Body.Bytes(), &cResp)

	confirmReq := map[string]interface{}{"challengePin": cResp.Data.ChallengePin}
	b, _ := json.Marshal(confirmReq)
	reqConf := httptest.NewRequest("POST", "/api/v1/bookings/"+env.booking.ID.String()+"/handover-confirm", bytes.NewReader(b))
	reqConf.Header.Set("Authorization", "Bearer "+env.renterToken)
	reqConf.Header.Set("Content-Type", "application/json")
	wConf := httptest.NewRecorder()
	env.handler.ServeHTTP(wConf, reqConf)

	// Step 1: Renter reports breakdown
	bdReq := map[string]interface{}{
		"faultCategory":          "HYDRAULIC_CIRCUIT",
		"faultDescription":       "Severe fluid leak on main boom ram",
		"diagnosticTroubleCodes": []string{"ERR-HYD-01"},
		"currentHourMeter":       1012.0,
		"incidentPhotos":         []string{"s3://photos/leak.jpg"},
	}
	b, _ = json.Marshal(bdReq)
	reqBD := httptest.NewRequest("POST", "/api/v1/bookings/"+env.booking.ID.String()+"/breakdown", bytes.NewReader(b))
	reqBD.Header.Set("Authorization", "Bearer "+env.renterToken)
	reqBD.Header.Set("Content-Type", "application/json")
	wBD := httptest.NewRecorder()
	env.handler.ServeHTTP(wBD, reqBD)

	if wBD.Code != http.StatusOK {
		t.Fatalf("expected 200 reporting breakdown, got %d: %s", wBD.Code, wBD.Body.String())
	}

	var bdResp struct {
		Success bool `json:"success"`
		Data    struct {
			IncidentId string `json:"incidentId"`
			Status     string `json:"status"`
		} `json:"data"`
	}
	_ = json.Unmarshal(wBD.Body.Bytes(), &bdResp)
	if bdResp.Data.Status != "ACTIVE_BREAKDOWN_PAUSED" {
		t.Fatalf("unexpected breakdown response: %+v", bdResp)
	}

	// Step 2: Owner resolves breakdown
	resReq := map[string]interface{}{
		"mechanicName":    "Abebe Kebede",
		"mechanicPhone":   "+251911223344",
		"resolutionNotes": "Replaced O-rings and refilled ISO 46 hydraulic oil",
	}
	b, _ = json.Marshal(resReq)
	reqRes := httptest.NewRequest("POST", "/api/v1/bookings/"+env.booking.ID.String()+"/breakdown/resolve", bytes.NewReader(b))
	reqRes.Header.Set("Authorization", "Bearer "+env.ownerToken)
	reqRes.Header.Set("Content-Type", "application/json")
	wRes := httptest.NewRecorder()
	env.handler.ServeHTTP(wRes, reqRes)

	if wRes.Code != http.StatusOK {
		t.Fatalf("expected 200 resolving breakdown, got %d: %s", wRes.Code, wRes.Body.String())
	}

	var resResp struct {
		Success bool `json:"success"`
		Data    struct {
			Status string `json:"status"`
		} `json:"data"`
	}
	_ = json.Unmarshal(wRes.Body.Bytes(), &resResp)
	if resResp.Data.Status != "RESOLVED" {
		t.Errorf("unexpected resolution response: %+v", resResp)
	}
}

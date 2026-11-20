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
	logisticsApp "github.com/dinkmachinery/dink-core/internal/modules/logistics/application"
	logisticsDomain "github.com/dinkmachinery/dink-core/internal/modules/logistics/domain"
	logisticsInfra "github.com/dinkmachinery/dink-core/internal/modules/logistics/infrastructure"
	"github.com/dinkmachinery/dink-core/internal/platform/i18n"
	transport "github.com/dinkmachinery/dink-core/internal/transport/http"
	"github.com/google/uuid"
)

type logisticsTestEnv struct {
	handler         http.Handler
	ownerToken      string
	renterToken     string
	carrier1Token   string
	carrier2Token   string
	parentBookingID uuid.UUID
	lowboy1FleetID  uuid.UUID
	lowboy2FleetID  uuid.UUID
	undersizedID    uuid.UUID
}

func setupLogisticsHTTPTestEnv(t *testing.T) *logisticsTestEnv {
	t.Helper()
	ctx := context.Background()

	userRepo := newMockUserRepo()
	cache := infrastructure.NewMemoryCache()
	smsGateway := infrastructure.NewMockSMSGateway(nil)
	otpService := application.NewOTPService(cache, smsGateway, application.OTPServiceConfig{MockOTP: true})
	tokenService, err := application.NewTokenService(nil, cache)
	if err != nil {
		t.Fatalf("failed to create TokenService: %v", err)
	}
	authService := application.NewAuthService(userRepo, otpService, tokenService)

	mkUser := func(phoneStr, name string, role authDomain.UserRole) (*authDomain.User, string) {
		p, _ := authDomain.ParsePhoneNumber(phoneStr, "ET")
		u := &authDomain.User{
			ID:            uuid.New(),
			PhoneNumber:   p,
			CountryCode:   "ET",
			FullName:      name,
			PrimaryRole:   role,
			AccountStatus: authDomain.AccountStatusActive,
			KYCStatus:     authDomain.KYCStatusVerified,
			TokenVersion:  1,
		}
		_, _ = userRepo.Create(ctx, u)
		tok, _ := tokenService.MintTokenPair(ctx, u, "")
		return u, tok.AccessToken
	}

	ownerUser, ownerTok := mkUser("+251922334455", "Almaz Owner", authDomain.RoleOwner)
	renterUser, renterTok := mkUser("+251911442910", "Dawit Renter", authDomain.RoleRenter)
	carrier1User, carrier1Tok := mkUser("+251911449120", "Girma Carrier", authDomain.RoleOwner)
	carrier2User, carrier2Tok := mkUser("+251911773300", "Samuel AltCarrier", authDomain.RoleOwner)

	bookingRepo := bookingInfra.NewMemoryBookingRepository()
	logRepo := logisticsInfra.NewMemoryLogisticsRepository(bookingRepo)
	logSvc := logisticsApp.NewLogisticsService(logRepo, bookingRepo)

	machineFleetID := uuid.New()
	lowboy1ID := uuid.New()
	lowboy2ID := uuid.New()
	undersizedID := uuid.New()
	parentID := uuid.New()

	logRepo.SeedFleetProfile(&logisticsDomain.FleetLogisticsProfile{
		FleetID:          machineFleetID,
		OwnerID:          ownerUser.ID,
		Make:             "Caterpillar",
		Model:            "320D",
		RootCategoryCode: "EXCAVATION",
		SelfMoveCapable:  false,
		DailyRateMinor:   3500000,
		Status:           "ACTIVE",
		CountryCode:      "ET",
		Location:         logisticsDomain.GeoPoint{Latitude: 8.9806, Longitude: 38.7578, LocationName: "Bole Depot"},
		OwnerFullName:    ownerUser.FullName,
		OwnerPhoneNumber: ownerUser.PhoneNumber.String(),
		MachineSpecs: logisticsDomain.MachineLogisticsSpecs{
			OperatingWeightKg: 21800,
			TransportLengthMm: 9530,
			TransportWidthMm:  2980,
			TransportHeightMm: 3150,
		},
	})

	logRepo.SeedFleetProfile(&logisticsDomain.FleetLogisticsProfile{
		FleetID:          lowboy1ID,
		OwnerID:          carrier1User.ID,
		Make:             "Sinotruk",
		Model:            "Howo Tri-Axle Lowbed",
		RootCategoryCode: "LOGISTICS",
		SelfMoveCapable:  true,
		DailyRateMinor:   1800000,
		Status:           "ACTIVE",
		CountryCode:      "ET",
		Location:         logisticsDomain.GeoPoint{Latitude: 8.9124, Longitude: 38.7412, LocationName: "Kality Depot"},
		OwnerFullName:    carrier1User.FullName,
		OwnerPhoneNumber: carrier1User.PhoneNumber.String(),
		TrailerProfile: logisticsDomain.MatchedCarrier{
			FleetID:               lowboy1ID,
			OwnerID:               carrier1User.ID,
			Make:                  "Sinotruk",
			Model:                 "Howo Tri-Axle Lowbed",
			OwnerTradingName:      "Girma Heavy Haulage PLC",
			OwnerPhoneNumber:      carrier1User.PhoneNumber.String(),
			PayloadCapacityKg:     34000,
			DeckWellLengthMm:      11000,
			DeckWidthMm:           3100,
			OutriggersExtensionMm: 250,
			DeckHeightGroundMm:    900,
			TareWeightKg:          12500,
			AxleCount:             3,
			DailyRateMinor:        1800000,
			Currency:              "ETB",
			AverageRating:         4.95,
		},
	})

	logRepo.SeedFleetProfile(&logisticsDomain.FleetLogisticsProfile{
		FleetID:          lowboy2ID,
		OwnerID:          carrier2User.ID,
		Make:             "Mercedes-Benz",
		Model:            "Actros Lowbed",
		RootCategoryCode: "LOGISTICS",
		SelfMoveCapable:  true,
		DailyRateMinor:   1950000,
		Status:           "ACTIVE",
		CountryCode:      "ET",
		Location:         logisticsDomain.GeoPoint{Latitude: 8.7900, Longitude: 38.9100, LocationName: "Bishoftu Yard"},
		OwnerFullName:    carrier2User.FullName,
		OwnerPhoneNumber: carrier2User.PhoneNumber.String(),
		TrailerProfile: logisticsDomain.MatchedCarrier{
			FleetID:               lowboy2ID,
			OwnerID:               carrier2User.ID,
			Make:                  "Mercedes-Benz",
			Model:                 "Actros Lowbed",
			OwnerTradingName:      "Corridor Lowboy Logistics",
			OwnerPhoneNumber:      carrier2User.PhoneNumber.String(),
			PayloadCapacityKg:     36000,
			DeckWellLengthMm:      11500,
			DeckWidthMm:           3200,
			OutriggersExtensionMm: 300,
			DeckHeightGroundMm:    900,
			TareWeightKg:          13000,
			AxleCount:             3,
			DailyRateMinor:        1950000,
			Currency:              "ETB",
			AverageRating:         4.88,
		},
	})

	logRepo.SeedFleetProfile(&logisticsDomain.FleetLogisticsProfile{
		FleetID:          undersizedID,
		OwnerID:          carrier2User.ID,
		Make:             "Isuzu",
		Model:            "Light Flatbed",
		RootCategoryCode: "LOGISTICS",
		SelfMoveCapable:  true,
		DailyRateMinor:   900000,
		Status:           "ACTIVE",
		CountryCode:      "ET",
		Location:         logisticsDomain.GeoPoint{Latitude: 8.9100, Longitude: 38.7500, LocationName: "Addis"},
		OwnerFullName:    carrier2User.FullName,
		OwnerPhoneNumber: carrier2User.PhoneNumber.String(),
		TrailerProfile: logisticsDomain.MatchedCarrier{
			FleetID:            undersizedID,
			OwnerID:            carrier2User.ID,
			PayloadCapacityKg:  10000,
			DeckWellLengthMm:   6500,
			DeckWidthMm:        2400,
			DeckHeightGroundMm: 900,
			TareWeightKg:       5000,
			AxleCount:          2,
			DailyRateMinor:     900000,
			Currency:           "ETB",
		},
	})

	now := time.Now().UTC()
	_ = bookingRepo.CreateBooking(ctx, &bookingDomain.Booking{
		ID:                        parentID,
		BookingReference:          "BK-HTTP-LOG-001",
		RenterID:                  renterUser.ID,
		OwnerID:                   ownerUser.ID,
		FleetID:                   machineFleetID,
		BookingKind:               "MACHINERY",
		StartDate:                 now.AddDate(0, 0, 2),
		EndDate:                   now.AddDate(0, 0, 7),
		DurationDays:              5,
		DailyRateAppliedMinor:     3500000,
		BaseHireFeeMinor:          17500000,
		GrossRentalAmountMinor:    17500000,
		PrepaymentAmountMinor:     8750000,
		PlatformCommissionRateBps: 1000,
		PlatformCommissionMinor:   1750000,
		Currency:                  "ETB",
		SiteLocationAddress:       "Mojo Site B",
		Status:                    bookingDomain.StatusConfirmed,
		Version:                   1,
		CreatedAt:                 now,
		UpdatedAt:                 now,
	})

	server := transport.NewServer(authService, tokenService, nil, nil, nil, nil, nil)
	server.SetLogisticsService(logSvc)
	mux := http.NewServeMux()
	server.RegisterRoutes(mux)

	return &logisticsTestEnv{
		handler:         i18n.LanguageMiddleware(mux),
		ownerToken:      ownerTok,
		renterToken:     renterTok,
		carrier1Token:   carrier1Tok,
		carrier2Token:   carrier2Tok,
		parentBookingID: parentID,
		lowboy1FleetID:  lowboy1ID,
		lowboy2FleetID:  lowboy2ID,
		undersizedID:    undersizedID,
	}
}

func TestLogisticsHTTP_EndToEnd(t *testing.T) {
	env := setupLogisticsHTTPTestEnv(t)
	var subBookingID string

	t.Run("POST logistics/match requires auth (401)", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/bookings/"+env.parentBookingID.String()+"/logistics/match", bytes.NewReader([]byte(`{"searchRadiusKm":75,"pickupDate":"2026-11-25T06:00:00Z"}`)))
		w := httptest.NewRecorder()
		env.handler.ServeHTTP(w, req)
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", w.Code)
		}
	})

	t.Run("POST logistics/match returns compatible lowboy carriers", func(t *testing.T) {
		body := `{"searchRadiusKm":75,"pickupDate":"2026-11-25T06:00:00Z"}`
		req := httptest.NewRequest("POST", "/api/v1/bookings/"+env.parentBookingID.String()+"/logistics/match", bytes.NewReader([]byte(body)))
		req.Header.Set("Authorization", "Bearer "+env.ownerToken)
		w := httptest.NewRecorder()
		env.handler.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d body: %s", w.Code, w.Body.String())
		}
		var resp map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		data := resp["data"].(map[string]interface{})
		if int(data["matchingCarriersCount"].(float64)) != 2 {
			t.Fatalf("expected 2 matching carriers, got %v", data["matchingCarriersCount"])
		}
	})

	t.Run("POST logistics/sub-booking rejects undersized trailer with 422", func(t *testing.T) {
		payload := map[string]interface{}{
			"logisticsFleetId": env.undersizedID.String(),
			"pickupDateTime":   "2026-11-25T06:30:00Z",
			"originCoordinates": map[string]interface{}{
				"latitude": 8.9124, "longitude": 38.7412, "locationName": "Kality",
			},
			"destinationCoordinates": map[string]interface{}{
				"latitude": 8.5412, "longitude": 39.2718, "locationName": "Mojo",
			},
		}
		b, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/v1/bookings/"+env.parentBookingID.String()+"/logistics/sub-booking", bytes.NewReader(b))
		req.Header.Set("Authorization", "Bearer "+env.ownerToken)
		w := httptest.NewRecorder()
		env.handler.ServeHTTP(w, req)
		if w.Code != http.StatusUnprocessableEntity {
			t.Fatalf("expected 422, got %d body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("POST logistics/sub-booking creates linked sub-booking with 2.5% commission", func(t *testing.T) {
		payload := map[string]interface{}{
			"logisticsFleetId": env.lowboy1FleetID.String(),
			"pickupDateTime":   "2026-11-25T06:30:00Z",
			"originCoordinates": map[string]interface{}{
				"latitude": 8.9124, "longitude": 38.7412, "locationName": "Kality Depot",
			},
			"destinationCoordinates": map[string]interface{}{
				"latitude": 8.5412, "longitude": 39.2718, "locationName": "Mojo Site B",
			},
		}
		b, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/v1/bookings/"+env.parentBookingID.String()+"/logistics/sub-booking", bytes.NewReader(b))
		req.Header.Set("Authorization", "Bearer "+env.ownerToken)
		w := httptest.NewRecorder()
		env.handler.ServeHTTP(w, req)
		if w.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d body: %s", w.Code, w.Body.String())
		}
		var resp map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		data := resp["data"].(map[string]interface{})
		subBookingID = data["subBookingId"].(string)
		if int(data["platformCommissionBps"].(float64)) != 250 {
			t.Fatalf("expected 250 bps, got %v", data["platformCommissionBps"])
		}
		if int64(data["commissionOwedMinor"].(float64)) != 90000 {
			t.Fatalf("expected 90000 commissionOwedMinor, got %v", data["commissionOwedMinor"])
		}
	})

	t.Run("POST logistics/dispatch advances stages and rejects missing winching photos (422)", func(t *testing.T) {
		// 1. DISPATCHED
		req1 := httptest.NewRequest("POST", "/api/v1/bookings/"+subBookingID+"/logistics/dispatch", bytes.NewReader([]byte(`{"targetStage":"DISPATCHED"}`)))
		req1.Header.Set("Authorization", "Bearer "+env.carrier1Token)
		w1 := httptest.NewRecorder()
		env.handler.ServeHTTP(w1, req1)
		if w1.Code != http.StatusOK {
			t.Fatalf("expected 200 on DISPATCHED, got %d body: %s", w1.Code, w1.Body.String())
		}

		// 2. LOADING_INSPECTION with only 2 photos -> 422
		badLoad := `{"targetStage":"LOADING_INSPECTION","winchingPhotos":[{"angle":"RIG_PROFILE","storageUri":"https://cdn.dinkmachinery.com/1.webp"},{"angle":"CHAIN_BINDERS","storageUri":"https://cdn.dinkmachinery.com/2.webp"}]}`
		reqBad := httptest.NewRequest("POST", "/api/v1/bookings/"+subBookingID+"/logistics/dispatch", bytes.NewReader([]byte(badLoad)))
		reqBad.Header.Set("Authorization", "Bearer "+env.carrier1Token)
		wBad := httptest.NewRecorder()
		env.handler.ServeHTTP(wBad, reqBad)
		if wBad.Code != http.StatusUnprocessableEntity {
			t.Fatalf("expected 422 on incomplete winching photos, got %d body: %s", wBad.Code, wBad.Body.String())
		}

		// 3. LOADING_INSPECTION with all 4 photos -> 200
		goodLoad := `{"targetStage":"LOADING_INSPECTION","winchingPhotos":[{"angle":"RIG_PROFILE","storageUri":"https://cdn.dinkmachinery.com/1.webp"},{"angle":"CHAIN_BINDERS","storageUri":"https://cdn.dinkmachinery.com/2.webp"},{"angle":"BOOM_PINNED","storageUri":"https://cdn.dinkmachinery.com/3.webp"},{"angle":"RAMPS_LOCKED","storageUri":"https://cdn.dinkmachinery.com/4.webp"}]}`
		reqGood := httptest.NewRequest("POST", "/api/v1/bookings/"+subBookingID+"/logistics/dispatch", bytes.NewReader([]byte(goodLoad)))
		reqGood.Header.Set("Authorization", "Bearer "+env.carrier1Token)
		wGood := httptest.NewRecorder()
		env.handler.ServeHTTP(wGood, reqGood)
		if wGood.Code != http.StatusOK {
			t.Fatalf("expected 200 on LOADING_INSPECTION, got %d body: %s", wGood.Code, wGood.Body.String())
		}

		// 4. IN_TRANSIT -> 200
		reqTransit := httptest.NewRequest("POST", "/api/v1/bookings/"+subBookingID+"/logistics/dispatch", bytes.NewReader([]byte(`{"targetStage":"IN_TRANSIT"}`)))
		reqTransit.Header.Set("Authorization", "Bearer "+env.carrier1Token)
		wTransit := httptest.NewRecorder()
		env.handler.ServeHTTP(wTransit, reqTransit)
		if wTransit.Code != http.StatusOK {
			t.Fatalf("expected 200 on IN_TRANSIT, got %d body: %s", wTransit.Code, wTransit.Body.String())
		}
	})

	t.Run("POST logistics/transshipment-alert & reverse-cascade Force Majeure", func(t *testing.T) {
		alertBody := `{"currentLocation":{"latitude":8.75,"longitude":38.90},"failureCategory":"PRIME_MOVER_ENGINE_BLOWN","failureDescription":"Turbocharger shaft sheared on highway shoulder.","photos":["https://cdn.dinkmachinery.com/err.webp"]}`
		reqAlert := httptest.NewRequest("POST", "/api/v1/bookings/"+subBookingID+"/logistics/transshipment-alert", bytes.NewReader([]byte(alertBody)))
		reqAlert.Header.Set("Authorization", "Bearer "+env.carrier1Token)
		wAlert := httptest.NewRecorder()
		env.handler.ServeHTTP(wAlert, reqAlert)
		if wAlert.Code != http.StatusOK {
			t.Fatalf("expected 200 on transshipment-alert, got %d body: %s", wAlert.Code, wAlert.Body.String())
		}

		// GET overview reflects active transshipment
		reqOv := httptest.NewRequest("GET", "/api/v1/bookings/"+env.parentBookingID.String()+"/logistics", nil)
		reqOv.Header.Set("Authorization", "Bearer "+env.renterToken)
		wOv := httptest.NewRecorder()
		env.handler.ServeHTTP(wOv, reqOv)
		if wOv.Code != http.StatusOK {
			t.Fatalf("expected 200 on GET logistics overview, got %d body: %s", wOv.Code, wOv.Body.String())
		}

		// Trigger Force Majeure Reverse Cascade
		reqCascade := httptest.NewRequest("POST", "/api/v1/bookings/"+env.parentBookingID.String()+"/logistics/reverse-cascade", bytes.NewReader([]byte(`{"reason":"60m transshipment window failed"}`)))
		reqCascade.Header.Set("Authorization", "Bearer "+env.ownerToken)
		wCascade := httptest.NewRecorder()
		env.handler.ServeHTTP(wCascade, reqCascade)
		if wCascade.Code != http.StatusOK {
			t.Fatalf("expected 200 on reverse-cascade, got %d body: %s", wCascade.Code, wCascade.Body.String())
		}

		var resp map[string]interface{}
		_ = json.Unmarshal(wCascade.Body.Bytes(), &resp)
		data := resp["data"].(map[string]interface{})
		if data["parentStatus"] != "CANCELLED" || data["subStatus"] != "CANCELLED_CARRIER_DEFAULT" {
			t.Fatalf("unexpected reverse cascade statuses: %v", data)
		}
		if int64(data["ownerPenaltyLeviedMinor"].(float64)) != 0 {
			t.Fatalf("expected 0 owner penalty, got %v", data["ownerPenaltyLeviedMinor"])
		}
	})
}

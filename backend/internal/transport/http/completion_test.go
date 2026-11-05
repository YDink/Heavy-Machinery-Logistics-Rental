package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	bookingDomain "github.com/dinkmachinery/dink-core/internal/modules/booking/domain"
)

func TestCompletionHTTP_FullReturnLifecycle(t *testing.T) {
	env := setupRentalTestEnv(t)
	ctx := context.Background()

	// Advance booking to ACTIVE state with check-in meter
	checkInMeter := 1250.0
	env.booking.Status = bookingDomain.StatusActive
	env.booking.CheckInMeterHours = &checkInMeter
	env.booking.DailyRateAppliedMinor = 2400000 // 2,400.00 ETB / day (300.00 ETB / hr)
	env.booking.BaseHireFeeMinor = 7200000      // 7,200.00 ETB
	env.booking.AttachmentsFeeMinor = 800000    // 800.00 ETB
	env.booking.PrepaymentAmountMinor = 4000000 // 4,000.00 ETB
	env.booking.DurationDays = 3
	_ = env.bookingRepo.UpdateBookingStatusOCC(ctx, env.booking, env.booking.Version)

	var claimID string

	// 1. Return Initiate
	t.Run("POST return-initiate unauthorized without token", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/bookings/"+env.booking.ID.String()+"/return-initiate", nil)
		w := httptest.NewRecorder()
		env.handler.ServeHTTP(w, req)
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401 Unauthorized, got %d", w.Code)
		}
	})

	t.Run("POST return-initiate success by renter", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/bookings/"+env.booking.ID.String()+"/return-initiate", nil)
		req.Header.Set("Authorization", "Bearer "+env.renterToken)
		w := httptest.NewRecorder()
		env.handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d body: %s", w.Code, w.Body.String())
		}

		var resp map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		data := resp["data"].(map[string]interface{})
		if data["status"] != "RETURN_INITIATED" {
			t.Errorf("expected status RETURN_INITIATED, got %v", data["status"])
		}
		if data["checkInMeterHours"].(float64) != 1250.0 {
			t.Errorf("expected checkInMeterHours 1250.0, got %v", data["checkInMeterHours"])
		}
	})

	// 2. Return Inspect
	t.Run("POST return-inspect rejects regressive meter reading", func(t *testing.T) {
		payload := map[string]interface{}{
			"checkOutMeterHours": 1240.0, // Regresses below 1250.0
			"inspectionPhotos": []map[string]string{
				{"angle": "FRONT_IMPLEMENT", "storageUri": "s3://photos/front.jpg"},
				{"angle": "REAR_COUNTERWEIGHT_ENGINE", "storageUri": "s3://photos/rear.jpg"},
				{"angle": "LEFT_SIDE_TRACKS", "storageUri": "s3://photos/left.jpg"},
				{"angle": "RIGHT_SIDE_TRACKS", "storageUri": "s3://photos/right.jpg"},
				{"angle": "HOUR_METER", "storageUri": "s3://photos/meter.jpg"},
				{"angle": "CABIN_INTERIOR", "storageUri": "s3://photos/cabin.jpg"},
			},
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/v1/bookings/"+env.booking.ID.String()+"/return-inspect", bytes.NewReader(body))
		req.Header.Set("Authorization", "Bearer "+env.ownerToken)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		env.handler.ServeHTTP(w, req)

		if w.Code != http.StatusUnprocessableEntity {
			t.Fatalf("expected 422 Unprocessable Entity, got %d body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("POST return-inspect rejects missing mandatory angles", func(t *testing.T) {
		payload := map[string]interface{}{
			"checkOutMeterHours": 1276.5,
			"inspectionPhotos": []map[string]string{
				{"angle": "FRONT_IMPLEMENT", "storageUri": "s3://photos/front.jpg"},
				// Missing the remaining 5 angles
			},
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/v1/bookings/"+env.booking.ID.String()+"/return-inspect", bytes.NewReader(body))
		req.Header.Set("Authorization", "Bearer "+env.ownerToken)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		env.handler.ServeHTTP(w, req)

		if w.Code != http.StatusUnprocessableEntity {
			t.Fatalf("expected 422 Unprocessable Entity, got %d body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("POST return-inspect succeeds with 6 angles and valid monotonic meter", func(t *testing.T) {
		payload := map[string]interface{}{
			"checkOutMeterHours":  1276.5, // 1276.5 - 1250.0 = 26.5h (> 24.0h allowance -> 2.5h excess)
			"fluidLevelsOk":       true,
			"safetyGlassIntact":   true,
			"accelerometerSteady": true,
			"notes":               "Machine returned in standard operational condition",
			"inspectionPhotos": []map[string]string{
				{"angle": "FRONT_IMPLEMENT", "storageUri": "s3://photos/front.jpg"},
				{"angle": "REAR_COUNTERWEIGHT_ENGINE", "storageUri": "s3://photos/rear.jpg"},
				{"angle": "LEFT_SIDE_TRACKS", "storageUri": "s3://photos/left.jpg"},
				{"angle": "RIGHT_SIDE_TRACKS", "storageUri": "s3://photos/right.jpg"},
				{"angle": "HOUR_METER", "storageUri": "s3://photos/meter.jpg"},
				{"angle": "CABIN_INTERIOR", "storageUri": "s3://photos/cabin.jpg"},
			},
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/v1/bookings/"+env.booking.ID.String()+"/return-inspect", bytes.NewReader(body))
		req.Header.Set("Authorization", "Bearer "+env.ownerToken)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		env.handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d body: %s", w.Code, w.Body.String())
		}

		var resp map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		data := resp["data"].(map[string]interface{})
		if data["status"] != "DAMAGE_PENDING_REVIEW" {
			t.Errorf("expected status DAMAGE_PENDING_REVIEW, got %v", data["status"])
		}
		if data["photosIngestedCount"].(float64) != 6 {
			t.Errorf("expected 6 ingested photos, got %v", data["photosIngestedCount"])
		}
	})

	// 3. Lodge Damage Claim
	t.Run("POST damage-claims succeeds when filed by owner", func(t *testing.T) {
		quoteURI := "https://s3.example.com/repair-quote.pdf"
		payload := map[string]interface{}{
			"componentCategory":  "TRACKS_UNDERCARRIAGE",
			"incidentNarrative":  "Deep gouge in left steel track shoe and cracked roller bracket",
			"claimedAmountMinor": 450000, // 4,500.00 ETB
			"currency":           "ETB",
			"repairQuoteUri":     quoteURI,
			"damagePhotos":       []string{"https://s3.example.com/damage1.jpg"},
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/v1/bookings/"+env.booking.ID.String()+"/damage-claims", bytes.NewReader(body))
		req.Header.Set("Authorization", "Bearer "+env.ownerToken)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		env.handler.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("expected 201 Created, got %d body: %s", w.Code, w.Body.String())
		}

		var resp map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		data := resp["data"].(map[string]interface{})
		claimID = data["claimId"].(string)
		if data["status"] != "PENDING_REVIEW" {
			t.Errorf("expected claim status PENDING_REVIEW, got %v", data["status"])
		}
		if int64(data["claimedAmountMinor"].(float64)) != 450000 {
			t.Errorf("expected claimedAmountMinor 450000, got %v", data["claimedAmountMinor"])
		}
	})

	// 4. Respond to Damage Claim
	t.Run("POST damage-claims respond succeeds when contractor accepts", func(t *testing.T) {
		payload := map[string]interface{}{
			"action":              "ACCEPT",
			"acceptedAmountMinor": 450000,
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/v1/bookings/"+env.booking.ID.String()+"/damage-claims/"+claimID+"/respond", bytes.NewReader(body))
		req.Header.Set("Authorization", "Bearer "+env.renterToken)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		env.handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d body: %s", w.Code, w.Body.String())
		}

		var resp map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		data := resp["data"].(map[string]interface{})
		if data["status"] != "ACCEPTED" {
			t.Errorf("expected claim status ACCEPTED, got %v", data["status"])
		}
		if int64(data["acceptedAmountMinor"].(float64)) != 450000 {
			t.Errorf("expected acceptedAmountMinor 450000, got %v", data["acceptedAmountMinor"])
		}
	})

	// 5. Final Reconcile
	t.Run("POST final-reconcile computes universal balance statement", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/bookings/"+env.booking.ID.String()+"/final-reconcile", nil)
		req.Header.Set("Authorization", "Bearer "+env.renterToken)
		w := httptest.NewRecorder()
		env.handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d body: %s", w.Code, w.Body.String())
		}

		var resp map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		data := resp["data"].(map[string]interface{})

		// Base (7,200) + Attachments (800) = 8,000.00 ETB
		// Excess hours: 26.5h - 24.0h = 2.5h * 300.00 = 750.00 ETB (75,000 minor)
		// Accepted damage: 450,000 minor
		// Gross: 800,000 + 75,000 + 450,000 = 1,325,000 minor
		// Prepayment: 400,000 minor
		// Net: 1,325,000 - 400,000 = 925,000 minor > 0 -> PENDING_FINAL_PAYMENT
		if data["status"] != "PENDING_FINAL_PAYMENT" {
			t.Errorf("expected settlement status PENDING_FINAL_PAYMENT, got %v", data["status"])
		}
		if int64(data["excessHoursSurchargeMinor"].(float64)) != 750000 {
			t.Errorf("expected excess surcharge 750000, got %v", data["excessHoursSurchargeMinor"])
		}
		if int64(data["acceptedDamageMinor"].(float64)) != 450000 {
			t.Errorf("expected accepted damage 450000, got %v", data["acceptedDamageMinor"])
		}
		if int64(data["netBalanceMinor"].(float64)) <= 0 {
			t.Errorf("expected positive net balance, got %v", data["netBalanceMinor"])
		}
	})

	// 6. Double-Blind Reviews
	t.Run("POST reviews validates rating boundaries and enforces double-blind unblinding", func(t *testing.T) {
		// Advance booking to COMPLETED so reviews can be posted
		env.booking.Status = bookingDomain.StatusCompleted
		_ = env.bookingRepo.UpdateBookingStatusOCC(ctx, env.booking, env.booking.Version)

		// Invalid rating (< 1.0)
		invalidPayload := map[string]interface{}{
			"rating":   0.5,
			"feedback": "Too low",
		}
		invBody, _ := json.Marshal(invalidPayload)
		reqInv := httptest.NewRequest("POST", "/api/v1/bookings/"+env.booking.ID.String()+"/reviews", bytes.NewReader(invBody))
		reqInv.Header.Set("Authorization", "Bearer "+env.renterToken)
		reqInv.Header.Set("Content-Type", "application/json")
		wInv := httptest.NewRecorder()
		env.handler.ServeHTTP(wInv, reqInv)
		if wInv.Code != http.StatusUnprocessableEntity {
			t.Fatalf("expected 422 for rating < 1.0, got %d", wInv.Code)
		}

		// Renter posts first review -> isPublished must be false
		renterPayload := map[string]interface{}{
			"rating":   4.9,
			"feedback": "Powerful excavator, maintained high productivity throughout job",
			"dimensionRatings": map[string]float64{
				"EQUIPMENT_CONDITION": 5.0,
				"OPERATOR_ACCURACY":   4.8,
			},
		}
		rBody, _ := json.Marshal(renterPayload)
		reqRenter := httptest.NewRequest("POST", "/api/v1/bookings/"+env.booking.ID.String()+"/reviews", bytes.NewReader(rBody))
		reqRenter.Header.Set("Authorization", "Bearer "+env.renterToken)
		reqRenter.Header.Set("Content-Type", "application/json")
		wRenter := httptest.NewRecorder()
		env.handler.ServeHTTP(wRenter, reqRenter)

		if wRenter.Code != http.StatusCreated {
			t.Fatalf("expected 201 Created for renter review, got %d body: %s", wRenter.Code, wRenter.Body.String())
		}
		var rResp map[string]interface{}
		_ = json.Unmarshal(wRenter.Body.Bytes(), &rResp)
		rData := rResp["data"].(map[string]interface{})
		if rData["isPublished"].(bool) {
			t.Errorf("expected first review to remain unpublished (blinded), got isPublished=true")
		}

		// Owner posts second review -> bilateral unblinding (isPublished = true, booking -> RATED)
		ownerPayload := map[string]interface{}{
			"rating":   5.0,
			"feedback": "Flawless site coordination and prompt off-hire payment release",
			"dimensionRatings": map[string]float64{
				"PAYMENT_PUNCTUALITY": 5.0,
				"SITE_SAFETY":         5.0,
			},
		}
		oBody, _ := json.Marshal(ownerPayload)
		reqOwner := httptest.NewRequest("POST", "/api/v1/bookings/"+env.booking.ID.String()+"/reviews", bytes.NewReader(oBody))
		reqOwner.Header.Set("Authorization", "Bearer "+env.ownerToken)
		reqOwner.Header.Set("Content-Type", "application/json")
		wOwner := httptest.NewRecorder()
		env.handler.ServeHTTP(wOwner, reqOwner)

		if wOwner.Code != http.StatusCreated {
			t.Fatalf("expected 201 Created for owner review, got %d body: %s", wOwner.Code, wOwner.Body.String())
		}
		var oResp map[string]interface{}
		_ = json.Unmarshal(wOwner.Body.Bytes(), &oResp)
		oData := oResp["data"].(map[string]interface{})
		if !oData["isPublished"].(bool) {
			t.Errorf("expected second review to unblind reviews (isPublished=true), got false")
		}

		// Verify booking status transitioned to RATED
		updatedBooking, _ := env.bookingRepo.GetBookingByID(ctx, env.booking.ID)
		if updatedBooking.Status != bookingDomain.StatusRated {
			t.Errorf("expected booking status RATED, got %s", updatedBooking.Status)
		}
	})
}

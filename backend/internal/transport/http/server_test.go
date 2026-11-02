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
	"github.com/dinkmachinery/dink-core/internal/modules/auth/application"
	"github.com/dinkmachinery/dink-core/internal/modules/auth/domain"
	"github.com/dinkmachinery/dink-core/internal/modules/auth/infrastructure"
	transport "github.com/dinkmachinery/dink-core/internal/transport/http"
)

type mockUserRepo struct {
	users map[uuid.UUID]*domain.User
	byPh  map[string]*domain.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{
		users: make(map[uuid.UUID]*domain.User),
		byPh:  make(map[string]*domain.User),
	}
}

func (m *mockUserRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	u, exists := m.users[id]
	if !exists {
		return nil, domain.ErrUserNotFound
	}
	return u, nil
}

func (m *mockUserRepo) GetByPhoneNumber(ctx context.Context, phone string) (*domain.User, error) {
	u, exists := m.byPh[phone]
	if !exists {
		return nil, domain.ErrUserNotFound
	}
	return u, nil
}

func (m *mockUserRepo) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	m.users[user.ID] = user
	m.byPh[user.PhoneNumber.String()] = user
	return user, nil
}

func (m *mockUserRepo) MarkPhoneVerified(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	u, exists := m.users[id]
	if !exists {
		return nil, domain.ErrUserNotFound
	}
	now := time.Now()
	u.AccountStatus = domain.AccountStatusActive
	u.PhoneVerifiedAt = &now
	return u, nil
}

func (m *mockUserRepo) UpdateAccountStatus(ctx context.Context, id uuid.UUID, status domain.AccountStatus) (*domain.User, error) {
	u, exists := m.users[id]
	if !exists {
		return nil, domain.ErrUserNotFound
	}
	u.AccountStatus = status
	return u, nil
}

func (m *mockUserRepo) UpdateKYCStatus(ctx context.Context, id uuid.UUID, status domain.KYCStatus) (*domain.User, error) {
	u, exists := m.users[id]
	if !exists {
		return nil, domain.ErrUserNotFound
	}
	u.KYCStatus = status
	return u, nil
}

func (m *mockUserRepo) IncrementTokenVersion(ctx context.Context, id uuid.UUID) (int32, error) {
	u, exists := m.users[id]
	if !exists {
		return 0, domain.ErrUserNotFound
	}
	u.TokenVersion++
	return u.TokenVersion, nil
}

func setupTestRouter(t *testing.T) (http.Handler, *application.TokenService) {
	userRepo := newMockUserRepo()
	cache := infrastructure.NewMemoryCache()
	smsGateway := infrastructure.NewMockSMSGateway(nil)

	otpService := application.NewOTPService(cache, smsGateway, application.OTPServiceConfig{MockOTP: true})
	tokenService, err := application.NewTokenService(nil, cache)
	if err != nil {
		t.Fatalf("failed to create TokenService: %v", err)
	}

	authService := application.NewAuthService(userRepo, otpService, tokenService)
	server := transport.NewServer(authService, tokenService, nil, nil, nil)

	mux := http.NewServeMux()
	server.RegisterRoutes(mux)

	return mux, tokenService
}

func TestHTTP_AuthFullLifecycle(t *testing.T) {
	handler, _ := setupTestRouter(t)

	// 1. Request OTP
	reqBody := map[string]string{
		"phoneNumber": "+251911442910",
		"purpose":     "LOGIN",
	}
	reqBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/otp/request", bytes.NewReader(reqBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Country-Code", "ET")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("RequestOtp expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var otpResp struct {
		Success bool `json:"success"`
		Data    struct {
			VerificationSessionID string `json:"verificationSessionId"`
			TimeoutSeconds        int    `json:"timeoutSeconds"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &otpResp); err != nil {
		t.Fatalf("failed to parse RequestOtp response: %v", err)
	}
	if !otpResp.Success || otpResp.Data.VerificationSessionID == "" {
		t.Fatalf("expected valid verificationSessionId in response")
	}

	sessionID := otpResp.Data.VerificationSessionID

	// 2. Verify OTP with correct code ("123456" in mock mode)
	verifyReqBody := map[string]string{
		"verificationSessionId": sessionID,
		"code":                  "123456",
	}
	verifyBytes, _ := json.Marshal(verifyReqBody)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/otp/verify", bytes.NewReader(verifyBytes))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("VerifyOtp expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var verifyResp struct {
		Success bool `json:"success"`
		Data    struct {
			AccessToken      string `json:"accessToken"`
			RefreshToken     string `json:"refreshToken"`
			TokenType        string `json:"tokenType"`
			ExpiresInSeconds int    `json:"expiresInSeconds"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &verifyResp); err != nil {
		t.Fatalf("failed to parse VerifyOtp response: %v", err)
	}
	if verifyResp.Data.AccessToken == "" || verifyResp.Data.RefreshToken == "" {
		t.Fatalf("expected non-empty tokens in VerifyOtp response")
	}

	accessToken := verifyResp.Data.AccessToken
	refreshToken := verifyResp.Data.RefreshToken

	// 3. Authenticated Profile Inspection (GET /auth/users/me)
	req = httptest.NewRequest(http.MethodGet, "/api/v1/auth/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	w = httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("GetCurrentUserProfile expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var profileResp struct {
		Success bool `json:"success"`
		Data    struct {
			PhoneNumber string `json:"phoneNumber"`
			Role        string `json:"role"`
			Status      string `json:"status"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &profileResp); err != nil {
		t.Fatalf("failed to parse profile response: %v", err)
	}
	if profileResp.Data.PhoneNumber != "+251911442910" {
		t.Errorf("expected phone +251911442910, got %s", profileResp.Data.PhoneNumber)
	}

	// 4. Token Refresh (POST /auth/session/refresh)
	refreshReqBody := map[string]string{
		"refreshToken": refreshToken,
	}
	refreshBytes, _ := json.Marshal(refreshReqBody)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/session/refresh", bytes.NewReader(refreshBytes))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("RefreshSession expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var refreshResp struct {
		Success bool `json:"success"`
		Data    struct {
			AccessToken  string `json:"accessToken"`
			RefreshToken string `json:"refreshToken"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &refreshResp); err != nil {
		t.Fatalf("failed to parse refresh response: %v", err)
	}
	if refreshResp.Data.AccessToken == "" || refreshResp.Data.AccessToken == accessToken {
		t.Fatalf("expected new access token upon refresh")
	}

	// 5. Test Token Reuse Detection (Re-using old refresh token)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/session/refresh", bytes.NewReader(refreshBytes))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401 for compromised token reuse, got %d", w.Code)
	}
}

func TestHTTP_RequestOtp_InvalidPhone(t *testing.T) {
	handler, _ := setupTestRouter(t)

	reqBody := map[string]string{
		"phoneNumber": "not-a-phone",
		"purpose":     "LOGIN",
	}
	reqBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/otp/request", bytes.NewReader(reqBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 for invalid phone, got %d", w.Code)
	}

	var errResp struct {
		Success bool `json:"success"`
		Error   struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &errResp)
	if errResp.Error.Code != "ERR_INVALID_PHONE_NUMBER" {
		t.Errorf("expected code ERR_INVALID_PHONE_NUMBER, got %s", errResp.Error.Code)
	}
}

package rbac_test

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/dinkmachinery/dink-core/internal/modules/auth/application"
	"github.com/dinkmachinery/dink-core/internal/modules/auth/domain"
	"github.com/dinkmachinery/dink-core/internal/modules/auth/infrastructure"
	"github.com/dinkmachinery/dink-core/internal/platform/auth/rbac"
	"github.com/dinkmachinery/dink-core/internal/platform/i18n"
	"github.com/dinkmachinery/dink-core/internal/platform/response"
)

func setupTestAuth(t *testing.T) (*application.TokenService, *ecdsa.PrivateKey) {
	t.Helper()
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}
	cache := infrastructure.NewMemoryCache()
	tokenService, err := application.NewTokenService(privKey, cache)
	if err != nil {
		t.Fatalf("failed to create token service: %v", err)
	}
	return tokenService, privKey
}

func TestRequirePermission_DefaultDeny(t *testing.T) {
	tokenService, _ := setupTestAuth(t)
	checker := rbac.NewMemoryPermissionChecker()

	userID := uuid.New()
	targetPerm := "fleets:approve"

	testUser := &domain.User{
		ID:           userID,
		CountryCode:  "ET",
		TokenVersion: 1,
		PrimaryRole:  domain.RoleOwner,
	}

	// Create valid token for user
	tokenPair, err := tokenService.MintTokenPair(context.Background(), testUser, "")
	if err != nil {
		t.Fatalf("failed to mint token: %v", err)
	}

	successHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok": true}`))
	})

	middleware := rbac.RequirePermission(tokenService, checker, targetPerm)
	protectedRoute := i18n.LanguageMiddleware(middleware(successHandler))

	// 1. User without permission receives 403 Forbidden
	req := httptest.NewRequest("POST", "/api/v1/fleets/approve", nil)
	req.Header.Set("Authorization", "Bearer "+tokenPair.AccessToken)
	rr := httptest.NewRecorder()

	protectedRoute.ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", rr.Code)
	}

	var errResp response.ErrorEnvelope
	if err := json.NewDecoder(rr.Body).Decode(&errResp); err != nil {
		t.Fatalf("failed to parse error response: %v", err)
	}
	if errResp.Error.Code != "ERR_FORBIDDEN_INSUFFICIENT_PERMISSIONS" {
		t.Errorf("expected ERR_FORBIDDEN_INSUFFICIENT_PERMISSIONS, got %q", errResp.Error.Code)
	}

	// 2. Localized Amharic 403 Forbidden
	reqAm := httptest.NewRequest("POST", "/api/v1/fleets/approve", nil)
	reqAm.Header.Set("Authorization", "Bearer "+tokenPair.AccessToken)
	reqAm.Header.Set("Accept-Language", "am-ET,am;q=0.9")
	rrAm := httptest.NewRecorder()

	protectedRoute.ServeHTTP(rrAm, reqAm)

	if rrAm.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", rrAm.Code)
	}

	var errRespAm response.ErrorEnvelope
	if err := json.NewDecoder(rrAm.Body).Decode(&errRespAm); err != nil {
		t.Fatalf("failed to parse error response: %v", err)
	}
	if errRespAm.Error.Title != "የተከለከለ እርምጃ" {
		t.Errorf("expected Amharic title 'የተከለከለ እርምጃ', got %q", errRespAm.Error.Title)
	}

	// 3. Granted permission allows request through
	checker.GrantPermission(userID, targetPerm)

	reqAllowed := httptest.NewRequest("POST", "/api/v1/fleets/approve", nil)
	reqAllowed.Header.Set("Authorization", "Bearer "+tokenPair.AccessToken)
	rrAllowed := httptest.NewRecorder()

	protectedRoute.ServeHTTP(rrAllowed, reqAllowed)

	if rrAllowed.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rrAllowed.Code)
	}
}

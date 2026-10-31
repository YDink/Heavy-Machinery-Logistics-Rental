package rbac

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/dinkmachinery/dink-core/internal/modules/auth/application"
	"github.com/dinkmachinery/dink-core/internal/platform/response"
)

type userContextKey string

const AuthenticatedUserKey userContextKey = "dink_authenticated_user_id"

// PermissionChecker verifies whether an identified user holds an authorization permission key.
type PermissionChecker interface {
	HasPermission(ctx context.Context, userID uuid.UUID, permission string) (bool, error)
	GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]string, error)
}

// MemoryPermissionChecker provides a thread-safe in-memory permission store for testing and fast local caching.
type MemoryPermissionChecker struct {
	permissions map[uuid.UUID]map[string]bool
}

func NewMemoryPermissionChecker() *MemoryPermissionChecker {
	return &MemoryPermissionChecker{
		permissions: make(map[uuid.UUID]map[string]bool),
	}
}

func (m *MemoryPermissionChecker) GrantPermission(userID uuid.UUID, permission string) {
	if m.permissions[userID] == nil {
		m.permissions[userID] = make(map[string]bool)
	}
	m.permissions[userID][permission] = true
}

func (m *MemoryPermissionChecker) HasPermission(ctx context.Context, userID uuid.UUID, permission string) (bool, error) {
	if perms, exists := m.permissions[userID]; exists {
		return perms[permission], nil
	}
	return false, nil
}

func (m *MemoryPermissionChecker) GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]string, error) {
	var result []string
	if perms, exists := m.permissions[userID]; exists {
		for p := range perms {
			result = append(result, p)
		}
	}
	return result, nil
}

// RequirePermission constructs a default-deny HTTP middleware asserting the caller holds the targeted permission.
func RequirePermission(tokenService *application.TokenService, checker PermissionChecker, permissionCode string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// 1. Resolve User ID from context or Authorization header
			var userID uuid.UUID
			if val, ok := ctx.Value(AuthenticatedUserKey).(uuid.UUID); ok {
				userID = val
			} else {
				authHeader := r.Header.Get("Authorization")
				if !strings.HasPrefix(authHeader, "Bearer ") {
					response.LocalizedError(w, r, http.StatusUnauthorized,
						"ERR_UNAUTHORIZED",
						"errors.unauthorized_title",
						"errors.unauthorized_detail",
						"Unauthorized",
						"Authentication credentials are missing, invalid, or expired.")
					return
				}

				tokenString := strings.TrimPrefix(authHeader, "Bearer ")
				claims, err := tokenService.VerifyAccessToken(ctx, tokenString)
				if err != nil {
					response.LocalizedError(w, r, http.StatusUnauthorized,
						"ERR_UNAUTHORIZED",
						"errors.unauthorized_title",
						"errors.unauthorized_detail",
						"Unauthorized",
						"Access token is invalid, expired, or revoked.")
					return
				}

				parsedID, err := uuid.Parse(claims.Subject)
				if err != nil {
					response.LocalizedError(w, r, http.StatusUnauthorized,
						"ERR_INVALID_SUBJECT",
						"errors.unauthorized_title",
						"errors.unauthorized_detail",
						"Unauthorized",
						"Subject token claim is not a valid UUID.")
					return
				}
				userID = parsedID
				ctx = context.WithValue(ctx, AuthenticatedUserKey, userID)
			}

			// 2. Evaluate Dynamic RBAC Permission (Default-Deny)
			hasPerm, err := checker.HasPermission(ctx, userID, permissionCode)
			if err != nil || !hasPerm {
				response.LocalizedError(w, r, http.StatusForbidden,
					"ERR_FORBIDDEN_INSUFFICIENT_PERMISSIONS",
					"errors.forbidden_permission_title",
					"errors.forbidden_permission_detail",
					"Forbidden",
					"You do not have the required permissions to perform this action.")
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

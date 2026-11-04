package http

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/dinkmachinery/dink-core/internal/platform/response"
	"github.com/dinkmachinery/dink-core/internal/transport/http/generated"
)

// GetCurrentUserProfile handles GET /api/v1/auth/users/me
func (s *Server) GetCurrentUserProfile(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		response.Error(w, r, http.StatusUnauthorized, response.ProblemDetails{
			Title:  "Unauthorized",
			Code:   "ERR_UNAUTHORIZED",
			Detail: "Missing or malformed Authorization header.",
		})
		return
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := s.tokenService.VerifyAccessToken(r.Context(), tokenStr)
	if err != nil {
		response.Error(w, r, http.StatusUnauthorized, response.ProblemDetails{
			Title:  "Unauthorized",
			Code:   "ERR_INVALID_TOKEN",
			Detail: "Access token is invalid or expired.",
		})
		return
	}

	userUUID, err := uuid.Parse(claims.Subject)
	if err != nil {
		response.Error(w, r, http.StatusUnauthorized, response.ProblemDetails{
			Title:  "Unauthorized",
			Code:   "ERR_INVALID_SUBJECT",
			Detail: "Subject claim is not a valid UUID.",
		})
		return
	}

	user, err := s.authService.GetUserProfile(r.Context(), userUUID)
	if err != nil {
		response.Error(w, r, http.StatusNotFound, response.ProblemDetails{
			Title:  "User Not Found",
			Code:   "ERR_USER_NOT_FOUND",
			Detail: "User profile corresponding to token could not be found.",
		})
		return
	}

	userProfile := generated.UserProfile{
		Id:          user.ID,
		PhoneNumber: user.PhoneNumber.String(),
		FullName:    &user.FullName,
		Role:        generated.UserProfileRole(user.PrimaryRole),
		Status:      generated.UserProfileStatus(user.AccountStatus),
		CreatedAt:   user.CreatedAt,
	}

	response.JSON(w, r, http.StatusOK, userProfile)
}

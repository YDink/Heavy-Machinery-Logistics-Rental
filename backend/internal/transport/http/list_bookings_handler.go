package http

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/dinkmachinery/dink-core/internal/platform/response"
	"github.com/dinkmachinery/dink-core/internal/transport/http/generated"
)

// ListBookings handles GET /api/v1/bookings
func (s *Server) ListBookings(w http.ResponseWriter, r *http.Request) {
	if s.bookingService == nil {
		s.notImplemented(w, r, "05-05")
		return
	}

	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		response.Error(w, r, http.StatusUnauthorized, response.ProblemDetails{
			Title:  "Unauthorized",
			Code:   "ERR_UNAUTHORIZED",
			Detail: "Missing or malformed Authorization header.",
		})
		return
	}

	claims, err := s.tokenService.VerifyAccessToken(r.Context(), strings.TrimPrefix(authHeader, "Bearer "))
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

	statusFilter := r.URL.Query().Get("status")
	bookings, err := s.bookingService.ListUserBookings(r.Context(), userUUID, claims.Role, statusFilter, 50)
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, response.ProblemDetails{
			Title:  "Internal Server Error",
			Code:   "ERR_INTERNAL",
			Detail: "Failed to list user bookings.",
		})
		return
	}

	dtos := make([]generated.Booking, len(bookings))
	for i, b := range bookings {
		dtos[i] = toGeneratedBooking(b)
	}

	response.JSON(w, r, http.StatusOK, dtos)
}

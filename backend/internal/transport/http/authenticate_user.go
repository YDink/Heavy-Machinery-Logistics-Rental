package http

import (
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func (s *Server) authenticateUser(r *http.Request) (uuid.UUID, error) {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return uuid.Nil, errors.New("missing or malformed Authorization header")
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := s.tokenService.VerifyAccessToken(r.Context(), tokenStr)
	if err != nil {
		return uuid.Nil, errors.New("invalid or expired token")
	}

	userUUID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, errors.New("invalid subject in token")
	}

	return userUUID, nil
}

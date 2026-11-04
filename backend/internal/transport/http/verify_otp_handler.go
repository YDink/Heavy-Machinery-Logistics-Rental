package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dinkmachinery/dink-core/internal/modules/auth/application"
	"github.com/dinkmachinery/dink-core/internal/modules/auth/domain"
	"github.com/dinkmachinery/dink-core/internal/platform/response"
	"github.com/dinkmachinery/dink-core/internal/transport/http/generated"
)

// VerifyOtp handles POST /api/v1/auth/otp/verify
func (s *Server) VerifyOtp(w http.ResponseWriter, r *http.Request) {
	var body generated.VerifyOtpJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Bad Request",
			Code:   "ERR_MALFORMED_JSON",
			Detail: "Request body could not be parsed as valid JSON.",
		})
		return
	}

	input := application.VerifyOTPInput{
		VerificationSessionID: body.VerificationSessionId.String(),
		Code:                  body.Code,
	}

	tokenPair, _, err := s.authService.VerifyOTP(r.Context(), input)
	if err != nil {
		if errors.Is(err, domain.ErrOTPExpired) {
			response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
				Title:  "OTP Expired",
				Code:   "ERR_OTP_EXPIRED",
				Detail: "The verification code has expired. Please request a new code.",
			})
			return
		}
		if errors.Is(err, domain.ErrOTPMaxAttemptsExceeded) {
			response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
				Title:  "Max Attempts Exceeded",
				Code:   "ERR_OTP_MAX_ATTEMPTS_EXCEEDED",
				Detail: "Maximum verification attempts exceeded. Challenge has been invalidated.",
			})
			return
		}
		if errors.Is(err, domain.ErrOTPInvalid) {
			response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
				Title:  "Invalid Code",
				Code:   "ERR_OTP_INVALID",
				Detail: "The verification code entered is incorrect.",
			})
			return
		}
		if errors.Is(err, domain.ErrAccountSuspended) {
			response.Error(w, r, http.StatusForbidden, response.ProblemDetails{
				Title:  "Account Suspended",
				Code:   "ERR_ACCOUNT_SUSPENDED",
				Detail: "Your account is temporarily suspended. Please contact support.",
			})
			return
		}

		response.Error(w, r, http.StatusInternalServerError, response.ProblemDetails{
			Title:  "Internal Server Error",
			Code:   "ERR_INTERNAL",
			Detail: "Failed to verify OTP code.",
		})
		return
	}

	data := generated.AuthTokenPair{
		AccessToken:      tokenPair.AccessToken,
		RefreshToken:     tokenPair.RefreshToken,
		TokenType:        tokenPair.TokenType,
		ExpiresInSeconds: tokenPair.ExpiresInSeconds,
	}

	response.JSON(w, r, http.StatusOK, data)
}

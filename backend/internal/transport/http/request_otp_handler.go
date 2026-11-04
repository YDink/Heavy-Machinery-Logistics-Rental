package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/dinkmachinery/dink-core/internal/modules/auth/application"
	"github.com/dinkmachinery/dink-core/internal/modules/auth/domain"
	"github.com/dinkmachinery/dink-core/internal/platform/response"
	"github.com/dinkmachinery/dink-core/internal/transport/http/generated"
)

// RequestOtp handles POST /api/v1/auth/otp/request
func (s *Server) RequestOtp(w http.ResponseWriter, r *http.Request) {
	var body generated.RequestOtpJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Bad Request",
			Code:   "ERR_MALFORMED_JSON",
			Detail: "Request body could not be parsed as valid JSON.",
		})
		return
	}

	countryCode := r.Header.Get("X-Country-Code")
	if countryCode == "" {
		countryCode = "ET"
	}

	input := application.RequestOTPInput{
		PhoneNumber: body.PhoneNumber,
		CountryCode: countryCode,
		Purpose:     domain.OTPPurpose(body.Purpose),
	}

	out, err := s.authService.RequestOTP(r.Context(), input)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidPhoneNumber) || errors.Is(err, domain.ErrUnsupportedCountryCode) {
			response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
				Title:  "Invalid Phone Number",
				Code:   "ERR_INVALID_PHONE_NUMBER",
				Detail: err.Error(),
			})
			return
		}
		if errors.Is(err, domain.ErrRateLimitExceeded) {
			response.Error(w, r, http.StatusTooManyRequests, response.ProblemDetails{
				Title:  "Too Many Requests",
				Code:   "ERR_RATE_LIMIT_EXCEEDED",
				Detail: "You have exceeded the maximum number of OTP requests. Please wait 15 minutes before retrying.",
			})
			return
		}

		response.Error(w, r, http.StatusInternalServerError, response.ProblemDetails{
			Title:  "Internal Server Error",
			Code:   "ERR_INTERNAL",
			Detail: "Failed to dispatch OTP verification code.",
		})
		return
	}

	parsedUUID, _ := uuid.Parse(out.VerificationSessionID)
	data := map[string]interface{}{
		"verificationSessionId": parsedUUID,
		"timeoutSeconds":        out.TimeoutSeconds,
	}

	response.JSON(w, r, http.StatusOK, data)
}

package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/dinkmachinery/dink-core/internal/modules/auth/application"
	"github.com/dinkmachinery/dink-core/internal/modules/auth/domain"
	"github.com/dinkmachinery/dink-core/internal/platform/response"
	"github.com/dinkmachinery/dink-core/internal/transport/http/generated"
)

type Server struct {
	authService  *application.AuthService
	tokenService *application.TokenService
}

func NewServer(
	authService *application.AuthService,
	tokenService *application.TokenService,
) *Server {
	return &Server{
		authService:  authService,
		tokenService: tokenService,
	}
}

// Ensure Server implements generated.ServerInterface
var _ generated.ServerInterface = (*Server)(nil)

// RequestOtp handles POST /auth/otp/request
func (s *Server) RequestOtp(w http.ResponseWriter, r *http.Request, params generated.RequestOtpParams) {
	var body generated.RequestOtpJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Bad Request",
			Code:   "ERR_MALFORMED_JSON",
			Detail: "Request body could not be parsed as valid JSON.",
		})
		return
	}

	countryCode := "ET"
	if params.XCountryCode != nil && *params.XCountryCode != "" {
		countryCode = string(*params.XCountryCode)
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

// VerifyOtp handles POST /auth/otp/verify
func (s *Server) VerifyOtp(w http.ResponseWriter, r *http.Request, params generated.VerifyOtpParams) {
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
				Detail: "Account is suspended or locked. Please contact compliance support.",
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

// RefreshSession handles POST /auth/session/refresh
func (s *Server) RefreshSession(w http.ResponseWriter, r *http.Request) {
	var body generated.RefreshSessionJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Bad Request",
			Code:   "ERR_MALFORMED_JSON",
			Detail: "Request body could not be parsed as valid JSON.",
		})
		return
	}

	tokenPair, err := s.authService.RefreshSession(r.Context(), body.RefreshToken)
	if err != nil {
		if errors.Is(err, domain.ErrTokenFamilyCompromised) {
			response.Error(w, r, http.StatusUnauthorized, response.ProblemDetails{
				Title:  "Token Compromised",
				Code:   "ERR_TOKEN_FAMILY_COMPROMISED",
				Detail: "A previously used refresh token was presented. All active sessions in this family have been terminated for security.",
			})
			return
		}
		if errors.Is(err, domain.ErrInvalidToken) || errors.Is(err, domain.ErrUserNotFound) {
			response.Error(w, r, http.StatusUnauthorized, response.ProblemDetails{
				Title:  "Invalid Token",
				Code:   "ERR_INVALID_TOKEN",
				Detail: "The refresh token is invalid or does not exist.",
			})
			return
		}

		response.Error(w, r, http.StatusInternalServerError, response.ProblemDetails{
			Title:  "Internal Server Error",
			Code:   "ERR_INTERNAL",
			Detail: "Failed to refresh session.",
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

// GetCurrentUserProfile handles GET /auth/users/me
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

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := s.tokenService.VerifyAccessToken(r.Context(), tokenString)
	if err != nil {
		response.Error(w, r, http.StatusUnauthorized, response.ProblemDetails{
			Title:  "Unauthorized",
			Code:   "ERR_INVALID_ACCESS_TOKEN",
			Detail: "Access token is invalid, expired, or revoked.",
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

// Unimplemented Stubs for Future Partitions
func (s *Server) notImplemented(w http.ResponseWriter, r *http.Request, partition string) {
	response.Error(w, r, http.StatusNotImplemented, response.ProblemDetails{
		Title:  "Not Implemented",
		Code:   "ERR_PARTITION_PENDING",
		Detail: "This endpoint will be activated during partition " + partition + ".",
	})
}

func (s *Server) ListCountryConfigs(w http.ResponseWriter, r *http.Request) { s.notImplemented(w, r, "05-11") }
func (s *Server) GetCountryConfigByCode(w http.ResponseWriter, r *http.Request, countryCode string) { s.notImplemented(w, r, "05-11") }
func (s *Server) UpdateCountryConfig(w http.ResponseWriter, r *http.Request, countryCode string, params generated.UpdateCountryConfigParams) { s.notImplemented(w, r, "05-11") }
func (s *Server) GetAdminUiConfig(w http.ResponseWriter, r *http.Request) { s.notImplemented(w, r, "05-11") }
func (s *Server) ListBookings(w http.ResponseWriter, r *http.Request, params generated.ListBookingsParams) { s.notImplemented(w, r, "05-05") }
func (s *Server) CreateBooking(w http.ResponseWriter, r *http.Request, params generated.CreateBookingParams) { s.notImplemented(w, r, "05-05") }
func (s *Server) GetBookingById(w http.ResponseWriter, r *http.Request, id openapi_types.UUID) { s.notImplemented(w, r, "05-05") }
func (s *Server) SubmitBookingInspection(w http.ResponseWriter, r *http.Request, id openapi_types.UUID, params generated.SubmitBookingInspectionParams) { s.notImplemented(w, r, "05-07") }
func (s *Server) CancelBooking(w http.ResponseWriter, r *http.Request, id openapi_types.UUID, params generated.CancelBookingParams) { s.notImplemented(w, r, "05-05") }
func (s *Server) ConfirmBookingPayment(w http.ResponseWriter, r *http.Request, id openapi_types.UUID, params generated.ConfirmBookingPaymentParams) { s.notImplemented(w, r, "05-05") }
func (s *Server) ListCategories(w http.ResponseWriter, r *http.Request) { s.notImplemented(w, r, "05-04") }
func (s *Server) ListFleets(w http.ResponseWriter, r *http.Request, params generated.ListFleetsParams) { s.notImplemented(w, r, "05-04") }
func (s *Server) CreateFleetAsset(w http.ResponseWriter, r *http.Request, params generated.CreateFleetAssetParams) { s.notImplemented(w, r, "05-02") }
func (s *Server) GetFleetById(w http.ResponseWriter, r *http.Request, id openapi_types.UUID) { s.notImplemented(w, r, "05-04") }
func (s *Server) SearchFleetsByRadius(w http.ResponseWriter, r *http.Request, params generated.SearchFleetsByRadiusParams) { s.notImplemented(w, r, "05-04") }
func (s *Server) InitiatePayment(w http.ResponseWriter, r *http.Request, params generated.InitiatePaymentParams) { s.notImplemented(w, r, "05-06") }
func (s *Server) GetMyWallet(w http.ResponseWriter, r *http.Request) { s.notImplemented(w, r, "05-09") }
func (s *Server) ListWalletTransactions(w http.ResponseWriter, r *http.Request, params generated.ListWalletTransactionsParams) { s.notImplemented(w, r, "05-09") }

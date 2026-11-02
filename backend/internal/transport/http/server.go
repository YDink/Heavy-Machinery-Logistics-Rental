package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/dinkmachinery/dink-core/internal/modules/auth/application"
	"github.com/dinkmachinery/dink-core/internal/modules/auth/domain"
	machineryApp "github.com/dinkmachinery/dink-core/internal/modules/machinery/application"
	machineryDomain "github.com/dinkmachinery/dink-core/internal/modules/machinery/domain"
	"github.com/dinkmachinery/dink-core/internal/platform/i18n"
	"github.com/dinkmachinery/dink-core/internal/platform/response"
	"github.com/dinkmachinery/dink-core/internal/transport/http/generated"
)

type Server struct {
	authService       *application.AuthService
	tokenService      *application.TokenService
	machineryService  *machineryApp.MachineryService
	attachmentService *machineryApp.AttachmentService
	discoveryService  *machineryApp.DiscoveryService
}

func NewServer(
	authService *application.AuthService,
	tokenService *application.TokenService,
	machineryService *machineryApp.MachineryService,
	attachmentService *machineryApp.AttachmentService,
	discoveryService *machineryApp.DiscoveryService,
) *Server {
	return &Server{
		authService:       authService,
		tokenService:      tokenService,
		machineryService:  machineryService,
		attachmentService: attachmentService,
		discoveryService:  discoveryService,
	}
}

// RegisterRoutes registers all API endpoints directly onto the Go standard library http.ServeMux.
func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	// Authentication (Partition 05-01)
	mux.HandleFunc("POST /api/v1/auth/otp/request", s.RequestOtp)
	mux.HandleFunc("POST /api/v1/auth/otp/verify", s.VerifyOtp)
	mux.HandleFunc("POST /api/v1/auth/session/refresh", s.RefreshSession)
	mux.HandleFunc("GET /api/v1/auth/users/me", s.GetCurrentUserProfile)

	// Machinery Categories & Dynamic Specs (Partition 05-02)
	mux.HandleFunc("GET /api/v1/categories", s.ListCategories)
	mux.HandleFunc("POST /api/v1/fleets", s.CreateFleetAsset)
	mux.HandleFunc("GET /api/v1/fleets", s.ListFleets)
	mux.HandleFunc("GET /api/v1/fleets/{id}", s.GetFleetById)

	// Attachments & Bundling (Partition 05-03)
	mux.HandleFunc("POST /api/v1/attachments", s.CreateAttachment)
	mux.HandleFunc("GET /api/v1/attachments/inventory", s.ListAttachmentInventory)
	mux.HandleFunc("POST /api/v1/fleets/{id}/attachments/link", s.LinkFleetAttachment)
	mux.HandleFunc("POST /api/v1/bookings/check-bundled-availability", s.CheckBundledAvailability)

	// Spatial Discovery & Dynamic Comparison (Partition 05-04)
	mux.HandleFunc("GET /api/v1/fleets/search", s.SearchFleets)
	mux.HandleFunc("POST /api/v1/fleets/compare", s.CompareFleets)
	mux.HandleFunc("GET /api/v1/fleets/{id}/availability", s.GetFleetAvailability)
	mux.HandleFunc("GET /api/v1/geo/fleets/radius-search", s.SearchFleetsByRadius)

	// Stubs for Upcoming Partitions
	mux.HandleFunc("GET /api/v1/country-configs", s.ListCountryConfigs)
	mux.HandleFunc("GET /api/v1/country-configs/{code}", s.GetCountryConfigByCode)
	mux.HandleFunc("PUT /api/v1/country-configs/{code}", s.UpdateCountryConfig)
	mux.HandleFunc("GET /api/v1/admin/ui-config", s.GetAdminUiConfig)
	mux.HandleFunc("GET /api/v1/bookings", s.ListBookings)
	mux.HandleFunc("POST /api/v1/bookings", s.CreateBooking)
	mux.HandleFunc("GET /api/v1/bookings/{id}", s.GetBookingById)
	mux.HandleFunc("POST /api/v1/bookings/{id}/inspections", s.SubmitBookingInspection)
	mux.HandleFunc("POST /api/v1/bookings/{id}/cancel", s.CancelBooking)
	mux.HandleFunc("POST /api/v1/bookings/{id}/confirm-payment", s.ConfirmBookingPayment)
	mux.HandleFunc("POST /api/v1/payments", s.InitiatePayment)
	mux.HandleFunc("GET /api/v1/wallets/me", s.GetMyWallet)
	mux.HandleFunc("GET /api/v1/wallets/me/transactions", s.ListWalletTransactions)
}

// RequestOtp handles POST /auth/otp/request
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

// VerifyOtp handles POST /auth/otp/verify
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

// ==================== Partition 05-02 Machinery & Categories ====================

// ListCategories handles GET /categories
func (s *Server) ListCategories(w http.ResponseWriter, r *http.Request) {
	if s.machineryService == nil {
		s.notImplemented(w, r, "05-02")
		return
	}

	lang := i18n.GetLanguage(r.Context())
	categories, err := s.machineryService.ListCategories(r.Context())
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, response.ProblemDetails{
			Title:  "Internal Server Error",
			Code:   "ERR_INTERNAL",
			Detail: "Failed to list machinery categories.",
		})
		return
	}

	result := make([]generated.Category, len(categories))
	for i, c := range categories {
		name := c.NameEN
		if lang == "am" && c.NameAM != "" {
			name = c.NameAM
		}
		var parentID *openapi_types.UUID
		if c.ParentID != nil {
			pid := openapi_types.UUID(*c.ParentID)
			parentID = &pid
		}
		result[i] = generated.Category{
			Id:       openapi_types.UUID(c.ID),
			Code:     c.Code,
			Name:     name,
			IconUrl:  c.IconURL,
			ParentId: parentID,
		}
	}

	response.JSON(w, r, http.StatusOK, result)
}

// GetFleetById handles GET /fleets/{id}
func (s *Server) GetFleetById(w http.ResponseWriter, r *http.Request) {
	if s.machineryService == nil {
		s.notImplemented(w, r, "05-02")
		return
	}

	idStr := r.PathValue("id")
	fleetID, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Invalid Fleet ID",
			Code:   "ERR_INVALID_UUID",
			Detail: "The provided fleet ID is not a valid UUID.",
		})
		return
	}

	asset, err := s.machineryService.GetFleetAsset(r.Context(), fleetID)
	if err != nil {
		if errors.Is(err, machineryDomain.ErrFleetNotFound) {
			response.Error(w, r, http.StatusNotFound, response.ProblemDetails{
				Title:  "Fleet Asset Not Found",
				Code:   "ERR_FLEET_NOT_FOUND",
				Detail: "The requested machinery fleet asset could not be found.",
			})
			return
		}
		response.Error(w, r, http.StatusInternalServerError, response.ProblemDetails{
			Title:  "Internal Server Error",
			Code:   "ERR_INTERNAL",
			Detail: "Failed to retrieve fleet asset.",
		})
		return
	}

	dto := generated.FleetAsset{
		Id:                 openapi_types.UUID(asset.ID),
		OwnerId:            openapi_types.UUID(asset.OwnerID),
		CategoryId:         openapi_types.UUID(asset.SubcategoryID),
		Make:               &asset.Make,
		Model:              asset.Model,
		ManufacturingYear:  &asset.ManufactureYear,
		RegistrationNumber: asset.PlateNumber,
		DailyRate: generated.MoneyAmount{
			AmountMinor: asset.DailyRateMinor,
			Currency:    "ETB",
		},
		Status:  generated.FleetAssetStatus(asset.Status),
		Version: asset.Version,
	}
	dto.Location.Latitude = asset.Latitude
	dto.Location.Longitude = asset.Longitude
	dto.Location.AddressText = asset.AddressText

	response.JSON(w, r, http.StatusOK, dto)
}

// CreateFleetAsset handles POST /fleets
func (s *Server) CreateFleetAsset(w http.ResponseWriter, r *http.Request) {
	if s.machineryService == nil {
		s.notImplemented(w, r, "05-02")
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

	var body struct {
		generated.CreateFleetRequest
		ChassisVin          string  `json:"chassisVin"`
		BaselineMeterHours  float64 `json:"baselineMeterHours"`
		FuelType            string  `json:"fuelType"`
		SelfMoveCapable     bool    `json:"selfMoveCapable"`
		ExcessHourRateMinor int64   `json:"excessHourRateMinor"`
		MinRentalDays       int     `json:"minRentalDays"`
		PaymentTimeoutHours int     `json:"paymentTimeoutHours"`
		OperatorIncluded    bool    `json:"operatorIncluded"`
		Photos              []struct {
			Angle      string `json:"angle"`
			StorageURI string `json:"storageUri"`
			FileHash   string `json:"fileHash"`
		} `json:"photos"`
		SpecValues []struct {
			FieldKey     string   `json:"fieldKey"`
			ValueString  *string  `json:"valueString"`
			ValueNumeric *float64 `json:"valueNumeric"`
			ValueBoolean *bool    `json:"valueBoolean"`
		} `json:"specValues"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Bad Request",
			Code:   "ERR_MALFORMED_JSON",
			Detail: "Request body could not be parsed as valid JSON.",
		})
		return
	}

	var photos []machineryApp.CreateFleetPhotoInput
	for _, p := range body.Photos {
		photos = append(photos, machineryApp.CreateFleetPhotoInput{
			Angle:      machineryDomain.PhotoAngle(p.Angle),
			StorageURI: p.StorageURI,
			FileHash:   p.FileHash,
		})
	}

	var specValues []machineryApp.CreateFleetSpecValueInput
	for _, v := range body.SpecValues {
		specValues = append(specValues, machineryApp.CreateFleetSpecValueInput{
			FieldKey:     v.FieldKey,
			ValueString:  v.ValueString,
			ValueNumeric: v.ValueNumeric,
			ValueBoolean: v.ValueBoolean,
		})
	}

	fuel := machineryDomain.FuelTypeDiesel
	if body.FuelType != "" {
		fuel = machineryDomain.FuelType(body.FuelType)
	}

	input := machineryApp.CreateFleetInput{
		OwnerID:              userUUID,
		SubcategoryID:        uuid.UUID(body.CategoryId),
		Make:                 body.Make,
		Model:                body.Model,
		ManufactureYear:      body.ManufacturingYear,
		PlateNumber:          body.RegistrationNumber,
		ChassisVIN:           body.ChassisVin,
		BaselineMeterHours:   body.BaselineMeterHours,
		DailyRateMinor:       body.DailyRate.AmountMinor,
		FuelType:             fuel,
		SelfMoveCapable:      body.SelfMoveCapable,
		ExcessHourRateMinor:  body.ExcessHourRateMinor,
		MinRentalDays:        body.MinRentalDays,
		PaymentTimeoutHours:  body.PaymentTimeoutHours,
		OperatorIncluded:     body.OperatorIncluded,
		Latitude:             body.Location.Latitude,
		Longitude:            body.Location.Longitude,
		AddressText:          body.Location.AddressText,
		Photos:               photos,
		SpecValues:          specValues,
	}

	out, err := s.machineryService.CreateFleet(r.Context(), input)
	if err != nil {
		if errors.Is(err, machineryDomain.ErrInvalidBaselineMeter) {
			response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
				Title:  "Invalid Baseline Meter",
				Code:   "ERR_INVALID_BASELINE_METER",
				Detail: "Baseline engine meter reading must be non-negative.",
			})
			return
		}
		if errors.Is(err, machineryDomain.ErrIncompletePhotoWalkaround) {
			response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
				Title:  "Incomplete Photo Walkaround",
				Code:   "ERR_INCOMPLETE_PHOTO_WALKAROUND",
				Detail: "All 6 walkaround angles (FRONT, REAR, LEFT, RIGHT, CABIN, METER) are required.",
			})
			return
		}
		if errors.Is(err, machineryDomain.ErrCategoryNotFound) {
			response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
				Title:  "Category Not Found",
				Code:   "ERR_CATEGORY_NOT_FOUND",
				Detail: "The specified machinery category does not exist.",
			})
			return
		}
		response.Error(w, r, http.StatusInternalServerError, response.ProblemDetails{
			Title:  "Internal Server Error",
			Code:   "ERR_INTERNAL",
			Detail: "Failed to create fleet asset.",
		})
		return
	}

	createdAsset, err := s.machineryService.GetFleetAsset(r.Context(), out.FleetID)
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, response.ProblemDetails{
			Title:  "Internal Server Error",
			Code:   "ERR_INTERNAL",
			Detail: "Failed to fetch created fleet asset.",
		})
		return
	}

	dto := generated.FleetAsset{
		Id:                 openapi_types.UUID(createdAsset.ID),
		OwnerId:            openapi_types.UUID(createdAsset.OwnerID),
		CategoryId:         openapi_types.UUID(createdAsset.SubcategoryID),
		Make:               &createdAsset.Make,
		Model:              createdAsset.Model,
		ManufacturingYear:  &createdAsset.ManufactureYear,
		RegistrationNumber: createdAsset.PlateNumber,
		DailyRate: generated.MoneyAmount{
			AmountMinor: createdAsset.DailyRateMinor,
			Currency:    "ETB",
		},
		Status:  generated.FleetAssetStatus(createdAsset.Status),
		Version: createdAsset.Version,
	}
	dto.Location.Latitude = createdAsset.Latitude
	dto.Location.Longitude = createdAsset.Longitude
	dto.Location.AddressText = createdAsset.AddressText

	response.JSON(w, r, http.StatusCreated, dto)
}

// ListFleets handles GET /fleets
func (s *Server) ListFleets(w http.ResponseWriter, r *http.Request) {
	if s.machineryService == nil {
		s.notImplemented(w, r, "05-04")
		return
	}
	response.JSON(w, r, http.StatusOK, []generated.FleetAsset{})
}

// ==================== Attachments Handlers (Partition 05-03) ====================

// CreateAttachment handles POST /attachments (ATT-01 to ATT-14)
func (s *Server) CreateAttachment(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		response.Error(w, r, http.StatusUnauthorized, response.ProblemDetails{
			Title:  "Unauthorized",
			Code:   "ERR_UNAUTHORIZED",
			Detail: "Missing or invalid Bearer authorization token.",
		})
		return
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := s.tokenService.VerifyAccessToken(r.Context(), tokenStr)
	if err != nil {
		response.Error(w, r, http.StatusUnauthorized, response.ProblemDetails{
			Title:  "Unauthorized",
			Code:   "ERR_UNAUTHORIZED",
			Detail: "Access token is expired or signature is invalid.",
		})
		return
	}
	userUUID, err := uuid.Parse(claims.Subject)
	if err != nil {
		response.Error(w, r, http.StatusUnauthorized, response.ProblemDetails{
			Title:  "Unauthorized",
			Code:   "ERR_UNAUTHORIZED",
			Detail: "Invalid token subject claim.",
		})
		return
	}

	var body generated.CreateAttachmentJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Bad Request",
			Code:   "ERR_MALFORMED_JSON",
			Detail: "Request body could not be parsed as valid JSON.",
		})
		return
	}

	photos := make([]machineryDomain.AttachmentPhoto, 0)
	if body.Photos != nil {
		for _, p := range *body.Photos {
			photos = append(photos, machineryDomain.AttachmentPhoto{
				Angle:      machineryDomain.AttachmentPhotoAngle(p.Angle),
				StorageURI: p.StorageUri,
			})
		}
	}

	secDeposit := int64(0)
	if body.SecurityDepositMinor != nil {
		secDeposit = *body.SecurityDepositMinor
	}
	minDays := 1
	if body.MinRentalDays != nil {
		minDays = *body.MinRentalDays
	}

	att, err := s.attachmentService.RegisterAttachment(r.Context(), machineryApp.RegisterAttachmentInput{
		OwnerID:                    userUUID,
		AttachmentType:             machineryDomain.AttachmentType(body.AttachmentType),
		Make:                       body.Make,
		Model:                      body.Model,
		SubTypeSize:                body.SubTypeSize,
		SerialNumber:               body.SerialNumber,
		OperatingWeightKg:          float64(body.OperatingWeightKg),
		MinCarrierWeightTonnes:     float64(body.MinCarrierWeightTonnes),
		MaxCarrierWeightTonnes:     float64(body.MaxCarrierWeightTonnes),
		RequiredFlowMinLmin:        float64(body.RequiredFlowMinLmin),
		RequiredFlowMaxLmin:        float64(body.RequiredFlowMaxLmin),
		RequiredPressureBar:        float64(body.RequiredPressureBar),
		PinDiameterMm:              float64(body.PinDiameterMm),
		PinSpreadMm:                float64(body.PinSpreadMm),
		CompatibleSubcategories:    body.CompatibleSubcategories,
		StandaloneDailyRateMinor:   body.StandaloneDailyRateMinor,
		SecurityDepositMinor:       secDeposit,
		MinRentalDays:              minDays,
		PaymentEvidenceWindowHours: 24,
		Photos:                     photos,
	})
	if err != nil {
		if errors.Is(err, machineryDomain.ErrIncompleteAttachmentPhotos) {
			response.Error(w, r, http.StatusUnprocessableEntity, response.ProblemDetails{
				Title:  "Incomplete Attachment Photos",
				Code:   "ERR_INCOMPLETE_PHOTOS",
				Detail: err.Error(),
			})
			return
		}
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Invalid Attachment Registration",
			Code:   "ERR_INVALID_ATTACHMENT",
			Detail: err.Error(),
		})
		return
	}

	data := generated.AttachmentResponseData{
		AttachmentId:             openapi_types.UUID(att.ID),
		Status:                   generated.AttachmentResponseDataStatus(att.Status),
		AttachmentType:           string(att.AttachmentType),
		Make:                     att.Make,
		Model:                    att.Model,
		StandaloneDailyRateMinor: att.StandaloneDailyRateMinor,
		CreatedAt:                att.CreatedAt,
	}

	response.JSON(w, r, http.StatusCreated, data)
}

// ListAttachmentInventory handles GET /attachments/inventory
func (s *Server) ListAttachmentInventory(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		response.Error(w, r, http.StatusUnauthorized, response.ProblemDetails{
			Title:  "Unauthorized",
			Code:   "ERR_UNAUTHORIZED",
			Detail: "Missing or invalid Bearer authorization token.",
		})
		return
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := s.tokenService.VerifyAccessToken(r.Context(), tokenStr)
	if err != nil {
		response.Error(w, r, http.StatusUnauthorized, response.ProblemDetails{
			Title:  "Unauthorized",
			Code:   "ERR_UNAUTHORIZED",
			Detail: "Access token is expired or signature is invalid.",
		})
		return
	}
	userUUID, err := uuid.Parse(claims.Subject)
	if err != nil {
		response.Error(w, r, http.StatusUnauthorized, response.ProblemDetails{
			Title:  "Unauthorized",
			Code:   "ERR_UNAUTHORIZED",
			Detail: "Invalid token subject claim.",
		})
		return
	}

	statusFilter := r.URL.Query().Get("status")

	items, err := s.attachmentService.ListOwnerInventory(r.Context(), userUUID, statusFilter)
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, response.ProblemDetails{
			Title:  "Internal Server Error",
			Code:   "ERR_INTERNAL",
			Detail: "Failed to list attachment inventory.",
		})
		return
	}

	inventory := make([]generated.AttachmentInventoryItem, 0, len(items))
	for _, it := range items {
		var primaryPhoto *string
		if len(it.Photos) > 0 {
			primaryPhoto = &it.Photos[0].StorageURI
		}
		inventory = append(inventory, generated.AttachmentInventoryItem{
			AttachmentId:             openapi_types.UUID(it.ID),
			AttachmentType:           string(it.AttachmentType),
			Make:                     it.Make,
			Model:                    it.Model,
			SubTypeSize:              it.SubTypeSize,
			Status:                   string(it.Status),
			StandaloneDailyRateMinor: it.StandaloneDailyRateMinor,
			IsLinkedToFleet:          false,
			PrimaryPhotoUri:          primaryPhoto,
		})
	}

	response.JSON(w, r, http.StatusOK, inventory)
}

// LinkFleetAttachment handles POST /fleets/{id}/attachments/link (FL-32, FL-33)
func (s *Server) LinkFleetAttachment(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	fleetUUID, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Invalid Fleet ID",
			Code:   "ERR_INVALID_UUID",
			Detail: "The provided fleet ID is not a valid UUID.",
		})
		return
	}

	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		response.Error(w, r, http.StatusUnauthorized, response.ProblemDetails{
			Title:  "Unauthorized",
			Code:   "ERR_UNAUTHORIZED",
			Detail: "Missing or invalid Bearer authorization token.",
		})
		return
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := s.tokenService.VerifyAccessToken(r.Context(), tokenStr)
	if err != nil {
		response.Error(w, r, http.StatusUnauthorized, response.ProblemDetails{
			Title:  "Unauthorized",
			Code:   "ERR_UNAUTHORIZED",
			Detail: "Access token is expired or signature is invalid.",
		})
		return
	}
	userUUID, err := uuid.Parse(claims.Subject)
	if err != nil {
		response.Error(w, r, http.StatusUnauthorized, response.ProblemDetails{
			Title:  "Unauthorized",
			Code:   "ERR_UNAUTHORIZED",
			Detail: "Invalid token subject claim.",
		})
		return
	}

	var body generated.LinkFleetAttachmentJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Bad Request",
			Code:   "ERR_MALFORMED_JSON",
			Detail: "Request body could not be parsed as valid JSON.",
		})
		return
	}

	attUUID := uuid.UUID(body.AttachmentId)

	link, err := s.attachmentService.LinkAttachmentToFleet(r.Context(), userUUID, machineryApp.LinkAttachmentInput{
		FleetID:               fleetUUID,
		AttachmentID:          attUUID,
		BundledDailyRateMinor: body.BundledDailyRateMinor,
	})
	if err != nil {
		if errors.Is(err, machineryDomain.ErrBundledRateExceedsStandalone) {
			response.Error(w, r, http.StatusUnprocessableEntity, response.ProblemDetails{
				Title:  "Invalid Bundled Attachment Rate",
				Code:   "ERR_BUNDLED_RATE_EXCEEDS_STANDALONE",
				Detail: "The configured bundled daily rate cannot exceed the standalone inventory daily rate.",
			})
			return
		}
		if errors.Is(err, machineryDomain.ErrFleetNotFound) || errors.Is(err, machineryDomain.ErrAttachmentNotFound) {
			response.Error(w, r, http.StatusNotFound, response.ProblemDetails{
				Title:  "Not Found",
				Code:   "ERR_NOT_FOUND",
				Detail: err.Error(),
			})
			return
		}
		if errors.Is(err, machineryDomain.ErrUnauthorizedFleetOwner) || errors.Is(err, machineryDomain.ErrUnauthorizedAttachmentOwner) {
			response.Error(w, r, http.StatusForbidden, response.ProblemDetails{
				Title:  "Forbidden",
				Code:   "ERR_FORBIDDEN",
				Detail: "You do not have permission to link this equipment.",
			})
			return
		}
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Invalid Link Request",
			Code:   "ERR_INVALID_LINK",
			Detail: err.Error(),
		})
		return
	}

	discountPct := float32(0)
	if link.StandaloneDailyRate > 0 {
		discountPct = float32(float64(link.StandaloneDailyRate-link.BundledDailyRateMinor) / float64(link.StandaloneDailyRate) * 100.0)
	}

	data := generated.LinkFleetAttachmentResponseData{
		FleetId:               openapi_types.UUID(link.FleetID),
		AttachmentId:          openapi_types.UUID(link.AttachmentID),
		BundledDailyRateMinor: link.BundledDailyRateMinor,
		DiscountPercentage:    discountPct,
		Status:                "LINKED",
	}

	response.JSON(w, r, http.StatusOK, data)
}

// CheckBundledAvailability handles POST /bookings/check-bundled-availability (FL-34)
func (s *Server) CheckBundledAvailability(w http.ResponseWriter, r *http.Request) {
	var body generated.CheckBundledAvailabilityJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Bad Request",
			Code:   "ERR_MALFORMED_JSON",
			Detail: "Request body could not be parsed as valid JSON.",
		})
		return
	}

	fleetUUID := uuid.UUID(body.FleetId)
	var attUUIDs []uuid.UUID
	for _, a := range body.AttachmentIds {
		attUUIDs = append(attUUIDs, uuid.UUID(a))
	}

	startDate := body.StartDate.Time
	endDate := body.EndDate.Time

	out, err := s.attachmentService.CheckBundledAvailability(r.Context(), machineryApp.CheckBundledAvailabilityInput{
		FleetID:       fleetUUID,
		AttachmentIDs: attUUIDs,
		StartDate:     startDate,
		EndDate:       endDate,
	})
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Availability Check Failed",
			Code:   "ERR_AVAILABILITY_CHECK",
			Detail: err.Error(),
		})
		return
	}

	statuses := make([]generated.AttachmentAvailabilityStatus, len(out.AttachmentStatuses))
	for i, st := range out.AttachmentStatuses {
		statuses[i] = generated.AttachmentAvailabilityStatus{
			AttachmentId:          openapi_types.UUID(st.AttachmentID),
			IsAvailable:           st.IsAvailable,
			DailyRateAppliedMinor: st.DailyRateAppliedMinor,
		}
	}

	data := generated.CheckBundledAvailabilityResponseData{
		FleetAvailable:              out.FleetAvailable,
		AllAttachmentsAvailable:     out.AllAttachmentsAvailable,
		AttachmentStatuses:          statuses,
		TotalCombinedDailyRateMinor: out.TotalCombinedDailyRateMinor,
	}

	response.JSON(w, r, http.StatusOK, data)
}

// ==================== Partition 05-04 Spatial Discovery & Comparison ====================

// SearchFleets handles GET /fleets/search (SRS §05.04, SR-01 to SR-05, SR-07)
func (s *Server) SearchFleets(w http.ResponseWriter, r *http.Request) {
	if s.discoveryService == nil {
		s.notImplemented(w, r, "05-04")
		return
	}

	q := r.URL.Query()
	lat, _ := strconv.ParseFloat(q.Get("latitude"), 64)
	lng, _ := strconv.ParseFloat(q.Get("longitude"), 64)

	filter := machineryDomain.FleetSearchFilter{
		Latitude:  lat,
		Longitude: lng,
	}

	if radStr := q.Get("radiusKm"); radStr != "" {
		if rad, err := strconv.ParseFloat(radStr, 64); err == nil {
			filter.RadiusKm = rad
		}
	}
	if catStr := q.Get("categoryId"); catStr != "" {
		if id, err := uuid.Parse(catStr); err == nil {
			filter.CategoryID = &id
		}
	}
	if subStr := q.Get("subcategoryId"); subStr != "" {
		if id, err := uuid.Parse(subStr); err == nil {
			filter.SubcategoryID = &id
		}
	}
	if startStr := q.Get("startDate"); startStr != "" {
		if t, err := time.Parse("2006-01-02", startStr); err == nil {
			filter.StartDate = &t
		}
	}
	if endStr := q.Get("endDate"); endStr != "" {
		if t, err := time.Parse("2006-01-02", endStr); err == nil {
			filter.EndDate = &t
		}
	}
	if minStr := q.Get("minRateMinor"); minStr != "" {
		if m, err := strconv.ParseInt(minStr, 10, 64); err == nil {
			filter.MinRateMinor = &m
		}
	}
	if maxStr := q.Get("maxRateMinor"); maxStr != "" {
		if m, err := strconv.ParseInt(maxStr, 10, 64); err == nil {
			filter.MaxRateMinor = &m
		}
	}
	if opStr := q.Get("operatorIncluded"); opStr != "" {
		if b, err := strconv.ParseBool(opStr); err == nil {
			filter.OperatorIncluded = &b
		}
	}
	if fuelStr := q.Get("fuelType"); fuelStr != "" {
		filter.FuelType = &fuelStr
	}
	if smStr := q.Get("selfMoveCapable"); smStr != "" {
		if b, err := strconv.ParseBool(smStr); err == nil {
			filter.SelfMoveCapable = &b
		}
	}
	if queryStr := q.Get("query"); queryStr != "" {
		filter.Query = &queryStr
	}
	if sortStr := q.Get("sortBy"); sortStr != "" {
		filter.SortBy = machineryDomain.SortBy(sortStr)
	}
	if cc := r.Header.Get("X-Country-Code"); cc != "" {
		filter.CountryCode = cc
	}
	if limStr := q.Get("limit"); limStr != "" {
		if l, err := strconv.Atoi(limStr); err == nil {
			filter.Limit = l
		}
	}

	out, err := s.discoveryService.SearchFleets(r.Context(), filter)
	if err != nil {
		if errors.Is(err, machineryDomain.ErrSpatialRadiusOutOfBounds) {
			response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
				Title:  "Invalid Search Radius",
				Code:   "ERR_SPATIAL_RADIUS_OUT_OF_BOUNDS",
				Detail: "Search radius must be between 1.0 and 500.0 km.",
			})
			return
		}
		if errors.Is(err, machineryDomain.ErrInvalidRateRange) || errors.Is(err, machineryDomain.ErrInvalidDateRange) {
			response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
				Title:  "Invalid Filter Bounds",
				Code:   "ERR_INVALID_FILTER",
				Detail: err.Error(),
			})
			return
		}

		response.Error(w, r, http.StatusInternalServerError, response.ProblemDetails{
			Title:  "Internal Server Error",
			Code:   "ERR_INTERNAL",
			Detail: "Failed to execute fleet search.",
		})
		return
	}

	items := make([]generated.FleetSearchResultItem, len(out.Results))
	for i, res := range out.Results {
		var specs generated.FleetSearchKeySpecs
		if res.KeySpecs.OperatingWeightKg != nil {
			specs.OperatingWeightKg = res.KeySpecs.OperatingWeightKg
		}
		if res.KeySpecs.PayloadCapacityKg != nil {
			specs.PayloadCapacityKg = res.KeySpecs.PayloadCapacityKg
		}
		if res.KeySpecs.GrossPowerHp != nil {
			specs.GrossPowerHp = res.KeySpecs.GrossPowerHp
		}

		subEn := res.SubcategoryNameEN
		subAm := res.SubcategoryNameAM
		excess := res.ExcessHourRateMinor

		subcatID := openapi_types.UUID(res.SubcategoryID)

		items[i] = generated.FleetSearchResultItem{
			FleetId:             openapi_types.UUID(res.FleetID),
			SubcategoryId:       &subcatID,
			SubcategoryNameEn:   &subEn,
			SubcategoryNameAm:   &subAm,
			Make:                res.Make,
			Model:               res.Model,
			ManufactureYear:     res.ManufactureYear,
			DailyRateMinor:      res.DailyRateMinor,
			ExcessHourRateMinor: &excess,
			DistanceMeters:      res.DistanceMeters,
			AvailabilityStatus:  generated.FleetSearchResultItemAvailabilityStatus(res.AvailabilityStatus),
			SelfMoveCapable:     res.SelfMoveCapable,
			OperatorIncluded:    res.OperatorIncluded,
			FuelType:            res.FuelType,
			PrimaryPhotoUri:     res.PrimaryPhotoURI,
			Owner: generated.FleetSearchOwnerSummary{
				OwnerId:               openapi_types.UUID(res.Owner.OwnerID),
				TradingName:           res.Owner.TradingName,
				AverageRating:         res.Owner.AverageRating,
				TotalCompletedRentals: res.Owner.TotalCompletedRentals,
				IsVerified:            res.Owner.IsVerified,
			},
			KeySpecs: &specs,
		}
	}

	respData := generated.FleetSearchData{
		TotalMatches: out.TotalMatches,
		SearchCenter: generated.FleetSearchCenter{
			Latitude:  out.Latitude,
			Longitude: out.Longitude,
			RadiusKm:  out.RadiusKm,
		},
		Results: items,
	}

	response.JSON(w, r, http.StatusOK, respData)
}

// CompareFleets handles POST /fleets/compare (SRS §05.04, SR-08, FL-15)
func (s *Server) CompareFleets(w http.ResponseWriter, r *http.Request) {
	if s.discoveryService == nil {
		s.notImplemented(w, r, "05-04")
		return
	}

	var req generated.FleetCompareRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Bad Request",
			Code:   "ERR_MALFORMED_JSON",
			Detail: "Request body could not be parsed as valid JSON.",
		})
		return
	}

	if len(req.FleetIds) < 2 || len(req.FleetIds) > 3 {
		response.Error(w, r, http.StatusUnprocessableEntity, response.ProblemDetails{
			Title:  "Invalid Comparison Fleet Count",
			Code:   "ERR_INVALID_COMPARISON_COUNT",
			Detail: "Comparison requires between 2 and 3 fleets.",
		})
		return
	}

	ids := make([]uuid.UUID, len(req.FleetIds))
	for i, id := range req.FleetIds {
		ids[i] = uuid.UUID(id)
	}

	matrix, err := s.discoveryService.CompareFleets(r.Context(), ids)
	if err != nil {
		if errors.Is(err, machineryDomain.ErrNonComparableCategories) {
			response.Error(w, r, http.StatusUnprocessableEntity, response.ProblemDetails{
				Title:  "Taxonomic Category Mismatch in Comparison",
				Code:   "ERR_NON_COMPARABLE_CATEGORIES",
				Detail: "Comparison matrix mandates identical subcategories across all compared assets.",
			})
			return
		}
		if errors.Is(err, machineryDomain.ErrInvalidComparisonCount) {
			response.Error(w, r, http.StatusUnprocessableEntity, response.ProblemDetails{
				Title:  "Invalid Fleet Count",
				Code:   "ERR_INVALID_COMPARISON_COUNT",
				Detail: err.Error(),
			})
			return
		}
		if errors.Is(err, machineryDomain.ErrFleetNotFound) {
			response.Error(w, r, http.StatusNotFound, response.ProblemDetails{
				Title:  "Fleet Not Found",
				Code:   "ERR_FLEET_NOT_FOUND",
				Detail: "One or more fleet assets selected for comparison were not found.",
			})
			return
		}

		response.Error(w, r, http.StatusInternalServerError, response.ProblemDetails{
			Title:  "Internal Server Error",
			Code:   "ERR_INTERNAL",
			Detail: "Failed to build comparison matrix.",
		})
		return
	}

	headers := make([]generated.FleetCompareItemHeader, len(matrix.Fleets))
	for i, f := range matrix.Fleets {
		excess := f.ExcessHourRateMinor
		fuel := f.FuelType
		headers[i] = generated.FleetCompareItemHeader{
			Id:                  openapi_types.UUID(f.ID),
			Title:               f.Title,
			DailyRateMinor:      f.DailyRateMinor,
			ExcessHourRateMinor: &excess,
			PrimaryPhotoUri:     f.PrimaryPhotoURI,
			SelfMoveCapable:     f.SelfMoveCapable,
			OperatorIncluded:    f.OperatorIncluded,
			FuelType:            &fuel,
		}
	}

	rows := make([]generated.FleetCompareSpecRow, len(matrix.ComparisonMatrix))
	for i, row := range matrix.ComparisonMatrix {
		fieldKey := row.FieldKey
		fieldType := row.FieldType
		rows[i] = generated.FleetCompareSpecRow{
			FieldKey:    fieldKey,
			LabelEn:     row.LabelEN,
			LabelAm:     row.LabelAM,
			UnitSymbol:  row.UnitSymbol,
			FieldType:   &fieldType,
			Values:      row.Values,
			IsDifferent: row.IsDifferent,
		}
	}

	data := generated.FleetCompareData{
		SubcategoryId:     openapi_types.UUID(matrix.SubcategoryID),
		SubcategoryNameEn: matrix.SubcategoryNameEN,
		SubcategoryNameAm: matrix.SubcategoryNameAM,
		Fleets:            headers,
		ComparisonMatrix:  rows,
	}

	response.JSON(w, r, http.StatusOK, data)
}

// GetFleetAvailability handles GET /fleets/{id}/availability (SRS §05.04, SR-07)
func (s *Server) GetFleetAvailability(w http.ResponseWriter, r *http.Request) {
	if s.discoveryService == nil {
		s.notImplemented(w, r, "05-04")
		return
	}

	idStr := r.PathValue("id")
	fleetID, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Invalid Fleet ID",
			Code:   "ERR_INVALID_UUID",
			Detail: "The provided fleet ID is not a valid UUID.",
		})
		return
	}

	q := r.URL.Query()
	start, err := time.Parse("2006-01-02", q.Get("startDate"))
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Invalid Start Date",
			Code:   "ERR_INVALID_DATE_FORMAT",
			Detail: "Query parameter 'startDate' must be in YYYY-MM-DD format.",
		})
		return
	}
	end, err := time.Parse("2006-01-02", q.Get("endDate"))
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Invalid End Date",
			Code:   "ERR_INVALID_DATE_FORMAT",
			Detail: "Query parameter 'endDate' must be in YYYY-MM-DD format.",
		})
		return
	}

	info, err := s.discoveryService.GetFleetAvailability(r.Context(), fleetID, start, end)
	if err != nil {
		if errors.Is(err, machineryDomain.ErrFleetNotFound) {
			response.Error(w, r, http.StatusNotFound, response.ProblemDetails{
				Title:  "Fleet Not Found",
				Code:   "ERR_FLEET_NOT_FOUND",
				Detail: "The requested fleet asset was not found.",
			})
			return
		}
		if errors.Is(err, machineryDomain.ErrInvalidDateRange) {
			response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
				Title:  "Invalid Date Range",
				Code:   "ERR_INVALID_DATE_RANGE",
				Detail: "Start date must be before or equal to end date.",
			})
			return
		}

		response.Error(w, r, http.StatusInternalServerError, response.ProblemDetails{
			Title:  "Internal Server Error",
			Code:   "ERR_INTERNAL",
			Detail: "Failed to evaluate fleet availability.",
		})
		return
	}

	var blackouts []generated.FleetAvailabilityDateRange
	if len(info.BlackoutRanges) > 0 {
		blackouts = make([]generated.FleetAvailabilityDateRange, len(info.BlackoutRanges))
		for i, b := range info.BlackoutRanges {
			reason := b.Reason
			blackouts[i] = generated.FleetAvailabilityDateRange{
				StartDate: openapi_types.Date{Time: b.StartDate},
				EndDate:   openapi_types.Date{Time: b.EndDate},
				Reason:    &reason,
			}
		}
	}

	var reserved []generated.FleetAvailabilityDateRange
	if len(info.ReservedRanges) > 0 {
		reserved = make([]generated.FleetAvailabilityDateRange, len(info.ReservedRanges))
		for i, res := range info.ReservedRanges {
			state := res.State
			reserved[i] = generated.FleetAvailabilityDateRange{
				StartDate: openapi_types.Date{Time: res.StartDate},
				EndDate:   openapi_types.Date{Time: res.EndDate},
				State:     &state,
			}
		}
	}

	data := generated.FleetAvailabilityData{
		FleetId:              openapi_types.UUID(info.FleetID),
		IsAvailable:          info.IsAvailable,
		UnavailableReason:    info.UnavailableReason,
		HoldExpiresInSeconds: info.HoldExpiresInSeconds,
		BlackoutRanges:       &blackouts,
		ReservedRanges:       &reserved,
	}

	response.JSON(w, r, http.StatusOK, data)
}

// SearchFleetsByRadius handles legacy GET /geo/fleets/radius-search
func (s *Server) SearchFleetsByRadius(w http.ResponseWriter, r *http.Request) {
	s.SearchFleets(w, r)
}

// ==================== Unimplemented Stubs for Future Partitions ====================

func (s *Server) notImplemented(w http.ResponseWriter, r *http.Request, partition string) {
	response.Error(w, r, http.StatusNotImplemented, response.ProblemDetails{
		Title:  "Not Implemented",
		Code:   "ERR_PARTITION_PENDING",
		Detail: "This endpoint will be activated during partition " + partition + ".",
	})
}

func (s *Server) ListCountryConfigs(w http.ResponseWriter, r *http.Request) { s.notImplemented(w, r, "05-11") }
func (s *Server) GetCountryConfigByCode(w http.ResponseWriter, r *http.Request) { s.notImplemented(w, r, "05-11") }
func (s *Server) UpdateCountryConfig(w http.ResponseWriter, r *http.Request) { s.notImplemented(w, r, "05-11") }
func (s *Server) GetAdminUiConfig(w http.ResponseWriter, r *http.Request) { s.notImplemented(w, r, "05-11") }
func (s *Server) ListBookings(w http.ResponseWriter, r *http.Request) { s.notImplemented(w, r, "05-05") }
func (s *Server) CreateBooking(w http.ResponseWriter, r *http.Request) { s.notImplemented(w, r, "05-05") }
func (s *Server) GetBookingById(w http.ResponseWriter, r *http.Request) { s.notImplemented(w, r, "05-05") }
func (s *Server) SubmitBookingInspection(w http.ResponseWriter, r *http.Request) { s.notImplemented(w, r, "05-07") }
func (s *Server) CancelBooking(w http.ResponseWriter, r *http.Request) { s.notImplemented(w, r, "05-05") }
func (s *Server) ConfirmBookingPayment(w http.ResponseWriter, r *http.Request) { s.notImplemented(w, r, "05-05") }
func (s *Server) InitiatePayment(w http.ResponseWriter, r *http.Request) { s.notImplemented(w, r, "05-06") }
func (s *Server) GetMyWallet(w http.ResponseWriter, r *http.Request) { s.notImplemented(w, r, "05-09") }
func (s *Server) ListWalletTransactions(w http.ResponseWriter, r *http.Request) { s.notImplemented(w, r, "05-09") }

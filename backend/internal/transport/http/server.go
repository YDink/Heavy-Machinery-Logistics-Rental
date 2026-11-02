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

// ==================== Partition 05-02 & 05-04 Fleet & Categories ====================

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
func (s *Server) GetFleetById(w http.ResponseWriter, r *http.Request, id openapi_types.UUID) {
	if s.machineryService == nil {
		s.notImplemented(w, r, "05-02")
		return
	}

	asset, err := s.machineryService.GetFleetAsset(r.Context(), uuid.UUID(id))
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
func (s *Server) CreateFleetAsset(w http.ResponseWriter, r *http.Request, params generated.CreateFleetAssetParams) {
	if s.machineryService == nil {
		s.notImplemented(w, r, "05-02")
		return
	}

	// 1. Authenticate user from Bearer Token
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

	// 2. Decode JSON body with extended fleet properties
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
		OwnerID:             userUUID,
		SubcategoryID:       uuid.UUID(body.CategoryId),
		Make:                body.Make,
		Model:               body.Model,
		ManufactureYear:     body.ManufacturingYear,
		PlateNumber:         body.RegistrationNumber,
		ChassisVIN:          body.ChassisVin,
		BaselineMeterHours:  body.BaselineMeterHours,
		FuelType:            fuel,
		SelfMoveCapable:     body.SelfMoveCapable,
		DailyRateMinor:      body.DailyRate.AmountMinor,
		ExcessHourRateMinor: body.ExcessHourRateMinor,
		MinRentalDays:       body.MinRentalDays,
		PaymentTimeoutHours: body.PaymentTimeoutHours,
		OperatorIncluded:    body.OperatorIncluded,
		Latitude:            body.Location.Latitude,
		Longitude:           body.Location.Longitude,
		AddressText:         body.Location.AddressText,
		Photos:              photos,
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
		Version:            createdAsset.Version,
	}
	dto.Location.Latitude = createdAsset.Latitude
	dto.Location.Longitude = createdAsset.Longitude
	dto.Location.AddressText = createdAsset.AddressText

	response.JSON(w, r, http.StatusCreated, dto)
}

// ListFleets handles GET /fleets
func (s *Server) ListFleets(w http.ResponseWriter, r *http.Request, params generated.ListFleetsParams) {
	if s.machineryService == nil {
		s.notImplemented(w, r, "05-04")
		return
	}
	response.JSON(w, r, http.StatusOK, []generated.FleetAsset{})
}

// ==================== Attachments Handlers (Partition 05-03) ====================

// CreateAttachment handles POST /attachments (ATT-01 to ATT-14)
func (s *Server) CreateAttachment(w http.ResponseWriter, r *http.Request, params generated.CreateAttachmentParams) {
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
func (s *Server) ListAttachmentInventory(w http.ResponseWriter, r *http.Request, params generated.ListAttachmentInventoryParams) {
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

	statusFilter := ""
	if params.Status != nil {
		statusFilter = string(*params.Status)
	}

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

	resp := generated.AttachmentInventoryListResponseData{
		TotalCount:  len(inventory),
		Attachments: inventory,
	}

	response.JSON(w, r, http.StatusOK, resp)
}

// LinkFleetAttachment handles POST /fleets/{id}/attachments/link (FL-32, FL-33)
func (s *Server) LinkFleetAttachment(w http.ResponseWriter, r *http.Request, id openapi_types.UUID) {
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

	fleetUUID := uuid.UUID(id)
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

// ==================== Unimplemented Stubs for Future Partitions ====================

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
// SearchFleets handles GET /fleets/search (SRS §05.04, SR-01 to SR-05, SR-07)
func (s *Server) SearchFleets(w http.ResponseWriter, r *http.Request, params generated.SearchFleetsParams) {
	if s.discoveryService == nil {
		s.notImplemented(w, r, "05-04")
		return
	}

	filter := machineryDomain.FleetSearchFilter{
		Latitude:  params.Latitude,
		Longitude: params.Longitude,
	}

	if params.RadiusKm != nil {
		filter.RadiusKm = float64(*params.RadiusKm)
	}
	if params.CategoryId != nil {
		id := uuid.UUID(*params.CategoryId)
		filter.CategoryID = &id
	}
	if params.SubcategoryId != nil {
		id := uuid.UUID(*params.SubcategoryId)
		filter.SubcategoryID = &id
	}
	if params.StartDate != nil {
		t := params.StartDate.Time
		filter.StartDate = &t
	}
	if params.EndDate != nil {
		t := params.EndDate.Time
		filter.EndDate = &t
	}
	if params.MinRateMinor != nil {
		filter.MinRateMinor = params.MinRateMinor
	}
	if params.MaxRateMinor != nil {
		filter.MaxRateMinor = params.MaxRateMinor
	}
	if params.OperatorIncluded != nil {
		filter.OperatorIncluded = params.OperatorIncluded
	}
	if params.FuelType != nil {
		ft := string(*params.FuelType)
		filter.FuelType = &ft
	}
	if params.SelfMoveCapable != nil {
		filter.SelfMoveCapable = params.SelfMoveCapable
	}
	if params.Query != nil {
		filter.Query = params.Query
	}
	if params.SortBy != nil {
		filter.SortBy = machineryDomain.SortBy(*params.SortBy)
	}
	if params.XCountryCode != nil {
		filter.CountryCode = string(*params.XCountryCode)
	}
	if params.Limit != nil {
		filter.Limit = *params.Limit
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
func (s *Server) GetFleetAvailability(w http.ResponseWriter, r *http.Request, id openapi_types.UUID, params generated.GetFleetAvailabilityParams) {
	if s.discoveryService == nil {
		s.notImplemented(w, r, "05-04")
		return
	}

	fleetID := uuid.UUID(id)
	start := params.StartDate.Time
	end := params.EndDate.Time

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
func (s *Server) SearchFleetsByRadius(w http.ResponseWriter, r *http.Request, params generated.SearchFleetsByRadiusParams) {
	radius := params.RadiusKm
	searchParams := generated.SearchFleetsParams{
		Latitude:  params.Latitude,
		Longitude: params.Longitude,
		RadiusKm:  &radius,
		Limit:     params.Limit,
		Cursor:    params.Cursor,
	}
	s.SearchFleets(w, r, searchParams)
}
func (s *Server) InitiatePayment(w http.ResponseWriter, r *http.Request, params generated.InitiatePaymentParams) { s.notImplemented(w, r, "05-06") }
func (s *Server) GetMyWallet(w http.ResponseWriter, r *http.Request) { s.notImplemented(w, r, "05-09") }
func (s *Server) ListWalletTransactions(w http.ResponseWriter, r *http.Request, params generated.ListWalletTransactionsParams) { s.notImplemented(w, r, "05-09") }

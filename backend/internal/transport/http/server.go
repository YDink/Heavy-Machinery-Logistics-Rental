package http

import (
	"net/http"

	"github.com/dinkmachinery/dink-core/internal/modules/auth/application"
	bookingApp "github.com/dinkmachinery/dink-core/internal/modules/booking/application"
	machineryApp "github.com/dinkmachinery/dink-core/internal/modules/machinery/application"
	rentalApp "github.com/dinkmachinery/dink-core/internal/modules/rental/application"
)

type Server struct {
	authService       *application.AuthService
	tokenService      *application.TokenService
	machineryService  *machineryApp.MachineryService
	attachmentService *machineryApp.AttachmentService
	discoveryService  *machineryApp.DiscoveryService
	bookingService    *bookingApp.BookingService
	rentalService     *rentalApp.RentalService
	swaggerEnabled    bool
}

func NewServer(
	authService *application.AuthService,
	tokenService *application.TokenService,
	machineryService *machineryApp.MachineryService,
	attachmentService *machineryApp.AttachmentService,
	discoveryService *machineryApp.DiscoveryService,
	bookingService *bookingApp.BookingService,
	rentalService *rentalApp.RentalService,
) *Server {
	return &Server{
		authService:       authService,
		tokenService:      tokenService,
		machineryService:  machineryService,
		attachmentService: attachmentService,
		discoveryService:  discoveryService,
		bookingService:    bookingService,
		rentalService:     rentalService,
		swaggerEnabled:    true,
	}
}

// SetSwaggerEnabled explicitly configures whether Swagger UI and OpenAPI documentation routes are active.
func (s *Server) SetSwaggerEnabled(enabled bool) {
	s.swaggerEnabled = enabled
}

// RegisterRoutes registers all API endpoints directly onto the Go standard library http.ServeMux.
func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	// Swagger UI & Interactive Documentation
	if s.swaggerEnabled {
		mux.HandleFunc("GET /swagger/", s.SwaggerUI)
		mux.HandleFunc("GET /swagger/index.html", s.SwaggerUI)
		mux.HandleFunc("GET /swagger/openapi.yaml", s.SwaggerSpec)
		mux.HandleFunc("GET /swagger", s.SwaggerRedirect)
	}
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

	// Booking & Lifecycle State Machine (Partition 05-05)
	mux.HandleFunc("GET /api/v1/bookings", s.ListBookings)
	mux.HandleFunc("POST /api/v1/bookings", s.CreateBooking)
	mux.HandleFunc("GET /api/v1/bookings/{id}", s.GetBookingById)
	mux.HandleFunc("POST /api/v1/bookings/{id}/payment-evidence", s.SubmitBookingPaymentEvidence)
	mux.HandleFunc("POST /api/v1/bookings/{id}/confirm-payment", s.ConfirmBookingPayment)
	mux.HandleFunc("POST /api/v1/bookings/{id}/decline-payment", s.DeclineBookingPayment)
	mux.HandleFunc("POST /api/v1/bookings/{id}/dispute", s.DisputeBookingPayment)
	mux.HandleFunc("POST /api/v1/bookings/{id}/cancel", s.CancelBooking)
	mux.HandleFunc("POST /api/v1/bookings/{id}/handover-check-in", s.HandoverCheckIn)
	mux.HandleFunc("POST /api/v1/bookings/{id}/return-check-out", s.ReturnCheckOut)
	mux.HandleFunc("POST /api/v1/bookings/{id}/inspections", s.SubmitBookingInspection)
	mux.HandleFunc("POST /api/v1/bookings/{action}", s.DispatchBookingAction)
	mux.HandleFunc("POST /api/v1/payments", s.InitiatePayment)
	mux.HandleFunc("GET /api/v1/wallets/me", s.GetMyWallet)
	mux.HandleFunc("GET /api/v1/wallets/me/transactions", s.ListWalletTransactions)

	// Active Rental Management & Shift Telemetry (Partition 05-06)
	mux.HandleFunc("POST /api/v1/bookings/{id}/handover-inspect", s.HandoverInspect)
	mux.HandleFunc("POST /api/v1/bookings/{id}/handover-challenge", s.GenerateHandoverChallenge)
	mux.HandleFunc("POST /api/v1/bookings/{id}/handover-confirm", s.ConfirmHandover)
	mux.HandleFunc("GET /api/v1/bookings/{id}/hud", s.GetActiveRentalHUD)
	mux.HandleFunc("POST /api/v1/bookings/{id}/extend", s.RequestRentalExtension)
	mux.HandleFunc("POST /api/v1/bookings/{id}/breakdown", s.ReportMechanicalBreakdown)
	mux.HandleFunc("POST /api/v1/bookings/{id}/breakdown/resolve", s.ResolveMechanicalBreakdown)

	// Rental Completion, Return Inspection, Damage & Penalty Adjudication (Partition 05-07)
	mux.HandleFunc("POST /api/v1/bookings/{id}/return-initiate", s.InitiateReturn)
	mux.HandleFunc("POST /api/v1/bookings/{id}/return-inspect", s.InspectReturn)
	mux.HandleFunc("POST /api/v1/bookings/{id}/damage-claims", s.LodgeDamageClaim)
	mux.HandleFunc("POST /api/v1/bookings/{id}/damage-claims/{claimId}/respond", s.RespondDamageClaim)
	mux.HandleFunc("POST /api/v1/bookings/{id}/final-reconcile", s.FinalReconcile)
	mux.HandleFunc("POST /api/v1/bookings/{id}/reviews", s.SubmitBookingReview)
}

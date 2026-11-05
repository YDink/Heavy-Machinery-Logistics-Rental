package http

import (
	"encoding/json"
	"errors"
	"net/http"

	rentalApp "github.com/dinkmachinery/dink-core/internal/modules/rental/application"
	rentalDomain "github.com/dinkmachinery/dink-core/internal/modules/rental/domain"
	"github.com/dinkmachinery/dink-core/internal/platform/response"
	"github.com/dinkmachinery/dink-core/internal/transport/http/generated"
	"github.com/google/uuid"
)

// SubmitBookingReview handles POST /api/v1/bookings/{id}/reviews
func (s *Server) SubmitBookingReview(w http.ResponseWriter, r *http.Request) {
	userUUID, err := s.authenticateUser(r)
	if err != nil {
		response.Error(w, r, http.StatusUnauthorized, response.ProblemDetails{
			Title:  "Unauthorized",
			Code:   "ERR_UNAUTHORIZED",
			Detail: err.Error(),
		})
		return
	}

	bookingID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Bad Request",
			Code:   "ERR_INVALID_BOOKING_ID",
			Detail: "Invalid booking ID in path.",
		})
		return
	}

	var req generated.BookingReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, r, http.StatusBadRequest, response.ProblemDetails{
			Title:  "Bad Request",
			Code:   "ERR_MALFORMED_JSON",
			Detail: "Malformed request payload.",
		})
		return
	}

	feedback := ""
	if req.Feedback != nil {
		feedback = *req.Feedback
	}
	dims := map[string]float64{}
	if req.DimensionRatings != nil {
		dims = *req.DimensionRatings
	}

	cmd := rentalApp.SubmitReviewCommand{
		BookingID:        bookingID,
		AuthorID:         userUUID,
		Rating:           req.Rating,
		Feedback:         feedback,
		DimensionRatings: dims,
	}

	review, err := s.rentalService.SubmitBookingReview(r.Context(), cmd)
	if err != nil {
		if errors.Is(err, rentalDomain.ErrSelfReviewForbidden) || errors.Is(err, rentalDomain.ErrInvalidReviewRating) {
			response.Error(w, r, http.StatusUnprocessableEntity, response.ProblemDetails{
				Title:  "Unprocessable Entity",
				Code:   "ERR_REVIEW_VALIDATION_FAILED",
				Detail: err.Error(),
			})
			return
		}
		response.Error(w, r, http.StatusConflict, response.ProblemDetails{
			Title:  "Conflict",
			Code:   "ERR_SUBMIT_REVIEW_FAILED",
			Detail: err.Error(),
		})
		return
	}

	resp := generated.BookingReviewResponseData{
		ReviewId:    review.ID,
		BookingId:   review.BookingID,
		AuthorId:    review.AuthorID,
		TargetId:    review.TargetID,
		Role:        review.Role,
		Rating:      review.Rating,
		Feedback:    &review.Feedback,
		IsPublished: review.IsPublished,
		CreatedAt:   review.CreatedAt,
	}

	response.JSON(w, r, http.StatusCreated, resp)
}

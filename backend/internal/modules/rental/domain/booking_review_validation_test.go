package domain_test

import (
	"testing"

	"github.com/google/uuid"

	"github.com/dinkmachinery/dink-core/internal/modules/rental/domain"
)

func TestValidateBookingReview(t *testing.T) {
	authorID := uuid.New()
	targetID := uuid.New()

	t.Run("valid review passes", func(t *testing.T) {
		if err := domain.ValidateBookingReview(authorID, targetID, 4.8, "RENTER"); err != nil {
			t.Errorf("expected valid review to pass, got: %v", err)
		}
	})

	t.Run("self-review fails with ErrSelfReviewForbidden", func(t *testing.T) {
		if err := domain.ValidateBookingReview(authorID, authorID, 5.0, "OWNER"); err != domain.ErrSelfReviewForbidden {
			t.Fatalf("expected ErrSelfReviewForbidden, got: %v", err)
		}
	})

	t.Run("out-of-range rating fails with ErrInvalidReviewRating", func(t *testing.T) {
		if err := domain.ValidateBookingReview(authorID, targetID, 0.5, "RENTER"); err != domain.ErrInvalidReviewRating {
			t.Errorf("expected ErrInvalidReviewRating for 0.5, got: %v", err)
		}
		if err := domain.ValidateBookingReview(authorID, targetID, 5.5, "RENTER"); err != domain.ErrInvalidReviewRating {
			t.Errorf("expected ErrInvalidReviewRating for 5.5, got: %v", err)
		}
	})
}

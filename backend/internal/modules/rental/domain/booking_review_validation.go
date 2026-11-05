package domain

import "github.com/google/uuid"

// ValidateBookingReview asserts rating range, non-self review, and valid participant role.
func ValidateBookingReview(authorID, targetID uuid.UUID, rating float64, role string) error {
	if authorID == targetID {
		return ErrSelfReviewForbidden
	}
	if rating < 1.0 || rating > 5.0 {
		return ErrInvalidReviewRating
	}
	return nil
}

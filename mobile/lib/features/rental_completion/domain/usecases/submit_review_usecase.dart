import '../entities/booking_review.dart';
import '../repositories/rental_completion_repository.dart';

class SubmitReviewParams {
  final String bookingId;
  final double rating;
  final String? feedback;
  final Map<String, double>? dimensionRatings;

  const SubmitReviewParams({
    required this.bookingId,
    required this.rating,
    this.feedback,
    this.dimensionRatings,
  });
}

class SubmitReviewUseCase {
  final RentalCompletionRepository repository;

  SubmitReviewUseCase(this.repository);

  Future<BookingReview> call(SubmitReviewParams params) {
    return repository.submitReview(
      bookingId: params.bookingId,
      rating: params.rating,
      feedback: params.feedback,
      dimensionRatings: params.dimensionRatings,
    );
  }
}

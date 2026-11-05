import '../../domain/entities/booking_review.dart';

class BookingReviewModel extends BookingReview {
  const BookingReviewModel({
    required super.reviewId,
    required super.bookingId,
    required super.authorId,
    required super.targetId,
    required super.role,
    required super.rating,
    super.feedback,
    required super.isPublished,
    required super.createdAt,
  });

  factory BookingReviewModel.fromJson(Map<String, dynamic> json) {
    return BookingReviewModel(
      reviewId: json['reviewId'] as String? ?? '',
      bookingId: json['bookingId'] as String? ?? '',
      authorId: json['authorId'] as String? ?? '',
      targetId: json['targetId'] as String? ?? '',
      role: json['role'] as String? ?? '',
      rating: (json['rating'] as num?)?.toDouble() ?? 5.0,
      feedback: json['feedback'] as String?,
      isPublished: json['isPublished'] as bool? ?? false,
      createdAt: DateTime.tryParse(json['createdAt'] as String? ?? '') ?? DateTime.now(),
    );
  }
}

import 'package:equatable/equatable.dart';

class BookingReview extends Equatable {
  final String reviewId;
  final String bookingId;
  final String authorId;
  final String targetId;
  final String role;
  final double rating;
  final String? feedback;
  final bool isPublished;
  final DateTime createdAt;

  const BookingReview({
    required this.reviewId,
    required this.bookingId,
    required this.authorId,
    required this.targetId,
    required this.role,
    required this.rating,
    this.feedback,
    required this.isPublished,
    required this.createdAt,
  });

  @override
  List<Object?> get props => [
        reviewId,
        bookingId,
        authorId,
        targetId,
        role,
        rating,
        feedback,
        isPublished,
        createdAt,
      ];
}

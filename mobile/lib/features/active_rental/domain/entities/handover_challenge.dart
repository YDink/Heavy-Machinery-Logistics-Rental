import 'package:equatable/equatable.dart';

class HandoverChallenge extends Equatable {
  final String challengePin;
  final String bookingId;
  final DateTime expiresAt;
  final int maxAttempts;
  final int attemptCount;

  const HandoverChallenge({
    required this.challengePin,
    required this.bookingId,
    required this.expiresAt,
    required this.maxAttempts,
    required this.attemptCount,
  });

  @override
  List<Object?> get props => [challengePin, bookingId, expiresAt, maxAttempts, attemptCount];
}

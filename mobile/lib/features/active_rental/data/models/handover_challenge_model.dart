import '../../domain/entities/handover_challenge.dart';

class HandoverChallengeModel extends HandoverChallenge {
  const HandoverChallengeModel({
    required super.challengePin,
    required super.bookingId,
    required super.expiresAt,
    required super.maxAttempts,
    required super.attemptCount,
  });

  factory HandoverChallengeModel.fromJson(Map<String, dynamic> json) {
    return HandoverChallengeModel(
      challengePin: json['challengePin'] as String? ?? '',
      bookingId: json['bookingId'] as String? ?? '',
      expiresAt: DateTime.tryParse(json['expiresAt'] as String? ?? '') ?? DateTime.now(),
      maxAttempts: (json['maxAttempts'] as num?)?.toInt() ?? 3,
      attemptCount: (json['attemptCount'] as num?)?.toInt() ?? 0,
    );
  }
}

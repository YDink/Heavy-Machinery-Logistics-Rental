import 'package:equatable/equatable.dart';

class FleetSearchOwnerSummary extends Equatable {
  final String ownerId;
  final String tradingName;
  final double averageRating;
  final int totalCompletedRentals;
  final bool isVerified;

  const FleetSearchOwnerSummary({
    required this.ownerId,
    required this.tradingName,
    required this.averageRating,
    required this.totalCompletedRentals,
    required this.isVerified,
  });

  @override
  List<Object?> get props => [
        ownerId,
        tradingName,
        averageRating,
        totalCompletedRentals,
        isVerified,
      ];
}

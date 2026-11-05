import 'package:equatable/equatable.dart';

class DamageClaim extends Equatable {
  final String claimId;
  final String bookingId;
  final String componentCategory;
  final int claimedAmountMinor;
  final int acceptedAmountMinor;
  final String currency;
  final String status;
  final DateTime createdAt;

  const DamageClaim({
    required this.claimId,
    required this.bookingId,
    required this.componentCategory,
    required this.claimedAmountMinor,
    required this.acceptedAmountMinor,
    required this.currency,
    required this.status,
    required this.createdAt,
  });

  @override
  List<Object?> get props => [
        claimId,
        bookingId,
        componentCategory,
        claimedAmountMinor,
        acceptedAmountMinor,
        currency,
        status,
        createdAt,
      ];
}

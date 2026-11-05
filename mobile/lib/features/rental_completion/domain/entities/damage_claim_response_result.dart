import 'package:equatable/equatable.dart';

class DamageClaimResponseResult extends Equatable {
  final String claimId;
  final String bookingId;
  final int acceptedAmountMinor;
  final String status;
  final DateTime updatedAt;

  const DamageClaimResponseResult({
    required this.claimId,
    required this.bookingId,
    required this.acceptedAmountMinor,
    required this.status,
    required this.updatedAt,
  });

  @override
  List<Object?> get props => [
        claimId,
        bookingId,
        acceptedAmountMinor,
        status,
        updatedAt,
      ];
}

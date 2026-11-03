import 'package:equatable/equatable.dart';

class RentalExtension extends Equatable {
  final String bookingId;
  final int extensionDays;
  final int incrementalPrepaymentRequiredMinor;
  final int additionalGrossMinor;
  final String currency;
  final String status;
  final bool ownerApprovalRequired;

  const RentalExtension({
    required this.bookingId,
    required this.extensionDays,
    required this.incrementalPrepaymentRequiredMinor,
    required this.additionalGrossMinor,
    required this.currency,
    required this.status,
    required this.ownerApprovalRequired,
  });

  @override
  List<Object?> get props => [
        bookingId,
        extensionDays,
        incrementalPrepaymentRequiredMinor,
        additionalGrossMinor,
        currency,
        status,
        ownerApprovalRequired,
      ];
}

import 'package:equatable/equatable.dart';

class HandoverInspection extends Equatable {
  final String bookingId;
  final DateTime inspectedAt;
  final bool meterVerified;
  final int photosIngestedCount;
  final String status;

  const HandoverInspection({
    required this.bookingId,
    required this.inspectedAt,
    required this.meterVerified,
    required this.photosIngestedCount,
    required this.status,
  });

  @override
  List<Object?> get props => [bookingId, inspectedAt, meterVerified, photosIngestedCount, status];
}

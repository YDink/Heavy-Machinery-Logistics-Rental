import 'package:equatable/equatable.dart';

class BookingAttachmentEntity extends Equatable {
  final String id;
  final String attachmentId;
  final String? attachmentType;
  final String? make;
  final String? model;
  final int bundledDailyRateMinor;
  final String currency;

  const BookingAttachmentEntity({
    required this.id,
    required this.attachmentId,
    this.attachmentType,
    this.make,
    this.model,
    required this.bundledDailyRateMinor,
    this.currency = 'ETB',
  });

  @override
  List<Object?> get props => [
        id,
        attachmentId,
        attachmentType,
        make,
        model,
        bundledDailyRateMinor,
        currency,
      ];
}

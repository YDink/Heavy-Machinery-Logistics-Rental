import 'package:equatable/equatable.dart';

class BankSlipEntity extends Equatable {
  final String id;
  final String bookingId;
  final String uploaderId;
  final String storageUri;
  final String fileHash;
  final String? transactionReference;
  final int amountMinor;
  final String currency;
  final String status;
  final DateTime createdAt;

  const BankSlipEntity({
    required this.id,
    required this.bookingId,
    required this.uploaderId,
    required this.storageUri,
    required this.fileHash,
    this.transactionReference,
    required this.amountMinor,
    this.currency = 'ETB',
    required this.status,
    required this.createdAt,
  });

  @override
  List<Object?> get props => [
        id,
        bookingId,
        uploaderId,
        storageUri,
        fileHash,
        transactionReference,
        amountMinor,
        currency,
        status,
        createdAt,
      ];
}

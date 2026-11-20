import '../../domain/entities/external_transport.dart';

class ExternalTransportModel {
  final String parentBookingId;
  final String status;
  final String carrierName;
  final String driverFullName;
  final String driverPhoneNumber;
  final String primeMoverPlateNumber;
  final DateTime estimatedArrivalDateTime;

  const ExternalTransportModel({
    required this.parentBookingId,
    required this.status,
    required this.carrierName,
    required this.driverFullName,
    required this.driverPhoneNumber,
    required this.primeMoverPlateNumber,
    required this.estimatedArrivalDateTime,
  });

  factory ExternalTransportModel.fromJson(Map<String, dynamic> json) {
    return ExternalTransportModel(
      parentBookingId: json['parentBookingId'] as String? ?? '',
      status: json['status'] as String? ?? 'READY_FOR_PICKUP',
      carrierName: json['carrierName'] as String? ?? '',
      driverFullName: json['driverFullName'] as String? ?? '',
      driverPhoneNumber: json['driverPhoneNumber'] as String? ?? '',
      primeMoverPlateNumber: json['primeMoverPlateNumber'] as String? ?? '',
      estimatedArrivalDateTime: DateTime.tryParse(
            json['estimatedArrivalDateTime'] as String? ?? '',
          ) ??
          DateTime.now().toUtc(),
    );
  }

  ExternalTransport toEntity() {
    return ExternalTransport(
      parentBookingId: parentBookingId,
      status: status,
      carrierName: carrierName,
      driverFullName: driverFullName,
      driverPhoneNumber: driverPhoneNumber,
      primeMoverPlateNumber: primeMoverPlateNumber,
      estimatedArrivalDateTime: estimatedArrivalDateTime,
    );
  }
}

import '../../domain/entities/matched_carrier.dart';

class MatchedCarrierModel {
  final String fleetId;
  final String make;
  final String model;
  final String ownerTradingName;
  final double payloadCapacityKg;
  final int deckWellLengthMm;
  final int deckWidthMm;
  final double tareWeightKg;
  final int axleCount;
  final int dailyRateMinor;
  final String currency;
  final double distanceFromDepotMeters;
  final double averageRating;

  const MatchedCarrierModel({
    required this.fleetId,
    required this.make,
    required this.model,
    required this.ownerTradingName,
    required this.payloadCapacityKg,
    required this.deckWellLengthMm,
    required this.deckWidthMm,
    required this.tareWeightKg,
    required this.axleCount,
    required this.dailyRateMinor,
    required this.currency,
    required this.distanceFromDepotMeters,
    required this.averageRating,
  });

  factory MatchedCarrierModel.fromJson(Map<String, dynamic> json) {
    return MatchedCarrierModel(
      fleetId: json['fleetId'] as String? ?? '',
      make: json['make'] as String? ?? '',
      model: json['model'] as String? ?? '',
      ownerTradingName: json['ownerTradingName'] as String? ?? '',
      payloadCapacityKg: (json['payloadCapacityKg'] as num?)?.toDouble() ?? 0.0,
      deckWellLengthMm: (json['deckWellLengthMm'] as num?)?.toInt() ?? 0,
      deckWidthMm: (json['deckWidthMm'] as num?)?.toInt() ?? 0,
      tareWeightKg: (json['tareWeightKg'] as num?)?.toDouble() ?? 0.0,
      axleCount: (json['axleCount'] as num?)?.toInt() ?? 3,
      dailyRateMinor: (json['dailyRateMinor'] as num?)?.toInt() ?? 0,
      currency: json['currency'] as String? ?? 'ETB',
      distanceFromDepotMeters: (json['distanceFromDepotMeters'] as num?)?.toDouble() ?? 0.0,
      averageRating: (json['averageRating'] as num?)?.toDouble() ?? 5.0,
    );
  }

  MatchedCarrier toEntity() {
    return MatchedCarrier(
      fleetId: fleetId,
      make: make,
      model: model,
      ownerTradingName: ownerTradingName,
      payloadCapacityKg: payloadCapacityKg,
      deckWellLengthMm: deckWellLengthMm,
      deckWidthMm: deckWidthMm,
      tareWeightKg: tareWeightKg,
      axleCount: axleCount,
      dailyRateMinor: dailyRateMinor,
      currency: currency,
      distanceFromDepotMeters: distanceFromDepotMeters,
      averageRating: averageRating,
    );
  }
}

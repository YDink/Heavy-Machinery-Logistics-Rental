class MatchedCarrier {
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

  const MatchedCarrier({
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
}

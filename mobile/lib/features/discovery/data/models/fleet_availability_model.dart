import 'package:dink_mobile/features/discovery/domain/entities/fleet_availability.dart';

class FleetAvailabilityDataModel extends FleetAvailabilityData {
  const FleetAvailabilityDataModel({
    required super.fleetId,
    required super.isAvailable,
    super.unavailableReason,
    super.holdExpiresInSeconds,
    required super.blackoutRanges,
    required super.reservedRanges,
  });

  factory FleetAvailabilityDataModel.fromJson(Map<String, dynamic> json) {
    final rawBlackouts = json['blackoutRanges'] as List<dynamic>? ?? [];
    final blackouts = rawBlackouts.map((e) {
      final m = e as Map<String, dynamic>;
      return FleetAvailabilityDateRange(
        startDate: DateTime.tryParse(m['startDate'] as String? ?? '') ?? DateTime.now(),
        endDate: DateTime.tryParse(m['endDate'] as String? ?? '') ?? DateTime.now(),
        reason: m['reason'] as String?,
        state: m['state'] as String?,
      );
    }).toList();

    final rawReserved = json['reservedRanges'] as List<dynamic>? ?? [];
    final reserved = rawReserved.map((e) {
      final m = e as Map<String, dynamic>;
      return FleetAvailabilityDateRange(
        startDate: DateTime.tryParse(m['startDate'] as String? ?? '') ?? DateTime.now(),
        endDate: DateTime.tryParse(m['endDate'] as String? ?? '') ?? DateTime.now(),
        reason: m['reason'] as String?,
        state: m['state'] as String?,
      );
    }).toList();

    return FleetAvailabilityDataModel(
      fleetId: json['fleetId'] as String? ?? '',
      isAvailable: json['isAvailable'] as bool? ?? true,
      unavailableReason: json['unavailableReason'] as String?,
      holdExpiresInSeconds: (json['holdExpiresInSeconds'] as num?)?.toInt(),
      blackoutRanges: blackouts,
      reservedRanges: reserved,
    );
  }
}

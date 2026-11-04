import '../../domain/entities/active_rental_hud.dart';

class ActiveRentalHUDModel extends ActiveRentalHUD {
  const ActiveRentalHUDModel({
    required super.bookingId,
    required super.status,
    required super.fleetId,
    required super.fleetTitle,
    required super.plateNumber,
    super.operator,
    required super.currentRentalDay,
    required super.totalRentalDays,
    required super.handoverBaselineHours,
    required super.accumulatedHoursToday,
    required super.allowedHoursPerDay,
    required super.excessHoursTotal,
    required super.shiftTimerPaused,
    required super.canExtend,
    required super.canReportBreakdown,
    required super.canInitiateReturn,
  });

  factory ActiveRentalHUDModel.fromJson(Map<String, dynamic> json) {
    final fleet = json['fleet'] as Map<String, dynamic>? ?? {};
    final operatorMap = json['operator'] as Map<String, dynamic>?;
    final quickActions = json['quickActionsEnabled'] as Map<String, dynamic>? ?? {};
    final shiftTelemetry = json['shiftTelemetry'] as Map<String, dynamic>? ?? {};

    HUDOperator? op;
    if (operatorMap != null && operatorMap['fullName'] != null) {
      op = HUDOperator(
        id: operatorMap['id'] as String?,
        fullName: operatorMap['fullName'] as String? ?? '',
        phoneNumber: operatorMap['phoneNumber'] as String? ?? '',
        licenseClass: operatorMap['licenseClass'] as String? ?? '',
      );
    }

    return ActiveRentalHUDModel(
      bookingId: json['bookingId'] as String? ?? '',
      status: json['status'] as String? ?? '',
      fleetId: fleet['id'] as String? ?? '',
      fleetTitle: fleet['title'] as String? ?? '',
      plateNumber: fleet['plateNumber'] as String? ?? '',
      operator: op,
      currentRentalDay: (shiftTelemetry['currentRentalDay'] as num?)?.toInt() ?? 1,
      totalRentalDays: (shiftTelemetry['totalRentalDays'] as num?)?.toInt() ?? 1,
      handoverBaselineHours: (shiftTelemetry['handoverBaselineHours'] as num?)?.toDouble() ?? 0.0,
      accumulatedHoursToday: (shiftTelemetry['accumulatedHoursToday'] as num?)?.toDouble() ?? 0.0,
      allowedHoursPerDay: (shiftTelemetry['allowedHoursPerDay'] as num?)?.toDouble() ?? 8.0,
      excessHoursTotal: (shiftTelemetry['excessHoursTotal'] as num?)?.toDouble() ?? 0.0,
      shiftTimerPaused: shiftTelemetry['shiftTimerPaused'] as bool? ?? false,
      canExtend: quickActions['canExtend'] as bool? ?? false,
      canReportBreakdown: quickActions['canReportBreakdown'] as bool? ?? false,
      canInitiateReturn: quickActions['canInitiateReturn'] as bool? ?? false,
    );
  }
}

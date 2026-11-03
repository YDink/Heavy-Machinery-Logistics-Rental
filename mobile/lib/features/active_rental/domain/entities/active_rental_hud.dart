import 'package:equatable/equatable.dart';

class ActiveRentalHUD extends Equatable {
  final String bookingId;
  final String status;
  final String fleetId;
  final String fleetTitle;
  final String plateNumber;
  final HUDOperator? operator;
  final int currentRentalDay;
  final int totalRentalDays;
  final double handoverBaselineHours;
  final double accumulatedHoursToday;
  final double allowedHoursPerDay;
  final double excessHoursTotal;
  final bool shiftTimerPaused;
  final bool canExtend;
  final bool canReportBreakdown;
  final bool canInitiateReturn;

  const ActiveRentalHUD({
    required this.bookingId,
    required this.status,
    required this.fleetId,
    required this.fleetTitle,
    required this.plateNumber,
    this.operator,
    required this.currentRentalDay,
    required this.totalRentalDays,
    required this.handoverBaselineHours,
    required this.accumulatedHoursToday,
    required this.allowedHoursPerDay,
    required this.excessHoursTotal,
    required this.shiftTimerPaused,
    required this.canExtend,
    required this.canReportBreakdown,
    required this.canInitiateReturn,
  });

  double get shiftPacingRatio =>
      allowedHoursPerDay > 0 ? (accumulatedHoursToday / allowedHoursPerDay) : 0.0;

  @override
  List<Object?> get props => [
        bookingId,
        status,
        fleetId,
        fleetTitle,
        plateNumber,
        operator,
        currentRentalDay,
        totalRentalDays,
        handoverBaselineHours,
        accumulatedHoursToday,
        allowedHoursPerDay,
        excessHoursTotal,
        shiftTimerPaused,
        canExtend,
        canReportBreakdown,
        canInitiateReturn,
      ];
}

class HUDOperator extends Equatable {
  final String? id;
  final String fullName;
  final String phoneNumber;
  final String licenseClass;

  const HUDOperator({
    this.id,
    required this.fullName,
    required this.phoneNumber,
    required this.licenseClass,
  });

  @override
  List<Object?> get props => [id, fullName, phoneNumber, licenseClass];
}

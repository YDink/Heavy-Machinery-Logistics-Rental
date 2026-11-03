import '../entities/mechanical_breakdown.dart';
import '../repositories/active_rental_repository.dart';

class ResolveBreakdownUseCase {
  final ActiveRentalRepository repository;

  ResolveBreakdownUseCase(this.repository);

  Future<MechanicalBreakdown> call({
    required String bookingId,
    required String mechanicName,
    required String mechanicPhone,
    required String resolutionNotes,
  }) {
    return repository.resolveBreakdown(
      bookingId: bookingId,
      mechanicName: mechanicName,
      mechanicPhone: mechanicPhone,
      resolutionNotes: resolutionNotes,
    );
  }
}

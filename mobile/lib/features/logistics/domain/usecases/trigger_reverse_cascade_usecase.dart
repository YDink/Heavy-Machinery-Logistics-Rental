import '../entities/reverse_cascade_result.dart';
import '../repositories/logistics_repository.dart';

class TriggerReverseCascadeUseCase {
  final LogisticsRepository repository;

  TriggerReverseCascadeUseCase(this.repository);

  Future<ReverseCascadeResult> call({
    required String parentBookingId,
    String? reason,
  }) {
    return repository.triggerReverseCascade(
      parentBookingId: parentBookingId,
      reason: reason,
    );
  }
}

import '../entities/logistics_dispatch_result.dart';
import '../entities/winching_photo.dart';
import '../repositories/logistics_repository.dart';

class AdvanceLogisticsDispatchUseCase {
  final LogisticsRepository repository;

  AdvanceLogisticsDispatchUseCase(this.repository);

  Future<LogisticsDispatchResult> call({
    required String subBookingId,
    required String targetStage,
    List<WinchingPhoto> winchingPhotos = const [],
    String? notes,
  }) {
    return repository.advanceDispatchStage(
      subBookingId: subBookingId,
      targetStage: targetStage,
      winchingPhotos: winchingPhotos,
      notes: notes,
    );
  }
}

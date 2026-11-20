import '../entities/logistics_match_result.dart';
import '../repositories/logistics_repository.dart';

class MatchLogisticsCarriersUseCase {
  final LogisticsRepository repository;

  MatchLogisticsCarriersUseCase(this.repository);

  Future<LogisticsMatchResult> call({
    required String bookingId,
    required int searchRadiusKm,
    required DateTime pickupDate,
  }) {
    return repository.matchCarriers(
      bookingId: bookingId,
      searchRadiusKm: searchRadiusKm,
      pickupDate: pickupDate,
    );
  }
}

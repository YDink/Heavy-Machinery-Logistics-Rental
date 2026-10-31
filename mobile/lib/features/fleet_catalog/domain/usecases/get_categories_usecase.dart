import 'package:dink_mobile/features/fleet_catalog/domain/entities/category.dart';
import 'package:dink_mobile/features/fleet_catalog/domain/repositories/fleet_repository.dart';

class GetCategoriesUseCase {
  final FleetRepository repository;

  GetCategoriesUseCase(this.repository);

  Future<List<Category>> call() {
    return repository.getCategories();
  }
}

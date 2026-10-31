import 'package:dink_mobile/features/attachment/domain/entities/attachment_inventory_item.dart';
import 'package:dink_mobile/features/attachment/domain/repositories/attachment_repository.dart';

class GetAttachmentInventoryUseCase {
  final AttachmentRepository repository;

  GetAttachmentInventoryUseCase(this.repository);

  Future<List<AttachmentInventoryItem>> call({String? status}) {
    return repository.getInventory(status: status);
  }
}

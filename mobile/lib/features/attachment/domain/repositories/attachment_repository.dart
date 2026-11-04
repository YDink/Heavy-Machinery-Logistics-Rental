import '../entities/attachment.dart';
import '../entities/attachment_inventory_item.dart';
import '../entities/attachment_registration_request.dart';
import '../entities/bundled_attachment_link.dart';

abstract class AttachmentRepository {
  Future<Attachment> registerAttachment(AttachmentRegistrationRequest request);
  Future<List<AttachmentInventoryItem>> getInventory({String? status});
  Future<BundledAttachmentLink> linkAttachmentToFleet({
    required String fleetId,
    required String attachmentId,
    required int bundledDailyRateMinor,
  });
  Future<void> saveDraft(AttachmentRegistrationRequest request);
  Future<AttachmentRegistrationRequest?> getDraft();
  Future<void> clearDraft();
}

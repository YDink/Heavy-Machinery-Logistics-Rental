import 'package:dink_mobile/features/attachment/domain/entities/bundled_attachment_link.dart';
import 'package:dink_mobile/features/attachment/domain/repositories/attachment_repository.dart';

class LinkAttachmentUseCase {
  final AttachmentRepository repository;

  LinkAttachmentUseCase(this.repository);

  Future<BundledAttachmentLink> call({
    required String fleetId,
    required String attachmentId,
    required int bundledDailyRateMinor,
  }) {
    return repository.linkAttachmentToFleet(
      fleetId: fleetId,
      attachmentId: attachmentId,
      bundledDailyRateMinor: bundledDailyRateMinor,
    );
  }
}

import 'package:dink_mobile/features/attachment/domain/entities/attachment.dart';
import 'package:dink_mobile/features/attachment/domain/repositories/attachment_repository.dart';

class RegisterAttachmentUseCase {
  final AttachmentRepository repository;

  RegisterAttachmentUseCase(this.repository);

  Future<Attachment> call(AttachmentRegistrationRequest request) {
    return repository.registerAttachment(request);
  }
}

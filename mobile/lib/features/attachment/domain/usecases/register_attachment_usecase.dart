import '../entities/attachment.dart';
import '../entities/attachment_registration_request.dart';
import '../repositories/attachment_repository.dart';

class RegisterAttachmentUseCase {
  final AttachmentRepository repository;

  RegisterAttachmentUseCase(this.repository);

  Future<Attachment> call(AttachmentRegistrationRequest request) {
    return repository.registerAttachment(request);
  }
}

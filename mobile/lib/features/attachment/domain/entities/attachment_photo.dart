import 'attachment_photo_angle.dart';

class AttachmentPhoto {
  final AttachmentPhotoAngle angle;
  final String storageUri;
  final String? localFilePath;

  const AttachmentPhoto({
    required this.angle,
    required this.storageUri,
    this.localFilePath,
  });
}

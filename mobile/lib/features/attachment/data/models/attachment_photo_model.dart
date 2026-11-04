import '../../domain/entities/attachment_photo.dart';
import '../../domain/entities/attachment_photo_angle.dart';

class AttachmentPhotoModel extends AttachmentPhoto {
  const AttachmentPhotoModel({
    required super.angle,
    required super.storageUri,
    super.localFilePath,
  });

  factory AttachmentPhotoModel.fromJson(Map<String, dynamic> json) {
    return AttachmentPhotoModel(
      angle: AttachmentPhotoAngle.fromCode(json['angle'] as String? ?? ''),
      storageUri: json['storageUri'] as String? ?? '',
      localFilePath: json['localFilePath'] as String?,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'angle': angle.code,
      'storageUri': storageUri,
      if (localFilePath != null) 'localFilePath': localFilePath,
    };
  }

  factory AttachmentPhotoModel.fromEntity(AttachmentPhoto entity) {
    return AttachmentPhotoModel(
      angle: entity.angle,
      storageUri: entity.storageUri,
      localFilePath: entity.localFilePath,
    );
  }
}

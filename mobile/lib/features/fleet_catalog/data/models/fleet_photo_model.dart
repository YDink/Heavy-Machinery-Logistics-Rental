import '../../domain/entities/fleet_photo.dart';
import '../../domain/entities/photo_angle.dart';

class FleetPhotoModel extends FleetPhoto {
  const FleetPhotoModel({
    required super.angle,
    required super.storageUri,
    required super.fileHash,
    super.localFilePath,
  });

  factory FleetPhotoModel.fromJson(Map<String, dynamic> json) {
    return FleetPhotoModel(
      angle: PhotoAngle.fromCode(json['angle'] as String? ?? '') ?? PhotoAngle.front,
      storageUri: json['storageUri'] as String? ?? '',
      fileHash: json['fileHash'] as String? ?? '',
      localFilePath: json['localFilePath'] as String?,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'angle': angle.code,
      'storageUri': storageUri,
      'fileHash': fileHash,
      if (localFilePath != null) 'localFilePath': localFilePath,
    };
  }

  factory FleetPhotoModel.fromEntity(FleetPhoto entity) {
    return FleetPhotoModel(
      angle: entity.angle,
      storageUri: entity.storageUri,
      fileHash: entity.fileHash,
      localFilePath: entity.localFilePath,
    );
  }
}

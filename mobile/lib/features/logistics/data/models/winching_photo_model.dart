import '../../domain/entities/winching_photo.dart';

class WinchingPhotoModel {
  final String angle;
  final String storageUri;

  const WinchingPhotoModel({
    required this.angle,
    required this.storageUri,
  });

  factory WinchingPhotoModel.fromEntity(WinchingPhoto entity) {
    return WinchingPhotoModel(
      angle: entity.angle,
      storageUri: entity.storageUri,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'angle': angle,
      'storageUri': storageUri,
    };
  }
}

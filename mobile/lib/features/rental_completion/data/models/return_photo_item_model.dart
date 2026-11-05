import '../../domain/entities/return_photo_item.dart';

class ReturnPhotoItemModel extends ReturnPhotoItem {
  const ReturnPhotoItemModel({
    required super.angle,
    required super.storageUri,
  });

  factory ReturnPhotoItemModel.fromJson(Map<String, dynamic> json) {
    return ReturnPhotoItemModel(
      angle: json['angle'] as String? ?? '',
      storageUri: json['storageUri'] as String? ?? '',
    );
  }

  factory ReturnPhotoItemModel.fromEntity(ReturnPhotoItem entity) {
    return ReturnPhotoItemModel(
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

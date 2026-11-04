import '../../domain/entities/handover_photo_item.dart';

class HandoverPhotoItemModel extends HandoverPhotoItem {
  const HandoverPhotoItemModel({
    required super.angle,
    required super.storageUri,
  });

  factory HandoverPhotoItemModel.fromJson(Map<String, dynamic> json) {
    return HandoverPhotoItemModel(
      angle: json['angle'] as String? ?? '',
      storageUri: json['storageUri'] as String? ?? '',
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'angle': angle,
      'storageUri': storageUri,
    };
  }

  factory HandoverPhotoItemModel.fromEntity(HandoverPhotoItem entity) {
    return HandoverPhotoItemModel(
      angle: entity.angle,
      storageUri: entity.storageUri,
    );
  }
}

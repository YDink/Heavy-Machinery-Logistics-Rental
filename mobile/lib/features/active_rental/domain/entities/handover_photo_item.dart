import 'package:equatable/equatable.dart';

class HandoverPhotoItem extends Equatable {
  final String angle;
  final String storageUri;

  const HandoverPhotoItem({
    required this.angle,
    required this.storageUri,
  });

  @override
  List<Object?> get props => [angle, storageUri];
}

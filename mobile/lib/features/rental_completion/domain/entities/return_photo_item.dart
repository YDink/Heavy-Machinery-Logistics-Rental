import 'package:equatable/equatable.dart';

class ReturnPhotoItem extends Equatable {
  final String angle;
  final String storageUri;

  const ReturnPhotoItem({
    required this.angle,
    required this.storageUri,
  });

  @override
  List<Object?> get props => [angle, storageUri];
}

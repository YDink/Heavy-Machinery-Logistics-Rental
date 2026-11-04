import 'package:equatable/equatable.dart';
import 'photo_angle.dart';

class FleetPhoto extends Equatable {
  final PhotoAngle angle;
  final String storageUri;
  final String fileHash;
  final String? localFilePath;

  const FleetPhoto({
    required this.angle,
    required this.storageUri,
    required this.fileHash,
    this.localFilePath,
  });

  @override
  List<Object?> get props => [angle, storageUri, fileHash, localFilePath];
}

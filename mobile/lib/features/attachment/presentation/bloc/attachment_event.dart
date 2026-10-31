import 'package:equatable/equatable.dart';
import 'package:dink_mobile/features/attachment/domain/entities/attachment.dart';

abstract class AttachmentEvent extends Equatable {
  const AttachmentEvent();

  @override
  List<Object?> get props => [];
}

class LoadAttachmentDraftEvent extends AttachmentEvent {
  const LoadAttachmentDraftEvent();
}

class SaveAttachmentDraftEvent extends AttachmentEvent {
  final AttachmentRegistrationRequest draft;
  const SaveAttachmentDraftEvent(this.draft);

  @override
  List<Object?> get props => [draft];
}

class PhotoCapturedEvent extends AttachmentEvent {
  final AttachmentPhotoAngle angle;
  final String localFilePath;
  final String storageUri;

  const PhotoCapturedEvent({
    required this.angle,
    required this.localFilePath,
    required this.storageUri,
  });

  @override
  List<Object?> get props => [angle, localFilePath, storageUri];
}

class SubmitAttachmentRegistrationEvent extends AttachmentEvent {
  final AttachmentRegistrationRequest request;
  const SubmitAttachmentRegistrationEvent(this.request);

  @override
  List<Object?> get props => [request];
}

class LoadAttachmentInventoryEvent extends AttachmentEvent {
  final String? status;
  const LoadAttachmentInventoryEvent({this.status});

  @override
  List<Object?> get props => [status];
}

class LinkAttachmentToFleetEvent extends AttachmentEvent {
  final String fleetId;
  final String attachmentId;
  final int bundledDailyRateMinor;

  const LinkAttachmentToFleetEvent({
    required this.fleetId,
    required this.attachmentId,
    required this.bundledDailyRateMinor,
  });

  @override
  List<Object?> get props => [fleetId, attachmentId, bundledDailyRateMinor];
}

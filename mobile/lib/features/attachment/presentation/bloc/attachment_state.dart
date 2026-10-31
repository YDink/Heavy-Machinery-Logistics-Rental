import 'package:equatable/equatable.dart';
import 'package:dink_mobile/features/attachment/domain/entities/attachment.dart';
import 'package:dink_mobile/features/attachment/domain/entities/attachment_inventory_item.dart';
import 'package:dink_mobile/features/attachment/domain/entities/bundled_attachment_link.dart';

enum AttachmentStatusState { initial, loading, draftLoaded, draftSaved, registered, loaded, linked, error }

class AttachmentState extends Equatable {
  final AttachmentStatusState status;
  final AttachmentRegistrationRequest? draft;
  final Map<AttachmentPhotoAngle, AttachmentPhoto> capturedPhotos;
  final List<AttachmentInventoryItem> inventory;
  final BundledAttachmentLink? lastLinked;
  final String? errorMessage;

  const AttachmentState({
    this.status = AttachmentStatusState.initial,
    this.draft,
    this.capturedPhotos = const {},
    this.inventory = const [],
    this.lastLinked,
    this.errorMessage,
  });

  bool get areAllPhotosCaptured =>
      capturedPhotos.containsKey(AttachmentPhotoAngle.overall) &&
      capturedPhotos.containsKey(AttachmentPhotoAngle.bracket) &&
      capturedPhotos.containsKey(AttachmentPhotoAngle.pins) &&
      capturedPhotos.containsKey(AttachmentPhotoAngle.toolTip);

  AttachmentState copyWith({
    AttachmentStatusState? status,
    AttachmentRegistrationRequest? draft,
    Map<AttachmentPhotoAngle, AttachmentPhoto>? capturedPhotos,
    List<AttachmentInventoryItem>? inventory,
    BundledAttachmentLink? lastLinked,
    String? errorMessage,
  }) {
    return AttachmentState(
      status: status ?? this.status,
      draft: draft ?? this.draft,
      capturedPhotos: capturedPhotos ?? this.capturedPhotos,
      inventory: inventory ?? this.inventory,
      lastLinked: lastLinked ?? this.lastLinked,
      errorMessage: errorMessage ?? this.errorMessage,
    );
  }

  @override
  List<Object?> get props => [
        status,
        draft,
        capturedPhotos,
        inventory,
        lastLinked,
        errorMessage,
      ];
}

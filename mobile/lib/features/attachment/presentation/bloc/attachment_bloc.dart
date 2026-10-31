import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:dink_mobile/features/attachment/domain/entities/attachment.dart';
import 'package:dink_mobile/features/attachment/domain/usecases/get_attachment_inventory_usecase.dart';
import 'package:dink_mobile/features/attachment/domain/usecases/link_attachment_usecase.dart';
import 'package:dink_mobile/features/attachment/domain/usecases/register_attachment_usecase.dart';
import 'package:dink_mobile/features/attachment/domain/repositories/attachment_repository.dart';
import 'package:dink_mobile/features/attachment/presentation/bloc/attachment_event.dart';
import 'package:dink_mobile/features/attachment/presentation/bloc/attachment_state.dart';

class AttachmentBloc extends Bloc<AttachmentEvent, AttachmentState> {
  final RegisterAttachmentUseCase registerAttachmentUseCase;
  final GetAttachmentInventoryUseCase getInventoryUseCase;
  final LinkAttachmentUseCase linkAttachmentUseCase;
  final AttachmentRepository attachmentRepository;

  AttachmentBloc({
    required this.registerAttachmentUseCase,
    required this.getInventoryUseCase,
    required this.linkAttachmentUseCase,
    required this.attachmentRepository,
  }) : super(const AttachmentState()) {
    on<LoadAttachmentDraftEvent>(_onLoadDraft);
    on<SaveAttachmentDraftEvent>(_onSaveDraft);
    on<PhotoCapturedEvent>(_onPhotoCaptured);
    on<SubmitAttachmentRegistrationEvent>(_onSubmitRegistration);
    on<LoadAttachmentInventoryEvent>(_onLoadInventory);
    on<LinkAttachmentToFleetEvent>(_onLinkAttachment);
  }

  Future<void> _onLoadDraft(
    LoadAttachmentDraftEvent event,
    Emitter<AttachmentState> emit,
  ) async {
    final draft = await attachmentRepository.getDraft();
    if (draft != null) {
      final photoMap = <AttachmentPhotoAngle, AttachmentPhoto>{};
      for (final p in draft.photos) {
        photoMap[p.angle] = p;
      }
      emit(state.copyWith(
        status: AttachmentStatusState.draftLoaded,
        draft: draft,
        capturedPhotos: photoMap,
      ));
    }
  }

  Future<void> _onSaveDraft(
    SaveAttachmentDraftEvent event,
    Emitter<AttachmentState> emit,
  ) async {
    await attachmentRepository.saveDraft(event.draft);
    emit(state.copyWith(
      status: AttachmentStatusState.draftSaved,
      draft: event.draft,
    ));
  }

  void _onPhotoCaptured(
    PhotoCapturedEvent event,
    Emitter<AttachmentState> emit,
  ) {
    final updated = Map<AttachmentPhotoAngle, AttachmentPhoto>.from(state.capturedPhotos);
    updated[event.angle] = AttachmentPhoto(
      angle: event.angle,
      storageUri: event.storageUri,
      localFilePath: event.localFilePath,
    );
    emit(state.copyWith(capturedPhotos: updated));
  }

  Future<void> _onSubmitRegistration(
    SubmitAttachmentRegistrationEvent event,
    Emitter<AttachmentState> emit,
  ) async {
    emit(state.copyWith(status: AttachmentStatusState.loading));
    try {
      await registerAttachmentUseCase(event.request);
      emit(state.copyWith(
        status: AttachmentStatusState.registered,
        capturedPhotos: {},
      ));
    } catch (e) {
      emit(state.copyWith(
        status: AttachmentStatusState.error,
        errorMessage: e.toString(),
      ));
    }
  }

  Future<void> _onLoadInventory(
    LoadAttachmentInventoryEvent event,
    Emitter<AttachmentState> emit,
  ) async {
    emit(state.copyWith(status: AttachmentStatusState.loading));
    try {
      final items = await getInventoryUseCase(status: event.status);
      emit(state.copyWith(
        status: AttachmentStatusState.loaded,
        inventory: items,
      ));
    } catch (e) {
      emit(state.copyWith(
        status: AttachmentStatusState.error,
        errorMessage: e.toString(),
      ));
    }
  }

  Future<void> _onLinkAttachment(
    LinkAttachmentToFleetEvent event,
    Emitter<AttachmentState> emit,
  ) async {
    emit(state.copyWith(status: AttachmentStatusState.loading));
    try {
      final link = await linkAttachmentUseCase(
        fleetId: event.fleetId,
        attachmentId: event.attachmentId,
        bundledDailyRateMinor: event.bundledDailyRateMinor,
      );
      emit(state.copyWith(
        status: AttachmentStatusState.linked,
        lastLinked: link,
      ));
    } catch (e) {
      emit(state.copyWith(
        status: AttachmentStatusState.error,
        errorMessage: e.toString(),
      ));
    }
  }
}

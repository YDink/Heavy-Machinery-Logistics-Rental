import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:dink_mobile/core/errors/failure.dart';
import 'package:dink_mobile/features/fleet_catalog/domain/entities/fleet_asset.dart';
import 'package:dink_mobile/features/fleet_catalog/domain/repositories/fleet_repository.dart';
import 'package:dink_mobile/features/fleet_catalog/domain/usecases/get_categories_usecase.dart';
import 'package:dink_mobile/features/fleet_catalog/domain/usecases/get_spec_definitions_usecase.dart';
import 'package:dink_mobile/features/fleet_catalog/domain/usecases/register_fleet_usecase.dart';
import 'package:dink_mobile/features/fleet_catalog/presentation/bloc/fleet_registration_event.dart';
import 'package:dink_mobile/features/fleet_catalog/presentation/bloc/fleet_registration_state.dart';

class FleetRegistrationBloc
    extends Bloc<FleetRegistrationEvent, FleetRegistrationState> {
  final GetCategoriesUseCase getCategoriesUseCase;
  final GetSpecDefinitionsUseCase getSpecDefinitionsUseCase;
  final RegisterFleetUseCase registerFleetUseCase;
  final FleetRepository fleetRepository;

  FleetRegistrationBloc({
    required this.getCategoriesUseCase,
    required this.getSpecDefinitionsUseCase,
    required this.registerFleetUseCase,
    required this.fleetRepository,
  }) : super(const FleetRegistrationInitial()) {
    on<LoadCategoriesEvent>(_onLoadCategories);
    on<SelectCategoryEvent>(_onSelectCategory);
    on<UpdatePhotoEvent>(_onUpdatePhoto);
    on<UpdateSpecValueEvent>(_onUpdateSpecValue);
    on<SaveDraftEvent>(_onSaveDraft);
    on<SubmitFleetRegistrationEvent>(_onSubmitFleetRegistration);
  }

  Future<void> _onLoadCategories(
    LoadCategoriesEvent event,
    Emitter<FleetRegistrationState> emit,
  ) async {
    emit(const CategoriesLoading());
    try {
      final categories = await getCategoriesUseCase();
      emit(CategoriesLoaded(categories: categories));
    } on Failure catch (f) {
      emit(FleetSubmitFailure(message: f.message, code: f.code));
    } catch (e) {
      emit(FleetSubmitFailure(message: e.toString(), code: 'ERR_UNKNOWN'));
    }
  }

  Future<void> _onSelectCategory(
    SelectCategoryEvent event,
    Emitter<FleetRegistrationState> emit,
  ) async {
    final current = state;
    if (current is! CategoriesLoaded) return;

    try {
      final specGroups = await getSpecDefinitionsUseCase(event.categoryId);
      final draft = await fleetRepository.getDraft(event.categoryId);

      Map<PhotoAngle, FleetPhoto> photos = Map.from(current.photos);
      Map<String, FleetSpecValueInput> specValues =
          Map.from(current.specValues);

      if (draft != null) {
        for (final p in draft.photos) {
          photos[p.angle] = p;
        }
        for (final s in draft.specValues) {
          specValues[s.fieldKey] = s;
        }
      }

      emit(current.copyWith(
        selectedCategoryId: event.categoryId,
        specGroups: specGroups,
        photos: photos,
        specValues: specValues,
      ));
    } catch (_) {
      emit(current.copyWith(
        selectedCategoryId: event.categoryId,
        specGroups: [],
      ));
    }
  }

  void _onUpdatePhoto(
    UpdatePhotoEvent event,
    Emitter<FleetRegistrationState> emit,
  ) {
    final current = state;
    if (current is! CategoriesLoaded) return;

    final updatedPhotos = Map<PhotoAngle, FleetPhoto>.from(current.photos);
    updatedPhotos[event.photo.angle] = event.photo;

    emit(current.copyWith(photos: updatedPhotos));
  }

  void _onUpdateSpecValue(
    UpdateSpecValueEvent event,
    Emitter<FleetRegistrationState> emit,
  ) {
    final current = state;
    if (current is! CategoriesLoaded) return;

    final updatedSpecs =
        Map<String, FleetSpecValueInput>.from(current.specValues);
    updatedSpecs[event.specValue.fieldKey] = event.specValue;

    emit(current.copyWith(specValues: updatedSpecs));
  }

  Future<void> _onSaveDraft(
    SaveDraftEvent event,
    Emitter<FleetRegistrationState> emit,
  ) async {
    final current = state;
    if (current is! CategoriesLoaded) return;

    await fleetRepository.saveDraft(event.request);
    emit(current.copyWith(isDraftSaved: true));
  }

  Future<void> _onSubmitFleetRegistration(
    SubmitFleetRegistrationEvent event,
    Emitter<FleetRegistrationState> emit,
  ) async {
    emit(const FleetSubmitting());
    try {
      final asset = await registerFleetUseCase(event.request);
      emit(FleetSubmitSuccess(asset));
    } on Failure catch (f) {
      emit(FleetSubmitFailure(message: f.message, code: f.code));
    } catch (e) {
      emit(FleetSubmitFailure(message: e.toString(), code: 'ERR_UNKNOWN'));
    }
  }
}

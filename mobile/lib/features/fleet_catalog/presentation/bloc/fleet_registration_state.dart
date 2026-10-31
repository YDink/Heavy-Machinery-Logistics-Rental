import 'package:equatable/equatable.dart';
import 'package:dink_mobile/features/fleet_catalog/domain/entities/category.dart';
import 'package:dink_mobile/features/fleet_catalog/domain/entities/fleet_asset.dart';
import 'package:dink_mobile/features/fleet_catalog/domain/entities/spec_field_group.dart';

abstract class FleetRegistrationState extends Equatable {
  const FleetRegistrationState();

  @override
  List<Object?> get props => [];
}

class FleetRegistrationInitial extends FleetRegistrationState {
  const FleetRegistrationInitial();
}

class CategoriesLoading extends FleetRegistrationState {
  const CategoriesLoading();
}

class CategoriesLoaded extends FleetRegistrationState {
  final List<Category> categories;
  final String? selectedCategoryId;
  final List<SpecFieldGroup> specGroups;
  final Map<PhotoAngle, FleetPhoto> photos;
  final Map<String, FleetSpecValueInput> specValues;
  final bool isDraftSaved;

  const CategoriesLoaded({
    required this.categories,
    this.selectedCategoryId,
    this.specGroups = const [],
    this.photos = const {},
    this.specValues = const {},
    this.isDraftSaved = false,
  });

  bool get hasCompleteWalkaroundPhotos {
    return photos.keys.toSet().containsAll(PhotoAngle.values);
  }

  CategoriesLoaded copyWith({
    List<Category>? categories,
    String? selectedCategoryId,
    List<SpecFieldGroup>? specGroups,
    Map<PhotoAngle, FleetPhoto>? photos,
    Map<String, FleetSpecValueInput>? specValues,
    bool? isDraftSaved,
  }) {
    return CategoriesLoaded(
      categories: categories ?? this.categories,
      selectedCategoryId: selectedCategoryId ?? this.selectedCategoryId,
      specGroups: specGroups ?? this.specGroups,
      photos: photos ?? this.photos,
      specValues: specValues ?? this.specValues,
      isDraftSaved: isDraftSaved ?? this.isDraftSaved,
    );
  }

  @override
  List<Object?> get props => [
        categories,
        selectedCategoryId,
        specGroups,
        photos,
        specValues,
        isDraftSaved,
      ];
}

class FleetSubmitting extends FleetRegistrationState {
  const FleetSubmitting();
}

class FleetSubmitSuccess extends FleetRegistrationState {
  final FleetAsset asset;

  const FleetSubmitSuccess(this.asset);

  @override
  List<Object?> get props => [asset];
}

class FleetSubmitFailure extends FleetRegistrationState {
  final String message;
  final String code;

  const FleetSubmitFailure({required this.message, required this.code});

  @override
  List<Object?> get props => [message, code];
}

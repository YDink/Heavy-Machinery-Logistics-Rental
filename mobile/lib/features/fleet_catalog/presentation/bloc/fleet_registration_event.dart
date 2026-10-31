import 'package:equatable/equatable.dart';
import 'package:dink_mobile/features/fleet_catalog/domain/entities/fleet_asset.dart';

abstract class FleetRegistrationEvent extends Equatable {
  const FleetRegistrationEvent();

  @override
  List<Object?> get props => [];
}

class LoadCategoriesEvent extends FleetRegistrationEvent {
  const LoadCategoriesEvent();
}

class SelectCategoryEvent extends FleetRegistrationEvent {
  final String categoryId;

  const SelectCategoryEvent(this.categoryId);

  @override
  List<Object?> get props => [categoryId];
}

class UpdatePhotoEvent extends FleetRegistrationEvent {
  final FleetPhoto photo;

  const UpdatePhotoEvent(this.photo);

  @override
  List<Object?> get props => [photo];
}

class UpdateSpecValueEvent extends FleetRegistrationEvent {
  final FleetSpecValueInput specValue;

  const UpdateSpecValueEvent(this.specValue);

  @override
  List<Object?> get props => [specValue];
}

class SaveDraftEvent extends FleetRegistrationEvent {
  final FleetRegistrationRequest request;

  const SaveDraftEvent(this.request);

  @override
  List<Object?> get props => [request];
}

class SubmitFleetRegistrationEvent extends FleetRegistrationEvent {
  final FleetRegistrationRequest request;

  const SubmitFleetRegistrationEvent(this.request);

  @override
  List<Object?> get props => [request];
}

import 'package:equatable/equatable.dart';
import '../../domain/entities/active_rental_hud.dart';
import '../../domain/entities/handover_challenge.dart';

abstract class ActiveRentalState extends Equatable {
  const ActiveRentalState();

  @override
  List<Object?> get props => [];
}

class ActiveRentalInitial extends ActiveRentalState {}

class ActiveRentalLoading extends ActiveRentalState {}

class ActiveRentalHUDLoaded extends ActiveRentalState {
  final ActiveRentalHUD hud;

  const ActiveRentalHUDLoaded(this.hud);

  @override
  List<Object?> get props => [hud];
}

class HandoverChallengeGenerated extends ActiveRentalState {
  final HandoverChallenge challenge;

  const HandoverChallengeGenerated(this.challenge);

  @override
  List<Object?> get props => [challenge];
}

class ActiveRentalActionSuccess extends ActiveRentalState {
  final String message;

  const ActiveRentalActionSuccess(this.message);

  @override
  List<Object?> get props => [message];
}

class ActiveRentalError extends ActiveRentalState {
  final String message;

  const ActiveRentalError(this.message);

  @override
  List<Object?> get props => [message];
}

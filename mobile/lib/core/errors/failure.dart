import 'package:equatable/equatable.dart';

abstract class Failure extends Equatable {
  final String message;
  final String code;

  const Failure({required this.message, this.code = 'ERR_UNKNOWN'});

  @override
  List<Object?> get props => [message, code];
}

class ServerFailure extends Failure {
  const ServerFailure({required super.message, super.code = 'ERR_SERVER'});
}

class NetworkFailure extends Failure {
  const NetworkFailure({
    super.message = 'Network connection unavailable. Please check your connectivity.',
    super.code = 'ERR_NETWORK',
  });
}

class AuthFailure extends Failure {
  const AuthFailure({required super.message, super.code = 'ERR_AUTH'});
}

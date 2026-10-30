import 'package:equatable/equatable.dart';

enum UserRole {
  renter,
  owner,
  operator,
  admin,
  unknown,
}

enum AccountStatus {
  unverified,
  active,
  suspended,
  unknown,
}

class AuthUser extends Equatable {
  final String id;
  final String phoneNumber;
  final String fullName;
  final UserRole role;
  final AccountStatus status;

  const AuthUser({
    required this.id,
    required this.phoneNumber,
    required this.fullName,
    required this.role,
    required this.status,
  });

  @override
  List<Object?> get props => [id, phoneNumber, fullName, role, status];
}

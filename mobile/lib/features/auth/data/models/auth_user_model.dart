import 'package:dink_mobile/features/auth/domain/entities/auth_user.dart';

class AuthUserModel extends AuthUser {
  const AuthUserModel({
    required super.id,
    required super.phoneNumber,
    required super.fullName,
    required super.role,
    required super.status,
  });

  factory AuthUserModel.fromJson(Map<String, dynamic> json) {
    UserRole parseRole(String? r) {
      switch (r?.toUpperCase()) {
        case 'RENTER':
          return UserRole.renter;
        case 'OWNER':
          return UserRole.owner;
        case 'OPERATOR':
          return UserRole.operator;
        case 'ADMIN':
          return UserRole.admin;
        default:
          return UserRole.unknown;
      }
    }

    AccountStatus parseStatus(String? s) {
      switch (s?.toUpperCase()) {
        case 'ACTIVE':
          return AccountStatus.active;
        case 'UNVERIFIED':
          return AccountStatus.unverified;
        case 'SUSPENDED':
          return AccountStatus.suspended;
        default:
          return AccountStatus.unknown;
      }
    }

    return AuthUserModel(
      id: json['id'] as String,
      phoneNumber: json['phoneNumber'] as String,
      fullName: json['fullName'] as String? ?? '',
      role: parseRole(json['role'] as String?),
      status: parseStatus(json['status'] as String?),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'phoneNumber': phoneNumber,
      'fullName': fullName,
      'role': role.name.toUpperCase(),
      'status': status.name.toUpperCase(),
    };
  }
}

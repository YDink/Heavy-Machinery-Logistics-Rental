import 'package:dio/dio.dart';
import 'package:dink_mobile/core/network/api_client.dart';
import 'package:dink_mobile/features/auth/data/models/auth_tokens_model.dart';
import 'package:dink_mobile/features/auth/data/models/auth_user_model.dart';

abstract class AuthRemoteDataSource {
  Future<String> requestOtp({required String phoneNumber, required String purpose, String? countryCode});
  Future<AuthTokensModel> verifyOtp({required String verificationSessionId, required String code});
  Future<AuthTokensModel> refreshSession({required String refreshToken});
  Future<AuthUserModel> getCurrentUser();
}

class AuthRemoteDataSourceImpl implements AuthRemoteDataSource {
  final ApiClient apiClient;

  AuthRemoteDataSourceImpl({required this.apiClient});

  @override
  Future<String> requestOtp({
    required String phoneNumber,
    required String purpose,
    String? countryCode,
  }) async {
    try {
      final options = Options();
      if (countryCode != null && countryCode.isNotEmpty) {
        options.headers = {'X-Country-Code': countryCode};
      }

      final response = await apiClient.dio.post(
        '/auth/otp/request',
        data: {
          'phoneNumber': phoneNumber,
          'purpose': purpose,
        },
        options: options,
      );

      final data = response.data['data'] as Map<String, dynamic>;
      return data['verificationSessionId'] as String;
    } on DioException catch (e) {
      throw apiClient.mapDioExceptionToFailure(e);
    }
  }

  @override
  Future<AuthTokensModel> verifyOtp({
    required String verificationSessionId,
    required String code,
  }) async {
    try {
      final response = await apiClient.dio.post(
        '/auth/otp/verify',
        data: {
          'verificationSessionId': verificationSessionId,
          'code': code,
        },
      );

      final data = response.data['data'] as Map<String, dynamic>;
      return AuthTokensModel.fromJson(data);
    } on DioException catch (e) {
      throw apiClient.mapDioExceptionToFailure(e);
    }
  }

  @override
  Future<AuthTokensModel> refreshSession({required String refreshToken}) async {
    try {
      final response = await apiClient.dio.post(
        '/auth/session/refresh',
        data: {
          'refreshToken': refreshToken,
        },
      );

      final data = response.data['data'] as Map<String, dynamic>;
      return AuthTokensModel.fromJson(data);
    } on DioException catch (e) {
      throw apiClient.mapDioExceptionToFailure(e);
    }
  }

  @override
  Future<AuthUserModel> getCurrentUser() async {
    try {
      final response = await apiClient.dio.get('/auth/users/me');
      final data = response.data['data'] as Map<String, dynamic>;
      return AuthUserModel.fromJson(data);
    } on DioException catch (e) {
      throw apiClient.mapDioExceptionToFailure(e);
    }
  }
}

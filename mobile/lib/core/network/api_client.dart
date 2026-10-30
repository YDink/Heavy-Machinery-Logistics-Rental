import 'package:dio/dio.dart';
import 'package:dink_mobile/core/config/app_config.dart';
import 'package:dink_mobile/core/errors/failure.dart';
import 'package:dink_mobile/core/storage/secure_storage_service.dart';

class ApiClient {
  late final Dio dio;
  final SecureStorageService storageService;

  ApiClient({Dio? customDio, required this.storageService}) {
    dio = customDio ??
        Dio(
          BaseOptions(
            baseUrl: AppConfig.instance.apiBaseUrl,
            connectTimeout: const Duration(seconds: 10),
            receiveTimeout: const Duration(seconds: 10),
            headers: {
              'Content-Type': 'application/json',
              'Accept': 'application/json',
            },
          ),
        );

    dio.interceptors.add(
      InterceptorsWrapper(
        onRequest: (options, handler) async {
          final countryCode = await storageService.getCountryCode();
          options.headers['X-Country-Code'] = countryCode;

          final token = await storageService.getAccessToken();
          if (token != null && token.isNotEmpty) {
            options.headers['Authorization'] = 'Bearer $token';
          }
          return handler.next(options);
        },
        onError: (DioException e, handler) {
          return handler.next(e);
        },
      ),
    );
  }

  Failure mapDioExceptionToFailure(DioException e) {
    if (e.type == DioExceptionType.connectionTimeout ||
        e.type == DioExceptionType.receiveTimeout ||
        e.type == DioExceptionType.connectionError) {
      return const NetworkFailure();
    }

    if (e.response != null && e.response?.data is Map<String, dynamic>) {
      final data = e.response!.data as Map<String, dynamic>;
      if (data.containsKey('error') && data['error'] is Map<String, dynamic>) {
        final errorObj = data['error'] as Map<String, dynamic>;
        final detail = errorObj['detail'] as String? ?? 'An unexpected error occurred.';
        final code = errorObj['code'] as String? ?? 'ERR_API';
        return ServerFailure(message: detail, code: code);
      }
    }

    return ServerFailure(
      message: e.message ?? 'Server communication error',
      code: 'ERR_HTTP_${e.response?.statusCode ?? 500}',
    );
  }
}

import 'package:get_it/get_it.dart';
import 'package:dink_mobile/core/network/api_client.dart';
import 'package:dink_mobile/core/storage/secure_storage_service.dart';
import 'package:dink_mobile/features/auth/data/datasources/auth_remote_datasource.dart';
import 'package:dink_mobile/features/auth/data/repositories/auth_repository_impl.dart';
import 'package:dink_mobile/features/auth/domain/repositories/auth_repository.dart';
import 'package:dink_mobile/features/auth/domain/usecases/get_current_user_usecase.dart';
import 'package:dink_mobile/features/auth/domain/usecases/request_otp_usecase.dart';
import 'package:dink_mobile/features/auth/domain/usecases/verify_otp_usecase.dart';
import 'package:dink_mobile/features/auth/presentation/bloc/auth_bloc.dart';

final sl = GetIt.instance;

Future<void> initDependencies() async {
  // Core Services
  if (!sl.isRegistered<SecureStorageService>()) {
    sl.registerLazySingleton<SecureStorageService>(() => SecureStorageService());
  }

  if (!sl.isRegistered<ApiClient>()) {
    sl.registerLazySingleton<ApiClient>(
      () => ApiClient(storageService: sl<SecureStorageService>()),
    );
  }

  // Data Sources
  if (!sl.isRegistered<AuthRemoteDataSource>()) {
    sl.registerLazySingleton<AuthRemoteDataSource>(
      () => AuthRemoteDataSourceImpl(apiClient: sl<ApiClient>()),
    );
  }

  // Repositories
  if (!sl.isRegistered<AuthRepository>()) {
    sl.registerLazySingleton<AuthRepository>(
      () => AuthRepositoryImpl(
        remoteDataSource: sl<AuthRemoteDataSource>(),
        storageService: sl<SecureStorageService>(),
      ),
    );
  }

  // Use Cases
  if (!sl.isRegistered<RequestOtpUseCase>()) {
    sl.registerLazySingleton<RequestOtpUseCase>(
      () => RequestOtpUseCase(sl<AuthRepository>()),
    );
  }

  if (!sl.isRegistered<VerifyOtpUseCase>()) {
    sl.registerLazySingleton<VerifyOtpUseCase>(
      () => VerifyOtpUseCase(sl<AuthRepository>()),
    );
  }

  if (!sl.isRegistered<GetCurrentUserUseCase>()) {
    sl.registerLazySingleton<GetCurrentUserUseCase>(
      () => GetCurrentUserUseCase(sl<AuthRepository>()),
    );
  }

  // BLoC
  if (!sl.isRegistered<AuthBloc>()) {
    sl.registerFactory<AuthBloc>(
      () => AuthBloc(
        requestOtpUseCase: sl<RequestOtpUseCase>(),
        verifyOtpUseCase: sl<VerifyOtpUseCase>(),
        getCurrentUserUseCase: sl<GetCurrentUserUseCase>(),
        authRepository: sl<AuthRepository>(),
      ),
    );
  }
}

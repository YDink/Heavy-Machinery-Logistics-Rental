import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:get_it/get_it.dart';
import 'package:dink_mobile/core/localization/locale_controller.dart';
import 'package:dink_mobile/core/network/api_client.dart';
import 'package:dink_mobile/core/storage/secure_storage_service.dart';
import 'package:dink_mobile/features/auth/data/datasources/auth_remote_datasource.dart';
import 'package:dink_mobile/features/auth/data/repositories/auth_repository_impl.dart';
import 'package:dink_mobile/features/auth/domain/repositories/auth_repository.dart';
import 'package:dink_mobile/features/auth/domain/usecases/get_current_user_usecase.dart';
import 'package:dink_mobile/features/auth/domain/usecases/request_otp_usecase.dart';
import 'package:dink_mobile/features/auth/domain/usecases/verify_otp_usecase.dart';
import 'package:dink_mobile/features/auth/presentation/bloc/auth_bloc.dart';
import 'package:dink_mobile/features/fleet_catalog/data/datasources/fleet_local_datasource.dart';
import 'package:dink_mobile/features/fleet_catalog/data/datasources/fleet_remote_datasource.dart';
import 'package:dink_mobile/features/fleet_catalog/data/repositories/fleet_repository_impl.dart';
import 'package:dink_mobile/features/fleet_catalog/domain/repositories/fleet_repository.dart';
import 'package:dink_mobile/features/fleet_catalog/domain/usecases/get_categories_usecase.dart';
import 'package:dink_mobile/features/fleet_catalog/domain/usecases/get_spec_definitions_usecase.dart';
import 'package:dink_mobile/features/fleet_catalog/domain/usecases/register_fleet_usecase.dart';
import 'package:dink_mobile/features/fleet_catalog/presentation/bloc/fleet_registration_bloc.dart';

final sl = GetIt.instance;

Future<void> initDependencies() async {
  // Localization & State
  if (!sl.isRegistered<LocaleController>()) {
    sl.registerLazySingleton<LocaleController>(() => LocaleController());
  }

  // Core Services
  if (!sl.isRegistered<FlutterSecureStorage>()) {
    sl.registerLazySingleton<FlutterSecureStorage>(
        () => const FlutterSecureStorage());
  }

  if (!sl.isRegistered<SecureStorageService>()) {
    sl.registerLazySingleton<SecureStorageService>(
        () => SecureStorageService());
  }

  if (!sl.isRegistered<ApiClient>()) {
    sl.registerLazySingleton<ApiClient>(
      () => ApiClient(storageService: sl<SecureStorageService>()),
    );
  }

  // Auth Subsystem
  if (!sl.isRegistered<AuthRemoteDataSource>()) {
    sl.registerLazySingleton<AuthRemoteDataSource>(
      () => AuthRemoteDataSourceImpl(apiClient: sl<ApiClient>()),
    );
  }

  if (!sl.isRegistered<AuthRepository>()) {
    sl.registerLazySingleton<AuthRepository>(
      () => AuthRepositoryImpl(
        remoteDataSource: sl<AuthRemoteDataSource>(),
        storageService: sl<SecureStorageService>(),
      ),
    );
  }

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

  // Fleet Subsystem (Partition 05-02)
  if (!sl.isRegistered<FleetRemoteDataSource>()) {
    sl.registerLazySingleton<FleetRemoteDataSource>(
      () => FleetRemoteDataSourceImpl(apiClient: sl<ApiClient>()),
    );
  }

  if (!sl.isRegistered<FleetLocalDataSource>()) {
    sl.registerLazySingleton<FleetLocalDataSource>(
      () => FleetLocalDataSourceImpl(storage: sl<FlutterSecureStorage>()),
    );
  }

  if (!sl.isRegistered<FleetRepository>()) {
    sl.registerLazySingleton<FleetRepository>(
      () => FleetRepositoryImpl(
        remoteDataSource: sl<FleetRemoteDataSource>(),
        localDataSource: sl<FleetLocalDataSource>(),
      ),
    );
  }

  if (!sl.isRegistered<GetCategoriesUseCase>()) {
    sl.registerLazySingleton<GetCategoriesUseCase>(
      () => GetCategoriesUseCase(sl<FleetRepository>()),
    );
  }

  if (!sl.isRegistered<GetSpecDefinitionsUseCase>()) {
    sl.registerLazySingleton<GetSpecDefinitionsUseCase>(
      () => GetSpecDefinitionsUseCase(sl<FleetRepository>()),
    );
  }

  if (!sl.isRegistered<RegisterFleetUseCase>()) {
    sl.registerLazySingleton<RegisterFleetUseCase>(
      () => RegisterFleetUseCase(sl<FleetRepository>()),
    );
  }

  if (!sl.isRegistered<FleetRegistrationBloc>()) {
    sl.registerFactory<FleetRegistrationBloc>(
      () => FleetRegistrationBloc(
        getCategoriesUseCase: sl<GetCategoriesUseCase>(),
        getSpecDefinitionsUseCase: sl<GetSpecDefinitionsUseCase>(),
        registerFleetUseCase: sl<RegisterFleetUseCase>(),
        fleetRepository: sl<FleetRepository>(),
      ),
    );
  }
}

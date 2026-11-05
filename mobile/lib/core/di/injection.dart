import 'package:dink_mobile/features/discovery/data/datasources/discovery_remote_data_source.dart';
import 'package:dink_mobile/features/discovery/data/repositories/discovery_repository_impl.dart';
import 'package:dink_mobile/features/discovery/domain/entities/fleet_search_filter.dart';
import 'package:dink_mobile/features/discovery/domain/repositories/discovery_repository.dart';
import 'package:dink_mobile/features/discovery/domain/usecases/compare_fleets_usecase.dart';
import 'package:dink_mobile/features/discovery/domain/usecases/get_fleet_availability_usecase.dart';
import 'package:dink_mobile/features/discovery/domain/usecases/search_fleets_usecase.dart';
import 'package:dink_mobile/features/discovery/presentation/bloc/discovery_bloc.dart';
import 'package:dink_mobile/features/discovery/presentation/bloc/discovery_state.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:get_it/get_it.dart';
import 'package:dink_mobile/core/localization/locale_controller.dart';
import 'package:dink_mobile/core/network/api_client.dart';
import 'package:dink_mobile/core/storage/secure_storage_service.dart';
import 'package:dink_mobile/features/auth/data/datasources/auth_remote_data_source.dart';
import 'package:dink_mobile/features/auth/data/repositories/auth_repository_impl.dart';
import 'package:dink_mobile/features/auth/domain/repositories/auth_repository.dart';
import 'package:dink_mobile/features/auth/domain/usecases/get_current_user_usecase.dart';
import 'package:dink_mobile/features/auth/domain/usecases/request_otp_usecase.dart';
import 'package:dink_mobile/features/auth/domain/usecases/verify_otp_usecase.dart';
import 'package:dink_mobile/features/auth/presentation/bloc/auth_bloc.dart';
import 'package:dink_mobile/features/fleet_catalog/data/datasources/fleet_local_data_source.dart';
import 'package:dink_mobile/features/fleet_catalog/data/datasources/fleet_remote_data_source.dart';
import 'package:dink_mobile/features/fleet_catalog/data/repositories/fleet_repository_impl.dart';
import 'package:dink_mobile/features/fleet_catalog/domain/repositories/fleet_repository.dart';
import 'package:dink_mobile/features/fleet_catalog/domain/usecases/get_categories_usecase.dart';
import 'package:dink_mobile/features/fleet_catalog/domain/usecases/get_spec_definitions_usecase.dart';
import 'package:dink_mobile/features/fleet_catalog/domain/usecases/register_fleet_usecase.dart';
import 'package:dink_mobile/features/fleet_catalog/presentation/bloc/fleet_registration_bloc.dart';
import 'package:dink_mobile/features/attachment/data/datasources/attachment_local_data_source.dart';
import 'package:dink_mobile/features/attachment/data/datasources/attachment_remote_data_source.dart';
import 'package:dink_mobile/features/attachment/data/repositories/attachment_repository_impl.dart';
import 'package:dink_mobile/features/attachment/domain/repositories/attachment_repository.dart';
import 'package:dink_mobile/features/attachment/domain/usecases/get_attachment_inventory_usecase.dart';
import 'package:dink_mobile/features/attachment/domain/usecases/link_attachment_usecase.dart';
import 'package:dink_mobile/features/attachment/domain/usecases/register_attachment_usecase.dart';
import 'package:dink_mobile/features/attachment/presentation/bloc/attachment_bloc.dart';
import 'package:dink_mobile/features/booking/data/datasources/booking_remote_data_source.dart';
import 'package:dink_mobile/features/booking/data/repositories/booking_repository_impl.dart';
import 'package:dink_mobile/features/booking/domain/repositories/booking_repository.dart';
import 'package:dink_mobile/features/booking/domain/usecases/create_booking_usecase.dart';
import 'package:dink_mobile/features/booking/domain/usecases/get_booking_usecase.dart';
import 'package:dink_mobile/features/booking/domain/usecases/submit_payment_slip_usecase.dart';
import 'package:dink_mobile/features/booking/presentation/bloc/booking_bloc.dart';
import 'package:dink_mobile/features/active_rental/data/datasources/active_rental_remote_data_source.dart';
import 'package:dink_mobile/features/active_rental/data/repositories/active_rental_repository_impl.dart';
import 'package:dink_mobile/features/active_rental/domain/repositories/active_rental_repository.dart';
import 'package:dink_mobile/features/active_rental/domain/usecases/get_active_rental_hud_usecase.dart';
import 'package:dink_mobile/features/active_rental/domain/usecases/inspect_handover_usecase.dart';
import 'package:dink_mobile/features/active_rental/domain/usecases/generate_handover_pin_usecase.dart';
import 'package:dink_mobile/features/active_rental/domain/usecases/confirm_handover_usecase.dart';
import 'package:dink_mobile/features/active_rental/domain/usecases/request_extension_usecase.dart';
import 'package:dink_mobile/features/active_rental/domain/usecases/report_breakdown_usecase.dart';
import 'package:dink_mobile/features/active_rental/domain/usecases/resolve_breakdown_usecase.dart';
import 'package:dink_mobile/features/active_rental/presentation/bloc/active_rental_bloc.dart';
import 'package:dink_mobile/features/rental_completion/data/datasources/rental_completion_remote_data_source.dart';
import 'package:dink_mobile/features/rental_completion/data/repositories/rental_completion_repository_impl.dart';
import 'package:dink_mobile/features/rental_completion/domain/repositories/rental_completion_repository.dart';
import 'package:dink_mobile/features/rental_completion/domain/usecases/final_reconcile_usecase.dart';
import 'package:dink_mobile/features/rental_completion/domain/usecases/initiate_return_usecase.dart';
import 'package:dink_mobile/features/rental_completion/domain/usecases/inspect_return_usecase.dart';
import 'package:dink_mobile/features/rental_completion/domain/usecases/lodge_damage_claim_usecase.dart';
import 'package:dink_mobile/features/rental_completion/domain/usecases/respond_damage_claim_usecase.dart';
import 'package:dink_mobile/features/rental_completion/domain/usecases/submit_review_usecase.dart';
import 'package:dink_mobile/features/rental_completion/presentation/bloc/rental_completion_bloc.dart';

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

  // Attachment Subsystem (Partition 05-03)
  if (!sl.isRegistered<AttachmentRemoteDataSource>()) {
    sl.registerLazySingleton<AttachmentRemoteDataSource>(
      () => AttachmentRemoteDataSourceImpl(apiClient: sl<ApiClient>()),
    );
  }

  if (!sl.isRegistered<AttachmentLocalDataSource>()) {
    sl.registerLazySingleton<AttachmentLocalDataSource>(
      () => AttachmentLocalDataSourceImpl(storage: sl<FlutterSecureStorage>()),
    );
  }

  if (!sl.isRegistered<AttachmentRepository>()) {
    sl.registerLazySingleton<AttachmentRepository>(
      () => AttachmentRepositoryImpl(
        remoteDataSource: sl<AttachmentRemoteDataSource>(),
        localDataSource: sl<AttachmentLocalDataSource>(),
      ),
    );
  }

  if (!sl.isRegistered<RegisterAttachmentUseCase>()) {
    sl.registerLazySingleton<RegisterAttachmentUseCase>(
      () => RegisterAttachmentUseCase(sl<AttachmentRepository>()),
    );
  }

  if (!sl.isRegistered<GetAttachmentInventoryUseCase>()) {
    sl.registerLazySingleton<GetAttachmentInventoryUseCase>(
      () => GetAttachmentInventoryUseCase(sl<AttachmentRepository>()),
    );
  }

  if (!sl.isRegistered<LinkAttachmentUseCase>()) {
    sl.registerLazySingleton<LinkAttachmentUseCase>(
      () => LinkAttachmentUseCase(sl<AttachmentRepository>()),
    );
  }

  if (!sl.isRegistered<AttachmentBloc>()) {
    sl.registerFactory<AttachmentBloc>(
      () => AttachmentBloc(
        registerAttachmentUseCase: sl<RegisterAttachmentUseCase>(),
        getInventoryUseCase: sl<GetAttachmentInventoryUseCase>(),
        linkAttachmentUseCase: sl<LinkAttachmentUseCase>(),
        attachmentRepository: sl<AttachmentRepository>(),
      ),
    );
  }

  // Discovery Subsystem (Partition 05-04)
  if (!sl.isRegistered<DiscoveryRemoteDataSource>()) {
    sl.registerLazySingleton<DiscoveryRemoteDataSource>(
      () => DiscoveryRemoteDataSourceImpl(apiClient: sl<ApiClient>()),
    );
  }

  if (!sl.isRegistered<DiscoveryRepository>()) {
    sl.registerLazySingleton<DiscoveryRepository>(
      () => DiscoveryRepositoryImpl(remoteDataSource: sl<DiscoveryRemoteDataSource>()),
    );
  }

  if (!sl.isRegistered<SearchFleetsUseCase>()) {
    sl.registerLazySingleton<SearchFleetsUseCase>(
      () => SearchFleetsUseCase(sl<DiscoveryRepository>()),
    );
  }

  if (!sl.isRegistered<CompareFleetsUseCase>()) {
    sl.registerLazySingleton<CompareFleetsUseCase>(
      () => CompareFleetsUseCase(sl<DiscoveryRepository>()),
    );
  }

  if (!sl.isRegistered<GetFleetAvailabilityUseCase>()) {
    sl.registerLazySingleton<GetFleetAvailabilityUseCase>(
      () => GetFleetAvailabilityUseCase(sl<DiscoveryRepository>()),
    );
  }

  if (!sl.isRegistered<DiscoveryBloc>()) {
    sl.registerFactory<DiscoveryBloc>(
      () => DiscoveryBloc(
        searchFleetsUseCase: sl<SearchFleetsUseCase>(),
        compareFleetsUseCase: sl<CompareFleetsUseCase>(),
        getFleetAvailabilityUseCase: sl<GetFleetAvailabilityUseCase>(),
        initialState: const DiscoveryState(
          filter: FleetSearchFilter(latitude: 9.0300, longitude: 38.7400),
        ),
      ),
    );
  }

  // Booking Subsystem (Partition 05-05)
  if (!sl.isRegistered<BookingRemoteDataSource>()) {
    sl.registerLazySingleton<BookingRemoteDataSource>(
      () => BookingRemoteDataSourceImpl(apiClient: sl<ApiClient>()),
    );
  }

  if (!sl.isRegistered<BookingRepository>()) {
    sl.registerLazySingleton<BookingRepository>(
      () => BookingRepositoryImpl(remoteDataSource: sl<BookingRemoteDataSource>()),
    );
  }

  if (!sl.isRegistered<CreateBookingUseCase>()) {
    sl.registerLazySingleton<CreateBookingUseCase>(
      () => CreateBookingUseCase(sl<BookingRepository>()),
    );
  }

  if (!sl.isRegistered<GetBookingUseCase>()) {
    sl.registerLazySingleton<GetBookingUseCase>(
      () => GetBookingUseCase(sl<BookingRepository>()),
    );
  }

  if (!sl.isRegistered<SubmitPaymentSlipUseCase>()) {
    sl.registerLazySingleton<SubmitPaymentSlipUseCase>(
      () => SubmitPaymentSlipUseCase(sl<BookingRepository>()),
    );
  }

  if (!sl.isRegistered<BookingBloc>()) {
    sl.registerFactory<BookingBloc>(
      () => BookingBloc(repository: sl<BookingRepository>()),
    );
  }

  // Active Rental (Partition 05-06)
  if (!sl.isRegistered<ActiveRentalRemoteDataSource>()) {
    sl.registerLazySingleton<ActiveRentalRemoteDataSource>(
      () => ActiveRentalRemoteDataSourceImpl(apiClient: sl<ApiClient>()),
    );
  }

  if (!sl.isRegistered<ActiveRentalRepository>()) {
    sl.registerLazySingleton<ActiveRentalRepository>(
      () => ActiveRentalRepositoryImpl(remoteDataSource: sl<ActiveRentalRemoteDataSource>()),
    );
  }

  if (!sl.isRegistered<GetActiveRentalHUDUseCase>()) {
    sl.registerLazySingleton<GetActiveRentalHUDUseCase>(
      () => GetActiveRentalHUDUseCase(sl<ActiveRentalRepository>()),
    );
  }

  if (!sl.isRegistered<InspectHandoverUseCase>()) {
    sl.registerLazySingleton<InspectHandoverUseCase>(
      () => InspectHandoverUseCase(sl<ActiveRentalRepository>()),
    );
  }

  if (!sl.isRegistered<GenerateHandoverPINUseCase>()) {
    sl.registerLazySingleton<GenerateHandoverPINUseCase>(
      () => GenerateHandoverPINUseCase(sl<ActiveRentalRepository>()),
    );
  }

  if (!sl.isRegistered<ConfirmHandoverUseCase>()) {
    sl.registerLazySingleton<ConfirmHandoverUseCase>(
      () => ConfirmHandoverUseCase(sl<ActiveRentalRepository>()),
    );
  }

  if (!sl.isRegistered<RequestExtensionUseCase>()) {
    sl.registerLazySingleton<RequestExtensionUseCase>(
      () => RequestExtensionUseCase(sl<ActiveRentalRepository>()),
    );
  }

  if (!sl.isRegistered<ReportBreakdownUseCase>()) {
    sl.registerLazySingleton<ReportBreakdownUseCase>(
      () => ReportBreakdownUseCase(sl<ActiveRentalRepository>()),
    );
  }

  if (!sl.isRegistered<ResolveBreakdownUseCase>()) {
    sl.registerLazySingleton<ResolveBreakdownUseCase>(
      () => ResolveBreakdownUseCase(sl<ActiveRentalRepository>()),
    );
  }

  if (!sl.isRegistered<ActiveRentalBloc>()) {
    sl.registerFactory<ActiveRentalBloc>(
      () => ActiveRentalBloc(repository: sl<ActiveRentalRepository>()),
    );
  }

  // Rental Completion, Return & Damage (Partition 05-07)
  if (!sl.isRegistered<RentalCompletionRemoteDataSource>()) {
    sl.registerLazySingleton<RentalCompletionRemoteDataSource>(
      () => RentalCompletionRemoteDataSourceImpl(apiClient: sl<ApiClient>()),
    );
  }

  if (!sl.isRegistered<RentalCompletionRepository>()) {
    sl.registerLazySingleton<RentalCompletionRepository>(
      () => RentalCompletionRepositoryImpl(remoteDataSource: sl<RentalCompletionRemoteDataSource>()),
    );
  }

  if (!sl.isRegistered<InitiateReturnUseCase>()) {
    sl.registerLazySingleton<InitiateReturnUseCase>(
      () => InitiateReturnUseCase(sl<RentalCompletionRepository>()),
    );
  }

  if (!sl.isRegistered<InspectReturnUseCase>()) {
    sl.registerLazySingleton<InspectReturnUseCase>(
      () => InspectReturnUseCase(sl<RentalCompletionRepository>()),
    );
  }

  if (!sl.isRegistered<LodgeDamageClaimUseCase>()) {
    sl.registerLazySingleton<LodgeDamageClaimUseCase>(
      () => LodgeDamageClaimUseCase(sl<RentalCompletionRepository>()),
    );
  }

  if (!sl.isRegistered<RespondDamageClaimUseCase>()) {
    sl.registerLazySingleton<RespondDamageClaimUseCase>(
      () => RespondDamageClaimUseCase(sl<RentalCompletionRepository>()),
    );
  }

  if (!sl.isRegistered<FinalReconcileUseCase>()) {
    sl.registerLazySingleton<FinalReconcileUseCase>(
      () => FinalReconcileUseCase(sl<RentalCompletionRepository>()),
    );
  }

  if (!sl.isRegistered<SubmitReviewUseCase>()) {
    sl.registerLazySingleton<SubmitReviewUseCase>(
      () => SubmitReviewUseCase(sl<RentalCompletionRepository>()),
    );
  }

  if (!sl.isRegistered<RentalCompletionBloc>()) {
    sl.registerFactory<RentalCompletionBloc>(
      () => RentalCompletionBloc(
        initiateReturnUseCase: sl<InitiateReturnUseCase>(),
        inspectReturnUseCase: sl<InspectReturnUseCase>(),
        lodgeDamageClaimUseCase: sl<LodgeDamageClaimUseCase>(),
        respondDamageClaimUseCase: sl<RespondDamageClaimUseCase>(),
        finalReconcileUseCase: sl<FinalReconcileUseCase>(),
        submitReviewUseCase: sl<SubmitReviewUseCase>(),
      ),
    );
  }
}

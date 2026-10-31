import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:flutter_localizations/flutter_localizations.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:dink_mobile/features/attachment/domain/entities/attachment.dart';
import 'package:dink_mobile/features/attachment/domain/entities/attachment_inventory_item.dart';
import 'package:dink_mobile/features/attachment/domain/entities/bundled_attachment_link.dart';
import 'package:dink_mobile/features/attachment/domain/repositories/attachment_repository.dart';
import 'package:dink_mobile/features/attachment/domain/usecases/get_attachment_inventory_usecase.dart';
import 'package:dink_mobile/features/attachment/domain/usecases/link_attachment_usecase.dart';
import 'package:dink_mobile/features/attachment/domain/usecases/register_attachment_usecase.dart';
import 'package:dink_mobile/features/attachment/presentation/bloc/attachment_bloc.dart';
import 'package:dink_mobile/features/attachment/presentation/pages/attachment_registration_page.dart';
import 'package:dink_mobile/l10n/app_localizations.dart';

class StubAttachmentRepository implements AttachmentRepository {
  AttachmentRegistrationRequest? registeredRequest;
  AttachmentRegistrationRequest? draft;

  @override
  Future<Attachment> registerAttachment(AttachmentRegistrationRequest request) async {
    registeredRequest = request;
    return Attachment(
      id: 'att-999',
      ownerId: 'owner-1',
      attachmentType: request.attachmentType,
      make: request.make,
      model: request.model,
      subTypeSize: request.subTypeSize,
      serialNumber: request.serialNumber,
      operatingWeightKg: request.operatingWeightKg,
      minCarrierWeightTonnes: request.minCarrierWeightTonnes,
      maxCarrierWeightTonnes: request.maxCarrierWeightTonnes,
      requiredPressureBar: request.requiredPressureBar,
      requiredFlowMinLmin: request.requiredFlowMinLmin,
      requiredFlowMaxLmin: request.requiredFlowMaxLmin,
      pinDiameterMm: request.pinDiameterMm,
      pinSpreadMm: request.pinSpreadMm,
      compatibleSubcategories: request.compatibleSubcategories,
      standaloneDailyRateMinor: request.standaloneDailyRateMinor,
      securityDepositMinor: request.securityDepositMinor,
      minRentalDays: request.minRentalDays,
      status: 'AVAILABLE',
      photos: request.photos,
      createdAt: DateTime.utc(2026, 3, 31),
    );
  }

  @override
  Future<List<AttachmentInventoryItem>> getInventory({String? status}) async => [];

  @override
  Future<BundledAttachmentLink> linkAttachmentToFleet({
    required String fleetId,
    required String attachmentId,
    required int bundledDailyRateMinor,
  }) async {
    return BundledAttachmentLink(
      fleetId: fleetId,
      attachmentId: attachmentId,
      bundledDailyRateMinor: bundledDailyRateMinor,
      discountPercentage: 20.0,
      status: 'ACTIVE',
    );
  }

  @override
  Future<void> saveDraft(AttachmentRegistrationRequest request) async {
    draft = request;
  }

  @override
  Future<AttachmentRegistrationRequest?> getDraft() async => draft;

  @override
  Future<void> clearDraft() async {
    draft = null;
  }
}

Widget buildTestableWidget({
  required Widget child,
  required AttachmentBloc bloc,
}) {
  return MaterialApp(
    localizationsDelegates: const [
      AppLocalizations.delegate,
      GlobalMaterialLocalizations.delegate,
      GlobalWidgetsLocalizations.delegate,
      GlobalCupertinoLocalizations.delegate,
    ],
    supportedLocales: AppLocalizations.supportedLocales,
    home: BlocProvider<AttachmentBloc>.value(
      value: bloc,
      child: child,
    ),
  );
}

void main() {
  late StubAttachmentRepository stubRepo;
  late AttachmentBloc bloc;

  setUp(() {
    stubRepo = StubAttachmentRepository();
    bloc = AttachmentBloc(
      registerAttachmentUseCase: RegisterAttachmentUseCase(stubRepo),
      getInventoryUseCase: GetAttachmentInventoryUseCase(stubRepo),
      linkAttachmentUseCase: LinkAttachmentUseCase(stubRepo),
      attachmentRepository: stubRepo,
    );
  });

  tearDown(() {
    bloc.close();
  });

  testWidgets('renders registration form fields and reticle correctly', (tester) async {
    await tester.binding.setSurfaceSize(const Size(800, 2200));
    addTearDown(() => tester.binding.setSurfaceSize(null));

    await tester.pumpWidget(
      buildTestableWidget(
        child: const AttachmentRegistrationPage(),
        bloc: bloc,
      ),
    );
    await tester.pumpAndSettle();

    // Verify Title
    expect(find.text('Register Attachment'), findsOneWidget);

    // Verify make, model, subtype, serial text fields
    expect(find.byType(TextFormField), findsWidgets);
    expect(find.text('Make / Manufacturer'), findsOneWidget);
    expect(find.text('Model'), findsOneWidget);
    expect(find.text('Tool Sizing Specification'), findsOneWidget);
    expect(find.text('Serial PIN Number'), findsOneWidget);

    // Verify 4-angle conditioning photos reticle
    expect(find.text('Attachment Conditioning Photos'), findsOneWidget);
    expect(find.text('Overall Tool Profile'), findsOneWidget);
    expect(find.text('Mounting Bracket & Hitch'), findsOneWidget);
    expect(find.text('Pins & Bushings Wear'), findsOneWidget);
    expect(find.text('Working Tip / Chisel Point'), findsOneWidget);

    // Verify submit button
    expect(find.text('Submit for Verification'), findsOneWidget);
  });

  testWidgets('validates required fields when submitted empty', (tester) async {
    await tester.binding.setSurfaceSize(const Size(800, 2200));
    addTearDown(() => tester.binding.setSurfaceSize(null));

    await tester.pumpWidget(
      buildTestableWidget(
        child: const AttachmentRegistrationPage(),
        bloc: bloc,
      ),
    );
    await tester.pumpAndSettle();

    // Tap submit button without filling fields
    final submitButton = find.text('Submit for Verification');
    await tester.ensureVisible(submitButton);
    await tester.tap(submitButton);
    await tester.pumpAndSettle();

    // Verify error text
    expect(find.text('Please fill in all required fields'), findsWidgets);
  });
}

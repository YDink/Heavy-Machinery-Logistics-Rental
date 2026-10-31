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
import 'package:dink_mobile/features/attachment/presentation/pages/attachment_inventory_page.dart';
import 'package:dink_mobile/l10n/app_localizations.dart';

class StubAttachmentRepository implements AttachmentRepository {
  List<AttachmentInventoryItem> items = [];

  @override
  Future<Attachment> registerAttachment(AttachmentRegistrationRequest request) async {
    throw UnimplementedError();
  }

  @override
  Future<List<AttachmentInventoryItem>> getInventory({String? status}) async {
    return items;
  }

  @override
  Future<BundledAttachmentLink> linkAttachmentToFleet({
    required String fleetId,
    required String attachmentId,
    required int bundledDailyRateMinor,
  }) async {
    throw UnimplementedError();
  }

  @override
  Future<void> saveDraft(AttachmentRegistrationRequest request) async {}

  @override
  Future<AttachmentRegistrationRequest?> getDraft() async => null;

  @override
  Future<void> clearDraft() async {}
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

  testWidgets('renders empty inventory placeholder when list is empty', (tester) async {
    stubRepo.items = [];

    await tester.pumpWidget(
      buildTestableWidget(
        child: const AttachmentInventoryPage(),
        bloc: bloc,
      ),
    );
    await tester.runAsync(() async {
      await Future.delayed(const Duration(milliseconds: 50));
    });
    await tester.pumpAndSettle();

    expect(find.text('My Attachments'), findsOneWidget);
    expect(find.text('No attachments registered yet'), findsOneWidget);
  });

  testWidgets('renders list of attachments when inventory exists', (tester) async {
    stubRepo.items = [
      const AttachmentInventoryItem(
        id: 'att-1',
        attachmentType: 'HYDRAULIC_BREAKER',
        make: 'Soosan',
        model: 'SB70',
        subTypeSize: '135mm Chisel',
        standaloneDailyRateMinor: 1500000,
        status: 'ACTIVE',
        isLinkedToFleet: true,
      ),
    ];

    await tester.pumpWidget(
      buildTestableWidget(
        child: const AttachmentInventoryPage(),
        bloc: bloc,
      ),
    );
    await tester.runAsync(() async {
      await Future.delayed(const Duration(milliseconds: 50));
    });
    await tester.pumpAndSettle();

    expect(find.text('Soosan SB70'), findsOneWidget);
    expect(find.text('Available'), findsOneWidget);
    expect(find.text('Bundled with Host'), findsOneWidget);
  });
}

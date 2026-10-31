import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:flutter_localizations/flutter_localizations.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:dink_mobile/features/fleet_catalog/domain/entities/category.dart';
import 'package:dink_mobile/features/fleet_catalog/domain/entities/fleet_asset.dart';
import 'package:dink_mobile/features/fleet_catalog/domain/entities/spec_field_group.dart';
import 'package:dink_mobile/features/fleet_catalog/domain/repositories/fleet_repository.dart';
import 'package:dink_mobile/features/fleet_catalog/domain/usecases/get_categories_usecase.dart';
import 'package:dink_mobile/features/fleet_catalog/domain/usecases/get_spec_definitions_usecase.dart';
import 'package:dink_mobile/features/fleet_catalog/domain/usecases/register_fleet_usecase.dart';
import 'package:dink_mobile/features/fleet_catalog/presentation/bloc/fleet_registration_bloc.dart';
import 'package:dink_mobile/features/fleet_catalog/presentation/pages/fleet_registration_page.dart';
import 'package:dink_mobile/l10n/app_localizations.dart';

class StubFleetRepository implements FleetRepository {
  @override
  Future<List<Category>> getCategories() async {
    return [
      const Category(
        id: 'cat-1',
        code: 'EXCAVATOR',
        name: 'Crawler Excavator',
        iconUrl: 'https://cdn.example.com/exc.png',
      ),
    ];
  }

  @override
  Future<List<SpecFieldGroup>> getSpecDefinitions(String categoryId) async {
    return [
      const SpecFieldGroup(
        id: 'grp-1',
        categoryId: 'cat-1',
        sectionNameEN: 'Technical Weights',
        sectionNameAM: 'የቴክኒክ ክብደት',
        displayOrder: 1,
        fields: [
          SpecFieldDefinition(
            id: 'def-1',
            fieldKey: 'operating_weight_kg',
            labelEN: 'Operating Weight',
            labelAM: 'የስራ ክብደት',
            fieldType: SpecFieldType.decimal,
            unitSymbol: 'kg',
            isRequired: true,
            isLogisticsCritical: true,
          ),
        ],
      ),
    ];
  }

  @override
  Future<FleetAsset> registerFleet(FleetRegistrationRequest request) async {
    return FleetAsset(
      id: 'fleet-1',
      ownerId: 'owner-1',
      categoryId: request.categoryId,
      make: request.make,
      model: request.model,
      manufacturingYear: request.manufacturingYear,
      registrationNumber: request.registrationNumber,
      dailyRateMinor: request.dailyRateMinor,
      status: 'PENDING_APPROVAL',
      latitude: request.latitude,
      longitude: request.longitude,
      addressText: request.addressText,
      version: 1,
    );
  }

  @override
  Future<void> saveDraft(FleetRegistrationRequest draft) async {}

  @override
  Future<FleetRegistrationRequest?> getDraft(String categoryId) async => null;
}

Widget buildTestableWidget({
  required Widget child,
  Locale locale = const Locale('en'),
}) {
  final repo = StubFleetRepository();
  final bloc = FleetRegistrationBloc(
    getCategoriesUseCase: GetCategoriesUseCase(repo),
    getSpecDefinitionsUseCase: GetSpecDefinitionsUseCase(repo),
    registerFleetUseCase: RegisterFleetUseCase(repo),
    fleetRepository: repo,
  );

  return BlocProvider<FleetRegistrationBloc>.value(
    value: bloc,
    child: MaterialApp(
      locale: locale,
      localizationsDelegates: const [
        AppLocalizations.delegate,
        GlobalMaterialLocalizations.delegate,
        GlobalWidgetsLocalizations.delegate,
        GlobalCupertinoLocalizations.delegate,
      ],
      supportedLocales: const [
        Locale('en'),
        Locale('am'),
      ],
      home: child,
    ),
  );
}

void main() {
  testWidgets('FleetRegistrationPage renders in English with zero hardcoded strings',
      (tester) async {
    await tester.pumpWidget(
      buildTestableWidget(
        child: const FleetRegistrationPage(),
        locale: const Locale('en'),
      ),
    );

    // Initial load pump
    await tester.pumpAndSettle();

    // Verify AppBar title
    expect(find.text('Register Machinery'), findsOneWidget);

    // Verify core sections and fields
    expect(find.text('Machinery Details'), findsOneWidget);
    expect(find.text('Walkaround 6-Angle Photos'), findsOneWidget);
    expect(find.text('Submit for Verification'), findsOneWidget);
    expect(find.text('Save Offline Draft'), findsOneWidget);
  });

  testWidgets('FleetRegistrationPage renders fully localized in Amharic',
      (tester) async {
    await tester.pumpWidget(
      buildTestableWidget(
        child: const FleetRegistrationPage(),
        locale: const Locale('am'),
      ),
    );

    // Initial load pump
    await tester.pumpAndSettle();

    // Verify localized Amharic AppBar title
    expect(find.text('ማሽነሪ መመዝገቢያ'), findsOneWidget);

    // Verify localized section headings and actions
    expect(find.text('የማሽነሪው ዝርዝር መረጃ'), findsOneWidget);
    expect(find.text('የማሽነሪው 6 አቅጣጫ ፎቶዎች'), findsOneWidget);
    expect(find.text('ማሽነሪውን ለምርመራ አስገባ'), findsOneWidget);
    expect(find.text('ረቂቁን በስልኩ አስቀምጥ'), findsOneWidget);
  });
}

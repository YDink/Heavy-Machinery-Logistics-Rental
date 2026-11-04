import 'dart:convert';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:dink_mobile/features/fleet_catalog/domain/entities/fleet_asset.dart';

abstract class FleetLocalDataSource {
  Future<void> saveDraft(FleetRegistrationRequest draft);
  Future<FleetRegistrationRequest?> getDraft(String categoryId);
  Future<void> clearDraft(String categoryId);
}

class FleetLocalDataSourceImpl implements FleetLocalDataSource {
  final FlutterSecureStorage storage;

  FleetLocalDataSourceImpl({required this.storage});

  static const String _draftPrefix = 'fleet_draft_';

  @override
  Future<void> saveDraft(FleetRegistrationRequest draft) async {
    final key = '$_draftPrefix${draft.categoryId}';
    final map = {
      'categoryId': draft.categoryId,
      'make': draft.make,
      'model': draft.model,
      'manufacturingYear': draft.manufacturingYear,
      'registrationNumber': draft.registrationNumber,
      'chassisVin': draft.chassisVin,
      'baselineMeterHours': draft.baselineMeterHours,
      'fuelType': draft.fuelType,
      'selfMoveCapable': draft.selfMoveCapable,
      'dailyRateMinor': draft.dailyRateMinor,
      'excessHourRateMinor': draft.excessHourRateMinor,
      'minRentalDays': draft.minRentalDays,
      'paymentTimeoutHours': draft.paymentTimeoutHours,
      'operatorIncluded': draft.operatorIncluded,
      'latitude': draft.latitude,
      'longitude': draft.longitude,
      'addressText': draft.addressText,
      'photos': draft.photos
          .map((p) => {
                'angle': p.angle.code,
                'storageUri': p.storageUri,
                'fileHash': p.fileHash,
                'localFilePath': p.localFilePath,
              })
          .toList(),
      'specValues': draft.specValues
          .map((v) => {
                'fieldKey': v.fieldKey,
                'valueString': v.valueString,
                'valueNumeric': v.valueNumeric,
                'valueBoolean': v.valueBoolean,
              })
          .toList(),
    };
    await storage.write(key: key, value: jsonEncode(map));
  }

  @override
  Future<FleetRegistrationRequest?> getDraft(String categoryId) async {
    final key = '$_draftPrefix$categoryId';
    final jsonStr = await storage.read(key: key);
    if (jsonStr == null || jsonStr.isEmpty) return null;

    final map = jsonDecode(jsonStr) as Map<String, dynamic>;
    final photosRaw = (map['photos'] as List<dynamic>?) ?? [];
    final photos = photosRaw.map((p) {
      final angle =
          PhotoAngle.fromCode(p['angle'] as String) ?? PhotoAngle.front;
      return FleetPhoto(
        angle: angle,
        storageUri: p['storageUri'] as String? ?? '',
        fileHash: p['fileHash'] as String? ?? '',
        localFilePath: p['localFilePath'] as String?,
      );
    }).toList();

    final specsRaw = (map['specValues'] as List<dynamic>?) ?? [];
    final specValues = specsRaw.map((s) {
      return FleetSpecValueInput(
        fieldKey: s['fieldKey'] as String,
        valueString: s['valueString'] as String?,
        valueNumeric: (s['valueNumeric'] as num?)?.toDouble(),
        valueBoolean: s['valueBoolean'] as bool?,
      );
    }).toList();

    return FleetRegistrationRequest(
      categoryId: map['categoryId'] as String,
      make: map['make'] as String? ?? '',
      model: map['model'] as String? ?? '',
      manufacturingYear: map['manufacturingYear'] as int? ?? 2020,
      registrationNumber: map['registrationNumber'] as String? ?? '',
      chassisVin: map['chassisVin'] as String? ?? '',
      baselineMeterHours: (map['baselineMeterHours'] as num?)?.toDouble() ?? 0.0,
      fuelType: map['fuelType'] as String? ?? 'DIESEL',
      selfMoveCapable: map['selfMoveCapable'] as bool? ?? false,
      dailyRateMinor: map['dailyRateMinor'] as int? ?? 0,
      excessHourRateMinor: map['excessHourRateMinor'] as int? ?? 0,
      minRentalDays: map['minRentalDays'] as int? ?? 1,
      paymentTimeoutHours: map['paymentTimeoutHours'] as int? ?? 24,
      operatorIncluded: map['operatorIncluded'] as bool? ?? false,
      latitude: (map['latitude'] as num?)?.toDouble() ?? 0.0,
      longitude: (map['longitude'] as num?)?.toDouble() ?? 0.0,
      addressText: map['addressText'] as String? ?? '',
      photos: photos,
      specValues: specValues,
    );
  }

  @override
  Future<void> clearDraft(String categoryId) async {
    final key = '$_draftPrefix$categoryId';
    await storage.delete(key: key);
  }
}

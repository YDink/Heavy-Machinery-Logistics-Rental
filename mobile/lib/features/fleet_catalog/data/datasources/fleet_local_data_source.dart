import 'dart:convert';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import '../models/fleet_registration_request_model.dart';

abstract class FleetLocalDataSource {
  Future<void> saveDraft(FleetRegistrationRequestModel draft);
  Future<FleetRegistrationRequestModel?> getDraft(String categoryId);
  Future<void> clearDraft(String categoryId);
}

class FleetLocalDataSourceImpl implements FleetLocalDataSource {
  final FlutterSecureStorage storage;

  FleetLocalDataSourceImpl({required this.storage});

  static const String _draftPrefix = 'fleet_draft_';

  @override
  Future<void> saveDraft(FleetRegistrationRequestModel draft) async {
    final key = '$_draftPrefix${draft.categoryId}';
    await storage.write(key: key, value: jsonEncode(draft.toJson()));
  }

  @override
  Future<FleetRegistrationRequestModel?> getDraft(String categoryId) async {
    final key = '$_draftPrefix$categoryId';
    final jsonString = await storage.read(key: key);
    if (jsonString == null) return null;

    try {
      final map = jsonDecode(jsonString) as Map<String, dynamic>;
      return FleetRegistrationRequestModel.fromJson(map);
    } catch (_) {
      return null;
    }
  }

  @override
  Future<void> clearDraft(String categoryId) async {
    final key = '$_draftPrefix$categoryId';
    await storage.delete(key: key);
  }
}

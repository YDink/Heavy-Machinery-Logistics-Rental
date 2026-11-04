import 'dart:convert';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import '../models/attachment_registration_request_model.dart';

abstract class AttachmentLocalDataSource {
  Future<void> saveDraft(AttachmentRegistrationRequestModel draft);
  Future<AttachmentRegistrationRequestModel?> getDraft();
  Future<void> clearDraft();
}

class AttachmentLocalDataSourceImpl implements AttachmentLocalDataSource {
  final FlutterSecureStorage storage;

  AttachmentLocalDataSourceImpl({required this.storage});

  static const String _draftKey = 'attachment_registration_draft';

  @override
  Future<void> saveDraft(AttachmentRegistrationRequestModel draft) async {
    await storage.write(key: _draftKey, value: jsonEncode(draft.toJson()));
  }

  @override
  Future<AttachmentRegistrationRequestModel?> getDraft() async {
    final jsonString = await storage.read(key: _draftKey);
    if (jsonString == null) return null;

    try {
      final map = jsonDecode(jsonString) as Map<String, dynamic>;
      return AttachmentRegistrationRequestModel.fromJson(map);
    } catch (_) {
      return null;
    }
  }

  @override
  Future<void> clearDraft() async {
    await storage.delete(key: _draftKey);
  }
}

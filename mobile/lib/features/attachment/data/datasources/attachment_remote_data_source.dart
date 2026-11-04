import 'package:dio/dio.dart';
import 'package:uuid/uuid.dart';
import 'package:dink_mobile/core/network/api_client.dart';
import '../models/attachment_model.dart';
import '../models/attachment_inventory_item_model.dart';
import '../models/bundled_attachment_link_model.dart';
import '../models/attachment_registration_request_model.dart';

abstract class AttachmentRemoteDataSource {
  Future<AttachmentModel> registerAttachment(AttachmentRegistrationRequestModel request);
  Future<List<AttachmentInventoryItemModel>> getInventory({String? status});
  Future<BundledAttachmentLinkModel> linkAttachmentToFleet({
    required String fleetId,
    required String attachmentId,
    required int bundledDailyRateMinor,
  });
}

class AttachmentRemoteDataSourceImpl implements AttachmentRemoteDataSource {
  final ApiClient apiClient;

  AttachmentRemoteDataSourceImpl({required this.apiClient});

  @override
  Future<AttachmentModel> registerAttachment(AttachmentRegistrationRequestModel request) async {
    try {
      final response = await apiClient.dio.post(
        '/attachments',
        data: request.toJson(),
        options: Options(
          headers: {
            'Idempotency-Key': const Uuid().v4(),
          },
        ),
      );

      final data = response.data['data'] as Map<String, dynamic>;
      return AttachmentModel.fromJson(data);
    } on DioException catch (e) {
      throw apiClient.mapDioExceptionToFailure(e);
    }
  }

  @override
  Future<List<AttachmentInventoryItemModel>> getInventory({String? status}) async {
    try {
      final queryParams = <String, dynamic>{};
      if (status != null && status.isNotEmpty) {
        queryParams['status'] = status;
      }

      final response = await apiClient.dio.get(
        '/attachments/inventory',
        queryParameters: queryParams,
      );

      final data = response.data['data'] as Map<String, dynamic>;
      final list = data['attachments'] as List<dynamic>? ?? [];
      return list
          .map((item) => AttachmentInventoryItemModel.fromJson(item as Map<String, dynamic>))
          .toList();
    } on DioException catch (e) {
      throw apiClient.mapDioExceptionToFailure(e);
    }
  }

  @override
  Future<BundledAttachmentLinkModel> linkAttachmentToFleet({
    required String fleetId,
    required String attachmentId,
    required int bundledDailyRateMinor,
  }) async {
    try {
      final payload = {
        'attachmentId': attachmentId,
        'bundledDailyRateMinor': bundledDailyRateMinor,
      };

      final response = await apiClient.dio.post(
        '/fleets/$fleetId/attachments/link',
        data: payload,
      );

      final data = response.data['data'] as Map<String, dynamic>;
      return BundledAttachmentLinkModel.fromJson(data);
    } on DioException catch (e) {
      throw apiClient.mapDioExceptionToFailure(e);
    }
  }
}

import 'package:intl/intl.dart';
import 'package:dink_mobile/core/network/api_client.dart';
import '../models/booking_model.dart';

abstract class BookingRemoteDataSource {
  Future<BookingModel> createBooking({
    required String fleetId,
    required DateTime startDate,
    required DateTime endDate,
    required String siteLocationAddress,
    List<String> attachmentIds = const [],
  });

  Future<BookingModel> getBooking(String id);

  Future<List<BookingModel>> listBookings({String? status});

  Future<BankSlipModel> submitPaymentEvidence({
    required String bookingId,
    required String storageUri,
    required String fileHash,
    String? transactionReference,
    required int amountMinor,
    String currency = 'ETB',
  });

  Future<BookingModel> confirmPayment({
    required String bookingId,
    required int expectedVersion,
    String? paymentReference,
    String? notes,
  });

  Future<BookingModel> declinePayment({
    required String bookingId,
    required int expectedVersion,
    required String reasonCode,
    required String rejectionNotes,
  });

  Future<BookingModel> disputeBooking({
    required String bookingId,
    required int expectedVersion,
    required String disputeReason,
    List<String> evidenceUris = const [],
  });

  Future<BookingModel> cancelBooking({
    required String bookingId,
    required int expectedVersion,
    required String reason,
  });

  Future<BookingModel> handoverCheckIn({
    required String bookingId,
    required int expectedVersion,
    required String handoverPin,
    required double operatingHoursMeter,
    List<String> photoUris = const [],
  });

  Future<BookingModel> returnCheckOut({
    required String bookingId,
    required int expectedVersion,
    required double operatingHoursMeter,
    String? damageNotes,
    List<String> photoUris = const [],
  });
}

class BookingRemoteDataSourceImpl implements BookingRemoteDataSource {
  final ApiClient apiClient;

  BookingRemoteDataSourceImpl({required this.apiClient});

  final DateFormat _dateFormat = DateFormat('yyyy-MM-dd');

  @override
  Future<BookingModel> createBooking({
    required String fleetId,
    required DateTime startDate,
    required DateTime endDate,
    required String siteLocationAddress,
    List<String> attachmentIds = const [],
  }) async {
    final payload = <String, dynamic>{
      'fleetId': fleetId,
      'startDate': _dateFormat.format(startDate),
      'endDate': _dateFormat.format(endDate),
      'siteLocationAddress': siteLocationAddress,
    };
    if (attachmentIds.isNotEmpty) {
      payload['attachmentIds'] = attachmentIds;
    }

    final response = await apiClient.dio.post(
      '/bookings',
      data: payload,
    );
    final body = response.data as Map<String, dynamic>;
    final data = body['data'] as Map<String, dynamic>;
    return BookingModel.fromJson(data);
  }

  @override
  Future<BookingModel> getBooking(String id) async {
    final response = await apiClient.dio.get('/bookings/$id');
    final body = response.data as Map<String, dynamic>;
    final data = body['data'] as Map<String, dynamic>;
    return BookingModel.fromJson(data);
  }

  @override
  Future<List<BookingModel>> listBookings({String? status}) async {
    final queryParams = <String, dynamic>{};
    if (status != null && status.isNotEmpty) {
      queryParams['status'] = status;
    }

    final response = await apiClient.dio.get(
      '/bookings',
      queryParameters: queryParams,
    );
    final body = response.data as Map<String, dynamic>;
    final list = body['data'] as List<dynamic>;
    return list.map((e) => BookingModel.fromJson(e as Map<String, dynamic>)).toList();
  }

  @override
  Future<BankSlipModel> submitPaymentEvidence({
    required String bookingId,
    required String storageUri,
    required String fileHash,
    String? transactionReference,
    required int amountMinor,
    String currency = 'ETB',
  }) async {
    final payload = <String, dynamic>{
      'storageUri': storageUri,
      'fileHash': fileHash,
      'amountMinor': amountMinor,
      'currency': currency,
    };
    if (transactionReference != null && transactionReference.isNotEmpty) {
      payload['transactionReference'] = transactionReference;
    }

    final response = await apiClient.dio.post(
      '/bookings/$bookingId/payment-evidence',
      data: payload,
    );
    final body = response.data as Map<String, dynamic>;
    final data = body['data'] as Map<String, dynamic>;
    return BankSlipModel.fromJson(data);
  }

  @override
  Future<BookingModel> confirmPayment({
    required String bookingId,
    required int expectedVersion,
    String? paymentReference,
    String? notes,
  }) async {
    final payload = <String, dynamic>{
      'expectedVersion': expectedVersion,
    };
    if (paymentReference != null) payload['paymentReference'] = paymentReference;
    if (notes != null) payload['notes'] = notes;

    final response = await apiClient.dio.post(
      '/bookings/$bookingId:confirm-payment',
      data: payload,
    );
    final body = response.data as Map<String, dynamic>;
    final data = body['data'] as Map<String, dynamic>;
    return BookingModel.fromJson(data);
  }

  @override
  Future<BookingModel> declinePayment({
    required String bookingId,
    required int expectedVersion,
    required String reasonCode,
    required String rejectionNotes,
  }) async {
    final payload = <String, dynamic>{
      'expectedVersion': expectedVersion,
      'reasonCode': reasonCode,
      'rejectionNotes': rejectionNotes,
    };

    final response = await apiClient.dio.post(
      '/bookings/$bookingId:decline-payment',
      data: payload,
    );
    final body = response.data as Map<String, dynamic>;
    final data = body['data'] as Map<String, dynamic>;
    return BookingModel.fromJson(data);
  }

  @override
  Future<BookingModel> disputeBooking({
    required String bookingId,
    required int expectedVersion,
    required String disputeReason,
    List<String> evidenceUris = const [],
  }) async {
    final payload = <String, dynamic>{
      'expectedVersion': expectedVersion,
      'disputeReason': disputeReason,
      'evidenceUris': evidenceUris,
    };

    final response = await apiClient.dio.post(
      '/bookings/$bookingId:dispute',
      data: payload,
    );
    final body = response.data as Map<String, dynamic>;
    final data = body['data'] as Map<String, dynamic>;
    return BookingModel.fromJson(data);
  }

  @override
  Future<BookingModel> cancelBooking({
    required String bookingId,
    required int expectedVersion,
    required String reason,
  }) async {
    final payload = <String, dynamic>{
      'expectedVersion': expectedVersion,
      'reason': reason,
    };

    final response = await apiClient.dio.post(
      '/bookings/$bookingId:cancel',
      data: payload,
    );
    final body = response.data as Map<String, dynamic>;
    final data = body['data'] as Map<String, dynamic>;
    return BookingModel.fromJson(data);
  }

  @override
  Future<BookingModel> handoverCheckIn({
    required String bookingId,
    required int expectedVersion,
    required String handoverPin,
    required double operatingHoursMeter,
    List<String> photoUris = const [],
  }) async {
    final payload = <String, dynamic>{
      'expectedVersion': expectedVersion,
      'handoverPin': handoverPin,
      'operatingHoursMeter': operatingHoursMeter,
      'photoUris': photoUris,
    };

    final response = await apiClient.dio.post(
      '/bookings/$bookingId:handover-check-in',
      data: payload,
    );
    final body = response.data as Map<String, dynamic>;
    final data = body['data'] as Map<String, dynamic>;
    return BookingModel.fromJson(data);
  }

  @override
  Future<BookingModel> returnCheckOut({
    required String bookingId,
    required int expectedVersion,
    required double operatingHoursMeter,
    String? damageNotes,
    List<String> photoUris = const [],
  }) async {
    final payload = <String, dynamic>{
      'expectedVersion': expectedVersion,
      'operatingHoursMeter': operatingHoursMeter,
      'photoUris': photoUris,
    };
    if (damageNotes != null && damageNotes.isNotEmpty) {
      payload['damageNotes'] = damageNotes;
    }

    final response = await apiClient.dio.post(
      '/bookings/$bookingId:return-check-out',
      data: payload,
    );
    final body = response.data as Map<String, dynamic>;
    final data = body['data'] as Map<String, dynamic>;
    return BookingModel.fromJson(data);
  }
}

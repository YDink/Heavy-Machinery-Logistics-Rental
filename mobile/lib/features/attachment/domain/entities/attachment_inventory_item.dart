class AttachmentInventoryItem {
  final String id;
  final String attachmentType;
  final String make;
  final String model;
  final String subTypeSize;
  final String status;
  final int standaloneDailyRateMinor;
  final bool isLinkedToFleet;
  final String? linkedFleetId;
  final String? primaryPhotoUri;

  const AttachmentInventoryItem({
    required this.id,
    required this.attachmentType,
    required this.make,
    required this.model,
    required this.subTypeSize,
    required this.status,
    required this.standaloneDailyRateMinor,
    required this.isLinkedToFleet,
    this.linkedFleetId,
    this.primaryPhotoUri,
  });
}

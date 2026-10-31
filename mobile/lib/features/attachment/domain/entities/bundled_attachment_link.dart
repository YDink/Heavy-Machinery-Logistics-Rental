class BundledAttachmentLink {
  final String fleetId;
  final String attachmentId;
  final int bundledDailyRateMinor;
  final double discountPercentage;
  final String status;

  const BundledAttachmentLink({
    required this.fleetId,
    required this.attachmentId,
    required this.bundledDailyRateMinor,
    required this.discountPercentage,
    required this.status,
  });
}

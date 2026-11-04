enum AttachmentType {
  hydraulicBreaker('HYDRAULIC_BREAKER', 'Hydraulic Breaker'),
  standardBucket('STANDARD_BUCKET', 'Standard Bucket'),
  trenchingBucket('TRENCHING_BUCKET', 'Trenching Bucket'),
  auger('AUGER', 'Earth Auger'),
  ripperTooth('RIPPER_TOOTH', 'Ripper Tooth'),
  sortingGrapple('SORTING_GRAPPLE', 'Sorting Grapple'),
  crusher('CRUSHER', 'Demolition Crusher');

  final String code;
  final String displayName;
  const AttachmentType(this.code, this.displayName);

  static AttachmentType fromCode(String code) {
    return AttachmentType.values.firstWhere(
      (t) => t.code == code,
      orElse: () => AttachmentType.standardBucket,
    );
  }
}

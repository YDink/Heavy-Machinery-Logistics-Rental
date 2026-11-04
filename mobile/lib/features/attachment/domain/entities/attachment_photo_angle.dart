enum AttachmentPhotoAngle {
  overall('OVERALL', 'Overall Profile'),
  bracket('BRACKET', 'Mounting Bracket'),
  pins('PINS', 'Pins & Bushings'),
  toolTip('TOOL_TIP', 'Working Tip / Chisel');

  final String code;
  final String label;
  const AttachmentPhotoAngle(this.code, this.label);

  static AttachmentPhotoAngle fromCode(String code) {
    return AttachmentPhotoAngle.values.firstWhere(
      (a) => a.code == code,
      orElse: () => AttachmentPhotoAngle.overall,
    );
  }
}

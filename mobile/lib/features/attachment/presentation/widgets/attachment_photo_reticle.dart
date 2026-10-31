import 'package:flutter/material.dart';
import 'package:dink_mobile/features/attachment/domain/entities/attachment.dart';
import 'package:dink_mobile/l10n/app_localizations.dart';

class AttachmentPhotoReticle extends StatelessWidget {
  final Map<AttachmentPhotoAngle, AttachmentPhoto> capturedPhotos;
  final ValueChanged<AttachmentPhotoAngle> onCapturePhoto;

  const AttachmentPhotoReticle({
    super.key,
    required this.capturedPhotos,
    required this.onCapturePhoto,
  });

  String _getAngleLabel(AttachmentPhotoAngle angle, AppLocalizations l10n) {
    switch (angle) {
      case AttachmentPhotoAngle.overall:
        return l10n.photoAngleOverall;
      case AttachmentPhotoAngle.bracket:
        return l10n.photoAngleBracket;
      case AttachmentPhotoAngle.pins:
        return l10n.photoAnglePins;
      case AttachmentPhotoAngle.toolTip:
        return l10n.photoAngleToolTip;
    }
  }

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            Text(
              l10n.attachmentPhotosSection,
              style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 16),
            ),
            Text(
              l10n.photosCapturedCount(capturedPhotos.length, 4),
              style: TextStyle(
                fontWeight: FontWeight.bold,
                color: capturedPhotos.length == 4 ? Colors.green : Colors.orange,
              ),
            ),
          ],
        ),
        const SizedBox(height: 12),
        GridView.count(
          crossAxisCount: 2,
          shrinkWrap: true,
          physics: const NeverScrollableScrollPhysics(),
          mainAxisSpacing: 10,
          crossAxisSpacing: 10,
          childAspectRatio: 1.1,
          children: AttachmentPhotoAngle.values.map((angle) {
            final isCaptured = capturedPhotos.containsKey(angle);
            return _buildAngleCard(context, angle, isCaptured, l10n);
          }).toList(),
        ),
      ],
    );
  }

  Widget _buildAngleCard(
    BuildContext context,
    AttachmentPhotoAngle angle,
    bool isCaptured,
    AppLocalizations l10n,
  ) {
    return InkWell(
      onTap: () => onCapturePhoto(angle),
      borderRadius: BorderRadius.circular(10),
      child: Container(
        decoration: BoxDecoration(
          border: Border.all(
            color: isCaptured ? Colors.green : Colors.grey.shade400,
            width: isCaptured ? 2.0 : 1.0,
          ),
          borderRadius: BorderRadius.circular(10),
          color: isCaptured ? Colors.green.shade50 : Colors.grey.shade100,
        ),
        padding: const EdgeInsets.all(8),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(
              isCaptured ? Icons.check_circle : Icons.camera_alt_outlined,
              color: isCaptured ? Colors.green : Colors.blueGrey,
              size: 32,
            ),
            const SizedBox(height: 8),
            Text(
              _getAngleLabel(angle, l10n),
              textAlign: TextAlign.center,
              style: const TextStyle(fontWeight: FontWeight.w600, fontSize: 13),
            ),
            const SizedBox(height: 4),
            Text(
              isCaptured ? l10n.verifiedClean : l10n.tapToCapture,
              style: TextStyle(
                fontSize: 11,
                color: isCaptured ? Colors.green.shade700 : Colors.grey.shade600,
              ),
            ),
          ],
        ),
      ),
    );
  }
}

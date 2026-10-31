import 'package:flutter/material.dart';
import 'package:dink_mobile/features/fleet_catalog/domain/entities/fleet_asset.dart';
import 'package:dink_mobile/l10n/app_localizations.dart';

class WalkaroundCameraReticle extends StatelessWidget {
  final Map<PhotoAngle, FleetPhoto> photos;
  final ValueChanged<PhotoAngle> onAngleSelected;

  const WalkaroundCameraReticle({
    super.key,
    required this.photos,
    required this.onAngleSelected,
  });

  String _getAngleLabel(BuildContext context, PhotoAngle angle) {
    final l10n = AppLocalizations.of(context)!;
    switch (angle) {
      case PhotoAngle.front:
        return l10n.photoFront;
      case PhotoAngle.rear:
        return l10n.photoRear;
      case PhotoAngle.left:
        return l10n.photoLeft;
      case PhotoAngle.right:
        return l10n.photoRight;
      case PhotoAngle.cabin:
        return l10n.photoCabin;
      case PhotoAngle.meter:
        return l10n.photoMeter;
    }
  }

  IconData _getAngleIcon(PhotoAngle angle) {
    switch (angle) {
      case PhotoAngle.front:
        return Icons.arrow_upward;
      case PhotoAngle.rear:
        return Icons.arrow_downward;
      case PhotoAngle.left:
        return Icons.arrow_back;
      case PhotoAngle.right:
        return Icons.arrow_forward;
      case PhotoAngle.cabin:
        return Icons.airline_seat_recline_normal;
      case PhotoAngle.meter:
        return Icons.speed;
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final l10n = AppLocalizations.of(context)!;

    final capturedCount = photos.length;
    final totalCount = PhotoAngle.values.length;
    final isComplete = capturedCount == totalCount;

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            Text(
              l10n.walkaroundPhotos,
              style: theme.textTheme.titleMedium?.copyWith(
                fontWeight: FontWeight.bold,
              ),
            ),
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
              decoration: BoxDecoration(
                color: isComplete
                    ? Colors.green.withValues(alpha: 0.15)
                    : Colors.amber.withValues(alpha: 0.15),
                borderRadius: BorderRadius.circular(12),
              ),
              child: Text(
                '$capturedCount / $totalCount',
                style: TextStyle(
                  color: isComplete ? Colors.green.shade800 : Colors.amber.shade900,
                  fontWeight: FontWeight.w700,
                  fontSize: 12,
                ),
              ),
            ),
          ],
        ),
        const SizedBox(height: 12),
        GridView.builder(
          shrinkWrap: true,
          physics: const NeverScrollableScrollPhysics(),
          gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
            crossAxisCount: 3,
            crossAxisSpacing: 10,
            mainAxisSpacing: 10,
            childAspectRatio: 0.95,
          ),
          itemCount: PhotoAngle.values.length,
          itemBuilder: (context, index) {
            final angle = PhotoAngle.values[index];
            final photo = photos[angle];
            final isCaptured = photo != null;

            return InkWell(
              onTap: () => onAngleSelected(angle),
              borderRadius: BorderRadius.circular(12),
              child: Container(
                decoration: BoxDecoration(
                  color: isCaptured
                      ? theme.colorScheme.primaryContainer.withValues(alpha: 0.3)
                      : theme.colorScheme.surfaceContainerHighest.withValues(alpha: 0.4),
                  border: Border.all(
                    color: isCaptured
                        ? theme.colorScheme.primary
                        : theme.colorScheme.outlineVariant,
                    width: isCaptured ? 2 : 1,
                  ),
                  borderRadius: BorderRadius.circular(12),
                ),
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Stack(
                      alignment: Alignment.center,
                      children: [
                        CircleAvatar(
                          radius: 22,
                          backgroundColor: isCaptured
                              ? theme.colorScheme.primary
                              : theme.colorScheme.surface,
                          child: Icon(
                            _getAngleIcon(angle),
                            size: 20,
                            color: isCaptured
                                ? theme.colorScheme.onPrimary
                                : theme.colorScheme.onSurfaceVariant,
                          ),
                        ),
                        if (isCaptured)
                          const Positioned(
                            right: 0,
                            bottom: 0,
                            child: CircleAvatar(
                              radius: 8,
                              backgroundColor: Colors.green,
                              child: Icon(Icons.check, size: 10, color: Colors.white),
                            ),
                          ),
                      ],
                    ),
                    const SizedBox(height: 8),
                    Text(
                      _getAngleLabel(context, angle),
                      textAlign: TextAlign.center,
                      maxLines: 2,
                      overflow: TextOverflow.ellipsis,
                      style: theme.textTheme.bodySmall?.copyWith(
                        fontWeight: isCaptured ? FontWeight.w600 : FontWeight.normal,
                      ),
                    ),
                  ],
                ),
              ),
            );
          },
        ),
      ],
    );
  }
}

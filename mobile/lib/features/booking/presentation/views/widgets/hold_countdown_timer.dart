import 'dart:async';
import 'package:flutter/material.dart';
import 'package:dink_mobile/l10n/app_localizations.dart';

class HoldCountdownTimer extends StatefulWidget {
  final DateTime expiresAt;
  final VoidCallback? onExpired;

  const HoldCountdownTimer({
    super.key,
    required this.expiresAt,
    this.onExpired,
  });

  @override
  State<HoldCountdownTimer> createState() => _HoldCountdownTimerState();
}

class _HoldCountdownTimerState extends State<HoldCountdownTimer> {
  Timer? _timer;
  Duration _remaining = Duration.zero;

  @override
  void initState() {
    super.initState();
    _updateRemaining();
    if (_remaining > Duration.zero) {
      _timer = Timer.periodic(const Duration(seconds: 1), (_) {
        _updateRemaining();
      });
    }
  }

  void _updateRemaining() {
    final now = DateTime.now();
    if (now.isAfter(widget.expiresAt)) {
      setState(() {
        _remaining = Duration.zero;
      });
      _timer?.cancel();
      widget.onExpired?.call();
    } else {
      setState(() {
        _remaining = widget.expiresAt.difference(now);
      });
    }
  }

  @override
  void dispose() {
    _timer?.cancel();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;
    final theme = Theme.of(context);

    final minutes = _remaining.inMinutes.remainder(60).toString().padLeft(2, '0');
    final seconds = _remaining.inSeconds.remainder(60).toString().padLeft(2, '0');

    final isUrgent = _remaining.inMinutes < 15;

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 8),
      decoration: BoxDecoration(
        color: isUrgent
            ? theme.colorScheme.errorContainer.withValues(alpha: 0.5)
            : theme.colorScheme.secondaryContainer.withValues(alpha: 0.5),
        borderRadius: BorderRadius.circular(20),
        border: Border.all(
          color: isUrgent ? theme.colorScheme.error : theme.colorScheme.secondary,
        ),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(
            Icons.timer_outlined,
            size: 18,
            color: isUrgent ? theme.colorScheme.error : theme.colorScheme.secondary,
          ),
          const SizedBox(width: 6),
          Text(
            '${l10n.paymentHoldTimer}: $minutes:$seconds',
            style: theme.textTheme.labelMedium?.copyWith(
              fontWeight: FontWeight.bold,
              color: isUrgent ? theme.colorScheme.error : theme.colorScheme.secondary,
            ),
          ),
        ],
      ),
    );
  }
}

import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../../bloc/rental_completion_bloc.dart';
import '../../bloc/rental_completion_event.dart';
import '../../bloc/rental_completion_state.dart';

class BookingReviewModal extends StatefulWidget {
  final String bookingId;
  final String counterpartyName;
  final String role;

  const BookingReviewModal({
    super.key,
    required this.bookingId,
    this.counterpartyName = 'Almaz Heavy Equipment Fleet PLC',
    this.role = 'OWNER',
  });

  @override
  State<BookingReviewModal> createState() => _BookingReviewModalState();
}

class _BookingReviewModalState extends State<BookingReviewModal> {
  double _overallRating = 5.0;
  final _feedbackController = TextEditingController(text: 'Flawless site coordination, powerful excavator in certified working order.');
  double _dim1 = 5.0;
  double _dim2 = 4.8;
  double _dim3 = 5.0;

  void _submit() {
    final dims = {
      'EQUIPMENT_CONDITION': _dim1,
      'OPERATIONAL_PUNCTUALITY': _dim2,
      'SITE_COORDINATION': _dim3,
    };

    context.read<RentalCompletionBloc>().add(
          SubmitBookingReviewEvent(
            bookingId: widget.bookingId,
            rating: _overallRating,
            feedback: _feedbackController.text.trim(),
            dimensionRatings: dims,
          ),
        );
  }

  @override
  void dispose() {
    _feedbackController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    const bgDark = Color(0xFF111827);
    const cardDark = Color(0xFF1C1C1C);
    const primaryYellow = Color(0xFFFFB300);
    const borderGray = Color(0xFF393939);
    const successGreen = Color(0xFF10B981);

    return BlocConsumer<RentalCompletionBloc, RentalCompletionState>(
      listener: (context, state) {
        if (state is BookingReviewSubmittedState) {
          final isUnblinded = state.review.isPublished;
          Navigator.of(context).pop();
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text(isUnblinded
                  ? 'Bilateral review complete! Both reviews published to public profiles.'
                  : 'Review submitted! Stored in double-blind isolation until counterparty rates.'),
              backgroundColor: successGreen,
            ),
          );
        } else if (state is RentalCompletionError) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text(state.message),
              backgroundColor: const Color(0xFFEF4444),
            ),
          );
        }
      },
      builder: (context, state) {
        final isLoading = state is RentalCompletionLoading;

        return Material(
          color: bgDark,
          child: Container(
            padding: const EdgeInsets.all(20.0),
            child: SingleChildScrollView(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    const Text(
                      'MUTUAL PERFORMANCE REVIEW',
                      style: TextStyle(
                        fontFamily: 'Space Grotesk',
                        fontSize: 15,
                        fontWeight: FontWeight.w700,
                        color: Colors.white,
                      ),
                    ),
                    IconButton(
                      icon: const Icon(Icons.close, color: Colors.white70),
                      onPressed: () => Navigator.of(context).pop(),
                    ),
                  ],
                ),
                Text(
                  'Counterparty: ${widget.counterpartyName}',
                  style: const TextStyle(fontSize: 12, color: Colors.white60),
                ),
                const SizedBox(height: 14),

                // Double-Blind Guarantee Banner
                Container(
                  width: double.infinity,
                  padding: const EdgeInsets.all(12),
                  decoration: BoxDecoration(
                    color: const Color(0xFF1E293B),
                    border: Border.all(color: const Color(0xFF3B82F6)),
                    borderRadius: BorderRadius.zero,
                  ),
                  child: const Row(
                    children: [
                      Icon(Icons.lock_outline, size: 18, color: Color(0xFF60A5FA)),
                      SizedBox(width: 8),
                      Expanded(
                        child: Text(
                          'DOUBLE-BLIND GUARANTEE (RC-07): Your review remains hidden until both parties rate or the 14-day window expires. No retaliatory reviews.',
                          style: TextStyle(fontSize: 11, color: Colors.white, height: 1.3),
                        ),
                      ),
                    ],
                  ),
                ),
                const SizedBox(height: 20),

                // Overall Star Rating
                const Text(
                  'OVERALL STAR RATING',
                  style: TextStyle(
                    fontFamily: 'Space Grotesk',
                    fontSize: 12,
                    fontWeight: FontWeight.w700,
                    color: Colors.white,
                  ),
                ),
                const SizedBox(height: 8),
                Row(
                  children: List.generate(5, (index) {
                    final starVal = index + 1.0;
                    return IconButton(
                      iconSize: 32,
                      icon: Icon(
                        starVal <= _overallRating ? Icons.star : Icons.star_border,
                        color: primaryYellow,
                      ),
                      onPressed: () => setState(() => _overallRating = starVal),
                    );
                  }),
                ),
                const SizedBox(height: 16),

                // Dimensional Sliders
                _buildDimensionSlider('Equipment Operational Condition', _dim1, (v) => setState(() => _dim1 = v)),
                _buildDimensionSlider('Punctuality & Schedule Adherence', _dim2, (v) => setState(() => _dim2 = v)),
                _buildDimensionSlider('Site Safety & Communications', _dim3, (v) => setState(() => _dim3 = v)),
                const SizedBox(height: 16),

                // Feedback Field
                TextField(
                  controller: _feedbackController,
                  maxLines: 3,
                  style: const TextStyle(color: Colors.white, fontSize: 13),
                  decoration: const InputDecoration(
                    labelText: 'CONFIDENTIAL OPERATIONAL FEEDBACK',
                    labelStyle: TextStyle(color: Colors.white70),
                    filled: true,
                    fillColor: cardDark,
                    border: OutlineInputBorder(
                      borderRadius: BorderRadius.zero,
                      borderSide: BorderSide(color: borderGray),
                    ),
                  ),
                ),
                const SizedBox(height: 24),

                // 56dp Submit Button
                SizedBox(
                  width: double.infinity,
                  height: 56,
                  child: ElevatedButton(
                    onPressed: isLoading ? null : _submit,
                    style: ElevatedButton.styleFrom(
                      backgroundColor: primaryYellow,
                      foregroundColor: Colors.black,
                      shape: const RoundedRectangleBorder(borderRadius: BorderRadius.zero),
                    ),
                    child: isLoading
                        ? const CircularProgressIndicator(color: Colors.black)
                        : const Text(
                            'SUBMIT CONFIDENTIAL REVIEW →',
                            style: TextStyle(
                              fontFamily: 'Space Grotesk',
                              fontSize: 14,
                              fontWeight: FontWeight.w800,
                              letterSpacing: 0.5,
                            ),
                          ),
                  ),
                ),
              ],
            ),
          ),
        ),
      );
    },
  );
  }

  Widget _buildDimensionSlider(String label, double value, ValueChanged<double> onChanged) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            Text(label, style: const TextStyle(color: Colors.white70, fontSize: 11)),
            Text(
              '${value.toStringAsFixed(1)} ★',
              style: const TextStyle(
                fontFamily: 'JetBrains Mono',
                fontSize: 12,
                fontWeight: FontWeight.w700,
                color: Color(0xFFFFB300),
              ),
            ),
          ],
        ),
        SliderTheme(
          data: SliderTheme.of(context).copyWith(
            activeTrackColor: const Color(0xFFFFB300),
            inactiveTrackColor: const Color(0xFF393939),
            thumbColor: const Color(0xFFFFB300),
            overlayColor: const Color(0xFFFFB300).withOpacity(0.2),
          ),
          child: Slider(
            value: value,
            min: 1.0,
            max: 5.0,
            divisions: 8,
            onChanged: onChanged,
          ),
        ),
      ],
    );
  }
}

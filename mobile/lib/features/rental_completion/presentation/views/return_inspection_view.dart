import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../../domain/entities/return_photo_item.dart';
import '../bloc/rental_completion_bloc.dart';
import '../bloc/rental_completion_event.dart';
import '../bloc/rental_completion_state.dart';

class ReturnInspectionView extends StatefulWidget {
  final String bookingId;
  final String equipmentTitle;
  final double checkInMeterHours;

  const ReturnInspectionView({
    super.key,
    required this.bookingId,
    this.equipmentTitle = 'CAT 320D Hydraulic Excavator',
    this.checkInMeterHours = 1250.0,
  });

  @override
  State<ReturnInspectionView> createState() => _ReturnInspectionViewState();
}

class _ReturnInspectionViewState extends State<ReturnInspectionView> {
  final _meterController = TextEditingController(text: '1314.80');
  final _notesController = TextEditingController(text: 'Hydraulic cylinders dry; minor cosmetic scuff on right counterweight lower skirt within allowable ASTM tolerance.');
  bool _fluidLevelsOk = true;
  bool _safetyGlassIntact = true;
  bool _accelerometerSteady = true;

  final Map<String, String> _capturedPhotos = {
    'FRONT_IMPLEMENT': 's3://photos/front_implement.jpg',
    'REAR_COUNTERWEIGHT_ENGINE': 's3://photos/rear_counterweight.jpg',
    'LEFT_SIDE_TRACKS': 's3://photos/left_tracks.jpg',
    'RIGHT_SIDE_TRACKS': 's3://photos/right_tracks.jpg',
    'HOUR_METER': 's3://photos/hour_meter.jpg',
    'CABIN_INTERIOR': 's3://photos/cabin_interior.jpg',
  };

  static const _angles = [
    {'key': 'FRONT_IMPLEMENT', 'title': '01. FRONT IMPLEMENT', 'desc': 'Boom, Stick & Bucket Teeth'},
    {'key': 'REAR_COUNTERWEIGHT_ENGINE', 'title': '02. REAR COUNTERWEIGHT', 'desc': 'Hood & Radiator Grille'},
    {'key': 'LEFT_SIDE_TRACKS', 'title': '03. LEFT SIDE TRACKS', 'desc': 'Sprockets & Grouser Shoes'},
    {'key': 'RIGHT_SIDE_TRACKS', 'title': '04. RIGHT SIDE TRACKS', 'desc': 'Rollers & Idler Wheels'},
    {'key': 'HOUR_METER', 'title': '05. HOUR METER DIAL', 'desc': 'Macro Drum Numbers & Glass'},
    {'key': 'CABIN_INTERIOR', 'title': '06. CABIN INTERIOR', 'desc': 'ROPS/FOPS Glass & Joysticks'},
  ];

  double get _currentCheckoutMeter => double.tryParse(_meterController.text) ?? 0.0;
  bool get _isMonotonicValid => _currentCheckoutMeter >= widget.checkInMeterHours;
  double get _hoursDelta => _currentCheckoutMeter - widget.checkInMeterHours;

  void _submit() {
    if (!_isMonotonicValid) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('Return meter reading cannot regress below certified check-in reading.'),
          backgroundColor: Color(0xFFEF4444),
        ),
      );
      return;
    }

    final photos = _capturedPhotos.entries
        .map((e) => ReturnPhotoItem(angle: e.key, storageUri: e.value))
        .toList();

    context.read<RentalCompletionBloc>().add(
          SubmitReturnInspectionEvent(
            bookingId: widget.bookingId,
            checkOutMeterHours: _currentCheckoutMeter,
            fluidLevelsOk: _fluidLevelsOk,
            safetyGlassIntact: _safetyGlassIntact,
            accelerometerSteady: _accelerometerSteady,
            notes: _notesController.text.trim(),
            photos: photos,
          ),
        );
  }

  @override
  void dispose() {
    _meterController.dispose();
    _notesController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    const bgDark = Color(0xFF111827);
    const cardDark = Color(0xFF1C1C1C);
    const primaryYellow = Color(0xFFFFB300);
    const borderGray = Color(0xFF393939);
    const successGreen = Color(0xFF10B981);

    return Scaffold(
      backgroundColor: bgDark,
      appBar: AppBar(
        backgroundColor: bgDark,
        elevation: 0,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back, color: Colors.white),
          onPressed: () => Navigator.of(context).pop(),
        ),
        title: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text(
              'RETURN INSPECTION',
              style: TextStyle(
                fontFamily: 'Space Grotesk',
                fontSize: 16,
                fontWeight: FontWeight.w700,
                color: Colors.white,
                letterSpacing: 0.5,
              ),
            ),
            Text(
              widget.equipmentTitle,
              style: const TextStyle(
                fontSize: 12,
                color: Colors.white70,
              ),
            ),
          ],
        ),
        actions: [
          Container(
            margin: const EdgeInsets.only(right: 16, top: 12, bottom: 12),
            padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
            decoration: BoxDecoration(
              border: Border.all(color: primaryYellow, width: 1.5),
              borderRadius: BorderRadius.zero,
            ),
            child: const Center(
              child: Text(
                'RC-02 AUDIT',
                style: TextStyle(
                  fontFamily: 'JetBrains Mono',
                  fontSize: 11,
                  fontWeight: FontWeight.w700,
                  color: primaryYellow,
                ),
              ),
            ),
          ),
        ],
      ),
      body: BlocConsumer<RentalCompletionBloc, RentalCompletionState>(
        listener: (context, state) {
          if (state is ReturnInspectionSuccessState) {
            ScaffoldMessenger.of(context).showSnackBar(
              const SnackBar(
                content: Text('Return inspection submitted! 24-hour damage review window opened.'),
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

          return SingleChildScrollView(
            padding: const EdgeInsets.all(16.0),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                // Hardware Sensor Pill
                Container(
                  width: double.infinity,
                  padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
                  decoration: BoxDecoration(
                    color: const Color(0xFF1E293B),
                    border: Border.all(color: borderGray),
                    borderRadius: BorderRadius.zero,
                  ),
                  child: Row(
                    children: [
                      Container(
                        width: 8,
                        height: 8,
                        decoration: const BoxDecoration(
                          color: successGreen,
                          shape: BoxShape.circle,
                        ),
                      ),
                      const SizedBox(width: 8),
                      const Expanded(
                        child: Text(
                          'ACCELEROMETER < 0.05g STABLE • HARDWARE RETICLE VERIFIED',
                          style: TextStyle(
                            fontFamily: 'JetBrains Mono',
                            fontSize: 11,
                            fontWeight: FontWeight.w700,
                            color: Colors.white,
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
                const SizedBox(height: 16),

                // Mechanical Hour Meter Entry Card
                Container(
                  width: double.infinity,
                  padding: const EdgeInsets.all(16),
                  decoration: BoxDecoration(
                    color: cardDark,
                    border: Border.all(
                      color: _isMonotonicValid ? primaryYellow : const Color(0xFFEF4444),
                      width: 2,
                    ),
                    borderRadius: BorderRadius.zero,
                  ),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          const Text(
                            'ENGINE HOUR METER AUDIT',
                            style: TextStyle(
                              fontFamily: 'Space Grotesk',
                              fontSize: 14,
                              fontWeight: FontWeight.w700,
                              color: Colors.white,
                            ),
                          ),
                          Container(
                            padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
                            color: const Color(0xFF262626),
                            child: const Text(
                              'IMMUTABLE FLOOR',
                              style: TextStyle(
                                fontFamily: 'JetBrains Mono',
                                fontSize: 10,
                                color: Colors.white70,
                              ),
                            ),
                          ),
                        ],
                      ),
                      const SizedBox(height: 12),
                      Row(
                        children: [
                          Expanded(
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                const Text(
                                  'CERTIFIED CHECK-IN',
                                  style: TextStyle(
                                    fontFamily: 'JetBrains Mono',
                                    fontSize: 11,
                                    color: Colors.white60,
                                  ),
                                ),
                                const SizedBox(height: 4),
                                Text(
                                  '${widget.checkInMeterHours.toStringAsFixed(2)} hrs',
                                  style: const TextStyle(
                                    fontFamily: 'JetBrains Mono',
                                    fontSize: 18,
                                    fontWeight: FontWeight.w700,
                                    color: Colors.white,
                                  ),
                                ),
                              ],
                            ),
                          ),
                          Container(width: 1, height: 40, color: borderGray),
                          const SizedBox(width: 16),
                          Expanded(
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                const Text(
                                  'OPERATING DELTA',
                                  style: TextStyle(
                                    fontFamily: 'JetBrains Mono',
                                    fontSize: 11,
                                    color: Colors.white60,
                                  ),
                                ),
                                const SizedBox(height: 4),
                                Text(
                                  '+${_hoursDelta > 0 ? _hoursDelta.toStringAsFixed(2) : "0.00"} hrs',
                                  style: TextStyle(
                                    fontFamily: 'JetBrains Mono',
                                    fontSize: 18,
                                    fontWeight: FontWeight.w700,
                                    color: _isMonotonicValid ? primaryYellow : const Color(0xFFEF4444),
                                  ),
                                ),
                              ],
                            ),
                          ),
                        ],
                      ),
                      const SizedBox(height: 16),
                      TextField(
                        controller: _meterController,
                        keyboardType: const TextInputType.numberWithOptions(decimal: true),
                        onChanged: (_) => setState(() {}),
                        style: const TextStyle(
                          fontFamily: 'JetBrains Mono',
                          fontSize: 20,
                          fontWeight: FontWeight.w700,
                          color: Colors.white,
                        ),
                        decoration: InputDecoration(
                          labelText: 'RETURN CHECK-OUT METER READING',
                          labelStyle: const TextStyle(color: Colors.white70),
                          suffixText: 'HRS',
                          suffixStyle: const TextStyle(color: primaryYellow, fontFamily: 'JetBrains Mono'),
                          filled: true,
                          fillColor: const Color(0xFF121212),
                          border: const OutlineInputBorder(
                            borderRadius: BorderRadius.zero,
                            borderSide: BorderSide(color: borderGray),
                          ),
                          focusedBorder: OutlineInputBorder(
                            borderRadius: BorderRadius.zero,
                            borderSide: BorderSide(
                              color: _isMonotonicValid ? primaryYellow : const Color(0xFFEF4444),
                              width: 2,
                            ),
                          ),
                        ),
                      ),
                      const SizedBox(height: 8),
                      Row(
                        children: [
                          Icon(
                            _isMonotonicValid ? Icons.check_circle : Icons.warning_amber,
                            size: 16,
                            color: _isMonotonicValid ? successGreen : const Color(0xFFEF4444),
                          ),
                          const SizedBox(width: 6),
                          Expanded(
                            child: Text(
                              _isMonotonicValid
                                  ? 'MONOTONIC CHECK PASSED: Non-regressive meter reading verified.'
                                  : 'ERROR: Return reading cannot regress below check-in (${widget.checkInMeterHours} hrs).',
                              style: TextStyle(
                                fontFamily: 'JetBrains Mono',
                                fontSize: 11,
                                color: _isMonotonicValid ? successGreen : const Color(0xFFEF4444),
                              ),
                            ),
                          ),
                        ],
                      ),
                    ],
                  ),
                ),
                const SizedBox(height: 24),

                // 6-Angle Photographic Reticle Grid
                const Text(
                  'MANDATORY 6-ANGLE RETICLE AUDIT',
                  style: TextStyle(
                    fontFamily: 'Space Grotesk',
                    fontSize: 14,
                    fontWeight: FontWeight.w700,
                    color: Colors.white,
                    letterSpacing: 0.5,
                  ),
                ),
                const SizedBox(height: 12),
                GridView.builder(
                  shrinkWrap: true,
                  physics: const NeverScrollableScrollPhysics(),
                  gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                    crossAxisCount: 2,
                    crossAxisSpacing: 12,
                    mainAxisSpacing: 12,
                    childAspectRatio: 1.15,
                  ),
                  itemCount: _angles.length,
                  itemBuilder: (context, index) {
                    final angle = _angles[index];
                    final key = angle['key']!;
                    final isCaptured = _capturedPhotos.containsKey(key);

                    return Container(
                      padding: const EdgeInsets.all(10),
                      decoration: BoxDecoration(
                        color: cardDark,
                        border: Border.all(
                          color: isCaptured ? successGreen : borderGray,
                          width: 1.5,
                        ),
                        borderRadius: BorderRadius.zero,
                      ),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Row(
                            mainAxisAlignment: MainAxisAlignment.spaceBetween,
                            children: [
                              Icon(
                                isCaptured ? Icons.check_circle : Icons.camera_alt,
                                size: 18,
                                color: isCaptured ? successGreen : primaryYellow,
                              ),
                              Container(
                                padding: const EdgeInsets.symmetric(horizontal: 4, vertical: 2),
                                color: isCaptured ? const Color(0xFF064E3B) : const Color(0xFF262626),
                                child: Text(
                                  isCaptured ? 'CAPTURED' : 'PENDING',
                                  style: TextStyle(
                                    fontFamily: 'JetBrains Mono',
                                    fontSize: 9,
                                    fontWeight: FontWeight.w700,
                                    color: isCaptured ? successGreen : Colors.white70,
                                  ),
                                ),
                              ),
                            ],
                          ),
                          const Spacer(),
                          Text(
                            angle['title']!,
                            style: const TextStyle(
                              fontFamily: 'Space Grotesk',
                              fontSize: 12,
                              fontWeight: FontWeight.w700,
                              color: Colors.white,
                            ),
                          ),
                          const SizedBox(height: 2),
                          Text(
                            angle['desc']!,
                            maxLines: 1,
                            overflow: TextOverflow.ellipsis,
                            style: const TextStyle(
                              fontSize: 10,
                              color: Colors.white60,
                            ),
                          ),
                        ],
                      ),
                    );
                  },
                ),
                const SizedBox(height: 24),

                // Field Verification Toggles
                const Text(
                  'OPERATIONAL SAFETY VERIFICATION',
                  style: TextStyle(
                    fontFamily: 'Space Grotesk',
                    fontSize: 14,
                    fontWeight: FontWeight.w700,
                    color: Colors.white,
                  ),
                ),
                const SizedBox(height: 8),
                CheckboxListTile(
                  value: _fluidLevelsOk,
                  onChanged: (v) => setState(() => _fluidLevelsOk = v ?? false),
                  activeColor: primaryYellow,
                  checkColor: Colors.black,
                  tileColor: cardDark,
                  title: const Text(
                    'Fluid Levels Within Limits',
                    style: TextStyle(color: Colors.white, fontSize: 13),
                  ),
                  subtitle: const Text(
                    'Hydraulic fluid, engine oil, and coolant checked.',
                    style: TextStyle(color: Colors.white60, fontSize: 11),
                  ),
                ),
                const SizedBox(height: 8),
                CheckboxListTile(
                  value: _safetyGlassIntact,
                  onChanged: (v) => setState(() => _safetyGlassIntact = v ?? false),
                  activeColor: primaryYellow,
                  checkColor: Colors.black,
                  tileColor: cardDark,
                  title: const Text(
                    'Cabin Safety Glass Intact',
                    style: TextStyle(color: Colors.white, fontSize: 13),
                  ),
                  subtitle: const Text(
                    'ROPS/FOPS cabin safety glass and mirrors undamaged.',
                    style: TextStyle(color: Colors.white60, fontSize: 11),
                  ),
                ),
                const SizedBox(height: 16),

                // Notes Field
                TextField(
                  controller: _notesController,
                  maxLines: 2,
                  style: const TextStyle(color: Colors.white, fontSize: 13),
                  decoration: const InputDecoration(
                    labelText: 'OFF-HIRE FIELD OBSERVATION NOTES',
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

                // Submit Button (56dp touch target)
                SizedBox(
                  width: double.infinity,
                  height: 56,
                  child: ElevatedButton(
                    onPressed: (isLoading || !_isMonotonicValid) ? null : _submit,
                    style: ElevatedButton.styleFrom(
                      backgroundColor: primaryYellow,
                      foregroundColor: Colors.black,
                      disabledBackgroundColor: borderGray,
                      shape: const RoundedRectangleBorder(borderRadius: BorderRadius.zero),
                    ),
                    child: isLoading
                        ? const CircularProgressIndicator(color: Colors.black)
                        : const Text(
                            'SUBMIT RETURN INSPECTION →',
                            style: TextStyle(
                              fontFamily: 'Space Grotesk',
                              fontSize: 15,
                              fontWeight: FontWeight.w700,
                              letterSpacing: 0.5,
                            ),
                          ),
                  ),
                ),
              ],
            ),
          );
        },
      ),
    );
  }
}

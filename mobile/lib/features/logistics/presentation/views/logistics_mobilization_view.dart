import 'package:dink_mobile/l10n/app_localizations.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../../domain/entities/matched_carrier.dart';
import '../bloc/logistics_bloc.dart';
import '../bloc/logistics_event.dart';
import '../bloc/logistics_state.dart';

class LogisticsMobilizationView extends StatefulWidget {
  final String bookingId;
  final String equipmentTitle;

  const LogisticsMobilizationView({
    super.key,
    required this.bookingId,
    this.equipmentTitle = 'CAT 320D EXCAVATOR',
  });

  @override
  State<LogisticsMobilizationView> createState() =>
      _LogisticsMobilizationViewState();
}

class _LogisticsMobilizationViewState extends State<LogisticsMobilizationView> {
  int _selectedTabIndex = 0; // 0 = Platform Lowboy, 1 = Own Private Truck
  final _formKey = GlobalKey<FormState>();
  final _carrierNameController =
      TextEditingController(text: 'Private In-House Fleet');
  final _driverNameController = TextEditingController();
  final _driverPhoneController = TextEditingController();
  final _plateController = TextEditingController();

  @override
  void initState() {
    super.initState();
    context.read<LogisticsBloc>().add(
          MatchLogisticsCarriersRequested(
            bookingId: widget.bookingId,
            searchRadiusKm: 75,
            pickupDate: DateTime.now().toUtc().add(const Duration(hours: 24)),
          ),
        );
  }

  @override
  void dispose() {
    _carrierNameController.dispose();
    _driverNameController.dispose();
    _driverPhoneController.dispose();
    _plateController.dispose();
    super.dispose();
  }

  void _bookCarrier(MatchedCarrier carrier) {
    context.read<LogisticsBloc>().add(
          CreateLogisticsSubBookingRequested(
            bookingId: widget.bookingId,
            logisticsFleetId: carrier.fleetId,
            pickupDateTime:
                DateTime.now().toUtc().add(const Duration(hours: 24)),
            originLat: 8.9124,
            originLng: 38.7412,
            originName: 'Kality Industrial Depot',
            destLat: 8.5412,
            destLng: 39.2718,
            destName: 'Mojo Highway Site B',
          ),
        );
  }

  void _submitExternalDriver() {
    if (!(_formKey.currentState?.validate() ?? false)) {
      return;
    }
    context.read<LogisticsBloc>().add(
          RegisterExternalTransportRequested(
            bookingId: widget.bookingId,
            carrierName: _carrierNameController.text.trim(),
            driverFullName: _driverNameController.text.trim(),
            driverPhoneNumber: _driverPhoneController.text.trim(),
            primeMoverPlateNumber: _plateController.text.trim(),
            estimatedArrivalDateTime:
                DateTime.now().toUtc().add(const Duration(hours: 6)),
          ),
        );
  }

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;
    const bgDark = Color(0xFF111827);
    const cardDark = Color(0xFF1C1C1C);
    const borderGray = Color(0xFF3E3E3E);
    const primaryYellow = Color(0xFFFFB300);

    return Scaffold(
      backgroundColor: bgDark,
      appBar: AppBar(
        backgroundColor: bgDark,
        elevation: 0,
        title: Text(
          l10n.logisticsMobilizationTitle.toUpperCase(),
          style: const TextStyle(
            color: Color(0xFFF9FAFB),
            fontWeight: FontWeight.w700,
            letterSpacing: 1.0,
          ),
        ),
      ),
      body: BlocConsumer<LogisticsBloc, LogisticsState>(
        listener: (context, state) {
          if (state is LogisticsFailure) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                content: Text(state.message),
                backgroundColor: const Color(0xFFEF4444),
              ),
            );
          }
        },
        builder: (context, state) {
          return SingleChildScrollView(
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                // Booking Context Header Pill
                Container(
                  padding: const EdgeInsets.all(14),
                  decoration: BoxDecoration(
                    color: cardDark,
                    border: Border.all(color: borderGray, width: 2),
                  ),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        '${widget.equipmentTitle} • CONFIRMED',
                        style: const TextStyle(
                          color: primaryYellow,
                          fontWeight: FontWeight.w700,
                          fontSize: 13,
                          letterSpacing: 0.8,
                        ),
                      ),
                      if (state is LogisticsCarriersMatched) ...[
                        const SizedBox(height: 10),
                        Wrap(
                          spacing: 12,
                          runSpacing: 8,
                          children: [
                            _SpecBadge(
                              label: 'OP WEIGHT',
                              value:
                                  '${state.result.operatingWeightKg.toStringAsFixed(0)} kg',
                            ),
                            _SpecBadge(
                              label: 'MIN PAYLOAD (1.10x)',
                              value:
                                  '${state.result.requiredMinimumPayloadKg.toStringAsFixed(0)} kg',
                              highlight: true,
                            ),
                            _SpecBadge(
                              label: 'DIMENSIONS',
                              value:
                                  '${state.result.transportLengthMm}x${state.result.transportWidthMm} mm',
                            ),
                          ],
                        ),
                        const SizedBox(height: 8),
                        Text(
                          state.result.recommendedTrailerClass,
                          style: const TextStyle(
                            color: Color(0xFF10B981),
                            fontWeight: FontWeight.w600,
                            fontSize: 12,
                          ),
                        ),
                      ],
                    ],
                  ),
                ),
                const SizedBox(height: 16),

                // Two-Tab Mobilization Selector (>= 56dp height)
                Row(
                  children: [
                    Expanded(
                      child: SizedBox(
                        height: 56,
                        child: ElevatedButton(
                          onPressed: () =>
                              setState(() => _selectedTabIndex = 0),
                          style: ElevatedButton.styleFrom(
                            backgroundColor: _selectedTabIndex == 0
                                ? primaryYellow
                                : cardDark,
                            foregroundColor: _selectedTabIndex == 0
                                ? Colors.black
                                : const Color(0xFFF9FAFB),
                            shape: const RoundedRectangleBorder(
                              borderRadius: BorderRadius.zero,
                              side: BorderSide(color: borderGray, width: 2),
                            ),
                          ),
                          child: Text(
                            l10n.platformLowboyTab.toUpperCase(),
                            textAlign: TextAlign.center,
                            style: const TextStyle(
                              fontWeight: FontWeight.w700,
                              fontSize: 12,
                            ),
                          ),
                        ),
                      ),
                    ),
                    const SizedBox(width: 8),
                    Expanded(
                      child: SizedBox(
                        height: 56,
                        child: ElevatedButton(
                          onPressed: () =>
                              setState(() => _selectedTabIndex = 1),
                          style: ElevatedButton.styleFrom(
                            backgroundColor: _selectedTabIndex == 1
                                ? primaryYellow
                                : cardDark,
                            foregroundColor: _selectedTabIndex == 1
                                ? Colors.black
                                : const Color(0xFFF9FAFB),
                            shape: const RoundedRectangleBorder(
                              borderRadius: BorderRadius.zero,
                              side: BorderSide(color: borderGray, width: 2),
                            ),
                          ),
                          child: Text(
                            l10n.externalTruckTab.toUpperCase(),
                            textAlign: TextAlign.center,
                            style: const TextStyle(
                              fontWeight: FontWeight.w700,
                              fontSize: 12,
                            ),
                          ),
                        ),
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 16),

                if (state is LogisticsLoading)
                  const Center(
                    child: Padding(
                      padding: EdgeInsets.all(24),
                      child: CircularProgressIndicator(color: primaryYellow),
                    ),
                  )
                else if (_selectedTabIndex == 0 &&
                    state is LogisticsCarriersMatched)
                  ...state.result.availableLogisticsFleets.map(
                    (carrier) => _CarrierMatchCard(
                      carrier: carrier,
                      bookLabel: l10n.bookLowboyCarrier,
                      onBook: () => _bookCarrier(carrier),
                    ),
                  )
                else if (_selectedTabIndex == 1)
                  _buildExternalTruckForm(l10n, cardDark, borderGray, primaryYellow),

                if (state is LogisticsSubBookingCreated)
                  Container(
                    margin: const EdgeInsets.only(top: 16),
                    padding: const EdgeInsets.all(14),
                    decoration: BoxDecoration(
                      color: const Color(0xFF064E3B),
                      border: Border.all(
                          color: const Color(0xFF10B981), width: 2),
                    ),
                    child: Text(
                      'SUB-BOOKING ${state.subBooking.status} • 2.5% FEE: ${state.subBooking.commissionOwedMinor} ${state.subBooking.currency}',
                      style: const TextStyle(
                        color: Color(0xFFF9FAFB),
                        fontWeight: FontWeight.w700,
                      ),
                    ),
                  ),

                if (state is ExternalTransportRegistered)
                  Container(
                    margin: const EdgeInsets.only(top: 16),
                    padding: const EdgeInsets.all(14),
                    decoration: BoxDecoration(
                      color: const Color(0xFF064E3B),
                      border: Border.all(
                          color: const Color(0xFF10B981), width: 2),
                    ),
                    child: Text(
                      'EXTERNAL TRANSPORT REGISTERED: ${state.externalTransport.driverFullName} (${state.externalTransport.primeMoverPlateNumber})',
                      style: const TextStyle(
                        color: Color(0xFFF9FAFB),
                        fontWeight: FontWeight.w700,
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

  Widget _buildExternalTruckForm(
    AppLocalizations l10n,
    Color cardDark,
    Color borderGray,
    Color primaryYellow,
  ) {
    return Form(
      key: _formKey,
      child: Container(
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(
          color: cardDark,
          border: Border.all(color: borderGray, width: 2),
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            _IndustrialField(
              controller: _carrierNameController,
              label: l10n.carrierCompanyLabel,
            ),
            const SizedBox(height: 12),
            _IndustrialField(
              controller: _driverNameController,
              label: l10n.driverFullNameLabel,
            ),
            const SizedBox(height: 12),
            _IndustrialField(
              controller: _driverPhoneController,
              label: l10n.driverPhoneLabel,
            ),
            const SizedBox(height: 12),
            _IndustrialField(
              controller: _plateController,
              label: l10n.primeMoverPlateLabel,
            ),
            const SizedBox(height: 16),
            SizedBox(
              height: 56,
              child: ElevatedButton(
                onPressed: _submitExternalDriver,
                style: ElevatedButton.styleFrom(
                  backgroundColor: primaryYellow,
                  foregroundColor: Colors.black,
                  shape: const RoundedRectangleBorder(
                      borderRadius: BorderRadius.zero),
                ),
                child: Text(
                  l10n.declareExternalDriver.toUpperCase(),
                  style: const TextStyle(fontWeight: FontWeight.w700),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _SpecBadge extends StatelessWidget {
  final String label;
  final String value;
  final bool highlight;

  const _SpecBadge({
    required this.label,
    required this.value,
    this.highlight = false,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 6),
      decoration: BoxDecoration(
        color: const Color(0xFF111827),
        border: Border.all(
          color: highlight ? const Color(0xFFFFB300) : const Color(0xFF3E3E3E),
        ),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            label,
            style: const TextStyle(color: Color(0xFF9CA3AF), fontSize: 10),
          ),
          Text(
            value,
            style: TextStyle(
              color: highlight
                  ? const Color(0xFFFFB300)
                  : const Color(0xFFF9FAFB),
              fontWeight: FontWeight.w700,
              fontSize: 13,
            ),
          ),
        ],
      ),
    );
  }
}

class _CarrierMatchCard extends StatelessWidget {
  final MatchedCarrier carrier;
  final String bookLabel;
  final VoidCallback onBook;

  const _CarrierMatchCard({
    required this.carrier,
    required this.bookLabel,
    required this.onBook,
  });

  @override
  Widget build(BuildContext context) {
    final kmAway = (carrier.distanceFromDepotMeters / 1000.0).toStringAsFixed(1);
    return Container(
      margin: const EdgeInsets.only(bottom: 12),
      padding: const EdgeInsets.all(14),
      decoration: BoxDecoration(
        color: const Color(0xFF1C1C1C),
        border: Border.all(color: const Color(0xFF3E3E3E), width: 2),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          Text(
            '${carrier.make} ${carrier.model}',
            style: const TextStyle(
              color: Color(0xFFF9FAFB),
              fontWeight: FontWeight.w700,
              fontSize: 16,
            ),
          ),
          const SizedBox(height: 4),
          Text(
            '${carrier.ownerTradingName} • ${carrier.averageRating.toStringAsFixed(2)}★ • $kmAway km',
            style: const TextStyle(color: Color(0xFFD1D5DB), fontSize: 13),
          ),
          const SizedBox(height: 10),
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              _SpecBadge(
                label: 'PAYLOAD',
                value: '${carrier.payloadCapacityKg.toStringAsFixed(0)} kg',
              ),
              _SpecBadge(
                label: 'DECK',
                value: '${carrier.deckWellLengthMm} mm',
              ),
              _SpecBadge(
                label: 'RATE/DAY',
                value: '${carrier.dailyRateMinor ~/ 100} ${carrier.currency}',
              ),
            ],
          ),
          const SizedBox(height: 12),
          SizedBox(
            height: 56,
            child: ElevatedButton(
              onPressed: onBook,
              style: ElevatedButton.styleFrom(
                backgroundColor: const Color(0xFFFFB300),
                foregroundColor: Colors.black,
                shape: const RoundedRectangleBorder(
                  borderRadius: BorderRadius.zero,
                ),
              ),
              child: Text(
                bookLabel.toUpperCase(),
                style: const TextStyle(
                  fontWeight: FontWeight.w700,
                  letterSpacing: 0.8,
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }
}

class _IndustrialField extends StatelessWidget {
  final TextEditingController controller;
  final String label;

  const _IndustrialField({required this.controller, required this.label});

  @override
  Widget build(BuildContext context) {
    return ConstrainedBox(
      constraints: const BoxConstraints(minHeight: 56),
      child: TextFormField(
        controller: controller,
        style: const TextStyle(color: Color(0xFFF9FAFB)),
        validator: (v) =>
            (v == null || v.trim().isEmpty) ? '$label is required' : null,
        decoration: InputDecoration(
          labelText: label,
          labelStyle: const TextStyle(color: Color(0xFFD1D5DB)),
          filled: true,
          fillColor: const Color(0xFF111827),
          enabledBorder: const OutlineInputBorder(
            borderRadius: BorderRadius.zero,
            borderSide: BorderSide(color: Color(0xFF3E3E3E), width: 2),
          ),
          focusedBorder: const OutlineInputBorder(
            borderRadius: BorderRadius.zero,
            borderSide: BorderSide(color: Color(0xFFFFB300), width: 2),
          ),
          errorBorder: const OutlineInputBorder(
            borderRadius: BorderRadius.zero,
            borderSide: BorderSide(color: Color(0xFFEF4444), width: 2),
          ),
        ),
      ),
    );
  }
}

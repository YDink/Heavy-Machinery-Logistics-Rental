import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:dink_mobile/features/auth/presentation/bloc/auth_bloc.dart';
import 'package:dink_mobile/features/auth/presentation/bloc/auth_event.dart';
import 'package:dink_mobile/features/auth/presentation/bloc/auth_state.dart';
import 'package:dink_mobile/features/auth/presentation/pages/otp_verification_page.dart';
import 'package:dink_mobile/l10n/app_localizations.dart';

class PhoneInputPage extends StatefulWidget {
  const PhoneInputPage({super.key});

  @override
  State<PhoneInputPage> createState() => _PhoneInputPageState();
}

class _PhoneInputPageState extends State<PhoneInputPage> {
  final _phoneController = TextEditingController();
  String _selectedCountryCode = 'ET';
  String _dialPrefix = '+251';
  String _detectedCarrier = '';

  @override
  void initState() {
    super.initState();
    _phoneController.addListener(_onPhoneChanged);
  }

  @override
  void dispose() {
    _phoneController.removeListener(_onPhoneChanged);
    _phoneController.dispose();
    super.dispose();
  }

  void _onPhoneChanged() {
    final text = _phoneController.text.trim();
    String carrier = '';

    if (_selectedCountryCode == 'ET') {
      if (text.startsWith('9')) {
        carrier = 'Ethio Telecom';
      } else if (text.startsWith('7')) {
        carrier = 'Safaricom ET';
      }
    } else if (_selectedCountryCode == 'KE') {
      if (text.startsWith('7')) {
        carrier = 'Safaricom KE';
      } else if (text.startsWith('1')) {
        carrier = 'Airtel KE';
      }
    }

    if (carrier != _detectedCarrier) {
      setState(() {
        _detectedCarrier = carrier;
      });
    }
  }

  void _onCountryChanged(String? newCode) {
    if (newCode == null) return;
    setState(() {
      _selectedCountryCode = newCode;
      _dialPrefix = newCode == 'ET' ? '+251' : '+254';
      _phoneController.clear();
      _detectedCarrier = '';
    });
  }

  void _submit() {
    final l10n = AppLocalizations.of(context)!;
    final rawNumber = _phoneController.text.trim();
    if (rawNumber.isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text(l10n.errorInvalidPhone)),
      );
      return;
    }

    // Format full E.164 phone string
    final fullNumber = '$_dialPrefix$rawNumber';

    context.read<AuthBloc>().add(
          RequestOtpEvent(
            phoneNumber: fullNumber,
            countryCode: _selectedCountryCode,
          ),
        );
  }

  @override
  Widget build(BuildContext context) {
    const brandGold = Color(0xFFE5A93C);
    final l10n = AppLocalizations.of(context)!;

    return Scaffold(
      appBar: AppBar(
        title: Text(l10n.appName),
        backgroundColor: Colors.transparent,
        elevation: 0,
      ),
      body: BlocConsumer<AuthBloc, AuthState>(
        listener: (context, state) {
          if (state is OtpSentState) {
            Navigator.of(context).push(
              MaterialPageRoute(
                builder: (_) => BlocProvider.value(
                  value: context.read<AuthBloc>(),
                  child: OtpVerificationPage(
                    phoneNumber: state.phoneNumber,
                  ),
                ),
              ),
            );
          } else if (state is AuthErrorState) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                content: Text(state.message),
                backgroundColor: Colors.red.shade700,
              ),
            );
          }
        },
        builder: (context, state) {
          final isLoading = state is AuthLoadingState;

          return SafeArea(
            child: Padding(
              padding: const EdgeInsets.symmetric(horizontal: 24.0, vertical: 16.0),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.stretch,
                children: [
                  const SizedBox(height: 20),
                  const Icon(
                    Icons.construction,
                    size: 56,
                    color: brandGold,
                  ),
                  const SizedBox(height: 16),
                  Text(
                    l10n.phoneInputTitle,
                    textAlign: TextAlign.center,
                    style: const TextStyle(
                      fontSize: 26,
                      fontWeight: FontWeight.bold,
                      letterSpacing: -0.5,
                    ),
                  ),
                  const SizedBox(height: 8),
                  Text(
                    l10n.phoneInputSubtitle,
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      fontSize: 14,
                      color: Colors.grey.shade600,
                    ),
                  ),
                  const SizedBox(height: 36),

                  // Phone Input Container
                  Row(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      // Country Dropdown Picker
                      Container(
                        padding: const EdgeInsets.symmetric(horizontal: 12),
                        decoration: BoxDecoration(
                          border: Border.all(color: Colors.grey.shade400),
                          borderRadius: BorderRadius.circular(8),
                        ),
                        child: DropdownButtonHideUnderline(
                          child: DropdownButton<String>(
                            value: _selectedCountryCode,
                            onChanged: isLoading ? null : _onCountryChanged,
                            items: const [
                              DropdownMenuItem(
                                value: 'ET',
                                child: Text('🇪🇹 +251'),
                              ),
                              DropdownMenuItem(
                                value: 'KE',
                                child: Text('🇰🇪 +254'),
                              ),
                            ],
                          ),
                        ),
                      ),
                      const SizedBox(width: 12),

                      // Number Text Field
                      Expanded(
                        child: TextField(
                          controller: _phoneController,
                          enabled: !isLoading,
                          keyboardType: TextInputType.phone,
                          inputFormatters: [
                            FilteringTextInputFormatter.digitsOnly,
                            LengthLimitingTextInputFormatter(9),
                          ],
                          decoration: InputDecoration(
                            labelText: l10n.phoneNumberLabel,
                            hintText: _selectedCountryCode == 'ET' ? '911 223 344' : '712 345 678',
                            border: OutlineInputBorder(
                              borderRadius: BorderRadius.circular(8),
                            ),
                            contentPadding: const EdgeInsets.symmetric(horizontal: 14, vertical: 14),
                          ),
                        ),
                      ),
                    ],
                  ),

                  // Dynamic Carrier Ingress Badge (REG-02)
                  if (_detectedCarrier.isNotEmpty) ...[
                    const SizedBox(height: 10),
                    Align(
                      alignment: Alignment.centerLeft,
                      child: Container(
                        padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
                        decoration: BoxDecoration(
                          color: Colors.green.shade50,
                          borderRadius: BorderRadius.circular(12),
                          border: Border.all(color: Colors.green.shade300),
                        ),
                        child: Row(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            Icon(Icons.network_check, size: 14, color: Colors.green.shade700),
                            const SizedBox(width: 6),
                            Text(
                              '${l10n.carrierPrefix}: $_detectedCarrier',
                              style: TextStyle(
                                fontSize: 12,
                                fontWeight: FontWeight.w600,
                                color: Colors.green.shade800,
                              ),
                            ),
                          ],
                        ),
                      ),
                    ),
                  ],

                  const Spacer(),

                  Text(
                    l10n.termsAgreementNotice,
                    textAlign: TextAlign.center,
                    style: TextStyle(fontSize: 12, color: Colors.grey.shade600),
                  ),
                  const SizedBox(height: 16),

                  // Continue Button
                  SizedBox(
                    height: 50,
                    child: ElevatedButton(
                      onPressed: isLoading ? null : _submit,
                      style: ElevatedButton.styleFrom(
                        backgroundColor: brandGold,
                        foregroundColor: Colors.white,
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(8),
                        ),
                      ),
                      child: isLoading
                          ? const SizedBox(
                              height: 20,
                              width: 20,
                              child: CircularProgressIndicator(
                                strokeWidth: 2,
                                color: Colors.white,
                              ),
                            )
                          : Text(
                              l10n.sendOtpButton,
                              style: const TextStyle(fontSize: 16, fontWeight: FontWeight.bold),
                            ),
                    ),
                  ),
                  const SizedBox(height: 16),
                ],
              ),
            ),
          );
        },
      ),
    );
  }
}

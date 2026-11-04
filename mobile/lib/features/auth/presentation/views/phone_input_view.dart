import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../bloc/auth_bloc.dart';
import '../bloc/auth_event.dart';
import '../bloc/auth_state.dart';
import 'otp_verification_view.dart';
import 'package:dink_mobile/l10n/app_localizations.dart';

typedef PhoneInputPage = PhoneInputView;

class PhoneInputView extends StatefulWidget {
  const PhoneInputView({super.key});

  @override
  State<PhoneInputView> createState() => _PhoneInputViewState();
}

class _PhoneInputViewState extends State<PhoneInputView> {
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
    if (newCode == null || newCode == _selectedCountryCode) return;
    setState(() {
      _selectedCountryCode = newCode;
      _dialPrefix = newCode == 'ET' ? '+251' : '+254';
      _phoneController.clear();
      _detectedCarrier = '';
    });
  }

  void _submit() {
    final rawNumber = _phoneController.text.trim();
    if (rawNumber.isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text(AppLocalizations.of(context)!.errorInvalidPhone)),
      );
      return;
    }

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
                  child: OtpVerificationView(phoneNumber: state.phoneNumber),
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
                  Text(
                    l10n.phoneInputTitle,
                    style: const TextStyle(
                      fontSize: 24,
                      fontWeight: FontWeight.bold,
                      letterSpacing: -0.5,
                    ),
                  ),
                  const SizedBox(height: 8),
                  Text(
                    l10n.phoneInputSubtitle,
                    style: TextStyle(
                      fontSize: 14,
                      color: Colors.grey.shade600,
                    ),
                  ),
                  const SizedBox(height: 32),
                  Row(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Container(
                        height: 56,
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
                      Expanded(
                        child: TextField(
                          controller: _phoneController,
                          enabled: !isLoading,
                          keyboardType: TextInputType.phone,
                          inputFormatters: [
                            FilteringTextInputFormatter.digitsOnly,
                            LengthLimitingTextInputFormatter(_selectedCountryCode == 'ET' ? 9 : 10),
                          ],
                          decoration: InputDecoration(
                            labelText: l10n.phoneNumberLabel,
                            hintText: _selectedCountryCode == 'ET' ? '911 223 344' : '712 345 678',
                            border: OutlineInputBorder(
                              borderRadius: BorderRadius.circular(8),
                            ),
                            contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 18),
                          ),
                          onSubmitted: (_) => _submit(),
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 12),
                  if (_detectedCarrier.isNotEmpty)
                    Align(
                      alignment: Alignment.centerLeft,
                      child: Container(
                        padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
                        decoration: BoxDecoration(
                          color: Colors.amber.shade50,
                          borderRadius: BorderRadius.circular(12),
                          border: Border.all(color: Colors.amber.shade300),
                        ),
                        child: Row(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            Icon(Icons.network_cell, size: 14, color: Colors.amber.shade800),
                            const SizedBox(width: 4),
                            Text(
                              '${l10n.carrierPrefix}: $_detectedCarrier',
                              style: TextStyle(
                                fontSize: 12,
                                fontWeight: FontWeight.bold,
                                color: Colors.amber.shade900,
                              ),
                            ),
                          ],
                        ),
                      ),
                    ),
                  const Spacer(),
                  Text(
                    l10n.termsAgreementNotice,
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      fontSize: 12,
                      color: Colors.grey.shade600,
                    ),
                  ),
                  const SizedBox(height: 16),
                  SizedBox(
                    height: 56,
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

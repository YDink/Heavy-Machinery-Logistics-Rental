package domain

import (
	"fmt"
	"regexp"
	"strings"
)

type Carrier string

const (
	CarrierEthioTelecom Carrier = "ETHIO_TELECOM"
	CarrierSafaricomET  Carrier = "SAFARICOM_ET"
	CarrierSafaricomKE  Carrier = "SAFARICOM_KE"
	CarrierAirtelKE     Carrier = "AIRTEL_KE"
	CarrierUnknown      Carrier = "UNKNOWN"
)

var (
	// E.164 standard format regex
	e164Regex = regexp.MustCompile(`^\+[1-9]\d{1,14}$`)

	// Ethiopian mobile pattern: +251 followed by 9 digits starting with 9 or 7
	ethiopiaRegex = regexp.MustCompile(`^\+251[79]\d{8}$`)

	// Kenyan mobile pattern: +254 followed by 9 digits starting with 7 or 1
	kenyaRegex = regexp.MustCompile(`^\+254[17]\d{8}$`)
)

// PhoneNumber represents a validated sovereign E.164 mobile number (REG-01).
type PhoneNumber struct {
	raw         string
	countryCode string
	carrier     Carrier
}

// ParsePhoneNumber validates and normalizes an input phone number string into an E.164 PhoneNumber value object.
func ParsePhoneNumber(input string, defaultCountry string) (PhoneNumber, error) {
	cleaned := strings.TrimSpace(input)
	cleaned = strings.ReplaceAll(cleaned, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	cleaned = strings.ReplaceAll(cleaned, "(", "")
	cleaned = strings.ReplaceAll(cleaned, ")", "")

	// Handle domestic zero-prefixed numbers using default country
	if strings.HasPrefix(cleaned, "0") {
		switch strings.ToUpper(defaultCountry) {
		case "ET":
			cleaned = "+251" + cleaned[1:]
		case "KE":
			cleaned = "+254" + cleaned[1:]
		default:
			return PhoneNumber{}, fmt.Errorf("%w: cannot infer country code for domestic format", ErrInvalidPhoneNumber)
		}
	} else if !strings.HasPrefix(cleaned, "+") {
		// Handle numbers without '+' prefix
		if strings.HasPrefix(cleaned, "251") {
			cleaned = "+" + cleaned
		} else if strings.HasPrefix(cleaned, "254") {
			cleaned = "+" + cleaned
		} else {
			switch strings.ToUpper(defaultCountry) {
			case "ET":
				cleaned = "+251" + cleaned
			case "KE":
				cleaned = "+254" + cleaned
			default:
				return PhoneNumber{}, fmt.Errorf("%w: missing country dialing prefix", ErrInvalidPhoneNumber)
			}
		}
	}

	if !e164Regex.MatchString(cleaned) {
		return PhoneNumber{}, ErrInvalidPhoneNumber
	}

	var country string
	var carrier Carrier

	if ethiopiaRegex.MatchString(cleaned) {
		country = "ET"
		// Determine Ethiopian Carrier (REG-02)
		if strings.HasPrefix(cleaned, "+2519") {
			carrier = CarrierEthioTelecom
		} else if strings.HasPrefix(cleaned, "+2517") {
			carrier = CarrierSafaricomET
		}
	} else if kenyaRegex.MatchString(cleaned) {
		country = "KE"
		// Determine Kenyan Carrier (REG-02)
		if strings.HasPrefix(cleaned, "+2547") {
			carrier = CarrierSafaricomKE
		} else if strings.HasPrefix(cleaned, "+2541") {
			carrier = CarrierAirtelKE
		}
	} else {
		return PhoneNumber{}, ErrUnsupportedCountryCode
	}

	return PhoneNumber{
		raw:         cleaned,
		countryCode: country,
		carrier:     carrier,
	}, nil
}

// String returns the normalized E.164 string representation.
func (p PhoneNumber) String() string {
	return p.raw
}

// CountryCode returns the ISO 3166-1 alpha-2 country code (ET or KE).
func (p PhoneNumber) CountryCode() string {
	return p.countryCode
}

// Carrier returns the detected telco carrier.
func (p PhoneNumber) Carrier() Carrier {
	return p.carrier
}

// Masked returns a masked representation for zero-PII logging (e.g. "+25191****344").
func (p PhoneNumber) Masked() string {
	if len(p.raw) < 8 {
		return "***"
	}
	prefix := p.raw[:6]
	suffix := p.raw[len(p.raw)-3:]
	return prefix + "****" + suffix
}

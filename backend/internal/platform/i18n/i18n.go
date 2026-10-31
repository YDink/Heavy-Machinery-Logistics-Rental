package i18n

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const (
	LanguageContextKey contextKey = "dink_language"

	LangEN = "en"
	LangAM = "am"
	LangOM = "om"
	LangSW = "sw"
)

// SupportedLanguages defines the canonical supported language codes.
var SupportedLanguages = map[string]bool{
	LangEN: true,
	LangAM: true,
	LangOM: true,
	LangSW: true,
}

// Built-in fallback dictionaries for core platform messages and RFC 9457 error payloads.
var fallbackDictionaries = map[string]map[string]string{
	LangEN: {
		"errors.malformed_json_title":        "Bad Request",
		"errors.malformed_json_detail":       "Request body could not be parsed as valid JSON.",
		"errors.invalid_phone_title":         "Invalid Phone Number",
		"errors.invalid_phone_detail":        "The provided phone number is invalid for the selected country.",
		"errors.rate_limit_title":            "Too Many Requests",
		"errors.rate_limit_detail":           "Rate limit exceeded. Please wait before retrying.",
		"errors.forbidden_permission_title":  "Forbidden",
		"errors.forbidden_permission_detail": "You do not have the required permissions to perform this action.",
		"errors.unauthorized_title":          "Unauthorized",
		"errors.unauthorized_detail":         "Authentication credentials are missing, invalid, or expired.",
		"errors.user_not_found_title":        "User Not Found",
		"errors.user_not_found_detail":       "The requested user profile does not exist.",
		"errors.internal_title":              "Internal Server Error",
		"errors.internal_detail":             "An unexpected error occurred. Please contact support.",
	},
	LangAM: {
		"errors.malformed_json_title":        "ልክ ያልሆነ ጥያቄ",
		"errors.malformed_json_detail":       "የተላከው የጥያቄ መረጃ ትክክለኛ የJSON ቅርጸት የለውም።",
		"errors.invalid_phone_title":         "ልክ ያልሆነ ስልክ ቁጥር",
		"errors.invalid_phone_detail":        "የተሰጠው ስልክ ቁጥር ለተመረጠው ሀገር ትክክል አይደለም።",
		"errors.rate_limit_title":            "የሙከራ ገደብ አልፏል",
		"errors.rate_limit_detail":           "የተፈቀደው የጥያቄ ገደብ አልፏል። እባክዎን ትንሽ ቆይተው እንደገና ይሞክሩ።",
		"errors.forbidden_permission_title":  "የተከለከለ እርምጃ",
		"errors.forbidden_permission_detail": "ይህን ተግባር ለማከናወን የሚያስፈልገው ፈቃድ የለዎትም።",
		"errors.unauthorized_title":          "ያልተፈቀደ መዳረሻ",
		"errors.unauthorized_detail":         "የመግቢያ ፈቃድዎ አልተገኘም ወይም ጊዜው አልፏል።",
		"errors.user_not_found_title":        "ተጠቃሚው አልተገኘም",
		"errors.user_not_found_detail":       "የተጠየቀው የተጠቃሚ መገለጫ አልተገኘም።",
		"errors.internal_title":              "የውስጥ ሲስተም ስህተት",
		"errors.internal_detail":             "ያልተጠበቀ ስህተት አጋጥሟል። እባክዎን የድጋፍ ቡድናችንን ያነጋግሩ።",
	},
}

// ParseAcceptLanguage resolves the best matching language code from the Accept-Language header.
func ParseAcceptLanguage(header string) string {
	if header == "" {
		return LangEN
	}

	parts := strings.Split(header, ",")
	for _, part := range parts {
		tag := strings.TrimSpace(strings.Split(part, ";")[0])
		tag = strings.ToLower(tag)

		// Check full match (e.g. "am-et" -> "am")
		if len(tag) >= 2 {
			prefix := tag[:2]
			if SupportedLanguages[prefix] {
				return prefix
			}
		}
	}

	return LangEN
}

// LanguageMiddleware extracts the requested language from the Accept-Language header or user claims.
func LanguageMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lang := ParseAcceptLanguage(r.Header.Get("Accept-Language"))
		ctx := context.WithValue(r.Context(), LanguageContextKey, lang)
		w.Header().Set("Content-Language", lang)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetLanguage returns the language code stored in the context, defaulting to English.
func GetLanguage(ctx context.Context) string {
	if ctx == nil {
		return LangEN
	}
	if lang, ok := ctx.Value(LanguageContextKey).(string); ok && lang != "" {
		return lang
	}
	return LangEN
}

// Translate retrieves the localized string for a key, falling back to English and then to the provided fallback.
func Translate(ctx context.Context, key string, fallback string) string {
	lang := GetLanguage(ctx)

	if dict, exists := fallbackDictionaries[lang]; exists {
		if val, found := dict[key]; found {
			return val
		}
	}

	// Fallback to English
	if lang != LangEN {
		if enDict, exists := fallbackDictionaries[LangEN]; exists {
			if val, found := enDict[key]; found {
				return val
			}
		}
	}

	return fallback
}

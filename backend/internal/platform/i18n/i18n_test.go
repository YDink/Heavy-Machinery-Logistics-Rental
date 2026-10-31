package i18n

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParseAcceptLanguage(t *testing.T) {
	tests := []struct {
		header   string
		expected string
	}{
		{"", LangEN},
		{"am-ET,am;q=0.9,en;q=0.8", LangAM},
		{"am", LangAM},
		{"om-ET,en;q=0.5", LangOM},
		{"sw-KE,en;q=0.8", LangSW},
		{"fr-FR,de;q=0.5", LangEN}, // Unsupported falls back to en
		{"en-US,en;q=0.9", LangEN},
	}

	for _, tt := range tests {
		result := ParseAcceptLanguage(tt.header)
		if result != tt.expected {
			t.Errorf("ParseAcceptLanguage(%q) = %q; want %q", tt.header, result, tt.expected)
		}
	}
}

func TestTranslate_Fallback(t *testing.T) {
	ctxEN := context.WithValue(context.Background(), LanguageContextKey, LangEN)
	ctxAM := context.WithValue(context.Background(), LanguageContextKey, LangAM)

	// English test
	titleEN := Translate(ctxEN, "errors.invalid_phone_title", "Default Title")
	if titleEN != "Invalid Phone Number" {
		t.Errorf("expected 'Invalid Phone Number', got %q", titleEN)
	}

	// Amharic test
	titleAM := Translate(ctxAM, "errors.invalid_phone_title", "Default Title")
	if titleAM != "ልክ ያልሆነ ስልክ ቁጥር" {
		t.Errorf("expected Amharic title, got %q", titleAM)
	}

	// Non-existent key falls back to default
	unknown := Translate(ctxAM, "errors.non_existent", "Fallback String")
	if unknown != "Fallback String" {
		t.Errorf("expected 'Fallback String', got %q", unknown)
	}
}

func TestLanguageMiddleware(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Accept-Language", "am-ET,am;q=0.9")

	rr := httptest.NewRecorder()

	var capturedLang string
	handler := LanguageMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedLang = GetLanguage(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)

	if capturedLang != LangAM {
		t.Errorf("expected language %q, got %q", LangAM, capturedLang)
	}
	if rr.Header().Get("Content-Language") != LangAM {
		t.Errorf("expected Content-Language header %q, got %q", LangAM, rr.Header().Get("Content-Language"))
	}
}

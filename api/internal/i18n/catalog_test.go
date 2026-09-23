package i18n_test

import (
	"testing"

	"github.com/AyitiDev/Bornage/api/internal/i18n"
)

func TestParseLang(t *testing.T) {
	tests := []struct {
		input    string
		expected i18n.Lang
	}{
		{"fr", i18n.LangFrench},
		{"ht", i18n.LangHaitianCreole},
		{"en", i18n.DefaultLang}, 
		{"", i18n.DefaultLang},   
		{"es", i18n.DefaultLang}, 
	}

	for _, tt := range tests {
		t.Run("parse_"+tt.input, func(t *testing.T) {
			got := i18n.ParseLang(tt.input)
			if got != tt.expected {
				t.Errorf("ParseLang(%q): expected %q, got %q", tt.input, tt.expected, got)
			}
		})
	}
}

func TestTranslate_French(t *testing.T) {
	msg := i18n.T(i18n.LangFrench, i18n.KeyNotFound)
	if msg == "" || msg == string(i18n.KeyNotFound) {
		t.Errorf("expected French translation for KeyNotFound, got %q", msg)
	}

	// Must not contain English
	if msg == "Not found" {
		t.Errorf("expected French message, got English")
	}
}

func TestTranslate_HaitianCreole(t *testing.T) {
	msg := i18n.T(i18n.LangHaitianCreole, i18n.KeyNotFound)
	if msg == "" || msg == string(i18n.KeyNotFound) {
		t.Errorf("expected Haitian Creole translation for KeyNotFound, got %q", msg)
	}
}

func TestTranslate_FallbackToFrench(t *testing.T) {
	// An unsupported lang should effectively get French via ParseLang
	lang := i18n.ParseLang("zh")
	msg := i18n.T(lang, i18n.KeyInternalError)
	frMsg := i18n.T(i18n.LangFrench, i18n.KeyInternalError)
	if msg != frMsg {
		t.Errorf("expected fallback to French, got %q", msg)
	}
}

func TestTranslate_AllKeysHaveTranslations(t *testing.T) {
	keys := []i18n.Key{
		i18n.KeyInternalError, i18n.KeyNotFound, i18n.KeyBadRequest,
		i18n.KeyUnauthorized, i18n.KeyForbidden, i18n.KeyValidationFailed,
		i18n.KeyPartyNotFound, i18n.KeyPartyCreated, i18n.KeyPartyNameRequired, i18n.KeyPartyTypeInvalid,
		i18n.KeySpatialUnitNotFound, i18n.KeySpatialUnitCreated, i18n.KeyGeomRequired, i18n.KeyGeomInvalid, i18n.KeyOverlapDetected,
		i18n.KeyClaimNotFound, i18n.KeyClaimCreated, i18n.KeyClaimStatusInvalid, i18n.KeyClaimTransitionDenied, i18n.KeyObjectionPeriodActive,
		i18n.KeySourceNotFound, i18n.KeySourceUploaded, i18n.KeySourceHashMismatch,
		i18n.KeyWitnessNotFound, i18n.KeyWitnessCreated, i18n.KeyWitnessNameRequired,
	}

	langs := []i18n.Lang{i18n.LangFrench, i18n.LangHaitianCreole}

	for _, lang := range langs {
		for _, key := range keys {
			msg := i18n.T(lang, key)
			if msg == string(key) {
				t.Errorf("lang=%q key=%q: missing translation (returned raw key)", lang, key)
			}
			if msg == "" {
				t.Errorf("lang=%q key=%q: translation is empty", lang, key)
			}
		}
	}
}

func TestOpenClosedPrinciple_RegisterNewLanguage(t *testing.T) {
	langES := i18n.Lang("es")

	if i18n.ParseLang("es") != i18n.DefaultLang {
		t.Fatalf("expected 'es' to fall back before registration")
	}

	i18n.RegisterLanguage(langES, map[i18n.Key]string{
		i18n.KeyNotFound:      "El recurso solicitado no fue encontrado",
		i18n.KeyInternalError: "Ocurrió un error interno. Intente nuevamente más tarde",
	})

	if got := i18n.ParseLang("es"); got != langES {
		t.Errorf("expected ParseLang('es') to return %q, got %q", langES, got)
	}

	msg := i18n.T(langES, i18n.KeyNotFound)
	if msg != "El recurso solicitado no fue encontrado" {
		t.Errorf("expected Spanish translation, got %q", msg)
	}
}

package i18n

import "sync"

// Key represents a message key in the i18n catalog
type Key string

// All message keys used across the API
const (
	// Generic
	KeyInternalError    Key = "error.internal"
	KeyNotFound         Key = "error.not_found"
	KeyBadRequest       Key = "error.bad_request"
	KeyUnauthorized     Key = "error.unauthorized"
	KeyForbidden        Key = "error.forbidden"
	KeyValidationFailed Key = "error.validation_failed"

	// Party
	KeyPartyNotFound     Key = "party.not_found"
	KeyPartyCreated      Key = "party.created"
	KeyPartyNameRequired Key = "party.name_required"
	KeyPartyTypeInvalid  Key = "party.type_invalid"

	// Spatial Unit
	KeySpatialUnitNotFound Key = "spatial_unit.not_found"
	KeySpatialUnitCreated  Key = "spatial_unit.created"
	KeyGeomRequired        Key = "spatial_unit.geom_required"
	KeyGeomInvalid         Key = "spatial_unit.geom_invalid"
	KeyOverlapDetected     Key = "spatial_unit.overlap_detected"

	// Claim
	KeyClaimNotFound         Key = "claim.not_found"
	KeyClaimCreated          Key = "claim.created"
	KeyClaimStatusInvalid    Key = "claim.status_invalid"
	KeyClaimTransitionDenied Key = "claim.transition_denied"
	KeyObjectionPeriodActive Key = "claim.objection_period_active"

	// Source / Evidence
	KeySourceNotFound     Key = "source.not_found"
	KeySourceUploaded     Key = "source.uploaded"
	KeySourceHashMismatch Key = "source.hash_mismatch"

	// Witness
	KeyWitnessNotFound     Key = "witness.not_found"
	KeyWitnessCreated      Key = "witness.created"
	KeyWitnessNameRequired Key = "witness.name_required"
)

var (
	mu      sync.RWMutex
	catalog = map[Lang]map[Key]string{}
)

func init() {
	RegisterLanguage(LangFrench, frenchTranslations)
	RegisterLanguage(LangHaitianCreole, haitianCreoleTranslations)
}

// RegisterLanguage registers or adds translations for a language code
// Adheres strictly to the Open/Closed Principle (OCP):
// New languages can be registered at startup or runtime
// without modifying any existing translation logic or core functions
func RegisterLanguage(lang Lang, translations map[Key]string) {
	mu.Lock()
	defer mu.Unlock()

	if catalog[lang] == nil {
		catalog[lang] = make(map[Key]string)
	}
	for k, v := range translations {
		catalog[lang][k] = v
	}
}

// IsSupported checks if a language has registered translations in the catalog
func IsSupported(lang Lang) bool {
	mu.RLock()
	defer mu.RUnlock()
	_, exists := catalog[lang]
	return exists
}

// SupportedLanguages returns a slice of all currently registered language codes
func SupportedLanguages() []Lang {
	mu.RLock()
	defer mu.RUnlock()

	langs := make([]Lang, 0, len(catalog))
	for l := range catalog {
		langs = append(langs, l)
	}
	return langs
}

var frenchTranslations = map[Key]string{
	// Generic
	KeyInternalError:    "Une erreur interne est survenue. Veuillez réessayer plus tard",
	KeyNotFound:         "La ressource demandée est introuvable",
	KeyBadRequest:       "La requête est malformée ou contient des données invalides",
	KeyUnauthorized:     "Authentification requise pour accéder à cette ressource",
	KeyForbidden:        "Vous n'avez pas l'autorisation d'effectuer cette action",
	KeyValidationFailed: "La validation des données a échoué",

	// Party
	KeyPartyNotFound:     "Le titulaire foncier est introuvable",
	KeyPartyCreated:      "Le titulaire foncier a été créé avec succès",
	KeyPartyNameRequired: "Le prénom est obligatoire",
	KeyPartyTypeInvalid:  "Le type de titulaire est invalide",

	// Spatial Unit
	KeySpatialUnitNotFound: "La parcelle foncière est introuvable",
	KeySpatialUnitCreated:  "La parcelle foncière a été créée avec succès",
	KeyGeomRequired:        "La géométrie de la parcelle est obligatoire",
	KeyGeomInvalid:         "La géométrie fournie est invalide ou non reconnue",
	KeyOverlapDetected:     "La parcelle soumise chevauche une parcelle existante enregistrée",

	// Claim
	KeyClaimNotFound:         "La demande d'enregistrement est introuvable",
	KeyClaimCreated:          "La demande d'enregistrement a été créée avec succès",
	KeyClaimStatusInvalid:    "Le statut de la demande est invalide",
	KeyClaimTransitionDenied: "La transition de statut demandée n'est pas autorisée",
	KeyObjectionPeriodActive: "La période d'opposition publique est toujours en cours",

	// Source
	KeySourceNotFound:     "La preuve documentaire est introuvable",
	KeySourceUploaded:     "La preuve documentaire a été téléchargée avec succès",
	KeySourceHashMismatch: "L'intégrité du document est compromise : l'empreinte SHA-256 ne correspond pas",

	// Witness
	KeyWitnessNotFound:     "Le témoin de bornage est introuvable",
	KeyWitnessCreated:      "Le témoin de bornage a été enregistré avec succès",
	KeyWitnessNameRequired: "Le nom complet du témoin est obligatoire",
}

var haitianCreoleTranslations = map[Key]string{
	// Jeneral
	KeyInternalError:    "Gen yon pwoblèm entèn. Tanpri eseye ankò pita",
	KeyNotFound:         "Nou pa jwenn resous ou mande a",
	KeyBadRequest:       "Demann nan pa valab oswa li gen done ki pa valab",
	KeyUnauthorized:     "Ou dwe konekte pou w ka jwenn aksè ak resous sa a",
	KeyForbidden:        "Ou pa gen pèmisyon pou w fè aksyon sa a",
	KeyValidationFailed: "Done yo pa pase verifikasyon an",

	// Pati
	KeyPartyNotFound:     "Nou pa jwenn mèt tè a",
	KeyPartyCreated:      "Mèt tè a te anrejistre avèk siksè",
	KeyPartyNameRequired: "Prenon an obligatwa",
	KeyPartyTypeInvalid:  "Kalite mèt tè a pa valab",

	// Pasèl tè
	KeySpatialUnitNotFound: "Nou pa jwenn pasèl tè a",
	KeySpatialUnitCreated:  "Pasèl tè a te anrejistre avèk siksè",
	KeyGeomRequired:        "Jeyometri pasèl la obligatwa",
	KeyGeomInvalid:         "Jeyometri yo bay la pa valab oswa nou pa rekonèt li",
	KeyOverlapDetected:     "Pasèl ou soumèt la sipèpoze ak yon pasèl ki deja anrejistre",

	// Demann
	KeyClaimNotFound:         "Nou pa jwenn demann anrejistreman an",
	KeyClaimCreated:          "Demann anrejistreman an te kreye avèk siksè",
	KeyClaimStatusInvalid:    "Estati demann nan pa valab",
	KeyClaimTransitionDenied: "Ou pa gen pèmisyon pou chanje estati demann sa a",
	KeyObjectionPeriodActive: "Peryòd pou fè opozisyon piblik la poko fini",

	// Sous
	KeySourceNotFound:     "Nou pa jwenn prèv dokimantè a",
	KeySourceUploaded:     "Prèv dokimantè a te telechaje avèk siksè",
	KeySourceHashMismatch: "Dokiman an pa sanble ak dokiman orijinal la: anprent SHA-256 la pa koresponn",

	// Temwen
	KeyWitnessNotFound:     "Nou pa jwenn temwen limit tè a",
	KeyWitnessCreated:      "Temwen limit tè a te anrejistre avèk siksè",
	KeyWitnessNameRequired: "Non temwen an obligatwa",
}

// Translate returns the translated string for the given language and key
// Thread safe and falls back to DefaultLang if missing
func Translate(lang Lang, key Key) string {
	mu.RLock()
	defer mu.RUnlock()

	if translations, ok := catalog[lang]; ok {
		if msg, ok := translations[key]; ok {
			return msg
		}
	}
	// Fallback to DefaultLang
	if translations, ok := catalog[DefaultLang]; ok {
		if msg, ok := translations[key]; ok {
			return msg
		}
	}
	// Last resort: return key name to surface missing translations
	return string(key)
}

// T is a shorthand alias for Translate
func T(lang Lang, key Key) string {
	return Translate(lang, key)
}

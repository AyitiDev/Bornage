package i18n

// Lang represents a supported language code
type Lang string

const (
	LangFrench        Lang = "fr"
	LangHaitianCreole Lang = "ht"

	// DefaultLang is the fallback language per system requirement
	DefaultLang = LangFrench
)

// IsSupported returns true if the language is registered in the i18n catalog
func (l Lang) IsSupported() bool {
	return IsSupported(l)
}

// ParseLang converts a raw string to a Lang, falling back to DefaultLang
func ParseLang(raw string) Lang {
	l := Lang(raw)
	if l.IsSupported() {
		return l
	}
	return DefaultLang
}

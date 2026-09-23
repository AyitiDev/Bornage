package middleware

import (
	"strings"

	"github.com/AyitiDev/Bornage/api/internal/i18n"
	"github.com/gofiber/fiber/v2"
)

// LangKey is the Fiber locals key used to store the resolved language in the request context
const LangKey = "lang"

// I18n extracts the preferred language from:
//  1. Query parameter: ?lang=<code>
//  2. Accept Language HTTP header (first tag only)
//  3. Fallback: DefaultLang
//
// The resolved Lang is stored in Fiber locals and accessible via GetLang(c)
func I18n() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Priority 1: explicit ?lang= query param
		if raw := c.Query("lang"); raw != "" {
			c.Locals(LangKey, i18n.ParseLang(raw))
			return c.Next()
		}

		// Priority 2: Accept Language header (parse first tag, ignore quality values)
		if header := c.Get(fiber.HeaderAcceptLanguage); header != "" {
			// Accept Language: ht-HT,ht;q=0.9,fr;q=0.8  -> extract "ht-HT" -> take "ht"
			first := strings.SplitN(header, ",", 2)[0]                 // "ht-HT"
			tag := strings.SplitN(strings.TrimSpace(first), "-", 2)[0] // "ht"
			tag = strings.SplitN(tag, ";", 2)[0]                       // strip quality value if present
			c.Locals(LangKey, i18n.ParseLang(strings.TrimSpace(tag)))
			return c.Next()
		}

		// Priority 3: Fallback to DefaultLang
		c.Locals(LangKey, i18n.DefaultLang)
		return c.Next()
	}
}

// GetLang retrieves the resolved language from the Fiber request context
// Always returns a valid, supported Lang (never empty)
func GetLang(c *fiber.Ctx) i18n.Lang {
	if lang, ok := c.Locals(LangKey).(i18n.Lang); ok {
		return lang
	}
	return i18n.DefaultLang
}

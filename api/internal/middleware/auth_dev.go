package middleware

// WARNING: DEVELOPMENT ONLY. DO NOT USE IN PRODUCTION
// This middleware relies on untrusted HTTP headers (X-Dev-User-ID, X-Dev-Role)
// to bypass authenticating credentials during local development and testing
// IT IS STRICTLY FORBIDDEN TO MOUNT OR USE THIS MIDDLEWARE IN PRODUCTION

import (
	"os"

	"github.com/AyitiDev/Bornage/api/internal/domain"
	"github.com/AyitiDev/Bornage/api/internal/i18n"
	"github.com/gofiber/fiber/v2"
)

const (
	// HeaderDevUserID is the HTTP request header used in dev to pass user ID
	HeaderDevUserID = "X-Dev-User-ID"
	// HeaderDevRole is the HTTP request header used in dev to pass user role
	HeaderDevRole = "X-Dev-Role"

	// DevUserIDKey is the Fiber locals key used to store the development user ID
	DevUserIDKey = "dev_user_id"
	// DevRoleKey is the Fiber locals key used to store the development role
	DevRoleKey = "dev_role"
)

// DevAuthOptions holds configuration options for development authentication middleware
type DevAuthOptions struct {
	// Env specifies the runtime environment (e.x. "development", "production")
	// If empty, defaults to reading the ENV environment variable, falling back to "development"
	Env string
}

// RequireDevAuth creates a Fiber handler that authenticates requests using development headers
// If env is "development", it validates the presence of X-Dev-User-ID header
// If env is NOT "development" (e.x. "production"), it unconditionally blocks requests with HTTP 401 Unauthorized
func RequireDevAuth(opts ...DevAuthOptions) fiber.Handler {
	env := "development"
	if len(opts) > 0 && opts[0].Env != "" {
		env = opts[0].Env
	} else if e := os.Getenv("ENV"); e != "" {
		env = e
	}

	return func(c *fiber.Ctx) error {
		lang := GetLang(c)

		// Security guard: Dev auth must NEVER be usable outside of development environment
		if env != "development" {
			return c.Status(fiber.StatusUnauthorized).JSON(
				domain.NewErrorResponse(lang, i18n.KeyUnauthorized, nil),
			)
		}

		userID := c.Get(HeaderDevUserID)
		if userID == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(
				domain.NewErrorResponse(lang, i18n.KeyUnauthorized, nil),
			)
		}

		c.Locals(DevUserIDKey, userID)

		if role := c.Get(HeaderDevRole); role != "" {
			c.Locals(DevRoleKey, role)
		}

		return c.Next()
	}
}

// RequireDevRole creates a Fiber handler that enforces role based access control based on dev headers
// Returns HTTP 403 Forbidden if the user's role is missing or not contained in allowedRoles
func RequireDevRole(allowedRoles ...string) fiber.Handler {
	allowedMap := make(map[string]bool, len(allowedRoles))
	for _, r := range allowedRoles {
		allowedMap[r] = true
	}

	return func(c *fiber.Ctx) error {
		lang := GetLang(c)
		userRole := GetDevRole(c)

		if userRole == "" || !allowedMap[userRole] {
			return c.Status(fiber.StatusForbidden).JSON(
				domain.NewErrorResponse(lang, i18n.KeyForbidden, nil),
			)
		}

		return c.Next()
	}
}

// GetDevUserID retrieves the development user ID from Fiber context
func GetDevUserID(c *fiber.Ctx) string {
	if val, ok := c.Locals(DevUserIDKey).(string); ok {
		return val
	}
	return ""
}

// GetDevRole retrieves the development user role from Fiber context
func GetDevRole(c *fiber.Ctx) string {
	if val, ok := c.Locals(DevRoleKey).(string); ok {
		return val
	}
	return ""
}

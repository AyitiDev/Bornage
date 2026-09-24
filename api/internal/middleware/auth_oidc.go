package middleware

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"strings"

	"github.com/AyitiDev/Bornage/api/internal/domain"
	"github.com/AyitiDev/Bornage/api/internal/i18n"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

const (
	// UserIDKey is the Fiber locals key for storing the authenticated user ID
	UserIDKey = "user_id"
	// UserRolesKey is the Fiber locals key for storing user roles
	UserRolesKey = "user_roles"
	// ClaimsKey is the Fiber locals key for storing parsed JWT claims
	ClaimsKey = "jwt_claims"
)

// OIDCOptions contains configuration settings for JWT/OIDC authentication
type OIDCOptions struct {
	// Secret is the HMAC secret key used for verifying HS256/HS384/HS512 tokens
	Secret []byte
	// PublicKey is the RSA public key used for verifying RS256/RS384/RS512 tokens
	PublicKey *rsa.PublicKey
	// Issuer is the expected token issuer ("iss" claim)
	Issuer string
	// Audience is the expected token audience ("aud" claim)
	Audience string
}

// CustomClaims wraps standard claims and extracts Keycloak / OIDC roles
type CustomClaims struct {
	jwt.RegisteredClaims
	PreferredUsername string                 `json:"preferred_username,omitempty"`
	Email             string                 `json:"email,omitempty"`
	RealmAccess       RealmAccessClaim       `json:"realm_access,omitempty"`
	ResourceAccess    map[string]AccessClaim `json:"resource_access,omitempty"`
	Roles             []string               `json:"roles,omitempty"`
}

type RealmAccessClaim struct {
	Roles []string `json:"roles,omitempty"`
}

type AccessClaim struct {
	Roles []string `json:"roles,omitempty"`
}

// RequireOIDCAuth returns a Fiber middleware that validates JWT Bearer tokens
func RequireOIDCAuth(opts OIDCOptions) fiber.Handler {
	return func(c *fiber.Ctx) error {
		lang := GetLang(c)
		authHeader := c.Get(fiber.HeaderAuthorization)

		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(
				domain.NewErrorResponse(lang, i18n.KeyUnauthorized, nil),
			)
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(
				domain.NewErrorResponse(lang, i18n.KeyUnauthorized, "Malformed Authorization header"),
			)
		}

		tokenString := strings.TrimSpace(parts[1])

		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
			switch t.Method.(type) {
			case *jwt.SigningMethodHMAC:
				if len(opts.Secret) == 0 {
					return nil, errors.New("HMAC secret not configured")
				}
				return opts.Secret, nil
			case *jwt.SigningMethodRSA:
				if opts.PublicKey == nil {
					return nil, errors.New("RSA public key not configured")
				}
				return opts.PublicKey, nil
			default:
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(
				domain.NewErrorResponse(lang, i18n.KeyUnauthorized, nil),
			)
		}

		claims, ok := token.Claims.(*CustomClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(
				domain.NewErrorResponse(lang, i18n.KeyUnauthorized, nil),
			)
		}

		// Validate issuer if specified
		if opts.Issuer != "" && claims.Issuer != opts.Issuer {
			return c.Status(fiber.StatusUnauthorized).JSON(
				domain.NewErrorResponse(lang, i18n.KeyUnauthorized, "Invalid token issuer"),
			)
		}

		// Validate audience if specified
		if opts.Audience != "" {
			audMatched := false
			for _, aud := range claims.Audience {
				if aud == opts.Audience {
					audMatched = true
					break
				}
			}
			if !audMatched {
				return c.Status(fiber.StatusUnauthorized).JSON(
					domain.NewErrorResponse(lang, i18n.KeyUnauthorized, "Invalid token audience"),
				)
			}
		}

		// Extract subject / user ID
		userID := claims.Subject
		if userID == "" {
			userID = claims.PreferredUsername
		}

		// Aggregate roles from direct roles array, realm_access, and resource_access
		rolesMap := make(map[string]bool)
		for _, r := range claims.Roles {
			rolesMap[r] = true
		}
		for _, r := range claims.RealmAccess.Roles {
			rolesMap[r] = true
		}
		for _, res := range claims.ResourceAccess {
			for _, r := range res.Roles {
				rolesMap[r] = true
			}
		}

		roles := make([]string, 0, len(rolesMap))
		for r := range rolesMap {
			roles = append(roles, r)
		}

		c.Locals(UserIDKey, userID)
		c.Locals(UserRolesKey, roles)
		c.Locals(ClaimsKey, claims)

		return c.Next()
	}
}

// RequireOIDCRole returns a Fiber middleware enforcing role checks against OIDC claims
func RequireOIDCRole(allowedRoles ...string) fiber.Handler {
	allowedMap := make(map[string]bool, len(allowedRoles))
	for _, r := range allowedRoles {
		allowedMap[r] = true
	}

	return func(c *fiber.Ctx) error {
		lang := GetLang(c)
		userRoles := GetUserRoles(c)

		hasRole := false
		for _, role := range userRoles {
			if allowedMap[role] {
				hasRole = true
				break
			}
		}

		if !hasRole {
			return c.Status(fiber.StatusForbidden).JSON(
				domain.NewErrorResponse(lang, i18n.KeyForbidden, nil),
			)
		}

		return c.Next()
	}
}

// GetUserID retrieves the authenticated user ID from context
func GetUserID(c *fiber.Ctx) string {
	if val, ok := c.Locals(UserIDKey).(string); ok {
		return val
	}
	return ""
}

// GetUserRoles retrieves the user is roles from context
func GetUserRoles(c *fiber.Ctx) []string {
	if val, ok := c.Locals(UserRolesKey).([]string); ok {
		return val
	}
	return nil
}

// HasRole checks if the current context user possesses a specific role
func HasRole(c *fiber.Ctx, role string) bool {
	roles := GetUserRoles(c)
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

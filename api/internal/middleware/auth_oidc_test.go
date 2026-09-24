package middleware_test

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AyitiDev/Bornage/api/internal/domain"
	"github.com/AyitiDev/Bornage/api/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func generateHMACJWT(t *testing.T, secret []byte, claims middleware.CustomClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(secret)
	if err != nil {
		t.Fatalf("failed to sign HMAC JWT: %v", err)
	}
	return tokenStr
}

func generateRSAJWT(t *testing.T, key *rsa.PrivateKey, claims middleware.CustomClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenStr, err := token.SignedString(key)
	if err != nil {
		t.Fatalf("failed to sign RSA JWT: %v", err)
	}
	return tokenStr
}

func TestRequireOIDCAuth_ValidHMAC(t *testing.T) {
	secret := []byte("super-secret-key-1234567890-bornage")
	opts := middleware.OIDCOptions{
		Secret: secret,
	}

	app := fiber.New()
	app.Use(middleware.I18n())
	app.Use(middleware.RequireOIDCAuth(opts))
	app.Get("/test", func(c *fiber.Ctx) error {
		userID := middleware.GetUserID(c)
		roles := middleware.GetUserRoles(c)
		return c.JSON(fiber.Map{
			"user_id": userID,
			"roles":   roles,
		})
	})

	claims := middleware.CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "usr_gov_999",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
		RealmAccess: middleware.RealmAccessClaim{
			Roles: []string{"registrar", "admin"},
		},
	}
	tokenStr := generateHMACJWT(t, secret, claims)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var data map[string]interface{}
	_ = json.Unmarshal(body, &data)

	if data["user_id"] != "usr_gov_999" {
		t.Errorf("expected user_id 'usr_gov_999', got '%v'", data["user_id"])
	}
}

func TestRequireOIDCAuth_ValidRSA(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key pair: %v", err)
	}

	opts := middleware.OIDCOptions{
		PublicKey: &privateKey.PublicKey,
	}

	app := fiber.New()
	app.Use(middleware.I18n())
	app.Use(middleware.RequireOIDCAuth(opts))
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	claims := middleware.CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "usr_rsa_111",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}
	tokenStr := generateRSAJWT(t, privateKey, claims)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 for valid RSA token, got %d", resp.StatusCode)
	}
}

func TestRequireOIDCAuth_ExpiredToken(t *testing.T) {
	secret := []byte("super-secret-key-1234567890-bornage")
	opts := middleware.OIDCOptions{Secret: secret}

	app := fiber.New()
	app.Use(middleware.I18n())
	app.Use(middleware.RequireOIDCAuth(opts))
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	claims := middleware.CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "usr_expired",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // Expired 1h ago
		},
	}
	tokenStr := generateHMACJWT(t, secret, claims)

	req := httptest.NewRequest("GET", "/test?lang=ht", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status 401 Unauthorized for expired token, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var errResp domain.ErrorResponse
	_ = json.Unmarshal(body, &errResp)

	if errResp.Code != "error.unauthorized" {
		t.Errorf("expected code 'error.unauthorized', got '%s'", errResp.Code)
	}
	if errResp.Message != "Ou dwe konekte pou w ka jwenn aksè ak resous sa a" {
		t.Errorf("expected Haitian Creole unauthorized message, got '%s'", errResp.Message)
	}
}

func TestRequireOIDCAuth_MissingHeader(t *testing.T) {
	opts := middleware.OIDCOptions{Secret: []byte("secret")}

	app := fiber.New()
	app.Use(middleware.I18n())
	app.Use(middleware.RequireOIDCAuth(opts))
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	req := httptest.NewRequest("GET", "/test", nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status 401 for missing Authorization header, got %d", resp.StatusCode)
	}
}

func TestRequireOIDCRole_SuccessAndFailure(t *testing.T) {
	secret := []byte("secret")
	opts := middleware.OIDCOptions{Secret: secret}

	app := fiber.New()
	app.Use(middleware.I18n())
	app.Use(middleware.RequireOIDCAuth(opts))
	app.Use(middleware.RequireOIDCRole("cadastral_officer"))
	app.Get("/test", func(c *fiber.Ctx) error {
		hasRole := middleware.HasRole(c, "cadastral_officer")
		return c.JSON(fiber.Map{"has_role": hasRole})
	})

	// Authorized role
	claimsAuth := middleware.CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "usr_cadastre",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
		Roles: []string{"cadastral_officer"},
	}
	tokAuth := generateHMACJWT(t, secret, claimsAuth)

	req1 := httptest.NewRequest("GET", "/test", nil)
	req1.Header.Set("Authorization", "Bearer "+tokAuth)

	resp1, _ := app.Test(req1)
	if resp1.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK for cadastral_officer, got %d", resp1.StatusCode)
	}

	// Unauthorized role
	claimsUnauth := middleware.CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "usr_guest",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
		Roles: []string{"guest"},
	}
	tokUnauth := generateHMACJWT(t, secret, claimsUnauth)

	req2 := httptest.NewRequest("GET", "/test", nil)
	req2.Header.Set("Authorization", "Bearer "+tokUnauth)

	resp2, _ := app.Test(req2)
	if resp2.StatusCode != http.StatusForbidden {
		t.Errorf("expected 403 Forbidden for guest, got %d", resp2.StatusCode)
	}
}

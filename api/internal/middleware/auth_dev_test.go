package middleware_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AyitiDev/Bornage/api/internal/domain"
	"github.com/AyitiDev/Bornage/api/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func TestRequireDevAuth_DevelopmentMode_Success(t *testing.T) {
	app := fiber.New()
	app.Use(middleware.I18n())
	app.Use(middleware.RequireDevAuth(middleware.DevAuthOptions{Env: "development"}))
	app.Get("/test", func(c *fiber.Ctx) error {
		userID := middleware.GetDevUserID(c)
		role := middleware.GetDevRole(c)
		return c.JSON(fiber.Map{
			"user_id": userID,
			"role":    role,
		})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set(middleware.HeaderDevUserID, "usr_12345")
	req.Header.Set(middleware.HeaderDevRole, "surveyor")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var data map[string]string
	_ = json.Unmarshal(body, &data)

	if data["user_id"] != "usr_12345" {
		t.Errorf("expected user_id 'usr_12345', got '%s'", data["user_id"])
	}
	if data["role"] != "surveyor" {
		t.Errorf("expected role 'surveyor', got '%s'", data["role"])
	}
}

func TestRequireDevAuth_DevelopmentMode_MissingHeader(t *testing.T) {
	app := fiber.New()
	app.Use(middleware.I18n())
	app.Use(middleware.RequireDevAuth(middleware.DevAuthOptions{Env: "development"}))
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	req := httptest.NewRequest("GET", "/test?lang=fr", nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status 401 Unauthorized, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var errResp domain.ErrorResponse
	if err := json.Unmarshal(body, &errResp); err != nil {
		t.Fatalf("failed to unmarshal error response: %v", err)
	}

	if errResp.Code != "error.unauthorized" {
		t.Errorf("expected code 'error.unauthorized', got '%s'", errResp.Code)
	}
	if errResp.Message != "Authentification requise pour accéder à cette ressource" {
		t.Errorf("expected French unauthorized message, got '%s'", errResp.Message)
	}
}

func TestRequireDevAuth_ProductionMode_Blocks(t *testing.T) {
	app := fiber.New()
	app.Use(middleware.I18n())
	app.Use(middleware.RequireDevAuth(middleware.DevAuthOptions{Env: "production"}))
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set(middleware.HeaderDevUserID, "usr_12345")
	req.Header.Set(middleware.HeaderDevRole, "admin")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status 401 Unauthorized in production, got %d", resp.StatusCode)
	}
}

func TestRequireDevRole_SuccessAndFailure(t *testing.T) {
	app := fiber.New()
	app.Use(middleware.I18n())
	app.Use(middleware.RequireDevAuth(middleware.DevAuthOptions{Env: "development"}))
	app.Use(middleware.RequireDevRole("admin", "registrar"))
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("access_granted")
	})

	// Test 1: User with allowed role "admin"
	req1 := httptest.NewRequest("GET", "/test", nil)
	req1.Header.Set(middleware.HeaderDevUserID, "usr_1")
	req1.Header.Set(middleware.HeaderDevRole, "admin")

	resp1, err := app.Test(req1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp1.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 for admin, got %d", resp1.StatusCode)
	}

	// Test 2: User with disallowed role surveyor
	req2 := httptest.NewRequest("GET", "/test?lang=ht", nil)
	req2.Header.Set(middleware.HeaderDevUserID, "usr_2")
	req2.Header.Set(middleware.HeaderDevRole, "surveyor")

	resp2, err := app.Test(req2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp2.StatusCode != http.StatusForbidden {
		t.Errorf("expected status 403 Forbidden for surveyor, got %d", resp2.StatusCode)
	}

	body2, _ := io.ReadAll(resp2.Body)
	var errResp domain.ErrorResponse
	if err := json.Unmarshal(body2, &errResp); err != nil {
		t.Fatalf("failed to unmarshal error response: %v", err)
	}

	if errResp.Code != "error.forbidden" {
		t.Errorf("expected code 'error.forbidden', got '%s'", errResp.Code)
	}
	if errResp.Message != "Ou pa gen pèmisyon pou w fè aksyon sa a" {
		t.Errorf("expected Haitian Creole forbidden message, got '%s'", errResp.Message)
	}

	// Test 3: User with no role specified
	req3 := httptest.NewRequest("GET", "/test", nil)
	req3.Header.Set(middleware.HeaderDevUserID, "usr_3")

	resp3, err := app.Test(req3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp3.StatusCode != http.StatusForbidden {
		t.Errorf("expected status 403 Forbidden for missing role, got %d", resp3.StatusCode)
	}
}

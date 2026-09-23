package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AyitiDev/Bornage/api/internal/i18n"
	"github.com/AyitiDev/Bornage/api/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func newTestApp() *fiber.App {
	app := fiber.New()
	app.Use(middleware.I18n())
	app.Get("/lang", func(c *fiber.Ctx) error {
		return c.SendString(string(middleware.GetLang(c)))
	})
	return app
}

func TestI18n_QueryParam_fr(t *testing.T) {
	app := newTestApp()
	req := httptest.NewRequest(http.MethodGet, "/lang?lang=fr", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestI18n_QueryParam_ht(t *testing.T) {
	app := newTestApp()
	req := httptest.NewRequest(http.MethodGet, "/lang?lang=ht", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestI18n_QueryParam_Unsupported_FallsBackToFr(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/lang?lang=en", nil)

	var capturedLang i18n.Lang
	app := fiber.New()
	app.Use(middleware.I18n())
	app.Get("/lang", func(c *fiber.Ctx) error {
		capturedLang = middleware.GetLang(c)
		return c.SendString(string(capturedLang))
	})

	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if capturedLang != i18n.DefaultLang {
		t.Errorf("expected fallback to %q, got %q", i18n.DefaultLang, capturedLang)
	}
}

func TestI18n_AcceptLanguage_ht(t *testing.T) {
	var capturedLang i18n.Lang
	app := fiber.New()
	app.Use(middleware.I18n())
	app.Get("/lang", func(c *fiber.Ctx) error {
		capturedLang = middleware.GetLang(c)
		return c.SendString(string(capturedLang))
	})

	req := httptest.NewRequest(http.MethodGet, "/lang", nil)
	req.Header.Set("Accept-Language", "ht-HT,ht;q=0.9,fr;q=0.8")

	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if capturedLang != i18n.LangHaitianCreole {
		t.Errorf("expected ht, got %q", capturedLang)
	}
}

func TestI18n_AcceptLanguage_fr(t *testing.T) {
	var capturedLang i18n.Lang
	app := fiber.New()
	app.Use(middleware.I18n())
	app.Get("/lang", func(c *fiber.Ctx) error {
		capturedLang = middleware.GetLang(c)
		return c.SendString(string(capturedLang))
	})

	req := httptest.NewRequest(http.MethodGet, "/lang", nil)
	req.Header.Set("Accept-Language", "fr-FR,fr;q=0.9")

	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if capturedLang != i18n.LangFrench {
		t.Errorf("expected fr, got %q", capturedLang)
	}
}

func TestI18n_NoHeader_FallsBackToFr(t *testing.T) {
	var capturedLang i18n.Lang
	app := fiber.New()
	app.Use(middleware.I18n())
	app.Get("/lang", func(c *fiber.Ctx) error {
		capturedLang = middleware.GetLang(c)
		return c.SendString(string(capturedLang))
	})

	req := httptest.NewRequest(http.MethodGet, "/lang", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if capturedLang != i18n.DefaultLang {
		t.Errorf("expected default lang %q, got %q", i18n.DefaultLang, capturedLang)
	}
}

func TestI18n_QueryParamTakesPriorityOverHeader(t *testing.T) {
	var capturedLang i18n.Lang
	app := fiber.New()
	app.Use(middleware.I18n())
	app.Get("/lang", func(c *fiber.Ctx) error {
		capturedLang = middleware.GetLang(c)
		return c.SendString(string(capturedLang))
	})

	req := httptest.NewRequest(http.MethodGet, "/lang?lang=ht", nil)
	req.Header.Set("Accept-Language", "fr-FR") // header says fr but ?lang=ht should win

	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if capturedLang != i18n.LangHaitianCreole {
		t.Errorf("expected ?lang=ht to win over Accept-Language header, got %q", capturedLang)
	}
}

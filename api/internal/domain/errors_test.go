package domain_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/AyitiDev/Bornage/api/internal/domain"
	"github.com/AyitiDev/Bornage/api/internal/i18n"
)

func TestNewErrorResponse_French(t *testing.T) {
	resp := domain.NewErrorResponse(i18n.LangFrench, i18n.KeyNotFound, nil)

	if resp.Code != "error.not_found" {
		t.Errorf("expected code 'error.not_found', got '%s'", resp.Code)
	}
	if resp.Message != "La ressource demandée est introuvable" {
		t.Errorf("expected French message, got '%s'", resp.Message)
	}
	if resp.Details != nil {
		t.Errorf("expected nil details, got %v", resp.Details)
	}
}

func TestNewErrorResponse_HaitianCreole(t *testing.T) {
	details := map[string]string{"field": "first_name", "reason": "required"}
	resp := domain.NewErrorResponse(i18n.LangHaitianCreole, i18n.KeyPartyNameRequired, details)

	if resp.Code != "party.name_required" {
		t.Errorf("expected code 'party.name_required', got '%s'", resp.Code)
	}
	if resp.Message != "Prenon an obligatwa" {
		t.Errorf("expected Haitian Creole message, got '%s'", resp.Message)
	}
	if resp.Details == nil {
		t.Fatal("expected non-nil details")
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal ErrorResponse: %v", err)
	}

	expectedJSON := `{"code":"party.name_required","message":"Prenon an obligatwa","details":{"field":"first_name","reason":"required"}}`
	if string(data) != expectedJSON {
		t.Errorf("expected JSON %s, got %s", expectedJSON, string(data))
	}
}

func TestErrorResponse_OmitEmptyDetails(t *testing.T) {
	resp := domain.NewErrorResponse(i18n.LangFrench, i18n.KeyUnauthorized, nil)

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal ErrorResponse: %v", err)
	}

	expectedJSON := `{"code":"error.unauthorized","message":"Authentification requise pour accéder à cette ressource"}`
	if string(data) != expectedJSON {
		t.Errorf("expected JSON %s, got %s", expectedJSON, string(data))
	}
}

func TestDomainError_Methods(t *testing.T) {
	underlyingErr := errors.New("db connection failure")
	domErr := domain.NewDomainErrorWrapped(500, i18n.KeyInternalError, "details_info", underlyingErr)

	if domErr.HTTPStatus != 500 {
		t.Errorf("expected status 500, got %d", domErr.HTTPStatus)
	}
	if !errors.Is(domErr, underlyingErr) {
		t.Errorf("expected unwrapped error to equal underlying error")
	}

	resp := domErr.ToResponse(i18n.LangFrench)
	if resp.Code != "error.internal" {
		t.Errorf("expected code 'error.internal', got '%s'", resp.Code)
	}
}

func TestDomainError_Constructors(t *testing.T) {
	tests := []struct {
		name       string
		err        *domain.DomainError
		wantStatus int
		wantKey    i18n.Key
	}{
		{"NotFound", domain.ErrNotFound(i18n.KeyPartyNotFound, nil), 404, i18n.KeyPartyNotFound},
		{"BadRequest", domain.ErrBadRequest(i18n.KeyBadRequest, nil), 400, i18n.KeyBadRequest},
		{"Unauthorized", domain.ErrUnauthorized(i18n.KeyUnauthorized, nil), 401, i18n.KeyUnauthorized},
		{"Forbidden", domain.ErrForbidden(i18n.KeyForbidden, nil), 403, i18n.KeyForbidden},
		{"Internal", domain.ErrInternal(nil), 500, i18n.KeyInternalError},
		{"Validation", domain.ErrValidation(nil), 422, i18n.KeyValidationFailed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.HTTPStatus != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, tt.err.HTTPStatus)
			}
			if tt.err.Key != tt.wantKey {
				t.Errorf("expected key %s, got %s", tt.wantKey, tt.err.Key)
			}
		})
	}
}

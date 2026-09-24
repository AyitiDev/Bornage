package domain

import (
	"fmt"

	"github.com/AyitiDev/Bornage/api/internal/i18n"
)

// ErrorResponse represents the standardized JSON response structure for API errors
type ErrorResponse struct {
	Code    string `json:"code"`              // Internal error code / i18n key (e.g. "error.not_found")
	Message string `json:"message"`           // Localized human-readable error message
	Details any    `json:"details,omitempty"` // Optional error metadata or field validation details
}

// Error implements the standard Go error interface for ErrorResponse
func (e ErrorResponse) Error() string {
	return e.Message
}

// DomainError represents an internal domain level error carrying HTTP status and i18n key
type DomainError struct {
	HTTPStatus int
	Key        i18n.Key
	Details    any
	Err        error
}

// Error implements the standard Go error interface for DomainError
func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Key, e.Err)
	}
	return string(e.Key)
}

// Unwrap returns the underlying wrapped error, if any
func (e *DomainError) Unwrap() error {
	return e.Err
}

// ToResponse formats the DomainError into a localized ErrorResponse
func (e *DomainError) ToResponse(lang i18n.Lang) ErrorResponse {
	return NewErrorResponse(lang, e.Key, e.Details)
}

// NewErrorResponse builds a localized ErrorResponse using the specified language and i18n key
func NewErrorResponse(lang i18n.Lang, key i18n.Key, details any) ErrorResponse {
	return ErrorResponse{
		Code:    string(key),
		Message: i18n.T(lang, key),
		Details: details,
	}
}

// NewDomainError constructs a new DomainError instance
func NewDomainError(status int, key i18n.Key, details any) *DomainError {
	return &DomainError{
		HTTPStatus: status,
		Key:        key,
		Details:    details,
	}
}

// NewDomainErrorWrapped constructs a new DomainError instance wrapping an underlying error
func NewDomainErrorWrapped(status int, key i18n.Key, details any, err error) *DomainError {
	return &DomainError{
		HTTPStatus: status,
		Key:        key,
		Details:    details,
		Err:        err,
	}
}

// Standard domain error constructors for common HTTP status codes
func ErrNotFound(key i18n.Key, details any) *DomainError {
	return NewDomainError(404, key, details)
}

func ErrBadRequest(key i18n.Key, details any) *DomainError {
	return NewDomainError(400, key, details)
}

func ErrUnauthorized(key i18n.Key, details any) *DomainError {
	return NewDomainError(401, key, details)
}

func ErrForbidden(key i18n.Key, details any) *DomainError {
	return NewDomainError(403, key, details)
}

func ErrInternal(details any) *DomainError {
	return NewDomainError(500, i18n.KeyInternalError, details)
}

func ErrValidation(details any) *DomainError {
	return NewDomainError(422, i18n.KeyValidationFailed, details)
}

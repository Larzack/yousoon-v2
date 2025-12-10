// Package domain contains domain errors for the Discovery service.
package domain

import "errors"

// Domain errors
var (
	// Offer errors
	ErrOfferNotFound           = errors.New("offer not found")
	ErrOfferNotActive          = errors.New("offer is not active")
	ErrOfferNotAvailableNow    = errors.New("offer is not available at this time")
	ErrOfferFullyBooked        = errors.New("offer is fully booked")
	ErrOfferNotApproved        = errors.New("offer has not been approved")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrUserQuotaExceeded       = errors.New("user has exceeded their quota for this offer")
	ErrOfferAlreadyPublished   = errors.New("offer is already published")
	ErrOfferAlreadyArchived    = errors.New("offer is already archived")

	// Category errors
	ErrCategoryNotFound    = errors.New("category not found")
	ErrCategorySlugExists  = errors.New("category slug already exists")
	ErrCategoryHasChildren = errors.New("category has child categories")
	ErrCategoryHasOffers   = errors.New("category has associated offers")

	// Validation errors
	ErrInvalidDiscount = errors.New("invalid discount")
	ErrInvalidValidity = errors.New("invalid validity period")
	ErrInvalidSchedule = errors.New("invalid schedule")
	ErrInvalidQuota    = errors.New("invalid quota")
)

// ValidationError represents a validation error with field information.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e ValidationError) Error() string {
	return e.Field + ": " + e.Message
}

// NewValidationError creates a new validation error.
func NewValidationError(field, message string) ValidationError {
	return ValidationError{
		Field:   field,
		Message: message,
	}
}

// ValidationErrors represents multiple validation errors.
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

func (e ValidationErrors) Error() string {
	if len(e.Errors) == 0 {
		return "validation failed"
	}
	return e.Errors[0].Error()
}

// Add adds a validation error.
func (e *ValidationErrors) Add(field, message string) {
	e.Errors = append(e.Errors, ValidationError{
		Field:   field,
		Message: message,
	})
}

// HasErrors returns true if there are validation errors.
func (e ValidationErrors) HasErrors() bool {
	return len(e.Errors) > 0
}

// DomainError represents a domain-specific error with a code.
type DomainError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e DomainError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e DomainError) Unwrap() error {
	return e.Err
}

// NewDomainError creates a new domain error.
func NewDomainError(code, message string, err error) DomainError {
	return DomainError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Common error codes
const (
	ErrCodeNotFound         = "NOT_FOUND"
	ErrCodeValidation       = "VALIDATION_ERROR"
	ErrCodeUnauthorized     = "UNAUTHORIZED"
	ErrCodeForbidden        = "FORBIDDEN"
	ErrCodeConflict         = "CONFLICT"
	ErrCodeInternalError    = "INTERNAL_ERROR"
	ErrCodeInvalidOperation = "INVALID_OPERATION"
)

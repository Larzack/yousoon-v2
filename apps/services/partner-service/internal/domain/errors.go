package domain

import "errors"

// =============================================================================
// Domain Errors
// =============================================================================

var (
	// Partner errors
	ErrPartnerNotFound      = errors.New("partner not found")
	ErrPartnerNotVerified   = errors.New("partner is not verified")
	ErrAlreadyVerified      = errors.New("partner is already verified")
	ErrNoEstablishment      = errors.New("partner has no establishment")
	ErrInvalidPartnerStatus = errors.New("invalid partner status")

	// Company errors
	ErrCompanyNameRequired = errors.New("company name is required")
	ErrInvalidSIRET        = errors.New("invalid SIRET number")
	ErrSIRETAlreadyExists  = errors.New("SIRET already registered")

	// Contact errors
	ErrContactNameRequired  = errors.New("contact first and last name are required")
	ErrContactEmailRequired = errors.New("contact email is required")

	// Email errors
	ErrInvalidEmail = errors.New("invalid email format")

	// Geolocation errors
	ErrInvalidLongitude = errors.New("longitude must be between -180 and 180")
	ErrInvalidLatitude  = errors.New("latitude must be between -90 and 90")

	// Establishment errors
	ErrEstablishmentNotFound      = errors.New("establishment not found")
	ErrEstablishmentAlreadyExists = errors.New("establishment at this address already exists")

	// Team member errors
	ErrTeamMemberNotFound   = errors.New("team member not found")
	ErrTeamMemberExists     = errors.New("team member with this email already exists")
	ErrInvitationNotPending = errors.New("invitation is not pending")
	ErrInvalidTeamRole      = errors.New("invalid team role")

	// Permission errors
	ErrPermissionDenied = errors.New("permission denied")

	// Validation errors
	ErrValidation = errors.New("validation error")
)

// ValidationError represents a validation error with details.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// NewValidationError creates a new validation error.
func NewValidationError(field, message string) ValidationError {
	return ValidationError{
		Field:   field,
		Message: message,
	}
}

// Error returns the error message.
func (e ValidationError) Error() string {
	return e.Field + ": " + e.Message
}

// ValidationErrors is a collection of validation errors.
type ValidationErrors []ValidationError

// Error returns the error message.
func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		return "validation failed"
	}
	return e[0].Error()
}

// Add adds a validation error.
func (e *ValidationErrors) Add(field, message string) {
	*e = append(*e, NewValidationError(field, message))
}

// HasErrors returns true if there are validation errors.
func (e ValidationErrors) HasErrors() bool {
	return len(e) > 0
}

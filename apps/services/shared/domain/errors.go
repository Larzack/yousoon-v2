package domain

import (
	"errors"
	"fmt"
)

// =============================================================================
// Base Domain Errors
// =============================================================================

var (
	// ErrNotFound is returned when an entity is not found.
	ErrNotFound = errors.New("not found")

	// ErrAlreadyExists is returned when trying to create a duplicate entity.
	ErrAlreadyExists = errors.New("already exists")

	// ErrInvalidInput is returned when input validation fails.
	ErrInvalidInput = errors.New("invalid input")

	// ErrUnauthorized is returned when the user is not authenticated.
	ErrUnauthorized = errors.New("unauthorized")

	// ErrForbidden is returned when the user lacks permission.
	ErrForbidden = errors.New("forbidden")

	// ErrConflict is returned when there's a state conflict.
	ErrConflict = errors.New("conflict")

	// ErrInternal is returned for internal server errors.
	ErrInternal = errors.New("internal error")

	// ErrValidation is returned for validation errors.
	ErrValidation = errors.New("validation error")

	// ErrExpired is returned when something has expired (token, offer, etc).
	ErrExpired = errors.New("expired")

	// ErrQuotaExceeded is returned when a quota or limit is exceeded.
	ErrQuotaExceeded = errors.New("quota exceeded")
)

// =============================================================================
// Domain-specific Errors
// =============================================================================

// Identity Context Errors
var (
	ErrUserNotFound         = fmt.Errorf("user %w", ErrNotFound)
	ErrUserAlreadyExists    = fmt.Errorf("user %w", ErrAlreadyExists)
	ErrInvalidCredentials   = fmt.Errorf("invalid credentials: %w", ErrUnauthorized)
	ErrEmailNotVerified     = fmt.Errorf("email not verified: %w", ErrForbidden)
	ErrIdentityNotVerified  = fmt.Errorf("identity not verified: %w", ErrForbidden)
	ErrAlreadyVerified      = fmt.Errorf("already verified: %w", ErrConflict)
	ErrVerificationFailed   = fmt.Errorf("verification failed: %w", ErrValidation)
	ErrMaxVerificationTries = fmt.Errorf("max verification attempts reached: %w", ErrQuotaExceeded)
	ErrSubscriptionRequired = fmt.Errorf("subscription required: %w", ErrForbidden)
	ErrAlreadySubscribed    = fmt.Errorf("already subscribed: %w", ErrConflict)
	ErrSubscriptionNotFound = fmt.Errorf("subscription %w", ErrNotFound)
	ErrInvalidRefreshToken  = fmt.Errorf("invalid refresh token: %w", ErrUnauthorized)
	ErrTokenExpired         = fmt.Errorf("token %w", ErrExpired)
	ErrAccountSuspended     = fmt.Errorf("account suspended: %w", ErrForbidden)
	ErrAccountDeleted       = fmt.Errorf("account deleted: %w", ErrForbidden)
	ErrWeakPassword         = fmt.Errorf("password too weak: %w", ErrValidation)
	ErrInvalidEmail         = fmt.Errorf("invalid email: %w", ErrValidation)
	ErrInvalidPhone         = fmt.Errorf("invalid phone: %w", ErrValidation)
)

// Partner Context Errors
var (
	ErrPartnerNotFound       = fmt.Errorf("partner %w", ErrNotFound)
	ErrPartnerAlreadyExists  = fmt.Errorf("partner %w", ErrAlreadyExists)
	ErrPartnerNotVerified    = fmt.Errorf("partner not verified: %w", ErrForbidden)
	ErrEstablishmentNotFound = fmt.Errorf("establishment %w", ErrNotFound)
	ErrTeamMemberNotFound    = fmt.Errorf("team member %w", ErrNotFound)
	ErrTeamMemberExists      = fmt.Errorf("team member %w", ErrAlreadyExists)
	ErrInvalidSIRET          = fmt.Errorf("invalid SIRET: %w", ErrValidation)
	ErrNotTeamAdmin          = fmt.Errorf("not team admin: %w", ErrForbidden)
	ErrCannotRemoveSelf      = fmt.Errorf("cannot remove yourself: %w", ErrConflict)
	Err2FARequired           = fmt.Errorf("2FA required: %w", ErrForbidden)
)

// Discovery Context Errors
var (
	ErrOfferNotFound      = fmt.Errorf("offer %w", ErrNotFound)
	ErrOfferAlreadyExists = fmt.Errorf("offer %w", ErrAlreadyExists)
	ErrOfferNotPublished  = fmt.Errorf("offer not published: %w", ErrForbidden)
	ErrOfferExpired       = fmt.Errorf("offer %w", ErrExpired)
	ErrOfferFullyBooked   = fmt.Errorf("offer fully booked: %w", ErrQuotaExceeded)
	ErrCategoryNotFound   = fmt.Errorf("category %w", ErrNotFound)
	ErrInvalidSchedule    = fmt.Errorf("invalid schedule: %w", ErrValidation)
	ErrInvalidDiscount    = fmt.Errorf("invalid discount: %w", ErrValidation)
)

// Booking Context Errors
var (
	ErrOutingNotFound        = fmt.Errorf("outing %w", ErrNotFound)
	ErrOutingAlreadyExists   = fmt.Errorf("outing %w", ErrAlreadyExists)
	ErrOutingExpired         = fmt.Errorf("outing %w", ErrExpired)
	ErrOutingAlreadyUsed     = fmt.Errorf("outing already used: %w", ErrConflict)
	ErrOutingCancelled       = fmt.Errorf("outing cancelled: %w", ErrConflict)
	ErrInvalidQRCode         = fmt.Errorf("invalid QR code: %w", ErrValidation)
	ErrCannotCancelCheckedIn = fmt.Errorf("cannot cancel checked-in outing: %w", ErrConflict)
	ErrBookingLimitReached   = fmt.Errorf("booking limit reached: %w", ErrQuotaExceeded)
	ErrInvalidOutingStatus   = fmt.Errorf("invalid outing status: %w", ErrConflict)
)

// Engagement Context Errors
var (
	ErrReviewNotFound        = fmt.Errorf("review %w", ErrNotFound)
	ErrReviewAlreadyExists   = fmt.Errorf("review %w", ErrAlreadyExists)
	ErrCannotReviewOwnOffer  = fmt.Errorf("cannot review own offer: %w", ErrForbidden)
	ErrFavoriteNotFound      = fmt.Errorf("favorite %w", ErrNotFound)
	ErrFavoriteAlreadyExists = fmt.Errorf("favorite %w", ErrAlreadyExists)
	ErrConversationNotFound  = fmt.Errorf("conversation %w", ErrNotFound)
	ErrNotConversationMember = fmt.Errorf("not conversation member: %w", ErrForbidden)
	ErrInvalidRating         = fmt.Errorf("invalid rating (1-5): %w", ErrValidation)
)

// Notification Context Errors
var (
	ErrNotificationNotFound = fmt.Errorf("notification %w", ErrNotFound)
	ErrPushTokenInvalid     = fmt.Errorf("invalid push token: %w", ErrValidation)
	ErrNotificationFailed   = fmt.Errorf("notification failed: %w", ErrInternal)
)

// =============================================================================
// Error Wrapper
// =============================================================================

// DomainError wraps an error with additional context.
type DomainError struct {
	Err     error
	Message string
	Code    string
	Details map[string]any
}

// NewDomainError creates a new domain error.
func NewDomainError(err error, message string) *DomainError {
	return &DomainError{
		Err:     err,
		Message: message,
		Details: make(map[string]any),
	}
}

// WithCode adds an error code.
func (e *DomainError) WithCode(code string) *DomainError {
	e.Code = code
	return e
}

// WithDetail adds a detail to the error.
func (e *DomainError) WithDetail(key string, value any) *DomainError {
	e.Details[key] = value
	return e
}

// Error implements the error interface.
func (e *DomainError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Err.Error()
}

// Unwrap returns the underlying error.
func (e *DomainError) Unwrap() error {
	return e.Err
}

// Is checks if the error matches the target.
func (e *DomainError) Is(target error) bool {
	return errors.Is(e.Err, target)
}

// =============================================================================
// Error Helpers
// =============================================================================

// IsNotFound checks if the error is a not found error.
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// IsAlreadyExists checks if the error is an already exists error.
func IsAlreadyExists(err error) bool {
	return errors.Is(err, ErrAlreadyExists)
}

// IsUnauthorized checks if the error is an unauthorized error.
func IsUnauthorized(err error) bool {
	return errors.Is(err, ErrUnauthorized)
}

// IsForbidden checks if the error is a forbidden error.
func IsForbidden(err error) bool {
	return errors.Is(err, ErrForbidden)
}

// IsValidation checks if the error is a validation error.
func IsValidation(err error) bool {
	return errors.Is(err, ErrValidation) || errors.Is(err, ErrInvalidInput)
}

// IsConflict checks if the error is a conflict error.
func IsConflict(err error) bool {
	return errors.Is(err, ErrConflict)
}

// IsExpired checks if the error is an expiration error.
func IsExpired(err error) bool {
	return errors.Is(err, ErrExpired)
}

// IsQuotaExceeded checks if the error is a quota exceeded error.
func IsQuotaExceeded(err error) bool {
	return errors.Is(err, ErrQuotaExceeded)
}

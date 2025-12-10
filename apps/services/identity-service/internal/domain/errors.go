package domain

import (
	"errors"
)

// Domain errors for Identity context.
var (
	// User errors
	ErrUserNotFound        = errors.New("user not found")
	ErrEmailAlreadyExists  = errors.New("email already exists")
	ErrPhoneAlreadyExists  = errors.New("phone already exists")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrUserNotActive       = errors.New("user account is not active")
	ErrUserSuspended       = errors.New("user account is suspended")
	ErrUserDeleted         = errors.New("user account is deleted")
	ErrEmailNotVerified    = errors.New("email is not verified")
	ErrPhoneNotVerified    = errors.New("phone is not verified")
	ErrIdentityNotVerified = errors.New("identity is not verified")
	ErrAlreadyVerified     = errors.New("already verified")
	ErrMaxAttemptsExceeded = errors.New("maximum verification attempts exceeded")

	// Password errors
	ErrPasswordTooShort   = errors.New("password must be at least 8 characters")
	ErrPasswordTooWeak    = errors.New("password is too weak")
	ErrPasswordMismatch   = errors.New("passwords do not match")
	ErrInvalidOldPassword = errors.New("invalid old password")

	// Token errors
	ErrInvalidToken         = errors.New("invalid token")
	ErrTokenExpired         = errors.New("token has expired")
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
	ErrRefreshTokenRevoked  = errors.New("refresh token has been revoked")

	// Subscription errors
	ErrSubscriptionNotFound = errors.New("subscription not found")
	ErrAlreadySubscribed    = errors.New("user already has an active subscription")
	ErrSubscriptionExpired  = errors.New("subscription has expired")
	ErrInvalidPlan          = errors.New("invalid subscription plan")
	ErrPlanNotFound         = errors.New("subscription plan not found")

	// Social account errors
	ErrSocialAccountNotFound   = errors.New("social account not found")
	ErrSocialAccountExists     = errors.New("social account already linked")
	ErrCannotUnlinkLastAccount = errors.New("cannot unlink last authentication method")

	// Verification errors
	ErrVerificationNotFound = errors.New("verification not found")
	ErrVerificationPending  = errors.New("verification is pending")
	ErrVerificationRejected = errors.New("verification was rejected")
)

// IsNotFoundError returns true if the error is a "not found" type error.
func IsNotFoundError(err error) bool {
	return errors.Is(err, ErrUserNotFound) ||
		errors.Is(err, ErrSubscriptionNotFound) ||
		errors.Is(err, ErrPlanNotFound) ||
		errors.Is(err, ErrSocialAccountNotFound) ||
		errors.Is(err, ErrVerificationNotFound) ||
		errors.Is(err, ErrRefreshTokenNotFound)
}

// IsAuthenticationError returns true if the error is an authentication error.
func IsAuthenticationError(err error) bool {
	return errors.Is(err, ErrInvalidCredentials) ||
		errors.Is(err, ErrInvalidToken) ||
		errors.Is(err, ErrTokenExpired) ||
		errors.Is(err, ErrRefreshTokenRevoked)
}

// IsValidationError returns true if the error is a validation error.
func IsValidationError(err error) bool {
	return errors.Is(err, ErrPasswordTooShort) ||
		errors.Is(err, ErrPasswordTooWeak) ||
		errors.Is(err, ErrPasswordMismatch) ||
		errors.Is(err, ErrEmailAlreadyExists) ||
		errors.Is(err, ErrPhoneAlreadyExists)
}

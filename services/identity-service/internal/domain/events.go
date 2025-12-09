package domain

import (
	"time"

	"github.com/yousoon/services/shared/domain"
)

// =============================================================================
// User Events
// =============================================================================

// UserRegisteredEvent is raised when a new user registers.
type UserRegisteredEvent struct {
	domain.BaseEvent
	UserID UserID `json:"userId"`
	Email  Email  `json:"email"`
}

// NewUserRegisteredEvent creates a new UserRegisteredEvent.
func NewUserRegisteredEvent(userID UserID, email Email) UserRegisteredEvent {
	return UserRegisteredEvent{
		BaseEvent: domain.NewBaseEvent(
			"user.registered",
			userID.String(),
			"User",
			map[string]any{
				"email": email.String(),
			},
		),
		UserID: userID,
		Email:  email,
	}
}

// UserEmailVerifiedEvent is raised when a user's email is verified.
type UserEmailVerifiedEvent struct {
	domain.BaseEvent
	UserID UserID `json:"userId"`
}

// NewUserEmailVerifiedEvent creates a new UserEmailVerifiedEvent.
func NewUserEmailVerifiedEvent(userID UserID) UserEmailVerifiedEvent {
	return UserEmailVerifiedEvent{
		BaseEvent: domain.NewBaseEvent(
			"user.email_verified",
			userID.String(),
			"User",
			nil,
		),
		UserID: userID,
	}
}

// UserProfileUpdatedEvent is raised when a user updates their profile.
type UserProfileUpdatedEvent struct {
	domain.BaseEvent
	UserID UserID `json:"userId"`
}

// NewUserProfileUpdatedEvent creates a new UserProfileUpdatedEvent.
func NewUserProfileUpdatedEvent(userID UserID) UserProfileUpdatedEvent {
	return UserProfileUpdatedEvent{
		BaseEvent: domain.NewBaseEvent(
			"user.profile_updated",
			userID.String(),
			"User",
			nil,
		),
		UserID: userID,
	}
}

// IdentityVerificationSubmittedEvent is raised when identity verification is submitted.
type IdentityVerificationSubmittedEvent struct {
	domain.BaseEvent
	UserID         UserID         `json:"userId"`
	VerificationID VerificationID `json:"verificationId"`
}

// NewIdentityVerificationSubmittedEvent creates a new IdentityVerificationSubmittedEvent.
func NewIdentityVerificationSubmittedEvent(userID UserID, verificationID VerificationID) IdentityVerificationSubmittedEvent {
	return IdentityVerificationSubmittedEvent{
		BaseEvent: domain.NewBaseEvent(
			"identity.verification_submitted",
			userID.String(),
			"User",
			map[string]any{
				"verificationId": verificationID.String(),
			},
		),
		UserID:         userID,
		VerificationID: verificationID,
	}
}

// UserIdentityVerifiedEvent is raised when a user's identity is verified.
type UserIdentityVerifiedEvent struct {
	domain.BaseEvent
	UserID UserID `json:"userId"`
}

// NewUserIdentityVerifiedEvent creates a new UserIdentityVerifiedEvent.
func NewUserIdentityVerifiedEvent(userID UserID) UserIdentityVerifiedEvent {
	return UserIdentityVerifiedEvent{
		BaseEvent: domain.NewBaseEvent(
			"user.identity_verified",
			userID.String(),
			"User",
			nil,
		),
		UserID: userID,
	}
}

// UserSubscribedEvent is raised when a user subscribes.
type UserSubscribedEvent struct {
	domain.BaseEvent
	UserID         UserID         `json:"userId"`
	SubscriptionID SubscriptionID `json:"subscriptionId"`
}

// NewUserSubscribedEvent creates a new UserSubscribedEvent.
func NewUserSubscribedEvent(userID UserID, subscriptionID SubscriptionID) UserSubscribedEvent {
	return UserSubscribedEvent{
		BaseEvent: domain.NewBaseEvent(
			"user.subscribed",
			userID.String(),
			"User",
			map[string]any{
				"subscriptionId": subscriptionID.String(),
			},
		),
		UserID:         userID,
		SubscriptionID: subscriptionID,
	}
}

// UserSubscriptionCancelledEvent is raised when a user cancels their subscription.
type UserSubscriptionCancelledEvent struct {
	domain.BaseEvent
	UserID         UserID         `json:"userId"`
	SubscriptionID SubscriptionID `json:"subscriptionId"`
}

// NewUserSubscriptionCancelledEvent creates a new UserSubscriptionCancelledEvent.
func NewUserSubscriptionCancelledEvent(userID UserID, subscriptionID SubscriptionID) UserSubscriptionCancelledEvent {
	return UserSubscriptionCancelledEvent{
		BaseEvent: domain.NewBaseEvent(
			"user.subscription_cancelled",
			userID.String(),
			"User",
			map[string]any{
				"subscriptionId": subscriptionID.String(),
			},
		),
		UserID:         userID,
		SubscriptionID: subscriptionID,
	}
}

// UserDeletedEvent is raised when a user is deleted (soft delete).
type UserDeletedEvent struct {
	domain.BaseEvent
	UserID UserID `json:"userId"`
	Reason string `json:"reason,omitempty"`
}

// NewUserDeletedEvent creates a new UserDeletedEvent.
func NewUserDeletedEvent(userID UserID, reason string) UserDeletedEvent {
	return UserDeletedEvent{
		BaseEvent: domain.NewBaseEvent(
			"user.deleted",
			userID.String(),
			"User",
			map[string]any{
				"reason": reason,
			},
		),
		UserID: userID,
		Reason: reason,
	}
}

// UserLoggedInEvent is raised when a user logs in.
type UserLoggedInEvent struct {
	domain.BaseEvent
	UserID    UserID    `json:"userId"`
	Platform  Platform  `json:"platform"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"userAgent"`
	Timestamp time.Time `json:"timestamp"`
}

// NewUserLoggedInEvent creates a new UserLoggedInEvent.
func NewUserLoggedInEvent(userID UserID, platform Platform, ip, userAgent string) UserLoggedInEvent {
	now := time.Now()
	return UserLoggedInEvent{
		BaseEvent: domain.NewBaseEvent(
			"user.logged_in",
			userID.String(),
			"User",
			map[string]any{
				"platform":  string(platform),
				"ip":        ip,
				"userAgent": userAgent,
			},
		),
		UserID:    userID,
		Platform:  platform,
		IP:        ip,
		UserAgent: userAgent,
		Timestamp: now,
	}
}

// PasswordResetRequestedEvent is raised when a password reset is requested.
type PasswordResetRequestedEvent struct {
	domain.BaseEvent
	UserID UserID `json:"userId"`
	Email  Email  `json:"email"`
}

// NewPasswordResetRequestedEvent creates a new PasswordResetRequestedEvent.
func NewPasswordResetRequestedEvent(userID UserID, email Email) PasswordResetRequestedEvent {
	return PasswordResetRequestedEvent{
		BaseEvent: domain.NewBaseEvent(
			"user.password_reset_requested",
			userID.String(),
			"User",
			map[string]any{
				"email": email.String(),
			},
		),
		UserID: userID,
		Email:  email,
	}
}

// PasswordChangedEvent is raised when a user changes their password.
type PasswordChangedEvent struct {
	domain.BaseEvent
	UserID UserID `json:"userId"`
}

// NewPasswordChangedEvent creates a new PasswordChangedEvent.
func NewPasswordChangedEvent(userID UserID) PasswordChangedEvent {
	return PasswordChangedEvent{
		BaseEvent: domain.NewBaseEvent(
			"user.password_changed",
			userID.String(),
			"User",
			nil,
		),
		UserID: userID,
	}
}

// UserSuspendedEvent is raised when a user is suspended.
type UserSuspendedEvent struct {
	domain.BaseEvent
	UserID UserID `json:"userId"`
	Reason string `json:"reason"`
}

// NewUserSuspendedEvent creates a new UserSuspendedEvent.
func NewUserSuspendedEvent(userID UserID, reason string) UserSuspendedEvent {
	return UserSuspendedEvent{
		BaseEvent: domain.NewBaseEvent(
			"user.suspended",
			userID.String(),
			"User",
			map[string]any{
				"reason": reason,
			},
		),
		UserID: userID,
		Reason: reason,
	}
}

// UserReactivatedEvent is raised when a user is reactivated.
type UserReactivatedEvent struct {
	domain.BaseEvent
	UserID UserID `json:"userId"`
}

// NewUserReactivatedEvent creates a new UserReactivatedEvent.
func NewUserReactivatedEvent(userID UserID) UserReactivatedEvent {
	return UserReactivatedEvent{
		BaseEvent: domain.NewBaseEvent(
			"user.reactivated",
			userID.String(),
			"User",
			nil,
		),
		UserID: userID,
	}
}

// UserGradeUpgradedEvent is raised when a user's grade is upgraded.
type UserGradeUpgradedEvent struct {
	domain.BaseEvent
	UserID   UserID    `json:"userId"`
	OldGrade UserGrade `json:"oldGrade"`
	NewGrade UserGrade `json:"newGrade"`
}

// NewUserGradeUpgradedEvent creates a new UserGradeUpgradedEvent.
func NewUserGradeUpgradedEvent(userID UserID, oldGrade, newGrade UserGrade) UserGradeUpgradedEvent {
	return UserGradeUpgradedEvent{
		BaseEvent: domain.NewBaseEvent(
			"user.grade_upgraded",
			userID.String(),
			"User",
			map[string]any{
				"oldGrade": oldGrade.String(),
				"newGrade": newGrade.String(),
			},
		),
		UserID:   userID,
		OldGrade: oldGrade,
		NewGrade: newGrade,
	}
}

// UserLocationUpdatedEvent is raised when a user's location is updated.
type UserLocationUpdatedEvent struct {
	domain.BaseEvent
	UserID   UserID      `json:"userId"`
	Location GeoLocation `json:"location"`
}

// NewUserLocationUpdatedEvent creates a new UserLocationUpdatedEvent.
func NewUserLocationUpdatedEvent(userID UserID, location GeoLocation) UserLocationUpdatedEvent {
	return UserLocationUpdatedEvent{
		BaseEvent: domain.NewBaseEvent(
			"user.location_updated",
			userID.String(),
			"User",
			map[string]any{
				"longitude": location.Longitude(),
				"latitude":  location.Latitude(),
			},
		),
		UserID:   userID,
		Location: location,
	}
}

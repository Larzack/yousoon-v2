// Package domain contains the domain layer for the Identity bounded context.
// This includes User aggregate, value objects, domain events, and repository interfaces.
package domain

import (
	"time"

	"github.com/yousoon/services/shared/domain"
)

// =============================================================================
// User Aggregate Root
// =============================================================================

// User is the aggregate root for user identity management.
type User struct {
	domain.VersionedAggregateRoot

	// Identity
	ID           UserID `json:"id" bson:"_id"`
	Email        Email  `json:"email" bson:"email"`
	PasswordHash string `json:"-" bson:"passwordHash"`
	Phone        *Phone `json:"phone,omitempty" bson:"phone,omitempty"`

	// Profile
	Profile Profile `json:"profile" bson:"profile"`

	// Identity Verification
	Identity *IdentityVerification `json:"identity,omitempty" bson:"identity,omitempty"`

	// Subscription
	SubscriptionID *SubscriptionID `json:"subscriptionId,omitempty" bson:"subscriptionId,omitempty"`

	// Preferences
	Preferences Preferences `json:"preferences" bson:"preferences"`

	// Location
	LastLocation *GeoLocation `json:"lastLocation,omitempty" bson:"lastLocation,omitempty"`

	// Grade
	Grade UserGrade `json:"grade" bson:"grade"`

	// Social Accounts
	SocialAccounts []SocialAccount `json:"socialAccounts,omitempty" bson:"socialAccounts,omitempty"`

	// FCM Tokens for push notifications
	FCMTokens []FCMToken `json:"fcmTokens,omitempty" bson:"fcmTokens,omitempty"`

	// Status
	Status        UserStatus `json:"status" bson:"status"`
	EmailVerified bool       `json:"emailVerified" bson:"emailVerified"`
	PhoneVerified bool       `json:"phoneVerified" bson:"phoneVerified"`

	// Timestamps
	LastLoginAt *time.Time `json:"lastLoginAt,omitempty" bson:"lastLoginAt,omitempty"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty" bson:"deletedAt,omitempty"`
}

// NewUser creates a new user with required fields.
func NewUser(email Email, passwordHash string, profile Profile) (*User, error) {
	user := &User{
		VersionedAggregateRoot: domain.NewVersionedAggregateRoot(),
		ID:                     NewUserID(),
		Email:                  email,
		PasswordHash:           passwordHash,
		Profile:                profile,
		Preferences:            DefaultPreferences(),
		Grade:                  GradeExplorateur,
		Status:                 UserStatusActive,
		EmailVerified:          false,
		PhoneVerified:          false,
	}

	user.AddDomainEvent(NewUserRegisteredEvent(user.ID, email))
	return user, nil
}

// =============================================================================
// User Methods (Business Logic)
// =============================================================================

// CanBook checks if the user can make a booking.
func (u *User) CanBook() error {
	if u.Status != UserStatusActive {
		return domain.ErrAccountSuspended
	}
	if u.Identity == nil || u.Identity.Status != VerificationStatusVerified {
		return domain.ErrIdentityNotVerified
	}
	return nil
}

// VerifyEmail marks the email as verified.
func (u *User) VerifyEmail() {
	u.EmailVerified = true
	u.MarkUpdated()
	u.AddDomainEvent(NewUserEmailVerifiedEvent(u.ID))
}

// VerifyPhone marks the phone as verified.
func (u *User) VerifyPhone() {
	u.PhoneVerified = true
	u.MarkUpdated()
}

// SetPhone sets the user's phone number.
func (u *User) SetPhone(phone Phone) {
	u.Phone = &phone
	u.PhoneVerified = false
	u.MarkUpdated()
}

// UpdateProfile updates the user's profile.
func (u *User) UpdateProfile(profile Profile) {
	u.Profile = profile
	u.MarkUpdated()
	u.AddDomainEvent(NewUserProfileUpdatedEvent(u.ID))
}

// UpdatePreferences updates the user's preferences.
func (u *User) UpdatePreferences(prefs Preferences) {
	u.Preferences = prefs
	u.MarkUpdated()
}

// UpdateLocation updates the user's last known location.
func (u *User) UpdateLocation(location GeoLocation) {
	u.LastLocation = &location
	u.MarkUpdated()
}

// RecordLogin records a login event.
func (u *User) RecordLogin() {
	now := time.Now()
	u.LastLoginAt = &now
	u.MarkUpdated()
}

// SubmitIdentityVerification submits identity verification.
func (u *User) SubmitIdentityVerification(verification IdentityVerification) error {
	if u.Identity != nil && u.Identity.Status == VerificationStatusVerified {
		return domain.ErrAlreadyVerified
	}
	u.Identity = &verification
	u.MarkUpdated()
	u.AddDomainEvent(NewIdentityVerificationSubmittedEvent(u.ID, verification.ID))
	return nil
}

// ApproveIdentityVerification approves the identity verification.
func (u *User) ApproveIdentityVerification() error {
	if u.Identity == nil {
		return domain.ErrVerificationFailed
	}
	now := time.Now()
	u.Identity.Status = VerificationStatusVerified
	u.Identity.VerifiedAt = &now
	u.MarkUpdated()
	u.AddDomainEvent(NewUserIdentityVerifiedEvent(u.ID))
	return nil
}

// RejectIdentityVerification rejects the identity verification.
func (u *User) RejectIdentityVerification(reason string) error {
	if u.Identity == nil {
		return domain.ErrVerificationFailed
	}
	now := time.Now()
	u.Identity.Status = VerificationStatusRejected
	u.Identity.RejectedAt = &now
	u.Identity.RejectionReason = &reason
	u.MarkUpdated()
	return nil
}

// SetSubscription sets the user's subscription.
func (u *User) SetSubscription(subID SubscriptionID) {
	u.SubscriptionID = &subID
	u.MarkUpdated()
	u.AddDomainEvent(NewUserSubscribedEvent(u.ID, subID))
}

// ClearSubscription removes the user's subscription.
func (u *User) ClearSubscription() {
	if u.SubscriptionID != nil {
		subID := *u.SubscriptionID
		u.SubscriptionID = nil
		u.MarkUpdated()
		u.AddDomainEvent(NewUserSubscriptionCancelledEvent(u.ID, subID))
	}
}

// Suspend suspends the user account.
func (u *User) Suspend(reason string) {
	u.Status = UserStatusSuspended
	u.MarkUpdated()
}

// Activate activates the user account.
func (u *User) Activate() {
	u.Status = UserStatusActive
	u.MarkUpdated()
}

// SoftDelete marks the user for deletion.
func (u *User) SoftDelete(reason string) {
	now := time.Now()
	u.Status = UserStatusDeleted
	u.DeletedAt = &now
	u.MarkUpdated()
	u.AddDomainEvent(NewUserDeletedEvent(u.ID, reason))
}

// AddFCMToken adds an FCM token for push notifications.
func (u *User) AddFCMToken(token FCMToken) {
	// Remove existing token for the same platform if any
	tokens := make([]FCMToken, 0, len(u.FCMTokens))
	for _, t := range u.FCMTokens {
		if t.Token != token.Token {
			tokens = append(tokens, t)
		}
	}
	tokens = append(tokens, token)
	u.FCMTokens = tokens
	u.MarkUpdated()
}

// RemoveFCMToken removes an FCM token.
func (u *User) RemoveFCMToken(token string) {
	tokens := make([]FCMToken, 0, len(u.FCMTokens))
	for _, t := range u.FCMTokens {
		if t.Token != token {
			tokens = append(tokens, t)
		}
	}
	u.FCMTokens = tokens
	u.MarkUpdated()
}

// AddSocialAccount links a social account.
func (u *User) AddSocialAccount(account SocialAccount) error {
	for _, a := range u.SocialAccounts {
		if a.Provider == account.Provider && a.ProviderID == account.ProviderID {
			return domain.ErrAlreadyExists
		}
	}
	u.SocialAccounts = append(u.SocialAccounts, account)
	u.MarkUpdated()
	return nil
}

// UpgradeGrade upgrades the user's grade.
func (u *User) UpgradeGrade(newGrade UserGrade) {
	if newGrade > u.Grade {
		u.Grade = newGrade
		u.MarkUpdated()
	}
}

// HasVerifiedIdentity returns true if identity is verified.
func (u *User) HasVerifiedIdentity() bool {
	return u.Identity != nil && u.Identity.Status == VerificationStatusVerified
}

// HasActiveSubscription returns true if user has an active subscription.
func (u *User) HasActiveSubscription() bool {
	return u.SubscriptionID != nil
}

// IsActive returns true if the user account is active.
func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}

// FullName returns the user's full name.
func (u *User) FullName() string {
	return u.Profile.FullName()
}

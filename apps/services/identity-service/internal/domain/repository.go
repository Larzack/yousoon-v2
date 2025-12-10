package domain

import (
	"context"
)

// =============================================================================
// User Repository Interface
// =============================================================================

// UserRepository defines the interface for user persistence.
type UserRepository interface {
	// Create creates a new user.
	Create(ctx context.Context, user *User) error

	// Update updates an existing user.
	Update(ctx context.Context, user *User) error

	// Delete permanently deletes a user.
	Delete(ctx context.Context, id UserID) error

	// FindByID finds a user by ID.
	FindByID(ctx context.Context, id UserID) (*User, error)

	// FindByEmail finds a user by email.
	FindByEmail(ctx context.Context, email Email) (*User, error)

	// FindByPhone finds a user by phone.
	FindByPhone(ctx context.Context, phone Phone) (*User, error)

	// FindBySocialAccount finds a user by social account.
	FindBySocialAccount(ctx context.Context, provider SocialProvider, providerID string) (*User, error)

	// ExistsByEmail checks if a user with the email exists.
	ExistsByEmail(ctx context.Context, email Email) (bool, error)

	// ExistsByPhone checks if a user with the phone exists.
	ExistsByPhone(ctx context.Context, phone Phone) (bool, error)

	// FindUsersForDeletion finds users scheduled for deletion.
	FindUsersForDeletion(ctx context.Context, before int64) ([]*User, error)

	// Count returns the total number of users.
	Count(ctx context.Context) (int64, error)

	// CountActive returns the number of active users.
	CountActive(ctx context.Context) (int64, error)
}

// UserQuery represents query options for finding users.
type UserQuery struct {
	Status        *UserStatus
	EmailVerified *bool
	HasIdentity   *bool
	Limit         int
	Offset        int
	SortBy        string
	SortDesc      bool
}

// =============================================================================
// Subscription Repository Interface
// =============================================================================

// SubscriptionRepository defines the interface for subscription persistence.
type SubscriptionRepository interface {
	// Create creates a new subscription.
	Create(ctx context.Context, subscription *Subscription) error

	// Update updates an existing subscription.
	Update(ctx context.Context, subscription *Subscription) error

	// FindByID finds a subscription by ID.
	FindByID(ctx context.Context, id SubscriptionID) (*Subscription, error)

	// FindByUserID finds a subscription by user ID.
	FindByUserID(ctx context.Context, userID UserID) (*Subscription, error)

	// FindByTransactionID finds a subscription by transaction ID.
	FindByTransactionID(ctx context.Context, transactionID string) (*Subscription, error)

	// FindActive finds all active subscriptions.
	FindActive(ctx context.Context) ([]*Subscription, error)

	// FindExpiring finds subscriptions expiring before a date.
	FindExpiring(ctx context.Context, before int64) ([]*Subscription, error)
}

// =============================================================================
// Subscription Plan Repository Interface
// =============================================================================

// SubscriptionPlanRepository defines the interface for subscription plan persistence.
type SubscriptionPlanRepository interface {
	// Create creates a new subscription plan.
	Create(ctx context.Context, plan *SubscriptionPlan) error

	// Update updates an existing subscription plan.
	Update(ctx context.Context, plan *SubscriptionPlan) error

	// FindByID finds a subscription plan by ID.
	FindByID(ctx context.Context, id PlanID) (*SubscriptionPlan, error)

	// FindByCode finds a subscription plan by code.
	FindByCode(ctx context.Context, code string) (*SubscriptionPlan, error)

	// FindActive finds all active subscription plans.
	FindActive(ctx context.Context) ([]*SubscriptionPlan, error)
}

// =============================================================================
// Token Repository Interface (Redis-based)
// =============================================================================

// TokenRepository defines the interface for token storage.
type TokenRepository interface {
	// StoreRefreshToken stores a refresh token.
	StoreRefreshToken(ctx context.Context, userID UserID, token string, ttl int64) error

	// GetRefreshToken gets a refresh token.
	GetRefreshToken(ctx context.Context, userID UserID) (string, error)

	// DeleteRefreshToken deletes a refresh token.
	DeleteRefreshToken(ctx context.Context, userID UserID) error

	// StoreEmailVerificationToken stores an email verification token.
	StoreEmailVerificationToken(ctx context.Context, token string, userID UserID, ttl int64) error

	// GetEmailVerificationToken gets the user ID for a verification token.
	GetEmailVerificationToken(ctx context.Context, token string) (UserID, error)

	// DeleteEmailVerificationToken deletes an email verification token.
	DeleteEmailVerificationToken(ctx context.Context, token string) error

	// StorePasswordResetToken stores a password reset token.
	StorePasswordResetToken(ctx context.Context, token string, userID UserID, ttl int64) error

	// GetPasswordResetToken gets the user ID for a reset token.
	GetPasswordResetToken(ctx context.Context, token string) (UserID, error)

	// DeletePasswordResetToken deletes a password reset token.
	DeletePasswordResetToken(ctx context.Context, token string) error

	// BlacklistAccessToken blacklists an access token.
	BlacklistAccessToken(ctx context.Context, token string, ttl int64) error

	// IsAccessTokenBlacklisted checks if an access token is blacklisted.
	IsAccessTokenBlacklisted(ctx context.Context, token string) (bool, error)
}

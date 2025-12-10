package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/yousoon/shared/domain"
)

// =============================================================================
// Plan ID
// =============================================================================

// PlanID is a unique identifier for a subscription plan.
type PlanID string

// NewPlanID creates a new PlanID.
func NewPlanID() PlanID {
	return PlanID(uuid.New().String())
}

// String returns the string representation.
func (id PlanID) String() string {
	return string(id)
}

// =============================================================================
// Subscription Plan
// =============================================================================

// SubscriptionPlan represents a subscription plan.
type SubscriptionPlan struct {
	domain.BaseEntity

	ID          PlanID  `json:"id" bson:"_id"`
	Code        string  `json:"code" bson:"code"` // free, monthly, yearly
	Name        I18nMap `json:"name" bson:"name"`
	Description I18nMap `json:"description" bson:"description"`

	// Pricing
	Pricing PlanPricing `json:"pricing" bson:"pricing"`

	// Trial
	Trial PlanTrial `json:"trial" bson:"trial"`

	// Features
	Features []PlanFeature `json:"features" bson:"features"`

	// Limits
	Limits PlanLimits `json:"limits" bson:"limits"`

	// Display
	Display PlanDisplay `json:"display" bson:"display"`

	// Store IDs
	AppleProductID  string `json:"appleProductId,omitempty" bson:"appleProductId,omitempty"`
	GoogleProductID string `json:"googleProductId,omitempty" bson:"googleProductId,omitempty"`

	// Status
	IsActive bool `json:"isActive" bson:"isActive"`
}

// I18nMap represents internationalized text.
type I18nMap struct {
	FR string `json:"fr" bson:"fr"`
	EN string `json:"en" bson:"en"`
}

// PlanPricing represents plan pricing.
type PlanPricing struct {
	Amount        int64  `json:"amount" bson:"amount"`               // In cents
	Currency      string `json:"currency" bson:"currency"`           // EUR
	Interval      string `json:"interval" bson:"interval"`           // month, year
	IntervalCount int    `json:"intervalCount" bson:"intervalCount"` // 1, 3, 12
}

// PlanTrial represents trial configuration.
type PlanTrial struct {
	Enabled      bool `json:"enabled" bson:"enabled"`
	DurationDays int  `json:"durationDays" bson:"durationDays"`
}

// PlanFeature represents a plan feature.
type PlanFeature struct {
	Code     string  `json:"code" bson:"code"`
	Name     I18nMap `json:"name" bson:"name"`
	Included bool    `json:"included" bson:"included"`
	Limit    *int    `json:"limit,omitempty" bson:"limit,omitempty"`
}

// PlanLimits represents plan limits.
type PlanLimits struct {
	BookingsPerMonth *int `json:"bookingsPerMonth,omitempty" bson:"bookingsPerMonth,omitempty"`
	FavoritesMax     *int `json:"favoritesMax,omitempty" bson:"favoritesMax,omitempty"`
}

// PlanDisplay represents display configuration.
type PlanDisplay struct {
	Order       int     `json:"order" bson:"order"`
	Highlighted bool    `json:"highlighted" bson:"highlighted"`
	Badge       *string `json:"badge,omitempty" bson:"badge,omitempty"`
	Color       *string `json:"color,omitempty" bson:"color,omitempty"`
}

// =============================================================================
// Subscription
// =============================================================================

// Subscription represents a user's subscription.
type Subscription struct {
	domain.BaseEntity

	ID     SubscriptionID `json:"id" bson:"_id"`
	UserID UserID         `json:"userId" bson:"userId"`
	PlanID PlanID         `json:"planId" bson:"planId"`

	// In-App Purchase
	InAppPurchase *InAppPurchase `json:"inAppPurchase,omitempty" bson:"inAppPurchase,omitempty"`

	// Status
	Status SubscriptionStatus `json:"status" bson:"status"`

	// Trial
	Trial *SubscriptionTrial `json:"trial,omitempty" bson:"trial,omitempty"`

	// Current Period
	CurrentPeriod SubscriptionPeriod `json:"currentPeriod" bson:"currentPeriod"`

	// Cancellation
	Cancellation *SubscriptionCancellation `json:"cancellation,omitempty" bson:"cancellation,omitempty"`

	// Plan snapshot
	PlanSnapshot *SubscriptionPlanSnapshot `json:"_plan,omitempty" bson:"_plan,omitempty"`
}

// InAppPurchase represents in-app purchase details.
type InAppPurchase struct {
	Platform           Platform   `json:"platform" bson:"platform"`
	ProductID          string     `json:"productId" bson:"productId"`
	TransactionID      string     `json:"transactionId" bson:"transactionId"`
	OriginalTransID    string     `json:"originalTransactionId,omitempty" bson:"originalTransactionId,omitempty"`
	Receipt            string     `json:"receipt,omitempty" bson:"receipt,omitempty"`
	ReceiptValidatedAt *time.Time `json:"receiptValidatedAt,omitempty" bson:"receiptValidatedAt,omitempty"`
}

// SubscriptionStatus represents subscription status.
type SubscriptionStatus string

const (
	SubscriptionStatusTrialing  SubscriptionStatus = "trialing"
	SubscriptionStatusActive    SubscriptionStatus = "active"
	SubscriptionStatusPastDue   SubscriptionStatus = "past_due"
	SubscriptionStatusCancelled SubscriptionStatus = "cancelled"
	SubscriptionStatusExpired   SubscriptionStatus = "expired"
)

// SubscriptionTrial represents trial information.
type SubscriptionTrial struct {
	StartDate time.Time `json:"startDate" bson:"startDate"`
	EndDate   time.Time `json:"endDate" bson:"endDate"`
	Converted bool      `json:"converted" bson:"converted"`
}

// SubscriptionPeriod represents a billing period.
type SubscriptionPeriod struct {
	StartDate time.Time `json:"startDate" bson:"startDate"`
	EndDate   time.Time `json:"endDate" bson:"endDate"`
}

// SubscriptionCancellation represents cancellation details.
type SubscriptionCancellation struct {
	RequestedAt time.Time  `json:"requestedAt" bson:"requestedAt"`
	Reason      *string    `json:"reason,omitempty" bson:"reason,omitempty"`
	EffectiveAt *time.Time `json:"effectiveAt,omitempty" bson:"effectiveAt,omitempty"`
	Feedback    *string    `json:"feedback,omitempty" bson:"feedback,omitempty"`
}

// SubscriptionPlanSnapshot represents plan snapshot at subscription time.
type SubscriptionPlanSnapshot struct {
	Code     string `json:"code" bson:"code"`
	Name     string `json:"name" bson:"name"`
	Amount   int64  `json:"amount" bson:"amount"`
	Interval string `json:"interval" bson:"interval"`
}

// NewSubscription creates a new subscription.
func NewSubscription(userID UserID, planID PlanID, status SubscriptionStatus) *Subscription {
	now := time.Now()
	return &Subscription{
		BaseEntity: domain.NewBaseEntity(),
		ID:         NewSubscriptionID(),
		UserID:     userID,
		PlanID:     planID,
		Status:     status,
		CurrentPeriod: SubscriptionPeriod{
			StartDate: now,
			EndDate:   now.AddDate(0, 1, 0), // Default 1 month
		},
	}
}

// IsActive returns true if the subscription is active.
func (s *Subscription) IsActive() bool {
	return s.Status == SubscriptionStatusActive || s.Status == SubscriptionStatusTrialing
}

// IsTrialing returns true if the subscription is in trial.
func (s *Subscription) IsTrialing() bool {
	return s.Status == SubscriptionStatusTrialing
}

// Cancel cancels the subscription.
func (s *Subscription) Cancel(reason string) {
	now := time.Now()
	s.Cancellation = &SubscriptionCancellation{
		RequestedAt: now,
		Reason:      &reason,
		EffectiveAt: &s.CurrentPeriod.EndDate,
	}
	s.MarkUpdated()
}

// Expire expires the subscription.
func (s *Subscription) Expire() {
	s.Status = SubscriptionStatusExpired
	s.MarkUpdated()
}

// Renew renews the subscription for another period.
func (s *Subscription) Renew(period SubscriptionPeriod) {
	s.CurrentPeriod = period
	s.Status = SubscriptionStatusActive
	if s.Trial != nil && !s.Trial.Converted {
		s.Trial.Converted = true
	}
	s.MarkUpdated()
}

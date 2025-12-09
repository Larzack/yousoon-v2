// Package model contains GraphQL model types for the Discovery service.
package model

import (
	"time"
)

// =============================================================================
// Entity Types
// =============================================================================

// Offer represents an offer/deal in the GraphQL layer.
type Offer struct {
	ID                 string                 `json:"id"`
	PartnerID          string                 `json:"partnerId"`
	EstablishmentID    string                 `json:"establishmentId"`
	Title              string                 `json:"title"`
	Description        string                 `json:"description"`
	ShortDescription   string                 `json:"shortDescription"`
	CategoryID         string                 `json:"categoryId"`
	Tags               []string               `json:"tags"`
	Discount           *Discount              `json:"discount"`
	Conditions         []*Condition           `json:"conditions"`
	TermsAndConditions *string                `json:"termsAndConditions"`
	Validity           *Validity              `json:"validity"`
	Schedule           *Schedule              `json:"schedule"`
	Quota              *Quota                 `json:"quota"`
	Images             []*OfferImage          `json:"images"`
	Partner            *PartnerSnapshot       `json:"partner"`
	Establishment      *EstablishmentSnapshot `json:"establishment"`
	Stats              *OfferStats            `json:"stats"`
	Status             OfferStatus            `json:"status"`
	Moderation         *Moderation            `json:"moderation"`
	IsActive           bool                   `json:"isActive"`
	IsAvailableNow     bool                   `json:"isAvailableNow"`
	RemainingQuota     *int                   `json:"remainingQuota"`
	CreatedAt          time.Time              `json:"createdAt"`
	UpdatedAt          time.Time              `json:"updatedAt"`
	PublishedAt        *time.Time             `json:"publishedAt"`
}

// Category represents a category in the GraphQL layer.
type Category struct {
	ID          string           `json:"id"`
	Slug        string           `json:"slug"`
	Name        *LocalizedString `json:"name"`
	Description *LocalizedString `json:"description"`
	Icon        *string          `json:"icon"`
	Color       *string          `json:"color"`
	Image       *string          `json:"image"`
	ParentID    *string          `json:"parentId"`
	Order       int              `json:"order"`
	IsActive    bool             `json:"isActive"`
	OfferCount  int              `json:"offerCount"`
	CreatedAt   time.Time        `json:"createdAt"`
	UpdatedAt   time.Time        `json:"updatedAt"`
}

// =============================================================================
// Value Object Types
// =============================================================================

// Discount represents the discount information.
type Discount struct {
	Type            DiscountType `json:"type"`
	Value           int          `json:"value"`
	OriginalPrice   *Money       `json:"originalPrice"`
	DiscountedPrice *Money       `json:"discountedPrice"`
	Formula         *string      `json:"formula"`
}

// Money represents a monetary amount.
type Money struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}

// Condition represents an offer condition.
type Condition struct {
	Type  ConditionType `json:"type"`
	Value string        `json:"value"`
	Label string        `json:"label"`
}

// Validity represents the validity period.
type Validity struct {
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
	Timezone  string    `json:"timezone"`
}

// Schedule represents the offer schedule.
type Schedule struct {
	AllDay bool        `json:"allDay"`
	Slots  []*TimeSlot `json:"slots"`
}

// TimeSlot represents a time slot.
type TimeSlot struct {
	DayOfWeek int    `json:"dayOfWeek"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

// Quota represents offer quotas.
type Quota struct {
	Total     *int `json:"total"`
	PerUser   *int `json:"perUser"`
	PerDay    *int `json:"perDay"`
	Used      int  `json:"used"`
	Remaining *int `json:"remaining"`
}

// OfferImage represents an offer image.
type OfferImage struct {
	URL       string  `json:"url"`
	Alt       *string `json:"alt"`
	IsPrimary bool    `json:"isPrimary"`
	Order     int     `json:"order"`
}

// PartnerSnapshot represents denormalized partner data.
type PartnerSnapshot struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Logo     *string `json:"logo"`
	Category string  `json:"category"`
}

// EstablishmentSnapshot represents denormalized establishment data.
type EstablishmentSnapshot struct {
	ID       string       `json:"id"`
	Name     string       `json:"name"`
	Address  string       `json:"address"`
	City     string       `json:"city"`
	Location *GeoLocation `json:"location"`
}

// GeoLocation represents a geographic location.
type GeoLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// OfferStats represents offer statistics.
type OfferStats struct {
	Views       int     `json:"views"`
	Clicks      int     `json:"clicks"`
	Bookings    int     `json:"bookings"`
	Checkins    int     `json:"checkins"`
	Favorites   int     `json:"favorites"`
	AvgRating   float64 `json:"avgRating"`
	ReviewCount int     `json:"reviewCount"`
}

// Moderation represents moderation status.
type Moderation struct {
	Status     ModerationStatus `json:"status"`
	ReviewedBy *string          `json:"reviewedBy"`
	ReviewedAt *time.Time       `json:"reviewedAt"`
	Comment    *string          `json:"comment"`
}

// LocalizedString represents translated text.
type LocalizedString struct {
	FR string  `json:"fr"`
	EN *string `json:"en"`
}

// CategoryTree represents a category with children.
type CategoryTree struct {
	Category *Category       `json:"category"`
	Children []*CategoryTree `json:"children"`
}

// CategorySummary represents a category summary with offer count.
type CategorySummary struct {
	ID         string           `json:"id"`
	Name       *LocalizedString `json:"name"`
	Slug       string           `json:"slug"`
	Icon       *string          `json:"icon"`
	OfferCount int              `json:"offerCount"`
}

// =============================================================================
// Response Types
// =============================================================================

// OfferListResult represents a paginated list of offers.
type OfferListResult struct {
	Offers  []*Offer `json:"offers"`
	Total   int      `json:"total"`
	HasMore bool     `json:"hasMore"`
}

// OfferSearchResult represents search results.
type OfferSearchResult struct {
	Offers  []*OfferSummary `json:"offers"`
	Total   int             `json:"total"`
	HasMore bool            `json:"hasMore"`
}

// OfferSummary represents a lightweight offer summary.
type OfferSummary struct {
	ID                string       `json:"id"`
	PartnerID         string       `json:"partnerId"`
	EstablishmentID   string       `json:"establishmentId"`
	Title             string       `json:"title"`
	ShortDescription  string       `json:"shortDescription"`
	CategoryID        string       `json:"categoryId"`
	DiscountType      DiscountType `json:"discountType"`
	DiscountValue     int          `json:"discountValue"`
	Status            OfferStatus  `json:"status"`
	PartnerName       string       `json:"partnerName"`
	EstablishmentName string       `json:"establishmentName"`
	EstablishmentCity string       `json:"establishmentCity"`
	Location          *GeoLocation `json:"location"`
	Views             int          `json:"views"`
	AvgRating         float64      `json:"avgRating"`
	ReviewCount       int          `json:"reviewCount"`
	PublishedAt       *time.Time   `json:"publishedAt"`
	Distance          *float64     `json:"distance"`
}

// AutocompleteResult represents autocomplete suggestions.
type AutocompleteResult struct {
	Suggestions []string `json:"suggestions"`
}

// TrendingOffer represents a trending offer.
type TrendingOffer struct {
	Offer  *Offer  `json:"offer"`
	Score  float64 `json:"score"`
	Reason string  `json:"reason"`
}

// =============================================================================
// Enums
// =============================================================================

// OfferStatus represents the status of an offer.
type OfferStatus string

const (
	OfferStatusDraft    OfferStatus = "DRAFT"
	OfferStatusPending  OfferStatus = "PENDING"
	OfferStatusActive   OfferStatus = "ACTIVE"
	OfferStatusPaused   OfferStatus = "PAUSED"
	OfferStatusExpired  OfferStatus = "EXPIRED"
	OfferStatusArchived OfferStatus = "ARCHIVED"
)

func (e OfferStatus) IsValid() bool {
	switch e {
	case OfferStatusDraft, OfferStatusPending, OfferStatusActive,
		OfferStatusPaused, OfferStatusExpired, OfferStatusArchived:
		return true
	}
	return false
}

func (e OfferStatus) String() string {
	return string(e)
}

// ModerationStatus represents the moderation status.
type ModerationStatus string

const (
	ModerationStatusPending  ModerationStatus = "PENDING"
	ModerationStatusApproved ModerationStatus = "APPROVED"
	ModerationStatusRejected ModerationStatus = "REJECTED"
)

func (e ModerationStatus) IsValid() bool {
	switch e {
	case ModerationStatusPending, ModerationStatusApproved, ModerationStatusRejected:
		return true
	}
	return false
}

func (e ModerationStatus) String() string {
	return string(e)
}

// DiscountType represents the type of discount.
type DiscountType string

const (
	DiscountTypePercentage DiscountType = "PERCENTAGE"
	DiscountTypeFixed      DiscountType = "FIXED"
	DiscountTypeFormula    DiscountType = "FORMULA"
)

func (e DiscountType) IsValid() bool {
	switch e {
	case DiscountTypePercentage, DiscountTypeFixed, DiscountTypeFormula:
		return true
	}
	return false
}

func (e DiscountType) String() string {
	return string(e)
}

// ConditionType represents the type of condition.
type ConditionType string

const (
	ConditionTypeMinPurchase   ConditionType = "MIN_PURCHASE"
	ConditionTypeMinPeople     ConditionType = "MIN_PEOPLE"
	ConditionTypeFirstVisit    ConditionType = "FIRST_VISIT"
	ConditionTypeSpecificDays  ConditionType = "SPECIFIC_DAYS"
	ConditionTypeSpecificHours ConditionType = "SPECIFIC_HOURS"
	ConditionTypeOther         ConditionType = "OTHER"
)

func (e ConditionType) IsValid() bool {
	switch e {
	case ConditionTypeMinPurchase, ConditionTypeMinPeople, ConditionTypeFirstVisit,
		ConditionTypeSpecificDays, ConditionTypeSpecificHours, ConditionTypeOther:
		return true
	}
	return false
}

func (e ConditionType) String() string {
	return string(e)
}

// OfferSortBy represents sorting options for offers.
type OfferSortBy string

const (
	OfferSortByRelevance  OfferSortBy = "RELEVANCE"
	OfferSortByNewest     OfferSortBy = "NEWEST"
	OfferSortByPopularity OfferSortBy = "POPULARITY"
	OfferSortByRating     OfferSortBy = "RATING"
	OfferSortByDistance   OfferSortBy = "DISTANCE"
	OfferSortByDiscount   OfferSortBy = "DISCOUNT"
)

func (e OfferSortBy) IsValid() bool {
	switch e {
	case OfferSortByRelevance, OfferSortByNewest, OfferSortByPopularity,
		OfferSortByRating, OfferSortByDistance, OfferSortByDiscount:
		return true
	}
	return false
}

func (e OfferSortBy) String() string {
	return string(e)
}

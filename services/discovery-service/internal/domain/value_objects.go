// Package domain contains value objects for the Discovery service.
package domain

import (
	"errors"
	"time"
)

// =============================================================================
// ID Types
// =============================================================================

// OfferID represents a unique offer identifier.
type OfferID string

func (id OfferID) String() string { return string(id) }
func (id OfferID) IsEmpty() bool  { return id == "" }

// CategoryID represents a unique category identifier.
type CategoryID string

func (id CategoryID) String() string { return string(id) }
func (id CategoryID) IsEmpty() bool  { return id == "" }

// PartnerID represents a reference to a partner (cross-context).
type PartnerID string

func (id PartnerID) String() string { return string(id) }
func (id PartnerID) IsEmpty() bool  { return id == "" }

// EstablishmentID represents a reference to an establishment (cross-context).
type EstablishmentID string

func (id EstablishmentID) String() string { return string(id) }
func (id EstablishmentID) IsEmpty() bool  { return id == "" }

// UserID represents a reference to a user (cross-context).
type UserID string

func (id UserID) String() string { return string(id) }
func (id UserID) IsEmpty() bool  { return id == "" }

// =============================================================================
// Discount Value Object
// =============================================================================

// DiscountType represents the type of discount.
type DiscountType string

const (
	DiscountTypePercentage DiscountType = "percentage"
	DiscountTypeFixed      DiscountType = "fixed"
	DiscountTypeFormula    DiscountType = "formula"
)

// Discount represents a discount or reduction.
type Discount struct {
	Type          DiscountType `json:"type" bson:"type"`
	Value         int          `json:"value" bson:"value"`                  // Percentage (0-100) or cents
	OriginalPrice *int64       `json:"originalPrice" bson:"original_price"` // Original price in cents
	Formula       string       `json:"formula" bson:"formula"`              // E.g., "1 acheté = 1 offert"
}

// NewPercentageDiscount creates a percentage discount.
func NewPercentageDiscount(percentage int) Discount {
	return Discount{
		Type:  DiscountTypePercentage,
		Value: percentage,
	}
}

// NewFixedDiscount creates a fixed amount discount.
func NewFixedDiscount(amountCents int) Discount {
	return Discount{
		Type:  DiscountTypeFixed,
		Value: amountCents,
	}
}

// NewFormulaDiscount creates a formula-based discount.
func NewFormulaDiscount(formula string) Discount {
	return Discount{
		Type:    DiscountTypeFormula,
		Formula: formula,
	}
}

// Validate validates the discount.
func (d Discount) Validate() error {
	switch d.Type {
	case DiscountTypePercentage:
		if d.Value < 1 || d.Value > 100 {
			return errors.New("percentage must be between 1 and 100")
		}
	case DiscountTypeFixed:
		if d.Value < 1 {
			return errors.New("fixed discount must be positive")
		}
	case DiscountTypeFormula:
		if d.Formula == "" {
			return errors.New("formula is required")
		}
	default:
		return errors.New("invalid discount type")
	}
	return nil
}

// Apply applies the discount to a price.
func (d Discount) Apply(priceCents int64) int64 {
	switch d.Type {
	case DiscountTypePercentage:
		reduction := priceCents * int64(d.Value) / 100
		return priceCents - reduction
	case DiscountTypeFixed:
		result := priceCents - int64(d.Value)
		if result < 0 {
			return 0
		}
		return result
	case DiscountTypeFormula:
		// Formula discounts don't have a numeric calculation
		return priceCents
	}
	return priceCents
}

// GetDisplayText returns a human-readable discount text.
func (d Discount) GetDisplayText() string {
	switch d.Type {
	case DiscountTypePercentage:
		return "-" + string(rune(d.Value+'0')) + "%"
	case DiscountTypeFixed:
		euros := d.Value / 100
		return "-" + string(rune(euros+'0')) + "€"
	case DiscountTypeFormula:
		return d.Formula
	}
	return ""
}

// =============================================================================
// Condition Value Object
// =============================================================================

// ConditionType represents the type of condition.
type ConditionType string

const (
	ConditionTypeMinPurchase ConditionType = "min_purchase"
	ConditionTypeMinPeople   ConditionType = "min_people"
	ConditionTypeFirstVisit  ConditionType = "first_visit"
	ConditionTypeNewUser     ConditionType = "new_user"
	ConditionTypeOther       ConditionType = "other"
)

// Condition represents a condition for using an offer.
type Condition struct {
	Type  ConditionType `json:"type" bson:"type"`
	Value interface{}   `json:"value" bson:"value"`
	Label string        `json:"label" bson:"label"`
}

// NewMinPurchaseCondition creates a minimum purchase condition.
func NewMinPurchaseCondition(amountCents int, label string) Condition {
	return Condition{
		Type:  ConditionTypeMinPurchase,
		Value: amountCents,
		Label: label,
	}
}

// NewMinPeopleCondition creates a minimum people condition.
func NewMinPeopleCondition(count int, label string) Condition {
	return Condition{
		Type:  ConditionTypeMinPeople,
		Value: count,
		Label: label,
	}
}

// =============================================================================
// Validity Value Object
// =============================================================================

// Validity represents the validity period of an offer.
type Validity struct {
	StartDate time.Time `json:"startDate" bson:"start_date"`
	EndDate   time.Time `json:"endDate" bson:"end_date"`
	Timezone  string    `json:"timezone" bson:"timezone"`
}

// NewValidity creates a new validity period.
func NewValidity(startDate, endDate time.Time, timezone string) (Validity, error) {
	if endDate.Before(startDate) {
		return Validity{}, errors.New("end date must be after start date")
	}
	if timezone == "" {
		timezone = "Europe/Paris"
	}
	return Validity{
		StartDate: startDate,
		EndDate:   endDate,
		Timezone:  timezone,
	}, nil
}

// Validate validates the validity period.
func (v Validity) Validate() error {
	if v.EndDate.Before(v.StartDate) {
		return errors.New("end date must be after start date")
	}
	return nil
}

// IsExpired checks if the validity period has expired.
func (v Validity) IsExpired() bool {
	return time.Now().After(v.EndDate)
}

// HasStarted checks if the validity period has started.
func (v Validity) HasStarted() bool {
	return time.Now().After(v.StartDate)
}

// IsActive checks if the current time is within the validity period.
func (v Validity) IsActive() bool {
	now := time.Now()
	return now.After(v.StartDate) && now.Before(v.EndDate)
}

// DaysRemaining returns the number of days remaining.
func (v Validity) DaysRemaining() int {
	if v.IsExpired() {
		return 0
	}
	return int(v.EndDate.Sub(time.Now()).Hours() / 24)
}

// =============================================================================
// Schedule Value Object
// =============================================================================

// Schedule represents when the offer is available.
type Schedule struct {
	AllDay bool       `json:"allDay" bson:"all_day"`
	Slots  []TimeSlot `json:"slots" bson:"slots"`
}

// TimeSlot represents a time slot when the offer is available.
type TimeSlot struct {
	DayOfWeek int    `json:"dayOfWeek" bson:"day_of_week"` // 0 = Sunday
	StartTime string `json:"startTime" bson:"start_time"`  // "17:00"
	EndTime   string `json:"endTime" bson:"end_time"`      // "20:00"
}

// NewAllDaySchedule creates a schedule that's available all day, every day.
func NewAllDaySchedule() Schedule {
	return Schedule{
		AllDay: true,
		Slots:  nil,
	}
}

// NewScheduleWithSlots creates a schedule with specific time slots.
func NewScheduleWithSlots(slots []TimeSlot) Schedule {
	return Schedule{
		AllDay: false,
		Slots:  slots,
	}
}

// IsAvailableNow checks if the offer is available at the current time.
func (s Schedule) IsAvailableNow() bool {
	if s.AllDay {
		return true
	}

	now := time.Now()
	currentDayOfWeek := int(now.Weekday())
	currentTime := now.Format("15:04")

	for _, slot := range s.Slots {
		if slot.DayOfWeek == currentDayOfWeek {
			if currentTime >= slot.StartTime && currentTime <= slot.EndTime {
				return true
			}
		}
	}

	return false
}

// IsAvailableOn checks if the offer is available on a specific day and time.
func (s Schedule) IsAvailableOn(dayOfWeek int, timeStr string) bool {
	if s.AllDay {
		return true
	}

	for _, slot := range s.Slots {
		if slot.DayOfWeek == dayOfWeek {
			if timeStr >= slot.StartTime && timeStr <= slot.EndTime {
				return true
			}
		}
	}

	return false
}

// GetSlotsForDay returns all slots for a specific day.
func (s Schedule) GetSlotsForDay(dayOfWeek int) []TimeSlot {
	if s.AllDay {
		return []TimeSlot{{DayOfWeek: dayOfWeek, StartTime: "00:00", EndTime: "23:59"}}
	}

	result := make([]TimeSlot, 0)
	for _, slot := range s.Slots {
		if slot.DayOfWeek == dayOfWeek {
			result = append(result, slot)
		}
	}
	return result
}

// =============================================================================
// Quota Value Object
// =============================================================================

// Quota represents usage limits for an offer.
type Quota struct {
	Total   *int `json:"total" bson:"total"`      // Total limit (nil = unlimited)
	PerUser *int `json:"perUser" bson:"per_user"` // Limit per user (nil = unlimited)
	PerDay  *int `json:"perDay" bson:"per_day"`   // Limit per day (nil = unlimited)
	Used    int  `json:"used" bson:"used"`        // Current usage count
}

// NewUnlimitedQuota creates an unlimited quota.
func NewUnlimitedQuota() Quota {
	return Quota{}
}

// NewQuota creates a quota with limits.
func NewQuota(total, perUser, perDay *int) Quota {
	return Quota{
		Total:   total,
		PerUser: perUser,
		PerDay:  perDay,
		Used:    0,
	}
}

// IsExhausted checks if the quota is exhausted.
func (q Quota) IsExhausted() bool {
	if q.Total == nil {
		return false
	}
	return q.Used >= *q.Total
}

// Remaining returns the remaining quota.
func (q Quota) Remaining() int {
	if q.Total == nil {
		return -1 // Unlimited
	}
	remaining := *q.Total - q.Used
	if remaining < 0 {
		return 0
	}
	return remaining
}

// AvailabilityPercentage returns the percentage of quota still available.
func (q Quota) AvailabilityPercentage() float64 {
	if q.Total == nil || *q.Total == 0 {
		return 100.0
	}
	return float64(q.Remaining()) / float64(*q.Total) * 100
}

// =============================================================================
// OfferImage Value Object
// =============================================================================

// OfferImage represents an image attached to an offer.
type OfferImage struct {
	URL       string `json:"url" bson:"url"`
	Alt       string `json:"alt" bson:"alt"`
	IsPrimary bool   `json:"isPrimary" bson:"is_primary"`
	Order     int    `json:"order" bson:"order"`
}

// NewOfferImage creates a new offer image.
func NewOfferImage(url, alt string, isPrimary bool, order int) OfferImage {
	return OfferImage{
		URL:       url,
		Alt:       alt,
		IsPrimary: isPrimary,
		Order:     order,
	}
}

// =============================================================================
// OfferStats Value Object
// =============================================================================

// OfferStats represents statistics for an offer.
type OfferStats struct {
	Views     int `json:"views" bson:"views"`
	Bookings  int `json:"bookings" bson:"bookings"`
	Checkins  int `json:"checkins" bson:"checkins"`
	Favorites int `json:"favorites" bson:"favorites"`
}

// =============================================================================
// Moderation Value Object
// =============================================================================

// Moderation represents the moderation state of an offer.
type Moderation struct {
	Status     ModerationStatus `json:"status" bson:"status"`
	ReviewerID *string          `json:"reviewerId" bson:"reviewer_id"`
	ReviewedAt *time.Time       `json:"reviewedAt" bson:"reviewed_at"`
	Comment    *string          `json:"comment" bson:"comment"`
}

// =============================================================================
// Snapshot Value Objects (for cross-context data)
// =============================================================================

// PartnerSnapshot represents denormalized partner data.
type PartnerSnapshot struct {
	Name     string `json:"name" bson:"name"`
	Logo     string `json:"logo" bson:"logo"`
	Category string `json:"category" bson:"category"`
}

// EstablishmentSnapshot represents denormalized establishment data.
type EstablishmentSnapshot struct {
	Name     string      `json:"name" bson:"name"`
	Address  string      `json:"address" bson:"address"`
	City     string      `json:"city" bson:"city"`
	Location GeoLocation `json:"location" bson:"location"`
}

// GeoLocation represents a geographic location (GeoJSON Point).
type GeoLocation struct {
	Type        string    `json:"type" bson:"type"`               // Always "Point"
	Coordinates []float64 `json:"coordinates" bson:"coordinates"` // [longitude, latitude]
}

// NewGeoLocation creates a new GeoLocation.
func NewGeoLocation(longitude, latitude float64) (GeoLocation, error) {
	if longitude < -180 || longitude > 180 {
		return GeoLocation{}, errors.New("longitude must be between -180 and 180")
	}
	if latitude < -90 || latitude > 90 {
		return GeoLocation{}, errors.New("latitude must be between -90 and 90")
	}
	return GeoLocation{
		Type:        "Point",
		Coordinates: []float64{longitude, latitude},
	}, nil
}

// Longitude returns the longitude.
func (g GeoLocation) Longitude() float64 {
	if len(g.Coordinates) < 2 {
		return 0
	}
	return g.Coordinates[0]
}

// Latitude returns the latitude.
func (g GeoLocation) Latitude() float64 {
	if len(g.Coordinates) < 2 {
		return 0
	}
	return g.Coordinates[1]
}

// OfferSnapshot represents an immutable snapshot of an offer for bookings.
type OfferSnapshot struct {
	ID              OfferID         `json:"id" bson:"id"`
	PartnerID       PartnerID       `json:"partnerId" bson:"partner_id"`
	EstablishmentID EstablishmentID `json:"establishmentId" bson:"establishment_id"`
	Title           string          `json:"title" bson:"title"`
	Description     string          `json:"description" bson:"description"`
	Discount        Discount        `json:"discount" bson:"discount"`
	CategoryName    string          `json:"categoryName" bson:"category_name"`
	Location        GeoLocation     `json:"location" bson:"location"`
	CapturedAt      time.Time       `json:"capturedAt" bson:"captured_at"`
}

// =============================================================================
// Summary Value Objects (for lists/search)
// =============================================================================

// OfferSummary represents a summary of an offer for lists.
type OfferSummary struct {
	ID                OfferID     `json:"id" bson:"_id"`
	Title             string      `json:"title" bson:"title"`
	ShortDescription  string      `json:"shortDescription" bson:"short_description"`
	Discount          Discount    `json:"discount" bson:"discount"`
	PrimaryImage      string      `json:"primaryImage" bson:"primary_image"`
	PartnerName       string      `json:"partnerName" bson:"partner_name"`
	EstablishmentName string      `json:"establishmentName" bson:"establishment_name"`
	City              string      `json:"city" bson:"city"`
	Location          GeoLocation `json:"location" bson:"location"`
	CategoryID        CategoryID  `json:"categoryId" bson:"category_id"`
	Distance          *float64    `json:"distance,omitempty" bson:"distance,omitempty"` // km from user
}

// CategorySummary represents a summary of a category.
type CategorySummary struct {
	ID         CategoryID `json:"id" bson:"_id"`
	Slug       string     `json:"slug" bson:"slug"`
	Name       string     `json:"name" bson:"name"`
	Icon       string     `json:"icon" bson:"icon"`
	OfferCount int        `json:"offerCount" bson:"offer_count"`
}

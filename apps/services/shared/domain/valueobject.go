package domain

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// =============================================================================
// Email Value Object
// =============================================================================

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// Email represents a validated email address.
type Email struct {
	value string
}

// NewEmail creates a new Email value object.
func NewEmail(email string) (Email, error) {
	normalized := strings.ToLower(strings.TrimSpace(email))
	if !emailRegex.MatchString(normalized) {
		return Email{}, ErrInvalidEmail
	}
	return Email{value: normalized}, nil
}

// MustEmail creates an Email, panicking on error. Use only for tests.
func MustEmail(email string) Email {
	e, err := NewEmail(email)
	if err != nil {
		panic(err)
	}
	return e
}

// String returns the email address.
func (e Email) String() string {
	return e.value
}

// Equals checks equality with another email.
func (e Email) Equals(other Email) bool {
	return e.value == other.value
}

// IsZero returns true if the email is empty.
func (e Email) IsZero() bool {
	return e.value == ""
}

// Domain returns the domain part of the email.
func (e Email) Domain() string {
	parts := strings.Split(e.value, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

// =============================================================================
// Phone Value Object
// =============================================================================

var phoneRegex = regexp.MustCompile(`^\+[1-9]\d{1,14}$`)

// Phone represents a validated phone number in E.164 format.
type Phone struct {
	value string
}

// NewPhone creates a new Phone value object.
func NewPhone(phone string) (Phone, error) {
	normalized := strings.ReplaceAll(strings.TrimSpace(phone), " ", "")
	if !phoneRegex.MatchString(normalized) {
		return Phone{}, ErrInvalidPhone
	}
	return Phone{value: normalized}, nil
}

// String returns the phone number.
func (p Phone) String() string {
	return p.value
}

// Equals checks equality with another phone.
func (p Phone) Equals(other Phone) bool {
	return p.value == other.value
}

// IsZero returns true if the phone is empty.
func (p Phone) IsZero() bool {
	return p.value == ""
}

// CountryCode returns the country code part.
func (p Phone) CountryCode() string {
	if len(p.value) < 2 {
		return ""
	}
	// Simple extraction - could be improved with libphonenumber
	if strings.HasPrefix(p.value, "+33") {
		return "+33"
	}
	return p.value[:3]
}

// =============================================================================
// Money Value Object
// =============================================================================

// Money represents a monetary amount in cents.
type Money struct {
	Amount   int64  `json:"amount" bson:"amount"`     // In cents (990 = 9.90€)
	Currency string `json:"currency" bson:"currency"` // ISO 4217 (EUR, USD)
}

// NewMoney creates a new Money value object.
func NewMoney(amount int64, currency string) Money {
	return Money{
		Amount:   amount,
		Currency: strings.ToUpper(currency),
	}
}

// NewMoneyEUR creates a Money in EUR.
func NewMoneyEUR(amount int64) Money {
	return NewMoney(amount, "EUR")
}

// FromFloat creates Money from a float (e.g., 9.90 -> 990 cents).
func FromFloat(amount float64, currency string) Money {
	return NewMoney(int64(math.Round(amount*100)), currency)
}

// ToFloat returns the amount as a float.
func (m Money) ToFloat() float64 {
	return float64(m.Amount) / 100
}

// Add adds two Money values.
func (m Money) Add(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, fmt.Errorf("currency mismatch: %s vs %s", m.Currency, other.Currency)
	}
	return Money{Amount: m.Amount + other.Amount, Currency: m.Currency}, nil
}

// Subtract subtracts another Money value.
func (m Money) Subtract(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, fmt.Errorf("currency mismatch: %s vs %s", m.Currency, other.Currency)
	}
	return Money{Amount: m.Amount - other.Amount, Currency: m.Currency}, nil
}

// Multiply multiplies by a factor.
func (m Money) Multiply(factor float64) Money {
	return Money{Amount: int64(math.Round(float64(m.Amount) * factor)), Currency: m.Currency}
}

// Percentage applies a percentage discount.
func (m Money) Percentage(percent int) Money {
	reduction := m.Amount * int64(percent) / 100
	return Money{Amount: m.Amount - reduction, Currency: m.Currency}
}

// IsZero returns true if the amount is zero.
func (m Money) IsZero() bool {
	return m.Amount == 0
}

// IsPositive returns true if the amount is positive.
func (m Money) IsPositive() bool {
	return m.Amount > 0
}

// Equals checks equality.
func (m Money) Equals(other Money) bool {
	return m.Amount == other.Amount && m.Currency == other.Currency
}

// String returns a formatted string.
func (m Money) String() string {
	return fmt.Sprintf("%.2f %s", m.ToFloat(), m.Currency)
}

// =============================================================================
// GeoLocation Value Object
// =============================================================================

// GeoLocation represents a geographic coordinate.
type GeoLocation struct {
	Type        string    `json:"type" bson:"type"`               // Always "Point" for GeoJSON
	Coordinates []float64 `json:"coordinates" bson:"coordinates"` // [longitude, latitude]
}

// NewGeoLocation creates a new GeoLocation.
func NewGeoLocation(longitude, latitude float64) (GeoLocation, error) {
	if longitude < -180 || longitude > 180 {
		return GeoLocation{}, fmt.Errorf("invalid longitude: %f", longitude)
	}
	if latitude < -90 || latitude > 90 {
		return GeoLocation{}, fmt.Errorf("invalid latitude: %f", latitude)
	}
	return GeoLocation{
		Type:        "Point",
		Coordinates: []float64{longitude, latitude},
	}, nil
}

// MustGeoLocation creates a GeoLocation, panicking on error.
func MustGeoLocation(longitude, latitude float64) GeoLocation {
	loc, err := NewGeoLocation(longitude, latitude)
	if err != nil {
		panic(err)
	}
	return loc
}

// Longitude returns the longitude.
func (g GeoLocation) Longitude() float64 {
	if len(g.Coordinates) < 1 {
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

// DistanceTo calculates the distance to another location in kilometers.
// Uses the Haversine formula.
func (g GeoLocation) DistanceTo(other GeoLocation) float64 {
	const earthRadius = 6371 // km

	lat1 := g.Latitude() * math.Pi / 180
	lat2 := other.Latitude() * math.Pi / 180
	dLat := (other.Latitude() - g.Latitude()) * math.Pi / 180
	dLon := (other.Longitude() - g.Longitude()) * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

// IsZero returns true if the location is not set.
func (g GeoLocation) IsZero() bool {
	return len(g.Coordinates) == 0
}

// =============================================================================
// Address Value Object
// =============================================================================

// Address represents a physical address.
type Address struct {
	Street       string `json:"street" bson:"street"`
	StreetNumber string `json:"street_number,omitempty" bson:"street_number,omitempty"`
	Complement   string `json:"complement,omitempty" bson:"complement,omitempty"`
	PostalCode   string `json:"postal_code" bson:"postal_code"`
	City         string `json:"city" bson:"city"`
	Country      string `json:"country" bson:"country"` // ISO 3166-1 alpha-2
	Formatted    string `json:"formatted,omitempty" bson:"formatted,omitempty"`
}

// NewAddress creates a new Address.
func NewAddress(street, postalCode, city, country string) Address {
	addr := Address{
		Street:     street,
		PostalCode: postalCode,
		City:       city,
		Country:    strings.ToUpper(country),
	}
	addr.Formatted = addr.Format()
	return addr
}

// Format returns a formatted address string.
func (a Address) Format() string {
	parts := []string{}
	if a.StreetNumber != "" {
		parts = append(parts, a.StreetNumber)
	}
	if a.Street != "" {
		parts = append(parts, a.Street)
	}
	if a.Complement != "" {
		parts = append(parts, a.Complement)
	}
	if a.PostalCode != "" || a.City != "" {
		parts = append(parts, fmt.Sprintf("%s %s", a.PostalCode, a.City))
	}
	if a.Country != "" {
		parts = append(parts, a.Country)
	}
	return strings.Join(parts, ", ")
}

// IsZero returns true if the address is empty.
func (a Address) IsZero() bool {
	return a.Street == "" && a.City == ""
}

// =============================================================================
// QRCode Value Object
// =============================================================================

// QRCode represents a unique QR code for check-in.
type QRCode struct {
	Code      string    `json:"code" bson:"code"`
	Signature string    `json:"signature" bson:"signature"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	ExpiresAt time.Time `json:"expires_at" bson:"expires_at"`
}

// NewQRCode generates a new QR code with signature.
func NewQRCode(secretKey string, expiresIn time.Duration) QRCode {
	code := uuid.New().String()
	now := time.Now().UTC()

	// Create HMAC signature
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(code))
	signature := hex.EncodeToString(h.Sum(nil))[:16] // Truncate for readability

	return QRCode{
		Code:      code,
		Signature: signature,
		CreatedAt: now,
		ExpiresAt: now.Add(expiresIn),
	}
}

// FullCode returns the complete code with signature.
func (q QRCode) FullCode() string {
	return fmt.Sprintf("%s.%s", q.Code, q.Signature)
}

// Matches checks if a scanned code matches.
func (q QRCode) Matches(scanned string) bool {
	return q.Code == scanned || q.FullCode() == scanned
}

// IsExpired returns true if the QR code has expired.
func (q QRCode) IsExpired() bool {
	return time.Now().UTC().After(q.ExpiresAt)
}

// Validate checks if the code and signature match.
func (q QRCode) Validate(secretKey string) bool {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(q.Code))
	expectedSig := hex.EncodeToString(h.Sum(nil))[:16]
	return q.Signature == expectedSig
}

// =============================================================================
// Discount Value Object
// =============================================================================

// DiscountType defines the type of discount.
type DiscountType string

const (
	DiscountTypePercentage DiscountType = "percentage"
	DiscountTypeFixed      DiscountType = "fixed"
	DiscountTypeFormula    DiscountType = "formula"
)

// Discount represents a reduction on a price.
type Discount struct {
	Type          DiscountType `json:"type" bson:"type"`
	Value         int          `json:"value" bson:"value"`                                       // % or cents
	OriginalPrice *Money       `json:"original_price,omitempty" bson:"original_price,omitempty"` // Optional
	MaxDiscount   *Money       `json:"max_discount,omitempty" bson:"max_discount,omitempty"`     // Cap for percentage
	Formula       string       `json:"formula,omitempty" bson:"formula,omitempty"`               // e.g., "1 acheté = 1 offert"
}

// NewPercentageDiscount creates a percentage discount.
func NewPercentageDiscount(percent int) Discount {
	return Discount{
		Type:  DiscountTypePercentage,
		Value: percent,
	}
}

// NewFixedDiscount creates a fixed amount discount in cents.
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

// Apply applies the discount to a price.
func (d Discount) Apply(original Money) Money {
	switch d.Type {
	case DiscountTypePercentage:
		reduction := original.Amount * int64(d.Value) / 100
		if d.MaxDiscount != nil && reduction > d.MaxDiscount.Amount {
			reduction = d.MaxDiscount.Amount
		}
		return Money{Amount: original.Amount - reduction, Currency: original.Currency}
	case DiscountTypeFixed:
		return Money{Amount: original.Amount - int64(d.Value), Currency: original.Currency}
	default:
		return original
	}
}

// DisplayValue returns a display string for the discount.
func (d Discount) DisplayValue() string {
	switch d.Type {
	case DiscountTypePercentage:
		return fmt.Sprintf("-%d%%", d.Value)
	case DiscountTypeFixed:
		return fmt.Sprintf("-%.2f€", float64(d.Value)/100)
	case DiscountTypeFormula:
		return d.Formula
	default:
		return ""
	}
}

// IsValid validates the discount.
func (d Discount) IsValid() bool {
	switch d.Type {
	case DiscountTypePercentage:
		return d.Value > 0 && d.Value <= 100
	case DiscountTypeFixed:
		return d.Value > 0
	case DiscountTypeFormula:
		return d.Formula != ""
	default:
		return false
	}
}

// =============================================================================
// TimeSlot Value Object
// =============================================================================

// TimeSlot represents a time slot in a schedule.
type TimeSlot struct {
	DayOfWeek int    `json:"day_of_week" bson:"day_of_week"` // 0 = Sunday, 1 = Monday, ...
	StartTime string `json:"start_time" bson:"start_time"`   // "09:00"
	EndTime   string `json:"end_time" bson:"end_time"`       // "18:00"
}

// Schedule represents the validity schedule of an offer.
type Schedule struct {
	AllDay    bool       `json:"all_day" bson:"all_day"`
	StartDate time.Time  `json:"start_date" bson:"start_date"`
	EndDate   time.Time  `json:"end_date" bson:"end_date"`
	Timezone  string     `json:"timezone" bson:"timezone"`
	Slots     []TimeSlot `json:"slots,omitempty" bson:"slots,omitempty"`
}

// IsExpired returns true if the schedule has ended.
func (s Schedule) IsExpired() bool {
	return time.Now().UTC().After(s.EndDate)
}

// IsActiveNow returns true if the schedule is currently active.
func (s Schedule) IsActiveNow() bool {
	now := time.Now().UTC()
	return now.After(s.StartDate) && now.Before(s.EndDate)
}

// IsActiveAt checks if the schedule is active at a given time.
func (s Schedule) IsActiveAt(t time.Time) bool {
	if !t.After(s.StartDate) || !t.Before(s.EndDate) {
		return false
	}
	if s.AllDay {
		return true
	}
	// Check time slots
	dayOfWeek := int(t.Weekday())
	timeStr := t.Format("15:04")
	for _, slot := range s.Slots {
		if slot.DayOfWeek == dayOfWeek {
			if timeStr >= slot.StartTime && timeStr <= slot.EndTime {
				return true
			}
		}
	}
	return false
}

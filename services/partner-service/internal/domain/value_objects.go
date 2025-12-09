package domain

import (
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// =============================================================================
// ID Types
// =============================================================================

// PartnerID is the unique identifier for a partner.
type PartnerID string

// NewPartnerID generates a new partner ID.
func NewPartnerID() PartnerID {
	return PartnerID(uuid.New().String())
}

// String returns the string representation.
func (id PartnerID) String() string {
	return string(id)
}

// IsEmpty checks if the ID is empty.
func (id PartnerID) IsEmpty() bool {
	return string(id) == ""
}

// EstablishmentID is the unique identifier for an establishment.
type EstablishmentID string

// NewEstablishmentID generates a new establishment ID.
func NewEstablishmentID() EstablishmentID {
	return EstablishmentID(uuid.New().String())
}

// String returns the string representation.
func (id EstablishmentID) String() string {
	return string(id)
}

// IsEmpty checks if the ID is empty.
func (id EstablishmentID) IsEmpty() bool {
	return string(id) == ""
}

// TeamMemberID is the unique identifier for a team member.
type TeamMemberID string

// NewTeamMemberID generates a new team member ID.
func NewTeamMemberID() TeamMemberID {
	return TeamMemberID(uuid.New().String())
}

// String returns the string representation.
func (id TeamMemberID) String() string {
	return string(id)
}

// IsEmpty checks if the ID is empty.
func (id TeamMemberID) IsEmpty() bool {
	return string(id) == ""
}

// UserID is a cross-context reference to an Identity user.
type UserID string

// String returns the string representation.
func (id UserID) String() string {
	return string(id)
}

// IsEmpty checks if the ID is empty.
func (id UserID) IsEmpty() bool {
	return string(id) == ""
}

// =============================================================================
// Company Value Object
// =============================================================================

// Company contains the legal business information.
type Company struct {
	Name      string `json:"name" bson:"name"`           // Raison sociale
	TradeName string `json:"tradeName" bson:"tradeName"` // Nom commercial
	SIRET     string `json:"siret" bson:"siret"`
	VATNumber string `json:"vatNumber,omitempty" bson:"vatNumber,omitempty"`
	LegalForm string `json:"legalForm" bson:"legalForm"` // SARL, SAS, etc.
}

// NewCompany creates a new company value object.
func NewCompany(name, tradeName, siret, vatNumber, legalForm string) (Company, error) {
	if strings.TrimSpace(name) == "" {
		return Company{}, ErrCompanyNameRequired
	}
	if !isValidSIRET(siret) {
		return Company{}, ErrInvalidSIRET
	}
	return Company{
		Name:      strings.TrimSpace(name),
		TradeName: strings.TrimSpace(tradeName),
		SIRET:     siret,
		VATNumber: vatNumber,
		LegalForm: legalForm,
	}, nil
}

// isValidSIRET validates a French SIRET number.
func isValidSIRET(siret string) bool {
	// Remove spaces and check length
	siret = strings.ReplaceAll(siret, " ", "")
	if len(siret) != 14 {
		return false
	}
	// Check if all digits
	for _, c := range siret {
		if c < '0' || c > '9' {
			return false
		}
	}
	// Luhn algorithm check (optional, can be simplified)
	return true
}

// =============================================================================
// Branding Value Object
// =============================================================================

// Branding contains the partner's visual identity.
type Branding struct {
	Logo         string `json:"logo,omitempty" bson:"logo,omitempty"`
	CoverImage   string `json:"coverImage,omitempty" bson:"coverImage,omitempty"`
	PrimaryColor string `json:"primaryColor,omitempty" bson:"primaryColor,omitempty"`
	Description  string `json:"description,omitempty" bson:"description,omitempty"`
}

// NewBranding creates a new branding value object.
func NewBranding(logo, coverImage, primaryColor, description string) Branding {
	return Branding{
		Logo:         logo,
		CoverImage:   coverImage,
		PrimaryColor: primaryColor,
		Description:  description,
	}
}

// =============================================================================
// Contact Value Object
// =============================================================================

// Contact contains the primary contact information.
type Contact struct {
	FirstName string `json:"firstName" bson:"firstName"`
	LastName  string `json:"lastName" bson:"lastName"`
	Email     string `json:"email" bson:"email"`
	Phone     string `json:"phone" bson:"phone"`
	Role      string `json:"role,omitempty" bson:"role,omitempty"`
}

// NewContact creates a new contact value object.
func NewContact(firstName, lastName, email, phone, role string) (Contact, error) {
	if strings.TrimSpace(firstName) == "" || strings.TrimSpace(lastName) == "" {
		return Contact{}, ErrContactNameRequired
	}
	if strings.TrimSpace(email) == "" {
		return Contact{}, ErrContactEmailRequired
	}
	return Contact{
		FirstName: strings.TrimSpace(firstName),
		LastName:  strings.TrimSpace(lastName),
		Email:     strings.ToLower(strings.TrimSpace(email)),
		Phone:     phone,
		Role:      role,
	}, nil
}

// FullName returns the contact's full name.
func (c Contact) FullName() string {
	return c.FirstName + " " + c.LastName
}

// =============================================================================
// Address Value Object
// =============================================================================

// Address represents a physical address.
type Address struct {
	Street       string `json:"street" bson:"street"`
	StreetNumber string `json:"streetNumber,omitempty" bson:"streetNumber,omitempty"`
	Complement   string `json:"complement,omitempty" bson:"complement,omitempty"`
	PostalCode   string `json:"postalCode" bson:"postalCode"`
	City         string `json:"city" bson:"city"`
	Country      string `json:"country" bson:"country"` // ISO 3166-1 alpha-2
	Formatted    string `json:"formatted,omitempty" bson:"formatted,omitempty"`
}

// NewAddress creates a new address value object.
func NewAddress(street, streetNumber, complement, postalCode, city, country string) Address {
	addr := Address{
		Street:       street,
		StreetNumber: streetNumber,
		Complement:   complement,
		PostalCode:   postalCode,
		City:         city,
		Country:      strings.ToUpper(country),
	}
	addr.Formatted = addr.Format()
	return addr
}

// Format returns the formatted address string.
func (a Address) Format() string {
	parts := []string{}
	if a.StreetNumber != "" {
		parts = append(parts, a.StreetNumber)
	}
	parts = append(parts, a.Street)
	if a.Complement != "" {
		parts = append(parts, a.Complement)
	}
	street := strings.Join(parts, " ")
	return street + ", " + a.PostalCode + " " + a.City + ", " + a.Country
}

// Equals compares two addresses.
func (a Address) Equals(other Address) bool {
	return a.Street == other.Street &&
		a.StreetNumber == other.StreetNumber &&
		a.PostalCode == other.PostalCode &&
		a.City == other.City &&
		a.Country == other.Country
}

// =============================================================================
// GeoLocation Value Object
// =============================================================================

// GeoLocation represents geographic coordinates (GeoJSON format).
type GeoLocation struct {
	Type        string    `json:"type" bson:"type"`               // Always "Point"
	Coordinates []float64 `json:"coordinates" bson:"coordinates"` // [longitude, latitude]
}

// NewGeoLocation creates a new geolocation.
func NewGeoLocation(longitude, latitude float64) (GeoLocation, error) {
	if longitude < -180 || longitude > 180 {
		return GeoLocation{}, ErrInvalidLongitude
	}
	if latitude < -90 || latitude > 90 {
		return GeoLocation{}, ErrInvalidLatitude
	}
	return GeoLocation{
		Type:        "Point",
		Coordinates: []float64{longitude, latitude},
	}, nil
}

// Longitude returns the longitude.
func (g GeoLocation) Longitude() float64 {
	if len(g.Coordinates) >= 2 {
		return g.Coordinates[0]
	}
	return 0
}

// Latitude returns the latitude.
func (g GeoLocation) Latitude() float64 {
	if len(g.Coordinates) >= 2 {
		return g.Coordinates[1]
	}
	return 0
}

// =============================================================================
// Email Value Object
// =============================================================================

// Email represents a validated email address.
type Email struct {
	Value string `json:"value" bson:"value"`
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// NewEmail creates a new email value object.
func NewEmail(email string) (Email, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if !emailRegex.MatchString(email) {
		return Email{}, ErrInvalidEmail
	}
	return Email{Value: email}, nil
}

// String returns the string representation.
func (e Email) String() string {
	return e.Value
}

// Equals compares two emails.
func (e Email) Equals(other Email) bool {
	return e.Value == other.Value
}

// =============================================================================
// PartnerStats Value Object
// =============================================================================

// PartnerStats contains denormalized statistics for a partner.
type PartnerStats struct {
	TotalEstablishments int       `json:"totalEstablishments" bson:"totalEstablishments"`
	TotalOffers         int       `json:"totalOffers" bson:"totalOffers"`
	ActiveOffers        int       `json:"activeOffers" bson:"activeOffers"`
	TotalBookings       int       `json:"totalBookings" bson:"totalBookings"`
	TotalCheckins       int       `json:"totalCheckins" bson:"totalCheckins"`
	AvgRating           float64   `json:"avgRating" bson:"avgRating"`
	ReviewCount         int       `json:"reviewCount" bson:"reviewCount"`
	LastUpdated         time.Time `json:"lastUpdated" bson:"lastUpdated"`
}

// NewPartnerStats creates new partner stats with default values.
func NewPartnerStats() PartnerStats {
	return PartnerStats{
		TotalEstablishments: 0,
		TotalOffers:         0,
		ActiveOffers:        0,
		TotalBookings:       0,
		TotalCheckins:       0,
		AvgRating:           0.0,
		ReviewCount:         0,
		LastUpdated:         time.Now(),
	}
}

// =============================================================================
// PartnerStatus Enum
// =============================================================================

// PartnerStatus represents the status of a partner account.
type PartnerStatus string

const (
	PartnerStatusPending   PartnerStatus = "pending"   // Awaiting verification
	PartnerStatusActive    PartnerStatus = "active"    // Verified and active
	PartnerStatusSuspended PartnerStatus = "suspended" // Temporarily suspended
)

// AllPartnerStatuses returns all valid partner statuses.
func AllPartnerStatuses() []PartnerStatus {
	return []PartnerStatus{
		PartnerStatusPending,
		PartnerStatusActive,
		PartnerStatusSuspended,
	}
}

// IsValid checks if the status is valid.
func (s PartnerStatus) IsValid() bool {
	for _, valid := range AllPartnerStatuses() {
		if s == valid {
			return true
		}
	}
	return false
}

// String returns the string representation.
func (s PartnerStatus) String() string {
	return string(s)
}

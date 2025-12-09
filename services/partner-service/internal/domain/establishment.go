package domain

import (
	"time"
)

// =============================================================================
// Establishment Entity
// =============================================================================

// Establishment represents a physical location belonging to a partner.
// It's an entity within the Partner aggregate.
type Establishment struct {
	// Identity
	ID        EstablishmentID `json:"id" bson:"id"`
	PartnerID PartnerID       `json:"partnerId" bson:"partnerId"`

	// Basic Info
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`

	// Address
	Address Address `json:"address" bson:"address"`

	// Geolocation
	Location GeoLocation `json:"location" bson:"location"`

	// Contact
	Contact EstablishmentContact `json:"contact" bson:"contact"`

	// Opening Hours
	OpeningHours []OpeningHour `json:"openingHours" bson:"openingHours"`

	// Exceptional Closures
	Closures []Closure `json:"closures,omitempty" bson:"closures,omitempty"`

	// Media
	Images []Image `json:"images,omitempty" bson:"images,omitempty"`

	// Features/Amenities
	Features []string `json:"features,omitempty" bson:"features,omitempty"`

	// Type (more specific than partner category)
	Type string `json:"type" bson:"type"`

	// Price Range (1-4, € to €€€€)
	PriceRange int `json:"priceRange" bson:"priceRange"`

	// Status
	IsActive bool `json:"isActive" bson:"isActive"`

	// Timestamps
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

// =============================================================================
// Establishment Factory
// =============================================================================

// NewEstablishment creates a new establishment with required fields.
func NewEstablishment(name string, address Address, location GeoLocation) Establishment {
	now := time.Now()
	return Establishment{
		ID:           NewEstablishmentID(),
		Name:         name,
		Address:      address,
		Location:     location,
		OpeningHours: make([]OpeningHour, 0),
		Closures:     make([]Closure, 0),
		Images:       make([]Image, 0),
		Features:     make([]string, 0),
		PriceRange:   2, // Default: €€
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// =============================================================================
// Establishment Methods
// =============================================================================

// Update updates the establishment's basic info.
func (e *Establishment) Update(name, description string) {
	e.Name = name
	e.Description = description
	e.UpdatedAt = time.Now()
}

// SetContact sets the establishment's contact information.
func (e *Establishment) SetContact(contact EstablishmentContact) {
	e.Contact = contact
	e.UpdatedAt = time.Now()
}

// SetOpeningHours sets the establishment's opening hours.
func (e *Establishment) SetOpeningHours(hours []OpeningHour) {
	e.OpeningHours = hours
	e.UpdatedAt = time.Now()
}

// AddClosure adds an exceptional closure date.
func (e *Establishment) AddClosure(closure Closure) {
	e.Closures = append(e.Closures, closure)
	e.UpdatedAt = time.Now()
}

// RemoveClosure removes a closure by date.
func (e *Establishment) RemoveClosure(date time.Time) {
	for i, c := range e.Closures {
		if c.Date.Truncate(24 * time.Hour).Equal(date.Truncate(24 * time.Hour)) {
			e.Closures = append(e.Closures[:i], e.Closures[i+1:]...)
			e.UpdatedAt = time.Now()
			return
		}
	}
}

// AddImage adds an image to the establishment.
func (e *Establishment) AddImage(image Image) {
	// If this is the primary image, unset others
	if image.IsPrimary {
		for i := range e.Images {
			e.Images[i].IsPrimary = false
		}
	}
	e.Images = append(e.Images, image)
	e.UpdatedAt = time.Now()
}

// RemoveImage removes an image by URL.
func (e *Establishment) RemoveImage(imageURL string) {
	for i, img := range e.Images {
		if img.URL == imageURL {
			e.Images = append(e.Images[:i], e.Images[i+1:]...)
			e.UpdatedAt = time.Now()
			return
		}
	}
}

// SetPrimaryImage sets the primary image.
func (e *Establishment) SetPrimaryImage(imageURL string) {
	for i := range e.Images {
		e.Images[i].IsPrimary = e.Images[i].URL == imageURL
	}
	e.UpdatedAt = time.Now()
}

// GetPrimaryImage returns the primary image URL.
func (e *Establishment) GetPrimaryImage() string {
	for _, img := range e.Images {
		if img.IsPrimary {
			return img.URL
		}
	}
	if len(e.Images) > 0 {
		return e.Images[0].URL
	}
	return ""
}

// SetFeatures sets the establishment's features/amenities.
func (e *Establishment) SetFeatures(features []string) {
	e.Features = features
	e.UpdatedAt = time.Now()
}

// SetType sets the establishment type.
func (e *Establishment) SetType(t string) {
	e.Type = t
	e.UpdatedAt = time.Now()
}

// SetPriceRange sets the price range.
func (e *Establishment) SetPriceRange(priceRange int) {
	if priceRange < 1 {
		priceRange = 1
	}
	if priceRange > 4 {
		priceRange = 4
	}
	e.PriceRange = priceRange
	e.UpdatedAt = time.Now()
}

// SetAddress updates the address and location.
func (e *Establishment) SetAddress(address Address, location GeoLocation) {
	e.Address = address
	e.Location = location
	e.UpdatedAt = time.Now()
}

// Activate activates the establishment.
func (e *Establishment) Activate() {
	e.IsActive = true
	e.UpdatedAt = time.Now()
}

// Deactivate deactivates the establishment.
func (e *Establishment) Deactivate() {
	e.IsActive = false
	e.UpdatedAt = time.Now()
}

// IsOpenAt checks if the establishment is open at a given time.
func (e *Establishment) IsOpenAt(t time.Time) bool {
	// Check for exceptional closures
	for _, closure := range e.Closures {
		if closure.Date.Truncate(24 * time.Hour).Equal(t.Truncate(24 * time.Hour)) {
			return false
		}
	}

	// Check opening hours for the day of week
	dayOfWeek := int(t.Weekday())
	for _, oh := range e.OpeningHours {
		if oh.DayOfWeek == dayOfWeek && !oh.IsClosed {
			// Parse times and compare
			openTime, _ := time.Parse("15:04", oh.Open)
			closeTime, _ := time.Parse("15:04", oh.Close)
			currentTime, _ := time.Parse("15:04", t.Format("15:04"))

			if currentTime.After(openTime) && currentTime.Before(closeTime) {
				return true
			}
			// Handle midnight crossing (e.g., 22:00 - 02:00)
			if closeTime.Before(openTime) {
				if currentTime.After(openTime) || currentTime.Before(closeTime) {
					return true
				}
			}
		}
	}

	return false
}

// =============================================================================
// Opening Hour Value Object
// =============================================================================

// OpeningHour represents an opening period for a specific day.
type OpeningHour struct {
	DayOfWeek int    `json:"dayOfWeek" bson:"dayOfWeek"` // 0 = Sunday, 1 = Monday, ...
	Open      string `json:"open" bson:"open"`           // "09:00"
	Close     string `json:"close" bson:"close"`         // "23:00"
	IsClosed  bool   `json:"isClosed" bson:"isClosed"`
}

// NewOpeningHour creates a new opening hour.
func NewOpeningHour(dayOfWeek int, open, close string) OpeningHour {
	return OpeningHour{
		DayOfWeek: dayOfWeek,
		Open:      open,
		Close:     close,
		IsClosed:  false,
	}
}

// NewClosedDay creates a closed day entry.
func NewClosedDay(dayOfWeek int) OpeningHour {
	return OpeningHour{
		DayOfWeek: dayOfWeek,
		IsClosed:  true,
	}
}

// =============================================================================
// Closure Value Object
// =============================================================================

// Closure represents an exceptional closure date.
type Closure struct {
	Date   time.Time `json:"date" bson:"date"`
	Reason string    `json:"reason" bson:"reason"`
}

// NewClosure creates a new closure.
func NewClosure(date time.Time, reason string) Closure {
	return Closure{
		Date:   date.Truncate(24 * time.Hour),
		Reason: reason,
	}
}

// =============================================================================
// Image Value Object
// =============================================================================

// Image represents a media image with metadata.
type Image struct {
	URL       string `json:"url" bson:"url"`
	Alt       string `json:"alt" bson:"alt"`
	IsPrimary bool   `json:"isPrimary" bson:"isPrimary"`
	Order     int    `json:"order" bson:"order"`
}

// NewImage creates a new image.
func NewImage(url, alt string, isPrimary bool, order int) Image {
	return Image{
		URL:       url,
		Alt:       alt,
		IsPrimary: isPrimary,
		Order:     order,
	}
}

// =============================================================================
// Establishment Contact Value Object
// =============================================================================

// EstablishmentContact contains contact info for an establishment.
type EstablishmentContact struct {
	Phone   string `json:"phone,omitempty" bson:"phone,omitempty"`
	Email   string `json:"email,omitempty" bson:"email,omitempty"`
	Website string `json:"website,omitempty" bson:"website,omitempty"`
}

// NewEstablishmentContact creates a new establishment contact.
func NewEstablishmentContact(phone, email, website string) EstablishmentContact {
	return EstablishmentContact{
		Phone:   phone,
		Email:   email,
		Website: website,
	}
}

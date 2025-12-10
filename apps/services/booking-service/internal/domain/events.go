package domain

import "time"

// =============================================================================
// DOMAIN EVENTS
// =============================================================================

// OutingBooked is emitted when a new outing is created
type OutingBooked struct {
	OutingID        string    `json:"outing_id"`
	UserID          string    `json:"user_id"`
	OfferID         string    `json:"offer_id"`
	PartnerID       string    `json:"partner_id"`
	EstablishmentID string    `json:"establishment_id"`
	QRCode          string    `json:"qr_code"`
	ExpiresAt       time.Time `json:"expires_at"`
	Timestamp       time.Time `json:"timestamp"`
}

func (e OutingBooked) EventName() string     { return "outing.booked" }
func (e OutingBooked) OccurredAt() time.Time { return e.Timestamp }
func (e OutingBooked) AggregateID() string   { return e.OutingID }

// OutingCheckedIn is emitted when an outing is checked in
type OutingCheckedIn struct {
	OutingID        string    `json:"outing_id"`
	UserID          string    `json:"user_id"`
	OfferID         string    `json:"offer_id"`
	PartnerID       string    `json:"partner_id"`
	EstablishmentID string    `json:"establishment_id"`
	CheckedInBy     string    `json:"checked_in_by"`
	Method          string    `json:"method"`
	Timestamp       time.Time `json:"timestamp"`
}

func (e OutingCheckedIn) EventName() string     { return "outing.checked_in" }
func (e OutingCheckedIn) OccurredAt() time.Time { return e.Timestamp }
func (e OutingCheckedIn) AggregateID() string   { return e.OutingID }

// OutingCancelled is emitted when an outing is cancelled
type OutingCancelled struct {
	OutingID    string    `json:"outing_id"`
	UserID      string    `json:"user_id"`
	OfferID     string    `json:"offer_id"`
	PartnerID   string    `json:"partner_id"`
	CancelledBy string    `json:"cancelled_by"`
	Reason      string    `json:"reason"`
	Timestamp   time.Time `json:"timestamp"`
}

func (e OutingCancelled) EventName() string     { return "outing.cancelled" }
func (e OutingCancelled) OccurredAt() time.Time { return e.Timestamp }
func (e OutingCancelled) AggregateID() string   { return e.OutingID }

// OutingExpired is emitted when an outing expires
type OutingExpired struct {
	OutingID  string    `json:"outing_id"`
	UserID    string    `json:"user_id"`
	OfferID   string    `json:"offer_id"`
	PartnerID string    `json:"partner_id"`
	Timestamp time.Time `json:"timestamp"`
}

func (e OutingExpired) EventName() string     { return "outing.expired" }
func (e OutingExpired) OccurredAt() time.Time { return e.Timestamp }
func (e OutingExpired) AggregateID() string   { return e.OutingID }

package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// =============================================================================
// DOMAIN EVENTS
// =============================================================================

// OutingBooked is emitted when a new outing is created
type OutingBooked struct {
	ID              string    `json:"event_id"`
	OutingID        string    `json:"outing_id"`
	UserID          string    `json:"user_id"`
	OfferID         string    `json:"offer_id"`
	PartnerID       string    `json:"partner_id"`
	EstablishmentID string    `json:"establishment_id"`
	QRCode          string    `json:"qr_code"`
	ExpiresAt       time.Time `json:"expires_at"`
	Timestamp       time.Time `json:"timestamp"`
}

func NewOutingBookedEvent(outingID, userID, offerID, partnerID, establishmentID, qrCode string, expiresAt time.Time) OutingBooked {
	return OutingBooked{
		ID:              uuid.New().String(),
		OutingID:        outingID,
		UserID:          userID,
		OfferID:         offerID,
		PartnerID:       partnerID,
		EstablishmentID: establishmentID,
		QRCode:          qrCode,
		ExpiresAt:       expiresAt,
		Timestamp:       time.Now().UTC(),
	}
}

func (e OutingBooked) EventID() string          { return e.ID }
func (e OutingBooked) EventName() string        { return "outing.booked" }
func (e OutingBooked) OccurredAt() time.Time    { return e.Timestamp }
func (e OutingBooked) AggregateID() string      { return e.OutingID }
func (e OutingBooked) AggregateType() string    { return "Outing" }
func (e OutingBooked) Version() int             { return 1 }
func (e OutingBooked) Payload() ([]byte, error) { return json.Marshal(e) }

// OutingCheckedIn is emitted when an outing is checked in
type OutingCheckedIn struct {
	ID              string    `json:"event_id"`
	OutingID        string    `json:"outing_id"`
	UserID          string    `json:"user_id"`
	OfferID         string    `json:"offer_id"`
	PartnerID       string    `json:"partner_id"`
	EstablishmentID string    `json:"establishment_id"`
	CheckedInBy     string    `json:"checked_in_by"`
	Method          string    `json:"method"`
	Timestamp       time.Time `json:"timestamp"`
}

func NewOutingCheckedInEvent(outingID, userID, offerID, partnerID, establishmentID, checkedInBy, method string) OutingCheckedIn {
	return OutingCheckedIn{
		ID:              uuid.New().String(),
		OutingID:        outingID,
		UserID:          userID,
		OfferID:         offerID,
		PartnerID:       partnerID,
		EstablishmentID: establishmentID,
		CheckedInBy:     checkedInBy,
		Method:          method,
		Timestamp:       time.Now().UTC(),
	}
}

func (e OutingCheckedIn) EventID() string          { return e.ID }
func (e OutingCheckedIn) EventName() string        { return "outing.checked_in" }
func (e OutingCheckedIn) OccurredAt() time.Time    { return e.Timestamp }
func (e OutingCheckedIn) AggregateID() string      { return e.OutingID }
func (e OutingCheckedIn) AggregateType() string    { return "Outing" }
func (e OutingCheckedIn) Version() int             { return 1 }
func (e OutingCheckedIn) Payload() ([]byte, error) { return json.Marshal(e) }

// OutingCancelled is emitted when an outing is cancelled
type OutingCancelled struct {
	ID          string    `json:"event_id"`
	OutingID    string    `json:"outing_id"`
	UserID      string    `json:"user_id"`
	OfferID     string    `json:"offer_id"`
	PartnerID   string    `json:"partner_id"`
	CancelledBy string    `json:"cancelled_by"`
	Reason      string    `json:"reason"`
	Timestamp   time.Time `json:"timestamp"`
}

func NewOutingCancelledEvent(outingID, userID, offerID, partnerID, cancelledBy, reason string) OutingCancelled {
	return OutingCancelled{
		ID:          uuid.New().String(),
		OutingID:    outingID,
		UserID:      userID,
		OfferID:     offerID,
		PartnerID:   partnerID,
		CancelledBy: cancelledBy,
		Reason:      reason,
		Timestamp:   time.Now().UTC(),
	}
}

func (e OutingCancelled) EventID() string          { return e.ID }
func (e OutingCancelled) EventName() string        { return "outing.cancelled" }
func (e OutingCancelled) OccurredAt() time.Time    { return e.Timestamp }
func (e OutingCancelled) AggregateID() string      { return e.OutingID }
func (e OutingCancelled) AggregateType() string    { return "Outing" }
func (e OutingCancelled) Version() int             { return 1 }
func (e OutingCancelled) Payload() ([]byte, error) { return json.Marshal(e) }

// OutingExpired is emitted when an outing expires
type OutingExpired struct {
	ID        string    `json:"event_id"`
	OutingID  string    `json:"outing_id"`
	UserID    string    `json:"user_id"`
	OfferID   string    `json:"offer_id"`
	PartnerID string    `json:"partner_id"`
	Timestamp time.Time `json:"timestamp"`
}

func NewOutingExpiredEvent(outingID, userID, offerID, partnerID string) OutingExpired {
	return OutingExpired{
		ID:        uuid.New().String(),
		OutingID:  outingID,
		UserID:    userID,
		OfferID:   offerID,
		PartnerID: partnerID,
		Timestamp: time.Now().UTC(),
	}
}

func (e OutingExpired) EventID() string          { return e.ID }
func (e OutingExpired) EventName() string        { return "outing.expired" }
func (e OutingExpired) OccurredAt() time.Time    { return e.Timestamp }
func (e OutingExpired) AggregateID() string      { return e.OutingID }
func (e OutingExpired) AggregateType() string    { return "Outing" }
func (e OutingExpired) Version() int             { return 1 }
func (e OutingExpired) Payload() ([]byte, error) { return json.Marshal(e) }

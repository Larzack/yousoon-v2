// Package domain contains domain events for the Discovery service.
package domain

import "time"

// =============================================================================
// Offer Events
// =============================================================================

// OfferCreatedEvent is raised when an offer is created.
type OfferCreatedEvent struct {
	OfferID         OfferID         `json:"offerId"`
	PartnerID       PartnerID       `json:"partnerId"`
	EstablishmentID EstablishmentID `json:"establishmentId"`
	Title           string          `json:"title"`
	CategoryID      CategoryID      `json:"categoryId"`
	Timestamp       time.Time       `json:"timestamp"`
}

func (e OfferCreatedEvent) EventName() string     { return "discovery.offer.created" }
func (e OfferCreatedEvent) OccurredAt() time.Time { return e.Timestamp }
func (e OfferCreatedEvent) AggregateID() string   { return e.OfferID.String() }

// OfferSubmittedForReviewEvent is raised when an offer is submitted for moderation.
type OfferSubmittedForReviewEvent struct {
	OfferID   OfferID   `json:"offerId"`
	PartnerID PartnerID `json:"partnerId"`
	Timestamp time.Time `json:"timestamp"`
}

func (e OfferSubmittedForReviewEvent) EventName() string {
	return "discovery.offer.submitted_for_review"
}
func (e OfferSubmittedForReviewEvent) OccurredAt() time.Time { return e.Timestamp }
func (e OfferSubmittedForReviewEvent) AggregateID() string   { return e.OfferID.String() }

// OfferApprovedEvent is raised when an offer is approved by moderation.
type OfferApprovedEvent struct {
	OfferID    OfferID   `json:"offerId"`
	PartnerID  PartnerID `json:"partnerId"`
	ReviewerID string    `json:"reviewerId"`
	Timestamp  time.Time `json:"timestamp"`
}

func (e OfferApprovedEvent) EventName() string     { return "discovery.offer.approved" }
func (e OfferApprovedEvent) OccurredAt() time.Time { return e.Timestamp }
func (e OfferApprovedEvent) AggregateID() string   { return e.OfferID.String() }

// OfferRejectedEvent is raised when an offer is rejected by moderation.
type OfferRejectedEvent struct {
	OfferID   OfferID   `json:"offerId"`
	PartnerID PartnerID `json:"partnerId"`
	Reason    string    `json:"reason"`
	Timestamp time.Time `json:"timestamp"`
}

func (e OfferRejectedEvent) EventName() string     { return "discovery.offer.rejected" }
func (e OfferRejectedEvent) OccurredAt() time.Time { return e.Timestamp }
func (e OfferRejectedEvent) AggregateID() string   { return e.OfferID.String() }

// OfferPublishedEvent is raised when an offer is published (becomes active).
type OfferPublishedEvent struct {
	OfferID         OfferID         `json:"offerId"`
	PartnerID       PartnerID       `json:"partnerId"`
	EstablishmentID EstablishmentID `json:"establishmentId"`
	CategoryID      CategoryID      `json:"categoryId"`
	Location        GeoLocation     `json:"location"`
	Timestamp       time.Time       `json:"timestamp"`
}

func (e OfferPublishedEvent) EventName() string     { return "discovery.offer.published" }
func (e OfferPublishedEvent) OccurredAt() time.Time { return e.Timestamp }
func (e OfferPublishedEvent) AggregateID() string   { return e.OfferID.String() }

// OfferPausedEvent is raised when an offer is paused.
type OfferPausedEvent struct {
	OfferID   OfferID   `json:"offerId"`
	PartnerID PartnerID `json:"partnerId"`
	Timestamp time.Time `json:"timestamp"`
}

func (e OfferPausedEvent) EventName() string     { return "discovery.offer.paused" }
func (e OfferPausedEvent) OccurredAt() time.Time { return e.Timestamp }
func (e OfferPausedEvent) AggregateID() string   { return e.OfferID.String() }

// OfferResumedEvent is raised when an offer is resumed.
type OfferResumedEvent struct {
	OfferID   OfferID   `json:"offerId"`
	PartnerID PartnerID `json:"partnerId"`
	Timestamp time.Time `json:"timestamp"`
}

func (e OfferResumedEvent) EventName() string     { return "discovery.offer.resumed" }
func (e OfferResumedEvent) OccurredAt() time.Time { return e.Timestamp }
func (e OfferResumedEvent) AggregateID() string   { return e.OfferID.String() }

// OfferExpiredEvent is raised when an offer expires.
type OfferExpiredEvent struct {
	OfferID   OfferID   `json:"offerId"`
	PartnerID PartnerID `json:"partnerId"`
	Timestamp time.Time `json:"timestamp"`
}

func (e OfferExpiredEvent) EventName() string     { return "discovery.offer.expired" }
func (e OfferExpiredEvent) OccurredAt() time.Time { return e.Timestamp }
func (e OfferExpiredEvent) AggregateID() string   { return e.OfferID.String() }

// OfferArchivedEvent is raised when an offer is archived.
type OfferArchivedEvent struct {
	OfferID   OfferID   `json:"offerId"`
	PartnerID PartnerID `json:"partnerId"`
	Timestamp time.Time `json:"timestamp"`
}

func (e OfferArchivedEvent) EventName() string     { return "discovery.offer.archived" }
func (e OfferArchivedEvent) OccurredAt() time.Time { return e.Timestamp }
func (e OfferArchivedEvent) AggregateID() string   { return e.OfferID.String() }

// OfferBookedEvent is raised when an offer is booked (from Booking context).
type OfferBookedEvent struct {
	OfferID   OfferID   `json:"offerId"`
	UserID    UserID    `json:"userId"`
	OutingID  string    `json:"outingId"`
	Timestamp time.Time `json:"timestamp"`
}

func (e OfferBookedEvent) EventName() string     { return "discovery.offer.booked" }
func (e OfferBookedEvent) OccurredAt() time.Time { return e.Timestamp }
func (e OfferBookedEvent) AggregateID() string   { return e.OfferID.String() }

// OfferFavoritedEvent is raised when an offer is added to favorites.
type OfferFavoritedEvent struct {
	OfferID   OfferID   `json:"offerId"`
	UserID    UserID    `json:"userId"`
	Timestamp time.Time `json:"timestamp"`
}

func (e OfferFavoritedEvent) EventName() string     { return "discovery.offer.favorited" }
func (e OfferFavoritedEvent) OccurredAt() time.Time { return e.Timestamp }
func (e OfferFavoritedEvent) AggregateID() string   { return e.OfferID.String() }

// OfferUnfavoritedEvent is raised when an offer is removed from favorites.
type OfferUnfavoritedEvent struct {
	OfferID   OfferID   `json:"offerId"`
	UserID    UserID    `json:"userId"`
	Timestamp time.Time `json:"timestamp"`
}

func (e OfferUnfavoritedEvent) EventName() string     { return "discovery.offer.unfavorited" }
func (e OfferUnfavoritedEvent) OccurredAt() time.Time { return e.Timestamp }
func (e OfferUnfavoritedEvent) AggregateID() string   { return e.OfferID.String() }

// OfferViewedEvent is raised when an offer is viewed.
type OfferViewedEvent struct {
	OfferID   OfferID   `json:"offerId"`
	UserID    *UserID   `json:"userId,omitempty"` // nil if anonymous
	Timestamp time.Time `json:"timestamp"`
}

func (e OfferViewedEvent) EventName() string     { return "discovery.offer.viewed" }
func (e OfferViewedEvent) OccurredAt() time.Time { return e.Timestamp }
func (e OfferViewedEvent) AggregateID() string   { return e.OfferID.String() }

// =============================================================================
// Category Events
// =============================================================================

// CategoryCreatedEvent is raised when a category is created.
type CategoryCreatedEvent struct {
	CategoryID CategoryID `json:"categoryId"`
	Slug       string     `json:"slug"`
	Name       string     `json:"name"`
	Timestamp  time.Time  `json:"timestamp"`
}

func (e CategoryCreatedEvent) EventName() string     { return "discovery.category.created" }
func (e CategoryCreatedEvent) OccurredAt() time.Time { return e.Timestamp }
func (e CategoryCreatedEvent) AggregateID() string   { return e.CategoryID.String() }

// CategoryUpdatedEvent is raised when a category is updated.
type CategoryUpdatedEvent struct {
	CategoryID CategoryID `json:"categoryId"`
	Timestamp  time.Time  `json:"timestamp"`
}

func (e CategoryUpdatedEvent) EventName() string     { return "discovery.category.updated" }
func (e CategoryUpdatedEvent) OccurredAt() time.Time { return e.Timestamp }
func (e CategoryUpdatedEvent) AggregateID() string   { return e.CategoryID.String() }

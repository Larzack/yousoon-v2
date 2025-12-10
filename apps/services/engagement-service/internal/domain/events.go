package domain

import "time"

// =============================================================================
// DOMAIN EVENTS
// =============================================================================

// FavoriteAdded is emitted when a favorite is added
type FavoriteAdded struct {
	FavoriteID string    `json:"favorite_id"`
	UserID     string    `json:"user_id"`
	OfferID    string    `json:"offer_id"`
	Timestamp  time.Time `json:"timestamp"`
}

func (e FavoriteAdded) EventName() string     { return "favorite.added" }
func (e FavoriteAdded) OccurredAt() time.Time { return e.Timestamp }
func (e FavoriteAdded) AggregateID() string   { return e.FavoriteID }

// FavoriteRemoved is emitted when a favorite is removed
type FavoriteRemoved struct {
	FavoriteID string    `json:"favorite_id"`
	UserID     string    `json:"user_id"`
	OfferID    string    `json:"offer_id"`
	Timestamp  time.Time `json:"timestamp"`
}

func (e FavoriteRemoved) EventName() string     { return "favorite.removed" }
func (e FavoriteRemoved) OccurredAt() time.Time { return e.Timestamp }
func (e FavoriteRemoved) AggregateID() string   { return e.FavoriteID }

// ReviewSubmitted is emitted when a review is submitted
type ReviewSubmitted struct {
	ReviewID  string    `json:"review_id"`
	UserID    string    `json:"user_id"`
	OfferID   string    `json:"offer_id"`
	PartnerID string    `json:"partner_id"`
	Rating    int       `json:"rating"`
	Timestamp time.Time `json:"timestamp"`
}

func (e ReviewSubmitted) EventName() string     { return "review.submitted" }
func (e ReviewSubmitted) OccurredAt() time.Time { return e.Timestamp }
func (e ReviewSubmitted) AggregateID() string   { return e.ReviewID }

// ReviewApproved is emitted when a review is approved
type ReviewApproved struct {
	ReviewID    string    `json:"review_id"`
	ModeratorID string    `json:"moderator_id"`
	Timestamp   time.Time `json:"timestamp"`
}

func (e ReviewApproved) EventName() string     { return "review.approved" }
func (e ReviewApproved) OccurredAt() time.Time { return e.Timestamp }
func (e ReviewApproved) AggregateID() string   { return e.ReviewID }

// ReviewRejected is emitted when a review is rejected
type ReviewRejected struct {
	ReviewID    string    `json:"review_id"`
	ModeratorID string    `json:"moderator_id"`
	Reason      string    `json:"reason"`
	Timestamp   time.Time `json:"timestamp"`
}

func (e ReviewRejected) EventName() string     { return "review.rejected" }
func (e ReviewRejected) OccurredAt() time.Time { return e.Timestamp }
func (e ReviewRejected) AggregateID() string   { return e.ReviewID }

// ReviewReported is emitted when a review is reported
type ReviewReported struct {
	ReviewID   string    `json:"review_id"`
	ReporterID string    `json:"reporter_id"`
	Reason     string    `json:"reason"`
	Timestamp  time.Time `json:"timestamp"`
}

func (e ReviewReported) EventName() string     { return "review.reported" }
func (e ReviewReported) OccurredAt() time.Time { return e.Timestamp }
func (e ReviewReported) AggregateID() string   { return e.ReviewID }

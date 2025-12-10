package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// =============================================================================
// DOMAIN EVENTS
// =============================================================================

// FavoriteAdded is emitted when a favorite is added
type FavoriteAdded struct {
	ID         string    `json:"event_id"`
	FavoriteID string    `json:"favorite_id"`
	UserID     string    `json:"user_id"`
	OfferID    string    `json:"offer_id"`
	Timestamp  time.Time `json:"timestamp"`
}

func NewFavoriteAddedEvent(favoriteID, userID, offerID string) FavoriteAdded {
	return FavoriteAdded{
		ID:         uuid.New().String(),
		FavoriteID: favoriteID,
		UserID:     userID,
		OfferID:    offerID,
		Timestamp:  time.Now().UTC(),
	}
}

func (e FavoriteAdded) EventID() string          { return e.ID }
func (e FavoriteAdded) EventName() string        { return "favorite.added" }
func (e FavoriteAdded) OccurredAt() time.Time    { return e.Timestamp }
func (e FavoriteAdded) AggregateID() string      { return e.FavoriteID }
func (e FavoriteAdded) AggregateType() string    { return "Favorite" }
func (e FavoriteAdded) Version() int             { return 1 }
func (e FavoriteAdded) Payload() ([]byte, error) { return json.Marshal(e) }

// FavoriteRemoved is emitted when a favorite is removed
type FavoriteRemoved struct {
	ID         string    `json:"event_id"`
	FavoriteID string    `json:"favorite_id"`
	UserID     string    `json:"user_id"`
	OfferID    string    `json:"offer_id"`
	Timestamp  time.Time `json:"timestamp"`
}

func NewFavoriteRemovedEvent(favoriteID, userID, offerID string) FavoriteRemoved {
	return FavoriteRemoved{
		ID:         uuid.New().String(),
		FavoriteID: favoriteID,
		UserID:     userID,
		OfferID:    offerID,
		Timestamp:  time.Now().UTC(),
	}
}

func (e FavoriteRemoved) EventID() string          { return e.ID }
func (e FavoriteRemoved) EventName() string        { return "favorite.removed" }
func (e FavoriteRemoved) OccurredAt() time.Time    { return e.Timestamp }
func (e FavoriteRemoved) AggregateID() string      { return e.FavoriteID }
func (e FavoriteRemoved) AggregateType() string    { return "Favorite" }
func (e FavoriteRemoved) Version() int             { return 1 }
func (e FavoriteRemoved) Payload() ([]byte, error) { return json.Marshal(e) }

// ReviewSubmitted is emitted when a review is submitted
type ReviewSubmitted struct {
	ID        string    `json:"event_id"`
	ReviewID  string    `json:"review_id"`
	UserID    string    `json:"user_id"`
	OfferID   string    `json:"offer_id"`
	PartnerID string    `json:"partner_id"`
	Rating    int       `json:"rating"`
	Timestamp time.Time `json:"timestamp"`
}

func NewReviewSubmittedEvent(reviewID, userID, offerID, partnerID string, rating int) ReviewSubmitted {
	return ReviewSubmitted{
		ID:        uuid.New().String(),
		ReviewID:  reviewID,
		UserID:    userID,
		OfferID:   offerID,
		PartnerID: partnerID,
		Rating:    rating,
		Timestamp: time.Now().UTC(),
	}
}

func (e ReviewSubmitted) EventID() string          { return e.ID }
func (e ReviewSubmitted) EventName() string        { return "review.submitted" }
func (e ReviewSubmitted) OccurredAt() time.Time    { return e.Timestamp }
func (e ReviewSubmitted) AggregateID() string      { return e.ReviewID }
func (e ReviewSubmitted) AggregateType() string    { return "Review" }
func (e ReviewSubmitted) Version() int             { return 1 }
func (e ReviewSubmitted) Payload() ([]byte, error) { return json.Marshal(e) }

// ReviewApproved is emitted when a review is approved
type ReviewApproved struct {
	ID          string    `json:"event_id"`
	ReviewID    string    `json:"review_id"`
	ModeratorID string    `json:"moderator_id"`
	Timestamp   time.Time `json:"timestamp"`
}

func NewReviewApprovedEvent(reviewID, moderatorID string) ReviewApproved {
	return ReviewApproved{
		ID:          uuid.New().String(),
		ReviewID:    reviewID,
		ModeratorID: moderatorID,
		Timestamp:   time.Now().UTC(),
	}
}

func (e ReviewApproved) EventID() string          { return e.ID }
func (e ReviewApproved) EventName() string        { return "review.approved" }
func (e ReviewApproved) OccurredAt() time.Time    { return e.Timestamp }
func (e ReviewApproved) AggregateID() string      { return e.ReviewID }
func (e ReviewApproved) AggregateType() string    { return "Review" }
func (e ReviewApproved) Version() int             { return 1 }
func (e ReviewApproved) Payload() ([]byte, error) { return json.Marshal(e) }

// ReviewRejected is emitted when a review is rejected
type ReviewRejected struct {
	ID          string    `json:"event_id"`
	ReviewID    string    `json:"review_id"`
	ModeratorID string    `json:"moderator_id"`
	Reason      string    `json:"reason"`
	Timestamp   time.Time `json:"timestamp"`
}

func NewReviewRejectedEvent(reviewID, moderatorID, reason string) ReviewRejected {
	return ReviewRejected{
		ID:          uuid.New().String(),
		ReviewID:    reviewID,
		ModeratorID: moderatorID,
		Reason:      reason,
		Timestamp:   time.Now().UTC(),
	}
}

func (e ReviewRejected) EventID() string          { return e.ID }
func (e ReviewRejected) EventName() string        { return "review.rejected" }
func (e ReviewRejected) OccurredAt() time.Time    { return e.Timestamp }
func (e ReviewRejected) AggregateID() string      { return e.ReviewID }
func (e ReviewRejected) AggregateType() string    { return "Review" }
func (e ReviewRejected) Version() int             { return 1 }
func (e ReviewRejected) Payload() ([]byte, error) { return json.Marshal(e) }

// ReviewReported is emitted when a review is reported
type ReviewReported struct {
	ID         string    `json:"event_id"`
	ReviewID   string    `json:"review_id"`
	ReporterID string    `json:"reporter_id"`
	Reason     string    `json:"reason"`
	Timestamp  time.Time `json:"timestamp"`
}

func NewReviewReportedEvent(reviewID, reporterID, reason string) ReviewReported {
	return ReviewReported{
		ID:         uuid.New().String(),
		ReviewID:   reviewID,
		ReporterID: reporterID,
		Reason:     reason,
		Timestamp:  time.Now().UTC(),
	}
}

func (e ReviewReported) EventID() string          { return e.ID }
func (e ReviewReported) EventName() string        { return "review.reported" }
func (e ReviewReported) OccurredAt() time.Time    { return e.Timestamp }
func (e ReviewReported) AggregateID() string      { return e.ReviewID }
func (e ReviewReported) AggregateType() string    { return "Review" }
func (e ReviewReported) Version() int             { return 1 }
func (e ReviewReported) Payload() ([]byte, error) { return json.Marshal(e) }

package domain

import (
	"errors"
	"time"

	"github.com/yousoon/shared/domain"
)

// =============================================================================
// ERRORS
// =============================================================================

var (
	ErrFavoriteNotFound      = errors.New("favorite not found")
	ErrFavoriteAlreadyExists = errors.New("favorite already exists")
	ErrReviewNotFound        = errors.New("review not found")
	ErrReviewAlreadyExists   = errors.New("review already exists for this offer")
	ErrInvalidRating         = errors.New("rating must be between 1 and 5")
	ErrEmptyReviewContent    = errors.New("review content cannot be empty")
	ErrReviewAlreadyReported = errors.New("review already reported by this user")
	ErrCannotReviewOwnOffer  = errors.New("cannot review own offer")
	ErrMustHaveBooking       = errors.New("must have a booking to review")
)

// =============================================================================
// ENUMS
// =============================================================================

type ReviewStatus string

const (
	ReviewStatusPending  ReviewStatus = "pending"
	ReviewStatusApproved ReviewStatus = "approved"
	ReviewStatusRejected ReviewStatus = "rejected"
	ReviewStatusReported ReviewStatus = "reported"
)

func (s ReviewStatus) String() string {
	return string(s)
}

func (s ReviewStatus) IsValid() bool {
	switch s {
	case ReviewStatusPending, ReviewStatusApproved, ReviewStatusRejected, ReviewStatusReported:
		return true
	}
	return false
}

// =============================================================================
// AGGREGATE: Favorite
// =============================================================================

type Favorite struct {
	domain.AggregateRoot

	id      string
	userID  string
	offerID string

	// Denormalized data for quick display
	offerTitle    string
	offerImageURL string
	partnerName   string

	createdAt time.Time
}

func NewFavorite(userID, offerID, offerTitle, offerImageURL, partnerName string) (*Favorite, error) {
	now := time.Now()
	id := domain.NewID()

	favorite := &Favorite{
		id:            id,
		userID:        userID,
		offerID:       offerID,
		offerTitle:    offerTitle,
		offerImageURL: offerImageURL,
		partnerName:   partnerName,
		createdAt:     now,
	}

	favorite.AddDomainEvent(FavoriteAdded{
		FavoriteID: id,
		UserID:     userID,
		OfferID:    offerID,
		Timestamp:  now,
	})

	return favorite, nil
}

func ReconstructFavorite(id, userID, offerID, offerTitle, offerImageURL, partnerName string, createdAt time.Time) *Favorite {
	return &Favorite{
		id:            id,
		userID:        userID,
		offerID:       offerID,
		offerTitle:    offerTitle,
		offerImageURL: offerImageURL,
		partnerName:   partnerName,
		createdAt:     createdAt,
	}
}

func (f *Favorite) ID() string            { return f.id }
func (f *Favorite) UserID() string        { return f.userID }
func (f *Favorite) OfferID() string       { return f.offerID }
func (f *Favorite) OfferTitle() string    { return f.offerTitle }
func (f *Favorite) OfferImageURL() string { return f.offerImageURL }
func (f *Favorite) PartnerName() string   { return f.partnerName }
func (f *Favorite) CreatedAt() time.Time  { return f.createdAt }

// =============================================================================
// AGGREGATE: Review
// =============================================================================

type Review struct {
	domain.AggregateRoot

	id              string
	userID          string
	offerID         string
	partnerID       string
	establishmentID string
	bookingID       *string // Optional, for verified purchase

	// Content
	rating  int    // 1-5
	title   string // Optional
	content string
	images  []string

	// User info (denormalized)
	userFirstName string
	userAvatar    string

	// Offer info (denormalized)
	offerTitle string

	// Partner info (denormalized)
	partnerName string

	// Moderation
	status       ReviewStatus
	reports      []ReviewReport
	moderatedBy  *string
	moderatedAt  *time.Time
	rejectReason *string

	// Stats
	helpfulCount int

	// Flags
	isVerifiedPurchase bool

	createdAt time.Time
	updatedAt time.Time
}

type ReviewReport struct {
	UserID     string    `json:"user_id"`
	Reason     string    `json:"reason"`
	ReportedAt time.Time `json:"reported_at"`
}

func NewReview(
	userID, offerID, partnerID, establishmentID string,
	bookingID *string,
	rating int,
	title, content string,
	images []string,
	userFirstName, userAvatar string,
	offerTitle, partnerName string,
	isVerifiedPurchase bool,
) (*Review, error) {
	// Validate rating
	if rating < 1 || rating > 5 {
		return nil, ErrInvalidRating
	}

	// Validate content
	if content == "" {
		return nil, ErrEmptyReviewContent
	}

	now := time.Now()
	id := domain.NewID()

	review := &Review{
		id:                 id,
		userID:             userID,
		offerID:            offerID,
		partnerID:          partnerID,
		establishmentID:    establishmentID,
		bookingID:          bookingID,
		rating:             rating,
		title:              title,
		content:            content,
		images:             images,
		userFirstName:      userFirstName,
		userAvatar:         userAvatar,
		offerTitle:         offerTitle,
		partnerName:        partnerName,
		status:             ReviewStatusPending, // Needs moderation
		reports:            []ReviewReport{},
		helpfulCount:       0,
		isVerifiedPurchase: isVerifiedPurchase,
		createdAt:          now,
		updatedAt:          now,
	}

	review.AddDomainEvent(ReviewSubmitted{
		ReviewID:  id,
		UserID:    userID,
		OfferID:   offerID,
		PartnerID: partnerID,
		Rating:    rating,
		Timestamp: now,
	})

	return review, nil
}

func ReconstructReview(
	id, userID, offerID, partnerID, establishmentID string,
	bookingID *string,
	rating int,
	title, content string,
	images []string,
	userFirstName, userAvatar string,
	offerTitle, partnerName string,
	status ReviewStatus,
	reports []ReviewReport,
	moderatedBy *string,
	moderatedAt *time.Time,
	rejectReason *string,
	helpfulCount int,
	isVerifiedPurchase bool,
	createdAt, updatedAt time.Time,
) *Review {
	return &Review{
		id:                 id,
		userID:             userID,
		offerID:            offerID,
		partnerID:          partnerID,
		establishmentID:    establishmentID,
		bookingID:          bookingID,
		rating:             rating,
		title:              title,
		content:            content,
		images:             images,
		userFirstName:      userFirstName,
		userAvatar:         userAvatar,
		offerTitle:         offerTitle,
		partnerName:        partnerName,
		status:             status,
		reports:            reports,
		moderatedBy:        moderatedBy,
		moderatedAt:        moderatedAt,
		rejectReason:       rejectReason,
		helpfulCount:       helpfulCount,
		isVerifiedPurchase: isVerifiedPurchase,
		createdAt:          createdAt,
		updatedAt:          updatedAt,
	}
}

// Getters
func (r *Review) ID() string               { return r.id }
func (r *Review) UserID() string           { return r.userID }
func (r *Review) OfferID() string          { return r.offerID }
func (r *Review) PartnerID() string        { return r.partnerID }
func (r *Review) EstablishmentID() string  { return r.establishmentID }
func (r *Review) BookingID() *string       { return r.bookingID }
func (r *Review) Rating() int              { return r.rating }
func (r *Review) Title() string            { return r.title }
func (r *Review) Content() string          { return r.content }
func (r *Review) Images() []string         { return r.images }
func (r *Review) UserFirstName() string    { return r.userFirstName }
func (r *Review) UserAvatar() string       { return r.userAvatar }
func (r *Review) OfferTitle() string       { return r.offerTitle }
func (r *Review) PartnerName() string      { return r.partnerName }
func (r *Review) Status() ReviewStatus     { return r.status }
func (r *Review) Reports() []ReviewReport  { return r.reports }
func (r *Review) ModeratedBy() *string     { return r.moderatedBy }
func (r *Review) ModeratedAt() *time.Time  { return r.moderatedAt }
func (r *Review) RejectReason() *string    { return r.rejectReason }
func (r *Review) HelpfulCount() int        { return r.helpfulCount }
func (r *Review) IsVerifiedPurchase() bool { return r.isVerifiedPurchase }
func (r *Review) CreatedAt() time.Time     { return r.createdAt }
func (r *Review) UpdatedAt() time.Time     { return r.updatedAt }

// Commands
func (r *Review) Approve(moderatorID string) error {
	now := time.Now()
	r.status = ReviewStatusApproved
	r.moderatedBy = &moderatorID
	r.moderatedAt = &now
	r.updatedAt = now

	r.AddDomainEvent(ReviewApproved{
		ReviewID:    r.id,
		ModeratorID: moderatorID,
		Timestamp:   now,
	})

	return nil
}

func (r *Review) Reject(moderatorID, reason string) error {
	now := time.Now()
	r.status = ReviewStatusRejected
	r.moderatedBy = &moderatorID
	r.moderatedAt = &now
	r.rejectReason = &reason
	r.updatedAt = now

	r.AddDomainEvent(ReviewRejected{
		ReviewID:    r.id,
		ModeratorID: moderatorID,
		Reason:      reason,
		Timestamp:   now,
	})

	return nil
}

func (r *Review) Report(userID, reason string) error {
	// Check if already reported by this user
	for _, report := range r.reports {
		if report.UserID == userID {
			return ErrReviewAlreadyReported
		}
	}

	r.reports = append(r.reports, ReviewReport{
		UserID:     userID,
		Reason:     reason,
		ReportedAt: time.Now(),
	})

	// Auto-flag if too many reports
	if len(r.reports) >= 3 {
		r.status = ReviewStatusReported
	}

	r.updatedAt = time.Now()

	r.AddDomainEvent(ReviewReported{
		ReviewID:   r.id,
		ReporterID: userID,
		Reason:     reason,
		Timestamp:  time.Now(),
	})

	return nil
}

func (r *Review) MarkHelpful() {
	r.helpfulCount++
	r.updatedAt = time.Now()
}

func (r *Review) Update(rating int, title, content string, images []string) error {
	if rating < 1 || rating > 5 {
		return ErrInvalidRating
	}
	if content == "" {
		return ErrEmptyReviewContent
	}

	r.rating = rating
	r.title = title
	r.content = content
	r.images = images
	r.status = ReviewStatusPending // Re-moderation required
	r.updatedAt = time.Now()

	return nil
}

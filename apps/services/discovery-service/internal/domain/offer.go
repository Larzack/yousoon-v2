// Package domain contains the domain layer for the Discovery service.
// This file defines the Offer aggregate root.
package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// OfferStatus represents the status of an offer.
type OfferStatus string

const (
	OfferStatusDraft    OfferStatus = "draft"
	OfferStatusPending  OfferStatus = "pending"
	OfferStatusActive   OfferStatus = "active"
	OfferStatusPaused   OfferStatus = "paused"
	OfferStatusExpired  OfferStatus = "expired"
	OfferStatusArchived OfferStatus = "archived"
)

// ModerationStatus represents the moderation status of an offer.
type ModerationStatus string

const (
	ModerationStatusPending  ModerationStatus = "pending"
	ModerationStatusApproved ModerationStatus = "approved"
	ModerationStatusRejected ModerationStatus = "rejected"
)

// Offer is the aggregate root for offers/deals.
type Offer struct {
	id              OfferID
	partnerID       PartnerID
	establishmentID EstablishmentID

	// Core information
	title            string
	description      string
	shortDescription string
	categoryID       CategoryID
	tags             []string

	// Discount
	discount Discount

	// Conditions
	conditions         []Condition
	termsAndConditions string

	// Validity
	validity Validity

	// Schedule
	schedule Schedule

	// Quota
	quota Quota

	// Media
	images []OfferImage

	// Denormalized data for performance
	partnerSnapshot       PartnerSnapshot
	establishmentSnapshot EstablishmentSnapshot

	// Statistics
	stats OfferStats

	// Status
	status     OfferStatus
	moderation Moderation

	// Timestamps
	createdAt   time.Time
	updatedAt   time.Time
	publishedAt *time.Time
	deletedAt   *time.Time

	// Domain events
	events []interface{}
}

// NewOffer creates a new Offer.
func NewOffer(
	partnerID PartnerID,
	establishmentID EstablishmentID,
	title string,
	description string,
	categoryID CategoryID,
	discount Discount,
	validity Validity,
) (*Offer, error) {
	if title == "" {
		return nil, errors.New("title is required")
	}
	if partnerID == "" {
		return nil, errors.New("partner ID is required")
	}
	if establishmentID == "" {
		return nil, errors.New("establishment ID is required")
	}
	if categoryID == "" {
		return nil, errors.New("category ID is required")
	}

	now := time.Now()
	offer := &Offer{
		id:              OfferID(uuid.New().String()),
		partnerID:       partnerID,
		establishmentID: establishmentID,
		title:           title,
		description:     description,
		categoryID:      categoryID,
		discount:        discount,
		validity:        validity,
		schedule:        NewAllDaySchedule(),
		quota:           NewUnlimitedQuota(),
		status:          OfferStatusDraft,
		moderation: Moderation{
			Status: ModerationStatusPending,
		},
		stats:     OfferStats{},
		createdAt: now,
		updatedAt: now,
		events:    make([]interface{}, 0),
	}

	offer.events = append(offer.events, OfferCreatedEvent{
		OfferID:         offer.id,
		PartnerID:       offer.partnerID,
		EstablishmentID: offer.establishmentID,
		Title:           offer.title,
		CategoryID:      offer.categoryID,
		Timestamp:       now,
	})

	return offer, nil
}

// Getters

func (o *Offer) ID() OfferID                                  { return o.id }
func (o *Offer) PartnerID() PartnerID                         { return o.partnerID }
func (o *Offer) EstablishmentID() EstablishmentID             { return o.establishmentID }
func (o *Offer) Title() string                                { return o.title }
func (o *Offer) Description() string                          { return o.description }
func (o *Offer) ShortDescription() string                     { return o.shortDescription }
func (o *Offer) CategoryID() CategoryID                       { return o.categoryID }
func (o *Offer) Tags() []string                               { return o.tags }
func (o *Offer) Discount() Discount                           { return o.discount }
func (o *Offer) Conditions() []Condition                      { return o.conditions }
func (o *Offer) TermsAndConditions() string                   { return o.termsAndConditions }
func (o *Offer) Validity() Validity                           { return o.validity }
func (o *Offer) Schedule() Schedule                           { return o.schedule }
func (o *Offer) Quota() Quota                                 { return o.quota }
func (o *Offer) Images() []OfferImage                         { return o.images }
func (o *Offer) PartnerSnapshot() PartnerSnapshot             { return o.partnerSnapshot }
func (o *Offer) EstablishmentSnapshot() EstablishmentSnapshot { return o.establishmentSnapshot }
func (o *Offer) Stats() OfferStats                            { return o.stats }
func (o *Offer) Status() OfferStatus                          { return o.status }
func (o *Offer) Moderation() Moderation                       { return o.moderation }
func (o *Offer) CreatedAt() time.Time                         { return o.createdAt }
func (o *Offer) UpdatedAt() time.Time                         { return o.updatedAt }
func (o *Offer) PublishedAt() *time.Time                      { return o.publishedAt }
func (o *Offer) DeletedAt() *time.Time                        { return o.deletedAt }
func (o *Offer) Events() []interface{}                        { return o.events }
func (o *Offer) ClearEvents()                                 { o.events = make([]interface{}, 0) }

// IsActive checks if the offer is currently active and bookable.
func (o *Offer) IsActive() bool {
	if o.status != OfferStatusActive {
		return false
	}
	if o.validity.IsExpired() {
		return false
	}
	if !o.validity.HasStarted() {
		return false
	}
	return true
}

// CanBeBooked checks if the offer can be booked right now.
func (o *Offer) CanBeBooked() error {
	if !o.IsActive() {
		return ErrOfferNotActive
	}
	if o.quota.IsExhausted() {
		return ErrOfferFullyBooked
	}
	if !o.schedule.IsAvailableNow() {
		return ErrOfferNotAvailableNow
	}
	return nil
}

// CanUserBook checks if a specific user can book this offer.
func (o *Offer) CanUserBook(userBookingCount int) error {
	if err := o.CanBeBooked(); err != nil {
		return err
	}
	if o.quota.PerUser != nil && userBookingCount >= *o.quota.PerUser {
		return ErrUserQuotaExceeded
	}
	return nil
}

// =============================================================================
// Commands (State Changes)
// =============================================================================

// UpdateBasicInfo updates the basic information of the offer.
func (o *Offer) UpdateBasicInfo(title, description, shortDescription string) error {
	if title == "" {
		return errors.New("title is required")
	}
	o.title = title
	o.description = description
	o.shortDescription = shortDescription
	o.updatedAt = time.Now()
	return nil
}

// UpdateCategory updates the category of the offer.
func (o *Offer) UpdateCategory(categoryID CategoryID) error {
	if categoryID == "" {
		return errors.New("category ID is required")
	}
	o.categoryID = categoryID
	o.updatedAt = time.Now()
	return nil
}

// UpdateTags updates the tags of the offer.
func (o *Offer) UpdateTags(tags []string) {
	o.tags = tags
	o.updatedAt = time.Now()
}

// UpdateDiscount updates the discount of the offer.
func (o *Offer) UpdateDiscount(discount Discount) error {
	if err := discount.Validate(); err != nil {
		return err
	}
	o.discount = discount
	o.updatedAt = time.Now()
	return nil
}

// UpdateConditions updates the conditions of the offer.
func (o *Offer) UpdateConditions(conditions []Condition, terms string) {
	o.conditions = conditions
	o.termsAndConditions = terms
	o.updatedAt = time.Now()
}

// UpdateValidity updates the validity period of the offer.
func (o *Offer) UpdateValidity(validity Validity) error {
	if err := validity.Validate(); err != nil {
		return err
	}
	o.validity = validity
	o.updatedAt = time.Now()
	return nil
}

// UpdateSchedule updates the schedule of the offer.
func (o *Offer) UpdateSchedule(schedule Schedule) {
	o.schedule = schedule
	o.updatedAt = time.Now()
}

// UpdateQuota updates the quota of the offer.
func (o *Offer) UpdateQuota(quota Quota) {
	o.quota = quota
	o.updatedAt = time.Now()
}

// AddImage adds an image to the offer.
func (o *Offer) AddImage(image OfferImage) {
	// If this is the first image or marked as primary, set it as primary
	if len(o.images) == 0 || image.IsPrimary {
		// Unset other primary images
		for i := range o.images {
			o.images[i].IsPrimary = false
		}
		image.IsPrimary = true
	}
	image.Order = len(o.images)
	o.images = append(o.images, image)
	o.updatedAt = time.Now()
}

// RemoveImage removes an image from the offer.
func (o *Offer) RemoveImage(url string) {
	newImages := make([]OfferImage, 0, len(o.images))
	for _, img := range o.images {
		if img.URL != url {
			newImages = append(newImages, img)
		}
	}
	o.images = newImages
	// Reorder
	for i := range o.images {
		o.images[i].Order = i
	}
	o.updatedAt = time.Now()
}

// SetPartnerSnapshot sets the denormalized partner data.
func (o *Offer) SetPartnerSnapshot(snapshot PartnerSnapshot) {
	o.partnerSnapshot = snapshot
	o.updatedAt = time.Now()
}

// SetEstablishmentSnapshot sets the denormalized establishment data.
func (o *Offer) SetEstablishmentSnapshot(snapshot EstablishmentSnapshot) {
	o.establishmentSnapshot = snapshot
	o.updatedAt = time.Now()
}

// =============================================================================
// Status Transitions
// =============================================================================

// SubmitForReview submits the offer for moderation review.
func (o *Offer) SubmitForReview() error {
	if o.status != OfferStatusDraft {
		return ErrInvalidStatusTransition
	}

	o.status = OfferStatusPending
	o.moderation = Moderation{
		Status: ModerationStatusPending,
	}
	o.updatedAt = time.Now()

	o.events = append(o.events, OfferSubmittedForReviewEvent{
		OfferID:   o.id,
		PartnerID: o.partnerID,
		Timestamp: time.Now(),
	})

	return nil
}

// Approve approves the offer (admin action).
func (o *Offer) Approve(reviewerID string) error {
	if o.status != OfferStatusPending {
		return ErrInvalidStatusTransition
	}

	now := time.Now()
	o.moderation = Moderation{
		Status:     ModerationStatusApproved,
		ReviewerID: &reviewerID,
		ReviewedAt: &now,
	}
	o.updatedAt = now

	return nil
}

// Reject rejects the offer (admin action).
func (o *Offer) Reject(reviewerID, reason string) error {
	if o.status != OfferStatusPending {
		return ErrInvalidStatusTransition
	}

	now := time.Now()
	o.status = OfferStatusDraft // Back to draft for editing
	o.moderation = Moderation{
		Status:     ModerationStatusRejected,
		ReviewerID: &reviewerID,
		ReviewedAt: &now,
		Comment:    &reason,
	}
	o.updatedAt = now

	o.events = append(o.events, OfferRejectedEvent{
		OfferID:   o.id,
		PartnerID: o.partnerID,
		Reason:    reason,
		Timestamp: now,
	})

	return nil
}

// Publish publishes the offer (makes it active).
func (o *Offer) Publish() error {
	if o.status != OfferStatusPending && o.status != OfferStatusPaused {
		return ErrInvalidStatusTransition
	}
	if o.moderation.Status != ModerationStatusApproved {
		return ErrOfferNotApproved
	}

	now := time.Now()
	o.status = OfferStatusActive
	o.publishedAt = &now
	o.updatedAt = now

	o.events = append(o.events, OfferPublishedEvent{
		OfferID:         o.id,
		PartnerID:       o.partnerID,
		EstablishmentID: o.establishmentID,
		CategoryID:      o.categoryID,
		Location:        o.establishmentSnapshot.Location,
		Timestamp:       now,
	})

	return nil
}

// Pause pauses the offer temporarily.
func (o *Offer) Pause() error {
	if o.status != OfferStatusActive {
		return ErrInvalidStatusTransition
	}

	o.status = OfferStatusPaused
	o.updatedAt = time.Now()

	o.events = append(o.events, OfferPausedEvent{
		OfferID:   o.id,
		PartnerID: o.partnerID,
		Timestamp: time.Now(),
	})

	return nil
}

// Resume resumes a paused offer.
func (o *Offer) Resume() error {
	if o.status != OfferStatusPaused {
		return ErrInvalidStatusTransition
	}

	o.status = OfferStatusActive
	o.updatedAt = time.Now()

	return nil
}

// Expire expires the offer (called by system job).
func (o *Offer) Expire() error {
	if o.status == OfferStatusExpired || o.status == OfferStatusArchived {
		return nil // Already expired or archived
	}

	o.status = OfferStatusExpired
	o.updatedAt = time.Now()

	o.events = append(o.events, OfferExpiredEvent{
		OfferID:   o.id,
		PartnerID: o.partnerID,
		Timestamp: time.Now(),
	})

	return nil
}

// Archive archives the offer.
func (o *Offer) Archive() error {
	if o.status == OfferStatusArchived {
		return nil
	}

	o.status = OfferStatusArchived
	o.updatedAt = time.Now()

	return nil
}

// Delete soft-deletes the offer.
func (o *Offer) Delete() error {
	now := time.Now()
	o.deletedAt = &now
	o.updatedAt = now
	return nil
}

// =============================================================================
// Statistics
// =============================================================================

// IncrementViews increments the view count.
func (o *Offer) IncrementViews() {
	o.stats.Views++
}

// IncrementBookings increments the booking count and quota used.
func (o *Offer) IncrementBookings() error {
	if o.quota.IsExhausted() {
		return ErrOfferFullyBooked
	}
	o.stats.Bookings++
	o.quota.Used++
	return nil
}

// IncrementCheckins increments the check-in count.
func (o *Offer) IncrementCheckins() {
	o.stats.Checkins++
}

// IncrementFavorites increments the favorites count.
func (o *Offer) IncrementFavorites() {
	o.stats.Favorites++
}

// DecrementFavorites decrements the favorites count.
func (o *Offer) DecrementFavorites() {
	if o.stats.Favorites > 0 {
		o.stats.Favorites--
	}
}

// =============================================================================
// Helper Methods
// =============================================================================

// GetPrimaryImage returns the primary image URL.
func (o *Offer) GetPrimaryImage() string {
	for _, img := range o.images {
		if img.IsPrimary {
			return img.URL
		}
	}
	if len(o.images) > 0 {
		return o.images[0].URL
	}
	return ""
}

// CalculateFinalPrice calculates the final price after discount.
func (o *Offer) CalculateFinalPrice(originalPrice int64) int64 {
	return o.discount.Apply(originalPrice)
}

// ToSnapshot creates an immutable snapshot of the offer for bookings.
func (o *Offer) ToSnapshot() OfferSnapshot {
	return OfferSnapshot{
		ID:              o.id,
		PartnerID:       o.partnerID,
		EstablishmentID: o.establishmentID,
		Title:           o.title,
		Description:     o.shortDescription,
		Discount:        o.discount,
		CategoryName:    "", // To be filled by caller
		Location:        o.establishmentSnapshot.Location,
		CapturedAt:      time.Now(),
	}
}

// =============================================================================
// Reconstruction
// =============================================================================

// ReconstructOffer reconstructs an offer from persistence.
func ReconstructOffer(
	id OfferID,
	partnerID PartnerID,
	establishmentID EstablishmentID,
	title string,
	description string,
	shortDescription string,
	categoryID CategoryID,
	tags []string,
	discount Discount,
	conditions []Condition,
	termsAndConditions string,
	validity Validity,
	schedule Schedule,
	quota Quota,
	images []OfferImage,
	partnerSnapshot PartnerSnapshot,
	establishmentSnapshot EstablishmentSnapshot,
	stats OfferStats,
	status OfferStatus,
	moderation Moderation,
	createdAt time.Time,
	updatedAt time.Time,
	publishedAt *time.Time,
	deletedAt *time.Time,
) *Offer {
	return &Offer{
		id:                    id,
		partnerID:             partnerID,
		establishmentID:       establishmentID,
		title:                 title,
		description:           description,
		shortDescription:      shortDescription,
		categoryID:            categoryID,
		tags:                  tags,
		discount:              discount,
		conditions:            conditions,
		termsAndConditions:    termsAndConditions,
		validity:              validity,
		schedule:              schedule,
		quota:                 quota,
		images:                images,
		partnerSnapshot:       partnerSnapshot,
		establishmentSnapshot: establishmentSnapshot,
		stats:                 stats,
		status:                status,
		moderation:            moderation,
		createdAt:             createdAt,
		updatedAt:             updatedAt,
		publishedAt:           publishedAt,
		deletedAt:             deletedAt,
		events:                make([]interface{}, 0),
	}
}

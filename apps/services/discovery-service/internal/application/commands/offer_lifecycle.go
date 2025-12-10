// Package commands contains command handlers for offer lifecycle management.
package commands

import (
	"context"

	"github.com/yousoon/services/discovery/internal/domain"
)

// =============================================================================
// Submit Offer For Review Command
// =============================================================================

// SubmitOfferForReviewCommand submits an offer for moderation review.
type SubmitOfferForReviewCommand struct {
	OfferID string
}

// SubmitOfferForReviewHandler handles the submit for review command.
type SubmitOfferForReviewHandler struct {
	offerRepo domain.OfferRepository
}

// NewSubmitOfferForReviewHandler creates a new handler.
func NewSubmitOfferForReviewHandler(offerRepo domain.OfferRepository) *SubmitOfferForReviewHandler {
	return &SubmitOfferForReviewHandler{
		offerRepo: offerRepo,
	}
}

// Handle executes the submit for review command.
func (h *SubmitOfferForReviewHandler) Handle(ctx context.Context, cmd SubmitOfferForReviewCommand) (*domain.Offer, error) {
	offer, err := h.offerRepo.FindByID(ctx, domain.OfferID(cmd.OfferID))
	if err != nil {
		return nil, err
	}
	if offer == nil {
		return nil, domain.ErrOfferNotFound
	}

	if err := offer.SubmitForReview(); err != nil {
		return nil, err
	}

	if err := h.offerRepo.Save(ctx, offer); err != nil {
		return nil, err
	}

	return offer, nil
}

// =============================================================================
// Approve Offer Command (Admin)
// =============================================================================

// ApproveOfferCommand approves an offer.
type ApproveOfferCommand struct {
	OfferID    string
	ReviewerID string
}

// ApproveOfferHandler handles the approve offer command.
type ApproveOfferHandler struct {
	offerRepo domain.OfferRepository
}

// NewApproveOfferHandler creates a new handler.
func NewApproveOfferHandler(offerRepo domain.OfferRepository) *ApproveOfferHandler {
	return &ApproveOfferHandler{
		offerRepo: offerRepo,
	}
}

// Handle executes the approve offer command.
func (h *ApproveOfferHandler) Handle(ctx context.Context, cmd ApproveOfferCommand) (*domain.Offer, error) {
	offer, err := h.offerRepo.FindByID(ctx, domain.OfferID(cmd.OfferID))
	if err != nil {
		return nil, err
	}
	if offer == nil {
		return nil, domain.ErrOfferNotFound
	}

	if err := offer.Approve(cmd.ReviewerID); err != nil {
		return nil, err
	}

	if err := h.offerRepo.Save(ctx, offer); err != nil {
		return nil, err
	}

	return offer, nil
}

// =============================================================================
// Reject Offer Command (Admin)
// =============================================================================

// RejectOfferCommand rejects an offer.
type RejectOfferCommand struct {
	OfferID    string
	ReviewerID string
	Reason     string
}

// RejectOfferHandler handles the reject offer command.
type RejectOfferHandler struct {
	offerRepo domain.OfferRepository
}

// NewRejectOfferHandler creates a new handler.
func NewRejectOfferHandler(offerRepo domain.OfferRepository) *RejectOfferHandler {
	return &RejectOfferHandler{
		offerRepo: offerRepo,
	}
}

// Handle executes the reject offer command.
func (h *RejectOfferHandler) Handle(ctx context.Context, cmd RejectOfferCommand) (*domain.Offer, error) {
	offer, err := h.offerRepo.FindByID(ctx, domain.OfferID(cmd.OfferID))
	if err != nil {
		return nil, err
	}
	if offer == nil {
		return nil, domain.ErrOfferNotFound
	}

	if err := offer.Reject(cmd.ReviewerID, cmd.Reason); err != nil {
		return nil, err
	}

	if err := h.offerRepo.Save(ctx, offer); err != nil {
		return nil, err
	}

	return offer, nil
}

// =============================================================================
// Publish Offer Command
// =============================================================================

// PublishOfferCommand publishes an offer.
type PublishOfferCommand struct {
	OfferID string
}

// PublishOfferHandler handles the publish offer command.
type PublishOfferHandler struct {
	offerRepo domain.OfferRepository
}

// NewPublishOfferHandler creates a new handler.
func NewPublishOfferHandler(offerRepo domain.OfferRepository) *PublishOfferHandler {
	return &PublishOfferHandler{
		offerRepo: offerRepo,
	}
}

// Handle executes the publish offer command.
func (h *PublishOfferHandler) Handle(ctx context.Context, cmd PublishOfferCommand) (*domain.Offer, error) {
	offer, err := h.offerRepo.FindByID(ctx, domain.OfferID(cmd.OfferID))
	if err != nil {
		return nil, err
	}
	if offer == nil {
		return nil, domain.ErrOfferNotFound
	}

	if err := offer.Publish(); err != nil {
		return nil, err
	}

	if err := h.offerRepo.Save(ctx, offer); err != nil {
		return nil, err
	}

	return offer, nil
}

// =============================================================================
// Pause Offer Command
// =============================================================================

// PauseOfferCommand pauses an active offer.
type PauseOfferCommand struct {
	OfferID string
}

// PauseOfferHandler handles the pause offer command.
type PauseOfferHandler struct {
	offerRepo domain.OfferRepository
}

// NewPauseOfferHandler creates a new handler.
func NewPauseOfferHandler(offerRepo domain.OfferRepository) *PauseOfferHandler {
	return &PauseOfferHandler{
		offerRepo: offerRepo,
	}
}

// Handle executes the pause offer command.
func (h *PauseOfferHandler) Handle(ctx context.Context, cmd PauseOfferCommand) (*domain.Offer, error) {
	offer, err := h.offerRepo.FindByID(ctx, domain.OfferID(cmd.OfferID))
	if err != nil {
		return nil, err
	}
	if offer == nil {
		return nil, domain.ErrOfferNotFound
	}

	if err := offer.Pause(); err != nil {
		return nil, err
	}

	if err := h.offerRepo.Save(ctx, offer); err != nil {
		return nil, err
	}

	return offer, nil
}

// =============================================================================
// Resume Offer Command
// =============================================================================

// ResumeOfferCommand resumes a paused offer.
type ResumeOfferCommand struct {
	OfferID string
}

// ResumeOfferHandler handles the resume offer command.
type ResumeOfferHandler struct {
	offerRepo domain.OfferRepository
}

// NewResumeOfferHandler creates a new handler.
func NewResumeOfferHandler(offerRepo domain.OfferRepository) *ResumeOfferHandler {
	return &ResumeOfferHandler{
		offerRepo: offerRepo,
	}
}

// Handle executes the resume offer command.
func (h *ResumeOfferHandler) Handle(ctx context.Context, cmd ResumeOfferCommand) (*domain.Offer, error) {
	offer, err := h.offerRepo.FindByID(ctx, domain.OfferID(cmd.OfferID))
	if err != nil {
		return nil, err
	}
	if offer == nil {
		return nil, domain.ErrOfferNotFound
	}

	if err := offer.Resume(); err != nil {
		return nil, err
	}

	if err := h.offerRepo.Save(ctx, offer); err != nil {
		return nil, err
	}

	return offer, nil
}

// =============================================================================
// Archive Offer Command
// =============================================================================

// ArchiveOfferCommand archives an offer.
type ArchiveOfferCommand struct {
	OfferID string
}

// ArchiveOfferHandler handles the archive offer command.
type ArchiveOfferHandler struct {
	offerRepo domain.OfferRepository
}

// NewArchiveOfferHandler creates a new handler.
func NewArchiveOfferHandler(offerRepo domain.OfferRepository) *ArchiveOfferHandler {
	return &ArchiveOfferHandler{
		offerRepo: offerRepo,
	}
}

// Handle executes the archive offer command.
func (h *ArchiveOfferHandler) Handle(ctx context.Context, cmd ArchiveOfferCommand) (*domain.Offer, error) {
	offer, err := h.offerRepo.FindByID(ctx, domain.OfferID(cmd.OfferID))
	if err != nil {
		return nil, err
	}
	if offer == nil {
		return nil, domain.ErrOfferNotFound
	}

	if err := offer.Archive(); err != nil {
		return nil, err
	}

	if err := h.offerRepo.Save(ctx, offer); err != nil {
		return nil, err
	}

	return offer, nil
}

// =============================================================================
// Delete Offer Command
// =============================================================================

// DeleteOfferCommand soft-deletes an offer.
type DeleteOfferCommand struct {
	OfferID string
}

// DeleteOfferHandler handles the delete offer command.
type DeleteOfferHandler struct {
	offerRepo domain.OfferRepository
}

// NewDeleteOfferHandler creates a new handler.
func NewDeleteOfferHandler(offerRepo domain.OfferRepository) *DeleteOfferHandler {
	return &DeleteOfferHandler{
		offerRepo: offerRepo,
	}
}

// Handle executes the delete offer command.
func (h *DeleteOfferHandler) Handle(ctx context.Context, cmd DeleteOfferCommand) error {
	offer, err := h.offerRepo.FindByID(ctx, domain.OfferID(cmd.OfferID))
	if err != nil {
		return err
	}
	if offer == nil {
		return domain.ErrOfferNotFound
	}

	if err := offer.Delete(); err != nil {
		return err
	}

	return h.offerRepo.Save(ctx, offer)
}

// =============================================================================
// Add Offer Image Command
// =============================================================================

// AddOfferImageCommand adds an image to an offer.
type AddOfferImageCommand struct {
	OfferID   string
	URL       string
	Alt       string
	IsPrimary bool
}

// AddOfferImageHandler handles the add image command.
type AddOfferImageHandler struct {
	offerRepo domain.OfferRepository
}

// NewAddOfferImageHandler creates a new handler.
func NewAddOfferImageHandler(offerRepo domain.OfferRepository) *AddOfferImageHandler {
	return &AddOfferImageHandler{
		offerRepo: offerRepo,
	}
}

// Handle executes the add image command.
func (h *AddOfferImageHandler) Handle(ctx context.Context, cmd AddOfferImageCommand) (*domain.Offer, error) {
	offer, err := h.offerRepo.FindByID(ctx, domain.OfferID(cmd.OfferID))
	if err != nil {
		return nil, err
	}
	if offer == nil {
		return nil, domain.ErrOfferNotFound
	}

	offer.AddImage(domain.OfferImage{
		URL:       cmd.URL,
		Alt:       cmd.Alt,
		IsPrimary: cmd.IsPrimary,
	})

	if err := h.offerRepo.Save(ctx, offer); err != nil {
		return nil, err
	}

	return offer, nil
}

// =============================================================================
// Remove Offer Image Command
// =============================================================================

// RemoveOfferImageCommand removes an image from an offer.
type RemoveOfferImageCommand struct {
	OfferID  string
	ImageURL string
}

// RemoveOfferImageHandler handles the remove image command.
type RemoveOfferImageHandler struct {
	offerRepo domain.OfferRepository
}

// NewRemoveOfferImageHandler creates a new handler.
func NewRemoveOfferImageHandler(offerRepo domain.OfferRepository) *RemoveOfferImageHandler {
	return &RemoveOfferImageHandler{
		offerRepo: offerRepo,
	}
}

// Handle executes the remove image command.
func (h *RemoveOfferImageHandler) Handle(ctx context.Context, cmd RemoveOfferImageCommand) (*domain.Offer, error) {
	offer, err := h.offerRepo.FindByID(ctx, domain.OfferID(cmd.OfferID))
	if err != nil {
		return nil, err
	}
	if offer == nil {
		return nil, domain.ErrOfferNotFound
	}

	offer.RemoveImage(cmd.ImageURL)

	if err := h.offerRepo.Save(ctx, offer); err != nil {
		return nil, err
	}

	return offer, nil
}

// =============================================================================
// Increment Offer Stats Command
// =============================================================================

// IncrementOfferViewsCommand increments view count.
type IncrementOfferViewsCommand struct {
	OfferID string
}

// IncrementOfferViewsHandler handles the increment views command.
type IncrementOfferViewsHandler struct {
	offerRepo domain.OfferRepository
}

// NewIncrementOfferViewsHandler creates a new handler.
func NewIncrementOfferViewsHandler(offerRepo domain.OfferRepository) *IncrementOfferViewsHandler {
	return &IncrementOfferViewsHandler{
		offerRepo: offerRepo,
	}
}

// Handle executes the increment views command.
func (h *IncrementOfferViewsHandler) Handle(ctx context.Context, cmd IncrementOfferViewsCommand) error {
	offer, err := h.offerRepo.FindByID(ctx, domain.OfferID(cmd.OfferID))
	if err != nil {
		return err
	}
	if offer == nil {
		return domain.ErrOfferNotFound
	}

	offer.IncrementViews()

	return h.offerRepo.Save(ctx, offer)
}

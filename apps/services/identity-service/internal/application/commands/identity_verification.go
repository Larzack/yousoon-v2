package commands

import (
	"context"

	"github.com/yousoon/services/identity/internal/domain"
	"github.com/yousoon/shared/infrastructure/nats"
)

// SubmitIdentityVerificationCommand represents a command to submit identity verification.
type SubmitIdentityVerificationCommand struct {
	UserID       string
	DocumentType string // cni, passport, driving_license
	Method       string // internal_ocr, external
}

// SubmitIdentityVerificationResult represents the result.
type SubmitIdentityVerificationResult struct {
	VerificationID string
}

// SubmitIdentityVerificationHandler handles identity verification submission.
type SubmitIdentityVerificationHandler struct {
	userRepo       domain.UserRepository
	eventPublisher *nats.EventPublisher
}

// NewSubmitIdentityVerificationHandler creates a new handler.
func NewSubmitIdentityVerificationHandler(
	userRepo domain.UserRepository,
	eventPublisher *nats.EventPublisher,
) *SubmitIdentityVerificationHandler {
	return &SubmitIdentityVerificationHandler{
		userRepo:       userRepo,
		eventPublisher: eventPublisher,
	}
}

// Handle executes the command.
func (h *SubmitIdentityVerificationHandler) Handle(ctx context.Context, cmd SubmitIdentityVerificationCommand) (*SubmitIdentityVerificationResult, error) {
	// Parse user ID
	userID, err := domain.ParseUserID(cmd.UserID)
	if err != nil {
		return nil, err
	}

	// Find user
	user, err := h.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Check if already verified
	if user.Identity != nil && user.Identity.Status == domain.VerificationStatusVerified {
		return nil, domain.ErrAlreadyVerified
	}

	// Check max attempts (10)
	if user.Identity != nil && user.Identity.AttemptCount >= 10 {
		return nil, domain.ErrMaxAttemptsExceeded
	}

	// Create verification
	docType := domain.DocumentType(cmd.DocumentType)
	method := domain.VerificationMethod(cmd.Method)
	verification := domain.NewIdentityVerification(docType, method)

	// Submit verification
	if err := user.SubmitIdentityVerification(verification); err != nil {
		return nil, err
	}

	// Save user
	if err := h.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	// Publish events
	for _, event := range user.GetDomainEvents() {
		if err := h.eventPublisher.Publish(ctx, event); err != nil {
			// Log but don't fail
		}
	}
	user.ClearDomainEvents()

	return &SubmitIdentityVerificationResult{
		VerificationID: verification.ID.String(),
	}, nil
}

// VerifyIdentityCommand represents a command to verify a user's identity (admin action).
type VerifyIdentityCommand struct {
	UserID string
}

// VerifyIdentityHandler handles identity verification approval.
type VerifyIdentityHandler struct {
	userRepo       domain.UserRepository
	eventPublisher *nats.EventPublisher
}

// NewVerifyIdentityHandler creates a new handler.
func NewVerifyIdentityHandler(
	userRepo domain.UserRepository,
	eventPublisher *nats.EventPublisher,
) *VerifyIdentityHandler {
	return &VerifyIdentityHandler{
		userRepo:       userRepo,
		eventPublisher: eventPublisher,
	}
}

// Handle executes the command.
func (h *VerifyIdentityHandler) Handle(ctx context.Context, cmd VerifyIdentityCommand) error {
	// Parse user ID
	userID, err := domain.ParseUserID(cmd.UserID)
	if err != nil {
		return err
	}

	// Find user
	user, err := h.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// Verify identity
	if err := user.ApproveIdentityVerification(); err != nil {
		return err
	}

	// Save user
	if err := h.userRepo.Update(ctx, user); err != nil {
		return err
	}

	// Publish events
	for _, event := range user.GetDomainEvents() {
		if err := h.eventPublisher.Publish(ctx, event); err != nil {
			// Log but don't fail
		}
	}
	user.ClearDomainEvents()

	return nil
}

// RejectIdentityCommand represents a command to reject a user's identity verification.
type RejectIdentityCommand struct {
	UserID string
	Reason string
}

// RejectIdentityHandler handles identity verification rejection.
type RejectIdentityHandler struct {
	userRepo       domain.UserRepository
	eventPublisher *nats.EventPublisher
}

// NewRejectIdentityHandler creates a new handler.
func NewRejectIdentityHandler(
	userRepo domain.UserRepository,
	eventPublisher *nats.EventPublisher,
) *RejectIdentityHandler {
	return &RejectIdentityHandler{
		userRepo:       userRepo,
		eventPublisher: eventPublisher,
	}
}

// Handle executes the command.
func (h *RejectIdentityHandler) Handle(ctx context.Context, cmd RejectIdentityCommand) error {
	// Parse user ID
	userID, err := domain.ParseUserID(cmd.UserID)
	if err != nil {
		return err
	}

	// Find user
	user, err := h.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// Reject identity
	if err := user.RejectIdentityVerification(cmd.Reason); err != nil {
		return err
	}

	// Save user
	if err := h.userRepo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}

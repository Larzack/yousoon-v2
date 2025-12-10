package commands

import (
	"context"

	"github.com/yousoon/services/identity/internal/domain"
	"github.com/yousoon/shared/infrastructure/nats"
)

// UpdateProfileCommand represents a command to update user profile.
type UpdateProfileCommand struct {
	UserID      string
	FirstName   *string
	LastName    *string
	DisplayName *string
	Avatar      *string
	Gender      *string
}

// UpdateProfileHandler handles profile updates.
type UpdateProfileHandler struct {
	userRepo       domain.UserRepository
	eventPublisher *nats.EventPublisher
}

// NewUpdateProfileHandler creates a new UpdateProfileHandler.
func NewUpdateProfileHandler(
	userRepo domain.UserRepository,
	eventPublisher *nats.EventPublisher,
) *UpdateProfileHandler {
	return &UpdateProfileHandler{
		userRepo:       userRepo,
		eventPublisher: eventPublisher,
	}
}

// Handle executes the update profile command.
func (h *UpdateProfileHandler) Handle(ctx context.Context, cmd UpdateProfileCommand) error {
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

	// Update profile fields
	if cmd.FirstName != nil {
		user.Profile.FirstName = *cmd.FirstName
	}
	if cmd.LastName != nil {
		user.Profile.LastName = *cmd.LastName
	}
	if cmd.DisplayName != nil {
		user.Profile.DisplayName = *cmd.DisplayName
	}
	if cmd.Avatar != nil {
		user.Profile.Avatar = cmd.Avatar
	}
	if cmd.Gender != nil {
		gender := domain.Gender(*cmd.Gender)
		user.Profile.Gender = &gender
	}

	// Update user
	user.UpdateProfile(user.Profile)
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

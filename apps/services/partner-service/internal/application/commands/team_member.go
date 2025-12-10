package commands

import (
	"context"

	partnerdomain "github.com/yousoon/services/partner/internal/domain"
	sharedomain "github.com/yousoon/shared/domain"
)

// =============================================================================
// Invite Team Member Command
// =============================================================================

// InviteTeamMemberCommand represents a request to invite a team member.
type InviteTeamMemberCommand struct {
	PartnerID string `json:"partnerId"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Role      string `json:"role"`
}

// Validate validates the command.
func (c *InviteTeamMemberCommand) Validate() error {
	var errs partnerdomain.ValidationErrors

	if c.PartnerID == "" {
		errs.Add("partnerId", "partner ID is required")
	}
	if c.Email == "" {
		errs.Add("email", "email is required")
	}
	if c.FirstName == "" {
		errs.Add("firstName", "first name is required")
	}
	if c.LastName == "" {
		errs.Add("lastName", "last name is required")
	}
	if c.Role == "" {
		errs.Add("role", "role is required")
	}

	if errs.HasErrors() {
		return errs
	}
	return nil
}

// InviteTeamMemberHandler handles the InviteTeamMemberCommand.
type InviteTeamMemberHandler struct {
	repo      partnerdomain.PartnerRepository
	publisher sharedomain.EventPublisher
}

// NewInviteTeamMemberHandler creates a new handler.
func NewInviteTeamMemberHandler(repo partnerdomain.PartnerRepository, publisher sharedomain.EventPublisher) *InviteTeamMemberHandler {
	return &InviteTeamMemberHandler{
		repo:      repo,
		publisher: publisher,
	}
}

// Handle handles the command.
func (h *InviteTeamMemberHandler) Handle(ctx context.Context, cmd InviteTeamMemberCommand) (*partnerdomain.TeamMember, error) {
	// Validate command
	if err := cmd.Validate(); err != nil {
		return nil, err
	}

	// Get partner
	partner, err := h.repo.FindByID(ctx, partnerdomain.PartnerID(cmd.PartnerID))
	if err != nil {
		return nil, err
	}

	// Create email value object
	email, err := partnerdomain.NewEmail(cmd.Email)
	if err != nil {
		return nil, err
	}

	// Validate role
	role := partnerdomain.TeamRole(cmd.Role)
	if !role.IsValid() {
		return nil, partnerdomain.ErrInvalidTeamRole
	}

	// Create team member
	member := partnerdomain.NewTeamMember(email, cmd.FirstName, cmd.LastName, role)

	// Add to partner
	if err := partner.AddTeamMember(member); err != nil {
		return nil, err
	}

	// Save partner
	if err := h.repo.Save(ctx, partner); err != nil {
		return nil, err
	}

	// Publish domain events
	for _, event := range partner.GetDomainEvents() {
		if err := h.publisher.Publish(ctx, event); err != nil {
			// Log error but don't fail
		}
	}
	partner.ClearDomainEvents()

	// Return the created member
	created, _ := partner.GetTeamMember(member.ID)
	return created, nil
}

// =============================================================================
// Accept Team Invitation Command
// =============================================================================

// AcceptTeamInvitationCommand represents a request to accept an invitation.
type AcceptTeamInvitationCommand struct {
	PartnerID    string `json:"partnerId"`
	TeamMemberID string `json:"teamMemberId"`
	UserID       string `json:"userId"`
}

// Validate validates the command.
func (c *AcceptTeamInvitationCommand) Validate() error {
	var errs partnerdomain.ValidationErrors

	if c.PartnerID == "" {
		errs.Add("partnerId", "partner ID is required")
	}
	if c.TeamMemberID == "" {
		errs.Add("teamMemberId", "team member ID is required")
	}
	if c.UserID == "" {
		errs.Add("userId", "user ID is required")
	}

	if errs.HasErrors() {
		return errs
	}
	return nil
}

// AcceptTeamInvitationHandler handles the AcceptTeamInvitationCommand.
type AcceptTeamInvitationHandler struct {
	repo      partnerdomain.PartnerRepository
	publisher sharedomain.EventPublisher
}

// NewAcceptTeamInvitationHandler creates a new handler.
func NewAcceptTeamInvitationHandler(repo partnerdomain.PartnerRepository, publisher sharedomain.EventPublisher) *AcceptTeamInvitationHandler {
	return &AcceptTeamInvitationHandler{
		repo:      repo,
		publisher: publisher,
	}
}

// Handle handles the command.
func (h *AcceptTeamInvitationHandler) Handle(ctx context.Context, cmd AcceptTeamInvitationCommand) error {
	// Validate command
	if err := cmd.Validate(); err != nil {
		return err
	}

	// Get partner
	partner, err := h.repo.FindByID(ctx, partnerdomain.PartnerID(cmd.PartnerID))
	if err != nil {
		return err
	}

	// Accept invitation
	err = partner.AcceptTeamInvitation(
		partnerdomain.TeamMemberID(cmd.TeamMemberID),
		partnerdomain.UserID(cmd.UserID),
	)
	if err != nil {
		return err
	}

	// Save partner
	if err := h.repo.Save(ctx, partner); err != nil {
		return err
	}

	// Publish domain events
	for _, event := range partner.GetDomainEvents() {
		if err := h.publisher.Publish(ctx, event); err != nil {
			// Log error but don't fail
		}
	}
	partner.ClearDomainEvents()

	return nil
}

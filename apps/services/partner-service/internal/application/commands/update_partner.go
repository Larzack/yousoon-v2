package commands

import (
	"context"

	partnerdomain "github.com/yousoon/services/partner/internal/domain"
	sharedomain "github.com/yousoon/shared/domain"
)

// =============================================================================
// Update Partner Command
// =============================================================================

// UpdatePartnerCommand represents a request to update partner information.
type UpdatePartnerCommand struct {
	PartnerID string `json:"partnerId"`

	// Company (optional updates)
	CompanyName *string `json:"companyName,omitempty"`
	TradeName   *string `json:"tradeName,omitempty"`
	VATNumber   *string `json:"vatNumber,omitempty"`
	LegalForm   *string `json:"legalForm,omitempty"`

	// Branding (optional updates)
	Logo         *string `json:"logo,omitempty"`
	CoverImage   *string `json:"coverImage,omitempty"`
	PrimaryColor *string `json:"primaryColor,omitempty"`
	Description  *string `json:"description,omitempty"`

	// Contact (optional updates)
	ContactFirstName *string `json:"contactFirstName,omitempty"`
	ContactLastName  *string `json:"contactLastName,omitempty"`
	ContactEmail     *string `json:"contactEmail,omitempty"`
	ContactPhone     *string `json:"contactPhone,omitempty"`
	ContactRole      *string `json:"contactRole,omitempty"`

	// Category (optional updates)
	Category      *string   `json:"category,omitempty"`
	Subcategories *[]string `json:"subcategories,omitempty"`
}

// Validate validates the command.
func (c *UpdatePartnerCommand) Validate() error {
	var errs partnerdomain.ValidationErrors

	if c.PartnerID == "" {
		errs.Add("partnerId", "partner ID is required")
	}

	if errs.HasErrors() {
		return errs
	}
	return nil
}

// UpdatePartnerHandler handles the UpdatePartnerCommand.
type UpdatePartnerHandler struct {
	repo      partnerdomain.PartnerRepository
	publisher sharedomain.EventPublisher
}

// NewUpdatePartnerHandler creates a new handler.
func NewUpdatePartnerHandler(repo partnerdomain.PartnerRepository, publisher sharedomain.EventPublisher) *UpdatePartnerHandler {
	return &UpdatePartnerHandler{
		repo:      repo,
		publisher: publisher,
	}
}

// Handle handles the command.
func (h *UpdatePartnerHandler) Handle(ctx context.Context, cmd UpdatePartnerCommand) (*partnerdomain.Partner, error) {
	// Validate command
	if err := cmd.Validate(); err != nil {
		return nil, err
	}

	// Get partner
	partner, err := h.repo.FindByID(ctx, partnerdomain.PartnerID(cmd.PartnerID))
	if err != nil {
		return nil, err
	}

	// Update company if any field is provided
	if cmd.CompanyName != nil || cmd.TradeName != nil || cmd.VATNumber != nil || cmd.LegalForm != nil {
		company := partner.Company
		if cmd.CompanyName != nil {
			company.Name = *cmd.CompanyName
		}
		if cmd.TradeName != nil {
			company.TradeName = *cmd.TradeName
		}
		if cmd.VATNumber != nil {
			company.VATNumber = *cmd.VATNumber
		}
		if cmd.LegalForm != nil {
			company.LegalForm = *cmd.LegalForm
		}
		partner.UpdateCompany(company)
	}

	// Update branding if any field is provided
	if cmd.Logo != nil || cmd.CoverImage != nil || cmd.PrimaryColor != nil || cmd.Description != nil {
		branding := partner.Branding
		if cmd.Logo != nil {
			branding.Logo = *cmd.Logo
		}
		if cmd.CoverImage != nil {
			branding.CoverImage = *cmd.CoverImage
		}
		if cmd.PrimaryColor != nil {
			branding.PrimaryColor = *cmd.PrimaryColor
		}
		if cmd.Description != nil {
			branding.Description = *cmd.Description
		}
		partner.UpdateBranding(branding)
	}

	// Update contact if any field is provided
	if cmd.ContactFirstName != nil || cmd.ContactLastName != nil || cmd.ContactEmail != nil || cmd.ContactPhone != nil || cmd.ContactRole != nil {
		contact := partner.Contact
		if cmd.ContactFirstName != nil {
			contact.FirstName = *cmd.ContactFirstName
		}
		if cmd.ContactLastName != nil {
			contact.LastName = *cmd.ContactLastName
		}
		if cmd.ContactEmail != nil {
			contact.Email = *cmd.ContactEmail
		}
		if cmd.ContactPhone != nil {
			contact.Phone = *cmd.ContactPhone
		}
		if cmd.ContactRole != nil {
			contact.Role = *cmd.ContactRole
		}
		partner.UpdateContact(contact)
	}

	// Update category if provided
	if cmd.Category != nil {
		subcategories := partner.Subcategories
		if cmd.Subcategories != nil {
			subcategories = *cmd.Subcategories
		}
		partner.UpdateCategory(*cmd.Category, subcategories)
	} else if cmd.Subcategories != nil {
		partner.UpdateCategory(partner.Category, *cmd.Subcategories)
	}

	// Save partner
	if err := h.repo.Save(ctx, partner); err != nil {
		return nil, err
	}

	// Publish domain events
	for _, event := range partner.GetDomainEvents() {
		if err := h.publisher.Publish(event); err != nil {
			// Log error but don't fail
		}
	}
	partner.ClearDomainEvents()

	return partner, nil
}

// =============================================================================
// Verify Partner Command
// =============================================================================

// VerifyPartnerCommand represents a request to verify a partner.
type VerifyPartnerCommand struct {
	PartnerID    string `json:"partnerId"`
	VerifiedByID string `json:"verifiedById"`
}

// Validate validates the command.
func (c *VerifyPartnerCommand) Validate() error {
	var errs partnerdomain.ValidationErrors

	if c.PartnerID == "" {
		errs.Add("partnerId", "partner ID is required")
	}
	if c.VerifiedByID == "" {
		errs.Add("verifiedById", "verified by ID is required")
	}

	if errs.HasErrors() {
		return errs
	}
	return nil
}

// VerifyPartnerHandler handles the VerifyPartnerCommand.
type VerifyPartnerHandler struct {
	repo      partnerdomain.PartnerRepository
	publisher sharedomain.EventPublisher
}

// NewVerifyPartnerHandler creates a new handler.
func NewVerifyPartnerHandler(repo partnerdomain.PartnerRepository, publisher sharedomain.EventPublisher) *VerifyPartnerHandler {
	return &VerifyPartnerHandler{
		repo:      repo,
		publisher: publisher,
	}
}

// Handle handles the command.
func (h *VerifyPartnerHandler) Handle(ctx context.Context, cmd VerifyPartnerCommand) error {
	// Validate command
	if err := cmd.Validate(); err != nil {
		return err
	}

	// Get partner
	partner, err := h.repo.FindByID(ctx, partnerdomain.PartnerID(cmd.PartnerID))
	if err != nil {
		return err
	}

	// Verify partner
	if err := partner.Verify(cmd.VerifiedByID); err != nil {
		return err
	}

	// Save partner
	if err := h.repo.Save(ctx, partner); err != nil {
		return err
	}

	// Publish domain events
	for _, event := range partner.GetDomainEvents() {
		if err := h.publisher.Publish(event); err != nil {
			// Log error but don't fail
		}
	}
	partner.ClearDomainEvents()

	return nil
}

// =============================================================================
// Suspend Partner Command
// =============================================================================

// SuspendPartnerCommand represents a request to suspend a partner.
type SuspendPartnerCommand struct {
	PartnerID string `json:"partnerId"`
	Reason    string `json:"reason"`
}

// Validate validates the command.
func (c *SuspendPartnerCommand) Validate() error {
	var errs partnerdomain.ValidationErrors

	if c.PartnerID == "" {
		errs.Add("partnerId", "partner ID is required")
	}
	if c.Reason == "" {
		errs.Add("reason", "reason is required")
	}

	if errs.HasErrors() {
		return errs
	}
	return nil
}

// SuspendPartnerHandler handles the SuspendPartnerCommand.
type SuspendPartnerHandler struct {
	repo      partnerdomain.PartnerRepository
	publisher sharedomain.EventPublisher
}

// NewSuspendPartnerHandler creates a new handler.
func NewSuspendPartnerHandler(repo partnerdomain.PartnerRepository, publisher sharedomain.EventPublisher) *SuspendPartnerHandler {
	return &SuspendPartnerHandler{
		repo:      repo,
		publisher: publisher,
	}
}

// Handle handles the command.
func (h *SuspendPartnerHandler) Handle(ctx context.Context, cmd SuspendPartnerCommand) error {
	// Validate command
	if err := cmd.Validate(); err != nil {
		return err
	}

	// Get partner
	partner, err := h.repo.FindByID(ctx, partnerdomain.PartnerID(cmd.PartnerID))
	if err != nil {
		return err
	}

	// Suspend partner
	partner.Suspend(cmd.Reason)

	// Save partner
	if err := h.repo.Save(ctx, partner); err != nil {
		return err
	}

	// Publish domain events
	for _, event := range partner.GetDomainEvents() {
		if err := h.publisher.Publish(event); err != nil {
			// Log error but don't fail
		}
	}
	partner.ClearDomainEvents()

	return nil
}

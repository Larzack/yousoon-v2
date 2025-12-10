// Package commands contains the command handlers for the Partner bounded context.
package commands

import (
	"context"

	sharedomain "github.com/yousoon/shared/domain"
)

// =============================================================================
// Register Partner Command
// =============================================================================

// RegisterPartnerCommand represents a request to register a new partner.
type RegisterPartnerCommand struct {
	// Owner
	OwnerUserID string `json:"ownerUserId"`

	// Company
	CompanyName string `json:"companyName"`
	TradeName   string `json:"tradeName"`
	SIRET       string `json:"siret"`
	VATNumber   string `json:"vatNumber"`
	LegalForm   string `json:"legalForm"`

	// Contact
	ContactFirstName string `json:"contactFirstName"`
	ContactLastName  string `json:"contactLastName"`
	ContactEmail     string `json:"contactEmail"`
	ContactPhone     string `json:"contactPhone"`
	ContactRole      string `json:"contactRole"`

	// Category
	Category      string   `json:"category"`
	Subcategories []string `json:"subcategories"`
}

// Validate validates the command.
func (c *RegisterPartnerCommand) Validate() error {
	var errs partnerdomain.ValidationErrors

	if c.OwnerUserID == "" {
		errs.Add("ownerUserId", "owner user ID is required")
	}
	if c.CompanyName == "" {
		errs.Add("companyName", "company name is required")
	}
	if c.SIRET == "" {
		errs.Add("siret", "SIRET is required")
	}
	if c.ContactFirstName == "" {
		errs.Add("contactFirstName", "contact first name is required")
	}
	if c.ContactLastName == "" {
		errs.Add("contactLastName", "contact last name is required")
	}
	if c.ContactEmail == "" {
		errs.Add("contactEmail", "contact email is required")
	}
	if c.Category == "" {
		errs.Add("category", "category is required")
	}

	if errs.HasErrors() {
		return errs
	}
	return nil
}

// RegisterPartnerHandler handles the RegisterPartnerCommand.
type RegisterPartnerHandler struct {
	repo      partnerdomain.PartnerRepository
	publisher sharedomain.EventPublisher
}

// NewRegisterPartnerHandler creates a new handler.
func NewRegisterPartnerHandler(repo partnerdomain.PartnerRepository, publisher sharedomain.EventPublisher) *RegisterPartnerHandler {
	return &RegisterPartnerHandler{
		repo:      repo,
		publisher: publisher,
	}
}

// Handle handles the command.
func (h *RegisterPartnerHandler) Handle(ctx context.Context, cmd RegisterPartnerCommand) (*partnerdomain.Partner, error) {
	// Validate command
	if err := cmd.Validate(); err != nil {
		return nil, err
	}

	// Check if SIRET already exists
	existing, err := h.repo.FindBySIRET(ctx, cmd.SIRET)
	if err != nil && err != partnerdomain.ErrPartnerNotFound {
		return nil, err
	}
	if existing != nil {
		return nil, partnerdomain.ErrSIRETAlreadyExists
	}

	// Create company value object
	company, err := partnerdomain.NewCompany(
		cmd.CompanyName,
		cmd.TradeName,
		cmd.SIRET,
		cmd.VATNumber,
		cmd.LegalForm,
	)
	if err != nil {
		return nil, err
	}

	// Create contact value object
	contact, err := partnerdomain.NewContact(
		cmd.ContactFirstName,
		cmd.ContactLastName,
		cmd.ContactEmail,
		cmd.ContactPhone,
		cmd.ContactRole,
	)
	if err != nil {
		return nil, err
	}

	// Create partner
	partner, err := partnerdomain.NewPartner(
		partnerdomain.UserID(cmd.OwnerUserID),
		company,
		contact,
		cmd.Category,
	)
	if err != nil {
		return nil, err
	}

	// Set subcategories
	partner.Subcategories = cmd.Subcategories

	// Save partner
	if err := h.repo.Save(ctx, partner); err != nil {
		return nil, err
	}

	// Publish domain events
	for _, event := range partner.GetDomainEvents() {
		if err := h.publisher.Publish(ctx, event); err != nil {
			// Log error but don't fail the command
			// Events will be retried via outbox pattern
		}
	}
	partner.ClearDomainEvents()

	return partner, nil
}

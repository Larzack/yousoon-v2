package commands

import (
	"context"

	partnerdomain "github.com/yousoon/services/partner/internal/domain"
	sharedomain "github.com/yousoon/services/shared/domain"
)

// =============================================================================
// Add Establishment Command
// =============================================================================

// AddEstablishmentCommand represents a request to add an establishment.
type AddEstablishmentCommand struct {
	PartnerID string `json:"partnerId"`

	// Basic Info
	Name        string `json:"name"`
	Description string `json:"description"`

	// Address
	Street       string `json:"street"`
	StreetNumber string `json:"streetNumber"`
	Complement   string `json:"complement"`
	PostalCode   string `json:"postalCode"`
	City         string `json:"city"`
	Country      string `json:"country"`

	// Location
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`

	// Contact
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Website string `json:"website"`

	// Details
	Type       string   `json:"type"`
	Features   []string `json:"features"`
	PriceRange int      `json:"priceRange"`
}

// Validate validates the command.
func (c *AddEstablishmentCommand) Validate() error {
	var errs partnerdomain.ValidationErrors

	if c.PartnerID == "" {
		errs.Add("partnerId", "partner ID is required")
	}
	if c.Name == "" {
		errs.Add("name", "name is required")
	}
	if c.Street == "" {
		errs.Add("street", "street is required")
	}
	if c.PostalCode == "" {
		errs.Add("postalCode", "postal code is required")
	}
	if c.City == "" {
		errs.Add("city", "city is required")
	}
	if c.Country == "" {
		errs.Add("country", "country is required")
	}
	if c.Longitude == 0 && c.Latitude == 0 {
		errs.Add("location", "location is required")
	}

	if errs.HasErrors() {
		return errs
	}
	return nil
}

// AddEstablishmentHandler handles the AddEstablishmentCommand.
type AddEstablishmentHandler struct {
	repo      partnerdomain.PartnerRepository
	publisher sharedomain.EventPublisher
}

// NewAddEstablishmentHandler creates a new handler.
func NewAddEstablishmentHandler(repo partnerdomain.PartnerRepository, publisher sharedomain.EventPublisher) *AddEstablishmentHandler {
	return &AddEstablishmentHandler{
		repo:      repo,
		publisher: publisher,
	}
}

// Handle handles the command.
func (h *AddEstablishmentHandler) Handle(ctx context.Context, cmd AddEstablishmentCommand) (*partnerdomain.Establishment, error) {
	// Validate command
	if err := cmd.Validate(); err != nil {
		return nil, err
	}

	// Get partner
	partner, err := h.repo.FindByID(ctx, partnerdomain.PartnerID(cmd.PartnerID))
	if err != nil {
		return nil, err
	}

	// Create address
	address := partnerdomain.NewAddress(
		cmd.Street,
		cmd.StreetNumber,
		cmd.Complement,
		cmd.PostalCode,
		cmd.City,
		cmd.Country,
	)

	// Create location
	location, err := partnerdomain.NewGeoLocation(cmd.Longitude, cmd.Latitude)
	if err != nil {
		return nil, err
	}

	// Create establishment
	establishment := partnerdomain.NewEstablishment(cmd.Name, address, location)
	establishment.Description = cmd.Description
	establishment.SetContact(partnerdomain.NewEstablishmentContact(cmd.Phone, cmd.Email, cmd.Website))
	establishment.SetType(cmd.Type)
	establishment.SetFeatures(cmd.Features)
	establishment.SetPriceRange(cmd.PriceRange)

	// Add to partner
	if err := partner.AddEstablishment(establishment); err != nil {
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

	// Return the created establishment
	created, _ := partner.GetEstablishment(establishment.ID)
	return created, nil
}

// Package commands contains command handlers for creating offers.
package commands

import (
	"context"
	"errors"
	"time"

	"github.com/yousoon/discovery-service/internal/domain"
)

// =============================================================================
// Create Offer Command
// =============================================================================

// CreateOfferCommand represents a command to create a new offer.
type CreateOfferCommand struct {
	PartnerID          string
	EstablishmentID    string
	Title              string
	Description        string
	ShortDescription   string
	CategoryID         string
	Tags               []string
	Discount           DiscountInput
	Conditions         []ConditionInput
	TermsAndConditions string
	Validity           ValidityInput
	Schedule           ScheduleInput
	Quota              QuotaInput
	Images             []ImageInput

	// Denormalized data (provided by caller via ACL)
	PartnerName           string
	PartnerLogo           string
	PartnerCategory       string
	EstablishmentName     string
	EstablishmentAddress  string
	EstablishmentCity     string
	EstablishmentLocation LocationInput
}

// DiscountInput represents discount input.
type DiscountInput struct {
	Type          string
	Value         int
	OriginalPrice *int64
	Formula       string
}

// ConditionInput represents condition input.
type ConditionInput struct {
	Type  string
	Value interface{}
	Label string
}

// ValidityInput represents validity period input.
type ValidityInput struct {
	StartDate time.Time
	EndDate   time.Time
	Timezone  string
}

// ScheduleInput represents schedule input.
type ScheduleInput struct {
	AllDay bool
	Slots  []TimeSlotInput
}

// TimeSlotInput represents time slot input.
type TimeSlotInput struct {
	DayOfWeek int
	StartTime string
	EndTime   string
}

// QuotaInput represents quota input.
type QuotaInput struct {
	Total   *int
	PerUser *int
	PerDay  *int
}

// ImageInput represents image input.
type ImageInput struct {
	URL       string
	Alt       string
	IsPrimary bool
}

// LocationInput represents location input.
type LocationInput struct {
	Longitude float64
	Latitude  float64
}

// CreateOfferHandler handles the create offer command.
type CreateOfferHandler struct {
	offerRepo domain.OfferRepository
}

// NewCreateOfferHandler creates a new CreateOfferHandler.
func NewCreateOfferHandler(offerRepo domain.OfferRepository) *CreateOfferHandler {
	return &CreateOfferHandler{
		offerRepo: offerRepo,
	}
}

// Handle executes the create offer command.
func (h *CreateOfferHandler) Handle(ctx context.Context, cmd CreateOfferCommand) (*domain.Offer, error) {
	// Validate required fields
	if cmd.PartnerID == "" {
		return nil, errors.New("partner ID is required")
	}
	if cmd.EstablishmentID == "" {
		return nil, errors.New("establishment ID is required")
	}
	if cmd.Title == "" {
		return nil, errors.New("title is required")
	}
	if cmd.CategoryID == "" {
		return nil, errors.New("category ID is required")
	}

	// Create discount value object
	discount := domain.Discount{
		Type:          domain.DiscountType(cmd.Discount.Type),
		Value:         cmd.Discount.Value,
		OriginalPrice: cmd.Discount.OriginalPrice,
		Formula:       cmd.Discount.Formula,
	}
	if err := discount.Validate(); err != nil {
		return nil, err
	}

	// Create validity value object
	validity, err := domain.NewValidity(cmd.Validity.StartDate, cmd.Validity.EndDate, cmd.Validity.Timezone)
	if err != nil {
		return nil, err
	}

	// Create offer
	offer, err := domain.NewOffer(
		domain.PartnerID(cmd.PartnerID),
		domain.EstablishmentID(cmd.EstablishmentID),
		cmd.Title,
		cmd.Description,
		domain.CategoryID(cmd.CategoryID),
		discount,
		validity,
	)
	if err != nil {
		return nil, err
	}

	// Set optional fields
	if cmd.ShortDescription != "" {
		offer.UpdateBasicInfo(cmd.Title, cmd.Description, cmd.ShortDescription)
	}

	if len(cmd.Tags) > 0 {
		offer.UpdateTags(cmd.Tags)
	}

	// Set conditions
	if len(cmd.Conditions) > 0 {
		conditions := make([]domain.Condition, len(cmd.Conditions))
		for i, c := range cmd.Conditions {
			conditions[i] = domain.Condition{
				Type:  domain.ConditionType(c.Type),
				Value: c.Value,
				Label: c.Label,
			}
		}
		offer.UpdateConditions(conditions, cmd.TermsAndConditions)
	}

	// Set schedule
	if !cmd.Schedule.AllDay && len(cmd.Schedule.Slots) > 0 {
		slots := make([]domain.TimeSlot, len(cmd.Schedule.Slots))
		for i, s := range cmd.Schedule.Slots {
			slots[i] = domain.TimeSlot{
				DayOfWeek: s.DayOfWeek,
				StartTime: s.StartTime,
				EndTime:   s.EndTime,
			}
		}
		offer.UpdateSchedule(domain.NewScheduleWithSlots(slots))
	}

	// Set quota
	if cmd.Quota.Total != nil || cmd.Quota.PerUser != nil || cmd.Quota.PerDay != nil {
		offer.UpdateQuota(domain.NewQuota(cmd.Quota.Total, cmd.Quota.PerUser, cmd.Quota.PerDay))
	}

	// Set images
	for _, img := range cmd.Images {
		offer.AddImage(domain.OfferImage{
			URL:       img.URL,
			Alt:       img.Alt,
			IsPrimary: img.IsPrimary,
		})
	}

	// Set denormalized data
	offer.SetPartnerSnapshot(domain.PartnerSnapshot{
		Name:     cmd.PartnerName,
		Logo:     cmd.PartnerLogo,
		Category: cmd.PartnerCategory,
	})

	location, _ := domain.NewGeoLocation(cmd.EstablishmentLocation.Longitude, cmd.EstablishmentLocation.Latitude)
	offer.SetEstablishmentSnapshot(domain.EstablishmentSnapshot{
		Name:     cmd.EstablishmentName,
		Address:  cmd.EstablishmentAddress,
		City:     cmd.EstablishmentCity,
		Location: location,
	})

	// Save offer
	if err := h.offerRepo.Save(ctx, offer); err != nil {
		return nil, err
	}

	return offer, nil
}

// =============================================================================
// Update Offer Command
// =============================================================================

// UpdateOfferCommand represents a command to update an offer.
type UpdateOfferCommand struct {
	OfferID            string
	Title              *string
	Description        *string
	ShortDescription   *string
	CategoryID         *string
	Tags               []string
	Discount           *DiscountInput
	Conditions         []ConditionInput
	TermsAndConditions *string
	Validity           *ValidityInput
	Schedule           *ScheduleInput
	Quota              *QuotaInput
}

// UpdateOfferHandler handles the update offer command.
type UpdateOfferHandler struct {
	offerRepo domain.OfferRepository
}

// NewUpdateOfferHandler creates a new UpdateOfferHandler.
func NewUpdateOfferHandler(offerRepo domain.OfferRepository) *UpdateOfferHandler {
	return &UpdateOfferHandler{
		offerRepo: offerRepo,
	}
}

// Handle executes the update offer command.
func (h *UpdateOfferHandler) Handle(ctx context.Context, cmd UpdateOfferCommand) (*domain.Offer, error) {
	offer, err := h.offerRepo.FindByID(ctx, domain.OfferID(cmd.OfferID))
	if err != nil {
		return nil, err
	}
	if offer == nil {
		return nil, domain.ErrOfferNotFound
	}

	// Update basic info
	if cmd.Title != nil || cmd.Description != nil || cmd.ShortDescription != nil {
		title := offer.Title()
		description := offer.Description()
		shortDescription := offer.ShortDescription()

		if cmd.Title != nil {
			title = *cmd.Title
		}
		if cmd.Description != nil {
			description = *cmd.Description
		}
		if cmd.ShortDescription != nil {
			shortDescription = *cmd.ShortDescription
		}

		if err := offer.UpdateBasicInfo(title, description, shortDescription); err != nil {
			return nil, err
		}
	}

	// Update category
	if cmd.CategoryID != nil {
		if err := offer.UpdateCategory(domain.CategoryID(*cmd.CategoryID)); err != nil {
			return nil, err
		}
	}

	// Update tags
	if cmd.Tags != nil {
		offer.UpdateTags(cmd.Tags)
	}

	// Update discount
	if cmd.Discount != nil {
		discount := domain.Discount{
			Type:          domain.DiscountType(cmd.Discount.Type),
			Value:         cmd.Discount.Value,
			OriginalPrice: cmd.Discount.OriginalPrice,
			Formula:       cmd.Discount.Formula,
		}
		if err := offer.UpdateDiscount(discount); err != nil {
			return nil, err
		}
	}

	// Update conditions
	if cmd.Conditions != nil {
		conditions := make([]domain.Condition, len(cmd.Conditions))
		for i, c := range cmd.Conditions {
			conditions[i] = domain.Condition{
				Type:  domain.ConditionType(c.Type),
				Value: c.Value,
				Label: c.Label,
			}
		}
		terms := offer.TermsAndConditions()
		if cmd.TermsAndConditions != nil {
			terms = *cmd.TermsAndConditions
		}
		offer.UpdateConditions(conditions, terms)
	}

	// Update validity
	if cmd.Validity != nil {
		validity, err := domain.NewValidity(cmd.Validity.StartDate, cmd.Validity.EndDate, cmd.Validity.Timezone)
		if err != nil {
			return nil, err
		}
		if err := offer.UpdateValidity(validity); err != nil {
			return nil, err
		}
	}

	// Update schedule
	if cmd.Schedule != nil {
		if cmd.Schedule.AllDay {
			offer.UpdateSchedule(domain.NewAllDaySchedule())
		} else {
			slots := make([]domain.TimeSlot, len(cmd.Schedule.Slots))
			for i, s := range cmd.Schedule.Slots {
				slots[i] = domain.TimeSlot{
					DayOfWeek: s.DayOfWeek,
					StartTime: s.StartTime,
					EndTime:   s.EndTime,
				}
			}
			offer.UpdateSchedule(domain.NewScheduleWithSlots(slots))
		}
	}

	// Update quota
	if cmd.Quota != nil {
		offer.UpdateQuota(domain.NewQuota(cmd.Quota.Total, cmd.Quota.PerUser, cmd.Quota.PerDay))
	}

	// Save changes
	if err := h.offerRepo.Save(ctx, offer); err != nil {
		return nil, err
	}

	return offer, nil
}

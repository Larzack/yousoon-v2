// Package domain contains the domain layer for the Partner bounded context.
// This includes Partner aggregate, Establishment entity, TeamMember entity,
// value objects, domain events, and repository interfaces.
package domain

import (
	"time"

	"github.com/yousoon/services/shared/domain"
)

// =============================================================================
// Partner Aggregate Root
// =============================================================================

// Partner is the aggregate root for partner/business management.
// A Partner represents a business (bar, restaurant, activity organizer) that
// offers discounts to Yousoon users.
type Partner struct {
	domain.VersionedAggregateRoot

	// Identity
	ID          PartnerID `json:"id" bson:"_id"`
	OwnerUserID UserID    `json:"ownerUserId" bson:"ownerUserId"`

	// Company Information
	Company Company `json:"company" bson:"company"`

	// Branding
	Branding Branding `json:"branding" bson:"branding"`

	// Primary Contact
	Contact Contact `json:"contact" bson:"contact"`

	// Category
	Category      string   `json:"category" bson:"category"`
	Subcategories []string `json:"subcategories,omitempty" bson:"subcategories,omitempty"`

	// Establishments (part of the aggregate)
	Establishments []Establishment `json:"establishments,omitempty" bson:"establishments,omitempty"`

	// Team Members
	TeamMembers []TeamMember `json:"teamMembers,omitempty" bson:"teamMembers,omitempty"`

	// Statistics (denormalized for performance)
	Stats PartnerStats `json:"stats" bson:"stats"`

	// Status
	Status     PartnerStatus `json:"status" bson:"status"`
	VerifiedAt *time.Time    `json:"verifiedAt,omitempty" bson:"verifiedAt,omitempty"`

	// Timestamps
	DeletedAt *time.Time `json:"deletedAt,omitempty" bson:"deletedAt,omitempty"`
}

// =============================================================================
// Partner Factory
// =============================================================================

// NewPartner creates a new partner with required fields.
func NewPartner(ownerUserID UserID, company Company, contact Contact, category string) (*Partner, error) {
	partner := &Partner{
		VersionedAggregateRoot: domain.NewVersionedAggregateRoot(),
		ID:                     NewPartnerID(),
		OwnerUserID:            ownerUserID,
		Company:                company,
		Contact:                contact,
		Category:               category,
		Branding:               Branding{},
		Establishments:         make([]Establishment, 0),
		TeamMembers:            make([]TeamMember, 0),
		Stats:                  NewPartnerStats(),
		Status:                 PartnerStatusPending,
	}

	partner.AddDomainEvent(NewPartnerRegisteredEvent(partner.ID, ownerUserID, company.Name))
	return partner, nil
}

// =============================================================================
// Partner Methods (Business Logic)
// =============================================================================

// CanPublishOffer checks if the partner can publish offers.
func (p *Partner) CanPublishOffer() error {
	if p.Status != PartnerStatusActive {
		return ErrPartnerNotVerified
	}
	if len(p.Establishments) == 0 {
		return ErrNoEstablishment
	}
	return nil
}

// Verify verifies the partner account.
func (p *Partner) Verify(verifiedByUserID string) error {
	if p.Status == PartnerStatusActive {
		return ErrAlreadyVerified
	}
	now := time.Now()
	p.Status = PartnerStatusActive
	p.VerifiedAt = &now
	p.MarkUpdated()
	p.AddDomainEvent(NewPartnerVerifiedEvent(p.ID, verifiedByUserID))
	return nil
}

// Suspend suspends the partner account.
func (p *Partner) Suspend(reason string) {
	p.Status = PartnerStatusSuspended
	p.MarkUpdated()
	p.AddDomainEvent(NewPartnerSuspendedEvent(p.ID, reason))
}

// Activate activates a suspended partner.
func (p *Partner) Activate() error {
	if p.VerifiedAt == nil {
		return ErrPartnerNotVerified
	}
	p.Status = PartnerStatusActive
	p.MarkUpdated()
	return nil
}

// =============================================================================
// Company Management
// =============================================================================

// UpdateCompany updates the company information.
func (p *Partner) UpdateCompany(company Company) {
	p.Company = company
	p.MarkUpdated()
}

// UpdateBranding updates the branding information.
func (p *Partner) UpdateBranding(branding Branding) {
	p.Branding = branding
	p.MarkUpdated()
}

// UpdateContact updates the primary contact.
func (p *Partner) UpdateContact(contact Contact) {
	p.Contact = contact
	p.MarkUpdated()
}

// UpdateCategory updates the partner's category.
func (p *Partner) UpdateCategory(category string, subcategories []string) {
	p.Category = category
	p.Subcategories = subcategories
	p.MarkUpdated()
}

// =============================================================================
// Establishment Management
// =============================================================================

// AddEstablishment adds a new establishment to the partner.
func (p *Partner) AddEstablishment(est Establishment) error {
	// Check for duplicate address
	for _, existing := range p.Establishments {
		if existing.Address.Equals(est.Address) {
			return ErrEstablishmentAlreadyExists
		}
	}

	est.PartnerID = p.ID
	p.Establishments = append(p.Establishments, est)
	p.Stats.TotalEstablishments++
	p.MarkUpdated()
	p.AddDomainEvent(NewEstablishmentAddedEvent(p.ID, est.ID, est.Name, est.Location))
	return nil
}

// UpdateEstablishment updates an existing establishment.
func (p *Partner) UpdateEstablishment(estID EstablishmentID, updateFn func(*Establishment) error) error {
	for i := range p.Establishments {
		if p.Establishments[i].ID == estID {
			if err := updateFn(&p.Establishments[i]); err != nil {
				return err
			}
			p.MarkUpdated()
			return nil
		}
	}
	return ErrEstablishmentNotFound
}

// RemoveEstablishment removes an establishment from the partner.
func (p *Partner) RemoveEstablishment(estID EstablishmentID) error {
	for i, est := range p.Establishments {
		if est.ID == estID {
			p.Establishments = append(p.Establishments[:i], p.Establishments[i+1:]...)
			if p.Stats.TotalEstablishments > 0 {
				p.Stats.TotalEstablishments--
			}
			p.MarkUpdated()
			return nil
		}
	}
	return ErrEstablishmentNotFound
}

// GetEstablishment retrieves an establishment by ID.
func (p *Partner) GetEstablishment(estID EstablishmentID) (*Establishment, error) {
	for i := range p.Establishments {
		if p.Establishments[i].ID == estID {
			return &p.Establishments[i], nil
		}
	}
	return nil, ErrEstablishmentNotFound
}

// GetActiveEstablishments returns all active establishments.
func (p *Partner) GetActiveEstablishments() []Establishment {
	active := make([]Establishment, 0)
	for _, est := range p.Establishments {
		if est.IsActive {
			active = append(active, est)
		}
	}
	return active
}

// =============================================================================
// Team Management
// =============================================================================

// AddTeamMember adds a new team member (invitation).
func (p *Partner) AddTeamMember(member TeamMember) error {
	// Check for duplicate email
	for _, existing := range p.TeamMembers {
		if existing.Email.Equals(member.Email) {
			return ErrTeamMemberExists
		}
	}

	member.PartnerID = p.ID
	p.TeamMembers = append(p.TeamMembers, member)
	p.MarkUpdated()
	p.AddDomainEvent(NewTeamMemberInvitedEvent(p.ID, member.Email, member.Role))
	return nil
}

// UpdateTeamMember updates an existing team member.
func (p *Partner) UpdateTeamMember(memberID TeamMemberID, updateFn func(*TeamMember) error) error {
	for i := range p.TeamMembers {
		if p.TeamMembers[i].ID == memberID {
			if err := updateFn(&p.TeamMembers[i]); err != nil {
				return err
			}
			p.MarkUpdated()
			return nil
		}
	}
	return ErrTeamMemberNotFound
}

// RemoveTeamMember removes a team member.
func (p *Partner) RemoveTeamMember(memberID TeamMemberID) error {
	for i, member := range p.TeamMembers {
		if member.ID == memberID {
			p.TeamMembers = append(p.TeamMembers[:i], p.TeamMembers[i+1:]...)
			p.MarkUpdated()
			return nil
		}
	}
	return ErrTeamMemberNotFound
}

// GetTeamMember retrieves a team member by ID.
func (p *Partner) GetTeamMember(memberID TeamMemberID) (*TeamMember, error) {
	for i := range p.TeamMembers {
		if p.TeamMembers[i].ID == memberID {
			return &p.TeamMembers[i], nil
		}
	}
	return nil, ErrTeamMemberNotFound
}

// GetTeamMemberByEmail retrieves a team member by email.
func (p *Partner) GetTeamMemberByEmail(email Email) (*TeamMember, error) {
	for i := range p.TeamMembers {
		if p.TeamMembers[i].Email.Equals(email) {
			return &p.TeamMembers[i], nil
		}
	}
	return nil, ErrTeamMemberNotFound
}

// AcceptTeamInvitation accepts a team invitation and links to user.
func (p *Partner) AcceptTeamInvitation(memberID TeamMemberID, userID UserID) error {
	for i := range p.TeamMembers {
		if p.TeamMembers[i].ID == memberID {
			if p.TeamMembers[i].Status != TeamMemberStatusPending {
				return ErrInvitationNotPending
			}
			now := time.Now()
			p.TeamMembers[i].UserID = &userID
			p.TeamMembers[i].Status = TeamMemberStatusActive
			p.TeamMembers[i].JoinedAt = &now
			p.MarkUpdated()
			p.AddDomainEvent(NewTeamMemberJoinedEvent(p.ID, memberID, userID))
			return nil
		}
	}
	return ErrTeamMemberNotFound
}

// =============================================================================
// Statistics
// =============================================================================

// UpdateStats updates the partner's denormalized statistics.
func (p *Partner) UpdateStats(stats PartnerStats) {
	p.Stats = stats
	p.Stats.LastUpdated = time.Now()
	p.MarkUpdated()
}

// IncrementOfferCount increments the total offers count.
func (p *Partner) IncrementOfferCount() {
	p.Stats.TotalOffers++
	p.Stats.ActiveOffers++
	p.Stats.LastUpdated = time.Now()
}

// DecrementActiveOffers decrements the active offers count.
func (p *Partner) DecrementActiveOffers() {
	if p.Stats.ActiveOffers > 0 {
		p.Stats.ActiveOffers--
	}
	p.Stats.LastUpdated = time.Now()
}

// IncrementBookings increments the total bookings count.
func (p *Partner) IncrementBookings() {
	p.Stats.TotalBookings++
	p.Stats.LastUpdated = time.Now()
}

// IncrementCheckins increments the total check-ins count.
func (p *Partner) IncrementCheckins() {
	p.Stats.TotalCheckins++
	p.Stats.LastUpdated = time.Now()
}

// UpdateRating updates the average rating.
func (p *Partner) UpdateRating(avgRating float64, reviewCount int) {
	p.Stats.AvgRating = avgRating
	p.Stats.ReviewCount = reviewCount
	p.Stats.LastUpdated = time.Now()
}

// =============================================================================
// Soft Delete
// =============================================================================

// Delete soft-deletes the partner.
func (p *Partner) Delete() {
	now := time.Now()
	p.DeletedAt = &now
	p.Status = PartnerStatusSuspended
	p.MarkUpdated()
}

// IsDeleted returns true if the partner is soft-deleted.
func (p *Partner) IsDeleted() bool {
	return p.DeletedAt != nil
}

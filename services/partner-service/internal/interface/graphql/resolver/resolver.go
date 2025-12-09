// Package resolver contains the GraphQL resolvers for the Partner service.
package resolver

import (
	"context"
	"time"

	"github.com/yousoon/services/partner/internal/application/commands"
	"github.com/yousoon/services/partner/internal/application/queries"
	partnerdomain "github.com/yousoon/services/partner/internal/domain"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver is the root resolver for the Partner GraphQL schema.
type Resolver struct {
	// Repositories
	PartnerRepo     partnerdomain.PartnerRepository
	PartnerReadRepo partnerdomain.PartnerReadRepository

	// Command Handlers
	RegisterPartnerHandler      *commands.RegisterPartnerHandler
	UpdatePartnerHandler        *commands.UpdatePartnerHandler
	VerifyPartnerHandler        *commands.VerifyPartnerHandler
	SuspendPartnerHandler       *commands.SuspendPartnerHandler
	AddEstablishmentHandler     *commands.AddEstablishmentHandler
	InviteTeamMemberHandler     *commands.InviteTeamMemberHandler
	AcceptTeamInvitationHandler *commands.AcceptTeamInvitationHandler

	// Query Handlers
	GetPartnerHandler               *queries.GetPartnerHandler
	GetPartnerByOwnerHandler        *queries.GetPartnerByOwnerHandler
	ListPartnersHandler             *queries.ListPartnersHandler
	GetEstablishmentHandler         *queries.GetEstablishmentHandler
	ListEstablishmentsHandler       *queries.ListEstablishmentsHandler
	GetTeamMembersHandler           *queries.GetTeamMembersHandler
	GetPartnersForTeamMemberHandler *queries.GetPartnersForTeamMemberHandler
}

// =============================================================================
// Query Resolvers
// =============================================================================

// Query returns the QueryResolver implementation.
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Partner(ctx context.Context, id string) (*partnerdomain.Partner, error) {
	return r.GetPartnerHandler.Handle(ctx, queries.GetPartnerQuery{PartnerID: id})
}

func (r *queryResolver) PartnerByOwner(ctx context.Context, ownerUserID string) (*partnerdomain.Partner, error) {
	return r.GetPartnerByOwnerHandler.Handle(ctx, queries.GetPartnerByOwnerQuery{OwnerUserID: ownerUserID})
}

func (r *queryResolver) Partners(ctx context.Context, filter *PartnerFilterInput, first *int, after *string) (*PartnerConnection, error) {
	query := queries.ListPartnersQuery{
		Offset: 0,
		Limit:  20,
	}

	if filter != nil {
		if filter.Status != nil {
			status := string(*filter.Status)
			query.Status = &status
		}
		if filter.Category != nil {
			query.Category = filter.Category
		}
		if filter.Search != nil {
			query.Search = filter.Search
		}
	}

	if first != nil && *first > 0 {
		query.Limit = *first
	}

	// TODO: Implement cursor-based pagination
	// if after != nil {
	// 	query.Offset = decodeCursor(*after)
	// }

	result, err := r.ListPartnersHandler.Handle(ctx, query)
	if err != nil {
		return nil, err
	}

	edges := make([]*PartnerEdge, len(result.Partners))
	for i, p := range result.Partners {
		edges[i] = &PartnerEdge{
			Node:   p,
			Cursor: p.ID.String(), // Simple cursor using ID
		}
	}

	return &PartnerConnection{
		Edges:      edges,
		TotalCount: int(result.TotalCount),
		PageInfo: &PageInfo{
			HasNextPage:     result.Offset+result.Limit < int(result.TotalCount),
			HasPreviousPage: result.Offset > 0,
		},
	}, nil
}

func (r *queryResolver) Establishment(ctx context.Context, id string) (*partnerdomain.Establishment, error) {
	// Get partner by establishment ID
	partner, err := r.PartnerReadRepo.GetPartnerByEstablishmentID(ctx, partnerdomain.EstablishmentID(id))
	if err != nil {
		return nil, err
	}
	return partner.GetEstablishment(partnerdomain.EstablishmentID(id))
}

func (r *queryResolver) EstablishmentsNearby(ctx context.Context, input NearbyEstablishmentsInput) ([]*partnerdomain.EstablishmentSummary, error) {
	location, err := partnerdomain.NewGeoLocation(input.Longitude, input.Latitude)
	if err != nil {
		return nil, err
	}

	radiusKm := 10.0 // Default 10km
	if input.RadiusKm != nil {
		radiusKm = *input.RadiusKm
	}

	limit := 20
	if input.Limit != nil {
		limit = *input.Limit
	}

	summaries, err := r.PartnerReadRepo.GetEstablishmentsNearLocation(ctx, location, radiusKm, limit)
	if err != nil {
		return nil, err
	}

	// Convert to pointers
	result := make([]*partnerdomain.EstablishmentSummary, len(summaries))
	for i := range summaries {
		result[i] = &summaries[i]
	}
	return result, nil
}

func (r *queryResolver) SearchPartners(ctx context.Context, query string, limit *int) ([]*partnerdomain.PartnerSummary, error) {
	l := 20
	if limit != nil {
		l = *limit
	}
	return r.PartnerReadRepo.SearchPartners(ctx, query, l)
}

func (r *queryResolver) PartnersForTeamMember(ctx context.Context, email string) ([]*partnerdomain.Partner, error) {
	return r.GetPartnersForTeamMemberHandler.Handle(ctx, queries.GetPartnersForTeamMemberQuery{Email: email})
}

// =============================================================================
// Mutation Resolvers
// =============================================================================

// Mutation returns the MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) RegisterPartner(ctx context.Context, input RegisterPartnerInput) (*partnerdomain.Partner, error) {
	subcategories := []string{}
	if input.Subcategories != nil {
		subcategories = input.Subcategories
	}

	cmd := commands.RegisterPartnerCommand{
		OwnerUserID:      input.OwnerUserID,
		CompanyName:      input.CompanyName,
		TradeName:        derefString(input.TradeName),
		SIRET:            input.Siret,
		VATNumber:        derefString(input.VatNumber),
		LegalForm:        derefString(input.LegalForm),
		ContactFirstName: input.ContactFirstName,
		ContactLastName:  input.ContactLastName,
		ContactEmail:     input.ContactEmail,
		ContactPhone:     derefString(input.ContactPhone),
		ContactRole:      derefString(input.ContactRole),
		Category:         input.Category,
		Subcategories:    subcategories,
	}

	return r.RegisterPartnerHandler.Handle(ctx, cmd)
}

func (r *mutationResolver) UpdatePartner(ctx context.Context, id string, input UpdatePartnerInput) (*partnerdomain.Partner, error) {
	cmd := commands.UpdatePartnerCommand{
		PartnerID:        id,
		CompanyName:      input.CompanyName,
		TradeName:        input.TradeName,
		VATNumber:        input.VatNumber,
		LegalForm:        input.LegalForm,
		Logo:             input.Logo,
		CoverImage:       input.CoverImage,
		PrimaryColor:     input.PrimaryColor,
		Description:      input.Description,
		ContactFirstName: input.ContactFirstName,
		ContactLastName:  input.ContactLastName,
		ContactEmail:     input.ContactEmail,
		ContactPhone:     input.ContactPhone,
		ContactRole:      input.ContactRole,
		Category:         input.Category,
		Subcategories:    input.Subcategories,
	}

	return r.UpdatePartnerHandler.Handle(ctx, cmd)
}

func (r *mutationResolver) VerifyPartner(ctx context.Context, id string) (*partnerdomain.Partner, error) {
	// TODO: Get admin user ID from context
	adminID := "admin" // Placeholder

	cmd := commands.VerifyPartnerCommand{
		PartnerID:    id,
		VerifiedByID: adminID,
	}

	if err := r.VerifyPartnerHandler.Handle(ctx, cmd); err != nil {
		return nil, err
	}

	return r.PartnerRepo.FindByID(ctx, partnerdomain.PartnerID(id))
}

func (r *mutationResolver) SuspendPartner(ctx context.Context, id string, reason string) (*partnerdomain.Partner, error) {
	cmd := commands.SuspendPartnerCommand{
		PartnerID: id,
		Reason:    reason,
	}

	if err := r.SuspendPartnerHandler.Handle(ctx, cmd); err != nil {
		return nil, err
	}

	return r.PartnerRepo.FindByID(ctx, partnerdomain.PartnerID(id))
}

func (r *mutationResolver) DeletePartner(ctx context.Context, id string) (bool, error) {
	if err := r.PartnerRepo.Delete(ctx, partnerdomain.PartnerID(id)); err != nil {
		return false, err
	}
	return true, nil
}

func (r *mutationResolver) AddEstablishment(ctx context.Context, partnerID string, input AddEstablishmentInput) (*partnerdomain.Establishment, error) {
	cmd := commands.AddEstablishmentCommand{
		PartnerID:    partnerID,
		Name:         input.Name,
		Description:  derefString(input.Description),
		Street:       input.Street,
		StreetNumber: derefString(input.StreetNumber),
		Complement:   derefString(input.Complement),
		PostalCode:   input.PostalCode,
		City:         input.City,
		Country:      input.Country,
		Longitude:    input.Longitude,
		Latitude:     input.Latitude,
		Phone:        derefString(input.Phone),
		Email:        derefString(input.Email),
		Website:      derefString(input.Website),
		Type:         derefString(input.Type),
		Features:     input.Features,
		PriceRange:   derefInt(input.PriceRange, 2),
	}

	return r.AddEstablishmentHandler.Handle(ctx, cmd)
}

func (r *mutationResolver) UpdateEstablishment(ctx context.Context, partnerID string, establishmentID string, input UpdateEstablishmentInput) (*partnerdomain.Establishment, error) {
	partner, err := r.PartnerRepo.FindByID(ctx, partnerdomain.PartnerID(partnerID))
	if err != nil {
		return nil, err
	}

	err = partner.UpdateEstablishment(partnerdomain.EstablishmentID(establishmentID), func(est *partnerdomain.Establishment) error {
		if input.Name != nil {
			est.Name = *input.Name
		}
		if input.Description != nil {
			est.Description = *input.Description
		}
		if input.Type != nil {
			est.SetType(*input.Type)
		}
		if input.Features != nil {
			est.SetFeatures(input.Features)
		}
		if input.PriceRange != nil {
			est.SetPriceRange(*input.PriceRange)
		}
		if input.IsActive != nil {
			if *input.IsActive {
				est.Activate()
			} else {
				est.Deactivate()
			}
		}
		// TODO: Handle address and location updates
		return nil
	})

	if err != nil {
		return nil, err
	}

	if err := r.PartnerRepo.Save(ctx, partner); err != nil {
		return nil, err
	}

	return partner.GetEstablishment(partnerdomain.EstablishmentID(establishmentID))
}

func (r *mutationResolver) RemoveEstablishment(ctx context.Context, partnerID string, establishmentID string) (bool, error) {
	partner, err := r.PartnerRepo.FindByID(ctx, partnerdomain.PartnerID(partnerID))
	if err != nil {
		return false, err
	}

	if err := partner.RemoveEstablishment(partnerdomain.EstablishmentID(establishmentID)); err != nil {
		return false, err
	}

	if err := r.PartnerRepo.Save(ctx, partner); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) SetEstablishmentOpeningHours(ctx context.Context, partnerID string, establishmentID string, hours []*OpeningHourInput) (*partnerdomain.Establishment, error) {
	partner, err := r.PartnerRepo.FindByID(ctx, partnerdomain.PartnerID(partnerID))
	if err != nil {
		return nil, err
	}

	openingHours := make([]partnerdomain.OpeningHour, len(hours))
	for i, h := range hours {
		if h.IsClosed != nil && *h.IsClosed {
			openingHours[i] = partnerdomain.NewClosedDay(h.DayOfWeek)
		} else {
			openingHours[i] = partnerdomain.NewOpeningHour(h.DayOfWeek, derefString(h.Open), derefString(h.Close))
		}
	}

	err = partner.UpdateEstablishment(partnerdomain.EstablishmentID(establishmentID), func(est *partnerdomain.Establishment) error {
		est.SetOpeningHours(openingHours)
		return nil
	})

	if err != nil {
		return nil, err
	}

	if err := r.PartnerRepo.Save(ctx, partner); err != nil {
		return nil, err
	}

	return partner.GetEstablishment(partnerdomain.EstablishmentID(establishmentID))
}

func (r *mutationResolver) AddEstablishmentImage(ctx context.Context, partnerID string, establishmentID string, url string, alt *string, isPrimary *bool) (*partnerdomain.Establishment, error) {
	partner, err := r.PartnerRepo.FindByID(ctx, partnerdomain.PartnerID(partnerID))
	if err != nil {
		return nil, err
	}

	err = partner.UpdateEstablishment(partnerdomain.EstablishmentID(establishmentID), func(est *partnerdomain.Establishment) error {
		primary := false
		if isPrimary != nil {
			primary = *isPrimary
		}
		est.AddImage(partnerdomain.NewImage(url, derefString(alt), primary, len(est.Images)))
		return nil
	})

	if err != nil {
		return nil, err
	}

	if err := r.PartnerRepo.Save(ctx, partner); err != nil {
		return nil, err
	}

	return partner.GetEstablishment(partnerdomain.EstablishmentID(establishmentID))
}

func (r *mutationResolver) RemoveEstablishmentImage(ctx context.Context, partnerID string, establishmentID string, url string) (*partnerdomain.Establishment, error) {
	partner, err := r.PartnerRepo.FindByID(ctx, partnerdomain.PartnerID(partnerID))
	if err != nil {
		return nil, err
	}

	err = partner.UpdateEstablishment(partnerdomain.EstablishmentID(establishmentID), func(est *partnerdomain.Establishment) error {
		est.RemoveImage(url)
		return nil
	})

	if err != nil {
		return nil, err
	}

	if err := r.PartnerRepo.Save(ctx, partner); err != nil {
		return nil, err
	}

	return partner.GetEstablishment(partnerdomain.EstablishmentID(establishmentID))
}

func (r *mutationResolver) InviteTeamMember(ctx context.Context, partnerID string, input InviteTeamMemberInput) (*partnerdomain.TeamMember, error) {
	cmd := commands.InviteTeamMemberCommand{
		PartnerID: partnerID,
		Email:     input.Email,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Role:      string(input.Role),
	}

	return r.InviteTeamMemberHandler.Handle(ctx, cmd)
}

func (r *mutationResolver) AcceptTeamInvitation(ctx context.Context, partnerID string, teamMemberID string, userID string) (*partnerdomain.TeamMember, error) {
	cmd := commands.AcceptTeamInvitationCommand{
		PartnerID:    partnerID,
		TeamMemberID: teamMemberID,
		UserID:       userID,
	}

	if err := r.AcceptTeamInvitationHandler.Handle(ctx, cmd); err != nil {
		return nil, err
	}

	partner, err := r.PartnerRepo.FindByID(ctx, partnerdomain.PartnerID(partnerID))
	if err != nil {
		return nil, err
	}

	return partner.GetTeamMember(partnerdomain.TeamMemberID(teamMemberID))
}

func (r *mutationResolver) UpdateTeamMemberRole(ctx context.Context, partnerID string, teamMemberID string, role TeamRole) (*partnerdomain.TeamMember, error) {
	partner, err := r.PartnerRepo.FindByID(ctx, partnerdomain.PartnerID(partnerID))
	if err != nil {
		return nil, err
	}

	domainRole := partnerdomain.TeamRole(string(role))
	err = partner.UpdateTeamMember(partnerdomain.TeamMemberID(teamMemberID), func(m *partnerdomain.TeamMember) error {
		m.UpdateRole(domainRole)
		return nil
	})

	if err != nil {
		return nil, err
	}

	if err := r.PartnerRepo.Save(ctx, partner); err != nil {
		return nil, err
	}

	return partner.GetTeamMember(partnerdomain.TeamMemberID(teamMemberID))
}

func (r *mutationResolver) RemoveTeamMember(ctx context.Context, partnerID string, teamMemberID string) (bool, error) {
	partner, err := r.PartnerRepo.FindByID(ctx, partnerdomain.PartnerID(partnerID))
	if err != nil {
		return false, err
	}

	if err := partner.RemoveTeamMember(partnerdomain.TeamMemberID(teamMemberID)); err != nil {
		return false, err
	}

	if err := r.PartnerRepo.Save(ctx, partner); err != nil {
		return false, err
	}

	return true, nil
}

// =============================================================================
// Field Resolvers
// =============================================================================

// Partner field resolvers
func (r *Resolver) Partner() PartnerResolver {
	return &partnerResolver{r}
}

type partnerResolver struct{ *Resolver }

func (r *partnerResolver) EstablishmentCount(ctx context.Context, obj *partnerdomain.Partner) (int, error) {
	return len(obj.Establishments), nil
}

// Establishment field resolvers
func (r *Resolver) Establishment() EstablishmentResolver {
	return &establishmentResolver{r}
}

type establishmentResolver struct{ *Resolver }

func (r *establishmentResolver) PrimaryImage(ctx context.Context, obj *partnerdomain.Establishment) (*string, error) {
	img := obj.GetPrimaryImage()
	if img == "" {
		return nil, nil
	}
	return &img, nil
}

func (r *establishmentResolver) IsOpenNow(ctx context.Context, obj *partnerdomain.Establishment) (bool, error) {
	return obj.IsOpenAt(time.Now()), nil
}

// GeoLocation field resolvers
func (r *Resolver) GeoLocation() GeoLocationResolver {
	return &geoLocationResolver{r}
}

type geoLocationResolver struct{ *Resolver }

func (r *geoLocationResolver) Longitude(ctx context.Context, obj *partnerdomain.GeoLocation) (float64, error) {
	return obj.Longitude(), nil
}

func (r *geoLocationResolver) Latitude(ctx context.Context, obj *partnerdomain.GeoLocation) (float64, error) {
	return obj.Latitude(), nil
}

// Contact field resolvers
func (r *Resolver) Contact() ContactResolver {
	return &contactResolver{r}
}

type contactResolver struct{ *Resolver }

func (r *contactResolver) FullName(ctx context.Context, obj *partnerdomain.Contact) (string, error) {
	return obj.FullName(), nil
}

// TeamMember field resolvers
func (r *Resolver) TeamMember() TeamMemberResolver {
	return &teamMemberResolver{r}
}

type teamMemberResolver struct{ *Resolver }

func (r *teamMemberResolver) FullName(ctx context.Context, obj *partnerdomain.TeamMember) (string, error) {
	return obj.FullName(), nil
}

func (r *teamMemberResolver) Email(ctx context.Context, obj *partnerdomain.TeamMember) (string, error) {
	return obj.Email.String(), nil
}

// =============================================================================
// Federation Entity Resolvers
// =============================================================================

// Entity returns the EntityResolver implementation.
func (r *Resolver) Entity() EntityResolver {
	return &entityResolver{r}
}

type entityResolver struct{ *Resolver }

func (r *entityResolver) FindPartnerByID(ctx context.Context, id string) (*partnerdomain.Partner, error) {
	return r.PartnerRepo.FindByID(ctx, partnerdomain.PartnerID(id))
}

func (r *entityResolver) FindEstablishmentByID(ctx context.Context, id string) (*partnerdomain.Establishment, error) {
	partner, err := r.PartnerReadRepo.GetPartnerByEstablishmentID(ctx, partnerdomain.EstablishmentID(id))
	if err != nil {
		return nil, err
	}
	return partner.GetEstablishment(partnerdomain.EstablishmentID(id))
}

// User resolver (federated from Identity service)
func (r *Resolver) User() UserResolver {
	return &userResolver{r}
}

type userResolver struct{ *Resolver }

func (r *userResolver) OwnedPartner(ctx context.Context, obj *User) (*partnerdomain.Partner, error) {
	return r.PartnerRepo.FindByOwnerUserID(ctx, partnerdomain.UserID(obj.ID))
}

func (r *userResolver) PartnerMemberships(ctx context.Context, obj *User) ([]*partnerdomain.Partner, error) {
	// TODO: Implement finding partners where user is a team member
	return []*partnerdomain.Partner{}, nil
}

// =============================================================================
// Helper Functions
// =============================================================================

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func derefInt(i *int, defaultVal int) int {
	if i == nil {
		return defaultVal
	}
	return *i
}

// =============================================================================
// GraphQL Types (generated placeholder interfaces)
// =============================================================================

// These interfaces would be generated by gqlgen but are defined here as placeholders

type QueryResolver interface {
	Partner(ctx context.Context, id string) (*partnerdomain.Partner, error)
	PartnerByOwner(ctx context.Context, ownerUserID string) (*partnerdomain.Partner, error)
	Partners(ctx context.Context, filter *PartnerFilterInput, first *int, after *string) (*PartnerConnection, error)
	Establishment(ctx context.Context, id string) (*partnerdomain.Establishment, error)
	EstablishmentsNearby(ctx context.Context, input NearbyEstablishmentsInput) ([]*partnerdomain.EstablishmentSummary, error)
	SearchPartners(ctx context.Context, query string, limit *int) ([]*partnerdomain.PartnerSummary, error)
	PartnersForTeamMember(ctx context.Context, email string) ([]*partnerdomain.Partner, error)
}

type MutationResolver interface {
	RegisterPartner(ctx context.Context, input RegisterPartnerInput) (*partnerdomain.Partner, error)
	UpdatePartner(ctx context.Context, id string, input UpdatePartnerInput) (*partnerdomain.Partner, error)
	VerifyPartner(ctx context.Context, id string) (*partnerdomain.Partner, error)
	SuspendPartner(ctx context.Context, id string, reason string) (*partnerdomain.Partner, error)
	DeletePartner(ctx context.Context, id string) (bool, error)
	AddEstablishment(ctx context.Context, partnerID string, input AddEstablishmentInput) (*partnerdomain.Establishment, error)
	UpdateEstablishment(ctx context.Context, partnerID string, establishmentID string, input UpdateEstablishmentInput) (*partnerdomain.Establishment, error)
	RemoveEstablishment(ctx context.Context, partnerID string, establishmentID string) (bool, error)
	SetEstablishmentOpeningHours(ctx context.Context, partnerID string, establishmentID string, hours []*OpeningHourInput) (*partnerdomain.Establishment, error)
	AddEstablishmentImage(ctx context.Context, partnerID string, establishmentID string, url string, alt *string, isPrimary *bool) (*partnerdomain.Establishment, error)
	RemoveEstablishmentImage(ctx context.Context, partnerID string, establishmentID string, url string) (*partnerdomain.Establishment, error)
	InviteTeamMember(ctx context.Context, partnerID string, input InviteTeamMemberInput) (*partnerdomain.TeamMember, error)
	AcceptTeamInvitation(ctx context.Context, partnerID string, teamMemberID string, userID string) (*partnerdomain.TeamMember, error)
	UpdateTeamMemberRole(ctx context.Context, partnerID string, teamMemberID string, role TeamRole) (*partnerdomain.TeamMember, error)
	RemoveTeamMember(ctx context.Context, partnerID string, teamMemberID string) (bool, error)
}

type PartnerResolver interface {
	EstablishmentCount(ctx context.Context, obj *partnerdomain.Partner) (int, error)
}

type EstablishmentResolver interface {
	PrimaryImage(ctx context.Context, obj *partnerdomain.Establishment) (*string, error)
	IsOpenNow(ctx context.Context, obj *partnerdomain.Establishment) (bool, error)
}

type GeoLocationResolver interface {
	Longitude(ctx context.Context, obj *partnerdomain.GeoLocation) (float64, error)
	Latitude(ctx context.Context, obj *partnerdomain.GeoLocation) (float64, error)
}

type ContactResolver interface {
	FullName(ctx context.Context, obj *partnerdomain.Contact) (string, error)
}

type TeamMemberResolver interface {
	FullName(ctx context.Context, obj *partnerdomain.TeamMember) (string, error)
	Email(ctx context.Context, obj *partnerdomain.TeamMember) (string, error)
}

type EntityResolver interface {
	FindPartnerByID(ctx context.Context, id string) (*partnerdomain.Partner, error)
	FindEstablishmentByID(ctx context.Context, id string) (*partnerdomain.Establishment, error)
}

type UserResolver interface {
	OwnedPartner(ctx context.Context, obj *User) (*partnerdomain.Partner, error)
	PartnerMemberships(ctx context.Context, obj *User) ([]*partnerdomain.Partner, error)
}

// =============================================================================
// Input Types (generated placeholder types)
// =============================================================================

type PartnerFilterInput struct {
	Status   *PartnerStatus
	Category *string
	Search   *string
}

type PartnerStatus string

const (
	PartnerStatusPending   PartnerStatus = "PENDING"
	PartnerStatusActive    PartnerStatus = "ACTIVE"
	PartnerStatusSuspended PartnerStatus = "SUSPENDED"
)

type TeamRole string

const (
	TeamRoleAdmin   TeamRole = "ADMIN"
	TeamRoleManager TeamRole = "MANAGER"
	TeamRoleStaff   TeamRole = "STAFF"
	TeamRoleViewer  TeamRole = "VIEWER"
)

type RegisterPartnerInput struct {
	OwnerUserID      string
	CompanyName      string
	TradeName        *string
	Siret            string
	VatNumber        *string
	LegalForm        *string
	ContactFirstName string
	ContactLastName  string
	ContactEmail     string
	ContactPhone     *string
	ContactRole      *string
	Category         string
	Subcategories    []string
}

type UpdatePartnerInput struct {
	CompanyName      *string
	TradeName        *string
	VatNumber        *string
	LegalForm        *string
	Logo             *string
	CoverImage       *string
	PrimaryColor     *string
	Description      *string
	ContactFirstName *string
	ContactLastName  *string
	ContactEmail     *string
	ContactPhone     *string
	ContactRole      *string
	Category         *string
	Subcategories    *[]string
}

type AddEstablishmentInput struct {
	Name         string
	Description  *string
	Street       string
	StreetNumber *string
	Complement   *string
	PostalCode   string
	City         string
	Country      string
	Longitude    float64
	Latitude     float64
	Phone        *string
	Email        *string
	Website      *string
	Type         *string
	Features     []string
	PriceRange   *int
}

type UpdateEstablishmentInput struct {
	Name         *string
	Description  *string
	Street       *string
	StreetNumber *string
	Complement   *string
	PostalCode   *string
	City         *string
	Country      *string
	Longitude    *float64
	Latitude     *float64
	Phone        *string
	Email        *string
	Website      *string
	Type         *string
	Features     []string
	PriceRange   *int
	IsActive     *bool
}

type OpeningHourInput struct {
	DayOfWeek int
	Open      *string
	Close     *string
	IsClosed  *bool
}

type InviteTeamMemberInput struct {
	Email     string
	FirstName string
	LastName  string
	Role      TeamRole
}

type NearbyEstablishmentsInput struct {
	Longitude float64
	Latitude  float64
	RadiusKm  *float64
	Limit     *int
}

type PartnerConnection struct {
	Edges      []*PartnerEdge
	PageInfo   *PageInfo
	TotalCount int
}

type PartnerEdge struct {
	Node   *partnerdomain.Partner
	Cursor string
}

type PageInfo struct {
	HasNextPage     bool
	HasPreviousPage bool
	StartCursor     *string
	EndCursor       *string
}

type User struct {
	ID string
}

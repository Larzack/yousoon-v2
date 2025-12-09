// Package queries contains the query handlers for the Partner bounded context.
package queries

import (
	"context"

	partnerdomain "github.com/yousoon/services/partner/internal/domain"
)

// =============================================================================
// Get Partner Query
// =============================================================================

// GetPartnerQuery represents a query to get a partner by ID.
type GetPartnerQuery struct {
	PartnerID string `json:"partnerId"`
}

// GetPartnerHandler handles the GetPartnerQuery.
type GetPartnerHandler struct {
	repo partnerdomain.PartnerRepository
}

// NewGetPartnerHandler creates a new handler.
func NewGetPartnerHandler(repo partnerdomain.PartnerRepository) *GetPartnerHandler {
	return &GetPartnerHandler{repo: repo}
}

// Handle handles the query.
func (h *GetPartnerHandler) Handle(ctx context.Context, q GetPartnerQuery) (*partnerdomain.Partner, error) {
	if q.PartnerID == "" {
		return nil, partnerdomain.ErrPartnerNotFound
	}
	return h.repo.FindByID(ctx, partnerdomain.PartnerID(q.PartnerID))
}

// =============================================================================
// Get Partner By Owner Query
// =============================================================================

// GetPartnerByOwnerQuery represents a query to get a partner by owner user ID.
type GetPartnerByOwnerQuery struct {
	OwnerUserID string `json:"ownerUserId"`
}

// GetPartnerByOwnerHandler handles the GetPartnerByOwnerQuery.
type GetPartnerByOwnerHandler struct {
	repo partnerdomain.PartnerRepository
}

// NewGetPartnerByOwnerHandler creates a new handler.
func NewGetPartnerByOwnerHandler(repo partnerdomain.PartnerRepository) *GetPartnerByOwnerHandler {
	return &GetPartnerByOwnerHandler{repo: repo}
}

// Handle handles the query.
func (h *GetPartnerByOwnerHandler) Handle(ctx context.Context, q GetPartnerByOwnerQuery) (*partnerdomain.Partner, error) {
	if q.OwnerUserID == "" {
		return nil, partnerdomain.ErrPartnerNotFound
	}
	return h.repo.FindByOwnerUserID(ctx, partnerdomain.UserID(q.OwnerUserID))
}

// =============================================================================
// List Partners Query
// =============================================================================

// ListPartnersQuery represents a query to list partners.
type ListPartnersQuery struct {
	Status   *string `json:"status,omitempty"`
	Category *string `json:"category,omitempty"`
	Search   *string `json:"search,omitempty"`
	Offset   int     `json:"offset"`
	Limit    int     `json:"limit"`
}

// ListPartnersResult contains the result of a list partners query.
type ListPartnersResult struct {
	Partners   []*partnerdomain.Partner `json:"partners"`
	TotalCount int64                    `json:"totalCount"`
	Offset     int                      `json:"offset"`
	Limit      int                      `json:"limit"`
}

// ListPartnersHandler handles the ListPartnersQuery.
type ListPartnersHandler struct {
	repo partnerdomain.PartnerRepository
}

// NewListPartnersHandler creates a new handler.
func NewListPartnersHandler(repo partnerdomain.PartnerRepository) *ListPartnersHandler {
	return &ListPartnersHandler{repo: repo}
}

// Handle handles the query.
func (h *ListPartnersHandler) Handle(ctx context.Context, q ListPartnersQuery) (*ListPartnersResult, error) {
	// Build filter
	filter := partnerdomain.NewPartnerFilter()

	if q.Status != nil {
		status := partnerdomain.PartnerStatus(*q.Status)
		if status.IsValid() {
			filter = filter.WithStatus(status)
		}
	}
	if q.Category != nil {
		filter = filter.WithCategory(*q.Category)
	}
	if q.Search != nil {
		filter = filter.WithSearch(*q.Search)
	}
	if q.Limit > 0 {
		filter = filter.WithPagination(q.Offset, q.Limit)
	}

	// Execute query
	partners, total, err := h.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	return &ListPartnersResult{
		Partners:   partners,
		TotalCount: total,
		Offset:     filter.Offset,
		Limit:      filter.Limit,
	}, nil
}

// =============================================================================
// Get Establishment Query
// =============================================================================

// GetEstablishmentQuery represents a query to get an establishment.
type GetEstablishmentQuery struct {
	PartnerID       string `json:"partnerId"`
	EstablishmentID string `json:"establishmentId"`
}

// GetEstablishmentHandler handles the GetEstablishmentQuery.
type GetEstablishmentHandler struct {
	repo partnerdomain.PartnerRepository
}

// NewGetEstablishmentHandler creates a new handler.
func NewGetEstablishmentHandler(repo partnerdomain.PartnerRepository) *GetEstablishmentHandler {
	return &GetEstablishmentHandler{repo: repo}
}

// Handle handles the query.
func (h *GetEstablishmentHandler) Handle(ctx context.Context, q GetEstablishmentQuery) (*partnerdomain.Establishment, error) {
	if q.PartnerID == "" || q.EstablishmentID == "" {
		return nil, partnerdomain.ErrEstablishmentNotFound
	}

	partner, err := h.repo.FindByID(ctx, partnerdomain.PartnerID(q.PartnerID))
	if err != nil {
		return nil, err
	}

	return partner.GetEstablishment(partnerdomain.EstablishmentID(q.EstablishmentID))
}

// =============================================================================
// List Establishments Query
// =============================================================================

// ListEstablishmentsQuery represents a query to list establishments for a partner.
type ListEstablishmentsQuery struct {
	PartnerID  string `json:"partnerId"`
	ActiveOnly bool   `json:"activeOnly"`
}

// ListEstablishmentsHandler handles the ListEstablishmentsQuery.
type ListEstablishmentsHandler struct {
	repo partnerdomain.PartnerRepository
}

// NewListEstablishmentsHandler creates a new handler.
func NewListEstablishmentsHandler(repo partnerdomain.PartnerRepository) *ListEstablishmentsHandler {
	return &ListEstablishmentsHandler{repo: repo}
}

// Handle handles the query.
func (h *ListEstablishmentsHandler) Handle(ctx context.Context, q ListEstablishmentsQuery) ([]partnerdomain.Establishment, error) {
	if q.PartnerID == "" {
		return nil, partnerdomain.ErrPartnerNotFound
	}

	partner, err := h.repo.FindByID(ctx, partnerdomain.PartnerID(q.PartnerID))
	if err != nil {
		return nil, err
	}

	if q.ActiveOnly {
		return partner.GetActiveEstablishments(), nil
	}

	return partner.Establishments, nil
}

// =============================================================================
// Get Team Members Query
// =============================================================================

// GetTeamMembersQuery represents a query to get team members for a partner.
type GetTeamMembersQuery struct {
	PartnerID string `json:"partnerId"`
}

// GetTeamMembersHandler handles the GetTeamMembersQuery.
type GetTeamMembersHandler struct {
	repo partnerdomain.PartnerRepository
}

// NewGetTeamMembersHandler creates a new handler.
func NewGetTeamMembersHandler(repo partnerdomain.PartnerRepository) *GetTeamMembersHandler {
	return &GetTeamMembersHandler{repo: repo}
}

// Handle handles the query.
func (h *GetTeamMembersHandler) Handle(ctx context.Context, q GetTeamMembersQuery) ([]partnerdomain.TeamMember, error) {
	if q.PartnerID == "" {
		return nil, partnerdomain.ErrPartnerNotFound
	}

	partner, err := h.repo.FindByID(ctx, partnerdomain.PartnerID(q.PartnerID))
	if err != nil {
		return nil, err
	}

	return partner.TeamMembers, nil
}

// =============================================================================
// Get Partners For Team Member Query
// =============================================================================

// GetPartnersForTeamMemberQuery represents a query to get partners for a team member email.
type GetPartnersForTeamMemberQuery struct {
	Email string `json:"email"`
}

// GetPartnersForTeamMemberHandler handles the GetPartnersForTeamMemberQuery.
type GetPartnersForTeamMemberHandler struct {
	repo partnerdomain.PartnerRepository
}

// NewGetPartnersForTeamMemberHandler creates a new handler.
func NewGetPartnersForTeamMemberHandler(repo partnerdomain.PartnerRepository) *GetPartnersForTeamMemberHandler {
	return &GetPartnersForTeamMemberHandler{repo: repo}
}

// Handle handles the query.
func (h *GetPartnersForTeamMemberHandler) Handle(ctx context.Context, q GetPartnersForTeamMemberQuery) ([]*partnerdomain.Partner, error) {
	if q.Email == "" {
		return nil, partnerdomain.ErrInvalidEmail
	}

	email, err := partnerdomain.NewEmail(q.Email)
	if err != nil {
		return nil, err
	}

	return h.repo.FindByTeamMemberEmail(ctx, email)
}

package queries

import (
	"context"
	"fmt"

	"github.com/yousoon/apps/services/booking-service/internal/domain"
)

// =============================================================================
// GET OUTING BY ID
// =============================================================================

type GetOutingQuery struct {
	OutingID string
}

type GetOutingResult struct {
	Outing *domain.Outing
}

type GetOutingHandler struct {
	outingRepo domain.OutingRepository
}

func NewGetOutingHandler(outingRepo domain.OutingRepository) *GetOutingHandler {
	return &GetOutingHandler{
		outingRepo: outingRepo,
	}
}

func (h *GetOutingHandler) Handle(ctx context.Context, query GetOutingQuery) (*GetOutingResult, error) {
	outing, err := h.outingRepo.GetByID(ctx, query.OutingID)
	if err != nil {
		return nil, fmt.Errorf("failed to get outing: %w", err)
	}
	return &GetOutingResult{Outing: outing}, nil
}

// =============================================================================
// GET OUTING BY QR CODE
// =============================================================================

type GetOutingByQRQuery struct {
	QRCode string
}

type GetOutingByQRResult struct {
	Outing *domain.Outing
}

type GetOutingByQRHandler struct {
	outingRepo domain.OutingRepository
}

func NewGetOutingByQRHandler(outingRepo domain.OutingRepository) *GetOutingByQRHandler {
	return &GetOutingByQRHandler{
		outingRepo: outingRepo,
	}
}

func (h *GetOutingByQRHandler) Handle(ctx context.Context, query GetOutingByQRQuery) (*GetOutingByQRResult, error) {
	outing, err := h.outingRepo.GetByQRCode(ctx, query.QRCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get outing by QR: %w", err)
	}
	return &GetOutingByQRResult{Outing: outing}, nil
}

// =============================================================================
// LIST USER OUTINGS
// =============================================================================

type ListUserOutingsQuery struct {
	UserID string
	Filter domain.OutingFilter
}

type ListUserOutingsResult struct {
	Outings    []*domain.Outing
	TotalCount int64
}

type ListUserOutingsHandler struct {
	outingRepo domain.OutingRepository
}

func NewListUserOutingsHandler(outingRepo domain.OutingRepository) *ListUserOutingsHandler {
	return &ListUserOutingsHandler{
		outingRepo: outingRepo,
	}
}

func (h *ListUserOutingsHandler) Handle(ctx context.Context, query ListUserOutingsQuery) (*ListUserOutingsResult, error) {
	outings, count, err := h.outingRepo.GetByUserID(ctx, query.UserID, query.Filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list user outings: %w", err)
	}
	return &ListUserOutingsResult{
		Outings:    outings,
		TotalCount: count,
	}, nil
}

// =============================================================================
// LIST PARTNER OUTINGS
// =============================================================================

type ListPartnerOutingsQuery struct {
	PartnerID string
	Filter    domain.OutingFilter
}

type ListPartnerOutingsResult struct {
	Outings    []*domain.Outing
	TotalCount int64
}

type ListPartnerOutingsHandler struct {
	outingRepo domain.OutingRepository
}

func NewListPartnerOutingsHandler(outingRepo domain.OutingRepository) *ListPartnerOutingsHandler {
	return &ListPartnerOutingsHandler{
		outingRepo: outingRepo,
	}
}

func (h *ListPartnerOutingsHandler) Handle(ctx context.Context, query ListPartnerOutingsQuery) (*ListPartnerOutingsResult, error) {
	outings, count, err := h.outingRepo.GetByPartnerID(ctx, query.PartnerID, query.Filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list partner outings: %w", err)
	}
	return &ListPartnerOutingsResult{
		Outings:    outings,
		TotalCount: count,
	}, nil
}

// =============================================================================
// LIST ESTABLISHMENT OUTINGS
// =============================================================================

type ListEstablishmentOutingsQuery struct {
	EstablishmentID string
	Filter          domain.OutingFilter
}

type ListEstablishmentOutingsResult struct {
	Outings    []*domain.Outing
	TotalCount int64
}

type ListEstablishmentOutingsHandler struct {
	outingRepo domain.OutingRepository
}

func NewListEstablishmentOutingsHandler(outingRepo domain.OutingRepository) *ListEstablishmentOutingsHandler {
	return &ListEstablishmentOutingsHandler{
		outingRepo: outingRepo,
	}
}

func (h *ListEstablishmentOutingsHandler) Handle(ctx context.Context, query ListEstablishmentOutingsQuery) (*ListEstablishmentOutingsResult, error) {
	outings, count, err := h.outingRepo.GetByEstablishmentID(ctx, query.EstablishmentID, query.Filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list establishment outings: %w", err)
	}
	return &ListEstablishmentOutingsResult{
		Outings:    outings,
		TotalCount: count,
	}, nil
}

// =============================================================================
// GET BOOKING STATS
// =============================================================================

type GetBookingStatsQuery struct {
	PartnerID       string
	EstablishmentID string
	OfferID         string
}

type BookingStats struct {
	TotalBookings      int64
	TotalCheckIns      int64
	TotalCancelled     int64
	TotalExpired       int64
	ConversionRate     float64
	AverageCheckInTime float64
}

type GetBookingStatsResult struct {
	Stats *BookingStats
}

type GetBookingStatsHandler struct {
	outingRepo domain.OutingRepository
}

func NewGetBookingStatsHandler(outingRepo domain.OutingRepository) *GetBookingStatsHandler {
	return &GetBookingStatsHandler{
		outingRepo: outingRepo,
	}
}

func (h *GetBookingStatsHandler) Handle(ctx context.Context, query GetBookingStatsQuery) (*GetBookingStatsResult, error) {
	// This would typically involve aggregation queries
	// For now, return placeholder stats
	stats := &BookingStats{
		TotalBookings:      0,
		TotalCheckIns:      0,
		TotalCancelled:     0,
		TotalExpired:       0,
		ConversionRate:     0.0,
		AverageCheckInTime: 0.0,
	}

	return &GetBookingStatsResult{Stats: stats}, nil
}

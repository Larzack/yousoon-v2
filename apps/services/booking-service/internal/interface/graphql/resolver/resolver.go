package resolver

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"

	"github.com/yousoon/apps/services/booking-service/internal/application/commands"
	"github.com/yousoon/apps/services/booking-service/internal/application/queries"
	"github.com/yousoon/apps/services/booking-service/internal/domain"
	"github.com/yousoon/apps/services/booking-service/internal/interface/graphql/model"
)

// =============================================================================
// RESOLVER
// =============================================================================

type Resolver struct {
	// Command handlers
	bookOutingHandler   *commands.BookOutingHandler
	checkInHandler      *commands.CheckInOutingHandler
	cancelOutingHandler *commands.CancelOutingHandler

	// Query handlers
	getOutingHandler                *queries.GetOutingHandler
	getOutingByQRHandler            *queries.GetOutingByQRHandler
	listUserOutingsHandler          *queries.ListUserOutingsHandler
	listPartnerOutingsHandler       *queries.ListPartnerOutingsHandler
	listEstablishmentOutingsHandler *queries.ListEstablishmentOutingsHandler
	getBookingStatsHandler          *queries.GetBookingStatsHandler
}

func NewResolver(
	bookOutingHandler *commands.BookOutingHandler,
	checkInHandler *commands.CheckInOutingHandler,
	cancelOutingHandler *commands.CancelOutingHandler,
	getOutingHandler *queries.GetOutingHandler,
	getOutingByQRHandler *queries.GetOutingByQRHandler,
	listUserOutingsHandler *queries.ListUserOutingsHandler,
	listPartnerOutingsHandler *queries.ListPartnerOutingsHandler,
	listEstablishmentOutingsHandler *queries.ListEstablishmentOutingsHandler,
	getBookingStatsHandler *queries.GetBookingStatsHandler,
) *Resolver {
	return &Resolver{
		bookOutingHandler:               bookOutingHandler,
		checkInHandler:                  checkInHandler,
		cancelOutingHandler:             cancelOutingHandler,
		getOutingHandler:                getOutingHandler,
		getOutingByQRHandler:            getOutingByQRHandler,
		listUserOutingsHandler:          listUserOutingsHandler,
		listPartnerOutingsHandler:       listPartnerOutingsHandler,
		listEstablishmentOutingsHandler: listEstablishmentOutingsHandler,
		getBookingStatsHandler:          getBookingStatsHandler,
	}
}

// =============================================================================
// QUERY RESOLVERS
// =============================================================================

func (r *Resolver) Outing(ctx context.Context, id string) (*model.Outing, error) {
	result, err := r.getOutingHandler.Handle(ctx, queries.GetOutingQuery{OutingID: id})
	if err != nil {
		return nil, err
	}
	return mapOutingToModel(result.Outing), nil
}

func (r *Resolver) OutingByQRCode(ctx context.Context, qrCode string) (*model.Outing, error) {
	result, err := r.getOutingByQRHandler.Handle(ctx, queries.GetOutingByQRQuery{QRCode: qrCode})
	if err != nil {
		return nil, err
	}
	return mapOutingToModel(result.Outing), nil
}

func (r *Resolver) MyOutings(ctx context.Context, filter *model.OutingFilterInput, pagination *model.PaginationInput) (*model.OutingConnection, error) {
	// Get user ID from context
	userID := getUserIDFromContext(ctx)
	if userID == "" {
		return nil, fmt.Errorf("unauthorized")
	}

	domainFilter := mapFilterToDomain(filter, pagination)

	result, err := r.listUserOutingsHandler.Handle(ctx, queries.ListUserOutingsQuery{
		UserID: userID,
		Filter: domainFilter,
	})
	if err != nil {
		return nil, err
	}

	return buildOutingConnection(result.Outings, result.TotalCount, domainFilter), nil
}

func (r *Resolver) PartnerOutings(ctx context.Context, partnerID string, filter *model.OutingFilterInput, pagination *model.PaginationInput) (*model.OutingConnection, error) {
	domainFilter := mapFilterToDomain(filter, pagination)

	result, err := r.listPartnerOutingsHandler.Handle(ctx, queries.ListPartnerOutingsQuery{
		PartnerID: partnerID,
		Filter:    domainFilter,
	})
	if err != nil {
		return nil, err
	}

	return buildOutingConnection(result.Outings, result.TotalCount, domainFilter), nil
}

func (r *Resolver) EstablishmentOutings(ctx context.Context, establishmentID string, filter *model.OutingFilterInput, pagination *model.PaginationInput) (*model.OutingConnection, error) {
	domainFilter := mapFilterToDomain(filter, pagination)

	result, err := r.listEstablishmentOutingsHandler.Handle(ctx, queries.ListEstablishmentOutingsQuery{
		EstablishmentID: establishmentID,
		Filter:          domainFilter,
	})
	if err != nil {
		return nil, err
	}

	return buildOutingConnection(result.Outings, result.TotalCount, domainFilter), nil
}

func (r *Resolver) BookingStats(ctx context.Context, partnerID, establishmentID, offerID *string) (*model.BookingStats, error) {
	query := queries.GetBookingStatsQuery{}
	if partnerID != nil {
		query.PartnerID = *partnerID
	}
	if establishmentID != nil {
		query.EstablishmentID = *establishmentID
	}
	if offerID != nil {
		query.OfferID = *offerID
	}

	result, err := r.getBookingStatsHandler.Handle(ctx, query)
	if err != nil {
		return nil, err
	}

	return &model.BookingStats{
		TotalBookings:      int(result.Stats.TotalBookings),
		TotalCheckIns:      int(result.Stats.TotalCheckIns),
		TotalCancelled:     int(result.Stats.TotalCancelled),
		TotalExpired:       int(result.Stats.TotalExpired),
		ConversionRate:     result.Stats.ConversionRate,
		AverageCheckInTime: result.Stats.AverageCheckInTime,
	}, nil
}

// =============================================================================
// MUTATION RESOLVERS
// =============================================================================

func (r *Resolver) BookOffer(ctx context.Context, input model.BookOfferInput) (*model.BookOfferPayload, error) {
	userID := getUserIDFromContext(ctx)
	if userID == "" {
		return &model.BookOfferPayload{
			Success: false,
			Error: &model.BookingError{
				Code:    model.BookingErrorCodeUserNotVerified,
				Message: "User not authenticated",
			},
		}, nil
	}

	result, err := r.bookOutingHandler.Handle(ctx, commands.BookOutingCommand{
		UserID:  userID,
		OfferID: input.OfferID,
	})
	if err != nil {
		return &model.BookOfferPayload{
			Success: false,
			Error:   mapBookingError(err),
		}, nil
	}

	return &model.BookOfferPayload{
		Success: true,
		Outing:  mapOutingToModel(result.Outing),
	}, nil
}

func (r *Resolver) CheckInOuting(ctx context.Context, input model.CheckInInput) (*model.CheckInPayload, error) {
	staffUserID := getUserIDFromContext(ctx)

	result, err := r.checkInHandler.Handle(ctx, commands.CheckInOutingCommand{
		QRCode:      input.QRCode,
		StaffUserID: staffUserID,
		Latitude:    input.Latitude,
		Longitude:   input.Longitude,
	})
	if err != nil {
		return &model.CheckInPayload{
			Success: false,
			Error:   mapCheckInError(err),
		}, nil
	}

	return &model.CheckInPayload{
		Success: true,
		Outing:  mapOutingToModel(result.Outing),
	}, nil
}

func (r *Resolver) ManualCheckIn(ctx context.Context, input model.ManualCheckInInput) (*model.CheckInPayload, error) {
	staffUserID := getUserIDFromContext(ctx)

	result, err := r.checkInHandler.Handle(ctx, commands.CheckInOutingCommand{
		OutingID:    input.OutingID,
		StaffUserID: staffUserID,
		Latitude:    input.Latitude,
		Longitude:   input.Longitude,
	})
	if err != nil {
		return &model.CheckInPayload{
			Success: false,
			Error:   mapCheckInError(err),
		}, nil
	}

	return &model.CheckInPayload{
		Success: true,
		Outing:  mapOutingToModel(result.Outing),
	}, nil
}

func (r *Resolver) CancelOuting(ctx context.Context, input model.CancelOutingInput) (*model.CancelOutingPayload, error) {
	reason := ""
	if input.Reason != nil {
		reason = *input.Reason
	}

	result, err := r.cancelOutingHandler.Handle(ctx, commands.CancelOutingCommand{
		OutingID:    input.OutingID,
		CancelledBy: domain.CancellationActorUser,
		Reason:      reason,
	})
	if err != nil {
		return &model.CancelOutingPayload{
			Success: false,
			Error:   mapCancellationError(err),
		}, nil
	}

	return &model.CancelOutingPayload{
		Success: true,
		Outing:  mapOutingToModel(result.Outing),
	}, nil
}

// =============================================================================
// FEDERATION RESOLVERS
// =============================================================================

func (r *Resolver) FindOutingByID(ctx context.Context, id string) (*model.Outing, error) {
	return r.Outing(ctx, id)
}

func (r *Resolver) UserOutings(ctx context.Context, obj *model.User, filter *model.OutingFilterInput, pagination *model.PaginationInput) (*model.OutingConnection, error) {
	domainFilter := mapFilterToDomain(filter, pagination)

	result, err := r.listUserOutingsHandler.Handle(ctx, queries.ListUserOutingsQuery{
		UserID: obj.ID,
		Filter: domainFilter,
	})
	if err != nil {
		return nil, err
	}

	return buildOutingConnection(result.Outings, result.TotalCount, domainFilter), nil
}

func (r *Resolver) OfferBookings(ctx context.Context, obj *model.Offer, filter *model.OutingFilterInput, pagination *model.PaginationInput) (*model.OutingConnection, error) {
	domainFilter := mapFilterToDomain(filter, pagination)

	result, err := r.listPartnerOutingsHandler.Handle(ctx, queries.ListPartnerOutingsQuery{
		PartnerID: obj.ID, // This should be offer ID, need separate handler
		Filter:    domainFilter,
	})
	if err != nil {
		return nil, err
	}

	return buildOutingConnection(result.Outings, result.TotalCount, domainFilter), nil
}

func (r *Resolver) OfferActiveBookingsCount(ctx context.Context, obj *model.Offer) (int, error) {
	// Would need a count query
	return 0, nil
}

// =============================================================================
// HELPERS
// =============================================================================

func getUserIDFromContext(ctx context.Context) string {
	// Extract user ID from JWT context
	// This would be set by auth middleware
	if userID, ok := ctx.Value("user_id").(string); ok {
		return userID
	}
	return ""
}

func mapFilterToDomain(filter *model.OutingFilterInput, pagination *model.PaginationInput) domain.OutingFilter {
	df := domain.DefaultOutingFilter()

	if filter != nil {
		if len(filter.Status) > 0 {
			for _, s := range filter.Status {
				df.Status = append(df.Status, domain.OutingStatus(s.String()))
			}
		}
		if filter.StartDate != nil {
			df.StartDate = filter.StartDate
		}
		if filter.EndDate != nil {
			df.EndDate = filter.EndDate
		}
	}

	if pagination != nil {
		if pagination.First != nil {
			df.Limit = *pagination.First
		}
		if pagination.After != nil {
			// Decode cursor to offset
			if offset, err := decodeCursor(*pagination.After); err == nil {
				df.Offset = offset
			}
		}
	}

	return df
}

func buildOutingConnection(outings []*domain.Outing, total int64, filter domain.OutingFilter) *model.OutingConnection {
	edges := make([]*model.OutingEdge, len(outings))
	for i, o := range outings {
		edges[i] = &model.OutingEdge{
			Node:   mapOutingToModel(o),
			Cursor: encodeCursor(filter.Offset + i),
		}
	}

	hasNextPage := int64(filter.Offset+len(outings)) < total
	hasPrevPage := filter.Offset > 0

	var startCursor, endCursor *string
	if len(edges) > 0 {
		startCursor = &edges[0].Cursor
		endCursor = &edges[len(edges)-1].Cursor
	}

	return &model.OutingConnection{
		Edges: edges,
		PageInfo: &model.PageInfo{
			HasNextPage:     hasNextPage,
			HasPreviousPage: hasPrevPage,
			StartCursor:     startCursor,
			EndCursor:       endCursor,
		},
		TotalCount: int(total),
	}
}

func encodeCursor(offset int) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("cursor:%d", offset)))
}

func decodeCursor(cursor string) (int, error) {
	data, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return 0, err
	}
	var offset int
	_, err = fmt.Sscanf(string(data), "cursor:%d", &offset)
	return offset, err
}

func mapOutingToModel(o *domain.Outing) *model.Outing {
	if o == nil {
		return nil
	}

	desc := o.Offer().Description()
	category := o.Offer().Category()
	imageURL := o.Offer().ImageURL()

	outing := &model.Outing{
		ID:     o.ID(),
		UserID: o.UserID(),
		OfferSnapshot: &model.OfferSnapshot{
			OfferID:              o.Offer().OfferID(),
			PartnerID:            o.Offer().PartnerID(),
			EstablishmentID:      o.Offer().EstablishmentID(),
			Title:                o.Offer().Title(),
			Description:          &desc,
			DiscountType:         o.Offer().DiscountType(),
			DiscountValue:        o.Offer().DiscountValue(),
			Category:             &category,
			EstablishmentName:    o.Offer().EstablishmentName(),
			EstablishmentAddress: o.Offer().EstablishmentAddress(),
			Latitude:             o.Offer().Latitude(),
			Longitude:            o.Offer().Longitude(),
			ImageURL:             &imageURL,
			CapturedAt:           o.Offer().CapturedAt(),
		},
		QRCode: &model.QRCodeInfo{
			Code:      o.QRCode().Code(),
			FullCode:  o.QRCode().FullCode(),
			ExpiresAt: o.QRCode().ExpiresAt(),
			IsExpired: o.QRCode().IsExpired(),
		},
		Status:    model.OutingStatus(o.Status()),
		BookedAt:  o.BookedAt(),
		ExpiresAt: o.ExpiresAt(),
		CreatedAt: o.CreatedAt(),
		UpdatedAt: o.UpdatedAt(),
	}

	// Map timeline
	for _, entry := range o.Timeline() {
		outing.Timeline = append(outing.Timeline, &model.TimelineEntry{
			Status:    model.OutingStatus(entry.Status()),
			Timestamp: entry.Timestamp(),
			Actor:     entry.Actor(),
			Metadata:  entry.Metadata(),
		})
	}

	// Map check-in
	if o.CheckIn() != nil {
		outing.CheckIn = &model.CheckInInfo{
			CheckedInAt: o.CheckIn().CheckedInAt(),
			CheckedInBy: o.CheckIn().CheckedInBy(),
			Method:      model.CheckInMethod(o.CheckIn().Method()),
			Latitude:    o.CheckIn().Latitude(),
			Longitude:   o.CheckIn().Longitude(),
		}
	}

	// Map cancellation
	if o.Cancellation() != nil {
		reason := o.Cancellation().Reason()
		outing.Cancellation = &model.CancellationInfo{
			CancelledAt: o.Cancellation().CancelledAt(),
			CancelledBy: model.CancellationActor(o.Cancellation().CancelledBy()),
			Reason:      &reason,
		}
	}

	return outing
}

func mapBookingError(err error) *model.BookingError {
	switch err {
	case domain.ErrOfferNotBookable:
		return &model.BookingError{
			Code:    model.BookingErrorCodeOfferNotAvailable,
			Message: err.Error(),
		}
	case domain.ErrUserQuotaExceeded:
		return &model.BookingError{
			Code:    model.BookingErrorCodeUserQuotaExceeded,
			Message: err.Error(),
		}
	case domain.ErrOfferQuotaExceeded:
		return &model.BookingError{
			Code:    model.BookingErrorCodeOfferQuotaExceeded,
			Message: err.Error(),
		}
	case domain.ErrOutingAlreadyExists:
		return &model.BookingError{
			Code:    model.BookingErrorCodeAlreadyBooked,
			Message: err.Error(),
		}
	default:
		return &model.BookingError{
			Code:    model.BookingErrorCodeInternalError,
			Message: err.Error(),
		}
	}
}

func mapCheckInError(err error) *model.CheckInError {
	switch err {
	case domain.ErrOutingNotFound:
		return &model.CheckInError{
			Code:    model.CheckInErrorCodeOutingNotFound,
			Message: err.Error(),
		}
	case domain.ErrInvalidQRCode:
		return &model.CheckInError{
			Code:    model.CheckInErrorCodeInvalidQRCode,
			Message: err.Error(),
		}
	case domain.ErrOutingExpired:
		return &model.CheckInError{
			Code:    model.CheckInErrorCodeOutingExpired,
			Message: err.Error(),
		}
	case domain.ErrOutingAlreadyUsed:
		return &model.CheckInError{
			Code:    model.CheckInErrorCodeAlreadyCheckedIn,
			Message: err.Error(),
		}
	case domain.ErrOutingCancelled:
		return &model.CheckInError{
			Code:    model.CheckInErrorCodeOutingCancelled,
			Message: err.Error(),
		}
	default:
		return &model.CheckInError{
			Code:    model.CheckInErrorCodeInternalError,
			Message: err.Error(),
		}
	}
}

func mapCancellationError(err error) *model.CancellationError {
	switch err {
	case domain.ErrOutingNotFound:
		return &model.CancellationError{
			Code:    model.CancellationErrorCodeOutingNotFound,
			Message: err.Error(),
		}
	case domain.ErrOutingCancelled:
		return &model.CancellationError{
			Code:    model.CancellationErrorCodeAlreadyCancelled,
			Message: err.Error(),
		}
	case domain.ErrCannotCancelUsed:
		return &model.CancellationError{
			Code:    model.CancellationErrorCodeAlreadyCheckedIn,
			Message: err.Error(),
		}
	default:
		return &model.CancellationError{
			Code:    model.CancellationErrorCodeInternalError,
			Message: err.Error(),
		}
	}
}

// Unused import fix
var _ = strconv.Itoa

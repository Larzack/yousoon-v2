package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/yousoon/apps/services/booking-service/internal/domain"
)

// =============================================================================
// BOOK OUTING COMMAND
// =============================================================================

type BookOutingCommand struct {
	UserID  string
	OfferID string
}

type BookOutingResult struct {
	Outing *domain.Outing
}

type BookOutingHandler struct {
	outingRepo     domain.OutingRepository
	offerService   domain.OfferService
	userService    domain.UserService
	notifyService  domain.NotificationService
	expirationMins int
}

func NewBookOutingHandler(
	outingRepo domain.OutingRepository,
	offerService domain.OfferService,
	userService domain.UserService,
	notifyService domain.NotificationService,
	expirationMins int,
) *BookOutingHandler {
	return &BookOutingHandler{
		outingRepo:     outingRepo,
		offerService:   offerService,
		userService:    userService,
		notifyService:  notifyService,
		expirationMins: expirationMins,
	}
}

func (h *BookOutingHandler) Handle(ctx context.Context, cmd BookOutingCommand) (*BookOutingResult, error) {
	// 1. Check if user can book
	if err := h.userService.CanBook(ctx, cmd.UserID); err != nil {
		return nil, fmt.Errorf("user cannot book: %w", err)
	}

	// 2. Check if offer can be booked
	if err := h.offerService.CanBook(ctx, cmd.OfferID); err != nil {
		return nil, fmt.Errorf("offer cannot be booked: %w", err)
	}

	// 3. Check if user already has an active booking for this offer
	existing, err := h.outingRepo.GetActiveByUserAndOffer(ctx, cmd.UserID, cmd.OfferID)
	if err != nil && err != domain.ErrOutingNotFound {
		return nil, fmt.Errorf("failed to check existing booking: %w", err)
	}
	if existing != nil {
		return nil, domain.ErrOutingAlreadyExists
	}

	// 4. Get offer snapshot
	offerSnapshot, err := h.offerService.GetOfferSnapshot(ctx, cmd.OfferID)
	if err != nil {
		return nil, fmt.Errorf("failed to get offer snapshot: %w", err)
	}

	// 5. Get user snapshot
	userSnapshot, err := h.userService.GetUserSnapshot(ctx, cmd.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user snapshot: %w", err)
	}

	// 6. Create outing
	outing, err := domain.NewOuting(cmd.UserID, *offerSnapshot, *userSnapshot, h.expirationMins)
	if err != nil {
		return nil, fmt.Errorf("failed to create outing: %w", err)
	}

	// 7. Persist outing
	if err := h.outingRepo.Create(ctx, outing); err != nil {
		return nil, fmt.Errorf("failed to save outing: %w", err)
	}

	// 8. Increment offer booking count
	if err := h.offerService.IncrementBookingCount(ctx, cmd.OfferID); err != nil {
		// Log warning but don't fail
		fmt.Printf("warning: failed to increment booking count: %v\n", err)
	}

	// 9. Send notification (async, don't block)
	go func() {
		if err := h.notifyService.SendBookingConfirmation(context.Background(), outing); err != nil {
			fmt.Printf("warning: failed to send booking confirmation: %v\n", err)
		}
	}()

	return &BookOutingResult{Outing: outing}, nil
}

// =============================================================================
// CHECK IN OUTING COMMAND
// =============================================================================

type CheckInOutingCommand struct {
	OutingID    string
	QRCode      string
	StaffUserID string
	Latitude    *float64
	Longitude   *float64
}

type CheckInOutingResult struct {
	Outing *domain.Outing
}

type CheckInOutingHandler struct {
	outingRepo    domain.OutingRepository
	notifyService domain.NotificationService
}

func NewCheckInOutingHandler(
	outingRepo domain.OutingRepository,
	notifyService domain.NotificationService,
) *CheckInOutingHandler {
	return &CheckInOutingHandler{
		outingRepo:    outingRepo,
		notifyService: notifyService,
	}
}

func (h *CheckInOutingHandler) Handle(ctx context.Context, cmd CheckInOutingCommand) (*CheckInOutingResult, error) {
	// 1. Get outing by ID or QR code
	var outing *domain.Outing
	var err error

	if cmd.OutingID != "" {
		outing, err = h.outingRepo.GetByID(ctx, cmd.OutingID)
	} else if cmd.QRCode != "" {
		outing, err = h.outingRepo.GetByQRCode(ctx, cmd.QRCode)
	} else {
		return nil, fmt.Errorf("outing ID or QR code is required")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get outing: %w", err)
	}

	// 2. Perform check-in
	if cmd.QRCode != "" {
		if err := outing.CheckInWithQR(cmd.QRCode, cmd.StaffUserID, cmd.Latitude, cmd.Longitude); err != nil {
			return nil, err
		}
	} else {
		if err := outing.CheckInManual(cmd.StaffUserID, cmd.Latitude, cmd.Longitude); err != nil {
			return nil, err
		}
	}

	// 3. Update outing
	if err := h.outingRepo.Update(ctx, outing); err != nil {
		return nil, fmt.Errorf("failed to update outing: %w", err)
	}

	// 4. Send notification (async)
	go func() {
		if err := h.notifyService.SendCheckInConfirmation(context.Background(), outing); err != nil {
			fmt.Printf("warning: failed to send check-in confirmation: %v\n", err)
		}
	}()

	return &CheckInOutingResult{Outing: outing}, nil
}

// =============================================================================
// CANCEL OUTING COMMAND
// =============================================================================

type CancelOutingCommand struct {
	OutingID    string
	CancelledBy domain.CancellationActor
	Reason      string
}

type CancelOutingResult struct {
	Outing *domain.Outing
}

type CancelOutingHandler struct {
	outingRepo    domain.OutingRepository
	offerService  domain.OfferService
	notifyService domain.NotificationService
}

func NewCancelOutingHandler(
	outingRepo domain.OutingRepository,
	offerService domain.OfferService,
	notifyService domain.NotificationService,
) *CancelOutingHandler {
	return &CancelOutingHandler{
		outingRepo:    outingRepo,
		offerService:  offerService,
		notifyService: notifyService,
	}
}

func (h *CancelOutingHandler) Handle(ctx context.Context, cmd CancelOutingCommand) (*CancelOutingResult, error) {
	// 1. Get outing
	outing, err := h.outingRepo.GetByID(ctx, cmd.OutingID)
	if err != nil {
		return nil, fmt.Errorf("failed to get outing: %w", err)
	}

	// 2. Cancel outing
	if err := outing.Cancel(cmd.CancelledBy, cmd.Reason); err != nil {
		return nil, err
	}

	// 3. Update outing
	if err := h.outingRepo.Update(ctx, outing); err != nil {
		return nil, fmt.Errorf("failed to update outing: %w", err)
	}

	// 4. Decrement offer booking count
	if err := h.offerService.DecrementBookingCount(ctx, outing.Offer().OfferID()); err != nil {
		fmt.Printf("warning: failed to decrement booking count: %v\n", err)
	}

	// 5. Send notification (async)
	go func() {
		if err := h.notifyService.SendCancellationNotification(context.Background(), outing); err != nil {
			fmt.Printf("warning: failed to send cancellation notification: %v\n", err)
		}
	}()

	return &CancelOutingResult{Outing: outing}, nil
}

// =============================================================================
// EXPIRE OUTINGS COMMAND (CRON JOB)
// =============================================================================

type ExpireOutingsCommand struct {
	BatchSize int
}

type ExpireOutingsResult struct {
	ExpiredCount int
}

type ExpireOutingsHandler struct {
	outingRepo domain.OutingRepository
}

func NewExpireOutingsHandler(outingRepo domain.OutingRepository) *ExpireOutingsHandler {
	return &ExpireOutingsHandler{
		outingRepo: outingRepo,
	}
}

func (h *ExpireOutingsHandler) Handle(ctx context.Context, cmd ExpireOutingsCommand) (*ExpireOutingsResult, error) {
	batchSize := cmd.BatchSize
	if batchSize <= 0 {
		batchSize = 100
	}

	// Get expired outings
	outings, err := h.outingRepo.GetExpiredOutings(ctx, time.Now(), batchSize)
	if err != nil {
		return nil, fmt.Errorf("failed to get expired outings: %w", err)
	}

	expiredCount := 0
	for _, outing := range outings {
		if err := outing.MarkAsExpired(); err != nil {
			fmt.Printf("warning: failed to mark outing %s as expired: %v\n", outing.ID(), err)
			continue
		}

		if err := h.outingRepo.Update(ctx, outing); err != nil {
			fmt.Printf("warning: failed to update expired outing %s: %v\n", outing.ID(), err)
			continue
		}

		expiredCount++
	}

	return &ExpireOutingsResult{ExpiredCount: expiredCount}, nil
}

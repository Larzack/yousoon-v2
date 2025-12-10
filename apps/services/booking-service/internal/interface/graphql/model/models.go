package model

import "time"

// =============================================================================
// ENUMS
// =============================================================================

type OutingStatus string

const (
	OutingStatusPending   OutingStatus = "PENDING"
	OutingStatusConfirmed OutingStatus = "CONFIRMED"
	OutingStatusCheckedIn OutingStatus = "CHECKED_IN"
	OutingStatusCancelled OutingStatus = "CANCELLED"
	OutingStatusExpired   OutingStatus = "EXPIRED"
	OutingStatusNoShow    OutingStatus = "NO_SHOW"
)

func (s OutingStatus) IsValid() bool {
	switch s {
	case OutingStatusPending, OutingStatusConfirmed, OutingStatusCheckedIn,
		OutingStatusCancelled, OutingStatusExpired, OutingStatusNoShow:
		return true
	}
	return false
}

func (s OutingStatus) String() string {
	return string(s)
}

type CheckInMethod string

const (
	CheckInMethodQRScan CheckInMethod = "QR_SCAN"
	CheckInMethodManual CheckInMethod = "MANUAL"
)

type CancellationActor string

const (
	CancellationActorUser    CancellationActor = "USER"
	CancellationActorPartner CancellationActor = "PARTNER"
	CancellationActorSystem  CancellationActor = "SYSTEM"
)

type BookingErrorCode string

const (
	BookingErrorCodeOfferNotFound        BookingErrorCode = "OFFER_NOT_FOUND"
	BookingErrorCodeOfferNotAvailable    BookingErrorCode = "OFFER_NOT_AVAILABLE"
	BookingErrorCodeOfferQuotaExceeded   BookingErrorCode = "OFFER_QUOTA_EXCEEDED"
	BookingErrorCodeUserNotVerified      BookingErrorCode = "USER_NOT_VERIFIED"
	BookingErrorCodeUserQuotaExceeded    BookingErrorCode = "USER_QUOTA_EXCEEDED"
	BookingErrorCodeAlreadyBooked        BookingErrorCode = "ALREADY_BOOKED"
	BookingErrorCodeSubscriptionRequired BookingErrorCode = "SUBSCRIPTION_REQUIRED"
	BookingErrorCodeInternalError        BookingErrorCode = "INTERNAL_ERROR"
)

type CheckInErrorCode string

const (
	CheckInErrorCodeOutingNotFound   CheckInErrorCode = "OUTING_NOT_FOUND"
	CheckInErrorCodeInvalidQRCode    CheckInErrorCode = "INVALID_QR_CODE"
	CheckInErrorCodeOutingExpired    CheckInErrorCode = "OUTING_EXPIRED"
	CheckInErrorCodeAlreadyCheckedIn CheckInErrorCode = "ALREADY_CHECKED_IN"
	CheckInErrorCodeOutingCancelled  CheckInErrorCode = "OUTING_CANCELLED"
	CheckInErrorCodeInternalError    CheckInErrorCode = "INTERNAL_ERROR"
)

type CancellationErrorCode string

const (
	CancellationErrorCodeOutingNotFound   CancellationErrorCode = "OUTING_NOT_FOUND"
	CancellationErrorCodeAlreadyCancelled CancellationErrorCode = "ALREADY_CANCELLED"
	CancellationErrorCodeAlreadyCheckedIn CancellationErrorCode = "ALREADY_CHECKED_IN"
	CancellationErrorCodeInternalError    CancellationErrorCode = "INTERNAL_ERROR"
)

// =============================================================================
// TYPES
// =============================================================================

type Outing struct {
	ID            string            `json:"id"`
	UserID        string            `json:"userId"`
	OfferSnapshot *OfferSnapshot    `json:"offerSnapshot"`
	QRCode        *QRCodeInfo       `json:"qrCode"`
	Status        OutingStatus      `json:"status"`
	Timeline      []*TimelineEntry  `json:"timeline"`
	CheckIn       *CheckInInfo      `json:"checkIn,omitempty"`
	Cancellation  *CancellationInfo `json:"cancellation,omitempty"`
	BookedAt      time.Time         `json:"bookedAt"`
	ExpiresAt     time.Time         `json:"expiresAt"`
	CreatedAt     time.Time         `json:"createdAt"`
	UpdatedAt     time.Time         `json:"updatedAt"`
}

type OfferSnapshot struct {
	OfferID              string    `json:"offerId"`
	PartnerID            string    `json:"partnerId"`
	EstablishmentID      string    `json:"establishmentId"`
	Title                string    `json:"title"`
	Description          *string   `json:"description,omitempty"`
	DiscountType         string    `json:"discountType"`
	DiscountValue        int       `json:"discountValue"`
	Category             *string   `json:"category,omitempty"`
	EstablishmentName    string    `json:"establishmentName"`
	EstablishmentAddress string    `json:"establishmentAddress"`
	Latitude             float64   `json:"latitude"`
	Longitude            float64   `json:"longitude"`
	ImageURL             *string   `json:"imageUrl,omitempty"`
	CapturedAt           time.Time `json:"capturedAt"`
}

type QRCodeInfo struct {
	Code      string    `json:"code"`
	FullCode  string    `json:"fullCode"`
	ExpiresAt time.Time `json:"expiresAt"`
	IsExpired bool      `json:"isExpired"`
}

type TimelineEntry struct {
	Status    OutingStatus           `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Actor     string                 `json:"actor"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

type CheckInInfo struct {
	CheckedInAt time.Time     `json:"checkedInAt"`
	CheckedInBy string        `json:"checkedInBy"`
	Method      CheckInMethod `json:"method"`
	Latitude    *float64      `json:"latitude,omitempty"`
	Longitude   *float64      `json:"longitude,omitempty"`
}

type CancellationInfo struct {
	CancelledAt time.Time         `json:"cancelledAt"`
	CancelledBy CancellationActor `json:"cancelledBy"`
	Reason      *string           `json:"reason,omitempty"`
}

type BookingStats struct {
	TotalBookings      int     `json:"totalBookings"`
	TotalCheckIns      int     `json:"totalCheckIns"`
	TotalCancelled     int     `json:"totalCancelled"`
	TotalExpired       int     `json:"totalExpired"`
	ConversionRate     float64 `json:"conversionRate"`
	AverageCheckInTime float64 `json:"averageCheckInTime"`
}

// =============================================================================
// CONNECTION TYPES
// =============================================================================

type OutingConnection struct {
	Edges      []*OutingEdge `json:"edges"`
	PageInfo   *PageInfo     `json:"pageInfo"`
	TotalCount int           `json:"totalCount"`
}

type OutingEdge struct {
	Node   *Outing `json:"node"`
	Cursor string  `json:"cursor"`
}

type PageInfo struct {
	HasNextPage     bool    `json:"hasNextPage"`
	HasPreviousPage bool    `json:"hasPreviousPage"`
	StartCursor     *string `json:"startCursor,omitempty"`
	EndCursor       *string `json:"endCursor,omitempty"`
}

// =============================================================================
// INPUTS
// =============================================================================

type BookOfferInput struct {
	OfferID string `json:"offerId"`
}

type CheckInInput struct {
	QRCode    string   `json:"qrCode"`
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
}

type ManualCheckInInput struct {
	OutingID  string   `json:"outingId"`
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
}

type CancelOutingInput struct {
	OutingID string  `json:"outingId"`
	Reason   *string `json:"reason,omitempty"`
}

type OutingFilterInput struct {
	Status    []OutingStatus `json:"status,omitempty"`
	StartDate *time.Time     `json:"startDate,omitempty"`
	EndDate   *time.Time     `json:"endDate,omitempty"`
}

type PaginationInput struct {
	First  *int    `json:"first,omitempty"`
	After  *string `json:"after,omitempty"`
	Last   *int    `json:"last,omitempty"`
	Before *string `json:"before,omitempty"`
}

// =============================================================================
// PAYLOADS
// =============================================================================

type BookOfferPayload struct {
	Success bool          `json:"success"`
	Outing  *Outing       `json:"outing,omitempty"`
	Error   *BookingError `json:"error,omitempty"`
}

type CheckInPayload struct {
	Success bool          `json:"success"`
	Outing  *Outing       `json:"outing,omitempty"`
	Error   *CheckInError `json:"error,omitempty"`
}

type CancelOutingPayload struct {
	Success bool               `json:"success"`
	Outing  *Outing            `json:"outing,omitempty"`
	Error   *CancellationError `json:"error,omitempty"`
}

type BookingError struct {
	Code    BookingErrorCode `json:"code"`
	Message string           `json:"message"`
}

type CheckInError struct {
	Code    CheckInErrorCode `json:"code"`
	Message string           `json:"message"`
}

type CancellationError struct {
	Code    CancellationErrorCode `json:"code"`
	Message string                `json:"message"`
}

// =============================================================================
// FEDERATION TYPES
// =============================================================================

type User struct {
	ID string `json:"id"`
}

type Offer struct {
	ID string `json:"id"`
}

type Partner struct {
	ID string `json:"id"`
}

type Establishment struct {
	ID string `json:"id"`
}

package domain

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/yousoon/shared/domain"
)

// =============================================================================
// ERRORS
// =============================================================================

var (
	ErrOutingNotFound       = errors.New("outing not found")
	ErrOutingAlreadyExists  = errors.New("outing already exists")
	ErrOutingExpired        = errors.New("outing has expired")
	ErrOutingAlreadyUsed    = errors.New("outing has already been used")
	ErrOutingCancelled      = errors.New("outing has been cancelled")
	ErrInvalidQRCode        = errors.New("invalid QR code")
	ErrInvalidOutingStatus  = errors.New("invalid outing status")
	ErrCannotCancelUsed     = errors.New("cannot cancel used outing")
	ErrOfferNotBookable     = errors.New("offer is not bookable")
	ErrUserQuotaExceeded    = errors.New("user booking quota exceeded")
	ErrOfferQuotaExceeded   = errors.New("offer booking quota exceeded")
	ErrInvalidCheckInWindow = errors.New("check-in window has not started or has expired")
)

// =============================================================================
// ENUMS
// =============================================================================

type OutingStatus string

const (
	OutingStatusPending   OutingStatus = "pending"
	OutingStatusConfirmed OutingStatus = "confirmed"
	OutingStatusCheckedIn OutingStatus = "checked_in"
	OutingStatusCancelled OutingStatus = "cancelled"
	OutingStatusExpired   OutingStatus = "expired"
	OutingStatusNoShow    OutingStatus = "no_show"
)

func (s OutingStatus) String() string {
	return string(s)
}

func (s OutingStatus) IsValid() bool {
	switch s {
	case OutingStatusPending, OutingStatusConfirmed, OutingStatusCheckedIn,
		OutingStatusCancelled, OutingStatusExpired, OutingStatusNoShow:
		return true
	}
	return false
}

type CancellationActor string

const (
	CancellationActorUser    CancellationActor = "user"
	CancellationActorPartner CancellationActor = "partner"
	CancellationActorSystem  CancellationActor = "system"
)

type CheckInMethod string

const (
	CheckInMethodQRScan CheckInMethod = "qr_scan"
	CheckInMethodManual CheckInMethod = "manual"
)

// =============================================================================
// VALUE OBJECTS
// =============================================================================

// QRCode represents the unique QR code for check-in
type QRCode struct {
	code      string
	signature string
	createdAt time.Time
	expiresAt time.Time
}

var qrSecretKey = []byte("yousoon-qr-secret-key-change-in-prod")

func NewQRCode(expiresAt time.Time) (QRCode, error) {
	// Generate random code
	codeBytes := make([]byte, 16)
	if _, err := rand.Read(codeBytes); err != nil {
		return QRCode{}, err
	}
	code := hex.EncodeToString(codeBytes)

	// Create HMAC signature
	h := hmac.New(sha256.New, qrSecretKey)
	h.Write([]byte(code))
	signature := base64.URLEncoding.EncodeToString(h.Sum(nil))

	return QRCode{
		code:      code,
		signature: signature,
		createdAt: time.Now(),
		expiresAt: expiresAt,
	}, nil
}

func ReconstructQRCode(code, signature string, createdAt, expiresAt time.Time) QRCode {
	return QRCode{
		code:      code,
		signature: signature,
		createdAt: createdAt,
		expiresAt: expiresAt,
	}
}

func (q QRCode) Code() string         { return q.code }
func (q QRCode) Signature() string    { return q.signature }
func (q QRCode) CreatedAt() time.Time { return q.createdAt }
func (q QRCode) ExpiresAt() time.Time { return q.expiresAt }

func (q QRCode) FullCode() string {
	return fmt.Sprintf("%s.%s", q.code, q.signature)
}

func (q QRCode) IsExpired() bool {
	return time.Now().After(q.expiresAt)
}

func (q QRCode) Matches(scanned string) bool {
	return q.code == scanned || q.FullCode() == scanned
}

func (q QRCode) Verify(scanned string) bool {
	if !q.Matches(scanned) {
		return false
	}

	// Verify signature
	h := hmac.New(sha256.New, qrSecretKey)
	h.Write([]byte(q.code))
	expectedSignature := base64.URLEncoding.EncodeToString(h.Sum(nil))

	return hmac.Equal([]byte(q.signature), []byte(expectedSignature))
}

// OfferSnapshot captures offer details at booking time (immutable)
type OfferSnapshot struct {
	offerID              string
	partnerID            string
	establishmentID      string
	title                string
	description          string
	discountType         string
	discountValue        int
	category             string
	establishmentName    string
	establishmentAddress string
	latitude             float64
	longitude            float64
	imageURL             string
	capturedAt           time.Time
}

func NewOfferSnapshot(
	offerID, partnerID, establishmentID string,
	title, description string,
	discountType string, discountValue int,
	category string,
	establishmentName, establishmentAddress string,
	lat, lng float64,
	imageURL string,
) OfferSnapshot {
	return OfferSnapshot{
		offerID:              offerID,
		partnerID:            partnerID,
		establishmentID:      establishmentID,
		title:                title,
		description:          description,
		discountType:         discountType,
		discountValue:        discountValue,
		category:             category,
		establishmentName:    establishmentName,
		establishmentAddress: establishmentAddress,
		latitude:             lat,
		longitude:            lng,
		imageURL:             imageURL,
		capturedAt:           time.Now(),
	}
}

func ReconstructOfferSnapshot(
	offerID, partnerID, establishmentID string,
	title, description string,
	discountType string, discountValue int,
	category string,
	establishmentName, establishmentAddress string,
	lat, lng float64,
	imageURL string,
	capturedAt time.Time,
) OfferSnapshot {
	return OfferSnapshot{
		offerID:              offerID,
		partnerID:            partnerID,
		establishmentID:      establishmentID,
		title:                title,
		description:          description,
		discountType:         discountType,
		discountValue:        discountValue,
		category:             category,
		establishmentName:    establishmentName,
		establishmentAddress: establishmentAddress,
		latitude:             lat,
		longitude:            lng,
		imageURL:             imageURL,
		capturedAt:           capturedAt,
	}
}

func (s OfferSnapshot) OfferID() string              { return s.offerID }
func (s OfferSnapshot) PartnerID() string            { return s.partnerID }
func (s OfferSnapshot) EstablishmentID() string      { return s.establishmentID }
func (s OfferSnapshot) Title() string                { return s.title }
func (s OfferSnapshot) Description() string          { return s.description }
func (s OfferSnapshot) DiscountType() string         { return s.discountType }
func (s OfferSnapshot) DiscountValue() int           { return s.discountValue }
func (s OfferSnapshot) Category() string             { return s.category }
func (s OfferSnapshot) EstablishmentName() string    { return s.establishmentName }
func (s OfferSnapshot) EstablishmentAddress() string { return s.establishmentAddress }
func (s OfferSnapshot) Latitude() float64            { return s.latitude }
func (s OfferSnapshot) Longitude() float64           { return s.longitude }
func (s OfferSnapshot) ImageURL() string             { return s.imageURL }
func (s OfferSnapshot) CapturedAt() time.Time        { return s.capturedAt }

// UserSnapshot captures user details at booking time
type UserSnapshot struct {
	userID    string
	firstName string
	lastName  string
	email     string
}

func NewUserSnapshot(userID, firstName, lastName, email string) UserSnapshot {
	return UserSnapshot{
		userID:    userID,
		firstName: firstName,
		lastName:  lastName,
		email:     email,
	}
}

func (s UserSnapshot) UserID() string    { return s.userID }
func (s UserSnapshot) FirstName() string { return s.firstName }
func (s UserSnapshot) LastName() string  { return s.lastName }
func (s UserSnapshot) Email() string     { return s.email }

// CheckInInfo contains check-in details
type CheckInInfo struct {
	checkedInAt time.Time
	checkedInBy string // UserID of staff member
	method      CheckInMethod
	latitude    *float64
	longitude   *float64
}

func NewCheckInInfo(checkedInBy string, method CheckInMethod, lat, lng *float64) CheckInInfo {
	return CheckInInfo{
		checkedInAt: time.Now(),
		checkedInBy: checkedInBy,
		method:      method,
		latitude:    lat,
		longitude:   lng,
	}
}

func ReconstructCheckInInfo(checkedInAt time.Time, checkedInBy string, method CheckInMethod, lat, lng *float64) CheckInInfo {
	return CheckInInfo{
		checkedInAt: checkedInAt,
		checkedInBy: checkedInBy,
		method:      method,
		latitude:    lat,
		longitude:   lng,
	}
}

func (c CheckInInfo) CheckedInAt() time.Time { return c.checkedInAt }
func (c CheckInInfo) CheckedInBy() string    { return c.checkedInBy }
func (c CheckInInfo) Method() CheckInMethod  { return c.method }
func (c CheckInInfo) Latitude() *float64     { return c.latitude }
func (c CheckInInfo) Longitude() *float64    { return c.longitude }

// CancellationInfo contains cancellation details
type CancellationInfo struct {
	cancelledAt time.Time
	cancelledBy CancellationActor
	reason      string
}

func NewCancellationInfo(cancelledBy CancellationActor, reason string) CancellationInfo {
	return CancellationInfo{
		cancelledAt: time.Now(),
		cancelledBy: cancelledBy,
		reason:      reason,
	}
}

func ReconstructCancellationInfo(cancelledAt time.Time, cancelledBy CancellationActor, reason string) CancellationInfo {
	return CancellationInfo{
		cancelledAt: cancelledAt,
		cancelledBy: cancelledBy,
		reason:      reason,
	}
}

func (c CancellationInfo) CancelledAt() time.Time         { return c.cancelledAt }
func (c CancellationInfo) CancelledBy() CancellationActor { return c.cancelledBy }
func (c CancellationInfo) Reason() string                 { return c.reason }

// TimelineEntry represents a status change in the outing lifecycle
type TimelineEntry struct {
	status    OutingStatus
	timestamp time.Time
	actor     string
	metadata  map[string]interface{}
}

func NewTimelineEntry(status OutingStatus, actor string, metadata map[string]interface{}) TimelineEntry {
	return TimelineEntry{
		status:    status,
		timestamp: time.Now(),
		actor:     actor,
		metadata:  metadata,
	}
}

func ReconstructTimelineEntry(status OutingStatus, timestamp time.Time, actor string, metadata map[string]interface{}) TimelineEntry {
	return TimelineEntry{
		status:    status,
		timestamp: timestamp,
		actor:     actor,
		metadata:  metadata,
	}
}

func (t TimelineEntry) Status() OutingStatus             { return t.status }
func (t TimelineEntry) Timestamp() time.Time             { return t.timestamp }
func (t TimelineEntry) Actor() string                    { return t.actor }
func (t TimelineEntry) Metadata() map[string]interface{} { return t.metadata }

// =============================================================================
// AGGREGATE ROOT: Outing
// =============================================================================

type Outing struct {
	domain.AggregateRoot

	id     string
	userID string

	// Snapshots (immutable at booking time)
	offer OfferSnapshot
	user  UserSnapshot

	// QR Code
	qrCode QRCode

	// Status & Timeline
	status   OutingStatus
	timeline []TimelineEntry

	// Check-in details (optional)
	checkIn *CheckInInfo

	// Cancellation details (optional)
	cancellation *CancellationInfo

	// Timing
	bookedAt  time.Time
	expiresAt time.Time

	// Metadata
	createdAt time.Time
	updatedAt time.Time
}

// Factory
func NewOuting(
	userID string,
	offer OfferSnapshot,
	user UserSnapshot,
	expirationMinutes int,
) (*Outing, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(expirationMinutes) * time.Minute)

	qrCode, err := NewQRCode(expiresAt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code: %w", err)
	}

	id := domain.NewID()

	outing := &Outing{
		id:     id,
		userID: userID,
		offer:  offer,
		user:   user,
		qrCode: qrCode,
		status: OutingStatusConfirmed,
		timeline: []TimelineEntry{
			NewTimelineEntry(OutingStatusConfirmed, "system", map[string]interface{}{
				"action": "booking_created",
			}),
		},
		bookedAt:  now,
		expiresAt: expiresAt,
		createdAt: now,
		updatedAt: now,
	}

	outing.AddDomainEvent(OutingBooked{
		OutingID:        id,
		UserID:          userID,
		OfferID:         offer.OfferID(),
		PartnerID:       offer.PartnerID(),
		EstablishmentID: offer.EstablishmentID(),
		QRCode:          qrCode.FullCode(),
		ExpiresAt:       expiresAt,
		Timestamp:       now,
	})

	return outing, nil
}

// Reconstruction
func ReconstructOuting(
	id, userID string,
	offer OfferSnapshot,
	user UserSnapshot,
	qrCode QRCode,
	status OutingStatus,
	timeline []TimelineEntry,
	checkIn *CheckInInfo,
	cancellation *CancellationInfo,
	bookedAt, expiresAt time.Time,
	createdAt, updatedAt time.Time,
) *Outing {
	return &Outing{
		id:           id,
		userID:       userID,
		offer:        offer,
		user:         user,
		qrCode:       qrCode,
		status:       status,
		timeline:     timeline,
		checkIn:      checkIn,
		cancellation: cancellation,
		bookedAt:     bookedAt,
		expiresAt:    expiresAt,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}

// Getters
func (o *Outing) ID() string                      { return o.id }
func (o *Outing) UserID() string                  { return o.userID }
func (o *Outing) Offer() OfferSnapshot            { return o.offer }
func (o *Outing) User() UserSnapshot              { return o.user }
func (o *Outing) QRCode() QRCode                  { return o.qrCode }
func (o *Outing) Status() OutingStatus            { return o.status }
func (o *Outing) Timeline() []TimelineEntry       { return o.timeline }
func (o *Outing) CheckIn() *CheckInInfo           { return o.checkIn }
func (o *Outing) Cancellation() *CancellationInfo { return o.cancellation }
func (o *Outing) BookedAt() time.Time             { return o.bookedAt }
func (o *Outing) ExpiresAt() time.Time            { return o.expiresAt }
func (o *Outing) CreatedAt() time.Time            { return o.createdAt }
func (o *Outing) UpdatedAt() time.Time            { return o.updatedAt }

// Business Logic
func (o *Outing) CanCheckIn() error {
	if o.status == OutingStatusCancelled {
		return ErrOutingCancelled
	}
	if o.status == OutingStatusCheckedIn {
		return ErrOutingAlreadyUsed
	}
	if o.status == OutingStatusExpired || o.status == OutingStatusNoShow {
		return ErrOutingExpired
	}
	if time.Now().After(o.expiresAt) {
		return ErrOutingExpired
	}
	return nil
}

func (o *Outing) CheckInWithQR(scannedQR string, staffUserID string, lat, lng *float64) error {
	if err := o.CanCheckIn(); err != nil {
		return err
	}

	if !o.qrCode.Verify(scannedQR) {
		return ErrInvalidQRCode
	}

	now := time.Now()
	o.status = OutingStatusCheckedIn
	checkIn := NewCheckInInfo(staffUserID, CheckInMethodQRScan, lat, lng)
	o.checkIn = &checkIn
	o.updatedAt = now

	o.timeline = append(o.timeline, NewTimelineEntry(OutingStatusCheckedIn, staffUserID, map[string]interface{}{
		"method": string(CheckInMethodQRScan),
	}))

	o.AddDomainEvent(OutingCheckedIn{
		OutingID:        o.id,
		UserID:          o.userID,
		OfferID:         o.offer.OfferID(),
		PartnerID:       o.offer.PartnerID(),
		EstablishmentID: o.offer.EstablishmentID(),
		CheckedInBy:     staffUserID,
		Method:          string(CheckInMethodQRScan),
		Timestamp:       now,
	})

	return nil
}

func (o *Outing) CheckInManual(staffUserID string, lat, lng *float64) error {
	if err := o.CanCheckIn(); err != nil {
		return err
	}

	now := time.Now()
	o.status = OutingStatusCheckedIn
	checkIn := NewCheckInInfo(staffUserID, CheckInMethodManual, lat, lng)
	o.checkIn = &checkIn
	o.updatedAt = now

	o.timeline = append(o.timeline, NewTimelineEntry(OutingStatusCheckedIn, staffUserID, map[string]interface{}{
		"method": string(CheckInMethodManual),
	}))

	o.AddDomainEvent(OutingCheckedIn{
		OutingID:        o.id,
		UserID:          o.userID,
		OfferID:         o.offer.OfferID(),
		PartnerID:       o.offer.PartnerID(),
		EstablishmentID: o.offer.EstablishmentID(),
		CheckedInBy:     staffUserID,
		Method:          string(CheckInMethodManual),
		Timestamp:       now,
	})

	return nil
}

func (o *Outing) Cancel(actor CancellationActor, reason string) error {
	if o.status == OutingStatusCheckedIn {
		return ErrCannotCancelUsed
	}
	if o.status == OutingStatusCancelled {
		return ErrOutingCancelled
	}

	now := time.Now()
	o.status = OutingStatusCancelled
	cancellation := NewCancellationInfo(actor, reason)
	o.cancellation = &cancellation
	o.updatedAt = now

	o.timeline = append(o.timeline, NewTimelineEntry(OutingStatusCancelled, string(actor), map[string]interface{}{
		"reason": reason,
	}))

	o.AddDomainEvent(OutingCancelled{
		OutingID:    o.id,
		UserID:      o.userID,
		OfferID:     o.offer.OfferID(),
		PartnerID:   o.offer.PartnerID(),
		CancelledBy: string(actor),
		Reason:      reason,
		Timestamp:   now,
	})

	return nil
}

func (o *Outing) MarkAsExpired() error {
	if o.status == OutingStatusCheckedIn {
		return ErrOutingAlreadyUsed
	}
	if o.status == OutingStatusCancelled {
		return ErrOutingCancelled
	}

	now := time.Now()
	o.status = OutingStatusExpired
	o.updatedAt = now

	o.timeline = append(o.timeline, NewTimelineEntry(OutingStatusExpired, "system", map[string]interface{}{
		"action": "auto_expired",
	}))

	o.AddDomainEvent(OutingExpired{
		OutingID:  o.id,
		UserID:    o.userID,
		OfferID:   o.offer.OfferID(),
		PartnerID: o.offer.PartnerID(),
		Timestamp: now,
	})

	return nil
}

func (o *Outing) MarkAsNoShow() error {
	if o.status == OutingStatusCheckedIn {
		return ErrOutingAlreadyUsed
	}
	if o.status == OutingStatusCancelled {
		return ErrOutingCancelled
	}

	now := time.Now()
	o.status = OutingStatusNoShow
	o.updatedAt = now

	o.timeline = append(o.timeline, NewTimelineEntry(OutingStatusNoShow, "system", map[string]interface{}{
		"action": "marked_no_show",
	}))

	return nil
}

func (o *Outing) IsExpired() bool {
	return time.Now().After(o.expiresAt) || o.status == OutingStatusExpired
}

func (o *Outing) IsActive() bool {
	return o.status == OutingStatusConfirmed || o.status == OutingStatusPending
}

package domain

import (
	"errors"
	"time"

	"github.com/yousoon/shared/domain"
)

// =============================================================================
// ERRORS
// =============================================================================

var (
	ErrNotificationNotFound = errors.New("notification not found")
	ErrInvalidChannel       = errors.New("invalid notification channel")
	ErrTemplateNotFound     = errors.New("notification template not found")
	ErrPushTokenNotFound    = errors.New("push token not found")
	ErrSendFailed           = errors.New("notification send failed")
)

// =============================================================================
// ENUMS
// =============================================================================

type NotificationChannel string

const (
	ChannelPush  NotificationChannel = "push"
	ChannelEmail NotificationChannel = "email"
	ChannelSMS   NotificationChannel = "sms"
)

func (c NotificationChannel) String() string {
	return string(c)
}

func (c NotificationChannel) IsValid() bool {
	switch c {
	case ChannelPush, ChannelEmail, ChannelSMS:
		return true
	}
	return false
}

type NotificationStatus string

const (
	StatusPending   NotificationStatus = "pending"
	StatusSent      NotificationStatus = "sent"
	StatusDelivered NotificationStatus = "delivered"
	StatusFailed    NotificationStatus = "failed"
	StatusRead      NotificationStatus = "read"
)

func (s NotificationStatus) String() string {
	return string(s)
}

type NotificationType string

const (
	TypeBookingConfirmed     NotificationType = "booking_confirmed"
	TypeBookingReminder      NotificationType = "booking_reminder"
	TypeBookingExpiring      NotificationType = "booking_expiring"
	TypeCheckInConfirmed     NotificationType = "checkin_confirmed"
	TypeOfferNearby          NotificationType = "offer_nearby"
	TypeNewOfferFromFavorite NotificationType = "new_offer_favorite"
	TypeReviewApproved       NotificationType = "review_approved"
	TypeReviewRejected       NotificationType = "review_rejected"
	TypePartnerVerified      NotificationType = "partner_verified"
	TypeIdentityVerified     NotificationType = "identity_verified"
	TypeIdentityRejected     NotificationType = "identity_rejected"
	TypeSubscriptionRenewed  NotificationType = "subscription_renewed"
	TypeSubscriptionExpiring NotificationType = "subscription_expiring"
	TypeWelcome              NotificationType = "welcome"
	TypeMarketing            NotificationType = "marketing"
)

func (t NotificationType) String() string {
	return string(t)
}

// =============================================================================
// AGGREGATE: Notification
// =============================================================================

type Notification struct {
	domain.AggregateRoot

	id      string
	userID  string
	channel NotificationChannel
	ntype   NotificationType

	// Content
	title    string
	body     string
	imageURL *string
	data     map[string]interface{}

	// Status
	status      NotificationStatus
	sentAt      *time.Time
	deliveredAt *time.Time
	readAt      *time.Time
	error       *string

	// Related entity
	relatedType string // "offer", "booking", "partner", etc.
	relatedID   *string

	createdAt time.Time
}

func NewNotification(
	userID string,
	channel NotificationChannel,
	ntype NotificationType,
	title, body string,
	imageURL *string,
	data map[string]interface{},
	relatedType string,
	relatedID *string,
) (*Notification, error) {
	if !channel.IsValid() {
		return nil, ErrInvalidChannel
	}

	now := time.Now()
	id := domain.NewID()

	return &Notification{
		id:          id,
		userID:      userID,
		channel:     channel,
		ntype:       ntype,
		title:       title,
		body:        body,
		imageURL:    imageURL,
		data:        data,
		status:      StatusPending,
		relatedType: relatedType,
		relatedID:   relatedID,
		createdAt:   now,
	}, nil
}

func ReconstructNotification(
	id, userID string,
	channel NotificationChannel,
	ntype NotificationType,
	title, body string,
	imageURL *string,
	data map[string]interface{},
	status NotificationStatus,
	sentAt, deliveredAt, readAt *time.Time,
	errMsg *string,
	relatedType string,
	relatedID *string,
	createdAt time.Time,
) *Notification {
	return &Notification{
		id:          id,
		userID:      userID,
		channel:     channel,
		ntype:       ntype,
		title:       title,
		body:        body,
		imageURL:    imageURL,
		data:        data,
		status:      status,
		sentAt:      sentAt,
		deliveredAt: deliveredAt,
		readAt:      readAt,
		error:       errMsg,
		relatedType: relatedType,
		relatedID:   relatedID,
		createdAt:   createdAt,
	}
}

// Getters
func (n *Notification) ID() string                   { return n.id }
func (n *Notification) UserID() string               { return n.userID }
func (n *Notification) Channel() NotificationChannel { return n.channel }
func (n *Notification) Type() NotificationType       { return n.ntype }
func (n *Notification) Title() string                { return n.title }
func (n *Notification) Body() string                 { return n.body }
func (n *Notification) ImageURL() *string            { return n.imageURL }
func (n *Notification) Data() map[string]interface{} { return n.data }
func (n *Notification) Status() NotificationStatus   { return n.status }
func (n *Notification) SentAt() *time.Time           { return n.sentAt }
func (n *Notification) DeliveredAt() *time.Time      { return n.deliveredAt }
func (n *Notification) ReadAt() *time.Time           { return n.readAt }
func (n *Notification) Error() *string               { return n.error }
func (n *Notification) RelatedType() string          { return n.relatedType }
func (n *Notification) RelatedID() *string           { return n.relatedID }
func (n *Notification) CreatedAt() time.Time         { return n.createdAt }

// Commands
func (n *Notification) MarkAsSent() {
	now := time.Now()
	n.status = StatusSent
	n.sentAt = &now
}

func (n *Notification) MarkAsDelivered() {
	now := time.Now()
	n.status = StatusDelivered
	n.deliveredAt = &now
}

func (n *Notification) MarkAsRead() {
	now := time.Now()
	n.status = StatusRead
	n.readAt = &now
}

func (n *Notification) MarkAsFailed(errMsg string) {
	n.status = StatusFailed
	n.error = &errMsg
}

// =============================================================================
// ENTITY: PushToken
// =============================================================================

type PushToken struct {
	id        string
	userID    string
	token     string
	platform  string // "ios", "android"
	createdAt time.Time
	updatedAt time.Time
}

func NewPushToken(userID, token, platform string) *PushToken {
	now := time.Now()
	return &PushToken{
		id:        domain.NewID(),
		userID:    userID,
		token:     token,
		platform:  platform,
		createdAt: now,
		updatedAt: now,
	}
}

func ReconstructPushToken(id, userID, token, platform string, createdAt, updatedAt time.Time) *PushToken {
	return &PushToken{
		id:        id,
		userID:    userID,
		token:     token,
		platform:  platform,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (p *PushToken) ID() string           { return p.id }
func (p *PushToken) UserID() string       { return p.userID }
func (p *PushToken) Token() string        { return p.token }
func (p *PushToken) Platform() string     { return p.platform }
func (p *PushToken) CreatedAt() time.Time { return p.createdAt }
func (p *PushToken) UpdatedAt() time.Time { return p.updatedAt }

func (p *PushToken) UpdateToken(token string) {
	p.token = token
	p.updatedAt = time.Now()
}

// =============================================================================
// ENTITY: NotificationTemplate
// =============================================================================

type NotificationTemplate struct {
	id        string
	code      string // e.g., "booking_confirmed"
	channel   NotificationChannel
	title     map[string]string // {fr: "...", en: "..."}
	body      map[string]string
	imageURL  *string
	isActive  bool
	createdAt time.Time
	updatedAt time.Time
}

func NewNotificationTemplate(
	code string,
	channel NotificationChannel,
	title, body map[string]string,
	imageURL *string,
) (*NotificationTemplate, error) {
	now := time.Now()
	return &NotificationTemplate{
		id:        domain.NewID(),
		code:      code,
		channel:   channel,
		title:     title,
		body:      body,
		imageURL:  imageURL,
		isActive:  true,
		createdAt: now,
		updatedAt: now,
	}, nil
}

func (t *NotificationTemplate) ID() string                   { return t.id }
func (t *NotificationTemplate) Code() string                 { return t.code }
func (t *NotificationTemplate) Channel() NotificationChannel { return t.channel }
func (t *NotificationTemplate) Title() map[string]string     { return t.title }
func (t *NotificationTemplate) Body() map[string]string      { return t.body }
func (t *NotificationTemplate) ImageURL() *string            { return t.imageURL }
func (t *NotificationTemplate) IsActive() bool               { return t.isActive }
func (t *NotificationTemplate) CreatedAt() time.Time         { return t.createdAt }
func (t *NotificationTemplate) UpdatedAt() time.Time         { return t.updatedAt }

func (t *NotificationTemplate) GetTitle(lang string) string {
	if title, ok := t.title[lang]; ok {
		return title
	}
	if title, ok := t.title["fr"]; ok {
		return title
	}
	return ""
}

func (t *NotificationTemplate) GetBody(lang string) string {
	if body, ok := t.body[lang]; ok {
		return body
	}
	if body, ok := t.body["fr"]; ok {
		return body
	}
	return ""
}

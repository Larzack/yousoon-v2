package domain

import "context"

// =============================================================================
// NOTIFICATION REPOSITORY
// =============================================================================

type NotificationRepository interface {
	Create(ctx context.Context, notification *Notification) error
	Update(ctx context.Context, notification *Notification) error
	GetByID(ctx context.Context, id string) (*Notification, error)
	GetByUserID(ctx context.Context, userID string, filter NotificationFilter) ([]*Notification, int64, error)
	GetUnreadCount(ctx context.Context, userID string) (int64, error)
	MarkAllAsRead(ctx context.Context, userID string) error
	DeleteOld(ctx context.Context, olderThan int) (int64, error)
}

type NotificationFilter struct {
	Status    []NotificationStatus
	Channel   *NotificationChannel
	Type      []NotificationType
	Unread    bool
	Offset    int
	Limit     int
	SortOrder string
}

func DefaultNotificationFilter() NotificationFilter {
	return NotificationFilter{
		Offset:    0,
		Limit:     20,
		SortOrder: "desc",
	}
}

// =============================================================================
// PUSH TOKEN REPOSITORY
// =============================================================================

type PushTokenRepository interface {
	Create(ctx context.Context, token *PushToken) error
	Update(ctx context.Context, token *PushToken) error
	Delete(ctx context.Context, id string) error
	GetByUserID(ctx context.Context, userID string) ([]*PushToken, error)
	GetByToken(ctx context.Context, token string) (*PushToken, error)
	DeleteByToken(ctx context.Context, token string) error
}

// =============================================================================
// DEVICE TOKEN REPOSITORY
// =============================================================================

type DeviceTokenRepository interface {
	Create(ctx context.Context, token *DeviceToken) error
	UpdateUserID(ctx context.Context, token, userID string) error
	Delete(ctx context.Context, token string) error
	Deactivate(ctx context.Context, token string) error
	FindByUserID(ctx context.Context, userID string) ([]*DeviceToken, error)
	FindByToken(ctx context.Context, token string) (*DeviceToken, error)
}

// =============================================================================
// TEMPLATE REPOSITORY
// =============================================================================

type TemplateRepository interface {
	Create(ctx context.Context, template *NotificationTemplate) error
	Update(ctx context.Context, template *NotificationTemplate) error
	GetByCode(ctx context.Context, code string, channel NotificationChannel) (*NotificationTemplate, error)
	GetAll(ctx context.Context) ([]*NotificationTemplate, error)
}

// =============================================================================
// NOTIFICATION SENDERS
// =============================================================================

type PushSender interface {
	Send(ctx context.Context, tokens []string, title, body string, data map[string]interface{}) error
	SendToUser(ctx context.Context, userID, title, body string, data map[string]interface{}) error
}

type EmailSender interface {
	Send(ctx context.Context, to, subject, body string, isHTML bool) error
	SendTemplate(ctx context.Context, to, templateID string, data map[string]interface{}) error
}

type SMSSender interface {
	Send(ctx context.Context, phoneNumber, message string) error
}

// =============================================================================
// USER PREFERENCES SERVICE
// =============================================================================

type UserPreferencesService interface {
	GetUserPreferences(ctx context.Context, userID string) (*UserNotificationPreferences, error)
	GetUserLanguage(ctx context.Context, userID string) (string, error)
	GetUserEmail(ctx context.Context, userID string) (string, error)
	GetUserPhone(ctx context.Context, userID string) (string, error)
}

type UserNotificationPreferences struct {
	PushEnabled      bool
	EmailEnabled     bool
	SMSEnabled       bool
	MarketingEnabled bool
}

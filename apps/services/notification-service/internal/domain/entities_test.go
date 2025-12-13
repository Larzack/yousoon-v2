package domain

import (
	"testing"
)

// =============================================================================
// NotificationChannel Tests
// =============================================================================

func TestNotificationChannel_IsValid(t *testing.T) {
	tests := []struct {
		channel NotificationChannel
		valid   bool
	}{
		{ChannelPush, true},
		{ChannelEmail, true},
		{ChannelSMS, true},
		{NotificationChannel("invalid"), false},
		{NotificationChannel(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.channel), func(t *testing.T) {
			if tt.channel.IsValid() != tt.valid {
				t.Errorf("NotificationChannel(%v).IsValid() = %v, want %v", tt.channel, tt.channel.IsValid(), tt.valid)
			}
		})
	}
}

func TestNotificationChannel_String(t *testing.T) {
	tests := []struct {
		channel  NotificationChannel
		expected string
	}{
		{ChannelPush, "push"},
		{ChannelEmail, "email"},
		{ChannelSMS, "sms"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if tt.channel.String() != tt.expected {
				t.Errorf("NotificationChannel.String() = %v, want %v", tt.channel.String(), tt.expected)
			}
		})
	}
}

// =============================================================================
// NotificationStatus Tests
// =============================================================================

func TestNotificationStatus_String(t *testing.T) {
	tests := []struct {
		status   NotificationStatus
		expected string
	}{
		{StatusPending, "pending"},
		{StatusSent, "sent"},
		{StatusDelivered, "delivered"},
		{StatusFailed, "failed"},
		{StatusRead, "read"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if tt.status.String() != tt.expected {
				t.Errorf("NotificationStatus.String() = %v, want %v", tt.status.String(), tt.expected)
			}
		})
	}
}

// =============================================================================
// NotificationType Tests
// =============================================================================

func TestNotificationType_String(t *testing.T) {
	tests := []struct {
		notifType NotificationType
		expected  string
	}{
		{TypeBookingConfirmed, "booking_confirmed"},
		{TypeBookingReminder, "booking_reminder"},
		{TypeCheckInConfirmed, "checkin_confirmed"},
		{TypeOfferNearby, "offer_nearby"},
		{TypeWelcome, "welcome"},
		{TypeMarketing, "marketing"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if tt.notifType.String() != tt.expected {
				t.Errorf("NotificationType.String() = %v, want %v", tt.notifType.String(), tt.expected)
			}
		})
	}
}

// =============================================================================
// Platform Tests
// =============================================================================

func TestPlatform_IsValid(t *testing.T) {
	tests := []struct {
		platform Platform
		valid    bool
	}{
		{PlatformIOS, true},
		{PlatformAndroid, true},
		{Platform("invalid"), false},
		{Platform(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.platform), func(t *testing.T) {
			if tt.platform.IsValid() != tt.valid {
				t.Errorf("Platform(%v).IsValid() = %v, want %v", tt.platform, tt.platform.IsValid(), tt.valid)
			}
		})
	}
}

func TestPlatform_String(t *testing.T) {
	tests := []struct {
		platform Platform
		expected string
	}{
		{PlatformIOS, "ios"},
		{PlatformAndroid, "android"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if tt.platform.String() != tt.expected {
				t.Errorf("Platform.String() = %v, want %v", tt.platform.String(), tt.expected)
			}
		})
	}
}

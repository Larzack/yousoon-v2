// Package domain tests for value objects.
package domain

import (
	"testing"
)

// =============================================================================
// Email Tests
// =============================================================================

func TestNewEmail_Valid(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "simple email",
			input: "test@example.com",
			want:  "test@example.com",
		},
		{
			name:  "email with subdomain",
			input: "test@mail.example.com",
			want:  "test@mail.example.com",
		},
		{
			name:  "uppercase email should be lowercased",
			input: "TEST@EXAMPLE.COM",
			want:  "test@example.com",
		},
		{
			name:  "email with plus sign",
			input: "test+tag@example.com",
			want:  "test+tag@example.com",
		},
		{
			name:  "email with dots",
			input: "test.user@example.com",
			want:  "test.user@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email, err := NewEmail(tt.input)
			if err != nil {
				t.Fatalf("NewEmail(%q) error = %v, want nil", tt.input, err)
			}
			if email.String() != tt.want {
				t.Errorf("NewEmail(%q) = %v, want %v", tt.input, email.String(), tt.want)
			}
		})
	}
}

func TestNewEmail_Invalid(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "empty string",
			input: "",
		},
		{
			name:  "no at sign",
			input: "testexample.com",
		},
		{
			name:  "no domain",
			input: "test@",
		},
		{
			name:  "no username",
			input: "@example.com",
		},
		{
			name:  "no TLD",
			input: "test@example",
		},
		{
			name:  "double at sign",
			input: "test@@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewEmail(tt.input)
			if err == nil {
				t.Errorf("NewEmail(%q) error = nil, want error", tt.input)
			}
		})
	}
}

func TestEmail_Equals(t *testing.T) {
	email1, _ := NewEmail("test@example.com")
	email2, _ := NewEmail("test@example.com")
	email3, _ := NewEmail("other@example.com")

	if !email1.Equals(email2) {
		t.Error("Equal emails should return true")
	}
	if email1.Equals(email3) {
		t.Error("Different emails should return false")
	}
}

func TestEmail_Domain(t *testing.T) {
	email, _ := NewEmail("test@example.com")

	if domain := email.Domain(); domain != "example.com" {
		t.Errorf("Email.Domain() = %v, want example.com", domain)
	}
}

// =============================================================================
// Phone Tests
// =============================================================================

func TestNewPhone_Valid(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "French mobile",
			input: "+33612345678",
		},
		{
			name:  "US number",
			input: "+14155552671",
		},
		{
			name:  "UK number",
			input: "+442071234567",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			phone, err := NewPhone(tt.input)
			if err != nil {
				t.Fatalf("NewPhone(%q) error = %v, want nil", tt.input, err)
			}
			if phone.String() != tt.input {
				t.Errorf("NewPhone(%q) = %v, want %v", tt.input, phone.String(), tt.input)
			}
		})
	}
}

func TestNewPhone_Invalid(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "no plus sign",
			input: "33612345678",
		},
		{
			name:  "only plus sign",
			input: "+",
		},
		{
			name:  "contains letters",
			input: "+33ABC123",
		},
		{
			name:  "local format",
			input: "0612345678",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewPhone(tt.input)
			if err == nil {
				t.Errorf("NewPhone(%q) error = nil, want error", tt.input)
			}
		})
	}
}

// =============================================================================
// GeoLocation Tests
// =============================================================================

func TestNewGeoLocation_Valid(t *testing.T) {
	tests := []struct {
		name      string
		longitude float64
		latitude  float64
	}{
		{
			name:      "Paris",
			longitude: 2.3522,
			latitude:  48.8566,
		},
		{
			name:      "New York",
			longitude: -74.0060,
			latitude:  40.7128,
		},
		{
			name:      "origin",
			longitude: 0,
			latitude:  0,
		},
		{
			name:      "max longitude",
			longitude: 180,
			latitude:  0,
		},
		{
			name:      "min longitude",
			longitude: -180,
			latitude:  0,
		},
		{
			name:      "max latitude",
			longitude: 0,
			latitude:  90,
		},
		{
			name:      "min latitude",
			longitude: 0,
			latitude:  -90,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loc, err := NewGeoLocation(tt.longitude, tt.latitude)
			if err != nil {
				t.Fatalf("NewGeoLocation(%v, %v) error = %v", tt.longitude, tt.latitude, err)
			}
			if loc.Longitude() != tt.longitude {
				t.Errorf("Longitude() = %v, want %v", loc.Longitude(), tt.longitude)
			}
			if loc.Latitude() != tt.latitude {
				t.Errorf("Latitude() = %v, want %v", loc.Latitude(), tt.latitude)
			}
			if loc.Type != "Point" {
				t.Errorf("Type = %v, want Point", loc.Type)
			}
		})
	}
}

func TestNewGeoLocation_Invalid(t *testing.T) {
	tests := []struct {
		name      string
		longitude float64
		latitude  float64
	}{
		{
			name:      "longitude too high",
			longitude: 181,
			latitude:  0,
		},
		{
			name:      "longitude too low",
			longitude: -181,
			latitude:  0,
		},
		{
			name:      "latitude too high",
			longitude: 0,
			latitude:  91,
		},
		{
			name:      "latitude too low",
			longitude: 0,
			latitude:  -91,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewGeoLocation(tt.longitude, tt.latitude)
			if err == nil {
				t.Errorf("NewGeoLocation(%v, %v) error = nil, want error", tt.longitude, tt.latitude)
			}
		})
	}
}

// =============================================================================
// UserID Tests
// =============================================================================

func TestNewUserID(t *testing.T) {
	id1 := NewUserID()
	id2 := NewUserID()

	if id1.IsZero() {
		t.Error("NewUserID() should not return zero ID")
	}
	if id1 == id2 {
		t.Error("NewUserID() should generate unique IDs")
	}
}

func TestParseUserID_Valid(t *testing.T) {
	original := NewUserID()
	parsed, err := ParseUserID(original.String())

	if err != nil {
		t.Fatalf("ParseUserID() error = %v", err)
	}
	if parsed != original {
		t.Errorf("ParseUserID() = %v, want %v", parsed, original)
	}
}

func TestParseUserID_Invalid(t *testing.T) {
	_, err := ParseUserID("not-a-valid-uuid")

	if err == nil {
		t.Error("ParseUserID() should return error for invalid UUID")
	}
}

// =============================================================================
// Profile Tests
// =============================================================================

func TestNewProfile(t *testing.T) {
	profile := NewProfile("John", "Doe")

	if profile.FirstName != "John" {
		t.Errorf("FirstName = %v, want John", profile.FirstName)
	}
	if profile.LastName != "Doe" {
		t.Errorf("LastName = %v, want Doe", profile.LastName)
	}
	if profile.DisplayName != "John Doe" {
		t.Errorf("DisplayName = %v, want John Doe", profile.DisplayName)
	}
}

func TestProfile_FullName(t *testing.T) {
	profile := Profile{FirstName: "Jane", LastName: "Smith"}

	if fullName := profile.FullName(); fullName != "Jane Smith" {
		t.Errorf("FullName() = %v, want Jane Smith", fullName)
	}
}

// =============================================================================
// UserGrade Tests
// =============================================================================

func TestUserGrade_String(t *testing.T) {
	tests := []struct {
		grade UserGrade
		want  string
	}{
		{GradeExplorateur, "explorateur"},
		{GradeAventurier, "aventurier"},
		{GradeGrandVoyageur, "grand_voyageur"},
		{GradeConquerant, "conquerant"},
		{UserGrade(99), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.grade.String(); got != tt.want {
				t.Errorf("UserGrade.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

// =============================================================================
// DefaultPreferences Tests
// =============================================================================

func TestDefaultPreferences(t *testing.T) {
	prefs := DefaultPreferences()

	if prefs.Language != "fr" {
		t.Errorf("Language = %v, want fr", prefs.Language)
	}
	if prefs.MaxDistance != 10 {
		t.Errorf("MaxDistance = %v, want 10", prefs.MaxDistance)
	}
	if !prefs.NotificationSettings.Push {
		t.Error("Push notifications should be enabled by default")
	}
	if !prefs.NotificationSettings.Email {
		t.Error("Email notifications should be enabled by default")
	}
	if prefs.NotificationSettings.SMS {
		t.Error("SMS notifications should be disabled by default")
	}
	if prefs.NotificationSettings.Marketing {
		t.Error("Marketing notifications should be disabled by default")
	}
}

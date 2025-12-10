// Package domain tests for the Identity bounded context.
package domain

import (
	"testing"
	"time"
)

// =============================================================================
// User Creation Tests
// =============================================================================

func TestNewUser(t *testing.T) {
	email, _ := NewEmail("test@example.com")
	profile := Profile{
		FirstName:   "John",
		LastName:    "Doe",
		DisplayName: "John Doe",
	}

	user, err := NewUser(email, "hashedPassword123", profile)

	if err != nil {
		t.Fatalf("NewUser() error = %v, want nil", err)
	}
	if user == nil {
		t.Fatal("NewUser() returned nil user")
	}
	if user.Email != email {
		t.Errorf("NewUser() email = %v, want %v", user.Email, email)
	}
	if user.Status != UserStatusActive {
		t.Errorf("NewUser() status = %v, want %v", user.Status, UserStatusActive)
	}
	if user.Grade != GradeExplorateur {
		t.Errorf("NewUser() grade = %v, want %v", user.Grade, GradeExplorateur)
	}
	if user.EmailVerified {
		t.Error("NewUser() emailVerified = true, want false")
	}
	if user.PhoneVerified {
		t.Error("NewUser() phoneVerified = true, want false")
	}
	if len(user.GetDomainEvents()) == 0 {
		t.Error("NewUser() should emit UserRegistered event")
	}
}

func TestNewUser_DefaultPreferences(t *testing.T) {
	email, _ := NewEmail("test@example.com")
	profile := Profile{FirstName: "John", LastName: "Doe"}

	user, _ := NewUser(email, "hash", profile)

	if user.Preferences.Language != "fr" {
		t.Errorf("NewUser() preferences.language = %v, want fr", user.Preferences.Language)
	}
	if user.Preferences.MaxDistance != 10 {
		t.Errorf("NewUser() preferences.maxDistance = %v, want 10", user.Preferences.MaxDistance)
	}
}

// =============================================================================
// User Business Logic Tests
// =============================================================================

func TestUser_CanBook(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*User)
		wantError bool
	}{
		{
			name: "active user with verified identity can book",
			setup: func(u *User) {
				u.Status = UserStatusActive
				u.Identity = &IdentityVerification{
					Status: VerificationStatusVerified,
				}
			},
			wantError: false,
		},
		{
			name: "suspended user cannot book",
			setup: func(u *User) {
				u.Status = UserStatusSuspended
				u.Identity = &IdentityVerification{
					Status: VerificationStatusVerified,
				}
			},
			wantError: true,
		},
		{
			name: "user without identity verification cannot book",
			setup: func(u *User) {
				u.Status = UserStatusActive
				u.Identity = nil
			},
			wantError: true,
		},
		{
			name: "user with pending identity verification cannot book",
			setup: func(u *User) {
				u.Status = UserStatusActive
				u.Identity = &IdentityVerification{
					Status: VerificationStatusPending,
				}
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := createTestUser()
			tt.setup(user)

			err := user.CanBook()
			if (err != nil) != tt.wantError {
				t.Errorf("User.CanBook() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestUser_VerifyEmail(t *testing.T) {
	user := createTestUser()
	initialVersion := user.Version

	user.VerifyEmail()

	if !user.EmailVerified {
		t.Error("VerifyEmail() emailVerified = false, want true")
	}
	if user.Version <= initialVersion {
		t.Error("VerifyEmail() should increment version")
	}
}

func TestUser_SetPhone(t *testing.T) {
	user := createTestUser()
	user.PhoneVerified = true
	phone, _ := NewPhone("+33612345678")

	user.SetPhone(phone)

	if user.Phone == nil || *user.Phone != phone {
		t.Errorf("SetPhone() phone = %v, want %v", user.Phone, phone)
	}
	if user.PhoneVerified {
		t.Error("SetPhone() should reset phoneVerified to false")
	}
}

func TestUser_UpdateProfile(t *testing.T) {
	user := createTestUser()
	newProfile := Profile{
		FirstName:   "Jane",
		LastName:    "Smith",
		DisplayName: "Jane Smith",
	}

	user.UpdateProfile(newProfile)

	if user.Profile.FirstName != "Jane" {
		t.Errorf("UpdateProfile() firstName = %v, want Jane", user.Profile.FirstName)
	}
	if user.Profile.LastName != "Smith" {
		t.Errorf("UpdateProfile() lastName = %v, want Smith", user.Profile.LastName)
	}
}

func TestUser_UpdateLocation(t *testing.T) {
	user := createTestUser()
	location, _ := NewGeoLocation(2.3522, 48.8566) // Paris

	user.UpdateLocation(location)

	if user.LastLocation == nil {
		t.Fatal("UpdateLocation() lastLocation is nil")
	}
	if user.LastLocation.Longitude() != 2.3522 {
		t.Errorf("UpdateLocation() longitude = %v, want 2.3522", user.LastLocation.Longitude())
	}
	if user.LastLocation.Latitude() != 48.8566 {
		t.Errorf("UpdateLocation() latitude = %v, want 48.8566", user.LastLocation.Latitude())
	}
}

func TestUser_RecordLogin(t *testing.T) {
	user := createTestUser()
	before := time.Now()

	user.RecordLogin()

	if user.LastLoginAt == nil {
		t.Fatal("RecordLogin() lastLoginAt is nil")
	}
	if user.LastLoginAt.Before(before) {
		t.Error("RecordLogin() lastLoginAt should be after test start")
	}
}

// =============================================================================
// Identity Verification Tests
// =============================================================================

func TestUser_SubmitIdentityVerification(t *testing.T) {
	user := createTestUser()
	verification := IdentityVerification{
		ID:           NewVerificationID(),
		Status:       VerificationStatusPending,
		DocumentType: DocumentTypeCNI,
	}

	err := user.SubmitIdentityVerification(verification)

	if err != nil {
		t.Fatalf("SubmitIdentityVerification() error = %v, want nil", err)
	}
	if user.Identity == nil {
		t.Fatal("SubmitIdentityVerification() identity is nil")
	}
	if user.Identity.Status != VerificationStatusPending {
		t.Errorf("SubmitIdentityVerification() status = %v, want pending", user.Identity.Status)
	}
}

func TestUser_SubmitIdentityVerification_AlreadyVerified(t *testing.T) {
	user := createTestUser()
	user.Identity = &IdentityVerification{
		Status: VerificationStatusVerified,
	}

	verification := IdentityVerification{
		ID:     NewVerificationID(),
		Status: VerificationStatusPending,
	}

	err := user.SubmitIdentityVerification(verification)

	if err == nil {
		t.Error("SubmitIdentityVerification() should return error when already verified")
	}
}

func TestUser_ApproveIdentityVerification(t *testing.T) {
	user := createTestUser()
	user.Identity = &IdentityVerification{
		ID:     NewVerificationID(),
		Status: VerificationStatusPending,
	}

	err := user.ApproveIdentityVerification()

	if err != nil {
		t.Fatalf("ApproveIdentityVerification() error = %v, want nil", err)
	}
	if user.Identity.Status != VerificationStatusVerified {
		t.Errorf("ApproveIdentityVerification() status = %v, want verified", user.Identity.Status)
	}
	if user.Identity.VerifiedAt == nil {
		t.Error("ApproveIdentityVerification() verifiedAt should be set")
	}
}

func TestUser_ApproveIdentityVerification_NoIdentity(t *testing.T) {
	user := createTestUser()
	user.Identity = nil

	err := user.ApproveIdentityVerification()

	if err == nil {
		t.Error("ApproveIdentityVerification() should return error when no identity")
	}
}

func TestUser_RejectIdentityVerification(t *testing.T) {
	user := createTestUser()
	user.Identity = &IdentityVerification{
		ID:     NewVerificationID(),
		Status: VerificationStatusPending,
	}
	reason := "Document expired"

	err := user.RejectIdentityVerification(reason)

	if err != nil {
		t.Fatalf("RejectIdentityVerification() error = %v, want nil", err)
	}
	if user.Identity.Status != VerificationStatusRejected {
		t.Errorf("RejectIdentityVerification() status = %v, want rejected", user.Identity.Status)
	}
	if user.Identity.RejectionReason == nil || *user.Identity.RejectionReason != reason {
		t.Errorf("RejectIdentityVerification() reason = %v, want %v", user.Identity.RejectionReason, reason)
	}
}

// =============================================================================
// Subscription Tests
// =============================================================================

func TestUser_SetSubscription(t *testing.T) {
	user := createTestUser()
	subID := NewSubscriptionID()
	eventsCount := len(user.GetDomainEvents())

	user.SetSubscription(subID)

	if user.SubscriptionID == nil || *user.SubscriptionID != subID {
		t.Errorf("SetSubscription() subscriptionId = %v, want %v", user.SubscriptionID, subID)
	}
	if len(user.GetDomainEvents()) <= eventsCount {
		t.Error("SetSubscription() should emit UserSubscribed event")
	}
}

func TestUser_ClearSubscription(t *testing.T) {
	user := createTestUser()
	subID := NewSubscriptionID()
	user.SubscriptionID = &subID

	user.ClearSubscription()

	if user.SubscriptionID != nil {
		t.Error("ClearSubscription() subscriptionId should be nil")
	}
}

func TestUser_ClearSubscription_NoSubscription(t *testing.T) {
	user := createTestUser()
	user.SubscriptionID = nil
	eventsCount := len(user.GetDomainEvents())

	user.ClearSubscription()

	// Should not panic or emit event
	if len(user.GetDomainEvents()) != eventsCount {
		t.Error("ClearSubscription() should not emit event when no subscription")
	}
}

// =============================================================================
// Status Tests
// =============================================================================

func TestUser_Suspend(t *testing.T) {
	user := createTestUser()

	user.Suspend("Terms violation")

	if user.Status != UserStatusSuspended {
		t.Errorf("Suspend() status = %v, want suspended", user.Status)
	}
}

// =============================================================================
// Test Helpers
// =============================================================================

func createTestUser() *User {
	email, _ := NewEmail("test@example.com")
	profile := Profile{
		FirstName:   "Test",
		LastName:    "User",
		DisplayName: "Test User",
	}
	user, _ := NewUser(email, "hashedPassword", profile)
	// Clear events for testing
	user.ClearDomainEvents()
	return user
}

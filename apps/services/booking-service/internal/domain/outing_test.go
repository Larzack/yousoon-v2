package domain

import (
	"testing"
	"time"
)

// =============================================================================
// Outing Creation Tests
// =============================================================================

func TestNewOuting(t *testing.T) {
	offer := createTestOfferSnapshot()
	user := createTestUserSnapshot()

	outing, err := NewOuting("user-123", offer, user, 30)

	if err != nil {
		t.Fatalf("NewOuting() error = %v, want nil", err)
	}
	if outing == nil {
		t.Fatal("NewOuting() returned nil outing")
	}
	if outing.Status() != OutingStatusConfirmed {
		t.Errorf("NewOuting() status = %v, want %v", outing.Status(), OutingStatusConfirmed)
	}
	if outing.UserID() != "user-123" {
		t.Errorf("NewOuting() userID = %v, want user-123", outing.UserID())
	}
	if outing.QRCode().Code() == "" {
		t.Error("NewOuting() should generate QR code")
	}
	if outing.ExpiresAt().Before(time.Now()) {
		t.Error("NewOuting() expiresAt should be in the future")
	}
}

func TestNewOuting_ExpirationTime(t *testing.T) {
	offer := createTestOfferSnapshot()
	user := createTestUserSnapshot()
	expirationMinutes := 30

	before := time.Now()
	outing, _ := NewOuting("user-123", offer, user, expirationMinutes)
	after := time.Now()

	expectedMin := before.Add(time.Duration(expirationMinutes) * time.Minute)
	expectedMax := after.Add(time.Duration(expirationMinutes) * time.Minute)

	if outing.ExpiresAt().Before(expectedMin) || outing.ExpiresAt().After(expectedMax) {
		t.Errorf("NewOuting() expiresAt = %v, want between %v and %v",
			outing.ExpiresAt(), expectedMin, expectedMax)
	}
}

func TestNewOuting_Timeline(t *testing.T) {
	offer := createTestOfferSnapshot()
	user := createTestUserSnapshot()

	outing, _ := NewOuting("user-123", offer, user, 30)

	if len(outing.Timeline()) != 1 {
		t.Errorf("NewOuting() timeline length = %d, want 1", len(outing.Timeline()))
	}
	if outing.Timeline()[0].Status() != OutingStatusConfirmed {
		t.Errorf("NewOuting() initial timeline status = %v, want %v",
			outing.Timeline()[0].Status(), OutingStatusConfirmed)
	}
}

// =============================================================================
// QRCode Tests
// =============================================================================

func TestNewQRCode(t *testing.T) {
	expiresAt := time.Now().Add(30 * time.Minute)

	qr, err := NewQRCode(expiresAt)

	if err != nil {
		t.Fatalf("NewQRCode() error = %v, want nil", err)
	}
	if qr.Code() == "" {
		t.Error("NewQRCode() code should not be empty")
	}
	if qr.Signature() == "" {
		t.Error("NewQRCode() signature should not be empty")
	}
	if qr.ExpiresAt() != expiresAt {
		t.Errorf("NewQRCode() expiresAt = %v, want %v", qr.ExpiresAt(), expiresAt)
	}
}

func TestQRCode_FullCode(t *testing.T) {
	qr, _ := NewQRCode(time.Now().Add(30 * time.Minute))

	fullCode := qr.FullCode()
	expected := qr.Code() + "." + qr.Signature()

	if fullCode != expected {
		t.Errorf("QRCode.FullCode() = %v, want %v", fullCode, expected)
	}
}

func TestQRCode_IsExpired(t *testing.T) {
	tests := []struct {
		name      string
		expiresAt time.Time
		want      bool
	}{
		{
			name:      "not expired",
			expiresAt: time.Now().Add(30 * time.Minute),
			want:      false,
		},
		{
			name:      "expired",
			expiresAt: time.Now().Add(-1 * time.Minute),
			want:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qr := ReconstructQRCode("code", "sig", time.Now(), tt.expiresAt)
			if got := qr.IsExpired(); got != tt.want {
				t.Errorf("QRCode.IsExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQRCode_Matches(t *testing.T) {
	qr, _ := NewQRCode(time.Now().Add(30 * time.Minute))

	tests := []struct {
		name    string
		scanned string
		want    bool
	}{
		{
			name:    "matches code only",
			scanned: qr.Code(),
			want:    true,
		},
		{
			name:    "matches full code",
			scanned: qr.FullCode(),
			want:    true,
		},
		{
			name:    "does not match wrong code",
			scanned: "wrong-code",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := qr.Matches(tt.scanned); got != tt.want {
				t.Errorf("QRCode.Matches(%q) = %v, want %v", tt.scanned, got, tt.want)
			}
		})
	}
}

func TestQRCode_Verify(t *testing.T) {
	qr, _ := NewQRCode(time.Now().Add(30 * time.Minute))

	// Valid code should verify
	if !qr.Verify(qr.Code()) {
		t.Error("QRCode.Verify() should return true for valid code")
	}

	// Invalid code should not verify
	if qr.Verify("wrong-code") {
		t.Error("QRCode.Verify() should return false for invalid code")
	}
}

// =============================================================================
// Check-In Tests
// =============================================================================

func TestOuting_CanCheckIn(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*Outing)
		wantError bool
	}{
		{
			name:      "confirmed outing can check in",
			setup:     func(o *Outing) {},
			wantError: false,
		},
		{
			name: "cancelled outing cannot check in",
			setup: func(o *Outing) {
				o.status = OutingStatusCancelled
			},
			wantError: true,
		},
		{
			name: "already checked in outing cannot check in again",
			setup: func(o *Outing) {
				o.status = OutingStatusCheckedIn
			},
			wantError: true,
		},
		{
			name: "expired outing cannot check in",
			setup: func(o *Outing) {
				o.status = OutingStatusExpired
			},
			wantError: true,
		},
		{
			name: "no-show outing cannot check in",
			setup: func(o *Outing) {
				o.status = OutingStatusNoShow
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outing := createTestOuting()
			tt.setup(outing)

			err := outing.CanCheckIn()
			if (err != nil) != tt.wantError {
				t.Errorf("Outing.CanCheckIn() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestOuting_CheckInWithQR(t *testing.T) {
	outing := createTestOuting()
	qrCode := outing.QRCode().Code()
	staffID := "staff-123"

	err := outing.CheckInWithQR(qrCode, staffID, nil, nil)

	if err != nil {
		t.Fatalf("CheckInWithQR() error = %v, want nil", err)
	}
	if outing.Status() != OutingStatusCheckedIn {
		t.Errorf("CheckInWithQR() status = %v, want %v", outing.Status(), OutingStatusCheckedIn)
	}
	if outing.CheckIn() == nil {
		t.Error("CheckInWithQR() checkIn should not be nil")
	}
	if outing.CheckIn().Method() != CheckInMethodQRScan {
		t.Errorf("CheckInWithQR() method = %v, want %v", outing.CheckIn().Method(), CheckInMethodQRScan)
	}
	if outing.CheckIn().CheckedInBy() != staffID {
		t.Errorf("CheckInWithQR() checkedInBy = %v, want %v", outing.CheckIn().CheckedInBy(), staffID)
	}
}

func TestOuting_CheckInWithQR_InvalidQR(t *testing.T) {
	outing := createTestOuting()

	err := outing.CheckInWithQR("invalid-qr-code", "staff-123", nil, nil)

	if err != ErrInvalidQRCode {
		t.Errorf("CheckInWithQR() error = %v, want %v", err, ErrInvalidQRCode)
	}
	if outing.Status() != OutingStatusConfirmed {
		t.Errorf("CheckInWithQR() status should remain %v", OutingStatusConfirmed)
	}
}

func TestOuting_CheckInWithQR_AlreadyCheckedIn(t *testing.T) {
	outing := createTestOuting()
	qrCode := outing.QRCode().Code()

	// First check-in should succeed
	_ = outing.CheckInWithQR(qrCode, "staff-123", nil, nil)

	// Second check-in should fail
	err := outing.CheckInWithQR(qrCode, "staff-456", nil, nil)

	if err != ErrOutingAlreadyUsed {
		t.Errorf("CheckInWithQR() second attempt error = %v, want %v", err, ErrOutingAlreadyUsed)
	}
}

func TestOuting_CheckInManual(t *testing.T) {
	outing := createTestOuting()
	staffID := "staff-123"

	err := outing.CheckInManual(staffID, nil, nil)

	if err != nil {
		t.Fatalf("CheckInManual() error = %v, want nil", err)
	}
	if outing.Status() != OutingStatusCheckedIn {
		t.Errorf("CheckInManual() status = %v, want %v", outing.Status(), OutingStatusCheckedIn)
	}
	if outing.CheckIn().Method() != CheckInMethodManual {
		t.Errorf("CheckInManual() method = %v, want %v", outing.CheckIn().Method(), CheckInMethodManual)
	}
}

// =============================================================================
// Cancellation Tests
// =============================================================================

func TestOuting_Cancel(t *testing.T) {
	outing := createTestOuting()
	reason := "User requested cancellation"

	err := outing.Cancel(CancellationActorUser, reason)

	if err != nil {
		t.Fatalf("Cancel() error = %v, want nil", err)
	}
	if outing.Status() != OutingStatusCancelled {
		t.Errorf("Cancel() status = %v, want %v", outing.Status(), OutingStatusCancelled)
	}
	if outing.Cancellation() == nil {
		t.Error("Cancel() cancellation should not be nil")
	}
	if outing.Cancellation().CancelledBy() != CancellationActorUser {
		t.Errorf("Cancel() cancelledBy = %v, want %v", outing.Cancellation().CancelledBy(), CancellationActorUser)
	}
	if outing.Cancellation().Reason() != reason {
		t.Errorf("Cancel() reason = %v, want %v", outing.Cancellation().Reason(), reason)
	}
}

func TestOuting_Cancel_AlreadyCheckedIn(t *testing.T) {
	outing := createTestOuting()
	outing.status = OutingStatusCheckedIn

	err := outing.Cancel(CancellationActorUser, "reason")

	if err != ErrCannotCancelUsed {
		t.Errorf("Cancel() error = %v, want %v", err, ErrCannotCancelUsed)
	}
}

// =============================================================================
// Status Tests
// =============================================================================

func TestOutingStatus_IsValid(t *testing.T) {
	tests := []struct {
		status OutingStatus
		want   bool
	}{
		{OutingStatusPending, true},
		{OutingStatusConfirmed, true},
		{OutingStatusCheckedIn, true},
		{OutingStatusCancelled, true},
		{OutingStatusExpired, true},
		{OutingStatusNoShow, true},
		{OutingStatus("invalid"), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			if got := tt.status.IsValid(); got != tt.want {
				t.Errorf("OutingStatus.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

// =============================================================================
// OfferSnapshot Tests
// =============================================================================

func TestOfferSnapshot_Getters(t *testing.T) {
	snapshot := NewOfferSnapshot(
		"offer-123", "partner-456", "est-789",
		"Great Offer", "Description",
		"percentage", 20,
		"restaurant",
		"Restaurant XYZ", "123 Main St",
		48.8566, 2.3522,
		"http://example.com/image.jpg",
	)

	if snapshot.OfferID() != "offer-123" {
		t.Errorf("OfferID() = %v, want offer-123", snapshot.OfferID())
	}
	if snapshot.PartnerID() != "partner-456" {
		t.Errorf("PartnerID() = %v, want partner-456", snapshot.PartnerID())
	}
	if snapshot.EstablishmentID() != "est-789" {
		t.Errorf("EstablishmentID() = %v, want est-789", snapshot.EstablishmentID())
	}
	if snapshot.Title() != "Great Offer" {
		t.Errorf("Title() = %v, want Great Offer", snapshot.Title())
	}
	if snapshot.DiscountType() != "percentage" {
		t.Errorf("DiscountType() = %v, want percentage", snapshot.DiscountType())
	}
	if snapshot.DiscountValue() != 20 {
		t.Errorf("DiscountValue() = %v, want 20", snapshot.DiscountValue())
	}
}

// =============================================================================
// Test Helpers
// =============================================================================

func createTestOfferSnapshot() OfferSnapshot {
	return NewOfferSnapshot(
		"offer-123", "partner-456", "est-789",
		"Test Offer", "Test Description",
		"percentage", 20,
		"restaurant",
		"Test Restaurant", "123 Test St",
		48.8566, 2.3522,
		"http://example.com/image.jpg",
	)
}

func createTestUserSnapshot() UserSnapshot {
	return NewUserSnapshot("user-123", "John", "Doe", "john@example.com")
}

func createTestOuting() *Outing {
	offer := createTestOfferSnapshot()
	user := createTestUserSnapshot()
	outing, _ := NewOuting("user-123", offer, user, 30)
	outing.ClearDomainEvents()
	return outing
}

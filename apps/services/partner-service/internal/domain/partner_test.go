package domain

import (
	"testing"
)

// =============================================================================
// Partner Creation Tests
// =============================================================================

func TestNewPartner(t *testing.T) {
	company := Company{
		Name:      "Test Restaurant",
		TradeName: "Test Restaurant",
		Siret:     "12345678901234",
		LegalForm: "SARL",
	}
	contact := Contact{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@test.com",
		Phone:     "+33612345678",
	}

	partner, err := NewPartner("user-123", company, contact, "restaurant")

	if err != nil {
		t.Fatalf("NewPartner() error = %v, want nil", err)
	}
	if partner == nil {
		t.Fatal("NewPartner() returned nil partner")
	}
	if partner.Status != PartnerStatusPending {
		t.Errorf("NewPartner() status = %v, want %v", partner.Status, PartnerStatusPending)
	}
	if partner.OwnerUserID != "user-123" {
		t.Errorf("NewPartner() ownerUserID = %v, want user-123", partner.OwnerUserID)
	}
	if partner.Company.Name != "Test Restaurant" {
		t.Errorf("NewPartner() company.name = %v, want Test Restaurant", partner.Company.Name)
	}
	if len(partner.Establishments) != 0 {
		t.Errorf("NewPartner() establishments should be empty, got %d", len(partner.Establishments))
	}
}

func TestPartner_CanPublishOffer_NotVerified(t *testing.T) {
	company := Company{Name: "Test"}
	contact := Contact{Email: "test@test.com"}
	partner, _ := NewPartner("user-123", company, contact, "restaurant")

	err := partner.CanPublishOffer()

	if err != ErrPartnerNotVerified {
		t.Errorf("CanPublishOffer() error = %v, want %v", err, ErrPartnerNotVerified)
	}
}

func TestPartner_CanPublishOffer_NoEstablishment(t *testing.T) {
	company := Company{Name: "Test"}
	contact := Contact{Email: "test@test.com"}
	partner, _ := NewPartner("user-123", company, contact, "restaurant")
	partner.Status = PartnerStatusActive

	err := partner.CanPublishOffer()

	if err != ErrNoEstablishment {
		t.Errorf("CanPublishOffer() error = %v, want %v", err, ErrNoEstablishment)
	}
}

func TestPartner_Verify(t *testing.T) {
	company := Company{Name: "Test"}
	contact := Contact{Email: "test@test.com"}
	partner, _ := NewPartner("user-123", company, contact, "restaurant")

	err := partner.Verify("admin-123")

	if err != nil {
		t.Fatalf("Verify() error = %v, want nil", err)
	}
	if partner.Status != PartnerStatusActive {
		t.Errorf("Verify() status = %v, want %v", partner.Status, PartnerStatusActive)
	}
	if partner.VerifiedAt == nil {
		t.Error("Verify() should set verifiedAt")
	}
}

func TestPartner_Verify_AlreadyVerified(t *testing.T) {
	company := Company{Name: "Test"}
	contact := Contact{Email: "test@test.com"}
	partner, _ := NewPartner("user-123", company, contact, "restaurant")
	partner.Status = PartnerStatusActive

	err := partner.Verify("admin-123")

	if err != ErrAlreadyVerified {
		t.Errorf("Verify() error = %v, want %v", err, ErrAlreadyVerified)
	}
}

// =============================================================================
// PartnerID Tests
// =============================================================================

func TestNewPartnerID(t *testing.T) {
	id1 := NewPartnerID()
	id2 := NewPartnerID()

	if id1 == "" {
		t.Error("NewPartnerID() should not return empty string")
	}
	if id1 == id2 {
		t.Error("NewPartnerID() should generate unique IDs")
	}
}

// =============================================================================
// PartnerStats Tests
// =============================================================================

func TestNewPartnerStats(t *testing.T) {
	stats := NewPartnerStats()

	if stats.TotalOffers != 0 {
		t.Errorf("NewPartnerStats() totalOffers = %d, want 0", stats.TotalOffers)
	}
	if stats.ActiveOffers != 0 {
		t.Errorf("NewPartnerStats() activeOffers = %d, want 0", stats.ActiveOffers)
	}
	if stats.TotalBookings != 0 {
		t.Errorf("NewPartnerStats() totalBookings = %d, want 0", stats.TotalBookings)
	}
}

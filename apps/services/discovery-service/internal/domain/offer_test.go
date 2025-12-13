package domain

import (
	"testing"
	"time"
)

// =============================================================================
// Offer Creation Tests
// =============================================================================

func TestNewOffer(t *testing.T) {
	discount := NewPercentageDiscount(20, nil, nil)
	validity := NewValidity(time.Now(), time.Now().Add(30*24*time.Hour), "Europe/Paris")

	offer, err := NewOffer(
		"partner-123",
		"establishment-123",
		"Happy Hour -20%",
		"Profitez de 20% de réduction sur toutes les boissons",
		"category-123",
		discount,
		validity,
	)

	if err != nil {
		t.Fatalf("NewOffer() error = %v, want nil", err)
	}
	if offer == nil {
		t.Fatal("NewOffer() returned nil offer")
	}
	if offer.Title() != "Happy Hour -20%" {
		t.Errorf("NewOffer() title = %v, want Happy Hour -20%%", offer.Title())
	}
	if offer.Status() != OfferStatusDraft {
		t.Errorf("NewOffer() status = %v, want %v", offer.Status(), OfferStatusDraft)
	}
}

func TestNewOffer_EmptyTitle(t *testing.T) {
	discount := NewPercentageDiscount(20, nil, nil)
	validity := NewValidity(time.Now(), time.Now().Add(30*24*time.Hour), "Europe/Paris")

	_, err := NewOffer(
		"partner-123",
		"establishment-123",
		"", // Empty title
		"Description",
		"category-123",
		discount,
		validity,
	)

	if err == nil {
		t.Error("NewOffer() should return error for empty title")
	}
}

func TestNewOffer_EmptyPartnerID(t *testing.T) {
	discount := NewPercentageDiscount(20, nil, nil)
	validity := NewValidity(time.Now(), time.Now().Add(30*24*time.Hour), "Europe/Paris")

	_, err := NewOffer(
		"", // Empty partner ID
		"establishment-123",
		"Title",
		"Description",
		"category-123",
		discount,
		validity,
	)

	if err == nil {
		t.Error("NewOffer() should return error for empty partner ID")
	}
}

// =============================================================================
// OfferStatus Tests
// =============================================================================

func TestOfferStatus_Values(t *testing.T) {
	tests := []struct {
		status   OfferStatus
		expected string
	}{
		{OfferStatusDraft, "draft"},
		{OfferStatusPending, "pending"},
		{OfferStatusActive, "active"},
		{OfferStatusPaused, "paused"},
		{OfferStatusExpired, "expired"},
		{OfferStatusArchived, "archived"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if string(tt.status) != tt.expected {
				t.Errorf("OfferStatus = %v, want %v", tt.status, tt.expected)
			}
		})
	}
}

// =============================================================================
// Discount Tests
// =============================================================================

func TestNewPercentageDiscount(t *testing.T) {
	discount := NewPercentageDiscount(25, nil, nil)

	if discount.Type() != DiscountTypePercentage {
		t.Errorf("NewPercentageDiscount() type = %v, want %v", discount.Type(), DiscountTypePercentage)
	}
	if discount.Value() != 25 {
		t.Errorf("NewPercentageDiscount() value = %v, want 25", discount.Value())
	}
}

func TestNewFixedDiscount(t *testing.T) {
	originalPrice := 1000                             // 10.00€ in cents
	discount := NewFixedDiscount(200, &originalPrice) // 2.00€ discount

	if discount.Type() != DiscountTypeFixed {
		t.Errorf("NewFixedDiscount() type = %v, want %v", discount.Type(), DiscountTypeFixed)
	}
	if discount.Value() != 200 {
		t.Errorf("NewFixedDiscount() value = %v, want 200", discount.Value())
	}
}

// =============================================================================
// Validity Tests
// =============================================================================

func TestValidity_IsExpired(t *testing.T) {
	pastValidity := NewValidity(
		time.Now().Add(-48*time.Hour),
		time.Now().Add(-24*time.Hour),
		"Europe/Paris",
	)

	if !pastValidity.IsExpired() {
		t.Error("Validity.IsExpired() should return true for past dates")
	}
}

func TestValidity_IsActiveNow(t *testing.T) {
	futureValidity := NewValidity(
		time.Now().Add(24*time.Hour),
		time.Now().Add(48*time.Hour),
		"Europe/Paris",
	)

	if futureValidity.IsActiveNow() {
		t.Error("Validity.IsActiveNow() should return false for future dates")
	}

	currentValidity := NewValidity(
		time.Now().Add(-24*time.Hour),
		time.Now().Add(24*time.Hour),
		"Europe/Paris",
	)

	if !currentValidity.IsActiveNow() {
		t.Error("Validity.IsActiveNow() should return true for current period")
	}
}

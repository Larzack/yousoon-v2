package domain

import (
	"testing"
	"time"
)

// =============================================================================
// Favorite Tests
// =============================================================================

func TestNewFavorite(t *testing.T) {
	favorite, err := NewFavorite(
		"user-123",
		"offer-456",
		"Happy Hour -20%",
		"https://example.com/image.jpg",
		"Le Bar Sympa",
	)

	if err != nil {
		t.Fatalf("NewFavorite() error = %v, want nil", err)
	}
	if favorite == nil {
		t.Fatal("NewFavorite() returned nil favorite")
	}
	if favorite.UserID() != "user-123" {
		t.Errorf("NewFavorite() userID = %v, want user-123", favorite.UserID())
	}
	if favorite.OfferID() != "offer-456" {
		t.Errorf("NewFavorite() offerID = %v, want offer-456", favorite.OfferID())
	}
	if favorite.OfferTitle() != "Happy Hour -20%" {
		t.Errorf("NewFavorite() offerTitle = %v, want Happy Hour -20%%", favorite.OfferTitle())
	}
	if favorite.ID() == "" {
		t.Error("NewFavorite() should generate an ID")
	}
}

func TestReconstructFavorite(t *testing.T) {
	createdAt := time.Now().Add(-24 * time.Hour)

	favorite := ReconstructFavorite(
		"fav-123",
		"user-123",
		"offer-456",
		"Title",
		"https://example.com/image.jpg",
		"Partner",
		createdAt,
	)

	if favorite == nil {
		t.Fatal("ReconstructFavorite() returned nil")
	}
	if favorite.ID() != "fav-123" {
		t.Errorf("ReconstructFavorite() id = %v, want fav-123", favorite.ID())
	}
	if !favorite.CreatedAt().Equal(createdAt) {
		t.Errorf("ReconstructFavorite() createdAt = %v, want %v", favorite.CreatedAt(), createdAt)
	}
}

// =============================================================================
// ReviewStatus Tests
// =============================================================================

func TestReviewStatus_IsValid(t *testing.T) {
	tests := []struct {
		status ReviewStatus
		valid  bool
	}{
		{ReviewStatusPending, true},
		{ReviewStatusApproved, true},
		{ReviewStatusRejected, true},
		{ReviewStatusReported, true},
		{ReviewStatus("invalid"), false},
		{ReviewStatus(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			if tt.status.IsValid() != tt.valid {
				t.Errorf("ReviewStatus(%v).IsValid() = %v, want %v", tt.status, tt.status.IsValid(), tt.valid)
			}
		})
	}
}

func TestReviewStatus_String(t *testing.T) {
	tests := []struct {
		status   ReviewStatus
		expected string
	}{
		{ReviewStatusPending, "pending"},
		{ReviewStatusApproved, "approved"},
		{ReviewStatusRejected, "rejected"},
		{ReviewStatusReported, "reported"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if tt.status.String() != tt.expected {
				t.Errorf("ReviewStatus.String() = %v, want %v", tt.status.String(), tt.expected)
			}
		})
	}
}

// =============================================================================
// Error Tests
// =============================================================================

func TestDomainErrors(t *testing.T) {
	tests := []struct {
		err      error
		expected string
	}{
		{ErrFavoriteNotFound, "favorite not found"},
		{ErrFavoriteAlreadyExists, "favorite already exists"},
		{ErrReviewNotFound, "review not found"},
		{ErrReviewAlreadyExists, "review already exists for this offer"},
		{ErrInvalidRating, "rating must be between 1 and 5"},
		{ErrEmptyReviewContent, "review content cannot be empty"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if tt.err.Error() != tt.expected {
				t.Errorf("Error.Error() = %v, want %v", tt.err.Error(), tt.expected)
			}
		})
	}
}

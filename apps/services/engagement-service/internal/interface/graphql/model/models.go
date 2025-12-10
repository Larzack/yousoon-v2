package model

import "time"

// Favorite représente un favori
type Favorite struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	OfferID   string    `json:"offerId"`
	CreatedAt time.Time `json:"createdAt"`
}

// FavoritesConnection pour la pagination
type FavoritesConnection struct {
	Edges      []*FavoriteEdge `json:"edges"`
	PageInfo   *PageInfo       `json:"pageInfo"`
	TotalCount int             `json:"totalCount"`
}

// FavoriteEdge pour la pagination
type FavoriteEdge struct {
	Node   *Favorite `json:"node"`
	Cursor string    `json:"cursor"`
}

// Review représente un avis
type Review struct {
	ID                 string    `json:"id"`
	UserID             string    `json:"userId"`
	OfferID            string    `json:"offerId"`
	Rating             int       `json:"rating"`
	Title              *string   `json:"title,omitempty"`
	Content            string    `json:"content"`
	Images             []string  `json:"images"`
	HelpfulCount       int       `json:"helpfulCount"`
	IsVerifiedPurchase bool      `json:"isVerifiedPurchase"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
}

// ReviewsConnection pour la pagination
type ReviewsConnection struct {
	Edges         []*ReviewEdge `json:"edges"`
	PageInfo      *PageInfo     `json:"pageInfo"`
	TotalCount    int           `json:"totalCount"`
	AverageRating *float64      `json:"averageRating,omitempty"`
}

// ReviewEdge pour la pagination
type ReviewEdge struct {
	Node   *Review `json:"node"`
	Cursor string  `json:"cursor"`
}

// PageInfo pour la pagination
type PageInfo struct {
	HasNextPage     bool    `json:"hasNextPage"`
	HasPreviousPage bool    `json:"hasPreviousPage"`
	StartCursor     *string `json:"startCursor,omitempty"`
	EndCursor       *string `json:"endCursor,omitempty"`
}

// CreateReviewInput input pour créer un avis
type CreateReviewInput struct {
	OfferID string   `json:"offerId"`
	Rating  int      `json:"rating"`
	Title   *string  `json:"title,omitempty"`
	Content string   `json:"content"`
	Images  []string `json:"images,omitempty"`
}

// UpdateReviewInput input pour mettre à jour un avis
type UpdateReviewInput struct {
	Rating  *int     `json:"rating,omitempty"`
	Title   *string  `json:"title,omitempty"`
	Content *string  `json:"content,omitempty"`
	Images  []string `json:"images,omitempty"`
}

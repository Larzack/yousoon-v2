package resolver

import (
	"github.com/yousoon/apps/services/engagement-service/internal/domain"
)

// Resolver est le resolver principal GraphQL
type Resolver struct {
	favoriteRepo domain.FavoriteRepository
	reviewRepo   domain.ReviewRepository
}

// NewResolver crée un nouveau resolver
func NewResolver(
	favoriteRepo domain.FavoriteRepository,
	reviewRepo domain.ReviewRepository,
) *Resolver {
	return &Resolver{
		favoriteRepo: favoriteRepo,
		reviewRepo:   reviewRepo,
	}
}

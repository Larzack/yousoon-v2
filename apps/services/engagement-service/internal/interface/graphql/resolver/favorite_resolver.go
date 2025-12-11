package resolver

import (
	"context"

	"github.com/yousoon/apps/services/engagement-service/internal/domain"
	"github.com/yousoon/apps/services/engagement-service/internal/interface/graphql/model"
)

// AddFavorite ajoute une offre aux favoris
func (r *Resolver) AddFavorite(ctx context.Context, offerID string) (*model.Favorite, error) {
	userID := ctx.Value("userID").(string)

	// TODO: Récupérer les infos de l'offre via gRPC pour la dénormalisation
	// Pour l'instant, on utilise des valeurs vides
	favorite, err := domain.NewFavorite(userID, offerID, "", "", "")
	if err != nil {
		return nil, err
	}

	if err := r.favoriteRepo.Create(ctx, favorite); err != nil {
		return nil, err
	}

	return &model.Favorite{
		ID:        favorite.ID(),
		UserID:    favorite.UserID(),
		OfferID:   favorite.OfferID(),
		CreatedAt: favorite.CreatedAt(),
	}, nil
}

// RemoveFavorite supprime une offre des favoris
func (r *Resolver) RemoveFavorite(ctx context.Context, offerID string) (bool, error) {
	userID := ctx.Value("userID").(string)

	// Trouver le favori par userID et offerID
	favorite, err := r.favoriteRepo.GetByUserAndOffer(ctx, userID, offerID)
	if err != nil {
		return false, err
	}

	if err := r.favoriteRepo.Delete(ctx, favorite.ID()); err != nil {
		return false, err
	}

	return true, nil
}

// MyFavorites retourne les favoris de l'utilisateur
func (r *Resolver) MyFavorites(ctx context.Context, first *int, after *string) (*model.FavoritesConnection, error) {
	userID := ctx.Value("userID").(string)

	limit := 20
	if first != nil {
		limit = *first
	}

	offset := 0
	// TODO: Décoder le cursor after pour la pagination

	filter := domain.FavoriteFilter{
		Offset: offset,
		Limit:  limit,
	}

	favorites, total, err := r.favoriteRepo.GetByUserID(ctx, userID, filter)
	if err != nil {
		return nil, err
	}

	edges := make([]*model.FavoriteEdge, len(favorites))
	for i, f := range favorites {
		edges[i] = &model.FavoriteEdge{
			Node: &model.Favorite{
				ID:        f.ID(),
				UserID:    f.UserID(),
				OfferID:   f.OfferID(),
				CreatedAt: f.CreatedAt(),
			},
			Cursor: f.ID(), // Utiliser l'ID comme cursor simple
		}
	}

	hasNext := len(favorites) == limit && offset+limit < int(total)

	return &model.FavoritesConnection{
		Edges: edges,
		PageInfo: &model.PageInfo{
			HasNextPage:     hasNext,
			HasPreviousPage: offset > 0,
		},
		TotalCount: int(total),
	}, nil
}

// IsFavorited vérifie si une offre est en favori
func (r *Resolver) IsFavorited(ctx context.Context, offerID string) (bool, error) {
	userID, ok := ctx.Value("userID").(string)
	if !ok || userID == "" {
		return false, nil
	}

	return r.favoriteRepo.ExistsByUserAndOffer(ctx, userID, offerID)
}

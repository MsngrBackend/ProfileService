package usecase

import (
	"context"
	"fmt"

	"github.com/MsngrBackend/ProfileService/internal/domain"
)

type FavoriteUsecase struct {
	repo domain.FavoriteRepository
}

func NewFavoriteUsecase(repo domain.FavoriteRepository) *FavoriteUsecase {
	return &FavoriteUsecase{repo: repo}
}

func (uc *FavoriteUsecase) List(ctx context.Context, userID string) ([]domain.Favorite, error) {
	favs, err := uc.repo.List(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list favorites: %w", err)
	}
	// return empty slice instead of nil so JSON encodes as [] not null
	if favs == nil {
		favs = []domain.Favorite{}
	}
	return favs, nil
}

func (uc *FavoriteUsecase) Add(ctx context.Context, userID, chatID string) error {
	if chatID == "" {
		return fmt.Errorf("chat_id is required")
	}
	return uc.repo.Add(ctx, userID, chatID)
}

func (uc *FavoriteUsecase) Remove(ctx context.Context, userID, chatID string) error {
	if chatID == "" {
		return fmt.Errorf("chat_id is required")
	}
	return uc.repo.Remove(ctx, userID, chatID)
}

func (uc *FavoriteUsecase) IsFavorite(ctx context.Context, userID, chatID string) (bool, error) {
	favs, err := uc.repo.List(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("check favorite: %w", err)
	}
	for _, f := range favs {
		if f.ChatID == chatID {
			return true, nil
		}
	}
	return false, nil
}

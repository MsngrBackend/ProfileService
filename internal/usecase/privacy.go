package usecase

import (
	"context"

	"github.com/MsngrBackend/ProfileService/internal/domain"
)

type PrivacyUsecase struct {
	repo domain.PrivacyRepository
}

func NewPrivacyUsecase(repo domain.PrivacyRepository) *PrivacyUsecase {
	return &PrivacyUsecase{repo: repo}
}

func (uc *PrivacyUsecase) Get(ctx context.Context, userID string) (*domain.PrivacySettings, error) {
	return uc.repo.Get(ctx, userID)
}

func (uc *PrivacyUsecase) Update(ctx context.Context, s *domain.PrivacySettings) error {
	return uc.repo.Update(ctx, s)
}

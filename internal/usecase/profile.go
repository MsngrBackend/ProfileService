package usecase

import (
	"context"

	"github.com/MsngrBackend/ProfileService/internal/domain"
)

type ProfileUsecase struct {
	repo    domain.ProfileRepository
	storage domain.AvatarStorage
}

func NewProfileUsecase(repo domain.ProfileRepository, storage domain.AvatarStorage) *ProfileUsecase {
	return &ProfileUsecase{repo: repo, storage: storage}
}

func (uc *ProfileUsecase) CreateProfile(ctx context.Context, userID string) (*domain.Profile, error) {
	return uc.repo.Create(ctx, userID)
}

func (uc *ProfileUsecase) GetProfile(ctx context.Context, userID string) (*domain.Profile, error) {
	return uc.repo.GetByID(ctx, userID)
}

func (uc *ProfileUsecase) GetProfileByUsername(ctx context.Context, username string) (*domain.Profile, error) {
	return uc.repo.GetByUsername(ctx, username)
}

func (uc *ProfileUsecase) UpdateProfile(ctx context.Context, userID, firstName, lastName, username, bio string) (*domain.Profile, error) {
	p, err := uc.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	p.FirstName = &firstName
	p.LastName = &lastName
	p.Bio = &bio
	if username != "" {
		p.Username = &username
	}
	if err := uc.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (uc *ProfileUsecase) UploadAvatar(ctx context.Context, userID string, data []byte, contentType string) (string, error) {
	url, err := uc.storage.Upload(ctx, userID, data, contentType)
	if err != nil {
		return "", err
	}
	return url, uc.repo.UpdateAvatarURL(ctx, userID, url)
}

func (uc *ProfileUsecase) DeleteAvatar(ctx context.Context, userID string) error {
	if err := uc.storage.Delete(ctx, userID); err != nil {
		return err
	}
	return uc.repo.UpdateAvatarURL(ctx, userID, "")
}

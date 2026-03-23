package usecase

import (
	"context"

	"github.com/MsngrBackend/ProfileService/internal/domain"
)

type ProfileUsecase struct {
	repo     domain.ProfileRepository
	storage  domain.AvatarStorage
	contacts domain.ContactsRepository
	privacy  domain.PrivacyRepository
}

func NewProfileUsecase(
	repo domain.ProfileRepository,
	contacts domain.ContactsRepository,
	privacy domain.PrivacyRepository,
	storage domain.AvatarStorage,
) *ProfileUsecase {
	return &ProfileUsecase{repo: repo, contacts: contacts, privacy: privacy, storage: storage}
}

func canView(visibility string, isOwner, isContact bool) bool {
	switch visibility {
	case "everyone":
		return true
	case "contacts":
		return isOwner || isContact
	case "nobody":
		return isOwner
	default:
		return false
	}
}

func (uc *ProfileUsecase) applyPrivacy(ctx context.Context, p *domain.Profile, viewerID string) (*domain.Profile, error) {
	isOwner := p.UserID == viewerID

	settings, err := uc.privacy.Get(ctx, p.UserID)
	if err != nil {
		return nil, err
	}

	isContact := false
	if !isOwner {
		contactList, err := uc.contacts.List(ctx, p.UserID)
		if err != nil {
			return nil, err
		}
		for _, c := range contactList {
			if c.ContactID == viewerID {
				isContact = true
				break
			}
		}
	}

	if !canView(settings.ProfileVisibility, isOwner, isContact) {
		return nil, domain.ErrProfileHidden
	}
	if !canView(settings.AvatarVisibility, isOwner, isContact) {
		p.AvatarURL = nil
	}
	if !canView(settings.LastSeenVisibility, isOwner, isContact) {
		p.LastSeenAt = nil
	}

	return p, nil
}

func (uc *ProfileUsecase) CreateProfile(ctx context.Context, userID string) (*domain.Profile, error) {
	return uc.repo.Create(ctx, userID)
}

func (uc *ProfileUsecase) GetProfile(ctx context.Context, userID, viewerID string) (*domain.Profile, error) {
	p, err := uc.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return uc.applyPrivacy(ctx, p, viewerID)
}

func (uc *ProfileUsecase) GetProfileByUsername(ctx context.Context, username, viewerID string) (*domain.Profile, error) {
	p, err := uc.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return uc.applyPrivacy(ctx, p, viewerID)
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

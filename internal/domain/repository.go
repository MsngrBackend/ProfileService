package domain

import "context"

type ProfileRepository interface {
	Create(ctx context.Context, userID string) (*Profile, error)
	GetByID(ctx context.Context, userID string) (*Profile, error)
	GetByUsername(ctx context.Context, username string) (*Profile, error)
	Update(ctx context.Context, p *Profile) error
	UpdateAvatarURL(ctx context.Context, userID, url string) error
}

type ContactsRepository interface {
	List(ctx context.Context, ownerID string) ([]Contact, error)
	Add(ctx context.Context, c Contact) error
	Remove(ctx context.Context, ownerID, contactID string) error
}

type PrivacyRepository interface {
	Get(ctx context.Context, userID string) (*PrivacySettings, error)
	Update(ctx context.Context, s *PrivacySettings) error
}

type NotificationRepository interface {
	Get(ctx context.Context, userID string, chatID *string) (*NotificationSettings, error)
	Upsert(ctx context.Context, s *NotificationSettings) error
}

type FavoriteRepository interface {
	List(ctx context.Context, userID string) ([]Favorite, error)
	Add(ctx context.Context, userID, chatID string) error
	Remove(ctx context.Context, userID, chatID string) error
}

type AvatarStorage interface {
	Upload(ctx context.Context, userID string, data []byte, contentType string) (string, error)
	Delete(ctx context.Context, userID string) error
}

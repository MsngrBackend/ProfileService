package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/MsngrBackend/ProfileService/internal/domain"
)

type NotificationUsecase struct {
	repo domain.NotificationRepository
}

func NewNotificationUsecase(repo domain.NotificationRepository) *NotificationUsecase {
	return &NotificationUsecase{repo: repo}
}

func (uc *NotificationUsecase) Get(ctx context.Context, userID string, chatID *string) (*domain.NotificationSettings, error) {
	settings, err := uc.repo.Get(ctx, userID, chatID)
	if err != nil {
		return nil, fmt.Errorf("get notification settings: %w", err)
	}
	return settings, nil
}

func (uc *NotificationUsecase) GetForChat(ctx context.Context, userID, chatID string) (*domain.NotificationSettings, error) {
	settings, err := uc.repo.Get(ctx, userID, &chatID)
	if err != nil {
		// fallback to global settings
		settings, err = uc.repo.Get(ctx, userID, nil)
		if err != nil {
			return nil, fmt.Errorf("get notification settings: %w", err)
		}
	}
	return settings, nil
}

func (uc *NotificationUsecase) Update(ctx context.Context, userID string, muted bool, mutedUntil *string) error {
	t, err := parseMutedUntil(mutedUntil)
	if err != nil {
		return err
	}
	return uc.repo.Upsert(ctx, &domain.NotificationSettings{
		UserID:     userID,
		Muted:      muted,
		MutedUntil: t,
	})
}

func (uc *NotificationUsecase) UpdateForChat(ctx context.Context, userID, chatID string, muted bool, mutedUntil *string) error {
	t, err := parseMutedUntil(mutedUntil)
	if err != nil {
		return err
	}
	return uc.repo.Upsert(ctx, &domain.NotificationSettings{
		UserID:     userID,
		ChatID:     &chatID,
		Muted:      muted,
		MutedUntil: t,
	})
}

// Returns nil if input is nil or empty — meaning "no expiry".
func parseMutedUntil(mutedUntil *string) (*time.Time, error) {
	if mutedUntil == nil || *mutedUntil == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, *mutedUntil)
	if err != nil {
		return nil, fmt.Errorf("invalid muted_until format, expected RFC3339: %w", err)
	}
	return &t, nil
}

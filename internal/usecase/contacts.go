package usecase

import (
	"context"

	"github.com/MsngrBackend/ProfileService/internal/domain"
)

type ContactsUsecase struct {
	repo domain.ContactsRepository
}

func NewContactsUsecase(repo domain.ContactsRepository) *ContactsUsecase {
	return &ContactsUsecase{repo: repo}
}

func (uc *ContactsUsecase) GetAllContacts(ctx context.Context, ownerID string) ([]domain.Contact, error) {
	return uc.repo.List(ctx, ownerID)
}

func (uc *ContactsUsecase) AddContact(ctx context.Context, contact domain.Contact) error {
	if err := uc.repo.Add(ctx, contact); err != nil {
		return err
	}
	return nil
}

func (uc *ContactsUsecase) DeleteContact(ctx context.Context, ownerID, contactID string) error {
	if err := uc.repo.Remove(ctx, ownerID, contactID); err != nil {
		return err
	}
	return nil
}

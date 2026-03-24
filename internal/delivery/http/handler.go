package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/MsngrBackend/ProfileService/internal/domain"
)

type Handler struct {
	profileUC      profileUsecase
	contactsUC     contactsUsecase
	privacyUC      privacyUsecase
	favoriteUC     favoriteUsecase
	notificationUC notificationsUsecase
}

type profileUsecase interface {
	CreateProfile(ctx context.Context, userID string) (*domain.Profile, error)
	GetProfile(ctx context.Context, userID string) (*domain.Profile, error)
	GetProfileByUsername(ctx context.Context, username string) (*domain.Profile, error)
	UpdateProfile(ctx context.Context, userID, firstName, lastName, username, bio string) (*domain.Profile, error)
	UploadAvatar(ctx context.Context, userID string, data []byte, contentType string) (string, error)
	DeleteAvatar(ctx context.Context, userID string) error
}

type privacyUsecase interface {
	Get(ctx context.Context, userID string) (*domain.PrivacySettings, error)
	Update(ctx context.Context, s *domain.PrivacySettings) error
}

type notificationsUsecase interface {
	Get(ctx context.Context, userID string, chatID *string) (*domain.NotificationSettings, error)
	GetForChat(ctx context.Context, userID, chatID string) (*domain.NotificationSettings, error)
	Update(ctx context.Context, userID string, muted bool, mutedUntil *string) error
	UpdateForChat(ctx context.Context, userID, chatID string, muted bool, mutedUntil *string) error
}

type favoriteUsecase interface {
	List(ctx context.Context, userID string) ([]domain.Favorite, error)
	Add(ctx context.Context, userID, chatID string) error
	Remove(ctx context.Context, userID, chatID string) error
	IsFavorite(ctx context.Context, userID, chatID string) (bool, error)
}

type contactsUsecase interface {
	GetAllContacts(ctx context.Context, ownerID string) ([]domain.Contact, error)
	AddContact(ctx context.Context, contact domain.Contact) error
	DeleteContact(ctx context.Context, ownerID, contactID string) error
}

func NewHandler(
	profileUC profileUsecase,
	contactsUC contactsUsecase,
	privacyUC privacyUsecase,
	favoriteUC favoriteUsecase,
	notificationUC notificationsUsecase,
) *Handler {
	return &Handler{
		profileUC:      profileUC,
		contactsUC:     contactsUC,
		privacyUC:      privacyUC,
		favoriteUC:     favoriteUC,
		notificationUC: notificationUC,
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, ErrorResponse{Error: msg})
}

func userIDFromCtx(r *http.Request) string {
	return r.Context().Value(UserIDKey).(string)
}

func ownerIDFromCtx(r *http.Request) string {
	v, ok := r.Context().Value(UserIDKey).(string)
	if !ok || v == "" {
		return ""
	}
	return v
}

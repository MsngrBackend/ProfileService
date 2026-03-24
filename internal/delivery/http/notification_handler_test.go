package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/MsngrBackend/ProfileService/internal/domain"
	"github.com/MsngrBackend/ProfileService/internal/usecase"
)

// ─── Mock ─────────────────────────────────────────────────────────────────────

type mockNotificationUC struct {
	get           func(ctx context.Context, userID string, chatID *string) (*domain.NotificationSettings, error)
	getForChat    func(ctx context.Context, userID, chatID string) (*domain.NotificationSettings, error)
	update        func(ctx context.Context, userID string, muted bool, mutedUntil *string) error
	updateForChat func(ctx context.Context, userID, chatID string, muted bool, mutedUntil *string) error
}

func (m *mockNotificationUC) Get(ctx context.Context, userID string, chatID *string) (*domain.NotificationSettings, error) {
	return m.get(ctx, userID, chatID)
}
func (m *mockNotificationUC) GetForChat(ctx context.Context, userID, chatID string) (*domain.NotificationSettings, error) {
	return m.getForChat(ctx, userID, chatID)
}
func (m *mockNotificationUC) Update(ctx context.Context, userID string, muted bool, mutedUntil *string) error {
	return m.update(ctx, userID, muted, mutedUntil)
}
func (m *mockNotificationUC) UpdateForChat(ctx context.Context, userID, chatID string, muted bool, mutedUntil *string) error {
	return m.updateForChat(ctx, userID, chatID, muted, mutedUntil)
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func newNotificationTestHandler(uc *mockNotificationUC) *Handler {
	return NewHandler(
		&mockProfileUC{},
		&mockContactsUC{},
		&mockPrivacyUC{},
		&usecase.FavoriteUsecase{},
		uc,
	)
}

func notificationFixture() *domain.NotificationSettings {
	return &domain.NotificationSettings{
		UserID: testOwnerID,
		Muted:  false,
	}
}

// ─── GetNotifications ─────────────────────────────────────────────────────────

func TestGetNotifications_Success(t *testing.T) {
	settings := notificationFixture()
	uc := &mockNotificationUC{
		get: func(_ context.Context, userID string, chatID *string) (*domain.NotificationSettings, error) {
			assert.Equal(t, testOwnerID, userID)
			assert.Nil(t, chatID)
			return settings, nil
		},
	}

	req := withOwner(httptest.NewRequest(http.MethodGet, "/notifications", nil))
	w := httptest.NewRecorder()
	newNotificationTestHandler(uc).GetNotifications(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var got domain.NotificationSettings
	require.NoError(t, json.NewDecoder(w.Body).Decode(&got))
	assert.Equal(t, settings.UserID, got.UserID)
	assert.Equal(t, settings.Muted, got.Muted)
}

func TestGetNotifications_UCError_Returns500(t *testing.T) {
	uc := &mockNotificationUC{
		get: func(_ context.Context, _ string, _ *string) (*domain.NotificationSettings, error) {
			return nil, errors.New("db error")
		},
	}

	req := withOwner(httptest.NewRequest(http.MethodGet, "/notifications", nil))
	w := httptest.NewRecorder()
	newNotificationTestHandler(uc).GetNotifications(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ─── GetChatNotifications ─────────────────────────────────────────────────────

func TestGetChatNotifications_Success(t *testing.T) {
	settings := notificationFixture()
	uc := &mockNotificationUC{
		getForChat: func(_ context.Context, userID, chatID string) (*domain.NotificationSettings, error) {
			assert.Equal(t, testOwnerID, userID)
			assert.Equal(t, "chat-1", chatID)
			return settings, nil
		},
	}

	req := withOwner(httptest.NewRequest(http.MethodGet, "/notifications/chat-1", nil))
	req.SetPathValue("chat_id", "chat-1")
	w := httptest.NewRecorder()
	newNotificationTestHandler(uc).GetChatNotifications(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var got domain.NotificationSettings
	require.NoError(t, json.NewDecoder(w.Body).Decode(&got))
	assert.Equal(t, settings.UserID, got.UserID)
}

func TestGetChatNotifications_UCError_Returns500(t *testing.T) {
	uc := &mockNotificationUC{
		getForChat: func(_ context.Context, _, _ string) (*domain.NotificationSettings, error) {
			return nil, errors.New("db error")
		},
	}

	req := withOwner(httptest.NewRequest(http.MethodGet, "/notifications/chat-1", nil))
	req.SetPathValue("chat_id", "chat-1")
	w := httptest.NewRecorder()
	newNotificationTestHandler(uc).GetChatNotifications(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ─── UpdateNotifications ──────────────────────────────────────────────────────

func TestUpdateNotifications_Success_Muted(t *testing.T) {
	uc := &mockNotificationUC{
		update: func(_ context.Context, userID string, muted bool, mutedUntil *string) error {
			assert.Equal(t, testOwnerID, userID)
			assert.True(t, muted)
			assert.Nil(t, mutedUntil)
			return nil
		},
	}

	body, _ := json.Marshal(UpdateNotificationsRequest{Muted: true})
	req := withOwner(httptest.NewRequest(http.MethodPut, "/notifications", bytes.NewReader(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	newNotificationTestHandler(uc).UpdateNotifications(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body.String())
}

func TestUpdateNotifications_Success_WithMutedUntil(t *testing.T) {
	validTime := "2099-01-01T00:00:00Z"
	uc := &mockNotificationUC{
		update: func(_ context.Context, userID string, muted bool, mutedUntil *string) error {
			assert.Equal(t, testOwnerID, userID)
			assert.True(t, muted)
			require.NotNil(t, mutedUntil)
			assert.Equal(t, validTime, *mutedUntil)
			return nil
		},
	}

	body, _ := json.Marshal(UpdateNotificationsRequest{Muted: true, MutedUntil: validTime})
	req := withOwner(httptest.NewRequest(http.MethodPut, "/notifications", bytes.NewReader(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	newNotificationTestHandler(uc).UpdateNotifications(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestUpdateNotifications_InvalidJSON_Returns400(t *testing.T) {
	called := false
	uc := &mockNotificationUC{
		update: func(_ context.Context, _ string, _ bool, _ *string) error {
			called = true
			return nil
		},
	}

	req := withOwner(httptest.NewRequest(http.MethodPut, "/notifications", bytes.NewBufferString("{")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	newNotificationTestHandler(uc).UpdateNotifications(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.False(t, called, "UC не должен вызываться при невалидном теле")
}

func TestUpdateNotifications_InvalidMutedUntil_Returns400(t *testing.T) {
	// UC возвращает ошибку парсинга из parseMutedUntil — хендлер отдаёт 400
	uc := &mockNotificationUC{
		update: func(_ context.Context, _ string, _ bool, _ *string) error {
			return errors.New("invalid muted_until format, expected RFC3339: ...")
		},
	}

	body, _ := json.Marshal(UpdateNotificationsRequest{Muted: true, MutedUntil: "not-a-date"})
	req := withOwner(httptest.NewRequest(http.MethodPut, "/notifications", bytes.NewReader(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	newNotificationTestHandler(uc).UpdateNotifications(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ─── UpdateChatNotifications ──────────────────────────────────────────────────

func TestUpdateChatNotifications_Success(t *testing.T) {
	uc := &mockNotificationUC{
		updateForChat: func(_ context.Context, userID, chatID string, muted bool, mutedUntil *string) error {
			assert.Equal(t, testOwnerID, userID)
			assert.Equal(t, "chat-1", chatID)
			assert.True(t, muted)
			assert.Nil(t, mutedUntil)
			return nil
		},
	}

	body, _ := json.Marshal(UpdateNotificationsRequest{Muted: true})
	req := withOwner(httptest.NewRequest(http.MethodPut, "/notifications/chat-1", bytes.NewReader(body)))
	req.SetPathValue("chat_id", "chat-1")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	newNotificationTestHandler(uc).UpdateChatNotifications(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body.String())
}

func TestUpdateChatNotifications_WithMutedUntil(t *testing.T) {
	validTime := "2099-06-01T12:00:00Z"
	uc := &mockNotificationUC{
		updateForChat: func(_ context.Context, _, _ string, muted bool, mutedUntil *string) error {
			assert.True(t, muted)
			require.NotNil(t, mutedUntil)
			assert.Equal(t, validTime, *mutedUntil)
			return nil
		},
	}

	body, _ := json.Marshal(UpdateNotificationsRequest{Muted: true, MutedUntil: validTime})
	req := withOwner(httptest.NewRequest(http.MethodPut, "/notifications/chat-1", bytes.NewReader(body)))
	req.SetPathValue("chat_id", "chat-1")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	newNotificationTestHandler(uc).UpdateChatNotifications(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestUpdateChatNotifications_InvalidJSON_Returns400(t *testing.T) {
	called := false
	uc := &mockNotificationUC{
		updateForChat: func(_ context.Context, _, _ string, _ bool, _ *string) error {
			called = true
			return nil
		},
	}

	req := withOwner(httptest.NewRequest(http.MethodPut, "/notifications/chat-1", bytes.NewBufferString("{")))
	req.SetPathValue("chat_id", "chat-1")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	newNotificationTestHandler(uc).UpdateChatNotifications(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.False(t, called, "UC не должен вызываться при невалидном теле")
}

func TestUpdateChatNotifications_InvalidMutedUntil_Returns400(t *testing.T) {
	uc := &mockNotificationUC{
		updateForChat: func(_ context.Context, _, _ string, _ bool, _ *string) error {
			return errors.New("invalid muted_until format, expected RFC3339: ...")
		},
	}

	body, _ := json.Marshal(UpdateNotificationsRequest{Muted: true, MutedUntil: "not-a-date"})
	req := withOwner(httptest.NewRequest(http.MethodPut, "/notifications/chat-1", bytes.NewReader(body)))
	req.SetPathValue("chat_id", "chat-1")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	newNotificationTestHandler(uc).UpdateChatNotifications(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

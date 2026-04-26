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

type mockPrivacyUC struct {
	get    func(ctx context.Context, userID string) (*domain.PrivacySettings, error)
	update func(ctx context.Context, s *domain.PrivacySettings) error
}

func (m *mockPrivacyUC) Get(ctx context.Context, userID string) (*domain.PrivacySettings, error) {
	return m.get(ctx, userID)
}
func (m *mockPrivacyUC) Update(ctx context.Context, s *domain.PrivacySettings) error {
	return m.update(ctx, s)
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func newPrivacyTestHandler(uc *mockPrivacyUC) *Handler {
	return NewHandler(
		&mockProfileUC{},
		&mockContactsUC{},
		uc,
		&usecase.FavoriteUsecase{},
		&usecase.NotificationUsecase{},
	)
}

func privacyFixture() *domain.PrivacySettings {
	return &domain.PrivacySettings{
		UserID:             testOwnerID,
		ProfileVisibility:  "everyone",
		LastSeenVisibility: "contacts",
		AvatarVisibility:   "everyone",
	}
}

// ─── GetPrivacy ───────────────────────────────────────────────────────────────

func TestGetPrivacy_Success(t *testing.T) {
	settings := privacyFixture()
	uc := &mockPrivacyUC{
		get: func(_ context.Context, userID string) (*domain.PrivacySettings, error) {
			assert.Equal(t, testOwnerID, userID)
			return settings, nil
		},
	}

	req := withOwner(httptest.NewRequest(http.MethodGet, "/privacy", nil))
	w := httptest.NewRecorder()
	newPrivacyTestHandler(uc).GetPrivacy(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var got domain.PrivacySettings
	require.NoError(t, json.NewDecoder(w.Body).Decode(&got))
	assert.Equal(t, *settings, got)
}

func TestGetPrivacy_UCError_Returns500(t *testing.T) {
	uc := &mockPrivacyUC{
		get: func(_ context.Context, _ string) (*domain.PrivacySettings, error) {
			return nil, errors.New("db error")
		},
	}

	req := withOwner(httptest.NewRequest(http.MethodGet, "/privacy", nil))
	w := httptest.NewRecorder()
	newPrivacyTestHandler(uc).GetPrivacy(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ─── UpdatePrivacy ────────────────────────────────────────────────────────────

func TestUpdatePrivacy_Success(t *testing.T) {
	var received *domain.PrivacySettings
	uc := &mockPrivacyUC{
		update: func(_ context.Context, s *domain.PrivacySettings) error {
			received = s
			return nil
		},
	}

	body, _ := json.Marshal(UpdatePrivacyRequest{
		ProfileVisibility:  "everyone",
		LastSeenVisibility: "contacts",
		AvatarVisibility:   "nobody",
	})
	req := withOwner(httptest.NewRequest(http.MethodPut, "/privacy", bytes.NewReader(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	newPrivacyTestHandler(uc).UpdatePrivacy(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body.String())

	require.NotNil(t, received)
	assert.Equal(t, testOwnerID, received.UserID)
	assert.Equal(t, "everyone", received.ProfileVisibility)
	assert.Equal(t, "contacts", received.LastSeenVisibility)
	assert.Equal(t, "nobody", received.AvatarVisibility)
}

func TestUpdatePrivacy_InvalidJSON_Returns400(t *testing.T) {
	called := false
	uc := &mockPrivacyUC{
		update: func(_ context.Context, _ *domain.PrivacySettings) error {
			called = true
			return nil
		},
	}

	req := withOwner(httptest.NewRequest(http.MethodPut, "/privacy", bytes.NewBufferString("{")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	newPrivacyTestHandler(uc).UpdatePrivacy(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.False(t, called, "UC не должен вызываться при невалидном теле")
}

func TestUpdatePrivacy_UCError_Returns500(t *testing.T) {
	uc := &mockPrivacyUC{
		update: func(_ context.Context, _ *domain.PrivacySettings) error {
			return errors.New("db error")
		},
	}

	body, _ := json.Marshal(UpdatePrivacyRequest{
		ProfileVisibility: "everyone",
	})
	req := withOwner(httptest.NewRequest(http.MethodPut, "/privacy", bytes.NewReader(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	newPrivacyTestHandler(uc).UpdatePrivacy(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

package http

import (
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

type mockFavoriteUC struct {
	list       func(ctx context.Context, userID string) ([]domain.Favorite, error)
	add        func(ctx context.Context, userID, chatID string) error
	remove     func(ctx context.Context, userID, chatID string) error
	isFavorite func(ctx context.Context, userID, chatID string) (bool, error)
}

func (m *mockFavoriteUC) List(ctx context.Context, userID string) ([]domain.Favorite, error) {
	return m.list(ctx, userID)
}
func (m *mockFavoriteUC) Add(ctx context.Context, userID, chatID string) error {
	return m.add(ctx, userID, chatID)
}
func (m *mockFavoriteUC) Remove(ctx context.Context, userID, chatID string) error {
	return m.remove(ctx, userID, chatID)
}
func (m *mockFavoriteUC) IsFavorite(ctx context.Context, userID, chatID string) (bool, error) {
	return m.isFavorite(ctx, userID, chatID)
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func newFavoriteTestHandler(uc *mockFavoriteUC) *Handler {
	return NewHandler(
		&mockProfileUC{},
		&mockContactsUC{},
		&mockPrivacyUC{},
		uc,
		&usecase.NotificationUsecase{},
	)
}

// ─── GetFavorites ─────────────────────────────────────────────────────────────

func TestGetFavorites_Success(t *testing.T) {
	favorites := []domain.Favorite{
		{UserID: testOwnerID, ChatID: "chat-1"},
		{UserID: testOwnerID, ChatID: "chat-2"},
	}
	uc := &mockFavoriteUC{
		list: func(_ context.Context, userID string) ([]domain.Favorite, error) {
			assert.Equal(t, testOwnerID, userID)
			return favorites, nil
		},
	}

	req := withOwner(httptest.NewRequest(http.MethodGet, "/favorites", nil))
	w := httptest.NewRecorder()
	newFavoriteTestHandler(uc).GetFavorites(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var got []domain.Favorite
	require.NoError(t, json.NewDecoder(w.Body).Decode(&got))
	assert.Equal(t, favorites, got)
}

func TestGetFavorites_EmptyList(t *testing.T) {
	uc := &mockFavoriteUC{
		list: func(_ context.Context, _ string) ([]domain.Favorite, error) {
			return []domain.Favorite{}, nil
		},
	}

	req := withOwner(httptest.NewRequest(http.MethodGet, "/favorites", nil))
	w := httptest.NewRecorder()
	newFavoriteTestHandler(uc).GetFavorites(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var got []domain.Favorite
	require.NoError(t, json.NewDecoder(w.Body).Decode(&got))
	assert.Empty(t, got)
}

func TestGetFavorites_UCError_Returns500(t *testing.T) {
	uc := &mockFavoriteUC{
		list: func(_ context.Context, _ string) ([]domain.Favorite, error) {
			return nil, errors.New("db error")
		},
	}

	req := withOwner(httptest.NewRequest(http.MethodGet, "/favorites", nil))
	w := httptest.NewRecorder()
	newFavoriteTestHandler(uc).GetFavorites(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ─── AddFavorite ──────────────────────────────────────────────────────────────

func TestAddFavorite_Success(t *testing.T) {
	uc := &mockFavoriteUC{
		add: func(_ context.Context, userID, chatID string) error {
			assert.Equal(t, testOwnerID, userID)
			assert.Equal(t, "chat-1", chatID)
			return nil
		},
	}

	req := withOwner(httptest.NewRequest(http.MethodPost, "/favorites/chat-1", nil))
	req.SetPathValue("chat_id", "chat-1")
	w := httptest.NewRecorder()
	newFavoriteTestHandler(uc).AddFavorite(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Empty(t, w.Body.String())
}

func TestAddFavorite_UCError_Returns500(t *testing.T) {
	uc := &mockFavoriteUC{
		add: func(_ context.Context, _, _ string) error {
			return errors.New("db error")
		},
	}

	req := withOwner(httptest.NewRequest(http.MethodPost, "/favorites/chat-1", nil))
	req.SetPathValue("chat_id", "chat-1")
	w := httptest.NewRecorder()
	newFavoriteTestHandler(uc).AddFavorite(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ─── RemoveFavorite ───────────────────────────────────────────────────────────

func TestRemoveFavorite_Success(t *testing.T) {
	uc := &mockFavoriteUC{
		remove: func(_ context.Context, userID, chatID string) error {
			assert.Equal(t, testOwnerID, userID)
			assert.Equal(t, "chat-1", chatID)
			return nil
		},
	}

	req := withOwner(httptest.NewRequest(http.MethodDelete, "/favorites/chat-1", nil))
	req.SetPathValue("chat_id", "chat-1")
	w := httptest.NewRecorder()
	newFavoriteTestHandler(uc).RemoveFavorite(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body.String())
}

func TestRemoveFavorite_UCError_Returns500(t *testing.T) {
	uc := &mockFavoriteUC{
		remove: func(_ context.Context, _, _ string) error {
			return errors.New("db error")
		},
	}

	req := withOwner(httptest.NewRequest(http.MethodDelete, "/favorites/chat-1", nil))
	req.SetPathValue("chat_id", "chat-1")
	w := httptest.NewRecorder()
	newFavoriteTestHandler(uc).RemoveFavorite(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

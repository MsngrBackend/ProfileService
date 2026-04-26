package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/MsngrBackend/ProfileService/internal/domain"
	"github.com/MsngrBackend/ProfileService/internal/usecase"
)

// ─── Mock ─────────────────────────────────────────────────────────────────────

type mockProfileUC struct {
	createProfile        func(ctx context.Context, userID string) (*domain.Profile, error)
	getProfile           func(ctx context.Context, userID string) (*domain.Profile, error)
	getProfileByUsername func(ctx context.Context, username string) (*domain.Profile, error)
	updateProfile        func(ctx context.Context, userID, firstName, lastName, username, bio string) (*domain.Profile, error)
	uploadAvatar         func(ctx context.Context, userID string, data []byte, contentType string) (string, error)
	deleteAvatar         func(ctx context.Context, userID string) error
}

func (m *mockProfileUC) CreateProfile(ctx context.Context, userID string) (*domain.Profile, error) {
	return m.createProfile(ctx, userID)
}
func (m *mockProfileUC) GetProfile(ctx context.Context, userID string) (*domain.Profile, error) {
	return m.getProfile(ctx, userID)
}
func (m *mockProfileUC) GetProfileByUsername(ctx context.Context, username string) (*domain.Profile, error) {
	return m.getProfileByUsername(ctx, username)
}
func (m *mockProfileUC) UpdateProfile(ctx context.Context, userID, firstName, lastName, username, bio string) (*domain.Profile, error) {
	return m.updateProfile(ctx, userID, firstName, lastName, username, bio)
}
func (m *mockProfileUC) UploadAvatar(ctx context.Context, userID string, data []byte, contentType string) (string, error) {
	return m.uploadAvatar(ctx, userID, data, contentType)
}
func (m *mockProfileUC) DeleteAvatar(ctx context.Context, userID string) error {
	return m.deleteAvatar(ctx, userID)
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func ptr(s string) *string { return &s }

// profileFixture возвращает минимально заполненный Profile с указателями.
func profileFixture() *domain.Profile {
	return &domain.Profile{
		UserID:    testOwnerID,
		Username:  ptr("alice"),
		FirstName: ptr("Alice"),
		LastName:  ptr("Smith"),
		Bio:       ptr("hello"),
	}
}

func newProfileTestHandler(uc *mockProfileUC) *Handler {
	return NewHandler(
		uc,
		// contacts — пустой мок, в этих тестах не вызывается
		&mockContactsUC{},
		&usecase.PrivacyUsecase{},
		&usecase.FavoriteUsecase{},
		&usecase.NotificationUsecase{},
	)
}

// multipartRequest формирует multipart/form-data запрос с файлом.
func multipartRequest(t *testing.T, fieldName, fileName string, content []byte) *http.Request {
	t.Helper()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	part, err := w.CreateFormFile(fieldName, fileName)
	require.NoError(t, err)
	_, err = part.Write(content)
	require.NoError(t, err)
	require.NoError(t, w.Close())

	req := httptest.NewRequest(http.MethodPost, "/profile/avatar", &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	return req
}

// ─── CreateProfile ────────────────────────────────────────────────────────────

func TestCreateProfile_Success(t *testing.T) {
	profile := profileFixture()
	uc := &mockProfileUC{
		createProfile: func(_ context.Context, userID string) (*domain.Profile, error) {
			assert.Equal(t, testOwnerID, userID)
			return profile, nil
		},
	}

	body, _ := json.Marshal(map[string]string{"user_id": testOwnerID})
	req := httptest.NewRequest(http.MethodPost, "/profile", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	newProfileTestHandler(uc).CreateProfile(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var got domain.Profile
	require.NoError(t, json.NewDecoder(w.Body).Decode(&got))
	assert.Equal(t, profile.UserID, got.UserID)
	assert.Equal(t, profile.Username, got.Username)
}

func TestCreateProfile_MissingUserID_Returns400(t *testing.T) {
	called := false
	uc := &mockProfileUC{
		createProfile: func(_ context.Context, _ string) (*domain.Profile, error) {
			called = true
			return nil, nil
		},
	}

	body, _ := json.Marshal(map[string]string{}) // user_id отсутствует
	req := httptest.NewRequest(http.MethodPost, "/profile", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	newProfileTestHandler(uc).CreateProfile(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.False(t, called, "UC не должен вызываться без user_id")
}

func TestCreateProfile_InvalidJSON_Returns400(t *testing.T) {
	uc := &mockProfileUC{}

	req := httptest.NewRequest(http.MethodPost, "/profile", bytes.NewBufferString("{"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	newProfileTestHandler(uc).CreateProfile(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateProfile_UCError_Returns500(t *testing.T) {
	uc := &mockProfileUC{
		createProfile: func(_ context.Context, _ string) (*domain.Profile, error) {
			return nil, errors.New("db error")
		},
	}

	body, _ := json.Marshal(map[string]string{"user_id": testOwnerID})
	req := httptest.NewRequest(http.MethodPost, "/profile", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	newProfileTestHandler(uc).CreateProfile(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ─── GetMyProfile ─────────────────────────────────────────────────────────────

func TestGetMyProfile_Success(t *testing.T) {
	profile := profileFixture()
	uc := &mockProfileUC{
		getProfile: func(_ context.Context, userID string) (*domain.Profile, error) {
			assert.Equal(t, testOwnerID, userID)
			return profile, nil
		},
	}

	req := withOwner(httptest.NewRequest(http.MethodGet, "/profile/me", nil))
	w := httptest.NewRecorder()
	newProfileTestHandler(uc).GetMyProfile(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var got domain.Profile
	require.NoError(t, json.NewDecoder(w.Body).Decode(&got))
	assert.Equal(t, profile.UserID, got.UserID)
	assert.Equal(t, profile.Username, got.Username)
}

func TestGetMyProfile_UCError_Returns404(t *testing.T) {
	uc := &mockProfileUC{
		getProfile: func(_ context.Context, _ string) (*domain.Profile, error) {
			return nil, errors.New("not found")
		},
	}

	req := withOwner(httptest.NewRequest(http.MethodGet, "/profile/me", nil))
	w := httptest.NewRecorder()
	newProfileTestHandler(uc).GetMyProfile(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ─── GetProfileByID ───────────────────────────────────────────────────────────

func TestGetProfileByID_Success(t *testing.T) {
	profile := profileFixture()
	uc := &mockProfileUC{
		getProfile: func(_ context.Context, userID string) (*domain.Profile, error) {
			assert.Equal(t, "user-42", userID)
			return profile, nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/profile/user-42", nil)
	req.SetPathValue("user_id", "user-42")
	w := httptest.NewRecorder()
	newProfileTestHandler(uc).GetProfileByID(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var got domain.Profile
	require.NoError(t, json.NewDecoder(w.Body).Decode(&got))
	assert.Equal(t, profile.UserID, got.UserID)
}

func TestGetProfileByID_UCError_Returns404(t *testing.T) {
	uc := &mockProfileUC{
		getProfile: func(_ context.Context, _ string) (*domain.Profile, error) {
			return nil, errors.New("not found")
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/profile/unknown", nil)
	req.SetPathValue("user_id", "unknown")
	w := httptest.NewRecorder()
	newProfileTestHandler(uc).GetProfileByID(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ─── GetProfileByUsername ─────────────────────────────────────────────────────

func TestGetProfileByUsername_Success(t *testing.T) {
	profile := profileFixture()
	uc := &mockProfileUC{
		getProfileByUsername: func(_ context.Context, username string) (*domain.Profile, error) {
			assert.Equal(t, "alice", username)
			return profile, nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/profile/@alice", nil)
	req.SetPathValue("username", "alice")
	w := httptest.NewRecorder()
	newProfileTestHandler(uc).GetProfileByUsername(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var got domain.Profile
	require.NoError(t, json.NewDecoder(w.Body).Decode(&got))
	assert.Equal(t, profile.UserID, got.UserID)
}

func TestGetProfileByUsername_MissingUsername_Returns400(t *testing.T) {
	called := false
	uc := &mockProfileUC{
		getProfileByUsername: func(_ context.Context, _ string) (*domain.Profile, error) {
			called = true
			return nil, nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/profile/", nil)
	// SetPathValue не вызывается → PathValue вернёт ""
	w := httptest.NewRecorder()
	newProfileTestHandler(uc).GetProfileByUsername(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.False(t, called)
}

func TestGetProfileByUsername_UCError_Returns404(t *testing.T) {
	uc := &mockProfileUC{
		getProfileByUsername: func(_ context.Context, _ string) (*domain.Profile, error) {
			return nil, errors.New("not found")
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/profile/@ghost", nil)
	req.SetPathValue("username", "ghost")
	w := httptest.NewRecorder()
	newProfileTestHandler(uc).GetProfileByUsername(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ─── UpdateProfile ────────────────────────────────────────────────────────────

func TestUpdateProfile_Success(t *testing.T) {
	profile := profileFixture()
	uc := &mockProfileUC{
		updateProfile: func(_ context.Context, userID, firstName, lastName, username, bio string) (*domain.Profile, error) {
			assert.Equal(t, testOwnerID, userID)
			assert.Equal(t, "Alice", firstName)
			assert.Equal(t, "Smith", lastName)
			assert.Equal(t, "alice", username)
			assert.Equal(t, "hello", bio)
			return profile, nil
		},
	}

	body, _ := json.Marshal(UpdateProfileRequest{
		FirstName: "Alice",
		LastName:  "Smith",
		Username:  "alice",
		Bio:       "hello",
	})
	req := withOwner(httptest.NewRequest(http.MethodPut, "/profile", bytes.NewReader(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	newProfileTestHandler(uc).UpdateProfile(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var got domain.Profile
	require.NoError(t, json.NewDecoder(w.Body).Decode(&got))
	assert.Equal(t, profile.UserID, got.UserID)
}

func TestUpdateProfile_InvalidJSON_Returns400(t *testing.T) {
	uc := &mockProfileUC{}

	req := withOwner(httptest.NewRequest(http.MethodPut, "/profile", bytes.NewBufferString("{")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	newProfileTestHandler(uc).UpdateProfile(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateProfile_UsernameTaken_Returns409(t *testing.T) {
	uc := &mockProfileUC{
		updateProfile: func(_ context.Context, _, _, _, _, _ string) (*domain.Profile, error) {
			return nil, errors.New("username already taken")
		},
	}

	body, _ := json.Marshal(UpdateProfileRequest{Username: "taken"})
	req := withOwner(httptest.NewRequest(http.MethodPut, "/profile", bytes.NewReader(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	newProfileTestHandler(uc).UpdateProfile(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestUpdateProfile_UCError_Returns500(t *testing.T) {
	uc := &mockProfileUC{
		updateProfile: func(_ context.Context, _, _, _, _, _ string) (*domain.Profile, error) {
			return nil, errors.New("db error")
		},
	}

	body, _ := json.Marshal(UpdateProfileRequest{Username: "alice"})
	req := withOwner(httptest.NewRequest(http.MethodPut, "/profile", bytes.NewReader(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	newProfileTestHandler(uc).UpdateProfile(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ─── UploadAvatar ─────────────────────────────────────────────────────────────

func TestUploadAvatar_Success(t *testing.T) {
	fileContent := []byte("fake-image-data")
	uc := &mockProfileUC{
		uploadAvatar: func(_ context.Context, userID string, data []byte, _ string) (string, error) {
			assert.Equal(t, testOwnerID, userID)
			assert.Equal(t, fileContent, data)
			return "https://cdn.example.com/avatar.jpg", nil
		},
	}

	req := withOwner(multipartRequest(t, "avatar", "avatar.jpg", fileContent))
	w := httptest.NewRecorder()
	newProfileTestHandler(uc).UploadAvatar(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]string
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.Equal(t, "https://cdn.example.com/avatar.jpg", resp["avatar_url"])
}

func TestUploadAvatar_MissingFile_Returns400(t *testing.T) {
	uc := &mockProfileUC{}

	// запрос без multipart-тела — FormFile вернёт ошибку
	req := withOwner(httptest.NewRequest(http.MethodPost, "/profile/avatar", nil))
	w := httptest.NewRecorder()
	newProfileTestHandler(uc).UploadAvatar(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUploadAvatar_UCError_Returns500(t *testing.T) {
	uc := &mockProfileUC{
		uploadAvatar: func(_ context.Context, _ string, _ []byte, _ string) (string, error) {
			return "", errors.New("storage unavailable")
		},
	}

	req := withOwner(multipartRequest(t, "avatar", "avatar.jpg", []byte("data")))
	w := httptest.NewRecorder()
	newProfileTestHandler(uc).UploadAvatar(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ─── DeleteAvatar ─────────────────────────────────────────────────────────────

func TestDeleteAvatar_Success(t *testing.T) {
	uc := &mockProfileUC{
		deleteAvatar: func(_ context.Context, userID string) error {
			assert.Equal(t, testOwnerID, userID)
			return nil
		},
	}

	req := withOwner(httptest.NewRequest(http.MethodDelete, "/profile/avatar", nil))
	w := httptest.NewRecorder()
	newProfileTestHandler(uc).DeleteAvatar(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body.String())
}

func TestDeleteAvatar_UCError_Returns500(t *testing.T) {
	uc := &mockProfileUC{
		deleteAvatar: func(_ context.Context, _ string) error {
			return errors.New("storage error")
		},
	}

	req := withOwner(httptest.NewRequest(http.MethodDelete, "/profile/avatar", nil))
	w := httptest.NewRecorder()
	newProfileTestHandler(uc).DeleteAvatar(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

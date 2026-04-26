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

type mockContactsUC struct {
	getAllContacts func(ctx context.Context, ownerID string) ([]domain.Contact, error)
	addContact     func(ctx context.Context, contact domain.Contact) error
	deleteContact  func(ctx context.Context, ownerID, contactID string) error
}

func (m *mockContactsUC) GetAllContacts(ctx context.Context, ownerID string) ([]domain.Contact, error) {
	return m.getAllContacts(ctx, ownerID)
}
func (m *mockContactsUC) AddContact(ctx context.Context, contact domain.Contact) error {
	return m.addContact(ctx, contact)
}
func (m *mockContactsUC) DeleteContact(ctx context.Context, ownerID, contactID string) error {
	return m.deleteContact(ctx, ownerID, contactID)
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

const testOwnerID = "owner-111"

func withOwner(r *http.Request) *http.Request {
	ctx := context.WithValue(r.Context(), UserIDKey, testOwnerID)
	return r.WithContext(ctx)
}

func newTestHandler(uc *mockContactsUC) *Handler {
	return NewHandler(
		&usecase.ProfileUsecase{},
		uc,
		&usecase.PrivacyUsecase{},
		&usecase.FavoriteUsecase{},
		&usecase.NotificationUsecase{},
	)
}

func jsonBody(t *testing.T, v any) *bytes.Reader {
	t.Helper()
	b, err := json.Marshal(v)
	require.NoError(t, err)
	return bytes.NewReader(b)
}

// ─── GetContacts ──────────────────────────────────────────────────────────────

func TestGetContacts_Success(t *testing.T) {
	expected := []domain.Contact{
		{OwnerID: testOwnerID, ContactID: "c-1", Alias: "Alice"},
		{OwnerID: testOwnerID, ContactID: "c-2", Alias: "Bob"},
	}

	uc := &mockContactsUC{
		getAllContacts: func(_ context.Context, ownerID string) ([]domain.Contact, error) {
			assert.Equal(t, testOwnerID, ownerID)
			return expected, nil
		},
	}

	req := withOwner(httptest.NewRequest(http.MethodGet, "/contacts", nil))
	w := httptest.NewRecorder()
	newTestHandler(uc).GetContacts(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var got []domain.Contact
	require.NoError(t, json.NewDecoder(w.Body).Decode(&got))
	assert.Equal(t, expected, got)
}

func TestGetContacts_UCError_Returns404(t *testing.T) {
	uc := &mockContactsUC{
		getAllContacts: func(_ context.Context, _ string) ([]domain.Contact, error) {
			return nil, errors.New("db error")
		},
	}

	req := withOwner(httptest.NewRequest(http.MethodGet, "/contacts", nil))
	w := httptest.NewRecorder()
	newTestHandler(uc).GetContacts(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ─── AddContact ───────────────────────────────────────────────────────────────

func TestAddContact_Success(t *testing.T) {
	var received domain.Contact

	uc := &mockContactsUC{
		addContact: func(_ context.Context, c domain.Contact) error {
			received = c
			return nil
		},
	}

	body := jsonBody(t, map[string]string{"contact_id": "c-3", "alias": "Carol"})
	req := withOwner(httptest.NewRequest(http.MethodPost, "/contacts", body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	newTestHandler(uc).AddContact(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, testOwnerID, received.OwnerID)
	assert.Equal(t, "c-3", received.ContactID)
	assert.Equal(t, "Carol", received.Alias)

	var resp domain.Contact
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.Equal(t, received, resp)
}

func TestAddContact_InvalidJSON_Returns400(t *testing.T) {
	called := false
	uc := &mockContactsUC{
		addContact: func(_ context.Context, _ domain.Contact) error {
			called = true
			return nil
		},
	}

	req := withOwner(httptest.NewRequest(http.MethodPost, "/contacts", bytes.NewBufferString("{")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	newTestHandler(uc).AddContact(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.False(t, called, "UC не должен вызываться при невалидном теле")
}

func TestAddContact_MissingContactID_Returns400(t *testing.T) {
	called := false
	uc := &mockContactsUC{
		addContact: func(_ context.Context, _ domain.Contact) error {
			called = true
			return nil
		},
	}

	body := jsonBody(t, map[string]string{"alias": "NoID"})
	req := withOwner(httptest.NewRequest(http.MethodPost, "/contacts", body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	newTestHandler(uc).AddContact(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.False(t, called, "UC не должен вызываться без contact_id")
}

func TestAddContact_UCError_Returns409(t *testing.T) {
	uc := &mockContactsUC{
		addContact: func(_ context.Context, _ domain.Contact) error {
			return errors.New("duplicate")
		},
	}

	body := jsonBody(t, map[string]string{"contact_id": "c-3"})
	req := withOwner(httptest.NewRequest(http.MethodPost, "/contacts", body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	newTestHandler(uc).AddContact(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

// ─── DeleteContact ────────────────────────────────────────────────────────────

func TestDeleteContact_Success(t *testing.T) {
	uc := &mockContactsUC{
		deleteContact: func(_ context.Context, ownerID, contactID string) error {
			assert.Equal(t, testOwnerID, ownerID)
			assert.Equal(t, "c-4", contactID)
			return nil
		},
	}

	req := withOwner(httptest.NewRequest(http.MethodDelete, "/contacts/c-4", nil))
	req.SetPathValue("contact_id", "c-4")
	w := httptest.NewRecorder()
	newTestHandler(uc).DeleteContact(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body.String())
}

func TestDeleteContact_MissingPathParam_Returns400(t *testing.T) {
	called := false
	uc := &mockContactsUC{
		deleteContact: func(_ context.Context, _, _ string) error {
			called = true
			return nil
		},
	}

	req := withOwner(httptest.NewRequest(http.MethodDelete, "/contacts/", nil))
	w := httptest.NewRecorder()
	newTestHandler(uc).DeleteContact(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.False(t, called, "UC не должен вызываться без contact_id")
}

func TestDeleteContact_UCError_Returns404(t *testing.T) {
	uc := &mockContactsUC{
		deleteContact: func(_ context.Context, _, _ string) error {
			return errors.New("not found")
		},
	}

	req := withOwner(httptest.NewRequest(http.MethodDelete, "/contacts/c-999", nil))
	req.SetPathValue("contact_id", "c-999")
	w := httptest.NewRecorder()
	newTestHandler(uc).DeleteContact(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

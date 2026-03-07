package http

import (
	"encoding/json"
	"net/http"

	"github.com/MsngrBackend/ProfileService/internal/domain"
)

func (h *Handler) GetContacts(w http.ResponseWriter, r *http.Request) {
	contacts, err := h.contactsUC.GetAllContacts(r.Context(), ownerIDFromCtx(r))
	if err != nil {
		writeError(w, http.StatusNotFound, "contacts not found")
		return
	}
	writeJSON(w, http.StatusOK, contacts)
}

func (h *Handler) AddContact(w http.ResponseWriter, r *http.Request) {
	var req AddContactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.ContactID == "" {
		writeError(w, http.StatusBadRequest, "contact_id is required")
		return
	}

	contact := domain.Contact{
		OwnerID:   ownerIDFromCtx(r),
		ContactID: req.ContactID,
		Alias:     req.Alias,
	}

	if err := h.contactsUC.AddContact(r.Context(), contact); err != nil {
		writeError(w, http.StatusConflict, "failed to add contact")
		return
	}

	writeJSON(w, http.StatusCreated, contact)
}

func (h *Handler) DeleteContact(w http.ResponseWriter, r *http.Request) {
	contactID := r.PathValue("contact_id")
	if contactID == "" {
		writeError(w, http.StatusBadRequest, "contact_id is required")
		return
	}

	if err := h.contactsUC.DeleteContact(r.Context(), ownerIDFromCtx(r), contactID); err != nil {
		writeError(w, http.StatusNotFound, "contact not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

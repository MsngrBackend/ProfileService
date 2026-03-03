package http

import (
	"encoding/json"
	"net/http"

	"github.com/MsngrBackend/ProfileService/internal/domain"
)

func (h *Handler) GetPrivacy(w http.ResponseWriter, r *http.Request) {
	settings, err := h.privacyUC.Get(r.Context(), userIDFromCtx(r))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get privacy settings")
		return
	}
	writeJSON(w, http.StatusOK, settings)
}

func (h *Handler) UpdatePrivacy(w http.ResponseWriter, r *http.Request) {
	var req UpdatePrivacyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	err := h.privacyUC.Update(r.Context(), &domain.PrivacySettings{
		UserID:             userIDFromCtx(r),
		ProfileVisibility:  req.ProfileVisibility,
		LastSeenVisibility: req.LastSeenVisibility,
		AvatarVisibility:   req.AvatarVisibility,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "update failed")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

package http

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	settings, err := h.notificationUC.Get(r.Context(), userIDFromCtx(r), nil)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get settings")
		return
	}
	writeJSON(w, http.StatusOK, settings)
}

func (h *Handler) GetChatNotifications(w http.ResponseWriter, r *http.Request) {
	chatID := r.PathValue("chat_id")
	settings, err := h.notificationUC.GetForChat(r.Context(), userIDFromCtx(r), chatID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get settings")
		return
	}
	writeJSON(w, http.StatusOK, settings)
}

func (h *Handler) UpdateNotifications(w http.ResponseWriter, r *http.Request) {
	var req UpdateNotificationsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}

	var mutedUntil *string
	if req.MutedUntil != "" {
		mutedUntil = &req.MutedUntil
	}

	if err := h.notificationUC.Update(r.Context(), userIDFromCtx(r), req.Muted, mutedUntil); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UpdateChatNotifications(w http.ResponseWriter, r *http.Request) {
	var req UpdateNotificationsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}

	var mutedUntil *string
	if req.MutedUntil != "" {
		mutedUntil = &req.MutedUntil
	}

	if err := h.notificationUC.UpdateForChat(r.Context(), userIDFromCtx(r), r.PathValue("chat_id"), req.Muted, mutedUntil); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

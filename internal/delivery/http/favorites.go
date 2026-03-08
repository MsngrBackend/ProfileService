package http

import (
	"net/http"
)

func (h *Handler) GetFavorites(w http.ResponseWriter, r *http.Request) {
	favorites, err := h.favoriteUC.List(r.Context(), userIDFromCtx(r))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get favorites")
		return
	}
	writeJSON(w, http.StatusOK, favorites)
}

func (h *Handler) AddFavorite(w http.ResponseWriter, r *http.Request) {
	chatID := r.PathValue("chat_id")
	if err := h.favoriteUC.Add(r.Context(), userIDFromCtx(r), chatID); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to add favorite")
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) RemoveFavorite(w http.ResponseWriter, r *http.Request) {
	chatID := r.PathValue("chat_id")
	if err := h.favoriteUC.Remove(r.Context(), userIDFromCtx(r), chatID); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to remove favorite")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

package http

import (
	"encoding/json"
	"io"
	"net/http"
)

func (h *Handler) CreateProfile(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID string `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.UserID == "" {
		writeError(w, http.StatusBadRequest, "user_id required")
		return
	}
	profile, err := h.profileUC.CreateProfile(r.Context(), req.UserID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "create failed")
		return
	}
	writeJSON(w, http.StatusCreated, profile)
}

func (h *Handler) GetMyProfile(w http.ResponseWriter, r *http.Request) {
	profile, err := h.profileUC.GetProfile(r.Context(), userIDFromCtx(r))
	if err != nil {
		writeError(w, http.StatusNotFound, "profile not found")
		return
	}
	writeJSON(w, http.StatusOK, profile)
}

func (h *Handler) GetProfileByID(w http.ResponseWriter, r *http.Request) {
	profile, err := h.profileUC.GetProfile(r.Context(), r.PathValue("user_id"))
	if err != nil {
		writeError(w, http.StatusNotFound, "profile not found")
		return
	}
	writeJSON(w, http.StatusOK, profile)
}

func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	updated, err := h.profileUC.UpdateProfile(r.Context(), userIDFromCtx(r), req.FirstName, req.LastName, req.Username, req.Bio)
	if err != nil {
		if err.Error() == "username already taken" {
			writeError(w, http.StatusConflict, "username already taken")
			return
		}
		writeError(w, http.StatusInternalServerError, "update failed")
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func (h *Handler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20) // 10MB
	file, header, err := r.FormFile("avatar")
	if err != nil {
		writeError(w, http.StatusBadRequest, "missing file")
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "read error")
		return
	}

	url, err := h.profileUC.UploadAvatar(r.Context(), userIDFromCtx(r), data, header.Header.Get("Content-Type"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "upload failed")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"avatar_url": url})
}

func (h *Handler) DeleteAvatar(w http.ResponseWriter, r *http.Request) {
	if err := h.profileUC.DeleteAvatar(r.Context(), userIDFromCtx(r)); err != nil {
		writeError(w, http.StatusInternalServerError, "delete failed")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

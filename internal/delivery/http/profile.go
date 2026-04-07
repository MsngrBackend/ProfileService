package http

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/MsngrBackend/ProfileService/internal/domain"
)

func viewerFromCtx(ctx context.Context) string {
	if id, ok := ctx.Value("userID").(string); ok {
		return id
	}
	return ""
}

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
	userID := userIDFromCtx(r)
	profile, err := h.profileUC.GetProfile(r.Context(), userID, userID)
	if err != nil {
		writeError(w, http.StatusNotFound, "profile not found")
		return
	}
	writeJSON(w, http.StatusOK, profile)
}

func (h *Handler) GetProfileByID(w http.ResponseWriter, r *http.Request) {
	profile, err := h.profileUC.GetProfile(r.Context(), r.PathValue("user_id"), viewerFromCtx(r.Context()))
	if err != nil {
		if errors.Is(err, domain.ErrProfileHidden) {
			writeError(w, http.StatusForbidden, "profile is not visible")
			return
		}
		writeError(w, http.StatusNotFound, "profile not found")
		return
	}
	writeJSON(w, http.StatusOK, profile)
}

func (h *Handler) GetProfileByUsername(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")
	if username == "" {
		writeError(w, http.StatusBadRequest, "username is required")
		return
	}
	profile, err := h.profileUC.GetProfileByUsername(r.Context(), username, viewerFromCtx(r.Context()))
	if err != nil {
		if errors.Is(err, domain.ErrProfileHidden) {
			writeError(w, http.StatusForbidden, "profile is not visible")
			return
		}
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
	uid := userIDFromCtx(r)
	updated, err := h.profileUC.UpdateProfile(r.Context(), uid, req.FirstName, req.LastName, req.Username, req.Bio)
	if err != nil {
		if err.Error() == "username already taken" {
			writeError(w, http.StatusConflict, "username already taken")
			return
		}
		writeError(w, http.StatusInternalServerError, "update failed")
		return
	}
	if h.profileEvents != nil {
		h.profileEvents.PublishProfileUpdated(uid, "profile")
	}
	writeJSON(w, http.StatusOK, updated)
}

func (h *Handler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
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

	uid := userIDFromCtx(r)
	url, err := h.profileUC.UploadAvatar(r.Context(), uid, data, header.Header.Get("Content-Type"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "upload failed")
		return
	}
	if h.profileEvents != nil {
		h.profileEvents.PublishProfileUpdated(uid, "avatar")
	}
	writeJSON(w, http.StatusOK, map[string]string{"avatar_url": url})
}

func (h *Handler) DeleteAvatar(w http.ResponseWriter, r *http.Request) {
	uid := userIDFromCtx(r)
	if err := h.profileUC.DeleteAvatar(r.Context(), uid); err != nil {
		writeError(w, http.StatusInternalServerError, "delete failed")
		return
	}
	if h.profileEvents != nil {
		h.profileEvents.PublishProfileUpdated(uid, "avatar")
	}
	w.WriteHeader(http.StatusNoContent)
}

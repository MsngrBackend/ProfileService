package http

import (
	"encoding/json"
	"net/http"

	"github.com/MsngrBackend/ProfileService/internal/usecase"
)

type Handler struct {
	profileUC      *usecase.ProfileUsecase
	// contactUC      *usecase.ContactUsecase
	privacyUC      *usecase.PrivacyUsecase
	// favoriteUC     *usecase.FavoriteUsecase
	// notificationUC *usecase.NotificationUsecase
	jwtSecret      string
}

func NewHandler(
	profileUC *usecase.ProfileUsecase,
	// contactUC *usecase.ContactUsecase,
	privacyUC *usecase.PrivacyUsecase,
	// favoriteUC *usecase.FavoriteUsecase,
	// notificationUC *usecase.NotificationUsecase,
	jwtSecret string,
) *Handler {
	return &Handler{
		profileUC:      profileUC,
		// contactUC:      contactUC,
		privacyUC:      privacyUC,
		// favoriteUC:     favoriteUC,
		// notificationUC: notificationUC,
		jwtSecret:      jwtSecret,
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, ErrorResponse{Error: msg})
}

func userIDFromCtx(r *http.Request) string {
	return r.Context().Value(UserIDKey).(string)
}

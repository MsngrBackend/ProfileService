package http

import (
	"context"
	"net/http"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userID := r.Header.Get("X-User-ID")
		if userID == "" {
			writeError(w, http.StatusUnauthorized, "missing user id")
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

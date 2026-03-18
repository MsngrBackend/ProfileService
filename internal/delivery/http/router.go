package http

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-User-ID")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

func recovererMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %v\n%s", err, debug.Stack())
				writeError(w, http.StatusInternalServerError, "internal server error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func chain(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

func (h *Handler) NewRouter() http.Handler {
	mux := http.NewServeMux()

	// Internal — service-to-service only, not proxied by gateway
	mux.HandleFunc("POST /internal/profiles", h.CreateProfile)

	auth := h.AuthMiddleware

	// Profile
	mux.Handle("GET /me", auth(http.HandlerFunc(h.GetMyProfile)))
	mux.Handle("PATCH /me", auth(http.HandlerFunc(h.UpdateProfile)))
	mux.Handle("GET /by-username/{username}", auth(http.HandlerFunc(h.GetProfileByUsername)))
	mux.Handle("GET /{user_id}", auth(http.HandlerFunc(h.GetProfileByID)))

	// Avatar
	mux.Handle("POST /me/avatar", auth(http.HandlerFunc(h.UploadAvatar)))
	mux.Handle("DELETE /me/avatar", auth(http.HandlerFunc(h.DeleteAvatar)))

	// Privacy
	mux.Handle("GET /me/privacy", auth(http.HandlerFunc(h.GetPrivacy)))
	mux.Handle("PUT /me/privacy", auth(http.HandlerFunc(h.UpdatePrivacy)))

	// Contacts
	mux.Handle("GET /contacts", auth(http.HandlerFunc(h.GetContacts)))
	mux.Handle("POST /contacts", auth(http.HandlerFunc(h.AddContact)))
	mux.Handle("DELETE /contacts/{contact_id}", auth(http.HandlerFunc(h.DeleteContact)))

	// Favorites
	mux.Handle("GET /favorites", auth(http.HandlerFunc(h.GetFavorites)))
	mux.Handle("POST /favorites/{chat_id}", auth(http.HandlerFunc(h.AddFavorite)))
	mux.Handle("DELETE /favorites/{chat_id}", auth(http.HandlerFunc(h.RemoveFavorite)))

	// Notifications
	mux.Handle("GET /notifications", auth(http.HandlerFunc(h.GetNotifications)))
	mux.Handle("GET /notifications/{chat_id}", auth(http.HandlerFunc(h.GetChatNotifications)))
	mux.Handle("PUT /notifications", auth(http.HandlerFunc(h.UpdateNotifications)))
	mux.Handle("PUT /notifications/{chat_id}", auth(http.HandlerFunc(h.UpdateChatNotifications)))

	return chain(mux, recovererMiddleware, loggerMiddleware, corsMiddleware)
}

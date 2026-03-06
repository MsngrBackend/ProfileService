package http

import (
	"net/http"
)

func (h *Handler) NewRouter() http.Handler {
	mux := http.NewServeMux()

	// wrap all routes with auth middleware
	auth := h.AuthMiddleware

	// Profile
	mux.Handle("GET /api/v1/profile/me", auth(http.HandlerFunc(h.GetMyProfile)))
	mux.Handle("PATCH /api/v1/profile/me", auth(http.HandlerFunc(h.UpdateProfile)))
	mux.Handle("GET /api/v1/profile/{user_id}", auth(http.HandlerFunc(h.GetProfileByID)))

	// Avatar
	mux.Handle("POST /api/v1/profile/me/avatar", auth(http.HandlerFunc(h.UploadAvatar)))
	mux.Handle("DELETE /api/v1/profile/me/avatar", auth(http.HandlerFunc(h.DeleteAvatar)))

	// Privacy
	mux.Handle("GET /api/v1/profile/me/privacy", auth(http.HandlerFunc(h.GetPrivacy)))
	mux.Handle("PUT /api/v1/profile/me/privacy", auth(http.HandlerFunc(h.UpdatePrivacy)))

	// Contacts
	// mux.Handle("GET /api/v1/profile/contacts", auth(http.HandlerFunc(h.GetContacts)))
	// mux.Handle("POST /api/v1/profile/contacts", auth(http.HandlerFunc(h.AddContact)))
	// mux.Handle("DELETE /api/v1/profile/contacts/{user_id}", auth(http.HandlerFunc(h.RemoveContact)))

	// Favorites
	// mux.Handle("GET /api/v1/profile/favorites", auth(http.HandlerFunc(h.GetFavorites)))
	// mux.Handle("POST /api/v1/profile/favorites/{chat_id}", auth(http.HandlerFunc(h.AddFavorite)))
	// mux.Handle("DELETE /api/v1/profile/favorites/{chat_id}", auth(http.HandlerFunc(h.RemoveFavorite)))

	// Notifications
	// mux.Handle("GET /api/v1/profile/notifications", auth(http.HandlerFunc(h.GetNotifications)))
	// mux.Handle("PUT /api/v1/profile/notifications", auth(http.HandlerFunc(h.UpdateNotifications)))
	// mux.Handle("PUT /api/v1/profile/notifications/{chat_id}", auth(http.HandlerFunc(h.UpdateChatNotifications)))

	return mux
}

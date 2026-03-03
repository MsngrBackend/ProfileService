package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (h *Handler) NewRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1/profile", func(r chi.Router) {
		r.Use(h.AuthMiddleware)

		// Profile
		r.Get("/me", h.GetMyProfile)
		r.Patch("/me", h.UpdateProfile)
		r.Get("/{user_id}", h.GetProfileByID)

		// Avatar
		r.Post("/me/avatar", h.UploadAvatar)
		r.Delete("/me/avatar", h.DeleteAvatar)

		// Privacy
		r.Get("/me/privacy", h.GetPrivacy)
		r.Put("/me/privacy", h.UpdatePrivacy)

		// Contacts
		// r.Get("/contacts", h.GetContacts)
		// r.Post("/contacts", h.AddContact)
		// r.Delete("/contacts/{user_id}", h.RemoveContact)

		// Favorites
		// r.Get("/favorites", h.GetFavorites)
		// r.Post("/favorites/{chat_id}", h.AddFavorite)
		// r.Delete("/favorites/{chat_id}", h.RemoveFavorite)

		// Notifications
		// r.Get("/notifications", h.GetNotifications)
		// r.Put("/notifications", h.UpdateNotifications)
		// r.Put("/notifications/{chat_id}", h.UpdateChatNotifications)
	})

	return r
}

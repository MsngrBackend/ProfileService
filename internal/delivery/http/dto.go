package http

type UpdateProfileRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Bio       string `json:"bio"`
}

type AddContactRequest struct {
	ContactID string `json:"contact_id"`
	Alias     string `json:"alias,omitempty"`
}

type UpdatePrivacyRequest struct {
	ProfileVisibility  string `json:"profile_visibility"`
	LastSeenVisibility string `json:"last_seen_visibility"`
	AvatarVisibility   string `json:"avatar_visibility"`
}

type UpdateNotificationsRequest struct {
	Muted      bool   `json:"muted"`
	MutedUntil string `json:"muted_until,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

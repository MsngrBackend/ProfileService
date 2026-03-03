package domain

import "time"

type Profile struct {
	UserID    string     `db:"user_id"    json:"user_id"`
	FirstName string     `db:"first_name" json:"first_name"`
	LastName  string     `db:"last_name"  json:"last_name"`
	Username  string     `db:"username"   json:"username"`
	Bio       string     `db:"bio"        json:"bio"`
	AvatarURL string     `db:"avatar_url" json:"avatar_url"`
	LastSeenAt *time.Time `db:"last_seen_at" json:"last_seen_at,omitempty"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
}

type Contact struct {
	OwnerID   string    `db:"owner_id"   json:"owner_id"`
	ContactID string    `db:"contact_id" json:"contact_id"`
	Alias     string    `db:"alias"      json:"alias,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type PrivacySettings struct {
	UserID             string `db:"user_id"              json:"user_id"`
	ProfileVisibility  string `db:"profile_visibility"   json:"profile_visibility"`
	LastSeenVisibility string `db:"last_seen_visibility" json:"last_seen_visibility"`
	AvatarVisibility   string `db:"avatar_visibility"    json:"avatar_visibility"`
}

type NotificationSettings struct {
	ID         string     `db:"id"          json:"id"`
	UserID     string     `db:"user_id"     json:"user_id"`
	ChatID     *string    `db:"chat_id"     json:"chat_id,omitempty"`
	Muted      bool       `db:"muted"       json:"muted"`
	MutedUntil *time.Time `db:"muted_until" json:"muted_until,omitempty"`
}

type Favorite struct {
	UserID    string    `db:"user_id"   json:"user_id"`
	ChatID    string    `db:"chat_id"   json:"chat_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

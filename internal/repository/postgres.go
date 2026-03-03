package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/MsngrBackend/ProfileService/internal/domain"
)

// ---- Profile ----

type ProfilePostgres struct{ db *sqlx.DB }

func NewProfilePostgres(db *sqlx.DB) *ProfilePostgres {
	return &ProfilePostgres{db: db}
}

func (r *ProfilePostgres) GetByID(ctx context.Context, userID string) (*domain.Profile, error) {
	var p domain.Profile
	err := r.db.GetContext(ctx, &p,
		`SELECT * FROM profiles WHERE user_id = $1`, userID)
	return &p, err
}

func (r *ProfilePostgres) Update(ctx context.Context, p *domain.Profile) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE profiles
		 SET first_name = $1, last_name = $2, bio = $3, updated_at = now()
		 WHERE user_id = $4`,
		p.FirstName, p.LastName, p.Bio, p.UserID)
	return err
}

func (r *ProfilePostgres) UpdateAvatarURL(ctx context.Context, userID, url string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE profiles SET avatar_url = $1, updated_at = now() WHERE user_id = $2`,
		url, userID)
	return err
}

// ---- Contacts ----

type ContactPostgres struct{ db *sqlx.DB }

func NewContactPostgres(db *sqlx.DB) *ContactPostgres {
	return &ContactPostgres{db: db}
}

func (r *ContactPostgres) List(ctx context.Context, ownerID string) ([]domain.Contact, error) {
	var contacts []domain.Contact
	err := r.db.SelectContext(ctx, &contacts,
		`SELECT * FROM contacts WHERE owner_id = $1 ORDER BY created_at DESC`, ownerID)
	return contacts, err
}

func (r *ContactPostgres) Add(ctx context.Context, c domain.Contact) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO contacts (owner_id, contact_id, alias)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (owner_id, contact_id) DO UPDATE SET alias = $3`,
		c.OwnerID, c.ContactID, c.Alias)
	return err
}

func (r *ContactPostgres) Remove(ctx context.Context, ownerID, contactID string) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM contacts WHERE owner_id = $1 AND contact_id = $2`,
		ownerID, contactID)
	return err
}

// ---- Privacy ----

type PrivacyPostgres struct{ db *sqlx.DB }

func NewPrivacyPostgres(db *sqlx.DB) *PrivacyPostgres {
	return &PrivacyPostgres{db: db}
}

func (r *PrivacyPostgres) Get(ctx context.Context, userID string) (*domain.PrivacySettings, error) {
	var s domain.PrivacySettings
	err := r.db.GetContext(ctx, &s,
		`SELECT * FROM privacy_settings WHERE user_id = $1`, userID)
	return &s, err
}

func (r *PrivacyPostgres) Update(ctx context.Context, s *domain.PrivacySettings) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO privacy_settings (user_id, profile_visibility, last_seen_visibility, avatar_visibility)
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT (user_id) DO UPDATE
		 SET profile_visibility   = $2,
		     last_seen_visibility = $3,
		     avatar_visibility    = $4`,
		s.UserID, s.ProfileVisibility, s.LastSeenVisibility, s.AvatarVisibility)
	return err
}

// ---- Favorites ----

type FavoritePostgres struct{ db *sqlx.DB }

func NewFavoritePostgres(db *sqlx.DB) *FavoritePostgres {
	return &FavoritePostgres{db: db}
}

func (r *FavoritePostgres) List(ctx context.Context, userID string) ([]domain.Favorite, error) {
	var favs []domain.Favorite
	err := r.db.SelectContext(ctx, &favs,
		`SELECT * FROM favorites WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	return favs, err
}

func (r *FavoritePostgres) Add(ctx context.Context, userID, chatID string) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO favorites (user_id, chat_id)
		 VALUES ($1, $2)
		 ON CONFLICT (user_id, chat_id) DO NOTHING`,
		userID, chatID)
	return err
}

func (r *FavoritePostgres) Remove(ctx context.Context, userID, chatID string) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM favorites WHERE user_id = $1 AND chat_id = $2`,
		userID, chatID)
	return err
}

// ---- Notifications ----

type NotificationPostgres struct{ db *sqlx.DB }

func NewNotificationPostgres(db *sqlx.DB) *NotificationPostgres {
	return &NotificationPostgres{db: db}
}

func (r *NotificationPostgres) Get(ctx context.Context, userID string, chatID *string) (*domain.NotificationSettings, error) {
	var s domain.NotificationSettings
	var err error
	if chatID == nil {
		err = r.db.GetContext(ctx, &s,
			`SELECT * FROM notification_settings
			 WHERE user_id = $1 AND chat_id IS NULL`, userID)
	} else {
		err = r.db.GetContext(ctx, &s,
			`SELECT * FROM notification_settings
			 WHERE user_id = $1 AND chat_id = $2`, userID, *chatID)
	}
	return &s, err
}

func (r *NotificationPostgres) Upsert(ctx context.Context, s *domain.NotificationSettings) error {
	var mutedUntil *time.Time
	if s.MutedUntil != nil {
		mutedUntil = s.MutedUntil
	}
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO notification_settings (id, user_id, chat_id, muted, muted_until)
		 VALUES (gen_random_uuid(), $1, $2, $3, $4)
		 ON CONFLICT (user_id, chat_id)
		 DO UPDATE SET muted = $3, muted_until = $4`,
		s.UserID, s.ChatID, s.Muted, mutedUntil)
	return err
}

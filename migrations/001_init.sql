-- +goose Up

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE visibility_type AS ENUM ('everyone', 'contacts', 'nobody');

CREATE TABLE profiles (
    user_id    UUID PRIMARY KEY,
    first_name VARCHAR(100),
    last_name  VARCHAR(100),
    username   VARCHAR(64) UNIQUE,
    bio        TEXT,
    avatar_url TEXT,
    last_seen_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE privacy_settings (
    user_id             UUID PRIMARY KEY REFERENCES profiles(user_id) ON DELETE CASCADE,
    profile_visibility  visibility_type NOT NULL DEFAULT 'everyone',
    last_seen_visibility visibility_type NOT NULL DEFAULT 'everyone',
    avatar_visibility   visibility_type NOT NULL DEFAULT 'everyone'
);

CREATE TABLE contacts (
    owner_id   UUID REFERENCES profiles(user_id) ON DELETE CASCADE,
    contact_id UUID REFERENCES profiles(user_id) ON DELETE CASCADE,
    alias      VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (owner_id, contact_id)
);

CREATE TABLE favorites (
    user_id    UUID REFERENCES profiles(user_id) ON DELETE CASCADE,
    chat_id    UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, chat_id)
);

CREATE TABLE notification_settings (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID REFERENCES profiles(user_id) ON DELETE CASCADE,
    chat_id     UUID,
    muted       BOOLEAN NOT NULL DEFAULT false,
    muted_until TIMESTAMPTZ,
    UNIQUE (user_id, chat_id)
);

-- Indexes
CREATE INDEX idx_profiles_username       ON profiles(username);
CREATE INDEX idx_contacts_owner          ON contacts(owner_id);
CREATE INDEX idx_contacts_contact        ON contacts(contact_id);
CREATE INDEX idx_notif_user_chat         ON notification_settings(user_id, chat_id);
CREATE INDEX idx_favorites_user          ON favorites(user_id);

-- +goose Down
DROP TABLE IF EXISTS notification_settings;
DROP TABLE IF EXISTS favorites;
DROP TABLE IF EXISTS contacts;
DROP TABLE IF EXISTS privacy_settings;
DROP TABLE IF EXISTS profiles;
DROP TYPE  IF EXISTS visibility_type;

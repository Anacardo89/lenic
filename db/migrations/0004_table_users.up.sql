CREATE TABLE users (
    id UUID PRIMARY KEY,
    username VARCHAR(32) NOT NULL UNIQUE,
    display_name VARCHAR(64),
    email VARCHAR(128) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL DEFAULT '',
    profile_pic VARCHAR(128) NOT NULL DEFAULT '',
    bio TEXT,
    user_followers INTEGER NOT NULL DEFAULT 0,
    user_following INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    user_role user_role NOT NULL DEFAULT 'user',
    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);
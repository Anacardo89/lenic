CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    username VARCHAR(32) NOT NULL UNIQUE,
    display_name VARCHAR(64),
    email VARCHAR(128) NOT NULL UNIQUE,

    password_hash VARCHAR(128) NOT NULL DEFAULT '',
    profile_pic_path VARCHAR(128) NOT NULL DEFAULT '',
    bio TEXT,

    user_followers INTEGER NOT NULL DEFAULT 0,
    user_following INTEGER NOT NULL DEFAULT 0,

    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    role VARCHAR(16) NOT NULL DEFAULT 'user',

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_login_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);
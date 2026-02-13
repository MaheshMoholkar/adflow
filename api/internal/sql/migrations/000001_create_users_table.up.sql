CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    phone VARCHAR(20) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL DEFAULT '',
    phone_verified BOOLEAN NOT NULL DEFAULT false,
    name VARCHAR(255),
    business_name VARCHAR(255),
    city VARCHAR(255),
    address TEXT,
    location_url TEXT,
    plan VARCHAR(20) NOT NULL DEFAULT 'none',
    plan_started_at TIMESTAMPTZ,
    plan_expires_at TIMESTAMPTZ,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_phone ON users(phone);

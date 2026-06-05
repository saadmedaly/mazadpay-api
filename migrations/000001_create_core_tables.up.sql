-- MazadPay — Migration 000001: Core Tables
-- Run with: migrate -path migrations -database "$DATABASE_URL" up
-- Requires: PostgreSQL 14+ / Neon (sslmode=require)

-- Enable UUID generation
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ============================================================
-- CATEGORIES
-- ============================================================
CREATE TABLE categories (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name_ar    VARCHAR(100) NOT NULL,
    name_fr    VARCHAR(100) NOT NULL,
    slug       VARCHAR(100) UNIQUE NOT NULL,
    icon       TEXT,
    sort_order INTEGER      NOT NULL DEFAULT 0,
    is_active  BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_categories_slug      ON categories(slug);
CREATE INDEX idx_categories_is_active ON categories(is_active);

-- ============================================================
-- USERS
-- ============================================================
CREATE TABLE users (
    id             UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    phone          VARCHAR(20)  UNIQUE NOT NULL,
    email          VARCHAR(255) UNIQUE,
    full_name      VARCHAR(255) NOT NULL,
    password_hash  VARCHAR(255) NOT NULL,   -- bcrypt, set in Milestone 3
    avatar_url     TEXT,
    is_verified    BOOLEAN      NOT NULL DEFAULT FALSE,
    is_active      BOOLEAN      NOT NULL DEFAULT TRUE,
    wallet_balance NUMERIC(12,2) NOT NULL DEFAULT 0.00,
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_phone     ON users(phone);
CREATE INDEX idx_users_email     ON users(email);
CREATE INDEX idx_users_is_active ON users(is_active);

-- ============================================================
-- AUCTIONS
-- ============================================================
CREATE TABLE auctions (
    id             UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id        UUID          NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id    UUID          NOT NULL REFERENCES categories(id),
    title_ar       VARCHAR(255)  NOT NULL,
    title_fr       VARCHAR(255),
    description_ar TEXT,
    description_fr TEXT,
    starting_price NUMERIC(12,2) NOT NULL,
    current_price  NUMERIC(12,2) NOT NULL,
    reserve_price  NUMERIC(12,2),
    deposit_amount NUMERIC(12,2) NOT NULL DEFAULT 0.00,
    bid_increment  NUMERIC(12,2) NOT NULL DEFAULT 100.00,
    -- pending | approved | active | ended | cancelled | rejected
    status         VARCHAR(30)   NOT NULL DEFAULT 'pending',
    admin_note     TEXT,
    starts_at      TIMESTAMPTZ,
    ends_at        TIMESTAMPTZ,
    total_bids     INTEGER       NOT NULL DEFAULT 0,
    view_count     INTEGER       NOT NULL DEFAULT 0,
    winner_id      UUID          REFERENCES users(id),
    created_at     TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_auctions_user_id     ON auctions(user_id);
CREATE INDEX idx_auctions_category_id ON auctions(category_id);
CREATE INDEX idx_auctions_status      ON auctions(status);
CREATE INDEX idx_auctions_ends_at     ON auctions(ends_at);
CREATE INDEX idx_auctions_active      ON auctions(status, ends_at)
    WHERE status = 'active';

-- ============================================================
-- AUCTION IMAGES
-- ============================================================
CREATE TABLE auction_images (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    auction_id UUID        NOT NULL REFERENCES auctions(id) ON DELETE CASCADE,
    url        TEXT        NOT NULL,   -- Cloudflare R2 public URL
    key        TEXT        NOT NULL,   -- R2 object key (for deletion)
    is_primary BOOLEAN     NOT NULL DEFAULT FALSE,
    sort_order INTEGER     NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_auction_images_auction_id ON auction_images(auction_id);

-- ============================================================
-- BIDS
-- ============================================================
CREATE TABLE bids (
    id         UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    auction_id UUID          NOT NULL REFERENCES auctions(id) ON DELETE CASCADE,
    user_id    UUID          NOT NULL REFERENCES users(id),
    amount     NUMERIC(12,2) NOT NULL,
    is_winning BOOLEAN       NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_bids_auction_id ON bids(auction_id);
CREATE INDEX idx_bids_user_id    ON bids(user_id);
-- Fast lookup for highest bid per auction
CREATE INDEX idx_bids_auction_amount ON bids(auction_id, amount DESC);

-- ============================================================
-- FAVORITES
-- ============================================================
CREATE TABLE favorites (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    auction_id UUID        NOT NULL REFERENCES auctions(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, auction_id)
);

CREATE INDEX idx_favorites_user_id ON favorites(user_id);

-- ============================================================
-- CONTACT MESSAGES
-- ============================================================
CREATE TABLE contact_messages (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name    VARCHAR(255) NOT NULL,
    contact_info VARCHAR(255) NOT NULL,   -- email or phone
    message      TEXT        NOT NULL,
    -- unread | read | replied
    status       VARCHAR(30) NOT NULL DEFAULT 'unread',
    ip_address   INET,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_contact_messages_status     ON contact_messages(status);
CREATE INDEX idx_contact_messages_created_at ON contact_messages(created_at DESC);

-- ============================================================
-- AUDIT LOGS
-- ============================================================
CREATE TABLE audit_logs (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    admin_id   UUID,                      -- NULL for system actions
    action     VARCHAR(100) NOT NULL,     -- e.g. auction.approve
    entity     VARCHAR(100) NOT NULL,     -- e.g. auction, user
    entity_id  UUID,
    details    JSONB,
    ip_address INET,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_logs_admin_id  ON audit_logs(admin_id);
CREATE INDEX idx_audit_logs_entity    ON audit_logs(entity, entity_id);
CREATE INDEX idx_audit_logs_created   ON audit_logs(created_at DESC);

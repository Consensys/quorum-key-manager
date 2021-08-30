BEGIN;

CREATE TABLE IF NOT EXISTS secrets (
    pk INTEGER PRIMARY KEY,
    id TEXT NOT NULL,
    version TEXT NOT NULL,
    store_id TEXT NOT NULL,
    disabled BOOLEAN default false,
    tags JSONB,
    created_at TIMESTAMPTZ DEFAULT (now() at time zone 'utc') NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT (now() at time zone 'utc') NOT NULL,
    deleted_at TIMESTAMPTZ,
    UNIQUE(id, version, store_id)
);

CREATE TABLE IF NOT EXISTS keys (
    pk INTEGER PRIMARY KEY,
    id TEXT NOT NULL,
    store_id TEXT NOT NULL,
    public_key BYTEA NOT NULL,
    signing_algorithm TEXT NOT NULL,
    elliptic_curve TEXT NOT NULL,
    tags JSONB,
    annotations JSONB,
    disabled BOOLEAN default false,
    created_at TIMESTAMPTZ DEFAULT (now() at time zone 'utc') NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT (now() at time zone 'utc') NOT NULL,
    deleted_at TIMESTAMPTZ,
    UNIQUE(id, store_id)
);

CREATE TABLE IF NOT EXISTS eth_accounts (
    pk INTEGER PRIMARY KEY,
    address TEXT NOT NULL,
    store_id TEXT NOT NULL,
    key_id TEXT NOT NULL,
    public_key BYTEA NOT NULL,
    compressed_public_key BYTEA NOT NULL,
    tags JSONB,
    disabled BOOLEAN default false,
    created_at TIMESTAMPTZ DEFAULT (now() at time zone 'utc') NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT (now() at time zone 'utc') NOT NULL,
    deleted_at TIMESTAMPTZ,
    UNIQUE(address, store_id)
);

CREATE TABLE IF NOT EXISTS aliases (
    pk INTEGER PRIMARY KEY,
    key TEXT NOT NULL,
    registry_name TEXT NOT NULL,
    value TEXT NOT NULL,
    UNIQUE(key, registry_name)
);

COMMIT;

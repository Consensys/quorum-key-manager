BEGIN;

CREATE TABLE IF NOT EXISTS secrets (
    id VARCHAR (50) NOT NULL,
    version VARCHAR (50) NOT NULL,
    store_id VARCHAR (50) NOT NULL,
    disabled BOOLEAN default false,
    tags JSONB,
    created_at TIMESTAMPTZ DEFAULT (now() at time zone 'utc') NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT (now() at time zone 'utc') NOT NULL,
    deleted_at TIMESTAMPTZ,
    PRIMARY KEY (id, version, store_id)
);

CREATE TABLE IF NOT EXISTS keys (
    id VARCHAR (50) NOT NULL,
    store_id VARCHAR (50) NOT NULL,
    public_key BYTEA NOT NULL,
    signing_algorithm VARCHAR (50) NOT NULL,
    elliptic_curve VARCHAR (50) NOT NULL,
    tags JSONB,
    annotations JSONB,
    disabled BOOLEAN default false,
    created_at TIMESTAMPTZ DEFAULT (now() at time zone 'utc') NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT (now() at time zone 'utc') NOT NULL,
    deleted_at TIMESTAMPTZ,
    PRIMARY KEY (id, store_id)
);

CREATE TABLE IF NOT EXISTS eth_accounts (
    address VARCHAR (42) NOT NULL,
    store_id VARCHAR (50) NOT NULL,
    key_id VARCHAR (50) NOT NULL,
    public_key BYTEA NOT NULL,
    compressed_public_key BYTEA NOT NULL,
    tags JSONB,
    disabled BOOLEAN default false,
    created_at TIMESTAMPTZ DEFAULT (now() at time zone 'utc') NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT (now() at time zone 'utc') NOT NULL,
    deleted_at TIMESTAMPTZ,
    PRIMARY KEY (address, store_id)
);

CREATE TABLE IF NOT EXISTS aliases (
    key VARCHAR (50) NOT NULL,
    registry_name VARCHAR (50) NOT NULL,
    value TEXT NOT NULL,
    PRIMARY KEY (key, registry_name)
);

COMMIT;

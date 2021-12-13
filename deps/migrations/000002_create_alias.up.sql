BEGIN;

CREATE TABLE IF NOT EXISTS aliases (
    pk SERIAL PRIMARY KEY,
    registry_name INTEGER NOT NULL REFERENCES registries(name) ON DELETE CASCADE,
    key TEXT NOT NULL,
    value JSON NOT NULL
    created_at TIMESTAMPTZ DEFAULT (now() at time zone 'utc') NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT (now() at time zone 'utc') NOT NULL,
    UNIQUE(key)
);

CREATE TABLE IF NOT EXISTS registries (
    pk SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
    allowed_tenants TEXT []
    created_at TIMESTAMPTZ DEFAULT (now() at time zone 'utc') NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT (now() at time zone 'utc') NOT NULL,
    UNIQUE(name)
);

COMMIT;

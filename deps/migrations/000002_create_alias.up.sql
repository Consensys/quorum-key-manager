BEGIN;

CREATE TABLE IF NOT EXISTS aliases (
    pk SERIAL PRIMARY KEY,
    registry_id INTEGER NOT NULL REFERENCES registries(pk) ON DELETE CASCADE,
    key TEXT NOT NULL,
    value JSON NOT NULL
    created_at TIMESTAMPTZ DEFAULT (now() at time zone 'utc') NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT (now() at time zone 'utc') NOT NULL,
    UNIQUE(key)
);

CREATE TABLE IF NOT EXISTS registries (
    pk SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
    created_at TIMESTAMPTZ DEFAULT (now() at time zone 'utc') NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT (now() at time zone 'utc') NOT NULL,
    UNIQUE(name)
);

COMMIT;

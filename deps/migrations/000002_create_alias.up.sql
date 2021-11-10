BEGIN;

CREATE TABLE IF NOT EXISTS aliases (
    pk SERIAL PRIMARY KEY,
    key TEXT NOT NULL,
    registry_name TEXT NOT NULL,
    value JSON NOT NULL,
    UNIQUE(key, registry_name)
);

COMMIT;

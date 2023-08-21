DROP TABLE IF EXISTS transaction_overrides;

CREATE TABLE transaction_overrides (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    uid VARCHAR(128),
    original_type INTEGER,
    overridden_type INTEGER,
    created_at TEXT,
    updated_at TEXT,
    UNIQUE(uid)
);

CREATE INDEX idx_uid ON transaction_overrides (uid);
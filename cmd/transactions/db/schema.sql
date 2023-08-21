DROP TABLE IF EXISTS calculator_transaction_overrides;

CREATE TABLE calculator_transaction_overrides (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    uid VARCHAR(128),
    overridden_type INTEGER,
    created_at TEXT,
    updated_at TEXT,
    UNIQUE(uid)
);

CREATE INDEX idx_uid ON calculator_transaction_overrides (uid);
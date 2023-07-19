DROP TABLE IF EXISTS markets;

CREATE TABLE markets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    `from` VARCHAR(10),
    `to` VARCHAR(10),
    `timestamp` INTEGER,
    open REAL,
    high REAL,
    low REAL,
    close REAL,
    UNIQUE(`timestamp`,`from`,`to`)
);

CREATE INDEX idx_timestamp ON markets (timestamp);
CREATE INDEX idx_from ON markets (`from`);
CREATE INDEX idx_to ON markets (`to`);
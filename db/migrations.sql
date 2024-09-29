-- 28 Sep 2024
CREATE TABLE IF NOT EXISTS urls (
    id VARCHAR(8) PRIMARY KEY,
    long_url TEXT NOT NULL UNIQUE
);

CREATE INDEX IF NOT EXISTS long_Url ON urls (
    long_url
);
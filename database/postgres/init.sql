CREATE TABLE short_links
(
    short_path VARCHAR NOT NULL UNIQUE,
    real_url VARCHAR NOT NULL
);
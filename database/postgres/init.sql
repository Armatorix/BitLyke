CREATE TABLE shorten_links
(
    shorten_path VARCHAR NOT NULL UNIQUE,
    real_url VARCHAR NOT NULL
);
CREATE TABLE short_links
(
    short_path VARCHAR NOT NULL UNIQUE,
    real_url VARCHAR NOT NULL
);

CREATE UNIQUE INDEX short_path_idx ON short_links(short_path);
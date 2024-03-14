CREATE TABLE IF NOT EXISTS actors (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    gender TEXT NOT NULL,
    date_of_birth TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

CREATE TABLE IF NOT EXISTS films (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    rating FLOAT NOT NULL
);

CREATE TABLE IF NOT EXISTS actors_from_films (
    actors_id BIGINT REFERENCES actors(id),
    films_id BIGINT REFERENCES films(id),
    PRIMARY KEY(actors_id, films_id)
);
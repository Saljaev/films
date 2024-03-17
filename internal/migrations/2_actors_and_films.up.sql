CREATE TABLE IF NOT EXISTS actors (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    gender TEXT NOT NULL,
    date_of_birth TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

ALTER TABLE actors ADD CONSTRAINT actor_unique UNIQUE (first_name, last_name, gender, date_of_birth);
CREATE INDEX first_name_idx ON actors(first_name);
CREATE INDEX last_name_idx ON actors(last_name);

CREATE TABLE IF NOT EXISTS films (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    name TEXT UNIQUE NOT NULL,
    description TEXT NOT NULL,
    rating FLOAT NOT NULL,
    release_date TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

CREATE UNIQUE INDEX film_name_idx ON films(name);

CREATE TABLE IF NOT EXISTS actors_from_films (
    actors_id BIGINT REFERENCES actors(id),
    films_id BIGINT REFERENCES films(id),
    PRIMARY KEY(actors_id, films_id)
);
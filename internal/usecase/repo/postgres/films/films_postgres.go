package filmsrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"
	"tiny/internal/entities"
	"tiny/internal/usecase"
)

type FilmsRepo struct {
	*sql.DB
}

func NewFilmsRepo(db *sql.DB) *FilmsRepo {
	return &FilmsRepo{db}
}

var _ usecase.FilmsRepo = (*FilmsRepo)(nil)

func (fr *FilmsRepo) Add(ctx context.Context, f entities.Films) (int, error) {
	const op = "FilmsRepo - Add"

	var filmId int

	query := "INSERT INTO films(name, description, rating, release_date) " +
		"VALUES($1, $2, $3, $4) RETURNING id"
	err := fr.QueryRowContext(ctx, query, f.Name, f.Description, f.Rating, f.ReleaseDate).Scan(&filmId)
	if err != nil {
		return 0, fmt.Errorf("%s - fr.QueryRowContext - films: %w", op, err)
	}

	tx, err := fr.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("%s - fr.BeginTx: %w", op, err)
	}

	defer tx.Rollback()

	for i := range f.Actors {
		var id int

		query = "INSERT INTO actors(first_name,last_name,gender,date_of_birth) " +
			"VALUES($1, $2, $3, $4) " +
			"ON CONFLICT DO NOTHING " +
			"RETURNING id"
		err = tx.QueryRowContext(ctx, query, f.Actors[i].FirstName, f.Actors[i].LastName, f.Actors[i].Gender, f.Actors[i].DateOfBirth).Scan(&id)
		if errors.Is(err, sql.ErrNoRows) {
			query = "SELECT id FROM actors " +
				"WHERE first_name = $1 AND last_name = $2 " +
				"AND gender = $3 AND date_of_birth = $4"
			err = fr.QueryRowContext(ctx, query, f.Actors[i].FirstName, f.Actors[i].LastName, f.Actors[i].Gender, f.Actors[i].DateOfBirth).Scan(&id)
			if err != nil {
				return 0, fmt.Errorf("%s - tx.QueryRowContext - actors: %w", op, err)
			}
		} else if err != nil {
			return 0, fmt.Errorf("%s - tx.QueryRowContext - actors: %w", op, err)
		}

		query = "INSERT INTO actors_from_films(actors_id, films_id) " +
			"VALUES($1,$2) " +
			"ON CONFLICT(actors_id, films_id) " +
			"DO NOTHING"
		_, err = tx.ExecContext(ctx, query, id, filmId)
		if err != nil {
			return 0, fmt.Errorf("%s - tx.ExecContext: %w", op, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("%s - tx.Commit: %w", op, err)
	}

	return filmId, nil
}

func (fr *FilmsRepo) GetById(ctx context.Context, id int) (*entities.Films, error) {
	const op = "FilmsRepo - Add"

	film := entities.Films{}

	query := "SELECT * FROM films WHERE id = $1"

	row := fr.QueryRowContext(ctx, query, id)
	err := row.Scan(&film.Id, &film.Name, &film.Description, &film.Rating, &film.ReleaseDate)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, os.ErrNotExist
	} else if err != nil {
		return nil, fmt.Errorf("%s - fr.QueryRowContext: %w", op, err)
	}

	return &film, nil

}

func (fr *FilmsRepo) Update(ctx context.Context, f entities.Films) error {
	const op = "FilmsRepo - Update"

	query := "UPDATE films SET " +
		"name = COALESCE(NULLIF($1, ''), name)," +
		"description = COALESCE(NULLIF($2, ''), description)," +
		"rating = COALESCE(NULLIF($3, 0.0), rating)," +
		"release_date = COALESCE(NULLIF($4, $5)::timestamp, release_date)" +
		"WHERE id = $6"

	var zeroTime time.Time

	res, err := fr.ExecContext(ctx, query, f.Name, f.Description, f.Rating, f.ReleaseDate, zeroTime, f.Id)
	if err != nil {
		return fmt.Errorf("%s - fr.QueryRowConext: %w", op, err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s - res.RowsAffected: %w", op, err)
	}

	// TODO: maybe change error
	if rowsAffected == 0 {
		return fmt.Errorf(errors.New("no data updating").Error())
	}

	tx, err := fr.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s - fr.BeginTx: %w", op, err)
	}

	defer tx.Rollback()

	for i := range f.Actors {
		var id int

		query = "INSERT INTO actors(first_name,last_name,gender,date_of_birth) " +
			"VALUES($1, $2, $3, $4) " +
			"ON CONFLICT DO NOTHING " +
			"RETURNING id"
		err = tx.QueryRowContext(ctx, query, f.Actors[i].FirstName, f.Actors[i].LastName, f.Actors[i].Gender, f.Actors[i].DateOfBirth).Scan(&id)
		if errors.Is(err, sql.ErrNoRows) {
			query = "SELECT id FROM actors " +
				"WHERE first_name = $1 AND last_name = $2 " +
				"AND gender = $3 AND date_of_birth = $4"
			err = fr.QueryRowContext(ctx, query, f.Actors[i].FirstName, f.Actors[i].LastName, f.Actors[i].Gender, f.Actors[i].DateOfBirth).Scan(&id)
			if err != nil {
				return fmt.Errorf("%s - tx.QueryRowContext - actors: %w", op, err)
			}
		} else if err != nil {
			return fmt.Errorf("%s - tx.QueryRowContext - actors: %w", op, err)
		}

		query = "INSERT INTO actors_from_films(actors_id, films_id) " +
			"VALUES($1,$2) " +
			"ON CONFLICT(actors_id, films_id) DO NOTHING"
		_, err = tx.ExecContext(ctx, query, id, f.Id)
		if err != nil {
			return fmt.Errorf("%s - tx.ExecContext: %w", op, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("%s - tx.Commit: %w", op, err)
	}

	return nil
}

func (fr *FilmsRepo) Delete(ctx context.Context, id int) error {
	const op = "FilmsRepo - Delete"

	tx, err := fr.BeginTx(ctx, nil)
	if err != nil {
		fmt.Errorf("%s - fr.BeginTx: %w", op, err)
	}

	defer tx.Rollback()

	query := "DELETE FROM actors_from_films WHERE films_id = $1"

	_, err = tx.ExecContext(ctx, query, id)
	if err != nil {
		fmt.Errorf("%s - tx.ExecContext - actors_from_films: %w", op, err)
	}

	query = "DELETE FROM films WHERE id = $1"
	res, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		fmt.Errorf("%s - tx.ExecContext - films: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s - res.RowsAffected: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s - res.RowsAffected: %w", op, os.ErrNotExist)
	}

	err = tx.Commit()
	if err != nil {
		fmt.Errorf("%s - tx.Commit: %w", op, err)
	}

	return nil
}

func (fr *FilmsRepo) SearchByFilmName(ctx context.Context, name string) ([]*entities.Films, error) {
	const op = "FilmsRepo - SearchByFilmName"

	query := "SELECT films.*, actors.*" +
		"FROM films " +
		"LEFT JOIN actors_from_films actorsfilm ON actorsfilm.films_id=films.id " +
		"LEFT JOIN actors ON actors.id = actorsfilm.actors_id " +
		"WHERE films.name ILIKE '%' || $1 || '%' " +
		"ORDER BY films.name, films.id, actors.id;"

	rows, err := fr.Query(query, name)
	if err != nil {
		return nil, fmt.Errorf("%s - fr.Query: %w", op, err)
	}

	defer rows.Close()

	var films []*entities.Films

	for rows.Next() {
		film := entities.Films{}
		actor := entities.Actors{}
		rows.Scan(&film.Id, &film.Name, &film.Description, &film.Rating, &film.ReleaseDate,
			&actor.Id, &actor.FirstName, &actor.LastName, &actor.Gender, &actor.DateOfBirth)

		film.Actors = append(film.Actors, &actor)
		films = append(films, &film)
	}

	return films, nil
}

func (fr *FilmsRepo) SearchByActorName(ctx context.Context, firstName, lastName string) ([]*entities.Films, error) {
	const op = "FilmsRepo - SearchByActorName"

	query := "SELECT films.*, actors.*" +
		"FROM films " +
		"LEFT JOIN actors_from_films actorsfilm ON actorsfilm.films_id=films.id " +
		"LEFT JOIN actors ON actors.id = actorsfilm.actors_id " +
		"WHERE actors.first_name ILIKE '%' || $1 || '%' OR actors.last_name ILIKE '%' || $2 || '%' " +
		"ORDER BY films.id, actors.id;"

	rows, err := fr.Query(query, firstName, lastName)
	if err != nil {
		return nil, fmt.Errorf("%s - fr.Query: %w", op, err)
	}

	defer rows.Close()

	var films []*entities.Films

	for rows.Next() {
		film := entities.Films{}
		actor := entities.Actors{}
		rows.Scan(&film.Id, &film.Name, &film.Description, &film.Rating, &film.ReleaseDate,
			&actor.Id, &actor.FirstName, &actor.LastName, &actor.Gender, &actor.DateOfBirth)

		film.Actors = append(film.Actors, &actor)
		films = append(films, &film)
	}

	return films, nil
}

func (fr *FilmsRepo) RateByField(ctx context.Context, fragment string, increasing string) ([]*entities.Films, error) {
	const op = "FilmsRepo - RateByField"

	query := "SELECT * FROM films ORDER BY " + fragment + " " + increasing

	rows, err := fr.Query(query)
	if err != nil {
		return nil, fmt.Errorf("%s - fr.Query: %w", op, err)
	}

	defer rows.Close()

	var films []*entities.Films

	for rows.Next() {
		film := entities.Films{}

		rows.Scan(&film.Id, &film.Name, &film.Description, &film.Rating, &film.ReleaseDate)

		films = append(films, &film)
	}

	return films, nil
}

func (fr *FilmsRepo) DeleteActor(ctx context.Context, filmID, actorID int) error {
	const op = "FilmsRepo - DeleteActor"

	query := "DELETE FROM actors_from_films WHERE actors_id = $1 AND films_id = $2"

	_, err := fr.ExecContext(ctx, query, actorID, filmID)
	if err != nil {
		return fmt.Errorf("%s - fr.ExecContext: %w", op, err)
	}

	return nil
}

package actorsrepo

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

type ActorsRepo struct {
	*sql.DB
}

func NewActorsRepo(db *sql.DB) *ActorsRepo {
	return &ActorsRepo{db}
}

var _ usecase.ActorsRepo = (*ActorsRepo)(nil)

func (ar *ActorsRepo) Add(ctx context.Context, a entities.Actors) (int, error) {
	const op = "ActorsRepo - Add"

	query := "INSERT INTO actors(first_name, last_name, gender, date_of_birth) " +
		"VALUES($1, $2, $3, $4) " +
		"RETURNING ID"

	var id int

	err := ar.QueryRowContext(ctx, query, a.FirstName, a.LastName, a.Gender, a.DateOfBirth).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s - ar.QueryRowContext: %w", op, err)
	}

	return id, nil
}

func (ar *ActorsRepo) GetById(ctx context.Context, id int) (*entities.Actors, error) {
	const op = "ActorsRepo - GetById"

	actor := entities.Actors{}

	query := "SELECT * FROM actors WHERE id = $1"

	row := ar.QueryRowContext(ctx, query, id)
	err := row.Scan(&actor.Id, &actor.FirstName, &actor.LastName, &actor.Gender, &actor.DateOfBirth)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, os.ErrNotExist
	} else if err != nil {
		return nil, fmt.Errorf("%s - fr.QueryRowContext: %w", op, err)
	}

	return &actor, nil
}

func (ar *ActorsRepo) Update(ctx context.Context, a entities.Actors) error {
	const op = "ActorsRepo - Update"

	query := "UPDATE actors SET " +
		"first_name = COALESCE(NULLIF($1, ''), first_name) " +
		"last_name = COALESCE(NULLIF($2, ''), last_name) " +
		"gender = COALESCE(NULLIF($3, ''), gender) " +
		"date_of_birth = COALESCE(NULLIF($4, $5)::timestamp, date_of_birth) " +
		"WHERE id = $6"

	var zeroTime time.Time

	res, err := ar.ExecContext(ctx, query, a.FirstName, a.LastName, a.Gender, a.DateOfBirth, zeroTime, a.Id)
	if err != nil {
		return fmt.Errorf("%s - ar.ExecContext: %w", op, err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s - res.RowsAffected: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf(errors.New("no data updating").Error())
	}

	return nil
}

func (ar *ActorsRepo) Delete(ctx context.Context, id int) error {
	const op = "ActorsRepo - Delete"

	tx, err := ar.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s - ar.BeginTx: %w", op, err)
	}

	defer tx.Rollback()

	query := "DELETE FROM actors_from_films WHERE actors_id = $1"
	_, err = tx.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s - tx.ExecContext: %w", op, err)
	}

	query = "DELETE FROM actors WHERE id = $1"
	res, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s - tx.ExecContext: %w", op, err)
	}

	resAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s - res.RowsAffected: %w", op, err)
	}

	if resAffected == 0 {
		return fmt.Errorf("%s - res.RowsAffected: %w", op, os.ErrNotExist)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("%s - tx.Commit: %w", op, err)
	}

	return nil
}

func (ar *ActorsRepo) GetAll(ctx context.Context) ([]*entities.Actors, error) {
	const op = "ActorsRepo - GetAll"

	query := "SELECT actors.*, films.* " +
		"FROM actors" +
		"LEFT JOIN actors_from_films actorsfilm ON actorsfilm.actors_id = actors.id " +
		"LEFT JOIN films ON films.id = actorsfilm.films_id " +
		"ORDER BY actors.id, films.id"

	rows, err := ar.Query(query)
	if err != nil {
		return nil, fmt.Errorf("%s - ar.Query: %w", op, err)
	}

	defer rows.Close()

	var actors []*entities.Actors

	for rows.Next() {
		film := entities.Films{}
		actor := entities.Actors{}
		rows.Scan(&actor.Id, &actor.FirstName, &actor.LastName, &actor.Gender, &actor.DateOfBirth,
			&film.Id, &film.Name, &film.Description, &film.Description, &film.ReleaseDate)

		actor.Films = append(actor.Films, &film)
		actors = append(actors, &actor)
	}

	return actors, nil
}

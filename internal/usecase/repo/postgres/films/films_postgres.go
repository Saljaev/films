package filmsrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

	query := "INSERT INTO films(name, description, rating, release_date) VALUES($1, $2, $3, $4) RETURNING id"
	err := fr.QueryRowContext(ctx, query, f.Name, f.Description, f.Rating, f.ReleaseDate).Scan(&filmId)
	if err != nil {
		return 0, fmt.Errorf("%s - fr.QueryRowContext - films: %w", op, err)
	}

	for i := range f.Actors {
		var id int

		query = "INSERT INTO actors(first_name,last_name,gender,date_of_birth) VALUES($1, $2, $3, $4) ON CONFLICT DO NOTHING RETURNING id"
		err = fr.QueryRowContext(ctx, query, f.Actors[i].FirstName, f.Actors[i].LastName, f.Actors[i].Gender, f.Actors[i].DateOfBirth).Scan(&id)
		if errors.Is(err, sql.ErrNoRows) {
			query = "SELECT id FROM actors WHERE first_name = $1 AND last_name = $2 AND gender = $3 AND date_of_birth = $4"
			err = fr.QueryRowContext(ctx, query, f.Actors[i].FirstName, f.Actors[i].LastName, f.Actors[i].Gender, f.Actors[i].DateOfBirth).Scan(&id)
			if err != nil {
				return 0, fmt.Errorf("%s - tx.QueryRowContext - actors: %w", op, err)
			}
		} else if err != nil {
			return 0, fmt.Errorf("%s - tx.QueryRowContext - actors: %w", op, err)
		}

		query = "INSERT INTO actors_from_films(actors_id, films_id) VALUES($1,$2) ON CONFLICT(actors_id, films_id) DO NOTHING"
		_, err = fr.ExecContext(ctx, query, id, filmId)
		if err != nil {
			return 0, fmt.Errorf("%s - tx.ExecContext: %w", op, err)
		}
	}

	return filmId, nil
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
		return fmt.Errorf("%s - fr.QueryRowConext: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf(errors.New("no data updating").Error())
	}

	return nil
}

func (fr *FilmsRepo) Delete(ctx context.Context, id int) error {
	panic("implement me")
}

func (fr *FilmsRepo) SearchByFragment(ctx context.Context, fragment, owner string) ([]*entities.Films, error) {
	panic("implement me")
}

func (fr *FilmsRepo) RateByField(ctx context.Context, fragment string, increasing string) ([]*entities.Films, error) {
	panic("implement me")
}

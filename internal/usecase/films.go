package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"tiny/internal/entities"
	"tiny/internal/models"
)

type FilmsUseCase struct {
	repo FilmsRepo
}

var _ Films = (*FilmsUseCase)(nil)

func NewFilmsUseCase(repo FilmsRepo) *FilmsUseCase {
	return &FilmsUseCase{repo}
}

func (fs *FilmsUseCase) Add(ctx context.Context, f models.Films) (int, error) {
	const op = "FilmsUseCase - Add"

	actors := []*entities.Actors{}
	for i := range f.Actors {
		actors = append(actors, &entities.Actors{
			FirstName:   f.Actors[i].FirstName,
			LastName:    f.Actors[i].LastName,
			Gender:      f.Actors[i].Gender,
			DateOfBirth: f.Actors[i].DateOfBirth,
		})
	}

	newFilm := entities.Films{
		Name:        f.Name,
		Description: f.Description,
		Rating:      f.Rating,
		ReleaseDate: f.ReleaseDate,
		Actors:      actors,
	}

	id, err := fs.repo.Add(ctx, newFilm)
	if err != nil {
		return 0, fmt.Errorf("%s - fc.repo.Add: %w", op, err)
	}

	return id, nil
}

func (fs *FilmsUseCase) GetById(ctx context.Context, id int) (*models.Films, error) {
	const op = "FilmsUseCase - Add"

	films, err := fs.repo.GetById(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, os.ErrNotExist
	} else if err != nil {
		return nil, fmt.Errorf("%s - fr.repo.GetById: %w", op, err)
	}

	actors := []*models.Actors{}
	for i := range films.Actors {
		actors = append(actors, &models.Actors{
			FirstName:   films.Actors[i].FirstName,
			LastName:    films.Actors[i].LastName,
			Gender:      films.Actors[i].Gender,
			DateOfBirth: films.Actors[i].DateOfBirth,
		})
	}

	newFilm := models.Films{
		Id:          films.Id,
		Name:        films.Name,
		Description: films.Description,
		Rating:      films.Rating,
		ReleaseDate: films.ReleaseDate,
		Actors:      actors,
	}

	return &newFilm, nil
}

func (fs *FilmsUseCase) Update(ctx context.Context, f models.Films) error {
	const op = "FilmsUseCase - Update"

	var actors []*entities.Actors

	for i := range f.Actors {
		actor := entities.Actors{
			FirstName:   f.Actors[i].FirstName,
			LastName:    f.Actors[i].LastName,
			Gender:      f.Actors[i].Gender,
			DateOfBirth: f.Actors[i].DateOfBirth,
		}

		actors = append(actors, &actor)
	}

	film := entities.Films{
		Id:          f.Id,
		Name:        f.Name,
		Description: f.Description,
		Rating:      f.Rating,
		ReleaseDate: f.ReleaseDate,
		Actors:      actors,
	}

	err := fs.repo.Update(ctx, film)
	if err != nil {
		return fmt.Errorf("%s - fs.repo.Update: %w", op, err)
	}

	return nil

}

func (fs *FilmsUseCase) Delete(ctx context.Context, id int) error {
	const op = "FilmsUseCase - Delete"

	err := fs.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("%s - fs.repo.Delete: %w", op, err)
	}

	return nil
}

func (fs *FilmsUseCase) SearchByFilmName(ctx context.Context, name string) ([]*models.Films, error) {
	const op = "FilmsUseCase - SearchByFilmName"

	films, err := fs.repo.SearchByFilmName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("%s - fs.repo.SearchByFilmName: %w", op, err)
	}

	var res []*models.Films

	for i := range films {
		film := models.Films{
			Id:          films[i].Id,
			Name:        films[i].Name,
			Description: films[i].Description,
			Rating:      films[i].Rating,
			ReleaseDate: films[i].ReleaseDate,
		}

		actor := models.Actors{
			Id:          films[i].Actors[0].Id,
			FirstName:   films[i].Actors[0].FirstName,
			LastName:    films[i].Actors[0].LastName,
			Gender:      films[i].Actors[0].Gender,
			DateOfBirth: films[i].Actors[0].DateOfBirth,
		}

		film.Actors = append(film.Actors, &actor)

		res = append(res, &film)
	}

	return res, nil
}

func (fs *FilmsUseCase) SearchByActorName(ctx context.Context, firstName, lastName string) ([]*models.Films, error) {
	const op = "FilmsUseCase - SearchByActorName"

	films, err := fs.repo.SearchByActorName(ctx, firstName, lastName)
	if err != nil {
		return nil, fmt.Errorf("%s - fs.repo.SearchByActorName: %w", op, err)
	}

	var res []*models.Films

	for i := range films {
		film := models.Films{
			Id:          films[i].Id,
			Name:        films[i].Name,
			Description: films[i].Description,
			Rating:      films[i].Rating,
			ReleaseDate: films[i].ReleaseDate,
		}

		actor := models.Actors{
			Id:          films[i].Actors[0].Id,
			FirstName:   films[i].Actors[0].FirstName,
			LastName:    films[i].Actors[0].LastName,
			Gender:      films[i].Actors[0].Gender,
			DateOfBirth: films[i].Actors[0].DateOfBirth,
		}

		film.Actors = append(film.Actors, &actor)

		res = append(res, &film)
	}

	return res, nil
}

func (fs *FilmsUseCase) RateByField(ctx context.Context, field string, increasing bool) ([]*models.Films, error) {
	const op = "FilmsUseCase - RateByField"

	var order string

	if increasing {
		order = "ASC"
	} else {
		order = "DESC"
	}

	films, err := fs.repo.RateByField(ctx, field, order)
	if err != nil {
		fmt.Errorf("%s - fs.repo.RateByField: %w", op, err)
	}

	var res []*models.Films

	for i := range films {
		film := models.Films{
			Name:        films[i].Name,
			Description: films[i].Description,
			Rating:      films[i].Rating,
			ReleaseDate: films[i].ReleaseDate,
		}

		res = append(res, &film)
	}

	return res, nil
}

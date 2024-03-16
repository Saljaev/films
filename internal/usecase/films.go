package usecase

import (
	"context"
	"fmt"
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

func (fs *FilmsUseCase) SearchByFragment(ctx context.Context, fragment, owner string) ([]*models.Films, error) {
	const op = "FilmsUseCase - SearchByFragment"

	films, err := fs.repo.SearchByFragment(ctx, fragment, owner)
	if err != nil {
		return nil, fmt.Errorf("%s - fs.repo.SearchByFragment: %w", op, err)
	}

	var res []*models.Films

	for i := range films {
		var actors []*models.Actor

		for j := range films {
			actors = append(actors, &models.Actor{
				FirstName: films[i].Actors[j].FirstName,
				LastName:  films[i].Actors[j].LastName,
			})
		}

		film := models.Films{
			Name:        films[i].Name,
			Description: films[i].Description,
			Rating:      films[i].Rating,
			Actors:      actors,
		}

		res = append(res, &film)
	}

	return res, nil
}

func (fs *FilmsUseCase) RateByField(ctx context.Context, fragment string, increasing bool) ([]*models.Films, error) {
	const op = "FilmsUseCase - RateByField"

	var order string

	if increasing {
		order = "ASC"
	} else {
		order = "DESC"
	}

	films, err := fs.repo.RateByField(ctx, fragment, order)
	if err != nil {
		fmt.Errorf("%s - fs.repo.RateByField: %w", op, err)
	}

	var res []*models.Films

	for i := range films {
		var actors []*models.Actor

		for j := range films {
			actors = append(actors, &models.Actor{
				FirstName: films[i].Actors[j].FirstName,
				LastName:  films[i].Actors[j].LastName,
			})
		}

		film := models.Films{
			Name:        films[i].Name,
			Description: films[i].Description,
			Rating:      films[i].Rating,
			Actors:      actors,
		}

		res = append(res, &film)
	}

	return res, nil
}

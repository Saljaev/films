package usecase

import (
	"context"
	"fmt"
	"tiny/internal/entities"
	"tiny/internal/models"
)

type ActorsUseCase struct {
	repo ActorsRepo
}

var _ Actors = (*ActorsUseCase)(nil)

func NewActorsUseCase(repo ActorsRepo) *ActorsUseCase {
	return &ActorsUseCase{repo}
}

func (ac *ActorsUseCase) Add(ctx context.Context, a models.Actors) (int, error) {
	const op = "ActorsUseCase - Add"

	newActor := entities.Actors{
		FirstName:   a.FirstName,
		LastName:    a.LastName,
		Gender:      a.Gender,
		DateOfBirth: a.DateOfBirth,
	}

	id, err := ac.repo.Add(ctx, newActor)
	if err != nil {
		return 0, fmt.Errorf("%s - ac.repo.Add: %w", op, err)
	}

	return id, nil
}

func (ac *ActorsUseCase) GetById(ctx context.Context, id int) (*models.Actors, error) {
	const op = "ActorsUserCase - GetById"

	actor, err := ac.repo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s - ac.repo.GetById: %w", op, err)
	}

	newActor := models.Actors{
		Id:          actor.Id,
		FirstName:   actor.FirstName,
		LastName:    actor.LastName,
		Gender:      actor.Gender,
		DateOfBirth: actor.DateOfBirth,
	}

	return &newActor, nil

}

func (ac *ActorsUseCase) Update(ctx context.Context, a models.Actors) error {
	const op = "ActorsUseCase - Update"

	actor := entities.Actors{
		Id:          a.Id,
		FirstName:   a.FirstName,
		LastName:    a.LastName,
		Gender:      a.Gender,
		DateOfBirth: a.DateOfBirth,
	}

	err := ac.repo.Update(ctx, actor)
	if err != nil {
		return fmt.Errorf("%s - ac.repo.Update: %w", op, err)
	}

	return nil
}

func (ac *ActorsUseCase) Delete(ctx context.Context, id int) error {
	const op = "ActorsUseCase - Delete"

	err := ac.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("%s - ac.repo.Delete: %w", op, err)
	}

	return nil
}

func (ac *ActorsUseCase) GetAll(ctx context.Context) ([]*models.Actors, error) {
	const op = "ActorsUseCase - GetAll"

	actors, err := ac.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s 0 ac.repo.GetAll: %w", op, err)
	}

	var res []*models.Actors

	for i := range actors {
		var films []*models.Films

		for j := range actors[i].Films {
			films = append(films, &models.Films{
				Name:        actors[i].Films[j].Name,
				Description: actors[i].Films[j].Description,
				Rating:      actors[i].Films[j].Rating,
				ReleaseDate: actors[i].Films[j].ReleaseDate,
			})
		}

		actor := models.Actors{
			FirstName:   actors[i].FirstName,
			LastName:    actors[i].LastName,
			Gender:      actors[i].Gender,
			DateOfBirth: actors[i].DateOfBirth,
			Films:       films,
		}

		res = append(res, &actor)
	}

	return res, nil
}

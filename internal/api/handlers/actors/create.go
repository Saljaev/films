package actorshandler

import (
	"context"
	"net/http"
	"time"
	"tiny/internal/api/utilapi"
	"tiny/internal/models"
)

type ActorCreate struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"second_name"`
	DateOfBirth string `json:"date_of_birth"`
	Gender      string `json:"gender"`
}

type ActorCreateResponse struct {
	ActorID int `json:"actor_id"`
}

func (req *ActorCreate) IsValid() bool {
	if req.FirstName == "" || req.LastName == "" || req.Gender == "" || req.DateOfBirth == "" {
		return false
	}
	if !ValidateDate(req.DateOfBirth) {
		return false
	}
	return ValidateGender(req.Gender)
}

func (h *ActorsHandler) Create(ctx *utilapi.APIContext) {
	var req ActorCreate
	err := ctx.Decode(&req)
	if err != nil {
		return
	}

	ctx.Info("request decoded", "request", req)

	date, err := time.Parse(time.DateOnly, req.DateOfBirth)

	actor := models.Actors{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Gender:      req.Gender,
		DateOfBirth: date,
	}

	id, err := h.actors.Add(context.Background(), actor)
	if err != nil {
		ctx.Error("failed to add actor", err)
		ctx.WriteFailure(http.StatusInternalServerError, "failed to add actor")
		return
	}

	ctx.Info("actor add", "id", id)
	ctx.SuccessWithData(ActorCreateResponse{ActorID: id})
}

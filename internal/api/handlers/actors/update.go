package actorshandler

import (
	"context"
	"net/http"
	"strconv"
	"time"
	"tiny/internal/api/utilapi"
	"tiny/internal/models"
)

type ActorsUpdateRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"second_name"`
	DateOfBirth string `json:"date_of_birth"`
	Gender      string `json:"gender"`
}

type ActorsUpdateResponse struct {
	ActorID int `json:"actor_id"`
}

func (req *ActorsUpdateRequest) IsValid() bool {
	if req.FirstName == "" && req.LastName == "" &&
		req.Gender == "" && req.DateOfBirth == "" {
		return false
	}
	if req.Gender != "" && !ValidateGender(req.Gender) {
		return false
	}
	return ValidateDate(req.DateOfBirth)
}

func (h *ActorsHandler) Update(ctx *utilapi.APIContext) {
	var req ActorsUpdateRequest
	err := ctx.Decode(&req)
	if err != nil {
		return
	}

	rawActorID := ctx.GetURLParam("id")

	actorID, err := strconv.Atoi(rawActorID)
	if err != nil || actorID <= 0 {
		ctx.Error("invalid actor id from url param", err)
		ctx.WriteFailure(http.StatusBadRequest, "invalid actor id")
		return
	}

	date, _ := time.Parse(time.DateOnly, req.DateOfBirth)

	ctx.Info("request decoded", "request", req)

	actor := models.Actors{
		Id:          actorID,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Gender:      req.Gender,
		DateOfBirth: date,
	}

	err = h.actors.Update(context.Background(), actor)
	if err != nil {
		ctx.Error("failed to update actor", err)
		ctx.WriteFailure(http.StatusInternalServerError, "failed to update actor")
		return
	}

	ctx.Info("actor update", "id", actorID)
	ctx.SuccessWithData(ActorsUpdateResponse{ActorID: actorID})
}

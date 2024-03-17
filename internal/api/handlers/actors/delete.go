package actorshandler

import (
	"context"
	"net/http"
	"strconv"
	"tiny/internal/api/utilapi"
)

type ActorDeletResponse struct {
	ActorID int `json:"actor_id"`
}

func (h *ActorsHandler) Delete(ctx *utilapi.APIContext) {
	rawActorID := ctx.GetURLParam("id")

	actorID, err := strconv.Atoi(rawActorID)
	if err != nil || actorID <= 0 {
		ctx.Error("invalid actor id from url param", err)
		ctx.WriteFailure(http.StatusBadRequest, "invalid actor id")
		return
	}

	actor, err := h.actors.GetById(context.Background(), actorID)
	if err != nil {
		ctx.Error("failed to delete actor by id", err)
		ctx.WriteFailure(http.StatusBadRequest, "invalid actor id")
		return
	}

	err = h.actors.Delete(context.Background(), actor.Id)
	if err != nil {
		ctx.Error("failed to delete actor", err)
		ctx.WriteFailure(http.StatusInternalServerError, "server error")
		return
	}

	ctx.Info("actor delete", "actor_id", actorID)
	ctx.SuccessWithData(ActorDeletResponse{ActorID: actorID})
}

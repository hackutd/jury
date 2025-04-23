package router

import (
	"server/judging"
	"server/logging"
	"server/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type State struct {
	Db      *mongo.Database
	Clock   *models.SafeClock
	Comps   *judging.Comparisons
	Logger  *logging.Logger
	Limiter *Limiter
}

func NewState(db *mongo.Database, clock *models.SafeClock, comps *judging.Comparisons, logger *logging.Logger, limiter *Limiter) *State {
	return &State{
		Db:      db,
		Clock:   clock,
		Comps:   comps,
		Logger:  logger,
		Limiter: limiter,
	}
}

func GetState(ctx *gin.Context) *State {
	state := ctx.MustGet("state").(*State)
	return state
}

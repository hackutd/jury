package util

import (
	"crypto/rand"
	"math/big"
	"server/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// GenerateToken generates a random token of length 16
func GenerateToken() (string, error) {
	b := make([]rune, 16)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		b[i] = letters[n.Int64()]
	}
	return string(b), nil
}

// GetFullHostname returns the full hostname of the request, including scheme
func GetFullHostname(ctx *gin.Context) string {
	scheme := "http"
	if ctx.Request.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + ctx.Request.Host
}

// Now returns the current time as a primitive.DateTime
func Now() primitive.DateTime {
	return primitive.NewDateTimeFromTime(time.Now())
}

// FindSeenProjectIndex finds the index of a project in a judge's seen projects.
// Will return -1 if the project is not found.
func FindSeenProjectIndex(judge *models.Judge, projectId primitive.ObjectID) int {
	for i, seen := range judge.SeenProjects {
		if seen.ProjectId == projectId {
			return i
		}
	}
	return -1
}

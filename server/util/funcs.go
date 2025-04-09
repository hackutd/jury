package util

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"reflect"
	"server/models"
	"strings"
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

// Gets the group from the table number
func GroupFromTable(op *models.Options, table int64) int64 {
	group := int64(0)
	for group < int64(len(op.GroupSizes)) && table >= op.GroupSizes[group] {
		table -= op.GroupSizes[group]
		group++
	}
	if group == int64(len(op.GroupSizes)) {
		group--
	}
	return group
}

// RankingToString converts a ranking array to a string
func RankingToString(rankings []primitive.ObjectID) string {
	sb := strings.Builder{}
	sb.WriteString("[")
	for i, id := range rankings {
		if i != 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(id.Hex())
	}
	sb.WriteString("]")

	return sb.String()
}

// StructToString converts a struct to a string
func StructToString(s interface{}) string {
	sb := strings.Builder{}
	sb.WriteString("{")
	v := reflect.ValueOf(s)
	for i := 0; i < v.NumField(); i++ {
		if i != 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(v.Type().Field(i).Name)
		sb.WriteString(": ")
		sb.WriteString(fmt.Sprintf("%v", v.Field(i)))
	}
	sb.WriteString("}")

	return sb.String()
}

// StructToStringWithoutNils converts a struct to a string without nil values.
// Note that the interface fields MUST all be pointers.
func StructToStringWithoutNils(s interface{}) string {
	sb := strings.Builder{}
	sb.WriteString("{")
	v := reflect.ValueOf(s)
	hasSmth := false
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).IsNil() {
			continue
		}
		if hasSmth {
			sb.WriteString(", ")
		}
		hasSmth = true

		sb.WriteString(v.Type().Field(i).Name)
		sb.WriteString(": ")
		sb.WriteString(fmt.Sprintf("%v", v.Field(i).Elem()))
	}
	sb.WriteString("}")

	return sb.String()
}

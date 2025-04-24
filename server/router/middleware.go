package router

import (
	"encoding/base64"
	"fmt"
	"server/config"
	"server/database"
	"strings"

	"github.com/gin-gonic/gin"
)

// Authenticate Judge with Bearer token
func AuthenticateJudge() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			no("Authorization header required with Bearer token", ctx)
			return
		}

		// Make sure the auth header starts with "Bearer "
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			no("Invalid Authorization header, illegal format detected", ctx)
			return
		}

		// Extract the token
		token := authHeader[7:]

		// Make sure the token is valid (check for judge in database)
		state := GetState(ctx)
		judge, err := database.FindJudgeByToken(state.Db, token)
		if err != nil {
			ctx.AbortWithStatusJSON(500, gin.H{"error": fmt.Sprintf("Error finding judge in database: %s", err.Error())})
			return
		}
		if judge == nil {
			no("Invalid Authorization header, no judge with provided token", ctx)
			return
		}

		// Success! - set judge in context
		ctx.Set("judge", judge)
		ctx.Next()
	}
}

// Authenticate Admin with Basic auth
func AuthenticateAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			no("Authorization header required with Basic auth", ctx)
			return
		}

		// Make sure the auth header starts with "Basic "
		if len(authHeader) < 6 || authHeader[:6] != "Basic " {
			no("Invalid Authorization header, illegal format detected", ctx)
			return
		}

		// Extract the base64 encoded username:password and decode it
		userpassHash := authHeader[6:]
		userpass, err := base64.StdEncoding.DecodeString(userpassHash)
		if err != nil {
			no("Invalid Authorization header, invalid base64 encoding", ctx)
			return
		}

		// Split the username:password string
		authSplit := strings.Split(string(userpass), ":")
		if len(authSplit) != 2 {
			no("Invalid Authorization header, illegal format detected", ctx)
			return
		}

		// Username should be "admin" and password should be the admin password
		if authSplit[0] != "admin" || authSplit[1] != config.GetEnv("JURY_ADMIN_PASSWORD") {
			no("Invalid Authorization header, incorrect credentials", ctx)
			return
		}
	}
}

// When auth invalid, send a 401 error
func no(msg string, ctx *gin.Context) {
	ctx.AbortWithStatusJSON(401, gin.H{"error": msg})
}

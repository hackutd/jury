package main

import (
	"context"
	"tests/src"
)

func main() {
	// Initialize the logger (make sure to change the .gitignore if this filename ever changes)
	logger := src.NewLogger("test-log.txt")

	// Log start message with date and time
	logger.LogLn(src.Info, "\n===============\nTESTING STARTED\n===============")
	logger.Log(src.Info, "Date/time: %s\n", src.GetDateTime())

	// Initialize the database connection
	db := src.InitDb(logger)

	// Close the database connection
	defer db.Client().Disconnect(context.Background())

	// Wait for backend to load
	src.WaitForBackend(logger)

	// ### TIME TO START TESTING!!! ###
	apiEndpointTests(logger)
}

// apiEndpointTests tests the functionality of each API endpoint
func apiEndpointTests(logger *src.Logger) {
	logger.LogLn(src.Info, "\nAPI ENDPOINT TESTS\n------------------")

	// Heartbeat
	res := src.GetRequest(logger, "/", src.DefaultAuth())
	if !src.IsOk(res) {
		logger.LogLn(src.Error, "Error with heartbeat endpoint")
		return
	}

	// Invalid login to admin account
	// TODO: Move this to login stress tests
	res = src.PostRequest(logger, "/admin/login", src.H{"password": "THIS IS DEFINITELY THE WRONG PASSWORD"}, src.DefaultAuth())
	if src.IsOk(res) {
		logger.LogLn(src.Error, "Invalid login to admin account should not be successful")
		return
	}

	// Login to admin account
	password := src.GetEnv("ADMIN_PASSWORD")
	res = src.PostRequest(logger, "/admin/login", src.H{"password": password}, src.DefaultAuth())
	if !src.IsOk(res) {
		logger.LogLn(src.Error, "Error logging in as admin")
		return
	}

	// Check if authenticated with invalid token
	res = src.PostRequest(logger, "/admin/auth", nil, "Bearer INVALID_TOKEN")
	if src.IsOk(res) {
		logger.LogLn(src.Error, "Invalid token should not be authenticated")
		return
	}

	// Check if authenticated
	res = src.PostRequest(logger, "/admin/auth", nil, src.AdminAuth())
	if !src.IsOk(res) {
		logger.LogLn(src.Error, "Error checking if admin is authenticated")
		return
	}

	// Add a judge
	res = src.PostRequest(logger, "/admin/judge/new", src.H{"name": "Test Judge", "email": "test@gmail.com", "track": "", "notes": "test", "no_send": true}, src.AdminAuth())
	if !src.IsOk(res) {
		logger.LogLn(src.Error, "Error adding a judge")
		return
	}

	// TODO: Add all routes eventually

	// Get the clock
	res = src.GetRequest(logger, "/admin/clock", src.AdminAuth())
	if !src.IsValue(res, "running", src.BoolType, false) || !src.IsValue(res, "time", src.Float64Type, 0.0) {
		logger.LogLn(src.Error, "Error getting the clock")
		return
	}

	// Success!
	logger.LogLn(src.Info, "\tAll API ENDPOINT TESTS passed!")
}

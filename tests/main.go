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
	loginTests(logger)
}

func loginTests(logger *src.Logger) {
	logger.LogLn(src.Info, "\nLOGIN TESTS\n-----------")

	// Invalid login to admin account
	res := src.PostRequest(logger, "/admin/login", src.H{"password": "THIS IS DEFINITELY THE WRONG PASSWORD"}, src.DefaultAuth())
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

	// Check if authenticated
	res = src.PostRequest(logger, "/admin/auth", nil, src.AdminAuth())
	if !src.IsOk(res) {
		logger.LogLn(src.Error, "Error checking if admin is authenticated")
		return
	}

	// TODO: Potentially create a parser for this so I don't have to write all of this by hand
}

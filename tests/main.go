package main

import (
	"context"
	"tests/src"
	"tests/tests"
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

	// Create a context with the database and logger
	context := &src.Context{
		Db:     db,
		Logger: logger,
	}

	// Run all tests!
	tests.RunTests(context)
}

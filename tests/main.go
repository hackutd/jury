package main

import (
	"context"
	"os"
	"tests/tests"
	"tests/util"
)

func main() {
	// Initialize the logger (make sure to change the .gitignore if this filename ever changes)
	logger := util.NewLogger("test-log.txt")

	// Log start message with date and time
	logger.LogLn(util.Info, "\n===============\nTESTING STARTED\n===============")
	logger.Log(util.Info, "Date/time: %s\n", util.GetDateTime())

	// Initialize the database connection
	db := util.InitDb(logger)

	// Close the database connection
	defer db.Client().Disconnect(context.Background())

	// Wait for backend to load
	err := util.WaitForBackend(logger)
	if err != nil {
		os.Exit(1)
	}

	// Create a context with the database and logger
	context := &util.Context{
		Db:     db,
		Logger: logger,
	}

	// Run all tests!
	tests.RunTests(context)
}

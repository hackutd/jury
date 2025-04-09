package main

import (
	"errors"
	"fmt"
	"server/config"
	"server/database"
	"server/logging"
	"server/router"

	"github.com/joho/godotenv"
)

func main() {
	// Load the env file
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Did not load .env file (%s). This is expected when running in a Docker container\n", err.Error())
	}

	// Check for all necessary env variables)
	config.CheckEnv()

	// Connect to the database
	db := database.InitDb()

	// Load the logger
	logger, err := logging.NewLogger(db)
	if err != nil {
		panic(errors.New("failed to load logger: " + err.Error()))
	}
	logger.SystemLogf("Server restarted")

	// Create the router and attach variables
	r := router.NewRouter(db, logger)

	// Start the Gin server!
	r.Run(":" + config.GetOptEnv("PORT", "8080"))
}

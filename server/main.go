package main

import (
	"fmt"
	"server/config"
	"server/database"
	"server/logging"
	"server/router"

	"github.com/joho/godotenv"
)

func main() {
	// Load the logger
	logger := logging.NewLogger()
	logger.SystemLogf("Server started")

	// Load the env file
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Did not load .env file (%s). This is expected when running in a Docker container\n", err.Error())
	}

	// Check for all necessary env variables)
	config.CheckEnv()
	logger.SystemLogf("Successfully loaded environment variables")

	// Connect to the database
	db := database.InitDb()
	logger.SystemLogf("Connected to the database")

	// Create the router and attach variables
	r := router.NewRouter(db, logger)

	// Start the Gin server!
	r.Run(":" + config.GetOptEnv("PORT", "8080"))
}

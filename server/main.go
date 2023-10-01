package main

import (
	"log"
	"server/config"
	"server/database"
	"server/router"

	"github.com/joho/godotenv"
)

func main() {
	// Load the env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s\n", err.Error())
	}

	// Check for all necessary env files
	config.CheckEnv()

	// Connect to the database
	db := database.InitDb()

	// Create the router and attach variables
	r := router.NewRouter(db)

	// Start the Gin server!
	r.Run(":" + config.GetOptEnv("PORT", "8080"))
}

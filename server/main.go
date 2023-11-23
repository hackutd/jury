package main

import (
	"fmt"
	"server/config"
	"server/database"
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

	// Create the router and attach variables
	r := router.NewRouter(db)

	// Start the Gin server!
	r.Run(":" + config.GetOptEnv("PORT", "8080"))
}

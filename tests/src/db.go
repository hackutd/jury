package src

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InitDb initializes the database connection to MongoDB.
// This will proactively panic if any step of the connection protocol breaks
func InitDb(logger *Logger) *mongo.Database {
	logger.LogLn(Info, "Connecting to database...")

	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(GetEnv("MONGODB_URI")).SetServerAPIOptions(serverAPI).SetTimeout(30 * time.Second)

	// Connect to the database
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		logger.Log(Error, "Error connecting to database: %s\n", err.Error())
		os.Exit(1)
	}

	// Send ping to confirm successful connections
	var result bson.M
	if err := client.Database("admin").RunCommand(context.Background(), bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		logger.Log(Error, "Error pinging database: %s\n", err.Error())
		os.Exit(1)
	}
	logger.LogLn(Info, "Successfully connected to database!")

	// Check if replica set is initialized
	for {
		// Check if the replica set is initialized
		var result bson.M
		if err := client.Database("admin").RunCommand(context.Background(), bson.D{{Key: "replSetGetStatus", Value: 1}}).Decode(&result); err != nil {
			// If not initialized, log the error and wait 2 seconds before trying again
			logger.Log(Info, "Replica set not initialized, retrying in 2 seconds: %s\n", err.Error())
			time.Sleep(2 * time.Second)
			continue
		}
		break
	}
	logger.LogLn(Info, "Replica set is initialized!")

	// Return the "jury" database
	return client.Database("jury")
}

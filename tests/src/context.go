package src

import "go.mongodb.org/mongo-driver/mongo"

type Context struct {
	Db     *mongo.Database
	Logger *Logger
}

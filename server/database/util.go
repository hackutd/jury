package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

func WithTransaction(db *mongo.Database, fn func(mongo.SessionContext) (interface{}, error)) error {
	wc := writeconcern.Majority()
	txnOptions := options.Transaction().SetWriteConcern(wc)

	session, err := db.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(context.Background())

	_, err = session.WithTransaction(context.Background(), fn, txnOptions)
	return err
}

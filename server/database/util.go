package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

// WithTransaction runs a function with a transaction,
// returning only an error if the transaction fails.
func WithTransaction(db *mongo.Database, fn func(mongo.SessionContext) error) error {
	wc := writeconcern.Majority()
	txnOptions := options.Transaction().SetWriteConcern(wc)

	session, err := db.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(context.Background())

	wrappedFn := func(sc mongo.SessionContext) (interface{}, error) {
		return nil, fn(sc)
	}

	_, err = session.WithTransaction(context.Background(), wrappedFn, txnOptions)
	return err
}

// WithTransactionItem runs a function with a transaction, returning the result of the function
// and an error if the transaction fails.
func WithTransactionItem(db *mongo.Database, fn func(mongo.SessionContext) (interface{}, error)) (interface{}, error) {
	wc := writeconcern.Majority()
	txnOptions := options.Transaction().SetWriteConcern(wc)

	session, err := db.Client().StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(context.Background())

	return session.WithTransaction(context.Background(), fn, txnOptions)
}

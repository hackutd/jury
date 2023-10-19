package database

import (
	"context"
	"server/models"

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

func GetNextTableNum(o *models.Options) int64 {
	// Increment table number
	o.CurrTableNum++

	// If not using groups, just return the next table num
	if !o.UseGroups {
		return o.CurrTableNum
	}

	// If using groups, find the group that the table num belongs in
	for _, group := range o.Groups {
		// Exactly in a group
		if o.CurrTableNum >= group.Start && o.CurrTableNum <= group.End {
			return o.CurrTableNum
		}

		// Otherwise, the first group where the table number is less than the start of the group
		if o.CurrTableNum < group.Start {
			o.CurrTableNum = group.Start
			return group.Start
		}
	}

	// If greater than the last group, just use increment the table number
	return o.CurrTableNum
}

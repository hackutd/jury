---
sidebar_position: 3
title: Database
description: How Jury interacts and stores data in MongoDB
---

# MongoDB Database

Most of the database calling functionality should be located in the `database` module. Because MongoDB is a NoSQL database, a lot of the query handling is handled by the application and not the database. We use the MongoDB driver library to perform all calls with the database.

The connection to the database is done in `server/database/init.go`, using the provided [environmental variables](/docs/reference/envs) to connect to the database.

## Retrieving from Database

To retrieve an item from the database, you can follow the example from `FindJudgeByToken`:

```go
func FindJudgeByToken(db *mongo.Database, token string) (*models.Judge, error) {
	var judge models.Judge
	err := db.Collection("judges").FindOne(context.Background(), gin.H{"token": token}).Decode(&judge)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &judge, err
}
```

Note that we are using `Decode` to convert the resultant object into a defined model (`Judge` in this case). This method of `FindOne` is used across the app.

To retrieve multiple items from the database, we can use a cursor object to get all items:

```go
func FindAllJudges(db *mongo.Database, ctx context.Context) ([]*models.Judge, error) {
	judges := make([]*models.Judge, 0)
	cursor, err := db.Collection("judges").Find(ctx, gin.H{})
	if err != nil {
		return nil, err
	}
	for cursor.Next(ctx) {
		var judge models.Judge
		err := cursor.Decode(&judge)
		if err != nil {
			return nil, err
		}
		judges = append(judges, &judge)
	}
	return judges, nil
}
```

## Inserting into/Updating the Database

To update an item in the database, we can use the following format:

```go
func UpdateJudgeReadWelcome(db *mongo.Database, ctx context.Context, judgeId *primitive.ObjectID) error {
	_, err := db.Collection("judges").UpdateOne(ctx, gin.H{"_id": judgeId}, gin.H{"$set": gin.H{"read_welcome": true, "last_activity": util.Now()}})
	return err
}
```

Note that we will return only the `err` object in this case, ignoring the result object from the db update.

## Multiple Instances of Jury

Jury is not designed to run with multiple instances of the app connecting to the same database. This will cause many problems such as the following:

- Comparisons matrix will not be correct due to it being maintained on the server
- Judging Clock will not be synced across instances

## Transactions

A lot of the database functions run multiple database operations in a row. We want these operations to all happen at the same time otherwise there may be race conditions and inconsistencies among the DB. The solution to this is to use [transactions](https://www.mongodb.com/docs/manual/core/transactions/). We have written a wrapper function for this, so you only need to call the wrapper and make sure all database functions pass in the context with a `ctx context.Context` parameter:

```go
func PutJudgesInGroups(db *mongo.Database) error {
	return WithTransaction(db, func(sc mongo.SessionContext) error {
		// Get all judges
		judges, err := FindJudgesByTrack(db, sc, "")
		if err != nil {
			return err
		}

		// Get the next n groups to assign judges to
		groups, err := GetNextNJudgeGroups(db, sc, len(judges), true)
		if err != nil {
			return err
		}

		// Assign each judge to a group
		for i, judge := range judges {
			judge.Group = groups[i]
			judge.GroupSeen = 0
		}

		// Update the judges in the database
		err = UpdateJudgeGroups(db, sc, judges)
		return err
	})
}
```

Notice how here we are getting the judges, figuring out which groups to assign judges to, and updating all judges in one transaction. If any of these fail, the entire transaction should fail. Also, no DB operations are done between these operations.

## Exporting the Database

To do a full database export you can use the [mongoexport](https://www.mongodb.com/docs/database-tools/mongoexport/) tool. There is unfortunately no good way to do it through the driver and thus that functionality is not avaliable in Jury.

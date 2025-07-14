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

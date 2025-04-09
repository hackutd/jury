package logging

import (
	"context"
	"fmt"
	"server/database"
	"server/models"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Logger struct {
	Mutex  sync.Mutex
	Memory []string
	DbRef  *mongo.Database
}

type LogType int

const (
	System LogType = iota
	Admin
	Judge
)

const DbLogLimit = 10000

func NewLogger(db *mongo.Database) (*Logger, error) {
	// Get all logs from the db
	dbLogs, err := GetAllDbLogs(db)
	if err != nil {
		return nil, err
	}

	// Create a new logger and add all logs from the database
	memory := []string{}
	for _, log := range dbLogs {
		memory = append(memory, log.Entries...)
	}

	return &Logger{
		Mutex:  sync.Mutex{},
		Memory: memory,
		DbRef:  db,
	}, nil
}

// GetAllDbLogs retrieves all logs from the database.
func GetAllDbLogs(db *mongo.Database) ([]*models.Log, error) {
	// Get all logs from the database, sorted in ascending order of time
	cursor, err := db.Collection("logs").Find(context.Background(), gin.H{}, options.Find().SetSort(gin.H{"time": 1}))
	if err != nil {
		return nil, err
	}

	// Iterate through the cursor and decode each log
	logs := []*models.Log{}
	for cursor.Next(context.Background()) {
		var log models.Log
		err := cursor.Decode(&log)
		if err != nil {
			return nil, err
		}

		logs = append(logs, &log)
	}

	// If no logs, create one
	if len(logs) == 0 {
		_, err = db.Collection("logs").InsertOne(context.Background(), models.NewLog())
		if err != nil {
			return nil, err
		}
	}

	return logs, nil
}

// TODO: Pass in context as a parameter to the log functions, then remove the transaction wrapper from here
func (l *Logger) writeToDb(item string) error {
	return database.WithTransaction(l.DbRef, func(sc mongo.SessionContext) error {
		// Insert the log into the database
		res := l.DbRef.Collection("logs").FindOneAndUpdate(
			sc,
			gin.H{},
			gin.H{"$push": gin.H{"entries": item}, "$inc": gin.H{"count": 1}},
			options.FindOneAndUpdate().SetSort(gin.H{"time": -1}),
		)

		// Check for errors
		if res.Err() != nil {
			return res.Err()
		}

		// Get the log from the result
		var log models.Log
		err := res.Decode(&log)
		if err != nil {
			return err
		}

		// If the log is too long, create a new log in the db
		if log.Count > DbLogLimit {
			_, err = l.DbRef.Collection("logs").InsertOne(sc, models.NewLog())
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// Logf logs a message to the log file.
// The time is prepended to the message and a newline is appended.
func (l *Logger) Logf(t LogType, message string, args ...interface{}) error {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	// Form output string
	userInput := fmt.Sprintf(message, args...)
	output := fmt.Sprintf("[%s] %s | %s", time.Now().Format(time.RFC3339), typeToString(t), userInput)

	// Write to memory
	l.Memory = append(l.Memory, output)

	// Write to database
	l.writeToDb(output)

	return nil
}

// SystemLogf logs a message with the System LogType.
func (l *Logger) SystemLogf(message string, args ...interface{}) error {
	return l.Logf(System, message, args...)
}

// JudgeLogf logs a message with the Judge LogType.
// This will include a judge and project (if defined) in the log message.
func (l *Logger) JudgeLogf(judge *models.Judge, message string, args ...interface{}) error {
	judgeDisplay := ""
	if judge != nil {
		judgeDisplay = fmt.Sprintf("Judge %s (%s) | ", judge.Name, judge.Id.Hex())
	}

	judgeMessage := fmt.Sprintf("%s%s", judgeDisplay, message)
	return l.Logf(Judge, judgeMessage, args...)
}

// AdminLogf logs a message with the Admin LogType.
func (l *Logger) AdminLogf(message string, args ...interface{}) error {
	return l.Logf(Admin, message, args...)
}

// typeToString converts a LogType to a string.
func typeToString(t LogType) string {
	switch t {
	case System:
		return "SYSTEM"
	case Admin:
		return "ADMIN "
	case Judge:
		return "JUDGE "
	default:
		return "UNDEF "
	}
}

// Get returns the log file contents as a string.
func (l *Logger) Get() string {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	return strings.Join(l.Memory, "\n")
}

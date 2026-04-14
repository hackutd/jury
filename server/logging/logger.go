package logging

import (
	"context"
	"fmt"
	"server/models"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Logger struct {
	Mutex  sync.Mutex
	Memory []string
}

type LogType int

const (
	System LogType = iota
	Admin
	Judge
)

func NewLogger(db *mongo.Database) (*Logger, error) {
	// Load historical log entries from the database so the admin log
	// view still shows entries from before the current server session.
	dbLogs, err := GetAllDbLogs(db)
	if err != nil {
		return nil, err
	}

	memory := []string{}
	for _, log := range dbLogs {
		memory = append(memory, log.Entries...)
	}

	return &Logger{
		Mutex:  sync.Mutex{},
		Memory: memory,
	}, nil
}

// GetAllDbLogs retrieves all logs from the database.
func GetAllDbLogs(db *mongo.Database) ([]*models.Log, error) {
	cursor, err := db.Collection("logs").Find(context.Background(), map[string]interface{}{}, options.Find().SetSort(map[string]interface{}{"time": 1}))
	if err != nil {
		return nil, err
	}

	logs := []*models.Log{}
	for cursor.Next(context.Background()) {
		var log models.Log
		err := cursor.Decode(&log)
		if err != nil {
			return nil, err
		}
		logs = append(logs, &log)
	}

	return logs, nil
}

// Logf logs a message to the in-memory log.
// The time is prepended to the message and a newline is appended.
// Note: entries are no longer written to MongoDB on every call — doing so
// required a serialised majority-write transaction per API request and was
// the primary source of lock contention under concurrent judge load.
func (l *Logger) Logf(t LogType, message string, args ...interface{}) error {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	userInput := fmt.Sprintf(message, args...)
	output := fmt.Sprintf("[%s] %s | %s", time.Now().Format(time.RFC3339), typeToString(t), userInput)

	l.Memory = append(l.Memory, output)

	return nil
}

// SystemLogf logs a message with the System LogType.
func (l *Logger) SystemLogf(message string, args ...interface{}) error {
	return l.Logf(System, message, args...)
}

// JudgeLogf logs a message with the Judge LogType.
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

package logging

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"server/models"
	"strings"
	"sync"
	"time"
)

type Logger struct {
	Mutex      sync.Mutex
	File       *os.File
	FileWriter *bufio.Writer
	Memory     *strings.Builder
}

type LogType int

const (
	System LogType = iota
	Admin
	Judge
)

func NewLogger() *Logger {
	// Open the log file
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(errors.New("error opening log file: " + err.Error()))
	}

	// Read the log file into memory
	memory := &strings.Builder{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		memory.WriteString(scanner.Text() + "\n")
	}

	return &Logger{
		Mutex:      sync.Mutex{},
		File:       f,
		FileWriter: bufio.NewWriter(f),
		Memory:     &strings.Builder{},
	}
}

// Logf logs a message to the log file.
// The time is prepended to the message and a newline is appended.
func (l *Logger) Logf(t LogType, message string, args ...interface{}) error {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	// Form output string
	userInput := fmt.Sprintf(message, args...)
	output := fmt.Sprintf("[%s] %s | %s\n", time.Now().Format(time.RFC3339), typeToString(t), userInput)

	// Write to file
	_, err := l.FileWriter.WriteString(output)
	if err != nil {
		return err
	}

	// Flush file
	err = l.FileWriter.Flush()
	if err != nil {
		return err
	}

	// Write to memory
	_, err = l.Memory.WriteString(output)
	return err
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

	return l.Memory.String()
}

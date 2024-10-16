package src

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type LogLevel int

const (
	Error LogLevel = iota
	Warn
	Info
)

type Logger struct {
	OutFile *os.File
}

// NewLogger will create a new logger with the defined output sources
func NewLogger(outFile string) *Logger {
	// Open the file for writing
	file, err := os.OpenFile(outFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening log file: %s\n", err.Error())
	}

	return &Logger{OutFile: file}
}

// Log will log the message to the defined input sources
func (l *Logger) Log(level LogLevel, format string, a ...any) {
	if logLevelToInt(GetEnv("LOG_LEVEL")) < int(level) {
		return
	}

	message := fmt.Sprintf(format, a...)

	l.OutFile.WriteString(message)
	fmt.Println(message)
}

// LogLn will log the message to the defined input sources with a newline character
func (l *Logger) LogLn(level LogLevel, message string) {
	l.Log(level, "%s\n", message)
}

func logLevelToInt(level string) int {
	switch strings.ToLower(level) {
	case "error":
		return int(Error)
	case "warn":
		return int(Warn)
	case "info":
		return int(Info)
	default:
		log.Println("Invalid log level, defaulting to INFO")
		return int(Info)
	}
}

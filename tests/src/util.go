package src

import (
	"log"
	"os"
	"time"
)

// GetEnv returns the value of the environmental variable or panics if it does not exist
func GetEnv(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("ERROR: %s environmental variable not defined\n", key)
		return ""
	}
	return val
}

// GetOptEnv returns the value of the environmental variable or the default value if it does not exist
func GetOptEnv(key string, defaultVal string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultVal
	}
	return val
}

// GetDateTime returns a formatted datetime string
func GetDateTime() string {
	now := time.Now()
	return now.Format("January 2, 2006 @ 03:04 PM MST")
}

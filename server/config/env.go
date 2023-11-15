package config

import (
	"log"
	"os"
	"server/util"
)

var requiredEnvs = [...]string{"JURY_ADMIN_PASSWORD", "EMAIL_FROM"}
var smtpEnvs = []string{"EMAIL_HOST", "EMAIL_USERNAME", "EMAIL_PASSWORD"}
var sendgridEnvs = []string{"SENDGRID_API_KEY", "EMAIL_FROM_NAME"}

// Checks to see if all required environmental variables are defined
func CheckEnv() {
	for _, v := range requiredEnvs {
		if !hasEnv(v) {
			log.Fatalf("ERROR: %s environmental variable not defined\n", v)
		}
	}

	// Check to see if either all smtp envs are defined or all sendgrid envs are defined
	if !util.All(util.Map(smtpEnvs, hasEnv)) && !util.All(util.Map(sendgridEnvs, hasEnv)) {
		log.Fatalf("ERROR: either all envs for smtp or sendgrid must be defined (one of these sets): %v OR %v\n", smtpEnvs, sendgridEnvs)
	}
}

// Returns true if the environmental variable is defined and not empty
func hasEnv(key string) bool {
	val, ok := os.LookupEnv(key)
	if !ok {
		return false
	}
	return val != ""
}

func GetEnv(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("ERROR: %s environmental variable not defined\n", key)
		return ""
	}
	return val
}

func GetOptEnv(key string, defaultVal string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultVal
	}
	return val
}

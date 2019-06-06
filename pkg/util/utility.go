package util

import (
	"log"
	"os"
)

// FailOnError logs a fatal error, with text msg, if err is not nil.
func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// GetEnvironmentVariable retrieves environment variable, key.
// Logs and exits if not set and found.
func GetEnvironmentVariable(key string) string {
	val, set := os.LookupEnv(key)
	if !set {
		log.Fatalf("%s environment variable is not set", key)
	}
	return val
}

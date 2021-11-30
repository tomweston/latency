package utils

import "os"

// GetEnv returns the value of the environment variable named by the key. If the key value is empty the default value is returned.
func GetEnv(envName, valueDefault string) string {
	value := os.Getenv(envName)
	if value == "" {
		return valueDefault
	}
	return value
}

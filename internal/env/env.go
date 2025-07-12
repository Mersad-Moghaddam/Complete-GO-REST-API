package env

import (
	"os"
	"strconv"
)

// GetEnvString retrieves the value of an environment variable named
// 'key'. If the variable doesn't exist, the function returns the given
// 'defaultValue'.
func GetEnvString(key string, defaultValue string) string {
	if value, exist := os.LookupEnv(key); exist {
		return value
	}
	return defaultValue
}

func GetEnvInt(key string, defaultValue int) int {
	if value, exist := os.LookupEnv(key); exist {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

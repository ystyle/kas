package env

import (
	"os"
	"strconv"
)

func GetInt(key string, v int) int {
	if max := os.Getenv(key); max != "" {
		if i, err := strconv.Atoi(max); err != nil {
			return i
		}
	}
	return v
}

func GetBool(key string, v bool) bool {
	if b := os.Getenv(key); b != "" {
		switch b {
		case "1", "true", "True":
			return true
		case "0", "false", "False":
			return false
		}
	}
	return v
}

func GetString(key string, v string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return v
}

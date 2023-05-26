package utils

import "os"

// GetSecretKey get the secret key from the environment
func GetSecretKey(key string) string {
	return os.Getenv(key)
}

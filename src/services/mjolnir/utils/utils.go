package utils

import (
	"os"

	"go.uber.org/zap"

	"beckendrof/gaia/src/services/mjolnir/types"
)

var (
	Host      string
	Port      string
	LogLevels = []string{"all", "debug", "info", "warn", "error", "fatal", "panic"}
	Config    types.MjolnirConfig
	Loggers   = make(map[string]*zap.Logger)
)

func GetServices() []string {
	var topics []string
	for _, service := range Config.Services {
		topics = append(topics, service.Name)
	}
	return topics
}

func GetLogPath(service string) string {
	return Config.Services[service].LogPath
}

func GetS3Path(service string) string {
	return Config.Services[service].S3Path
}

func GetEnvKey(key string) string {
	return os.Getenv(key)
}

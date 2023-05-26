package utils

import (
	"encoding/json"
	"strings"

	"go.uber.org/zap"
	"golang.org/x/exp/slices"

	"beckendrof/gaia/src/services/mjolnir/types"
)

func Parser(msg string) (*types.MjolnirMessage, error) {
	logData := types.MjolnirMessage{}
	if err := json.Unmarshal([]byte(msg), &logData); err != nil {
		return nil, err
	} else {
		return &logData, nil
	}
}

func LogToFile(logData *types.MjolnirMessage) bool {
	if srvLogger, ok := Loggers[logData.Service]; ok {
		lvl := strings.ToLower(logData.Level)
		if slices.Contains(LogLevels, lvl) {
			switch lvl {
			case "debug":
				srvLogger.Debug(logData.Message, zap.String("caller", logData.Caller), zap.String("stacktrace", logData.StackTrace))
			case "info":
				srvLogger.Info(logData.Message, zap.String("caller", logData.Caller), zap.String("stacktrace", logData.StackTrace))
			case "warn":
				srvLogger.Warn(logData.Message, zap.String("caller", logData.Caller), zap.String("stacktrace", logData.StackTrace))
			case "error":
				srvLogger.Error(logData.Message, zap.String("caller", logData.Caller), zap.String("stacktrace", logData.StackTrace))
			case "fatal":
				srvLogger.Fatal(logData.Message, zap.String("caller", logData.Caller), zap.String("stacktrace", logData.StackTrace))
			case "panic":
				srvLogger.Panic(logData.Message, zap.String("caller", logData.Caller), zap.String("stacktrace", logData.StackTrace))
			}
			return true
		}
	}
	return false
}

package zap

import (
	"log"
	"os"
	"path/filepath"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"beckendrof/gaia/src/services/mjolnir/utils"
)

func InitializeLogger() bool {
	for _, service := range utils.GetServices() {
		zapConfig := zap.NewProductionEncoderConfig()
		zapConfig.EncodeTime = zapcore.RFC3339TimeEncoder
		// change encoder to RFC1123
		consoleEncoder := zapcore.NewConsoleEncoder(zapConfig)
		var logCores []zapcore.Core
		for _, level := range utils.LogLevels {
			fileEncoder := zapcore.NewJSONEncoder(zapConfig)
			writer := zapcore.AddSync(rotatorInit(level, service))
			log_level, err := zapcore.ParseLevel(level)
			if err != nil {
				logCores = append(logCores, zapcore.NewCore(fileEncoder, writer, zapcore.DebugLevel))
			} else {
				levelEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
					return lvl == log_level
				})
				logCores = append(logCores, zapcore.NewCore(fileEncoder, writer, levelEnabler))
			}
		}
		if utils.Config.StdLogsEnable {
			logCores = append(logCores, zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel))
		}
		core := zapcore.NewTee(logCores...)
		zapLevelLogger := zap.New(core)
		defer zapLevelLogger.Sync()
		utils.Loggers[service] = zapLevelLogger
	}
	return true
}

func rotatorInit(rtype string, topic string) *rotatelogs.RotateLogs {
	rotator, err := rotatelogs.New(
		filepath.Join(utils.GetLogPath(topic), rtype+"-%Y-%m-%d.log"),
		rotatelogs.WithMaxAge(15*24*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(24*time.Hour)))
	if err != nil {
		log.Printf("Rotator Init for %v logs failed for topic %v: %v", rtype, topic, err)
	}
	return rotator
}

// TODO: S3 push

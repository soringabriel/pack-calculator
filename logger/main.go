package logger

import (
	"os"
	"strconv"

	"github.com/google/uuid"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

var Instance *logrus.Entry

type LoggerSettings struct {
	LogFile       string
	LogMaxSizeMB  int
	LogMaxBackups int
	LogMaxAge     int
}

// Setup a logger instance
func SetupLogger() {
	if Instance == nil {
		logrus.SetFormatter(&logrus.JSONFormatter{})
		// Setup logger
		appName := os.Getenv("APP_NAME")
		Instance = logrus.WithFields(logrus.Fields{
			"appname": appName,
			"session": uuid.New().String(),
		})

		// Logging to file
		logFile := os.Getenv("LOG_FILE")
		if len(logFile) > 0 {
			logMaxSizeMB, err := strconv.Atoi(os.Getenv("LOG_MAX_SIZE_MB"))
			if err != nil {
				Instance.Fatal("Failed to convert LOG_MAX_SIZE_MB env variable to int", err)
			}
			logMaxBackups, err := strconv.Atoi(os.Getenv("LOG_MAX_BACKUPS"))
			if err != nil {
				Instance.Fatal("Failed to convert LOG_MAX_BACKUPS env variable to int", err)
			}
			logMaxAge, err := strconv.Atoi(os.Getenv("LOG_MAX_AGE"))
			if err != nil {
				Instance.Fatal("Failed to convert LOG_MAX_AGE env variable to int", err)
			}
			Instance.Logger.SetOutput(&lumberjack.Logger{
				Filename:   logFile,
				MaxSize:    logMaxSizeMB,
				MaxBackups: logMaxBackups,
				MaxAge:     logMaxAge,
			})
		}
	}
}

package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

func NewLogger() (*logrus.Entry, error) {
	var (
		level logrus.Level
		err   error
	)

	level = logrus.DebugLevel

	if logLevelEnv := os.Getenv("log_level"); logLevelEnv != "" {
		level, err = logrus.ParseLevel(logLevelEnv)
		if err != nil {
			return nil, err
		}

	}

	logger := logrus.New()
	logger.SetLevel(level)
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	entry := logrus.NewEntry(logger)

	return entry, nil

}

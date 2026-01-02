package debug

import (
	"log"
	"os"
)

const LogsDir = "./logs"

func NewLogger() (*log.Logger, error) {
	_, err := os.Stat(LogsDir)
	if err != nil {
		if err := os.Mkdir(LogsDir, 0744); err != nil {
			return &log.Logger{}, err
		}
	}

	file, err := os.OpenFile("logs/dev.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return &log.Logger{}, err
	}

	logger := log.Default()
	logger.SetOutput(file)

	return logger, nil
}

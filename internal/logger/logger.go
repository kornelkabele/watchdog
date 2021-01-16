package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

// Logger object
type Logger struct {
	file *rotatelogs.RotateLogs
}

// New logger
func New(logFile string) (*Logger, error) {
	file, err := rotatelogs.New(
		logFile+".%Y%m%d",
		rotatelogs.WithMaxAge(30*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create rotatelogs: %s", err)
	}

	logger := &Logger{file}

	mw := io.MultiWriter(os.Stdout, file)
	log.SetOutput(mw)
	log.Println("Started")
	log.Writer()

	return logger, nil
}

// Close closes logging
func (l *Logger) Close() error {
	log.Println("Stopped")
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

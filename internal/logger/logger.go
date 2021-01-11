package logger

import (
	"io"
	"log"
	"os"
)

// SetLogger sets logging
func SetLogger(logFile string) (*os.File, error) {
	file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	mw := io.MultiWriter(os.Stdout, file)
	log.SetOutput(mw)
	log.Println("Started")

	return file, nil
}

// StopLogger closes logging
func StopLogger(file *os.File) {
	log.Println("Stopped")
	file.Close()
}

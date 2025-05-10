package core

import (
	"io"
	"log/slog"
	"os"
)

func NewLogger(logsFile string) (*slog.Logger, error) {
	logFile, err := os.OpenFile(logsFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}

	multiWriter := io.MultiWriter(logFile, os.Stdout)
	l := slog.New(slog.NewJSONHandler(multiWriter, nil))

	return l, nil
}

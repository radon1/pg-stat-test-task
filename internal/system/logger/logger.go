package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
)

func New(logLevel string) (*zerolog.Logger, error) {
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		return nil, fmt.Errorf("failed to parse zerolog level: %w", err)
	}
	zeroLogger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: zerolog.TimeFormatUnixMs})
	zeroLogger.Level(level)

	return &zeroLogger, nil
}

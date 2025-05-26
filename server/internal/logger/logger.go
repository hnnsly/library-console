package logger

import (
	"os"
	"strings"
	"time"

	"github.com/hnnsly/library-console/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Setup(cfg *config.LogConfig) *zerolog.Logger {
	var l zerolog.Logger

	level, err := zerolog.ParseLevel(strings.ToLower(cfg.Level))
	if err != nil {
		log.Warn().Msgf("Invalid log level '%s', defaulting to 'info'", cfg.Level)
		level = zerolog.InfoLevel
	}

	if cfg.Format == "text" {
		l = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
			Level(level).
			With().
			Timestamp().
			Logger()
	} else {
		l = zerolog.New(os.Stderr).
			Level(level).
			With().
			Timestamp().
			Logger()
	}
	return &l
}

package logger

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hnnsly/library-console/internal/config"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Setup(cfg *config.LogConfig) *zerolog.Logger {
	zerolog.ParseLevel(cfg.Level)

	zerolog.TimeFieldFormat = time.DateTime

	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		// Получаем только последние два сегмента пути
		parts := strings.Split(file, "/")
		if len(parts) > 2 {
			file = strings.Join(parts[len(parts)-2:], "/")
		}
		return file + ":" + strconv.Itoa(line)
	}

	logFile, err := os.OpenFile(cfg.File, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to open log file")
	}

	loggerContext := zerolog.New(logFile).
		With().
		Timestamp().
		Logger().
		With().
		Caller().
		Logger()

	log.Info().Msg("logger setup complete!")

	return &loggerContext
}

func RequestLogger() fiber.Handler {
	logger := log.Logger
	return func(c *fiber.Ctx) error {
		logger.Info().
			Str("method", c.Method()).
			Str("url", c.OriginalURL()).
			Msg("incoming request")

		start := time.Now()
		defer func() {
			if time.Since(start) > time.Second*2 {
				logger.Warn().
					Str("method", c.Method()).
					Str("url", c.OriginalURL()).
					Dur("elapsed_ms", time.Since(start)).
					Msg("long response time")
			}
		}()

		return c.Next()
	}
}

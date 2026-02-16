package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Setup(logLevel string) {
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}

	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(level)

	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "15:04:05",
	})
}

func Get() zerolog.Logger {
	return log.With().Caller().Logger()
}

func GetWithContext(fields map[string]interface{}) zerolog.Logger {
	logger := Get()
	for key, value := range fields {
		logger = logger.With().Interface(key, value).Logger()
	}
	return logger
}

package logger

import (
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"
)

type Logger struct {
	logger zerolog.Logger
}

func NewLogger(level int) *Logger {
	zerolog.TimeFieldFormat = time.RFC3339

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	logger := zerolog.New(output).With().Timestamp().Logger()

	switch level {
	case 5:
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case 4:
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case 3:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case 2:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case 1:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case 0:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case -1:
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	return &Logger{
		logger: logger,
	}
}

func (l *Logger) Debug(msg string, fields ...interface{}) {
	l.logger.Debug().Fields(fieldsMap(fields...)).Msg(msg)
}

func (l *Logger) Info(msg string, fields ...interface{}) {
	l.logger.Info().Fields(fieldsMap(fields...)).Msg(msg)
}

func (l *Logger) Warn(msg string, fields ...interface{}) {
	l.logger.Warn().Fields(fieldsMap(fields...)).Msg(msg)
}

func (l *Logger) Error(msg string, err error, fields ...interface{}) {
	event := l.logger.Error().Err(err).Fields(fieldsMap(fields...))
	event.Msg(msg)
}

func (l *Logger) Fatal(msg string, err error, fields ...interface{}) {
	event := l.logger.Fatal().Err(err).Fields(fieldsMap(fields...))
	event.Msg(msg)
}

// Helper untuk convert slice to map[string]interface{}
func fieldsMap(fields ...interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	for i := 0; i < len(fields)-1; i += 2 {
		key, ok := fields[i].(string)
		if !ok {
			key = strconv.Itoa(i)
		}
		m[key] = fields[i+1]
	}
	return m
}

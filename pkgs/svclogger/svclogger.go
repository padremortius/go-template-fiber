package svclogger

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
)

// Logger -.
type Log struct {
	Logger *zerolog.Logger
	Level  string `env-required:"true" yaml:"level" json:"level" env:"LOG_LEVEL"`
}

// New -.
func New(level string) *Log {
	zerolog.SetGlobalLevel(GetLevelByString(level))
	zerolog.TimestampFieldName = "timestamp"

	locLog := zerolog.New(os.Stdout).With().Timestamp().Logger()

	return &Log{
		Logger: &locLog,
	}
}

func GetLevelByString(level string) zerolog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zerolog.DebugLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "info":
		return zerolog.InfoLevel
	case "trace":
		return zerolog.TraceLevel
	case "warn":
		return zerolog.WarnLevel
	default:
		return zerolog.InfoLevel
	}
}

func (l *Log) ChangeLogLevel(lvl string) {
	l.Logger.Info().Msgf("Change log level to %v", lvl)
	zerolog.SetGlobalLevel(GetLevelByString(lvl))
}

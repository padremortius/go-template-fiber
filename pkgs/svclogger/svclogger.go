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

func (l *Log) Debugf(aMsg string, args ...any) {
	l.Logger.Debug().Msgf(aMsg, args...)
}

func (l *Log) Infof(aMsg string, args ...any) {
	l.Logger.Info().Msgf(aMsg, args...)
}

func (l *Log) Errorf(aMsg string, args ...any) {
	l.Logger.Error().Msgf(aMsg, args...)
}

func (l *Log) Warnf(aMsg string, args ...any) {
	l.Logger.Warn().Msgf(aMsg, args...)
}

func (l *Log) Fatalf(aMsg string, args ...any) {
	l.Logger.Fatal().Msgf(aMsg, args...)
}

func (l *Log) Debug(aMsg string) {
	l.Logger.Debug().Msg(aMsg)
}

func (l *Log) Info(aMsg string) {
	l.Logger.Info().Msg(aMsg)
}

func (l *Log) Error(aMsg string) {
	l.Logger.Error().Msg(aMsg)
}

func (l *Log) Warn(aMsg string) {
	l.Logger.Warn().Msg(aMsg)
}

func (l *Log) Fatal(aMsg string) {
	l.Logger.Fatal().Msg(aMsg)
}

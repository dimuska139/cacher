package logging

import (
	"fmt"
	"github.com/dimuska139/cacher/pkg/config"
	"github.com/rs/zerolog"
	"time"
)

const (
	LogLevelDebug   = "debug"
	LogLevelInfo    = "info"
	LogLevelWarn    = "warn"
	LogLevelError   = "error"
	DefaultLogLevel = zerolog.DebugLevel
)

// Logger отвечает за логирование
type Logger struct {
	logger zerolog.Logger
}

// NewLogger создаёт логгер
func NewLogger(cfg *config.Config) *Logger {
	logger := zerolog.New(zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.TimeFormat = time.RFC3339
	})).
		With().
		Timestamp().
		Logger()

	if cfg != nil && cfg.Loglevel != "" {
		switch cfg.Loglevel {
		case LogLevelDebug:
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		case LogLevelInfo:
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		case LogLevelWarn:
			zerolog.SetGlobalLevel(zerolog.WarnLevel)
		case LogLevelError:
			zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		default:
			zerolog.SetGlobalLevel(DefaultLogLevel)
		}
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	return &Logger{logger}
}

func (z Logger) Debug(msg string, args ...interface{}) {
	z.logger.Debug().Timestamp().Fields(args).Msg(msg)
}

func (z Logger) Info(msg string, args ...interface{}) {
	z.logger.Info().Timestamp().Fields(args).Msg(msg)
}

func (z Logger) Warn(msg string, args ...interface{}) {
	z.logger.Warn().Timestamp().Fields(args).Msg(msg)
}

func (z Logger) Error(msg string, args ...interface{}) {
	z.logger.Error().Timestamp().Fields(args).Msg(msg)
}

func (z Logger) Panic(msg string, args ...interface{}) {
	z.logger.Panic().Timestamp().Fields(args).Msg(msg)
}

func (z Logger) Printf(format string, v ...interface{}) {
	z.logger.Printf(format, v...)
}

func (z Logger) Fatal(v ...interface{}) {
	z.logger.Fatal().Timestamp().Msg(fmt.Sprint(v...))
}

func (z Logger) Fatalf(format string, args ...interface{}) {
	z.logger.Fatal().Timestamp().Msgf(format, args...)
}

func (z Logger) Println(args ...interface{}) {
	z.logger.Info().Timestamp().Msgf("%v\r\n", args...)
}

func (z Logger) Print(args ...interface{}) {
	z.logger.Print(args...)
}

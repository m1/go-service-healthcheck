package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/m1/go-service-healthcheck/config"
)

// Default is the default logger
type Logger struct {
	*zap.Logger
}

func NewBasicLogger() *Logger {
	logger := &Logger{}
	cfg := zap.NewProductionConfig()
	l, _ := cfg.Build()
	logger.Logger = l
	return logger
}

// NewLogger Passes back a new logger
func NewLogger(config *config.Config) (*Logger, error) {
	var err error
	logger := &Logger{}
	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Sampling:    nil,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
	}

	l, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	logger.Logger = l.With([]zapcore.Field{
		zap.String("name", config.ServiceConfig.Name),
		zap.String("revision", config.ServiceConfig.GitCommit),
		zap.String("build_date", config.ServiceConfig.BuildDate),
		zap.String("api_port", config.APIConfig.Port),
		zap.String("api_domain", config.APIConfig.Domain),
		zap.String("env", config.Environment),
		zap.Bool("debug", config.Debug),
		zap.String("version", config.ServiceConfig.Version),
	}...)

	return logger, nil
}
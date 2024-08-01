package main

import (
	"log/slog"
	"strings"

	"github.com/mattn/go-colorable"
	slogzap "github.com/samber/slog-zap/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *slog.Logger

func parseLogLevel(s string) (slog.Level, error) {
	var level slog.Level
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	err := level.UnmarshalText([]byte(s))
	return level, err
}

func loggingSetup(level slog.Level) {
	encConf := zap.NewDevelopmentEncoderConfig()
	encConf.EncodeLevel = zapcore.CapitalColorLevelEncoder
	zapLogger := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(encConf),
		zapcore.AddSync(colorable.NewColorableStdout()),
		zapcore.DebugLevel,
	))
	logger := slog.New(slogzap.Option{Level: level, Logger: zapLogger}.NewZapHandler())
	globalLogger = logger
}

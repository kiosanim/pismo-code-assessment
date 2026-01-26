package logger

import (
	"context"
	"github.com/kiosanim/pismo-code-assessment/internal/core/config"
	"github.com/kiosanim/pismo-code-assessment/internal/core/contextkeys"
	"github.com/kiosanim/pismo-code-assessment/internal/core/contextutils"
	"github.com/kiosanim/pismo-code-assessment/internal/core/logger"
	"log/slog"
	"os"
	"strings"
)

type SlogLogger struct {
	l *slog.Logger
}

func NewSlogLogger(ctx context.Context, cfg *config.Configuration) *SlogLogger {
	var logLevel slog.Level
	logLevelFromConfig := cfg.App.LogLevel
	switch strings.ToLower(logLevelFromConfig) {
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	case "debug":
		logLevel = slog.LevelDebug
	default:
		logLevel = slog.LevelInfo
	}
	options := &slog.HandlerOptions{
		Level: logLevel,
	}
	traceID := contextutils.GetTraceID(ctx)
	sLogger := slog.New(slog.NewJSONHandler(os.Stdout, options))
	sLogger.With(
		contextkeys.TraceIDKey, traceID,
	)
	return &SlogLogger{l: sLogger}
}

func (s *SlogLogger) Info(msg string, args ...any) {
	s.l.Info(msg, args...)
}

func (s *SlogLogger) Warn(msg string, args ...any) {
	s.l.Warn(msg, args...)
}

func (s *SlogLogger) Debug(msg string, args ...any) {
	s.l.Debug(msg, args...)
}

func (s *SlogLogger) Error(msg string, args ...any) {
	s.l.Error(msg, args...)
}

func (s *SlogLogger) With(args ...any) logger.Logger {
	return &SlogLogger{l: s.l.With(args...)}
}

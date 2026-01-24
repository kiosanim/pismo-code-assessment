package logger

import (
	"github.com/kiosanim/pismo-code-assessment/internal/core/logger"
	"log/slog"
)

type SlogLogger struct {
	l *slog.Logger
}

func New(l *slog.Logger) *SlogLogger {
	return &SlogLogger{l: l}
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

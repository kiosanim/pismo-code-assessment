package mock

import (
	"github.com/kiosanim/pismo-code-assessment/internal/core/logger"
)

type MockLogger struct{}

func NewMockLogger() *MockLogger {
	return &MockLogger{}
}
func (m MockLogger) Info(msg string, args ...any)   {}
func (m MockLogger) Warn(msg string, args ...any)   {}
func (m MockLogger) Debug(msg string, args ...any)  {}
func (m MockLogger) Error(msg string, args ...any)  {}
func (m MockLogger) With(args ...any) logger.Logger { return MockLogger{} }

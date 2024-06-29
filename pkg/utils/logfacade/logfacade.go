package logfacade

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger интерфейс для логирования
type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
}

// LogFacade фасад для логирования
type LogFacade struct {
	logger Logger
}

// SetLogger устанавливает текущий логгер
func (l *LogFacade) SetLogger(logger Logger) {
	l.logger = logger
}

func (l *LogFacade) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *LogFacade) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *LogFacade) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

// Структуры для конкретных логгеров
type LogrusLogger struct {
	logger *logrus.Logger
}

func (ll *LogrusLogger) Info(args ...interface{}) {
	ll.logger.Info(args...)
}

func (ll *LogrusLogger) Error(args ...interface{}) {
	ll.logger.Error(args...)
}

func (ll *LogrusLogger) Fatal(args ...interface{}) {
	ll.logger.Fatal(args...)
}

// Logrus
func NewLogrusLogger() *LogrusLogger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)
	logger.SetOutput(os.Stdout)
	return &LogrusLogger{logger: logger}
}

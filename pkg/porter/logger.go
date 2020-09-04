package porter

import (
	"log"
	"os"
)

// Logger contains various instances of loggers for different log severity levels
type Logger struct {
	LogLevel int

	InfoLogger    *log.Logger
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
}

// Enumeration for various log levels
const (
	ERROR int = iota
	WARNING
	INFO
)

// NewLogger creates a new logger at a certain log level
func NewLogger(logLevel int) *Logger {
	return &Logger{
		LogLevel:      logLevel,
		InfoLogger:    log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile),
		WarningLogger: log.New(os.Stdout, "[WARNING] ", log.Ldate|log.Ltime|log.Lshortfile),
		ErrorLogger:   log.New(os.Stdout, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

// Log logs a message if the level matches the LogLevel of the configuration.
func (l Logger) Log(level int, strings ...interface{}) {
	switch level {
	case ERROR:
		l.ErrorLogger.Println(strings...)
	case WARNING:
		if l.LogLevel >= 1 {
			l.WarningLogger.Println(strings...)
		}
	case INFO:
		if l.LogLevel == 2 {
			l.InfoLogger.Println(strings...)
		}
	}
}

// Check checks if an error exists -- if it does, logs an error and panics.
func (l Logger) Check(err error, strings ...interface{}) {
	if err != nil {
		strings = append(strings, err.Error())
		l.Log(ERROR, strings...)
		panic(err)
	}
}

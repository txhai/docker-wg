package loggin

import (
	"log"
	"os"
)

type LogLevel int

const (
	LevelSilent LogLevel = iota
	LevelError
	LevelVerbose
)

type Logger struct {
	Printf func(format string, args ...interface{})
	Errorf func(format string, args ...interface{})
	Close  func()
}

// discardLogf Function for use in Logger for discarding logged lines.
func discardLogf(format string, args ...interface{}) {}

func NewLogger(level LogLevel) *Logger {
	logger := &Logger{discardLogf, discardLogf, func() {}}
	logf := func(prefix string) func(string, ...interface{}) {
		return log.New(os.Stdout, prefix+" ", log.Ldate|log.Ltime).Printf
	}
	if level >= LevelVerbose {
		logger.Printf = logf("DEBUG")
	}
	if level >= LevelError {
		logger.Errorf = logf("ERROR")
	}
	return logger
}

package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var _ Logger = (*StdLogger)(nil)

type StdLogger struct {
	level  int
	logger *log.Logger
}

func NewStdLogger(level int) *StdLogger {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	return &StdLogger{
		level:  level,
		logger: logger,
	}
}

func (l *StdLogger) Debug(msg string, v ...any) {
	if LevelDebug < l.level {
		return
	}
	l.logger.Printf("[DEBUG] "+msg, v...)
}

func (l *StdLogger) Info(msg string, v ...any) {
	if LevelInfo < l.level {
		return
	}
	l.logger.Printf("[INFO]  "+msg, v...)
}

func (l *StdLogger) Warn(msg string, v ...any) {
	if LevelWarn < l.level {
		return
	}
	l.logger.Printf("[WARN]  "+msg, v...)
}

func (l *StdLogger) Error(msg string, v ...any) {
	if LevelError < l.level {
		return
	}
	l.logger.Printf("[ERROR] "+msg, v...)
}

func (l *StdLogger) Fatal(msg string, v ...any) {
	if LevelFatal < l.level {
		return
	}
	l.logger.Printf("[FATAL] "+msg, v...)
	os.Exit(1)
}

func (l *StdLogger) Dump(path string, msg string, v ...any) {
	ext := filepath.Ext(path)
	name := strings.TrimSuffix(path, ext)
	uniquePath := fmt.Sprintf("%s_%d%s", name, time.Now().Unix(), ext)

	f, err := os.Create(uniquePath)
	if err != nil {
		return
	}
	defer f.Close()

	fmt.Fprintf(f, msg, v...)
}

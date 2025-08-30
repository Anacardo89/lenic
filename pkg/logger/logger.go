package logger

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/Anacardo89/lenic/config"
	"github.com/natefinch/lumberjack"
)

type Logger struct {
	log   *slog.Logger
	level slog.Level
}

func NewLogger(cfg config.Log) *Logger {
	level := slog.LevelInfo
	switch strings.ToLower(cfg.Level) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	case "fatal":
		level = slog.LevelError
	}
	lj := &lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/%s", cfg.Path, cfg.File),
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}
	fileJSONHandler := slog.NewJSONHandler(lj, &slog.HandlerOptions{AddSource: true})
	stderrHandler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{AddSource: true})
	multiHandler := NewMultiHandler(fileJSONHandler, stderrHandler)
	return &Logger{
		log:   slog.New(multiHandler),
		level: level,
	}
}

func (l *Logger) logWithLevel(lvl slog.Level, msg string, args ...any) {
	if lvl >= l.level {
		switch lvl {
		case slog.LevelDebug:
			l.log.Debug(msg, args...)
		case slog.LevelInfo:
			l.log.Info(msg, args...)
		case slog.LevelWarn:
			l.log.Warn(msg, args...)
		case slog.LevelError:
			l.log.Error(msg, args...)
		}
	}
}

func (l *Logger) Info(msg string, args ...any) {
	l.logWithLevel(slog.LevelInfo, msg, args...)
}
func (l *Logger) Error(msg string, args ...any) {
	l.logWithLevel(slog.LevelError, msg, args...)
}
func (l *Logger) Debug(msg string, args ...any) {
	l.logWithLevel(slog.LevelDebug, msg, args...)
}
func (l *Logger) Warn(msg string, args ...any) {
	l.logWithLevel(slog.LevelWarn, msg, args...)
}
func (l *Logger) Fatal(msg string, args ...any) {
	l.logWithLevel(slog.LevelError, msg, args...)
	os.Exit(1)
}

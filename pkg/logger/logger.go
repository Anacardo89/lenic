package logger

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/natefinch/lumberjack"

	"github.com/Anacardo89/lenic/config"
)

type Logger struct {
	*slog.Logger
}

func NewLogger(cfg *config.Log, homeDir string) (*Logger, error) {
	if err := os.MkdirAll(filepath.Join(homeDir, cfg.Path), 0755); err != nil {
		return nil, err
	}
	logLevel := strings.ToUpper(cfg.Level)
	level := slog.LevelInfo
	switch logLevel {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "WARN", "WARNING":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	}
	lj := &lumberjack.Logger{
		Filename:   strings.Split(cfg.Path, "/")[1],
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}
	fileJSONHandler := NewLoggerHandler(lj, level)
	stderrHandler := NewLoggerHandler(os.Stderr, level)
	multiHandler := NewMultiHandler(fileJSONHandler, stderrHandler)
	logger := slog.New(multiHandler)
	slog.SetDefault(logger)
	return &Logger{
		logger,
	}, nil
}

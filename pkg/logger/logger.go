package logger

import (
	"errors"
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

func NewLogger(cfg *config.Log, homeDir string, appEnv string) (*Logger, error) {
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
	var handler slog.Handler
	if appEnv == "aws" {
		handler = NewLoggerHandler(os.Stdout, level)
	} else {
		if err := os.MkdirAll(filepath.Join(homeDir, cfg.Path), 0755); err != nil {
			return nil, err
		}
		fileName := strings.Split(cfg.Path, "/")
		if len(fileName) < 2 {
			return nil, errors.New("badly formed log path, ensure the format is: <logDir>/<logFile>.log")
		}
		lj := &lumberjack.Logger{
			Filename:   fileName[1],
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		}
		fileJSONHandler := NewLoggerHandler(lj, level)
		stdoutHandler := NewLoggerHandler(os.Stdout, level)
		handler = NewMultiHandler(fileJSONHandler, stdoutHandler)
	}
	logger := slog.New(handler)
	slog.SetDefault(logger)
	return &Logger{
		logger,
	}, nil
}

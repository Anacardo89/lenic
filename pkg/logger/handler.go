package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type LoggerHandler struct {
	level   slog.Level
	handler slog.Handler
}

type LogRecord struct {
	Timestamp string         `json:"timestamp"`
	Level     string         `json:"level"`
	Location  string         `json:"location"`
	Msg       string         `json:"msg"`
	Attrs     map[string]any `json:"-"`
}

func NewLoggerHandler(out io.Writer, level slog.Level) *LoggerHandler {
	jsonHandler := slog.NewJSONHandler(out, &slog.HandlerOptions{
		AddSource: true,
		Level:     level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			return a
		},
	})
	return &LoggerHandler{
		level:   level,
		handler: jsonHandler,
	}
}

// Implements slog.Handler
func (h *LoggerHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *LoggerHandler) Handle(ctx context.Context, r slog.Record) error {
	location := ""
	pcs := make([]uintptr, 32)
	n := runtime.Callers(3, pcs)
	frames := runtime.CallersFrames(pcs[:n])
	for {
		f, more := frames.Next()
		if !more {
			break
		}
		if strings.HasPrefix(f.Function, "runtime.") ||
			strings.HasPrefix(f.Function, "log/slog.") {
			continue
		}
		parent := filepath.Base(filepath.Dir(f.File))
		base := filepath.Base(f.File)
		file := filepath.Join(parent, base)
		fn := ""
		if idx := strings.LastIndex(f.Function, "."); idx != -1 {
			fn = f.Function[idx+1:]
		} else {
			fn = f.Function
		}

		location = fmt.Sprintf("%s:%d %s", file, f.Line, fn)
		break
	}
	attrs := map[string]any{}
	r.Attrs(func(a slog.Attr) bool {
		v := a.Value.Any()
		if err, ok := v.(error); ok {
			attrs[a.Key] = err.Error()
		} else {
			attrs[a.Key] = v
		}
		return true
	})
	lr := LogRecord{
		Timestamp: r.Time.Format("2006/01/02 15:04:05"),
		Level:     r.Level.String(),
		Location:  location,
		Msg:       r.Message,
		Attrs:     attrs,
	}
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, `{"timestamp":%q,"level":%q,"location":%q,"msg":%q`,
		lr.Timestamp, lr.Level, lr.Location, lr.Msg)
	for k, v := range lr.Attrs {
		b, err := json.Marshal(v)
		if err != nil {
			continue
		}
		fmt.Fprintf(buf, `,"%s":%s`, k, b)
	}
	buf.WriteString("}\n")

	_, err := os.Stdout.Write(buf.Bytes())
	return err
}

func (h *LoggerHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &LoggerHandler{
		level:   h.level,
		handler: h.handler.WithAttrs(attrs),
	}
}

func (h *LoggerHandler) WithGroup(name string) slog.Handler {
	return &LoggerHandler{
		level:   h.level,
		handler: h.handler.WithGroup(name),
	}
}

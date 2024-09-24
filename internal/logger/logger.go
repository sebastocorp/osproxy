package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"slices"
)

// ----------------------------------------------------------------
// LOGGER
// ----------------------------------------------------------------

type LevelT int

const (
	DEBUG LevelT = LevelT(slog.LevelDebug)
	INFO  LevelT = LevelT(slog.LevelInfo)
	WARN  LevelT = LevelT(slog.LevelWarn)
	ERROR LevelT = LevelT(slog.LevelError)
)

type LoggerT struct {
	SLogger   *slog.Logger
	Context   context.Context
	ExtraAttr []any
}

var Log LoggerT

func InitLogger(ctx context.Context, level LevelT, extraAttr ...any) {
	opts := &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.Level(level),
	}
	jsonHandler := slog.NewJSONHandler(os.Stdout, opts)
	Log.SLogger = slog.New(jsonHandler)
	Log.Context = ctx
	Log.ExtraAttr = extraAttr
}

func GetLevel(levelStr string) (l LevelT, err error) {
	levelMap := map[string]LevelT{
		"debug": DEBUG,
		"info":  INFO,
		"warn":  WARN,
		"error": ERROR,
	}
	if l, ok := levelMap[levelStr]; ok {
		return l, err
	}
	l = INFO
	err = fmt.Errorf("log level '%s' not supported", levelStr)
	return l, err
}

func (l *LoggerT) Debugf(extra []any, format string, args ...any) {
	extra = slices.Concat(l.ExtraAttr, extra)
	if l.Context != nil {
		l.SLogger.DebugContext(l.Context, fmt.Sprintf(format, args...), extra...)
		return
	}
	l.SLogger.Debug(fmt.Sprintf(format, args...), extra...)
}

func (l *LoggerT) Infof(extra []any, format string, args ...any) {
	extra = slices.Concat(l.ExtraAttr, extra)
	if l.Context != nil {
		l.SLogger.InfoContext(l.Context, fmt.Sprintf(format, args...), extra...)
		return
	}
	l.SLogger.Info(fmt.Sprintf(format, args...), extra...)
}

func (l *LoggerT) Warnf(extra []any, format string, args ...any) {
	extra = slices.Concat(l.ExtraAttr, extra)
	if l.Context != nil {
		l.SLogger.WarnContext(l.Context, fmt.Sprintf(format, args...), extra...)
		return
	}
	l.SLogger.Warn(fmt.Sprintf(format, args...), extra...)
}

func (l *LoggerT) Errorf(extra []any, format string, args ...any) {
	extra = slices.Concat(l.ExtraAttr, extra)
	if l.Context != nil {
		l.SLogger.ErrorContext(l.Context, fmt.Sprintf(format, args...), extra...)
		return
	}
	l.SLogger.Error(fmt.Sprintf(format, args...), extra...)
}

func (l *LoggerT) Fatalf(extra []any, format string, args ...any) {
	extra = slices.Concat(l.ExtraAttr, extra)
	if l.Context != nil {
		l.SLogger.ErrorContext(l.Context, fmt.Sprintf(format, args...), extra...)
		os.Exit(1)
	}
	l.SLogger.Error(fmt.Sprintf(format, args...), extra...)
	os.Exit(1)
}

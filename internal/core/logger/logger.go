// Package logger provides a unified structured logging interface and context field propagation.
package logger

// 本文件实现统一业务日志入口和上下文字段透传能力。

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"gotribe/internal/core/config"
)

// contextKey 用于在上下文中存放日志字段的键类型。
type contextKey string

const (
	contextKeyRequestID contextKey = "request_id"
	contextKeyProjectID contextKey = "project_id"
	contextKeyUserID    contextKey = "user_id"
	contextKeyUsername  contextKey = "username"
	contextKeyTraceID   contextKey = "trace_id"
	contextKeySpanID    contextKey = "span_id"
)

// defaultLogger 全局 SugaredLogger 实例。
var defaultLogger *zap.SugaredLogger

func init() {
	logger, _ := zap.NewDevelopment()
	defaultLogger = logger.Sugar()
}

// Init 按运行环境初始化全局日志器。
func Init(app config.AppConfig) {
	var cfg zap.Config
	if app.IsProduction() {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
	}
	cfg.Level = zap.NewAtomicLevelAt(levelForEnv(app))

	logger, err := cfg.Build()
	if err != nil {
		return
	}
	defaultLogger = logger.Sugar()
}

// Sugared 返回全局 SugaredLogger 实例。
func Sugared() *zap.SugaredLogger {
	return defaultLogger
}

// Debug 输出调试级别日志。
func Debug(ctx context.Context, msg string, args ...any) {
	args = append(contextArgs(ctx), args...)
	defaultLogger.Debugw(msg, args...)
}

// Info 输出信息级别日志。
func Info(ctx context.Context, msg string, args ...any) {
	args = append(contextArgs(ctx), args...)
	defaultLogger.Infow(msg, args...)
}

// Warn 输出警告级别日志。
func Warn(ctx context.Context, msg string, args ...any) {
	args = append(contextArgs(ctx), args...)
	defaultLogger.Warnw(msg, args...)
}

// Error 输出错误级别日志。
func Error(ctx context.Context, msg string, args ...any) {
	args = append(contextArgs(ctx), args...)
	defaultLogger.Errorw(msg, args...)
}

// Fatal 输出致命级别日志并退出程序。
func Fatal(ctx context.Context, msg string, args ...any) {
	args = append(contextArgs(ctx), args...)
	defaultLogger.Fatalw(msg, args...)
}

// WithRequestID 向上下文写入请求 ID。
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return withString(ctx, contextKeyRequestID, requestID)
}

// WithProjectID 向上下文写入项目 ID。
func WithProjectID(ctx context.Context, projectID string) context.Context {
	return withString(ctx, contextKeyProjectID, projectID)
}

// WithUsername 向上下文写入用户名。
func WithUsername(ctx context.Context, username string) context.Context {
	return withString(ctx, contextKeyUsername, username)
}

// WithUserID 向上下文写入用户 ID。
func WithUserID(ctx context.Context, userID int64) context.Context {
	if userID <= 0 {
		return ctx
	}
	return context.WithValue(ctx, contextKeyUserID, userID)
}

// WithTrace 向上下文写入 trace_id 与 span_id。
func WithTrace(ctx context.Context, traceID, spanID string) context.Context {
	ctx = withString(ctx, contextKeyTraceID, traceID)
	return withString(ctx, contextKeySpanID, spanID)
}

// withString 向上下文写入字符串类型日志字段。
func withString(ctx context.Context, key contextKey, value string) context.Context {
	value = strings.TrimSpace(value)
	if value == "" {
		return ctx
	}
	return context.WithValue(ctx, key, value)
}

// levelForEnv 根据运行环境选择默认日志级别。
func levelForEnv(app config.AppConfig) zapcore.Level {
	switch {
	case app.IsDevelopment():
		return zapcore.DebugLevel
	case app.IsProduction():
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// InitWithRotation 初始化带日志轮转的全局日志器（供 Admin 侧使用）。
func InitWithRotation(path string, level zapcore.Level, maxSize, maxBackups, maxAge int, compress bool) {
	now := time.Now()
	infoLogFileName := fmt.Sprintf("%s/info/%04d-%02d-%02d.log", path, now.Year(), now.Month(), now.Day())
	errorLogFileName := fmt.Sprintf("%s/error/%04d-%02d-%02d.log", path, now.Year(), now.Month(), now.Day())

	encoderConfig := zapcore.EncoderConfig{
		MessageKey:    "msg",
		LevelKey:      "level",
		TimeKey:       "time",
		NameKey:       "name",
		CallerKey:     "file",
		FunctionKey:   "func",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zap.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zap.ErrorLevel && lvl >= zap.DebugLevel
	})

	if level >= zap.ErrorLevel {
		lowPriority = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return false
		})
	}

	infoFileSync := zapcore.AddSync(&lumberjack.Logger{
		Filename:   infoLogFileName,
		MaxSize:    maxSize,
		MaxAge:     maxAge,
		MaxBackups: maxBackups,
		LocalTime:  false,
		Compress:   compress,
	})
	infoFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(infoFileSync, zapcore.AddSync(os.Stdout)), lowPriority)

	errorFileSync := zapcore.AddSync(&lumberjack.Logger{
		Filename:   errorLogFileName,
		MaxSize:    maxSize,
		MaxAge:     maxAge,
		MaxBackups: maxBackups,
		LocalTime:  false,
		Compress:   compress,
	})
	errorFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(errorFileSync, zapcore.AddSync(os.Stdout)), highPriority)

	logger := zap.New(zapcore.NewTee(infoFileCore, errorFileCore), zap.AddCaller())
	defaultLogger = logger.Sugar()
	defaultLogger.Info("initialized zap logger with rotation")
}

// contextArgs 从上下文中提取需要写入日志的字段。
func contextArgs(ctx context.Context) []any {
	if ctx == nil {
		return nil
	}

	args := make([]any, 0, 8)
	if requestID, ok := ctx.Value(contextKeyRequestID).(string); ok && requestID != "" {
		args = append(args, "request_id", requestID)
	}
	if projectID, ok := ctx.Value(contextKeyProjectID).(string); ok && projectID != "" {
		args = append(args, "project_id", projectID)
	}
	if userID, ok := ctx.Value(contextKeyUserID).(int64); ok && userID > 0 {
		args = append(args, "user_id", userID)
	}
	if username, ok := ctx.Value(contextKeyUsername).(string); ok && username != "" {
		args = append(args, "username", username)
	}
	if traceID, ok := ctx.Value(contextKeyTraceID).(string); ok && traceID != "" {
		args = append(args, "trace_id", traceID)
	}
	if spanID, ok := ctx.Value(contextKeySpanID).(string); ok && spanID != "" {
		args = append(args, "span_id", spanID)
	}
	return args
}

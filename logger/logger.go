package logger

import (
	"bluebell/settings"
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	lgConfig    *settings.LoggerConfig
	lg          *zap.Logger
	AtomicLevel = zap.NewAtomicLevel()
)

const TraceIDKey = "trace_id"

// InitLogger 初始化

func InitLogger(cfg *settings.LoggerConfig) {
	AtomicLevel.SetLevel(getLogLevel(cfg.Level))

	// 过滤器：只记录 ERROR 及以上
	errorEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zap.ErrorLevel
	})

	var cores []zapcore.Core

	// 全量日志
	allWriter := getLogWriter(cfg.Filename, cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge, cfg.Compress)
	cores = append(cores, zapcore.NewCore(getFileEncoder(), allWriter, AtomicLevel))

	// 错误日志分流
	errWriter := getLogWriter(cfg.ErrorName, cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge, cfg.Compress)
	cores = append(cores, zapcore.NewCore(getFileEncoder(), errWriter, errorEnabler))

	if cfg.LogInConsole {
		cores = append(cores, zapcore.NewCore(getConsoleEncoder(), zapcore.AddSync(os.Stdout), AtomicLevel))
	}

	lg = zap.New(zapcore.NewTee(cores...), zap.AddCaller())
	zap.ReplaceGlobals(lg)
}

// L (Context Logger) 是核心技巧：从 Context 中提取 TraceID 并返回带字段的 Logger
func L(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return lg
	}
	// 兼容 Gin Context 和标准 Context
	var traceID string
	if gCtx, ok := ctx.(*gin.Context); ok {
		traceID = gCtx.GetString(TraceIDKey)
	} else {
		if id, ok := ctx.Value(TraceIDKey).(string); ok {
			traceID = id
		}
	}

	if traceID != "" {
		return lg.With(zap.String(TraceIDKey, traceID))
	}
	return lg
}

// --- 内部辅助函数 ---
func getFileEncoder() zapcore.Encoder {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(config)
}

func getConsoleEncoder() zapcore.Encoder {
	config := zap.NewDevelopmentEncoderConfig()
	config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	return zapcore.NewConsoleEncoder(config)
}

func getLogWriter(f string, s, b, a int, c bool) zapcore.WriteSyncer {
	return zapcore.AddSync(&lumberjack.Logger{
		Filename: f, MaxSize: s, MaxBackups: b, MaxAge: a, Compress: c,
	})
}

func getLogLevel(s string) zapcore.Level {
	var l zapcore.Level
	if err := l.UnmarshalText([]byte(s)); err != nil {
		return zap.InfoLevel
	}
	return l
}

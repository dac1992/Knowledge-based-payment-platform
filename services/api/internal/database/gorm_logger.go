package database

import (
	"context"
	"errors"
	"strings"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/yanmai/api/internal/config"
)

// gormZapLogger 把 GORM 的 SQL 日志桥接到 zap。
type gormZapLogger struct {
	log           *zap.Logger
	level         gormlogger.LogLevel
	slowThreshold time.Duration
}

func newGormLogger(log *zap.Logger, cfg *config.Config) gormlogger.Interface {
	level := gormlogger.Warn
	if cfg.App.IsDevelopment() && strings.EqualFold(cfg.Log.Level, "debug") {
		level = gormlogger.Info
	}
	return &gormZapLogger{
		log:           log.Named("gorm"),
		level:         level,
		slowThreshold: 500 * time.Millisecond,
	}
}

func (l *gormZapLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	clone := *l
	clone.level = level
	return &clone
}

func (l *gormZapLogger) Info(_ context.Context, msg string, data ...any) {
	if l.level >= gormlogger.Info {
		l.log.Sugar().Infof(msg, data...)
	}
}

func (l *gormZapLogger) Warn(_ context.Context, msg string, data ...any) {
	if l.level >= gormlogger.Warn {
		l.log.Sugar().Warnf(msg, data...)
	}
}

func (l *gormZapLogger) Error(_ context.Context, msg string, data ...any) {
	if l.level >= gormlogger.Error {
		l.log.Sugar().Errorf(msg, data...)
	}
}

func (l *gormZapLogger) Trace(_ context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.level <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()
	fields := []zap.Field{
		zap.String("sql", sql),
		zap.Int64("rows", rows),
		zap.Duration("elapsed", elapsed),
	}

	switch {
	case err != nil && !errors.Is(err, gorm.ErrRecordNotFound):
		l.log.Error("SQL 执行失败", append(fields, zap.Error(err))...)
	case elapsed > l.slowThreshold && l.slowThreshold > 0 && l.level >= gormlogger.Warn:
		l.log.Warn("SQL 执行缓慢", fields...)
	case l.level >= gormlogger.Info:
		l.log.Debug("SQL", fields...)
	}
}

package log

import (
	"context"
	"time"

	gormLogger "gorm.io/gorm/logger"
)

type GormLogger struct {
	SlowThreshold time.Duration
}

var _ gormLogger.Interface = (*GormLogger)(nil)

func NewGormLogger() *GormLogger {
	return &GormLogger{
		SlowThreshold: 200 * time.Millisecond, // 一般超过200毫秒就算慢查所以不使用配置进行更改
	}
}

var _ gormLogger.Interface = (*GormLogger)(nil)

func (l *GormLogger) LogMode(lev gormLogger.LogLevel) gormLogger.Interface {
	return &GormLogger{}
}
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	InfoF(msg, data)
}
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	WarnF(msg, data)
}
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	ErrorF(msg, data)
}
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	if err != nil {
		ErrorFTimes(5, "[SQL-Error]\t| err = %v \t| rows= %v \t| %v \t| %v", err, rows, elapsed, sql)
	}
	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		WarnFTimes(5, "[SQL-SlowLog]\t| rows= %v \t| %v \t| %v", rows, elapsed, sql)
	}
	InfoFTimes(5, "[SQL]\t| rows= %v \t| %v \t| %v", rows, elapsed, sql)
}

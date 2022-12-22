package database

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"runtime"
	"time"
	"ws/framework/plugin/logger"
)

var dataBaseLogger = logger.New("")
var (
	traceStr     = "[%.3fms] [rows:%v] %s"
	traceWarnStr = "%s[%.3fms] [rows:%v] %s"
	traceErrStr  = "%s %s"
)

type dbLoggerProxy struct {
	SlowThreshold time.Duration
	traceStr      string
	traceWarnStr  string
	traceErrStr   string
}

func newDBLoggerProxy() dbLoggerProxy {
	p := dbLoggerProxy{SlowThreshold: 200 * time.Millisecond}

	if runtime.GOOS == "windows" {
		p.traceStr = gormLogger.Yellow + "[%.3fms] " + gormLogger.BlueBold + "[rows:%v]" + gormLogger.Green + " %s" + gormLogger.Reset
		p.traceWarnStr = gormLogger.MagentaBold + "%s" + gormLogger.Red + "[%.3fms] " + gormLogger.BlueBold + "[rows:%v]" + gormLogger.Green + " %s" + gormLogger.Reset
		p.traceErrStr = gormLogger.MagentaBold + "%s" + gormLogger.Green + " %s" + gormLogger.Reset
	} else {
		p.traceStr = traceStr
		p.traceWarnStr = traceWarnStr
		p.traceErrStr = traceErrStr
	}

	return p
}

// LogMode .
func (l dbLoggerProxy) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	return l
}

// Info print info
func (l dbLoggerProxy) Info(ctx context.Context, msg string, data ...interface{}) {
	dataBaseLogger.Infof(msg, data)
}

// Warn print warn messages
func (l dbLoggerProxy) Warn(ctx context.Context, msg string, data ...interface{}) {
	dataBaseLogger.Warnf(msg, data)
}

// Error print error messages
func (l dbLoggerProxy) Error(ctx context.Context, msg string, data ...interface{}) {
	dataBaseLogger.Errorf(msg, data)
}

// Trace print sql message
func (l dbLoggerProxy) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	switch {
	case err != nil && (!errors.Is(err, gorm.ErrRecordNotFound)):
		sql, rows := fc()
		if rows == -1 {
			dataBaseLogger.Warnf(l.traceErrStr, err, sql)
		} else {
			dataBaseLogger.Warnf(l.traceErrStr, err, sql)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			dataBaseLogger.Warnf(l.traceWarnStr, slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			dataBaseLogger.Warnf(l.traceWarnStr, slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	default:
		if !logger.EnabledDebug() {
			return
		}

		sql, rows := fc()
		if rows == -1 {
			dataBaseLogger.Debugf(l.traceStr, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			dataBaseLogger.Debugf(l.traceStr, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}

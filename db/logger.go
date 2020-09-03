package db

import (
	"context"
	"github.com/karta0898098/kara/ctxutil"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"time"
)

// Colors
const (
	Reset       = "\033[0m"
	Red         = "\033[31m"
	Green       = "\033[32m"
	Yellow      = "\033[33m"
	Blue        = "\033[34m"
	Magenta     = "\033[35m"
	Cyan        = "\033[36m"
	White       = "\033[37m"
	MagentaBold = "\033[35;1m"
	RedBold     = "\033[31;1m"
	YellowBold  = "\033[33;1m"
)

type gormLogger struct {
	LogLevel                            logger.LogLevel
	Config                              logger.Config
	SlowThreshold                       time.Duration
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

func NewLogger(config logger.Config) logger.Interface {
	var (
		infoStr      = "%s\n[info] "
		warnStr      = "%s\n[warn] "
		errStr       = "%s\n[error] "
		traceStr     = "%s\n[%v]\n[rows:%d]\n%s"
		traceWarnStr = "%s\n[%v]\n[rows:%d]\n%s"
		traceErrStr  = "%s\n%s\n[%v]\n[rows:%d]\n%s"
	)

	if config.Colorful {
		infoStr = Green + "%s\n" + Reset + Green + "[info] " + Reset
		warnStr = Blue + "%s\n" + Reset + Magenta + "[warn] " + Reset
		errStr = Magenta + "%s\n" + Reset + Red + "[error] " + Reset
		traceStr = "\t" + Green + "%s\n" + Reset + Yellow + "\t[%.3fms]\n " + Blue + "\t[rows:%d]\n" + Reset + "\t%s"
		traceWarnStr = "\t" + Green + "%s\n" + Reset + RedBold + "\t[%.3fms]\n " + Yellow + "\t[rows:%d]\n" + Magenta + "\t%s" + Reset
		traceErrStr = "\t" + RedBold + "%s " + MagentaBold + "%s\n" + Reset + Yellow + "\t[%.3fms] " + Blue + "\t[rows:%d]\n" + Reset + "\t%s"
	}

	l := &gormLogger{
		LogLevel:      config.LogLevel,
		SlowThreshold: config.SlowThreshold,
		Config:        config,
		infoStr:       infoStr,
		warnStr:       warnStr,
		errStr:        errStr,
		traceStr:      traceStr,
		traceWarnStr:  traceWarnStr,
		traceErrStr:   traceErrStr,
	}

	return l
}

func (g *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newlogger := *g
	newlogger.LogLevel = level
	return g
}

func (g *gormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	currentLogger := log.With().Fields(map[string]interface{}{
		"request_id": ctxutil.GetRequestID(ctx),
		"db_log":     true,
	}).Logger()

	if g.LogLevel >= logger.Info {
		currentLogger.Info().Msgf(
			g.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

func (g *gormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	currentLogger := log.With().Fields(map[string]interface{}{
		"request_id": ctxutil.GetRequestID(ctx),
		"db_log":     true,
	}).Logger()

	if g.LogLevel >= logger.Warn {
		currentLogger.Warn().Msgf(
			g.warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

func (g *gormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	currentLogger := log.With().Fields(map[string]interface{}{
		"request_id": ctxutil.GetRequestID(ctx),
		"db_log":     true,
	}).Logger()

	if g.LogLevel >= logger.Error {
		currentLogger.Error().Msgf(g.errStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

func (g *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {

	currentLogger := log.With().Fields(map[string]interface{}{
		"request_id": ctxutil.GetRequestID(ctx),
		"db_log":     true,
	}).Logger()

	if g.LogLevel > 0 {
		elapsed := time.Since(begin)
		switch {
		case err != nil && g.LogLevel >= logger.Error:
			sql, rows := fc()
			currentLogger.Error().Msgf(g.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		case elapsed > g.SlowThreshold && g.SlowThreshold != 0 && g.LogLevel >= logger.Warn:
			sql, rows := fc()
			currentLogger.Warn().Msgf(g.traceWarnStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		case g.LogLevel >= logger.Info:
			sql, rows := fc()
			currentLogger.Info().Msgf(g.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}

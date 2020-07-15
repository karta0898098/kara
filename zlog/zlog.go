package zlog

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"

	"os"
	"strconv"
	"time"
)

var (
	// Teal ...
	Teal = Color("\033[1;36m%s\033[0m")
	// Yello ...
	Yello = Color("\033[35m%s\033[0m")
)

var Logger zerolog.Logger

// Color ...
func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}

// Graylog 的錯誤等級
const (
	levelEmerg   = int8(0)
	levelAlert   = int8(1)
	levelCrit    = int8(2)
	levelErr     = int8(3)
	levelWarning = int8(4)
	levelNotice  = int8(5)
	levelInfo    = int8(6)
	levelDebug   = int8(7)
)

type severityHook struct{}

// Run ...
func (h severityHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	lvl := int8(0)
	switch level {
	case zerolog.DebugLevel:
		lvl = levelDebug
	case zerolog.InfoLevel:
		lvl = levelInfo
	case zerolog.WarnLevel:
		lvl = levelWarning
	case zerolog.ErrorLevel:
		lvl = levelErr
	case zerolog.FatalLevel:
		lvl = levelCrit
	}
	e.Int8("log_level", lvl).
		Float64("timestamp", float64(time.Now().UnixNano()/int64(time.Millisecond))/1000)
	if msg == "" {
		e.Str("message", "no message")
	}
}

func Setup(config *Config) {
	zerolog.DisableSampling(true)
	zerolog.TimestampFieldName = "local_timestamp"
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	level := zerolog.InfoLevel
	if config.Debug {
		level = zerolog.DebugLevel
	}


	if config.Local {
		output := zerolog.ConsoleWriter{
			Out: os.Stdout,
		}
		output.FormatMessage = func(i interface{}) string {
			return fmt.Sprintf("[ %s ]", i)
		}
		output.FormatFieldName = func(i interface{}) string {
			return fmt.Sprintf("%s:", Teal(i))
		}
		output.FormatFieldValue = func(i interface{}) string {
			return fmt.Sprintf("%s", i)
		}
		output.FormatTimestamp = func(i interface{}) string {
			t := fmt.Sprintf("%s", i)
			millisecond, err := strconv.ParseInt(fmt.Sprintf("%s", i), 10, 64)
			if err == nil {
				t = time.Unix(int64(millisecond/1000), 0).Local().Format("2006/01/02 15:04:05")
			}
			return Yello(t)
		}
		Logger = zerolog.New(output)
	} else {
		Logger = zerolog.New(os.Stdout)
	}

	log.Logger = Logger.Hook(severityHook{}).
		With().
		Fields(map[string]interface{}{
			"app_id": config.AppID,
			"env":    config.Env,
		}).
		Timestamp().
		Logger().
		Level(level)
}

// Ctx wrap zerolog Ctx func, if ctx not setting Logger, return a default prevent for panic
func Ctx(ctx context.Context) *zerolog.Logger {
	defalutLogger := log.Logger
	if ctx == nil {
		defalutLogger.Warn().Msg("zlog func Ctx() not set context.Context in right way.")
		return &defalutLogger
	}

	return log.Ctx(ctx) // if ctx is not null and not set Logger yet. A disabled Logger is returned.
}

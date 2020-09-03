package zlog

import (
	"fmt"
	"strings"

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
	// Yellow ...
	Yellow = Color("\033[35m%s\033[0m")
)

const (
	colorBlack = iota + 30
	colorRed
	colorGreen
	colorYellow
	colorBlue
	colorMagenta
	colorCyan
	colorWhite

	colorBold     = 1
	colorDarkGray = 90
)

const spcae = "                    "

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
		output.FormatLevel = func(i interface{}) string {
			var l string
			if ll, ok := i.(string); ok {
				switch ll {
				case "trace":
					l = colorize("TRACE", colorMagenta)
				case "debug":
					l = colorize("DEBUG", colorYellow)
				case "info":
					l = colorize("INFO", colorGreen)
				case "warn":
					l = colorize("WARN", colorRed)
				case "error":
					l = colorize(colorize("ERROR", colorRed), colorBold)
				case "fatal":
					l = colorize(colorize("FATAL", colorRed), colorBold)
				case "panic":
					l = colorize(colorize("PANIC", colorRed), colorBold)
				default:
					l = colorize("???", colorBold)
				}
			} else {
				if i == nil {
					l = colorize("???", colorBold)
				} else {
					l = strings.ToUpper(fmt.Sprintf("%s", i))[0:3]
				}
			}
			return fmt.Sprintf("|%s|", l)
		}
		output.FormatMessage = func(i interface{}) string {
			return fmt.Sprintf("[%s]\n", i)
		}
		output.FormatFieldName = func(i interface{}) string {
			return fmt.Sprintf(spcae+"%s: ", Teal(i))
		}
		output.FormatFieldValue = func(i interface{}) string {
			return fmt.Sprintf("%s\n", i)
		}
		output.FormatTimestamp = func(i interface{}) string {
			t := fmt.Sprintf("%v", i)
			millisecond, err := strconv.ParseInt(fmt.Sprintf("%s", i), 10, 64)
			if err == nil {
				t = time.Unix(int64(millisecond/1000), 0).Local().Format("2006/01/02 15:04:05")
			}
			return Yellow(t)
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

// colorize returns the string s wrapped in ANSI code c, unless disabled is true.
func colorize(s interface{}, c int) string {
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", c, s)
}

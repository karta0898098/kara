package zlog

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

var (
	// Teal ...
	Teal = Color("\033[1;36m%s\033[0m")
	// Yellow ...
	Yellow = Color("\033[35m%s\033[0m")
	// Green
	Green = Color("\033[32m%s\033[0m")
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

var Logger zerolog.Logger

// Color ...
func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}

type severityHook struct{}

// Run ...
func (h severityHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	e.Float64("timestamp", float64(time.Now().UnixNano()/int64(time.Millisecond))/1000)
	if msg == "" {
		e.Str("message", "no message")
	}
}

func Setup(config Config) {
	zerolog.DisableSampling(true)
	zerolog.TimestampFieldName = "time"
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	level := config.Level
	if config.Debug {
		output := zerolog.ConsoleWriter{
			Out: os.Stdout,
		}
		output.FormatLevel = func(i interface{}) string {
			var l string
			if ll, ok := i.(string); ok {
				switch ll {
				case "trace":
					l = colorize("Trace", colorMagenta)
				case "debug":
					l = colorize("Debug", colorBlue)
				case "info":
					l = colorize("Info ", colorGreen)
				case "warn":
					l = colorize("Warn ", colorYellow)
				case "error":
					l = colorize(colorize("Error", colorRed), colorBold)
				case "fatal":
					l = colorize(colorize("Fatal", colorRed), colorBold)
				case "panic":
					l = colorize(colorize("Panic", colorRed), colorBold)
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
			return fmt.Sprintf("%-50s", i)
		}
		output.FormatFieldName = func(i interface{}) string {
			return fmt.Sprintf("%s = ", Teal(i))
		}
		output.FormatFieldValue = func(i interface{}) string {
			return fmt.Sprintf("%s", i)
		}
		output.FormatTimestamp = func(i interface{}) string {
			t := fmt.Sprintf("%v", i)
			millisecond, err := strconv.ParseInt(fmt.Sprintf("%s", i), 10, 64)
			if err == nil {
				t = time.Unix(int64(millisecond/1000), 0).Local().Format("2006/01/02 15:04:05")
			}
			return colorize(t, colorCyan)
		}
		output.FormatCaller = func(i interface{}) string {
			var c string
			if cc, ok := i.(string); ok {
				c = cc
			}
			if len(c) > 0 {
				cwd, err := os.Getwd()
				if err == nil {
					c = strings.TrimPrefix(c, cwd)
					c = strings.TrimPrefix(c, "/")
				}
				c = colorize(c, colorGreen)

				if c != "" {
					c = fmt.Sprintf("%s %s", " >", c)
				}
			}
			return c
		}

		output.PartsOrder = []string{
			zerolog.TimestampFieldName,
			zerolog.LevelFieldName,
			zerolog.MessageFieldName,
			zerolog.CallerFieldName,
		}
		Logger = zerolog.New(output)
	} else {
		Logger = zerolog.New(os.Stdout)
	}

	log.Logger = Logger.Hook(severityHook{}).
		With().
		Str("app_id", config.AppID).
		Str("env", config.Env).
		Timestamp().
		Logger().
		Level(zerolog.Level(level))
}

// colorize returns the string s wrapped in ANSI code c, unless disabled is true.
func colorize(s interface{}, c int) string {
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", c, s)
}

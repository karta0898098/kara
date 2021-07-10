package logging

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoggerMessageOutput(t *testing.T) {
	var tests = []struct {
		name string
		msg  string
	}{
		{
			name: "Success",
			msg:  "test",
		},
	}
	for _, tt := range tests {
		b := new(bytes.Buffer)
		logger := SetupWithOption(
			WithOutput(io.MultiWriter(b, os.Stdout)),
			WithDebug(true),
		)
		logger.Info().Msg(tt.msg)
		assert.Contains(t, b.String(), tt.msg)
	}
}

func TestLoggerHasEnvOutput(t *testing.T) {
	var tests = []struct {
		name string
		env  string
	}{
		{
			name: "Success",
			env:  "local",
		},
	}

	for _, tt := range tests {
		b := new(bytes.Buffer)
		logger := SetupWithOption(
			WithOutput(io.MultiWriter(b, os.Stdout)),
			WithEnv(tt.env),
			WithDebug(true),
		)
		logger.Info().Msg("")
		assert.Contains(t, b.String(), tt.env)
	}
}

func TestLoggerLevel(t *testing.T) {
	var tests = []struct {
		name   string
		msg    string
		level  Level
		hasLog bool
	}{
		{
			name:   "HasLog",
			msg:    "msg",
			level:  DebugLevel,
			hasLog: true,
		},
		{
			name:   "NoLog",
			msg:    "msg",
			level:  ErrorLevel,
			hasLog: false,
		},
	}

	for _, tt := range tests {
		b := new(bytes.Buffer)
		logger := SetupWithOption(
			WithOutput(io.MultiWriter(b, os.Stdout)),
			WithLevel(tt.level),
			WithDebug(true),
		)
		logger.Info().Msg(tt.msg)
		if tt.hasLog {
			assert.Contains(t, b.String(), tt.msg)
		} else {
			assert.Equal(t, b.Len(), 0)
		}
	}
}

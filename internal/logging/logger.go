package logging

import (
	"fmt"
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level string

type Data = map[string]interface{}

const (
	Debug   Level = "DEBUG"
	Info    Level = "INFO"
	Warning Level = "WARN"
	Error   Level = "ERROR"
)

func LevelFromString(level string) (Level, error) {
	switch level {
	case string(Debug):
		return Debug, nil
	case string(Info):
		return Info, nil
	case string(Warning):
		return Warning, nil
	case string(Error):
		return Error, nil
	default:
		return "", fmt.Errorf(
			"expected valid Level type, instead got \"%s\"", level,
		)
	}
}

var LevelToZapLevel = map[Level]zapcore.Level{
	Debug:   zapcore.DebugLevel,
	Info:    zapcore.InfoLevel,
	Warning: zapcore.WarnLevel,
	Error:   zapcore.ErrorLevel,
}

const stderrThresholdLevel = zapcore.ErrorLevel

func newConsoleCore(thresholdLevel Level) (zapcore.Core, error) {
	level := LevelToZapLevel[thresholdLevel]

	consoleOutFilter := zap.LevelEnablerFunc(
		func(lvl zapcore.Level) bool {
			return lvl >= level && lvl < stderrThresholdLevel
		},
	)
	consoleErrFilter := zap.LevelEnablerFunc(
		func(lvl zapcore.Level) bool {
			return lvl >= level && lvl >= stderrThresholdLevel
		},
	)

	streamErr := zapcore.Lock(os.Stderr)
	streamOut := zapcore.Lock(os.Stdout)

	jsonEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())

	consoleCore := zapcore.NewTee(
		zapcore.NewCore(jsonEncoder, streamOut, consoleOutFilter),
		zapcore.NewCore(jsonEncoder, streamErr, consoleErrFilter),
	)

	return consoleCore, nil
}

type Logger struct {
	logger *zap.SugaredLogger
}

func newLogger(core zapcore.Core) Logger {
	logger := Logger{
		logger: zap.New(core).Sugar(),
	}

	return logger
}

func (l Logger) Debug(msg string, params Data) {
	l.logger.Debugw(msg, mapParamsIntoKeyValues(params)...)
}

func (l Logger) Info(msg string, params Data) {
	l.logger.Infow(msg, mapParamsIntoKeyValues(params)...)
}

func (l Logger) Warn(msg string, params Data) {
	l.logger.Warnw(msg, mapParamsIntoKeyValues(params)...)
}

func (l Logger) Error(msg string, params Data) {
	l.logger.Errorw(msg, mapParamsIntoKeyValues(params)...)
}

func (l Logger) Sync() {
	l.logger.Sync()
}

func mapParamsIntoKeyValues(params Data) []interface{} {
	count := len(params) * 2
	args := make([]interface{}, count)
	index := 0

	for key, value := range params {
		args[index] = key
		args[index+1] = value
		index += 2
	}

	return args
}

func NewLogger(level string) *Logger {
	levelFromString, err := LevelFromString(level)
	if err != nil {
		log.Fatal(err)
	}

	consoleCore, err := newConsoleCore(levelFromString)
	if err != nil {
		log.Fatal(err)
	}

	logger := newLogger(consoleCore)
	return &logger
}

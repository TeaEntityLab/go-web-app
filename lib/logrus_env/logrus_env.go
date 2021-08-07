/*
Package logrus_env exports env variables for configuring logrus standard logger.


* If you only use `logrus` package level logger, just import this package to `_` and let import side effect do the job.

* If you need extra `*logrus.Logger`(s), use the `NewLogger()` function.
*/
package logrus_env

import (
	"io"
	"os"
	"strings"

	logrus "github.com/sirupsen/logrus"
)

var (
	LOG_OUTPUT = strings.ToLower(strings.TrimSpace(os.Getenv("LOG_OUTPUT"))) // LOG_OUTPUT accept one of following values: stderr (default), stdout
	LOG_FORMAT = strings.ToLower(strings.TrimSpace(os.Getenv("LOG_FORMAT"))) // LOG_FORMAT accept one of following values: text (default), json
	LOG_LEVEL  = strings.ToLower(strings.TrimSpace(os.Getenv("LOG_LEVEL")))  // LOG_LEVEL accept one of following values: panic, fatal, error, warn, warning, info (default), debug
)

func init() {
	logOutput, logFormat, logLevel := parseEnv()
	logrus.SetOutput(logOutput)
	logrus.SetFormatter(logFormat)
	logrus.SetLevel(logLevel)
}

// NewLogger returns a `*logrus.Logger` that is configured with env variables
func NewLogger() *logrus.Logger {
	logOutput, logFormat, logLevel := parseEnv()
	return &logrus.Logger{
		Out:       logOutput,
		Hooks:     nil,
		Formatter: logFormat,
		Level:     logLevel,
	}
}

func parseEnv() (logOutput io.Writer, logFormat logrus.Formatter, logLevel logrus.Level) {

	// LOG_OUTPUT

	switch LOG_OUTPUT {
	case "stdout":
		logOutput = os.Stdout
	case "stderr":
		fallthrough
	default:
		// logrus.New() default
		logOutput = os.Stderr
	}

	// LOG_FORMAT

	switch LOG_FORMAT {
	case "json":
		logFormat = new(logrus.JSONFormatter)
	case "text":
		fallthrough
	default:
		// logrus.New() default
		logFormat = new(logrus.TextFormatter)
	}

	// LOG_LEVEL

	var parseErr error
	logLevel, parseErr = logrus.ParseLevel(LOG_LEVEL)
	if parseErr != nil {
		// logrus.New() default
		logLevel = logrus.InfoLevel
	}

	return logOutput, logFormat, logLevel
}

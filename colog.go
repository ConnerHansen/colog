package colog

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

type logWriter struct {
	Color string
	Key   string
}

const (
	debugLevel = 4
	infoLevel  = 3
	warnLevel  = 2
	errorLevel = 1
	fatalLevel = 0
)

var (
	loggingLevel = -1

	// LoggingFormat the format to use for the loggers
	LoggingFormat = "2006-01-02 15:04:05 -0700"

	// Debug the debug logger
	Debug *log.Logger

	// Error the error logger
	Error *log.Logger

	// Info the info logger
	Info *log.Logger

	// Warn the warn logger
	Warn *log.Logger

	// LogColorError the console code for the error color
	LogColorError = "\033[31m"

	// LogColorDebug the console code for the debug color
	LogColorDebug = "\033[36m"

	// LogColorWarn the console code for the warn color
	LogColorWarn = "\033[33m"

	// LogColorInfo the console code for the info color
	LogColorInfo = "\033[39m"

	// LogColorDefault the color to reset the console
	// to at the end of the line
	LogColorDefault = "\033[39m"

	// LogLevelDebug the debugging log level
	LogLevelDebug = "debug"

	// LogLevelError the error log level
	LogLevelError = "error"

	// LogLevelFatal the fatal log level
	LogLevelFatal = "fatal"

	// LogLevelInfo the info log level
	LogLevelInfo = "info"

	// LogLevelWarn the warn log level
	LogLevelWarn = "warn"

	// LogPrefixDebug the prefix to use for debugging lines
	LogPrefixDebug = " [D] "

	// LogPrefixError the prefix to use for error lines
	LogPrefixError = " [E] "

	// LogPrefixInfo the prefix to use for info lines
	LogPrefixInfo = " [I] "

	// LogPrefixWarn the prefix to use for warn lines
	LogPrefixWarn = " [W] "
)

// Write writes the contents of the logging message
func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Print(writer.Color + time.Now().Format(LoggingFormat) + writer.Key + string(bytes) + LogColorDefault)
}

// SetLoggingLevel sets the current logging level
func SetLoggingLevel(level string) {
	switch strings.ToLower(level) {
	case LogLevelDebug:
		loggingLevel = debugLevel
	case LogLevelInfo:
		loggingLevel = infoLevel
	case LogLevelWarn:
		loggingLevel = warnLevel
	case LogLevelError:
		loggingLevel = errorLevel
	case LogLevelFatal:
		loggingLevel = fatalLevel
	default:
		// Dunno what you're trying to say, so you get INFO
		loggingLevel = infoLevel
	}
}

// SetupLoggers does the initial setup for the loggers
func SetupLoggers() (*log.Logger, *log.Logger, *log.Logger, *log.Logger) {
	if loggingLevel < 0 {
		loggingLevel = infoLevel
	}

	debugWriter := logWriter{
		Color: LogColorDebug,
		Key:   LogPrefixDebug,
	}

	errWriter := logWriter{
		Color: LogColorError,
		Key:   LogPrefixError,
	}

	infoWriter := logWriter{
		Color: LogColorInfo,
		Key:   LogPrefixInfo,
	}

	warnWriter := logWriter{
		Color: LogColorWarn,
		Key:   LogPrefixWarn,
	}

	// Setup each logger, make sure we disregard any levels that we
	// want to shut off
	Debug = setupLogger(&debugWriter, os.Stdout, debugLevel)
	Error = setupLogger(&errWriter, os.Stderr, errorLevel)
	Info = setupLogger(&infoWriter, os.Stdout, infoLevel)
	Warn = setupLogger(&warnWriter, os.Stdout, warnLevel)

	return Debug, Info, Warn, Error
}

func setupLogger(logWriter *logWriter, writer io.Writer, cutoff int) *log.Logger {
	var logger *log.Logger

	if loggingLevel >= cutoff {
		logger = log.New(writer, "", 0)
		logger.SetOutput(logWriter)
	} else {
		logger = log.New(ioutil.Discard, "", 0)
	}

	return logger
}

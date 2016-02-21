package main

/**
 * File: log.go
 *
 * Define some simple logging infastrucutre that can be used from anywhere else in
 * the program. To use simply:
 *
 *    Logger().LEVEL.Println("my log message")
 *
 * Where LEVEL can be any one of:
 *    - Trace
 *    - Info
 *    - Warning
 *    - Error
 */

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

var _loggers MyLoggers
var once sync.Once

// MyLoggers is a simple log container for my application
type MyLoggers struct {
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
}

// Logger returns a custom set of loggers
func Logger() MyLoggers {
	once.Do(func() {
		initLoggers()
	})
	return _loggers
}

func initLoggers() {
	infoHandle := os.Stderr
	warningHandle := os.Stderr
	errorHandle := os.Stderr

	var traceHandle io.Writer
	var logFormat int

	if Config().TraceLogging {
		traceHandle = os.Stderr
		logFormat = log.Ldate | log.Ltime | log.Lshortfile
	} else {
		traceHandle = ioutil.Discard
		logFormat = log.Ldate | log.Ltime
	}

	_loggers = MyLoggers{
		Trace:   log.New(traceHandle, "TRACE: ", logFormat),
		Info:    log.New(infoHandle, "INFO: ", logFormat),
		Warning: log.New(warningHandle, "WARNING: ", logFormat),
		Error:   log.New(errorHandle, "ERROR: ", logFormat),
	}
}

// TempLogger - returns a temporrary logger
func TempLogger(level string) *log.Logger {
	return log.New(os.Stderr, level+": ", log.Ldate|log.Ltime)
}

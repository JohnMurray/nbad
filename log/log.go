// Package log makes logging within the system simple
package log

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

	"github.com/JohnMurray/nbad/config"
)

var _loggers myLoggers
var once sync.Once

type myLoggers struct {
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
}

func logger() myLoggers {
	once.Do(func() {
		initLoggers()
	})
	return _loggers
}

// Trace returns a trace-level logger
func Trace() *log.Logger {
	return logger().Trace
}

// Info returns an info-level logger
func Info() *log.Logger {
	return logger().Info
}

// Warning returns a warning-level logger
func Warning() *log.Logger {
	return logger().Warning
}

// Error returns a error-level logger
func Error() *log.Logger {
	return logger().Error
}

func initLoggers() {
	infoHandle := os.Stderr
	warningHandle := os.Stderr
	errorHandle := os.Stderr

	var traceHandle io.Writer
	var logFormat int

	if config.TraceLogging() {
		traceHandle = os.Stderr
		logFormat = log.Ldate | log.Ltime | log.Lshortfile
	} else {
		traceHandle = ioutil.Discard
		logFormat = log.Ldate | log.Ltime
	}

	_loggers = myLoggers{
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

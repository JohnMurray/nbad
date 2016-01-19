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
	"io/ioutil"
	"log"
	"os"
	"sync"
)

var _loggers MyLoggers
var once sync.Once

type MyLoggers struct {
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
}

func Logger() MyLoggers {
	once.Do(func() {
		initLoggers()
	})
	return _loggers
}

func initLoggers() {
	traceHandle := ioutil.Discard
	infoHandle := os.Stderr
	warningHandle := os.Stderr
	errorHandle := os.Stderr

	_loggers = MyLoggers{
		Trace:   log.New(traceHandle, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile),
		Info:    log.New(infoHandle, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		Warning: log.New(warningHandle, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile),
		Error:   log.New(errorHandle, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

// Package logger is logger module for fxoss
package logger

import (
	"os"
	"log"
	"io/ioutil"
)

var (
	traceLogger *log.Logger
	debugLogger *log.Logger
)

func init(){
	traceLogger = log.New(ioutil.Discard, "[trace]: ", log.Ldate|log.Ltime|log.Lshortfile)
	debugLogger = log.New(os.Stdout, "[debug] : ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Mylogger(debug bool) *log.Logger {
	if debug {
		return debugLogger
	} else {
		return traceLogger
	}
}
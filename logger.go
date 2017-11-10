package main

import (
	"log"
	"os"
)

type nullWriter struct{}

func (nullWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

type Logger struct {
	Debug *log.Logger
	Info  *log.Logger
	Warn  *log.Logger
	Fatal *log.Logger
}

var logger *Logger

func InitLogger() {
	logger = NewLogger()
}

func DisableDebugLogging() {
	nw := nullWriter{}
	logger.Debug.SetOutput(nw)
}

func InitNullLogger() {
	nw := nullWriter{}
	InitLogger()
	logger.Debug.SetOutput(nw)
	logger.Info.SetOutput(nw)
	logger.Warn.SetOutput(nw)
	logger.Fatal.SetOutput(nw)
}

func NewLogger() *Logger {
	return &Logger{
		Debug: log.New(os.Stdout, "[debug] ", log.Ldate|log.Ltime|log.Lshortfile),
		Info:  log.New(os.Stdout, "[info] ", log.Ldate|log.Ltime|log.Lshortfile),
		Warn:  log.New(os.Stdout, "[warn] ", log.Ldate|log.Ltime|log.Lshortfile),
		Fatal: log.New(os.Stderr, "[fatal] ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

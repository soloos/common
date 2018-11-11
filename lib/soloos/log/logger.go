package log

import (
	"log"
	"os"
)

const (
	calldepth = 3
)

const (
	LDebug = iota
	LInfo
	LWarn
	LError
	LFatal
)

const (
	Ldate         = log.Ldate
	Ltime         = log.Ltime
	Lmicroseconds = log.Lmicroseconds
	Llongfile     = log.Llongfile
	Lshortfile    = log.Lshortfile
	LUTC          = log.LUTC
	LstdFlags     = Ldate | Ltime
)

var l Logger = &defaultLogger{LInfo, log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)}

type Logger interface {
	SetLevel(level int)
	SetFlags(flag int)

	Debug(v ...interface{})
	Debugf(format string, v ...interface{})

	Info(v ...interface{})
	Infof(format string, v ...interface{})

	Warn(v ...interface{})
	Warnf(format string, v ...interface{})

	Error(v ...interface{})
	Errorf(format string, v ...interface{})

	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})

	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
}

func SetLogger(logger Logger) {
	l = logger
}

func SetLevel(level int) {
	l.SetLevel(level)
}

func SetFlags(flag int) {
	l.SetFlags(flag)
}

func SetDummyLogger() {
	l = &dummyLogger{}
}

func Debug(v ...interface{}) {
	l.Debug(v...)
}
func Debugf(format string, v ...interface{}) {
	l.Debugf(format, v...)
}

func Info(v ...interface{}) {
	l.Info(v...)
}
func Infof(format string, v ...interface{}) {
	l.Infof(format, v...)
}

func Warn(v ...interface{}) {
	l.Warn(v...)
}
func Warnf(format string, v ...interface{}) {
	l.Warnf(format, v...)
}

func Error(v ...interface{}) {
	l.Error(v...)
}
func Errorf(format string, v ...interface{}) {
	l.Errorf(format, v...)
}

func Fatal(v ...interface{}) {
	l.Fatal(v...)
}
func Fatalf(format string, v ...interface{}) {
	l.Fatalf(format, v...)
}

func Panic(v ...interface{}) {
	l.Panic(v...)
}
func Panicf(format string, v ...interface{}) {
	l.Panicf(format, v...)
}

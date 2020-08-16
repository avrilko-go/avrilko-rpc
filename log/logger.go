package log

import (
	"log"
	"os"
)

const (
	CallDepth = 3 // 日志堆栈
)

var l = NewDefaultLogger()

type Logger interface {
	Debug(v ...interface{})
	DebugF(format string, v ...interface{})

	Info(v ...interface{})
	InfoF(format string, v ...interface{})

	Warn(v ...interface{})
	WarnF(format string, v ...interface{})

	Error(v ...interface{})
	ErrorF(format string, v ...interface{})

	Fatal(v ...interface{})
	FatalF(format string, v ...interface{})

	Panic(v ...interface{})
	PanicF(format string, v ...interface{})
}

type Handler interface {
	Handle(v ...interface{})
}

func NewDefaultLogger() Logger {
	return &defaultLogger{log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)}
}

// 自定义日志
func SetLogger(logger Logger) {
	l = logger
}

func Debug(v ...interface{}) {
	l.Debug(v...)
}
func DebugF(format string, v ...interface{}) {
	l.DebugF(format, v...)
}

func Info(v ...interface{}) {
	l.Info(v...)
}
func InfoF(format string, v ...interface{}) {
	l.InfoF(format, v...)
}

func Warn(v ...interface{}) {
	l.Warn(v...)
}
func WarnF(format string, v ...interface{}) {
	l.WarnF(format, v...)
}

func Error(v ...interface{}) {
	l.Error(v...)
}
func ErrorF(format string, v ...interface{}) {
	l.ErrorF(format, v...)
}

func Fatal(v ...interface{}) {
	l.Fatal(v...)
}
func FatalF(format string, v ...interface{}) {
	l.FatalF(format, v...)
}

func Panic(v ...interface{}) {
	l.Panic(v...)
}
func PanicF(format string, v ...interface{}) {
	l.PanicF(format, v...)
}

func Handle(v ...interface{}) {
	if handle, ok := l.(Handler); ok {
		handle.Handle(v)
	}
}

package log

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

const (
	FATAL = iota
	ERROR
	WARN
	INFO
	DEBUG
)

const (
	STDOUT = iota
	FILE
	EMAIL
)

type loggerType int

var adapter []func() LoggerInterface = make([]func() LoggerInterface, 5)

func register(lt loggerType, construct func() LoggerInterface) {
	adapter[lt] = construct
}

type LoggerInterface interface {
	Init(config string) error
	WriteMsg(level int, msg string) error
	Flush()
	Destroy()
}

type Logger struct {
	level     int
	callDepth int
	lg        LoggerInterface
}

func NewLogger(lt loggerType) *Logger {
	ret := new(Logger)
	construct := adapter[lt]
	lg := construct()
	ret.lg = lg
	ret.level = DEBUG
	ret.callDepth = 3

	return ret
}

func (self *Logger) SetLevel(level int) {
	self.level = level
}

func (self *Logger) SetLogger(config string) {
	self.lg.Init(config)
}

func (self *Logger) DeleteLogger() {
	self.lg.Destroy()
}

func (self *Logger) writeMsg(level int, msg string) {
	_, file, line, ok := runtime.Caller(self.callDepth)
	if !ok {
		file = "???"
		line = 0
	}
	msg = fmt.Sprintf("file %s line %d %s", file, line, msg)
	self.lg.WriteMsg(level, msg)
}

func (self *Logger) Fatal(format string, v ...interface{}) {
	if FATAL > self.level {
		return
	}
	msg := fmt.Sprintf("[FATAL] "+format, v...)
	self.writeMsg(FATAL, msg)
	time.Sleep(100 * time.Millisecond)
	os.Exit(1)
}

func (self *Logger) Error(format string, v ...interface{}) {
	if ERROR > self.level {
		return
	}
	msg := fmt.Sprintf("[ERROR] "+format, v...)
	self.writeMsg(ERROR, msg)
}

func (self *Logger) Warn(format string, v ...interface{}) {
	if WARN > self.level {
		return
	}
	msg := fmt.Sprintf("[WARN] "+format, v...)
	self.writeMsg(WARN, msg)
}

func (self *Logger) Info(format string, v ...interface{}) {
	if INFO > self.level {
		return
	}
	msg := fmt.Sprintf("[INFO] "+format, v...)
	self.writeMsg(INFO, msg)
}

func (self *Logger) Debug(format string, v ...interface{}) {
	if DEBUG > self.level {
		return
	}
	msg := fmt.Sprintf("[DEBUG] "+format, v...)
	self.writeMsg(DEBUG, msg)
}

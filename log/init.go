package log

import (
	"github.com/BigTong/common/file"
)

var stdLogger *Logger
var fileLogger *Logger
var emailLogger *Logger

var fileLoggerInited bool = false

func init() {
	stdLogger = NewLogger(STDOUT)
	fileLogger = NewLogger(FILE)
	emailLogger = NewLogger(EMAIL)
}

func InitFileLoggerFromConfigFile(fileName string, level int) {
	config := file.ReadFileToString(fileName)
	fileLogger.SetLogger(config)
	fileLogger.SetLevel(level)
	fileLoggerInited = true
}

func InitFileLoggerWithDefaultConfig() {
	fileLogger.SetLogger("")
	fileLogger.SetLevel(DEBUG)
	fileLoggerInited = true
}

func GetStdLogger() *Logger {
	return stdLogger
}

func GetFileLogger() *Logger {
	return fileLogger
}

func GetEmailLogger() *Logger {
	return emailLogger
}

func SInfo(format string, v ...interface{}) {
	stdLogger.Info(format, v...)
}
func SFatal(format string, v ...interface{}) {
	stdLogger.Fatal(format, v...)
}
func SError(format string, v ...interface{}) {
	stdLogger.Error(format, v...)
}
func SWarn(format string, v ...interface{}) {
	stdLogger.Warn(format, v...)
}
func SDebug(format string, v ...interface{}) {
	stdLogger.Debug(format, v...)
}

func FInfo(format string, v ...interface{}) {
	if !fileLoggerInited {
		return
	}
	fileLogger.Info(format, v...)
}
func FFatal(format string, v ...interface{}) {
	if !fileLoggerInited {
		return
	}
	fileLogger.Fatal(format, v...)
}
func FError(format string, v ...interface{}) {
	if !fileLoggerInited {
		return
	}
	fileLogger.Error(format, v...)
}
func FWarn(format string, v ...interface{}) {
	if !fileLoggerInited {
		return
	}
	fileLogger.Warn(format, v...)
}
func FDebug(format string, v ...interface{}) {
	if !fileLoggerInited {
		return
	}
	fileLogger.Debug(format, v...)
}

func EInfo(format string, v ...interface{}) {
	emailLogger.Info(format, v...)
}
func EFatal(format string, v ...interface{}) {
	emailLogger.Fatal(format, v...)
}
func EError(format string, v ...interface{}) {
	emailLogger.Error(format, v...)
}
func EWarn(format string, v ...interface{}) {
	emailLogger.Warn(format, v...)
}
func EDebug(format string, v ...interface{}) {
	emailLogger.Debug(format, v...)
}

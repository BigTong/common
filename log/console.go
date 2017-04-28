package log

import (
	"log"
	"os"
)

type Brush func(string) string

func NewBrush(color string) Brush {
	prefix := "\033["
	suffix := "\033[0m"
	return func(text string) string {
		return prefix + color + text + suffix
	}
}

var color = []Brush{
	NewBrush("0;35m"),
	NewBrush("0;31m"),
	NewBrush("0;33m"),
	NewBrush("0;37m"),
	NewBrush("0;34m"),
}

type ConsoleWriter struct {
	lg *log.Logger
}

func NewConsoleWriter() LoggerInterface {
	return &ConsoleWriter{
		lg: log.New(os.Stdout, "", log.LstdFlags),
	}
}

// don need to init console log, default use log
func (self *ConsoleWriter) Init(config string) error { return nil }

func (self *ConsoleWriter) WriteMsg(level int, msg string) error {
	self.lg.Println(color[level](msg))
	return nil
}

func (self *ConsoleWriter) Flush() {}

func (self *ConsoleWriter) Destroy() {}

func init() {
	register(STDOUT, NewConsoleWriter)
}

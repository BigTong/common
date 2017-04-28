package log

import (
	"fmt"
	"testing"
	"time"
)

func TestConsole(t *testing.T) {
	logger := NewLogger(STDOUT)
	logger.Debug("debug")
	logger.Info("info")
	logger.Warn("warn")
	logger.Error("error")
	//logger.Fatal("fatal")

	fmt.Println("test level")
	logger.SetLevel(INFO)
	logger.Debug("debug")
	logger.Info("info")
	logger.Warn("warn")
	logger.Error("error")

	fmt.Println("test fatal")
	logger.Fatal("fatal")
}

func TestFile(t *testing.T) {
	logger := NewLogger(FILE)
	logger.SetLogger(`{"file_name":"test.log","max_line_num":100000, "file_path":"/data/log/go/log_test/", "max_file_num":5}`)
	for i := 0; i < 1100000; i++ {
		logger.Info("%v", i)
	}
	//time.Sleep(1 * time.Second)
}

func TestEmail(t *testing.T) {
	logger := NewLogger(EMAIL)
	config := `
	{
		"username": "liuliu@yidian-inc.com",
		"password": "",
		"host": "10.111.0.15:25",
		"subject": "email logger test",
		"send_to": "liuliu@yidian-inc.com;zhangrentong@yidian-inc.com",
		"send_frequency": 1
	}
	`
	logger.SetLogger(config)
	logger.Debug("test email logger")
	logger.Info("test email logger")
	logger.Warn("test email logger")
	logger.Error("test email logger")
	time.Sleep(2 * time.Minute)
}

func TestSingleton(t *testing.T) {
	logger := GetFileLogger()
	logger.SetLogger(`{"file_name":"test.log","max_line_num":100000, "file_path":"/data/log/go/log_test/", "max_file_num":5}`)
	for i := 0; i < 1100000; i++ {
		logger.Info("%v", i)
	}
}

func TestLazyMode(t *testing.T) {
	/*
	 *	this need to init logger config in your project service entry
	 */
	logger := GetFileLogger()
	logger.SetLogger(`{"file_name":"test.log","max_line_num":100000, "file_path":"/data/log/go/log_test/", "max_file_num":5}`)

	for i := 0; i < 1100000; i++ {
		FDebug("%v", i)
		FInfo("%v", i)
	}
}

package log

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	kMAX_LINE_NUM = 1000000
	kMAX_FILE_NUM = 15
	kFILE_NAME    = "default.log"
	kFILE_PATH    = "./logs/"
	kBUF_SIZE     = 100
	kLineBufSize  = 1
	kTimeMax      = 99999999999999999
	kRotate       = false
)

const (
	defaultConfig = `{
		"file_name": "log4go.log",
		"file_path": "./logs/",
		"max_line_num": 1000000,
		"max_file_num": 30,
		"rotate": true
	}`
)

type FileWriter struct {
	lg *log.Logger
	mw *MuxWriter

	FileName    string `json:"file_name"`
	FilePath    string `json:"file_path"`
	MaxLineNum  int    `json:"max_line_num"`
	MaxFileNum  int    `json:"max_file_num"`
	Rotate      bool   `json:"rotate"`
	currentLine int
	writeBuf    chan string
	writeFlag   bool
}

//TODO: need to support multi consumer thread to consume writeBuf
type MuxWriter struct {
	lock *sync.RWMutex
	fd   *os.File
}

func (self *MuxWriter) Write(b []byte) (int, error) {
	self.lock.Lock()
	defer self.lock.Unlock()
	return self.fd.Write(b)
}

func NewFileWriter() LoggerInterface {
	return new(FileWriter)
}

func (self *FileWriter) setCurrentLine() {
	r := self.mw.fd
	buf := make([]byte, 8196)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		if err != nil && err != io.EOF {
			self.currentLine = count
			print(err.Error())
		}

		count += bytes.Count(buf[:c], lineSep)

		if err == io.EOF {
			break
		}
	}

	self.currentLine = count
}

func (self *FileWriter) createNewFile() (*os.File, error) {
	//fd, err := os.OpenFile(self.FileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	os.MkdirAll(self.FilePath, 0775)
	name := self.FilePath + self.FileName
	fd, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0660)
	return fd, err
}

func (self *FileWriter) Init(config string) error {
	if len(config) == 0 {
		config = defaultConfig
	}

	json.Unmarshal([]byte(config), self)

	self.writeBuf = make(chan string, 1000)
	self.mw = new(MuxWriter)
	self.mw.lock = &sync.RWMutex{}
	var err error
	self.mw.fd, err = self.createNewFile()
	if err != nil {
		panic("create new log file error")
	}
	self.lg = log.New(self.mw, "", 0)
	self.setCurrentLine()
	self.writeFlag = true

	go self.write()
	return nil
}

func (self *FileWriter) closeFile() {
	self.mw.fd.Close()
}

func (self *FileWriter) check() {
	remain := self.MaxLineNum - self.currentLine
	if remain == 0 {
		now := time.Now().Format("20060102150405")
		self.closeFile()
		name := self.FilePath + self.FileName
		os.Rename(name, name+"."+now)
		self.mw.fd, _ = self.createNewFile()
		self.currentLine = 0
	}

	if self.Rotate {
		self.rotate()
	}
}

func (self *FileWriter) rotate() {
	dir, _ := ioutil.ReadDir(self.FilePath)
	logs := []string{}
	for _, file := range dir {
		names := strings.Split(self.mw.fd.Name(), "/")
		name := names[len(names)-1]
		if strings.Contains(file.Name(), name) {
			logs = append(logs, file.Name())
		}
	}
	if len(logs) <= self.MaxFileNum {
		return
	}
	min := kTimeMax
	removeFileName := ""
	for _, l := range logs {
		str := strings.Split(l, ".")
		timestamp := str[len(str)-1]
		unix, err := strconv.Atoi(timestamp)
		if err != nil {
			continue
		}
		if min > unix {
			min = unix
			removeFileName = l
		}
	}
	os.Remove(self.FilePath + removeFileName)
}

func (self *FileWriter) getHeader() string {
	return time.Now().Format("2006/01/02 15:04:05")
}

func (self *FileWriter) write() {
	for {
		self.check()
		msg := <-self.writeBuf
		txt := self.getHeader() + " " + msg
		self.lg.Println(txt)
		self.currentLine++
	}
}

func (self *FileWriter) Sync() {
	for {
		if len(self.writeBuf) != 0 {
			time.Sleep(200 * time.Millisecond)
		} else {
			break
		}
	}
}

func (self *FileWriter) WriteMsg(level int, msg string) error {
	if self.writeFlag {
		self.writeBuf <- msg
	}
	return nil
}

func (self *FileWriter) Flush() {
	self.Sync()
	self.mw.fd.Sync()
}

func (self *FileWriter) Destroy() {
	self.writeFlag = false
	self.Flush()
	self.closeFile()
}

func init() {
	register(FILE, NewFileWriter)
}

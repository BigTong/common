package log

import (
	"encoding/json"
	"fmt"
	"net/smtp"
	"strings"
	"time"
)

const (
	kSendFrequency = 6
	kFormat        = "html"
)

type EmailWriter struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	Host          string `json:"host"`
	Subject       string `json:"subject"`
	SendTo        string `json:"send_to"`
	SendFrequency int    `json:"send_frequency"`
	writeFlag     bool
	msgBuf        []string
	auth          smtp.Auth
	ticker        *time.Ticker
}

func NewEmailWriter() LoggerInterface {
	ret := &EmailWriter{
		SendFrequency: kSendFrequency,
		ticker:        time.NewTicker(time.Duration(kSendFrequency) * time.Hour),
		writeFlag:     true,
	}

	go ret.checkBuf()
	return ret
}

func (self *EmailWriter) checkBuf() {
	for {
		select {
		case <-self.ticker.C:
			data := ""
			length := len(self.msgBuf)
			if length == 0 {
				continue
			}
			for _, msg := range self.msgBuf {
				data = data + msg + "</br>"
			}
			go self.sendMail(data)
			self.msgBuf = self.msgBuf[length:]
		}
	}
}

func (self *EmailWriter) Init(config string) error {
	err := json.Unmarshal([]byte(config), self)
	if err != nil {
		return err
	}
	self.setSmtpAuth()
	self.ticker = time.NewTicker(time.Duration(self.SendFrequency) * time.Hour)

	//this frequency just for test
	//self.ticker = time.NewTicker(time.Duration(self.SendFrequency) * time.Minute)
	return nil
}

func (self *EmailWriter) setSmtpAuth() {
	self.auth = smtp.PlainAuth("", self.Username, self.Password,
		strings.Split(self.Host, ":")[0])
}

func (self *EmailWriter) sendMail(msg string) {
	contentType := "Content-Type: text/" + kFormat + "; charset=UTF-8"
	body := []byte("To: " + self.SendTo + "\r\nFrom: " + self.Username +
		"<" + self.Username + ">\r\nSubject: " + self.Subject + "\r\n" +
		contentType + "\r\n\r\n" + msg)
	smtp.SendMail(self.Host, self.auth, self.Username,
		strings.Split(self.SendTo, ";"), body)
}

func (self *EmailWriter) WriteMsg(level int, msg string) error {
	msg = fmt.Sprintf("%s %s", time.Now().Format("2006-01-02 15:04:05"), msg)
	if self.writeFlag {
		self.msgBuf = append(self.msgBuf, msg)
	}
	return nil
}

func (self *EmailWriter) Flush() {
	//self.ticker = time.NewTicker(5 * time.Second)
	return
}

func (self *EmailWriter) Destroy() {
	self.writeFlag = false
	self.Flush()
}

func init() {
	register(EMAIL, NewEmailWriter)
}

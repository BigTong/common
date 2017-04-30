package encodings

import (
	"strings"

	"code.google.com/p/mahonia"
	"github.com/saintfish/chardet"
)

type Utf8Converter struct {
	detector *chardet.Detector
	gbk      mahonia.Decoder
	big5     mahonia.Decoder
}

func NewUtf8Converter() *Utf8Converter {
	return &Utf8Converter{
		detector: chardet.NewHtmlDetector(),
		gbk:      mahonia.NewDecoder("gb18030"),
		big5:     mahonia.NewDecoder("big5"),
	}
}

func (self *Utf8Converter) DetectCharset(src []byte) string {
	rets, err := self.detector.DetectAll(src)
	if err != nil {
		return ""
	}
	maxret := ""
	w := 0
	for _, ret := range rets {
		cs := strings.ToLower(ret.Charset)
		if strings.HasPrefix(cs, "gb") ||
			strings.HasPrefix(cs, "utf") {
			if w < ret.Confidence {
				w = ret.Confidence
				maxret = cs
			}
		}
		continue
	}
	return maxret
}

func (self *Utf8Converter) ToUTF8(src []byte) []byte {
	charset := self.DetectCharset(src)

	if !strings.HasPrefix(charset, "gb") &&
		!strings.HasPrefix(charset, "big") {
		charset = "utf-8"
	}
	switch charset {
	case "utf-8", "utf8":
		return src
	case "gb2312", "gb-2312", "gbk", "gb18030", "gb-18030":
		ret, ok := self.gbk.ConvertStringOK(string(src))
		if ok {
			return []byte(ret)
		}
	case "big5":
		ret, ok := self.big5.ConvertStringOK(string(src))
		if ok {
			return []byte(ret)
		}
	}
	return src
}

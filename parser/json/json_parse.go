package json_parser

import (
	"bytes"
	"io"

	"common/dlog"
	"common/parser"
	"github.com/bitly/go-simplejson"
)

const (
	kJsonFakeRoot = "json_fake_path"
	kMapToArray   = "map_to_array"
)

type JsonParser struct {
	jsonDoc *simplejson.Json
}

func NewJsonParserFromBytes(src []byte) *JsonParser {
	src = bytes.TrimPrefix(src, []byte("\xef\xbb\xbf"))
	r := bytes.NewReader(src)
	return NewJsonParser(r)
}

func NewJsonParser(r io.Reader) *JsonParser {
	jsonDoc, err := simplejson.NewFromReader(r)
	if err != nil {
		dlog.Error("new json parser get error:%s", err.Error())
		return nil
	}

	return &JsonParser{
		jsonDoc: jsonDoc,
	}
}

func (self *JsonParser) GetJsonDoc() *simplejson.Json {
	return self.jsonDoc
}

func (self *JsonParser) GetInt(s *parser.Selector) int {
	return GetIntValFromJsonDoc(s, self.jsonDoc)
}

func (self *JsonParser) GetInt64(s *parser.Selector) int64 {
	return GetInt64ValFromJsonDoc(s, self.jsonDoc)
}

func (self *JsonParser) GetFloat64(s *parser.Selector) float64 {
	return GetFloat64ValFromJsonDoc(s, self.jsonDoc)
}

func (self *JsonParser) GetString(s *parser.Selector) string {
	return GetStringValFromJsonDoc(s, self.jsonDoc)
}

func (self *JsonParser) GetJsonObjects(branch []string) []*simplejson.Json {
	return GetJsonObjectsFromDoc(branch, self.jsonDoc)
}

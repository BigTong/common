package terminator

import (
	"encoding/json"
	"strings"

	"common"
	"common/dlog"
	"common/parser"
	"common/parser/html"
	"common/parser/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/bitly/go-simplejson"
)

type LazyParser struct {
}

func NewLazyParser() *LazyParser {
	return &LazyParser{}
}

func (self *LazyParser) ParseFromStringConfig(src,
	conf string) (map[string]interface{}, error) {
	config := make(map[string]interface{})
	err := json.Unmarshal([]byte(conf), &config)
	if err != nil {
		return nil, err
	}
	return self.Parse(src, config), nil
}

func (self *LazyParser) ParseById(src, id string,
	config map[string]interface{}) map[string]interface{} {
	if config == nil {
		return nil
	}
	val, ok := config[id]
	if !ok {
		dlog.Error("not get config:%s", id)
		return nil
	}

	conf, err := InterfaceToMap(val)
	if err != nil {
		dlog.Error("interface to map get error:%s", err.Error())
		return nil
	}

	return self.Parse(src, conf)

}
func (self *LazyParser) Parse(src string,
	config map[string]interface{}) map[string]interface{} {
	input, ok := config[parser.KInputType]
	if !ok {
		dlog.Error("not set input type")
		return nil
	}
	//delete(config, parser.KInputType)

	inputType, err := MapToInputConfig(input)
	if err != nil {
		dlog.Error("not valid input type:%v, error:%s", input, err.Error())
		return nil
	}

	if len(inputType.Extractor) != 0 {
		_, src = common.StringExtract(src, inputType.Extractor)
	}

	switch inputType.ContentType {
	case parser.KContentJson:
		js := json_parser.NewJsonParser(strings.NewReader(src))
		return self.innerParseJson(js.GetJsonDoc(), config)
	case parser.KContentHtml:
		html := html_parser.NewHtmlParser(strings.NewReader(src))
		return self.innerParseHtml(html.GetData(), config)
	case parser.KContentString:
	default:
		dlog.Error("unknown input type")
		return nil
	}
	dlog.Warn("not supported input type:%s", inputType.ContentType)
	return nil
}

func (self *LazyParser) innerParseHtml(htmlDoc *goquery.Selection,
	config map[string]interface{}) map[string]interface{} {
	if config == nil {
		return nil
	}

	ret := make(map[string]interface{})
	for k, v := range config {
		if k == parser.KInputType {
			continue
		}

		elemConfig, err := MapToElemConfig(v)
		if err != nil {
			dlog.Warn("not valid elemconfig type:%v err:%s", v, err.Error())
			continue
		}

		if elemConfig.Selector != nil {
			switch elemConfig.Type {
			case parser.KInt:
				ret[k] = html_parser.ParseInt(htmlDoc, elemConfig.Selector)
			case parser.KFloat:
				ret[k] = html_parser.ParseFloat64(htmlDoc, elemConfig.Selector)
			default:
				ret[k] = html_parser.ParseString(htmlDoc, elemConfig.Selector)
			}
			continue
		}
		if elemConfig.ObjectConfig == nil {
			continue
		}

		if elemConfig.ObjectConfig.RootPath != nil {
			rootPath := elemConfig.ObjectConfig.RootPath
			htmlDocs := html_parser.GetSelectionFromDoc(htmlDoc, rootPath)
			arr := []map[string]interface{}{}
			for _, doc := range htmlDocs {
				arr = append(arr, self.innerParseHtml(doc,
					elemConfig.ObjectConfig.Fields))
			}
			ret[k] = arr
			continue
		}

		ret[k] = self.innerParseHtml(htmlDoc, elemConfig.ObjectConfig.Fields)
	}
	return ret
}

func (self *LazyParser) innerParseJson(jsonDoc *simplejson.Json,
	config map[string]interface{}) map[string]interface{} {
	if config == nil {
		return nil
	}

	ret := make(map[string]interface{})
	for k, v := range config {
		if k == parser.KInputType {
			continue
		}

		elemConfig, err := MapToElemConfig(v)
		if err != nil {
			dlog.Warn("not valid elem config type:%v err:%s", v, err.Error())
			continue
		}

		if elemConfig.Selector != nil {
			switch elemConfig.Type {
			case parser.KInt:
				ret[k] = json_parser.GetIntValFromJsonDoc(elemConfig.Selector,
					jsonDoc)
			case parser.KFloat:
				ret[k] = json_parser.GetFloat64ValFromJsonDoc(elemConfig.Selector,
					jsonDoc)
			default:
				ret[k] = json_parser.GetStringValFromJsonDoc(elemConfig.Selector,
					jsonDoc)
			}
			continue
		}

		if elemConfig.ObjectConfig == nil {
			continue
		}

		if elemConfig.ObjectConfig.RootPath != nil {
			branch := strings.Split(elemConfig.ObjectConfig.RootPath.JsonKey, " ")
			jsonDocs := json_parser.GetJsonObjectsFromDoc(branch, jsonDoc)
			arr := []map[string]interface{}{}
			for _, doc := range jsonDocs {
				arr = append(arr, self.innerParseJson(doc,
					elemConfig.ObjectConfig.Fields))
			}
			ret[k] = arr
			continue
		}

		ret[k] = self.innerParseJson(jsonDoc, elemConfig.ObjectConfig.Fields)
	}
	return ret
}

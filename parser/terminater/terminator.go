package terminater

import (
	"encoding/json"
	"strings"

	"github.com/BigTong/common/parser"
	"github.com/BigTong/common/parser/html"
	cjson "github.com/BigTong/common/parser/json"
	cstrings "github.com/BigTong/common/strings"
	"github.com/PuerkitoBio/goquery"
	"github.com/bitly/go-simplejson"
)

type Terminater struct {
}

func NewTerminater() *Terminater {
	return &Terminater{}
}

func (self *Terminater) ParseFromStringConfig(src,
	conf string) (map[string]interface{}, error) {
	config := make(map[string]interface{})
	err := json.Unmarshal([]byte(conf), &config)
	if err != nil {
		return nil, err
	}
	return self.Parse(src, config), nil
}

func (self *Terminater) ParseById(src, id string,
	config map[string]interface{}) map[string]interface{} {
	if config == nil {
		return nil
	}
	val, ok := config[id]
	if !ok {
		return nil
	}

	conf, err := InterfaceToMap(val)
	if err != nil {
		return nil
	}

	return self.Parse(src, conf)

}
func (self *Terminater) Parse(src string,
	config map[string]interface{}) map[string]interface{} {
	input, ok := config[parser.INPUT_TYPE]
	if !ok {
		return nil
	}
	//delete(config, parser.KInputType)

	inputType, err := MapToInputConfig(input)
	if err != nil {
		return nil
	}

	if len(inputType.Extractor) != 0 {
		_, src = cstrings.StringExtract(src, inputType.Extractor)
	}

	switch inputType.ContentType {
	case parser.CONTENT_TYPE_JSON:
		js := cjson.NewJsonParser(strings.NewReader(src))
		return self.innerParseJson(js.GetJsonDoc(), config)
	case parser.CONTENT_TYPE_HTML:
		html := html.NewHtmlParser(strings.NewReader(src))
		return self.innerParseHtml(html.GetData(), config)
	case parser.CONTENT_TYPE_STRING:
	default:
		return nil
	}
	return nil
}

func (self *Terminater) innerParseHtml(htmlDoc *goquery.Selection,
	config map[string]interface{}) map[string]interface{} {
	if config == nil {
		return nil
	}

	ret := make(map[string]interface{})
	for k, v := range config {
		if k == parser.INPUT_TYPE {
			continue
		}

		elemConfig, err := MapToElemConfig(v)
		if err != nil {
			continue
		}

		if elemConfig.Selector != nil {
			switch elemConfig.Type {
			case parser.TYPE_INT:
				ret[k] = html.ParseInt(htmlDoc, elemConfig.Selector)
			case parser.TYPE_FLOAT:
				ret[k] = html.ParseFloat64(htmlDoc, elemConfig.Selector)
			default:
				ret[k] = html.ParseString(htmlDoc, elemConfig.Selector)
			}
			continue
		}
		if elemConfig.ObjectConfig == nil {
			continue
		}

		if elemConfig.ObjectConfig.RootPath != nil {
			rootPath := elemConfig.ObjectConfig.RootPath
			htmlDocs := html.GetSelectionFromDoc(htmlDoc, rootPath)
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

func (self *Terminater) innerParseJson(jsonDoc *simplejson.Json,
	config map[string]interface{}) map[string]interface{} {
	if config == nil {
		return nil
	}

	ret := make(map[string]interface{})
	for k, v := range config {
		if k == parser.INPUT_TYPE {
			continue
		}

		elemConfig, err := MapToElemConfig(v)
		if err != nil {
			continue
		}

		if elemConfig.Selector != nil {
			switch elemConfig.Type {
			case parser.TYPE_INT:
				ret[k] = cjson.GetIntValFromJsonDoc(elemConfig.Selector,
					jsonDoc)
			case parser.TYPE_FLOAT:
				ret[k] = cjson.GetFloat64ValFromJsonDoc(elemConfig.Selector,
					jsonDoc)
			default:
				ret[k] = cjson.GetStringValFromJsonDoc(elemConfig.Selector,
					jsonDoc)
			}
			continue
		}

		if elemConfig.ObjectConfig == nil {
			continue
		}

		if elemConfig.ObjectConfig.RootPath != nil {
			branch := strings.Split(elemConfig.ObjectConfig.RootPath.JsonKey, " ")
			jsonDocs := cjson.GetJsonObjectsFromDoc(branch, jsonDoc)
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

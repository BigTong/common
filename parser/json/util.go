package json_parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"common"
	"common/dlog"
	"common/parser"
	"github.com/bitly/go-simplejson"
)

func keyIndex(buf string) (string, string) {
	kv := strings.Split(buf, "[")
	if len(kv) == 1 {
		return kv[0], ""
	} else if len(kv) == 2 {
		n := strings.Trim(kv[1], "]")
		return kv[0], n
	}
	return "", ""
}

func MapToArray(js *simplejson.Json) []*simplejson.Json {
	ret := []*simplejson.Json{}
	if js == nil {
		return ret
	}
	values, _ := js.Map()
	for k, v := range values {
		data, _ := json.Marshal(v)
		tmpJson, err := simplejson.NewFromReader(bytes.NewReader(data))
		if err != nil {
			dlog.Error("unmashal json get error:%s,%s", err.Error(), data)
			continue
		}
		tmpJson.Set("_key", k)
		ret = append(ret, tmpJson)
	}
	return ret
}

func GetInt64ValFromJsonDoc(s *parser.Selector,
	jsonDoc *simplejson.Json) int64 {
	if s == nil {
		return 0
	}

	branchs := strings.Split(s.JsonKey, ".")
	vals := GetJsonObjectsFromDoc(branchs, jsonDoc)
	if len(vals) == 1 {
		val, err := vals[0].Int64()
		if err != nil {
			dlog.Error("get int64 value from json get error:%s",
				err.Error())
			return 0
		}
		return val
	}
	return 0
}

func GetIntValFromJsonDoc(s *parser.Selector,
	jsonDoc *simplejson.Json) int {
	if s == nil {
		return 0
	}

	branchs := strings.Split(s.JsonKey, ".")
	vals := GetJsonObjectsFromDoc(branchs, jsonDoc)
	if len(vals) == 1 {
		val, err := vals[0].Int()
		if err != nil {
			dlog.Error("get int value from json get error:%s",
				err.Error())
			return 0
		}
		return val
	}
	return 0
}

func GetFloat64ValFromJsonDoc(s *parser.Selector,
	jsonDoc *simplejson.Json) float64 {
	if jsonDoc == nil {
		return 0.0
	}

	branchs := strings.Split(s.JsonKey, ".")
	vals := GetJsonObjectsFromDoc(branchs, jsonDoc)
	if len(vals) == 1 {
		val, err := vals[0].Float64()
		if err == nil {
			return val
		}
	}
	return 0.0
}

func GetStringValFromJsonDoc(s *parser.Selector,
	jsonDoc *simplejson.Json) string {
	ret := GetRawStringValFromJsonDoc(s, jsonDoc)
	if len(ret) != 0 && len(s.Extractor) != 0 {
		_, ret = common.StringExtract(ret, s.Extractor)
	}
	return ret

}

func GetRawStringValFromJsonDoc(s *parser.Selector,
	jsonDoc *simplejson.Json) string {
	if s == nil {
		return ""
	}
	branchs := strings.Split(s.JsonKey, ".")
	vals := GetJsonObjectsFromDoc(branchs, jsonDoc)
	if len(vals) == 1 {
		val, err := vals[0].String()
		if err == nil {
			return val
		}
		intVal, err := vals[0].Int()
		if err == nil {
			return fmt.Sprintf("%d", intVal)
		}
		floatVal, err := vals[0].Float64()
		if err == nil {
			return fmt.Sprintf("%f", floatVal)
		}
	}
	return ""
}

func GetJsonObjectsFromDoc(branch []string,
	jsonDoc *simplejson.Json) []*simplejson.Json {
	if jsonDoc == nil {
		return nil
	}
	if len(branch) == 0 {
		return []*simplejson.Json{jsonDoc}
	}

	b0 := branch[0]
	k, n := keyIndex(b0)
	value := jsonDoc.Get(k)
	if k == kJsonFakeRoot {
		value = jsonDoc
	} else if k == kMapToArray {
		array := MapToArray(jsonDoc)
		ret := []*simplejson.Json{}
		for _, v := range array {
			rs := GetJsonObjectsFromDoc(branch[1:len(branch)], v)
			for _, r := range rs {
				ret = append(ret, r)
			}
		}
		return ret
	}

	if len(n) == 0 {
		return GetJsonObjectsFromDoc(branch[1:len(branch)], value)
	}

	if n == "*" {
		values, err := value.Array()
		if err != nil {
			return nil
		}
		ret := []*simplejson.Json{}
		for i := 0; i < len(values); i++ {
			rs := GetJsonObjectsFromDoc(branch[1:len(branch)],
				value.GetIndex(i))
			for _, r := range rs {
				ret = append(ret, r)
			}
		}
		return ret
	}

	ni, err := strconv.Atoi(n)
	if err != nil {
		return nil
	}
	return GetJsonObjectsFromDoc(branch[1:len(branch)],
		value.GetIndex(ni))
}

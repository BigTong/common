package json

import (
	"strings"
	"testing"

	"github.com/BigTong/common/file"
	"github.com/BigTong/common/parser"
	"github.com/bmizerany/assert"
)

func TestJosnParse(t *testing.T) {
	testData := file.ReadFileToString("json_string")
	jsonParser := NewJsonParser(strings.NewReader(testData))
	s := &parser.Selector{JsonKey: "pagecount"}
	assert.Equal(t, jsonParser.GetInt(s), int(2))
	assert.Equal(t, jsonParser.GetInt64(s), int64(2))

	s.JsonKey = "test"
	assert.Equal(t, jsonParser.GetFloat64(s), float64(0.34))

	s.JsonKey = "articlelist[*]"
	result := jsonParser.GetJsonObjects(strings.Split(s.JsonKey, "."))
	assert.Equal(t, len(result), int(20))

	s.JsonKey = "title"
	assert.Equal(t, GetStringValFromJsonDoc(s, result[0]), "售3.79-4.79万元 比亚迪新款F0上市")
}

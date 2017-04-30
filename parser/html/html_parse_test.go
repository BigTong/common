package html

import (
	"strings"
	"testing"

	"github.com/BigTong/common/file"
	"github.com/BigTong/common/parser"
	"github.com/bmizerany/assert"
)

func TestHtmlParse(t *testing.T) {
	testData := file.ReadFileToString("html_string")
	htmlParser := NewHtmlParser(strings.NewReader(testData))

	s := &parser.Selector{
		Xpath: "div.sphinxsidebarwrapper div#searchbox h3",
	}
	assert.Equal(t, htmlParser.GetString(s), "Quick search")

	s = &parser.Selector{
		Xpath: "div.int_test",
	}
	assert.Equal(t, htmlParser.GetInt64(s), int64(99))

	s = &parser.Selector{
		Xpath: "div.float_test",
	}
	assert.Equal(t, htmlParser.GetFloat64(s), float64(99.99))

	s = &parser.Selector{
		Xpath: "div#id3.section div.toctree-wrapper.compound ul li.toctree-l1",
	}

	assert.Equal(t, len(htmlParser.GetSelections(s)), int(20))

	s = &parser.Selector{
		Xpath: "div.toctree-wrapper.compound ul li",
		Index: "1",
	}
	t.Log(htmlParser.GetString(s))
}

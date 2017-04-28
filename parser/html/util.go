package html_parser

import (
	"strconv"
	"strings"

	"common"
	"common/dlog"
	"common/parser"
	"github.com/PuerkitoBio/goquery"
)

func StringToInt64(src string) int64 {
	val, err := strconv.ParseInt(src, 10, 64)
	if err != nil {
		dlog.Error("parse int64 get error:%s, src:%s", err.Error(), src)
		return 0
	}
	return val
}

func StringToInt(src string) int {
	return int(StringToInt64(src))
}

func StringToFloat64(src string) float64 {
	val, err := strconv.ParseFloat(src, 64)
	if err != nil {
		dlog.Error("parse float64 get error:%s, src:%s", err.Error(), src)
		return 0.0
	}
	return val
}

func ParseInt(doc *goquery.Selection, selector *parser.Selector) int {
	src := ParseString(doc, selector)
	return StringToInt(src)
}

func ParseFloat64(doc *goquery.Selection, selector *parser.Selector) float64 {
	src := ParseString(doc, selector)
	return StringToFloat64(src)
}

func ParseInt64(doc *goquery.Selection, selector *parser.Selector) int64 {
	src := ParseString(doc, selector)
	return StringToInt64(src)
}

func ParseString(doc *goquery.Selection, selector *parser.Selector) string {
	find := false
	ret := ""
	if selector == nil || doc == nil {
		return ret
	}

	if len(selector.Xpath) != 0 {
		doc = doc.Find(selector.Xpath)
	}

	if doc == nil || doc.Size() == 0 {
		return ret
	}

	if len(selector.Index) != 0 {
		idx, _ := strconv.Atoi(selector.Index)
		doc.Each(func(i int, s *goquery.Selection) {
			if i == idx {
				find = true
				doc = s
				return
			}
		})
		if !find {
			return ""
		}
	}

	if len(selector.Contains) != 0 {
		doc.Each(func(i int, s *goquery.Selection) {
			matched := true
			for _, c := range selector.Contains {
				if !strings.Contains(s.Text(), c) {
					matched = false
				}
			}
			if matched {
				doc = s
			}
		})
	}

	if len(selector.Attr) > 0 {
		ret, _ = doc.Attr(selector.Attr)
	} else {
		ret = doc.Text()
	}

	if len(selector.Extractor) != 0 {
		_, ret = common.StringExtract(ret, selector.Extractor)
	}

	return common.CleanString(ret)
}

func GetOneSelection(doc *goquery.Selection,
	selector *parser.Selector) *goquery.Selection {
	if selector == nil {
		return nil
	}
	if len(selector.Xpath) == 0 {
		return nil
	}
	return doc.Find(selector.Xpath)
}

func GetSelectionFromDoc(doc *goquery.Selection,
	selector *parser.Selector) []*goquery.Selection {
	if selector == nil {
		return []*goquery.Selection{doc}
	}

	if len(selector.Xpath) == 0 {
		return []*goquery.Selection{doc}
	}

	ret := []*goquery.Selection{}
	doc.Find(selector.Xpath).Each(func(i int, s *goquery.Selection) {
		matched := true
		for _, c := range selector.Contains {
			if !strings.Contains(s.Text(), c) {
				matched = false
				break
			}
		}
		if matched {
			ret = append(ret, s)
		}
	})
	return ret
}

package html_parser

import (
	"io"

	"common/dlog"
	"common/parser"
	"github.com/PuerkitoBio/goquery"
)

type HtmlParser struct {
	doc *goquery.Selection
}

func NewHtmlParser(r io.Reader) *HtmlParser {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		dlog.Error("new doc get err:%s", err.Error())
		return nil
	}

	return &HtmlParser{
		doc: doc.Find("html"),
	}
}

func (self *HtmlParser) GetData() *goquery.Selection {
	return self.doc
}

func (self *HtmlParser) GetString(selector *parser.Selector) string {
	return ParseString(self.doc, selector)
}

func (self *HtmlParser) GetInt64(selector *parser.Selector) int64 {
	val := ParseString(self.doc, selector)
	return StringToInt64(val)
}

func (self *HtmlParser) GetInt(selector *parser.Selector) int {
	return int(self.GetInt64(selector))
}

func (self *HtmlParser) GetFloat64(selector *parser.Selector) float64 {
	val := ParseString(self.doc, selector)
	return StringToFloat64(val)
}

func (self *HtmlParser) GetOneSelection(
	selector *parser.Selector) *goquery.Selection {
	return GetOneSelection(self.doc, selector)
}

func (self *HtmlParser) GetSelections(
	selector *parser.Selector) []*goquery.Selection {
	return GetSelectionFromDoc(self.doc, selector)
}

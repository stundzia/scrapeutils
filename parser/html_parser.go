package parser

import (
	"github.com/PuerkitoBio/goquery"
	"io"
)

type HtmlParser struct {
	url string
	baseUrl string
	content string
	Doc *goquery.Document
}

func NewHtmlParser(url string, body io.Reader) (*HtmlParser, error) {
	var err error
	parser := &HtmlParser{
		url:     url,
	}
	parser.Doc, err = goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}
	return parser, nil
}

package parser

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/url"
)

type HtmlParser struct {
	Url string
	BaseUrl string
	content string
	Doc *goquery.Document
}

func NewHtmlParser(urlStr string, body io.Reader) (*HtmlParser, error) {
	var err error
	urlObj, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	parser := &HtmlParser{
		Url:     urlStr,
		BaseUrl: urlObj.Host,
	}
	parser.Doc, err = goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}
	return parser, nil
}

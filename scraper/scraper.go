package scraper

import (
	"f.oxy.works/paulius.stundzia/scrapeutils/parser"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"go.uber.org/zap"

	"f.oxy.works/paulius.stundzia/scrapeutils/proxy"
)

type Scraper struct {
	ProxyPool	*proxy.Pool
	logger *zap.Logger
}

func NewScraper() *Scraper {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("could not initiate zap logger: %s", err)
	}
	return &Scraper{
		ProxyPool:  proxy.NewProxyPool(logger),
		logger:     logger,
	}
}

func (scrap *Scraper) FetchContentBody(url string) (body io.ReadCloser, status int, err error) {
	c := &http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       20 * time.Second,
	}
	p := scrap.ProxyPool.GetRandomProxy()
	if p != nil {
		c.Transport = p.GetTransport()
	}
	resp, err := c.Get(url)
	if err != nil {
		scrap.logger.Error("error during GET", zap.String("error", err.Error()), zap.String("target url", url))
		return nil, 0, err
	}
	return resp.Body, resp.StatusCode, err
}

func (scrap *Scraper) FetchContentBytes(url string) ([]byte, int, error) {
	body, status, err := scrap.FetchContentBody(url)
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		scrap.logger.Error("error reading response body from GET", zap.String("error", err.Error()), zap.String("target url", url))
	}
	defer body.Close()
	return bodyBytes, status, nil
}


func (scrap *Scraper) FetchAndReturnParser(url string) (htmlParser *parser.HtmlParser, err error) {
	body, status, err := scrap.FetchContentBody(url)
	if err != nil || status != 200 {
		scrap.logger.Error("unable to get valid response for parsing", zap.String("target url", url), zap.Int("status", status))
	}
	htmlParser, err = parser.NewHtmlParser(url, body)
	if err != nil {
		scrap.logger.Error("unable to init parser", zap.String("target url", url), zap.String("error", err.Error()))
		return nil, err
	}
	return htmlParser, nil
}
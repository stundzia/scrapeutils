package scraper

import (
	"f.oxy.works/paulius.stundzia/scrapeutils/proxy"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Scraper struct {
	proxyPool	*proxy.Pool
	logger *zap.Logger
}

func NewScraper() *Scraper {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("could not initiate zap logger: %s", err)
	}
	return &Scraper{
		proxyPool:  proxy.NewProxyPool(logger),
		logger:     logger,
	}
}

func (scrap *Scraper) FetchContent(url string) ([]byte, error) {
	p := scrap.proxyPool.GetRandomProxy()
	if p != nil {

	}
	c := &http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       20 * time.Second,
	}
	resp, err := c.Get(url)
	if err != nil {
		scrap.logger.Error("error during GET", zap.String("error", err.Error()), zap.String("target url", url))
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		scrap.logger.Error("error reading response body from GET", zap.String("error", err.Error()), zap.String("target url", url))
	}
	return body, nil
}

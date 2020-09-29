package proxy

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"go.uber.org/zap"

	"f.oxy.works/paulius.stundzia/scrapeutils/db"
)

type Pool struct {
	proxies map[*proxy]int
	proxyDB *db.ProxyDB
	logger *zap.Logger
}

type proxy struct {
	urlStr string
	url *url.URL
}


func NewProxyPool(logger *zap.Logger) *Pool {
	return &Pool{
		proxies: make(map[*proxy]int),
		proxyDB: nil,
		logger:  logger,
	}
}

func (pp *Pool) InitMysqlProxyDB(user string, password string, name string, table string) {
	pp.proxyDB = db.NewProxyDB(user, password, name, table)
}

func (pp *Pool) GetRandomProxy() *proxy {
	for p := range pp.proxies {
		pp.proxies[p]++
		return p
	}
	return nil
}

func (pp *Pool) AddProxyViaUrlString(url string) {
	pp.proxies[GetProxyFromUrlString(url)] = 0
}

func (p *proxy) GetTransport() *http.Transport {
	transport := &http.Transport{
		Proxy: http.ProxyURL(p.url),
	}
	return transport
}

func GetProxyFromUrlString(s string) (p *proxy) {
	p = &proxy{
		urlStr:   s,
		url:      nil,
	}
	proxyURL, err := url.Parse(s)
	if err != nil {
		log.Fatal(err)
		return p
	}
	p.url = proxyURL
	return p
}

func (pp *Pool) AddProxyFromDBProxy(proxyRow *db.ProxyRow) {
	var proxyUrl string
	if len(proxyRow.Username) > 0 && len(proxyRow.Password) > 0 {
		proxyUrl = fmt.Sprintf("http://%s:%s@%s:%d", proxyRow.Username, proxyRow.Password, proxyRow.Host, proxyRow.Port)
	} else {
		proxyUrl = fmt.Sprintf("http://%s:%s", proxyRow.Host, proxyRow.Port)
	}
	pp.AddProxyViaUrlString(proxyUrl)
}

func (pp *Pool) PopulatePoolFromDB(basesource string, proxyType string, limit int) {
	pp.logger.Info("fetching proxies from DB", zap.Int("limit", limit), zap.String("basesource", basesource), zap.String("proxy type", proxyType))
	proxieRows, err := pp.proxyDB.GetProxies(basesource, proxyType, limit)
	if err != nil {
		pp.logger.Error("proxy get from db error", zap.String("error", err.Error()))
	}
	for _, r := range proxieRows {
		pp.AddProxyFromDBProxy(r)
	}
}
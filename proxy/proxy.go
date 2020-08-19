package proxy

import (
	"log"
	"net/http"
	"net/url"
)

type Pool struct {
	proxies map[*Proxy]int
}

func (pp *Pool) getRandomProxy() *Proxy {
	for p := range pp.proxies {
		pp.proxies[p]++
		return p
	}
	return nil
}

func (pp *Pool) AddProxyViaUrlString(s string) {
	pp.proxies[GetProxyFromUrlString(s)] = 0
}

type Proxy struct {
	urlStr string
	url *url.URL
}

func (p *Proxy) GetTransport() *http.Transport {
	transport := &http.Transport{
		Proxy: http.ProxyURL(p.url),
	}
	return transport
}

func GetProxyFromUrlString(s string) (proxy *Proxy) {
	proxy = &Proxy{
		urlStr:   s,
		url:      nil,
	}
	proxyURL, err := url.Parse(s)
	if err != nil {
		log.Fatal(err)
		return proxy
	}
	proxy.url = proxyURL
	return proxy
}

package config

import (
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/proxy"
)

var torProxyURL = "socks5://127.0.0.1:9150"

func SetScrappingProxy() colly.ProxyFunc {
	rotateProxies, err := proxy.RoundRobinProxySwitcher(torProxyURL)
	if err != nil {
		return nil
	}
	return rotateProxies
}

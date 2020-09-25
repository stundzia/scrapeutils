package scraper

import (
	"testing"
)

func TestScraper_FetchContentNoProxy(t *testing.T) {
	scrap := NewScraper()
	res, err := scrap.FetchContent("https://ip.oxylabs.io/")
	if err != nil {
		t.Errorf("error during fetch: %s", err.Error())
	}
	dotCount := 0
	for _, c := range res {
		if c == 46 {
			dotCount++
		}
	}
	if dotCount != 3 {
		t.Errorf("unexpected content:\n %s", string(res))
	}
}

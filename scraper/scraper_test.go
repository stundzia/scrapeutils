package scraper

import (
	"testing"
)

func TestScraper_FetchContentNoProxy(t *testing.T) {
	scrap := NewScraper()
	res, status, err := scrap.FetchContent("https://ip.oxylabs.io/")
	if err != nil {
		t.Errorf("error during fetch: %s", err.Error())
	}
	if status != 200 {
		t.Errorf("unexpected status code: %d", status)
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

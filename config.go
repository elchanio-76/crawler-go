package main

import (
	"net/url"
	"sync"
)

const MAX_CONCURRENCY = 5
const MAX_PAGES = 20
const REPORT_FILE = "report.json"

type config struct {
	pages              map[string]PageData
	baseURL            *url.URL
	concurrencyControl chan struct{}
	mux                *sync.Mutex
	wg                 *sync.WaitGroup
	maxPages			int
}

// AddPageIfNew atomically checks the dedup set and maxPages limit, then inserts.
// Returns false (and does not insert) if the page was already visited or the limit is reached.
func (c *config) AddPageIfNew(urlStr string, data PageData) bool {
	c.mux.Lock()
	defer c.mux.Unlock()
	if len(c.pages) >= c.maxPages {
		return false
	}
	if _, exists := c.pages[urlStr]; exists {
		return false
	}
	c.pages[urlStr] = data
	return true
}

func (c *config) IsPageVisited(urlStr string) bool {
	c.mux.Lock()
	defer c.mux.Unlock()
	_, ok := c.pages[urlStr]
	return ok
}

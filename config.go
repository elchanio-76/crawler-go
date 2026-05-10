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

func (c *config) AddPage(urlStr string, data PageData) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.pages[urlStr] = data
}

func (c *config) IsPageVisited(urlStr string) bool {
	c.mux.Lock()
	defer c.mux.Unlock()
	_, ok := c.pages[urlStr]
	return ok
}

package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"strconv"
)

func getHTML(rawURL string) (string, error) {
	// Returns the HTML of the given URL as a string.
	// If there is an error fetching the HTML, return an empty string and the error.

	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Accept", "text/html")
	req.Header.Add("User-Agent", "BootCrawlerGo/1.0")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode > 399 {
		return "", fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		return "", fmt.Errorf("content type error: %s", contentType)
	}
	html, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(html), nil

}

func (cfg *config) crawlPage(rawCurrentURL string) {
	// Crawls the given URL and updates the pages map with the number of times each URL is found.
	// If there is an error fetching the HTML, return without updating the pages map.
	// If the URL is not in the same domain as the base URL, return without updating the pages map.
	// If the URL has already been crawled, return without fetching the HTML, just update the pages map count

	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("error parsing current URL: %v", err)
		return
	}
	if cfg.baseURL.Hostname() != currentURL.Hostname() {
		return
	}
	normalizedURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("error normalizing URL: %v", err)
		return
	}
	// Optimistic early check to avoid an unnecessary HTTP request for an already-visited page.
	// The definitive dedup+maxPages enforcement happens atomically inside AddPageIfNew.
	if cfg.IsPageVisited(normalizedURL) {
		return
	}

	html, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("error fetching HTML: %v", err)
		return
	}
	fmt.Printf("fetched %s, length: %d\n", rawCurrentURL, len(html))

	pageData, err := extractPageData(html, rawCurrentURL)
	if err != nil {
		fmt.Printf("error extracting page data from %s: %v\n", rawCurrentURL, err)
		return
	}
	if !cfg.AddPageIfNew(normalizedURL, pageData) {
		return
	}
	urls, err := getURLsFromHTML(html, cfg.baseURL)
	fmt.Printf(" --> found %d URLs in %s\n", len(urls), rawCurrentURL)
	if err != nil {
		fmt.Printf("error getting URLs from HTML: %v", err)
		return
	}
	for _, url := range urls {
		// spawn go routine to crawl in parallel
		cfg.wg.Add(1)
		go func() {
			cfg.concurrencyControl <- struct{}{}
			defer cfg.wg.Done()
			defer func() { <-cfg.concurrencyControl }()
			cfg.crawlPage(url)
		}()
	}
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("no website provided")
		os.Exit(1)
	}
	if len(args) > 3 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}
	baseURL, err := url.Parse(args[0])
	if err != nil {
		fmt.Printf("error parsing base URL: %v", err)
		os.Exit(1)
	}
	maxConcurrency := MAX_CONCURRENCY
	maxPages := MAX_PAGES
	if len(args) >= 2 {
		maxConcurrency, err = strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("error parsing max pages: %v", err)
			os.Exit(1)
		}
	}
	if len(args) == 3 {
		maxPages, err = strconv.Atoi(args[2])
		if err != nil {
			fmt.Printf("error parsing max pages: %v", err)
			os.Exit(1)
		}
	}
	fmt.Printf("max pages: %d\n", maxPages)
	fmt.Printf("max concurrency: %d\n", maxConcurrency)
	fmt.Printf("starting crawl of: %s\n", baseURL)
	pages := make(map[string]PageData)

	// create a new config struct
	cfg := &config{
		pages:              pages,
		baseURL:            baseURL,
		concurrencyControl: make(chan struct{}, maxConcurrency),
		mux:                &sync.Mutex{},
		wg:                 &sync.WaitGroup{},
		maxPages:			maxPages,
	}

	cfg.crawlPage(baseURL.String())
	cfg.wg.Wait()

	if err := writeJSONReport(cfg.pages, REPORT_FILE); err != nil {
		fmt.Printf("error writing report: %v\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}

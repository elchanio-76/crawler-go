package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
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

func crawlPage(rawBAseURL, rawCurrentURL string, pages map[string]int) {
	// Crawls the given URL and updates the pages map with the number of times each URL is found.
	// If there is an error fetching the HTML, return without updating the pages map.
	// If the URL is not in the same domain as the base URL, return without updating the pages map.
	// If the URL has already been crawled, return without fetching the HTML, just update the pages map count
	baseURL, err := url.Parse(rawBAseURL)
	if err != nil {
		fmt.Printf("error parsing base URL: %v", err)
		return
	}
	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("error parsing current URL: %v", err)
		return
	}
	if baseURL.Hostname() != currentURL.Hostname() {
		return
	}
	normalizedURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("error normalizing URL: %v", err)
		return
	}
	if pages[normalizedURL] > 0 {
		pages[normalizedURL]++
		return
	}
	fmt.Printf("fetching %s\n", rawCurrentURL)
	html, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("error fetching HTML: %v", err)
		return
	}

	// Print a few characters of the HTML to confirm it was fetched.
	chars := 100
	if len(html)<200 {
		chars = len(html)
	}
	fmt.Printf("fetched %s\n", html[:chars])

	pages[normalizedURL] = 1
	urls, err := getURLsFromHTML(html, baseURL)
	if err != nil {
		fmt.Printf("error getting URLs from HTML: %v", err)
		return
	}
	for _, url := range urls {
		crawlPage(rawBAseURL, url, pages)
	}
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("no website provided")
		os.Exit(1)
	}
	if len(args) > 1 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}
	baseURL := args[0]
	fmt.Printf("starting crawl of: %s\n",baseURL)
	pages := make(map[string]int)
	crawlPage(baseURL, baseURL, pages)
	
	for url, count := range pages {
		fmt.Printf("Found %d internal links to %s\n", count, url)
	}
	
	os.Exit(0)
}

package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getHeadingFromHTML(html string) string {
	// Returns the first heading (H1) of the html,
	// or the first heading (H2) if there is no H1.
	// if no heading is found, return an empty string.

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}

	// Find the first H1 heading
	h1 := doc.Find("h1").First()
	if h1.Length() > 0 {
		return h1.Text()
	}

	// Find the first H2 heading
	h2 := doc.Find("h2").First()
	if h2.Length() > 0 {
		return h2.Text()
	}

	return ""
}

func getFirstParagraphFromHTML(html string) string {
	// Returns the first paragraph from the HTML after the <main> tag,
	// or the first paragraph in the HTML if there is no <main> tag.
	// If no paragraph is found, return an empty string.

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}

	// Find the first paragraph in the main tag
	main := doc.Find("main").First()
	if main.Length() > 0 {
		p := main.Find("p").First()
		if p.Length() > 0 {
			return p.Text()
		}
	}

	// Find the first paragraph in the HTML
	p := doc.Find("p").First()
	if p.Length() > 0 {
		return p.Text()
	}

	return ""
}

func getURLsFromHTML(html string, baseURL *url.URL) ([]string, error) {
	// Returns all the URLs from the HTML.
	// If no URLs are found, return an empty slice.
	// If there is an error parsing the HTML, return an empty slice
	// URLs are absolute URLs, not relative paths.

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return []string{}, fmt.Errorf("Error parsing html: %s", err)
	}

	urls := []string{}
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			// Resolve relative URLs to absolute URLs
			absoluteURL, err := baseURL.Parse(href)
			if err == nil {
				// Trim section markers
				absoluteURL.Fragment = ""
				urls = append(urls, absoluteURL.String())
			}
		}
	})

	return urls, nil
}

func getImagesFromHTML(htm string, baseURL *url.URL) ([]string, error) {
	// Returns all the image URLs from the HTML.
	// If no image URLs are found, return an empty slice.
	// If there is an error parsing the HTML, return an empty slice
	// Image URLs are absolute URLs, not relative paths.

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htm))
	if err != nil {
		return []string{}, fmt.Errorf("Error parsing html: %s", err)
	}

	urls := []string{}
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if exists {
			// Resolve relative URLs to absolute URLs
			if src != "" {
				absoluteURL, err := baseURL.Parse(src)
				if err == nil {
					urls = append(urls, absoluteURL.String())
				}
			}
		}
	})

	return urls, nil
}

type PageData struct {
	URL            string   `json:"url"`
	Heading        string   `json:"heading"`
	FirstParagraph string   `json:"first_paragraph"`
	OutgoingLinks  []string `json:"outgoing_links"`
	ImageURLs      []string `json:"image_urls"`
}

func extractPageData(html, pageURL string) (PageData, error) {
	baseURL, err := url.Parse(pageURL)
	if err != nil {
		return PageData{}, fmt.Errorf("error parsing page URL %s: %w", pageURL, err)
	}
	urls, err := getURLsFromHTML(html, baseURL)
	if err != nil {
		return PageData{}, fmt.Errorf("error getting URLs from %s: %w", pageURL, err)
	}
	imageURLs, err := getImagesFromHTML(html, baseURL)
	if err != nil {
		return PageData{}, fmt.Errorf("error getting images from %s: %w", pageURL, err)
	}

	return PageData{
		URL:            pageURL,
		Heading:        getHeadingFromHTML(html),
		FirstParagraph: getFirstParagraphFromHTML(html),
		OutgoingLinks:  urls,
		ImageURLs:      imageURLs,
	}, nil
}

func writeJSONReport(pages map[string]PageData, filename string) error {

	// Create a slice to hold the sorted keys
	keys := make([]string, 0, len(pages))
	for key := range pages {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Convert the map to a slice of PageData
	var pageDataSlice []PageData
	for _, k := range keys {
		pageDataSlice = append(pageDataSlice, pages[k])
	}

	// Marshal the slice to JSON
	jsonData, err := json.MarshalIndent(pageDataSlice, "", "  ")
	if err != nil {
		return fmt.Errorf("Error marshaling JSON: %s", err)
	}

	// Write the JSON data to the file
	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("Error writing JSON file: %s", err)
	}

	return nil
}

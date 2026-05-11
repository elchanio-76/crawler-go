package main

import (
	"net/url"
	"strings"
)

func normalizeURL(rawURL string) (string, error) {
	// TODO: implement this function
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	parsedURL.Scheme = ""
	parsedURL.Fragment = ""
	parsedURL.User = nil
	parsedURL.Host = strings.TrimSuffix(parsedURL.Host, "/")
	parsedURL.Host = strings.ToLower(parsedURL.Host)
	// Query params are intentionally stripped: URLs that differ only by query string
	// (e.g. /page?id=1 vs /page?id=2) are treated as the same page for deduplication.
	// This means param-variant URLs collected by getURLsFromHTML are silently de-duped here.
	parsedURL.RawQuery = ""
	parsedURL.Path = strings.TrimSuffix(parsedURL.Path, "/")

	normalizedURL := strings.TrimPrefix(parsedURL.String(), "//")

	return normalizedURL, nil
}

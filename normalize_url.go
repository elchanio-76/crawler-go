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
	parsedURL.RawQuery = ""
	parsedURL.Path = strings.TrimSuffix(parsedURL.Path, "/")

	normalizedURL := strings.TrimPrefix(parsedURL.String(), "//")

	return normalizedURL, nil
}

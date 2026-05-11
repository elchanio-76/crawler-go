## Crawler in Go

This project was built as an exercise in Go. Follows the [boot.dev practical course](https://www.boot.dev/courses/build-web-scraper-golang)

### Description

The project is a Go application. Needs Go v 1.22.1 or later. 

### Usage
First build the application:
`go build -o crawler`

Run with syntax `crawler <URL> <MAX_CONCURRENT_THREADS> <MAX_PAGES>` where:
- `<URL>` is the URL to scrape
- `<MAX_CONCURRENT_THREADS>` are the maximum concurrent threads to use while scraping
- `<MAX_PAGES>` are the maximum number of pages to crawl.

The program produces a JSON formatted report with the basic stats for each page crawled:
```
{
    "url": string
    "heading": string (first heading)
    "first_paragraph": string
    "outgoing links": [ string ] list of all links in the page
    "image_urls": [ string ] list of all image URLs in the page
}
```

The report is saved as `./report.json`

## License

MIT

## Motivation

The project was intended as a Go learning exercise in /net/http and go routines. 


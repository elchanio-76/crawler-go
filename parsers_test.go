package main

import (
	"net/url"
	"reflect"
	"testing"
)

func TestGetHeadingFromHTML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "no headings",
			input: `
			<html>
				<body>
					<a href="https://blog.boot.dev"><span>Boot.dev></span></a>
				</body>
			</html>
			`,
			expected: "",
		},
		{
			name: "heading in body",
			input: `
			<html>
				<body>
					<h1>Hello, world!</h1>
				</body>
			</html>
			`,
			expected: "Hello, world!",
		},
		{
			name: "heading in main",
			input: `
			<html>
				<body>
					<main>
						<div>
							<h1>Hello, world!</h1>
						</div>
					</main>
				</body>
			</html>
			`,
			expected: "Hello, world!",
		},
		{
			name: "heading in div",
			input: `
			<html>
				<body>
					<div>
						<h1>Hello, world!</h1>
					</div>
				</body>
			</html>
			`,
			expected: "Hello, world!",
		},
		{
			name: "H2 heading in div",
			input: `
			<html>
				<body>
					<div>
						<h2>Hello, world!</h2>
					</div>
				</body>
			</html>
			`,
			expected: "Hello, world!",
		},
		{
			name: "H2 before H1 heading in div",
			input: `
			<html>
				<body>
					<div>
						<h2>Hello, world, h2!</h2>
						<h1>Hello, world!</h1>
					</div>
				</body>
			</html>
			`,
			expected: "Hello, world!",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := getHeadingFromHTML(tc.input)
			if actual != tc.expected {
				t.Errorf("Test %v - %s FAIL: expected %v, got %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}

func TestGetFirstParagraphFromHTML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "no paragraphs",
			input: `
			<html>
				<body>
					<a href="https://blog.boot.dev"><span>Boot.dev></span></a>
				</body>
			</html>
			`,
			expected: "",
		},
		{
			name: "paragraph in body",
			input: `
			<html>
				<body>
					<p>Hello, world!</p>
				</body>
			</html>
			`,
			expected: "Hello, world!",
		},
		{
			name: "paragraph in main",
			input: `
			<html>
				<body>
					<main>
						<p>Hello, world!</p>
					</main>
				</body>
			</html>
			`,
			expected: "Hello, world!",
		},
		{
			name: "paragraph in div",
			input: `
			<html>
				<body>
					<div>
						<p>Hello, world!</p>
					</div>
				</body>
			</html>
			`,
			expected: "Hello, world!",
		},
		{
			name: "paragraph in div after main",
			input: `
			<html>
				<body>
					<main>
						<div>
							<p>Hello, world!</p>
						</div>
					</main>
					<div>
						<p>Hello, world inside!</p>
					</div>
				</body>
			</html>
			`,
			expected: "Hello, world!",
		},
		{
			name: "paragraph in div before main",
			input: `
			<html>
				<body>
					<div>
						<p>Hello, world!</p>
					</div>
					<main>
						<div>
							<p>Hello, world inside!</p>
						</div>
					</main>
				</body>
			</html>
			`,
			expected: "Hello, world inside!",
		},
	}
	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := getFirstParagraphFromHTML(tc.input)
			if actual != tc.expected {
				t.Errorf("Test %v - %s FAIL: expected %v, got %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}

func TestGetURLsFromHTML(t *testing.T) {
	tests := []struct {
		name      string
		inputURL  *url.URL
		inputBody string
		expected  []string
	}{
		{
			name: "absolute and relative URLs",
			inputURL: &url.URL{
				Scheme: "https",
				Host:   "blog.boot.dev",
			},
			inputBody: `
			<html>
				<body>
					<a href="/path/one">
						<span>Boot.dev</span>
					</a>
					<a href="https://other.com/path/one">
						<span>Boot.dev</span>
					</a>
				</body>
			</html>
			`,
			expected: []string{"https://blog.boot.dev/path/one", "https://other.com/path/one"},
		},
		{
			name: "relative URLs",
			inputURL: &url.URL{
				Scheme: "https",
				Host:   "blog.boot.dev",
			},
			inputBody: `
			<html>
				<body>
					<a href="/path/one">
						<span>Boot.dev</span>
					</a>
				</body>
			</html>
			`,
			expected: []string{"https://blog.boot.dev/path/one"},
		},
		{
			name: "absolute URLs",
			inputURL: &url.URL{
				Scheme: "https",
				Host:   "blog.boot.dev",
			},
			inputBody: `
			<html>
				<body>
					<a href="https://other.com/path/one">
						<span>Boot.dev</span>
					</a>
				</body>
			</html>
			`,
			expected: []string{"https://other.com/path/one"},
		},
		{
			name: "no URLs",
			inputURL: &url.URL{
				Scheme: "https",
				Host:   "blog.boot.dev",
			},
			inputBody: `
			<html>
				<body>
					<span>Boot.dev</span>
				</body>
			</html>
			`,
			expected: []string{},
		},
		{
			name: "URLs with fragments",
			inputURL: &url.URL{
				Scheme: "https",
				Host:   "blog.boot.dev",
			},
			inputBody: `
			<html>
				<body>
					<a href="https://other.com/path/one#section">
						<span>Boot.dev</span>
					</a>
				</body>
			</html>
			`,
			expected: []string{"https://other.com/path/one"},
		},
		{
			name: "URLs with query parameters",
			inputURL: &url.URL{
				Scheme: "https",
				Host:   "blog.boot.dev",
			},
			inputBody: `
			<html>
				<body>
					<a href="https://other.com/path/one?query=value">
						<span>Boot.dev</span>
					</a>
				</body>
			</html>
			`,
			expected: []string{"https://other.com/path/one?query=value"},
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := getURLsFromHTML(tc.inputBody, tc.inputURL)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			}
			if len(actual) != len(tc.expected) {
				t.Errorf("Test %v - '%s' FAIL: expected %v, got %v", i, tc.name, tc.expected, actual)
				return
			}
			for i, url := range actual {
				if url != tc.expected[i] {
					t.Errorf("Test %v - '%s' FAIL: expected URL %v, got %v", i, tc.name, tc.expected[i], url)
				}
			}
		})
	}
}

func TestGetImagesFromHTML(t *testing.T) {
	testURL := &url.URL{
		Scheme: "https",
		Host:   "blog.boot.dev",
	}
	tests := []struct {
		name      string
		inputURL  *url.URL
		inputBody string
		expected  []string
	}{
		{
			name:     "absolute and relative URLs",
			inputURL: testURL,
			inputBody: `
			<html>
				<body>
					<img src="https://blog.boot.dev/path/one" />
					<img src="/path/two" />
				</body>
			</html>
			`,
			expected: []string{"https://blog.boot.dev/path/one", "https://blog.boot.dev/path/two"},
		},
		{
			name:     "relative URLs",
			inputURL: testURL,
			inputBody: `
			<html>
				<body>
					<img src="/path/one" />
					<img src="/path/two" />
				</body>
			</html>
			`,
			expected: []string{"https://blog.boot.dev/path/one", "https://blog.boot.dev/path/two"},
		},
		{
			name:     "absolute URLs",
			inputURL: testURL,
			inputBody: `
			<html>
				<body>
					<img src="https://blog.boot.dev/path/one" />
					<img src="https://blog.boot.dev/path/two" />
				</body>
			</html>
			`,
			expected: []string{"https://blog.boot.dev/path/one", "https://blog.boot.dev/path/two"},
		},
		{
			name:     "no image URLs",
			inputURL: testURL,
			inputBody: `
			<html>
				<body>
					<h1>Hello World!</h1>
					<p>This is a paragraph.</p>
				</body>
			</html>
			`,
			expected: []string{},
		},
		{
			name:     "image URLs with query parameters",
			inputURL: testURL,
			inputBody: `
			<html>
				<body>
					<img src="https://blog.boot.dev/path/one?query=1" />
					<img src="/path/two?query=2" />
				</body>
			</html>
			`,
			expected: []string{"https://blog.boot.dev/path/one?query=1", "https://blog.boot.dev/path/two?query=2"},
		},
		{
			name:     "image URLs with missing src",
			inputURL: testURL,
			inputBody: `
			<html>
				<body>
					<img />
					<img src="" />
				</body>
			</html>
			`,
			expected: []string{},
		},
	}
	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := getImagesFromHTML(tc.inputBody, tc.inputURL)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			}
			if len(actual) != len(tc.expected) {
				t.Errorf("Test %v - '%s' FAIL: expected %v, got %v", i, tc.name, tc.expected, actual)
				return
			}
			for i, url := range actual {
				if url != tc.expected[i] {
					t.Errorf("Test %v - '%s' FAIL: expected URL %v, got %v", i, tc.name, tc.expected[i], url)
				}
			}
		})
	}
}


func TestExtractPageData(t *testing.T) {
	testURL := "https://blog.boot.dev"
	tests := []struct {
		name     string
		inputURL string
		inputBody string
		expected PageData
	}{
		{
			name:     "simple page",
			inputURL: testURL,
			inputBody: `
			<html>
				<body>
					<h1>Hello World!</h1>
					<p>This is a paragraph.</p>
					<a href="/path/one">Boot.dev</a>
					<a href="https://other.com/path/one">Something else</a>
					<img src="https://blog.boot.dev/path/one" />
					<img src="/path/two" />
				</body>
			</html>
			`,
			expected: PageData{
				URL:           testURL,
				Heading:       "Hello World!",
				FirstParagraph: "This is a paragraph.",
				OutgoingLinks: []string{"https://blog.boot.dev/path/one", "https://other.com/path/one"},
				ImageURLs:     []string{"https://blog.boot.dev/path/one", "https://blog.boot.dev/path/two"},
			},
		},
		{
			name:     "no images",
			inputURL: testURL,
			inputBody: `
			<html>
				<body>
					<h1>Hello World</h1>
					<p>This is a paragraph.</p>
					<a href="/path/one">Boot.dev</a>
					<a href="https://other.com/path/one">Something else</a>
				</body>
			</html>
			`,
			expected: PageData{
				URL:           testURL,
				Heading:       "Hello World",
				FirstParagraph: "This is a paragraph.",
				OutgoingLinks: []string{"https://blog.boot.dev/path/one", "https://other.com/path/one"},
				ImageURLs:     []string{},
			},
		},
		{
			name:     "no links",
			inputURL: testURL,
			inputBody: `
			<html>
				<body>
					<h1>Hello World</h1>
					<p>This is a paragraph.</p>
					<img src="https://blog.boot.dev/path/one" />
					<img src="/path/two" />
				</body>
			</html>
			`,
			expected: PageData{
				URL:           testURL,
				Heading:       "Hello World",
				FirstParagraph: "This is a paragraph.",
				OutgoingLinks: []string{},
				ImageURLs:     []string{"https://blog.boot.dev/path/one", "https://blog.boot.dev/path/two"},
			},
		},
		{
			name:     "no heading",
			inputURL: testURL,
			inputBody: `
			<html>
				<body>
					<p>This is a paragraph.</p>
					<a href="/path/one">Boot.dev</a>
					<a href="https://other.com/path/one">Something else</a>
					<img src="https://blog.boot.dev/path/one" />
					<img src="/path/two" />
				</body>
			</html>
			`,
			expected: PageData{
				URL:           testURL,
				Heading:       "",
				FirstParagraph: "This is a paragraph.",
				OutgoingLinks: []string{"https://blog.boot.dev/path/one", "https://other.com/path/one"},
				ImageURLs:     []string{"https://blog.boot.dev/path/one", "https://blog.boot.dev/path/two"},
			},
		},
		{
			name:     "no paragraph",
			inputURL: testURL,
			inputBody: `
			<html>
				<body>
					<h1>Hello World</h1>
					<a href="/path/one">Boot.dev</a>
					<a href="https://other.com/path/one">Something else</a>
					<img src="https://blog.boot.dev/path/one" />
					<img src="/path/two" />
				</body>
			</html>
			`,
			expected: PageData{
				URL:           testURL,
				Heading:       "Hello World",
				FirstParagraph: "",
				OutgoingLinks: []string{"https://blog.boot.dev/path/one", "https://other.com/path/one"},
				ImageURLs:     []string{"https://blog.boot.dev/path/one", "https://blog.boot.dev/path/two"},
			},
		},
	}
	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := extractPageData(tc.inputBody, tc.inputURL)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Test %v - PageData mismatch: expected %v, got %v", i, tc.expected, actual)
			}
		})
	}

}
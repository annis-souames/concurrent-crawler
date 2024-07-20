package crawler

import (
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// PageInfo holds information about a web page.
type PageInfo struct {
	Title   string
	UrlList []string
}

// FetchUrl fetches the content of a URL and returns it as a PageInfo struct.
func FetchUrl(url string) (PageInfo, error) {
	resp, err := http.Get(url)
	if err != nil {
		return PageInfo{}, fmt.Errorf("error fetching URL: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return PageInfo{}, fmt.Errorf("error reading response body: %w", err)
	}

	urls, err := ExtractUrls(body)
	if err != nil {
		return PageInfo{}, fmt.Errorf("error extracting URLs: %w", err)
	}
	title, err := ExtractTitle(body)
	if err != nil {
		return PageInfo{}, fmt.Errorf("error extracting title: %w", err)
	}
	return PageInfo{title, urls}, nil
}

// ExtractUrls extracts URLs from a page's HTML content.
func ExtractUrls(page []byte) ([]string, error) {
	doc, err := html.Parse(strings.NewReader(string(page)))
	if err != nil {
		return nil, fmt.Errorf("error parsing HTML: %w", err)
	}

	var urls []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					urls = append(urls, a.Val)
					break
				}
			}
		}
		for c := range n.Attr {
            f(c)
        }
	}
	f(doc)

	return urls, nil
}

// ExtractTitle extracts the title from a page's HTML content.
func ExtractTitle(page []byte) (string, error) {
	doc, err := html.Parse(strings.NewReader(string(page)))
	if err != nil {
		return "", fmt.Errorf("error parsing HTML: %w", err)
	}

	var title string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" && n.FirstChild != nil {
			title = n.FirstChild.Data
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	if title == "" {
		return "", fmt.Errorf("title not found")
	}
	return title, nil
}

// ResolveURL resolves a relative URL against a base URL.
func ResolveURL(inputURL, base string) (string, error) {
	baseURL, err := url.Parse(base)
	if err != nil {
		return "", fmt.Errorf("error parsing base URL: %w", err)
	}
	relURL, err := url.Parse(inputURL)
	if err != nil {
		return "", fmt.Errorf("error parsing input URL: %w", err)
	}
	resolvedURL := baseURL.ResolveReference(relURL)
	return resolvedURL.String(), nil
}

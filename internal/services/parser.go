package services

import (
	"fmt"
	"log"
	"strings"

	"golang.org/x/net/html"
)

// ExtractLinks extracts the valid links from the HTML content using XPath.
func ExtractLinks(content string) ([]string, error) {
	// Extract the links from the content using the XPath
	doc, err := html.Parse(strings.NewReader(content))
	if err != nil {
		log.Printf("failed to parse the HTML content: %v", err)
		return nil, fmt.Errorf("failed to parse the HTML content: %v", err)
	}

	var links []string
	var extract func(*html.Node)
	extract = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					links = append(links, a.Val)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extract(c)
		}
	}

	extract(doc)
	return links, nil
}

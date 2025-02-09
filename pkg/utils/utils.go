package utils

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/temoto/robotstxt"
)

// List of real browser User-Agents to rotate
var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
}

func IsValidDomain(domain string) bool {
	// Validate the domain
	var domainRegex = regexp.MustCompile(`^(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z]{2,}$`)
	return domainRegex.MatchString(domain)
}

func NormalizeURL(url string) string {
	// Normalize the URL
	if !regexp.MustCompile(`^https?://`).MatchString(url) {
		url = "http://" + url
	}
	return url
}

// HTTP Client is a wrapper around the http.Client with a timeout
// It is used to make HTTP requests
type HTTPClient struct {
	Client *http.Client
}

// NewHTTPClient creates a new HTTPClient with custom timeout
func NewHTTPClient(timeout int) *HTTPClient {
	return &HTTPClient{
		Client: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}
}

// FetchPage fetches the page from the given URL and returns the HTML Body
func (c *HTTPClient) FetchPage(url string) (string, error) {
	// Check if crawling is allowed
	if !c.crawlingAllowed(url) {
		return "", fmt.Errorf("crawling is not allowed for the URL: %s", url)
	}

	// Fetch the page from the given URL
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create a new request: %v", err)
	}

	// Rotate User-Agent
	req.Header.Set("User-Agent", userAgents[rand.Intn(len(userAgents))])
	req.Header.Set("Referer", "https://www.google.com/") // Pretend to be a Google referral
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := c.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch the URL %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("non-OK HTTP Status: %d, %v", resp.StatusCode, err)
	}

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read the response body: %v", err)
	}

	return string(respBody), nil
}

// crawlingAllowed checks if the URL is allowed for crawling
func (c *HTTPClient) crawlingAllowed(pageURL string) bool {
	// Check if the URL is allowed for crawling
	parsedURL, err := url.Parse(pageURL)
	if err != nil {
		return false
	}

	// Fetch robots.txt file
	robotsURL := fmt.Sprintf("%s://%s/robots.txt", parsedURL.Scheme, parsedURL.Host)
	resp, err := c.Client.Get(robotsURL)
	if err != nil {
		return true
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return true
	}
	robotsData, err := robotstxt.FromResponse(resp)
	if err != nil {
		return false
	}

	// Check if the User-Agent is allowed to crawl the URL
	return robotsData.TestAgent(parsedURL.Path, "User-Agent")
}

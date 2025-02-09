package crawler

import (
	"context"
	"log"
	"sync"

	"github.com/Vineet0197/shoppin-webcrawler/internal/services"
	"github.com/Vineet0197/shoppin-webcrawler/internal/storage"
	"github.com/Vineet0197/shoppin-webcrawler/pkg/utils"
)

var (
	maxWorkers = 5
)

// Crawler represents a web crawler that crawls the web pages.
type Crawler struct {
	kafka       *services.KafkaStorage
	jsonStorage *storage.JSONStore
	httpClient  *utils.HTTPClient
	wg          sync.WaitGroup
}

// NewCrawler creates a new Crawler with the given Kafka storage, JSON store, and HTTP client.
func NewCrawler(kafka *services.KafkaStorage, httpClient *utils.HTTPClient) *Crawler {
	// Create a new JSON store
	jsonStorage := storage.NewJSONStore("product_data.json")

	return &Crawler{
		kafka:       kafka,
		jsonStorage: jsonStorage,
		httpClient:  httpClient,
	}
}

// NormalizeURL normalizes the given URL and ensures the URL has a proper scheme.
func (c *Crawler) NormalizeURL(url string) string {
	return utils.NormalizeURL(url)
}

// StartCrawling starts crawling the web pages for the given domain.
func (c *Crawler) StartCrawling() {
	// Start crawling the domain
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := 0; i < maxWorkers; i++ {
		c.wg.Add(1)
		go c.crawl(ctx)
	}

	c.wg.Wait()
	c.kafka.Close()
	log.Println("Crawling completed")
}

// crawl crawls the web pages for the given domain.
func (c *Crawler) crawl(ctx context.Context) {
	defer c.wg.Done()

	// Use a map to keep track of processed domains to avoid duplicate processing
	processedDomains := make(map[string]bool)

	c.kafka.Consume(ctx, c.kafka.Topic, func(domain string) {
		// Check if the domain is already processed
		if processedDomains[domain] {
			log.Printf("Domain %s is already processed", domain)
			return
		}

		// Mark the domain as processed
		processedDomains[domain] = true

		// Normalize URL
		normalizedURL := c.NormalizeURL(domain)

		// Fetch the page content
		pageContent, err := c.httpClient.FetchPage(normalizedURL)
		if err != nil {
			log.Printf("failed to fetch the page %s: %v", normalizedURL, err)
			return
		}

		// Extract the product URL
		productURLs, _ := services.ExtractLinks(pageContent)

		// Save to JSON storage
		for _, productURL := range productURLs {
			c.jsonStorage.SaveToDomain(domain, productURL)
		}

		log.Printf("Fetched %d product URLs for domain %s", len(productURLs), domain)
	})
}

package storage

// ProductURL represents the product URL structure
type ProductURL struct {
	Domain string   `json:"domain"`
	URLs   []string `json:"urls"`
}

// ProductData holds extracted product URLs and their domain
type ProductData map[string][]string

// AddURL adds a new URL to the product data corresponding to the given domain
func (pd ProductData) AddURL(domain, url string) {
	if pd[domain] == nil {
		pd[domain] = []string{}
	}
	pd[domain] = append(pd[domain], url)
}

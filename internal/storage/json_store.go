package storage

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

// JSONStore represents a JSON file store for the product data.
type JSONStore struct {
	fileName string
	mu       sync.RWMutex
	data     map[string][]string
}

// NewJSONStore creates a new JSONStore with the given file name.
func NewJSONStore(fileName string) *JSONStore {
	return &JSONStore{
		fileName: fileName,
		data:     make(map[string][]string),
	}
}

// SaveToDomain stores URLs for the given domain.
func (js *JSONStore) SaveToDomain(domain, url string) {
	js.mu.Lock()
	defer js.mu.Unlock()

	js.data[domain] = append(js.data[domain], url)
	js.saveToFile()
}

// saveToFile writes the product data to the JSON file.
func (js *JSONStore) saveToFile() {
	file, err := os.Create(js.fileName)
	if err != nil {
		log.Printf("failed to create the JSON file: %v", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(js.data); err != nil {
		log.Printf("failed to encode the product data: %v", err)
	}
}

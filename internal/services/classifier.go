package services

import (
	"encoding/json"
	"fmt"
	"os"
)

// MLModel represents a simple model that stores product page patterns.
type MLModel struct {
	Patterns []string `json:"patterns"`
}

// LoadModel loads the ML model from a JSON file.
func LoadModel(filePath string) (*MLModel, error) {
	// Load the model from the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open the model file: %v", err)
	}
	defer file.Close()

	model := &MLModel{}
	err = json.NewDecoder(file).Decode(model)
	if err != nil {
		return nil, fmt.Errorf("failed to decode the model file: %v", err)
	}

	return model, nil
}

// Predict uses the ML model to predict the product page pattern.
func (m *MLModel) Predict(url string) string {
	// Predict the pattern based on the URL
	for _, pattern := range m.Patterns {
		if pattern == url {
			return pattern
		}
	}
	return ""
}

// IsProductURL checks if the given URL is a product page of any e-commerce website.
func (m *MLModel) IsProductURL(url string) bool {
	for _, pattern := range m.Patterns {
		if matchPattern(pattern, url) {
			return true
		}
	}
	return false
}

// matchPattern checks if the given URL matches the pattern.
func matchPattern(pattern, url string) bool {
	return len(pattern) > 0 && len(url) > 0 && (pattern == "*" || contains(url, pattern))
}

// contains checks if the URL contains the pattern.
func contains(url, pattern string) bool {
	return len(url) >= len(pattern) && url[:len(pattern)] == pattern
}

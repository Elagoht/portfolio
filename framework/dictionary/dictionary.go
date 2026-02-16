// Package dictionary provides key-value translation lookup for the Statigo framework.
package dictionary

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"strings"
)

// Dictionary manages translations.
type Dictionary struct {
	translations map[string]interface{}
}

// New creates a new Dictionary instance by loading translations from the given filesystem.
func New(translationsFS fs.FS, _ string) (*Dictionary, error) {
	data, err := fs.ReadFile(translationsFS, "en.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read en.json: %w", err)
	}

	var translations map[string]interface{}
	if err := json.Unmarshal(data, &translations); err != nil {
		return nil, fmt.Errorf("failed to parse en.json: %w", err)
	}

	return &Dictionary{
		translations: translations,
	}, nil
}

// GetRaw retrieves raw structured data (arrays, objects) from translations using dot notation.
// Example: GetRaw("features.descriptions") returns []interface{}
func (d *Dictionary) GetRaw(_ string, key string) interface{} {
	parts := strings.Split(key, ".")
	var current interface{} = d.translations

	for _, part := range parts {
		if currentMap, ok := current.(map[string]interface{}); ok {
			current = currentMap[part]
		} else {
			return nil
		}
	}
	return current
}

// Get retrieves a string translation for the given key.
// Returns the key itself if translation is not found.
func (d *Dictionary) Get(_ string, key string) string {
	value := d.GetRaw("", key)
	if str, ok := value.(string); ok {
		return str
	}
	return key
}

package parser

import (
	"encoding/json"
	"fmt"
	"os"
)

// LoadJSON loads and parses a JSON file
func LoadJSON(filename string) (any, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	var data any
	err = json.Unmarshal(file, &data)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	return data, nil
}

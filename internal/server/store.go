package server

import (
	"encoding/json"
	"fmt"
	"os"
)

func SaveHistory(filename string, history map[string][]string) error {
	data, err := json.MarshalIndent(history, "", " ")
	if err != nil {
		return fmt.Errorf("encode history: %w", err)
	}

	err = os.WriteFile(filename, data, 0o600)
	if err != nil {
		return fmt.Errorf("write history: %w", err)
	}

	return nil
}

func LoadHistory(filename string) (map[string][]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("read history: %w", err)
	}

	var history map[string][]string

	err = json.Unmarshal(data, &history)
	if err != nil {
		return nil, fmt.Errorf("decode history: %w", err)
	}

	if history == nil {
		return nil, fmt.Errorf("history must be a JSON object")
	}

	return history, nil
}

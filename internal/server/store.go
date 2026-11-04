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

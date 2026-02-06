package game

import (
	"ctf-tool/pkg/data"
	"encoding/json"
	"fmt"
)

func LoadConfig() (*Config, error) {
	raw := data.LoadRawData()

	var config Config
	err := json.Unmarshal(raw, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse game data: %w", err)
	}

	return &config, nil
}

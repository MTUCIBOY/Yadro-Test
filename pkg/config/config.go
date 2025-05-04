package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Laps        int    `json:"laps"`
	LapLen      int    `json:"lapLen"`
	PenaltyLen  int    `json:"penaltyLen"`
	FiringLines int    `json:"firingLines"`
	Start       string `json:"start"`
	StartDelta  string `json:"startDelta"`
}

func MustLoad(configPath string) *Config {
	file, err := os.Open(configPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		panic(err)
	}

	return &config
}

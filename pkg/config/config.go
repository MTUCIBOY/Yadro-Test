package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Laps        uint   `json:"laps"`
	LapLen      uint   `json:"lapLen"`
	PenaltyLen  uint   `json:"penaltyLen"`
	FiringLines uint   `json:"firingLines"`
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
